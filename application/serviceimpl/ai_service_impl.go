package serviceimpl

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/datatypes"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/cache"
	"gofiber-template/infrastructure/external/google"
	"gofiber-template/infrastructure/external/openai"
	"gofiber-template/pkg/logger"
)

type AIServiceImpl struct {
	sessionRepo  repositories.AIChatSessionRepository
	messageRepo  repositories.AIChatMessageRepository
	historyRepo  repositories.SearchHistoryRepository
	aiClient     *openai.AIClient
	googleSearch *google.SearchClient
	redisClient  *redis.Client
}

func NewAIService(
	sessionRepo repositories.AIChatSessionRepository,
	messageRepo repositories.AIChatMessageRepository,
	historyRepo repositories.SearchHistoryRepository,
	aiClient *openai.AIClient,
	googleSearch *google.SearchClient,
	redisClient *redis.Client,
) services.AIService {
	return &AIServiceImpl{
		sessionRepo:  sessionRepo,
		messageRepo:  messageRepo,
		historyRepo:  historyRepo,
		aiClient:     aiClient,
		googleSearch: googleSearch,
		redisClient:  redisClient,
	}
}

func (s *AIServiceImpl) AISearch(ctx context.Context, userID uuid.UUID, req *dto.AISearchRequest) (*dto.AISearchResponse, error) {
	startTime := time.Now()

	logger.InfoContext(ctx, "AI Search started",
		"user_id", userID.String(),
		"query", req.Query,
		"lang", req.Language,
	)

	// Check cache first (AI responses are expensive!)
	cacheKey := cache.SearchAIKey(req.Query)
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedResult dto.AISearchResponse
		if json.Unmarshal([]byte(cached), &cachedResult) == nil {
			logger.InfoContext(ctx, "AI Search cache hit",
				"user_id", userID.String(),
				"query", req.Query,
				"response_time_ms", time.Since(startTime).Milliseconds(),
			)
			// Save search history even for cached results
			if userID != uuid.Nil {
				history := &models.SearchHistory{
					UserID:      userID,
					Query:       req.Query,
					SearchType:  models.SearchTypeAI,
					ResultCount: len(cachedResult.Sources),
				}
				_ = s.historyRepo.Create(ctx, history)
			}
			return &cachedResult, nil
		}
	}

	// Cache miss - Get search results from Google
	searchResponse, err := s.googleSearch.SearchAll(ctx, req.Query, 1, 5)
	if err != nil {
		logger.ErrorContext(ctx, "AI Search - Google Search failed",
			"user_id", userID.String(),
			"query", req.Query,
			"error", err.Error(),
			"response_time_ms", time.Since(startTime).Milliseconds(),
		)
		return nil, err
	}

	logger.InfoContext(ctx, "AI Search - Google Search completed",
		"user_id", userID.String(),
		"query", req.Query,
		"results_count", len(searchResponse.Items),
	)

	// Prepare sources and search context
	var sources []dto.MessageSource
	var searchContext []openai.SearchResultContext
	for _, r := range searchResponse.Items {
		sources = append(sources, dto.MessageSource{
			Title:   r.Title,
			URL:     r.Link,
			Snippet: r.Snippet,
		})
		searchContext = append(searchContext, openai.SearchResultContext{
			Title:   r.Title,
			URL:     r.Link,
			Snippet: r.Snippet,
		})
	}

	// Get language from request, default to Thai
	lang := req.Language
	if lang == "" {
		lang = "th"
	}

	// Generate AI summary
	aiResponse, err := s.aiClient.GenerateTravelSummary(ctx, req.Query, searchContext, lang)
	if err != nil {
		logger.ErrorContext(ctx, "AI Search - OpenAI failed",
			"user_id", userID.String(),
			"query", req.Query,
			"lang", lang,
			"error", err.Error(),
			"response_time_ms", time.Since(startTime).Milliseconds(),
		)
		return nil, err
	}

	// Extract summary from response
	summary := ""
	if len(aiResponse.Choices) > 0 {
		summary = aiResponse.Choices[0].Message.Content
	}

	response := &dto.AISearchResponse{
		Query:   req.Query,
		Summary: summary,
		Sources: sources,
	}

	// Store in cache (6 hours for AI responses)
	if jsonData, err := json.Marshal(response); err == nil {
		s.redisClient.Set(ctx, cacheKey, jsonData, cache.TTLSearchAI)
	}

	// Save search history
	if userID != uuid.Nil {
		history := &models.SearchHistory{
			UserID:      userID,
			Query:       req.Query,
			SearchType:  models.SearchTypeAI,
			ResultCount: len(sources),
		}
		_ = s.historyRepo.Create(ctx, history)
	}

	logger.InfoContext(ctx, "AI Search completed",
		"user_id", userID.String(),
		"query", req.Query,
		"sources_count", len(sources),
		"response_time_ms", time.Since(startTime).Milliseconds(),
	)

	return response, nil
}

func (s *AIServiceImpl) CreateChatSession(ctx context.Context, userID uuid.UUID, req *dto.CreateAIChatRequest) (*dto.AIChatSessionDetailResponse, error) {
	startTime := time.Now()

	logger.InfoContext(ctx, "CreateChatSession started",
		"user_id", userID.String(),
		"query", req.Query,
		"lang", req.Lang,
	)

	// Get language from request, default to Thai
	lang := req.Lang
	if lang == "" {
		lang = "th"
	}

	// Create session
	session := &models.AIChatSession{
		UserID:       userID,
		Title:        truncateString(req.Query, 100),
		InitialQuery: req.Query,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		logger.ErrorContext(ctx, "CreateChatSession - session creation failed",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	logger.InfoContext(ctx, "CreateChatSession - session created",
		"user_id", userID.String(),
		"session_id", session.ID.String(),
	)

	// Create user message
	userMessage := &models.AIChatMessage{
		SessionID: session.ID,
		Role:      models.MessageRoleUser,
		Content:   req.Query,
		CreatedAt: time.Now(),
	}

	if err := s.messageRepo.Create(ctx, userMessage); err != nil {
		logger.ErrorContext(ctx, "CreateChatSession - user message creation failed",
			"user_id", userID.String(),
			"session_id", session.ID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	// Get search results
	searchResponse, err := s.googleSearch.SearchAll(ctx, req.Query, 1, 5)
	if err != nil {
		logger.ErrorContext(ctx, "CreateChatSession - Google Search failed",
			"user_id", userID.String(),
			"session_id", session.ID.String(),
			"query", req.Query,
			"error", err.Error(),
		)
		return nil, err
	}

	// Prepare sources and search context
	var sources []models.MessageSource
	var searchContext []openai.SearchResultContext
	for _, r := range searchResponse.Items {
		sources = append(sources, models.MessageSource{
			Title:   r.Title,
			URL:     r.Link,
			Snippet: r.Snippet,
		})
		searchContext = append(searchContext, openai.SearchResultContext{
			Title:   r.Title,
			URL:     r.Link,
			Snippet: r.Snippet,
		})
	}

	// Generate AI response with language
	aiResponse, err := s.aiClient.GenerateTravelSummary(ctx, req.Query, searchContext, lang)
	if err != nil {
		logger.ErrorContext(ctx, "CreateChatSession - OpenAI failed",
			"user_id", userID.String(),
			"session_id", session.ID.String(),
			"query", req.Query,
			"lang", lang,
			"error", err.Error(),
			"response_time_ms", time.Since(startTime).Milliseconds(),
		)
		return nil, err
	}

	// Extract content from response
	responseContent := ""
	if len(aiResponse.Choices) > 0 {
		responseContent = aiResponse.Choices[0].Message.Content
	}

	// Create assistant message
	sourcesJSON, _ := json.Marshal(sources)
	assistantMessage := &models.AIChatMessage{
		SessionID: session.ID,
		Role:      models.MessageRoleAssistant,
		Content:   responseContent,
		Sources:   datatypes.JSON(sourcesJSON),
		CreatedAt: time.Now(),
	}

	if err := s.messageRepo.Create(ctx, assistantMessage); err != nil {
		logger.ErrorContext(ctx, "CreateChatSession - assistant message creation failed",
			"user_id", userID.String(),
			"session_id", session.ID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	// Update session
	session.UpdatedAt = time.Now()
	_ = s.sessionRepo.Update(ctx, session.ID, session)

	// Get session with messages
	session, err = s.sessionRepo.GetByIDWithMessages(ctx, session.ID)
	if err != nil {
		logger.ErrorContext(ctx, "CreateChatSession - get session with messages failed",
			"user_id", userID.String(),
			"session_id", session.ID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	logger.InfoContext(ctx, "CreateChatSession completed",
		"user_id", userID.String(),
		"session_id", session.ID.String(),
		"response_time_ms", time.Since(startTime).Milliseconds(),
	)

	return dto.AIChatSessionToDetailResponse(session), nil
}

func (s *AIServiceImpl) GetChatSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (*dto.AIChatSessionDetailResponse, error) {
	logger.InfoContext(ctx, "GetChatSession",
		"user_id", userID.String(),
		"session_id", sessionID.String(),
	)

	session, err := s.sessionRepo.GetByIDWithMessages(ctx, sessionID)
	if err != nil {
		logger.WarnContext(ctx, "GetChatSession - session not found",
			"user_id", userID.String(),
			"session_id", sessionID.String(),
			"error", err.Error(),
		)
		return nil, errors.New("session not found")
	}

	if session.UserID != userID {
		logger.WarnContext(ctx, "GetChatSession - unauthorized access",
			"user_id", userID.String(),
			"session_id", sessionID.String(),
			"session_owner", session.UserID.String(),
		)
		return nil, errors.New("unauthorized")
	}

	return dto.AIChatSessionToDetailResponse(session), nil
}

func (s *AIServiceImpl) GetChatSessions(ctx context.Context, userID uuid.UUID, req *dto.GetAIChatSessionsRequest) (*dto.AIChatSessionListResponse, error) {
	logger.InfoContext(ctx, "GetChatSessions",
		"user_id", userID.String(),
		"page", req.Page,
		"page_size", req.PageSize,
	)

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	sessions, err := s.sessionRepo.GetByUserID(ctx, userID, offset, req.PageSize)
	if err != nil {
		logger.ErrorContext(ctx, "GetChatSessions failed",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	total, err := s.sessionRepo.CountByUserID(ctx, userID)
	if err != nil {
		logger.ErrorContext(ctx, "GetChatSessions - count failed",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	var sessionResponses []dto.AIChatSessionResponse
	for _, session := range sessions {
		sessionResponses = append(sessionResponses, *dto.AIChatSessionToResponse(session))
	}

	return &dto.AIChatSessionListResponse{
		Sessions: sessionResponses,
		Meta: dto.PaginationMeta{
			Total:  total,
			Offset: offset,
			Limit:  req.PageSize,
		},
	}, nil
}

func (s *AIServiceImpl) DeleteChatSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	logger.InfoContext(ctx, "DeleteChatSession",
		"user_id", userID.String(),
		"session_id", sessionID.String(),
	)

	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		logger.WarnContext(ctx, "DeleteChatSession - session not found",
			"user_id", userID.String(),
			"session_id", sessionID.String(),
			"error", err.Error(),
		)
		return errors.New("session not found")
	}

	if session.UserID != userID {
		logger.WarnContext(ctx, "DeleteChatSession - unauthorized access",
			"user_id", userID.String(),
			"session_id", sessionID.String(),
			"session_owner", session.UserID.String(),
		)
		return errors.New("unauthorized")
	}

	if err := s.sessionRepo.Delete(ctx, sessionID); err != nil {
		logger.ErrorContext(ctx, "DeleteChatSession failed",
			"user_id", userID.String(),
			"session_id", sessionID.String(),
			"error", err.Error(),
		)
		return err
	}

	logger.InfoContext(ctx, "DeleteChatSession completed",
		"user_id", userID.String(),
		"session_id", sessionID.String(),
	)

	return nil
}

func (s *AIServiceImpl) ClearAllChatSessions(ctx context.Context, userID uuid.UUID) error {
	logger.InfoContext(ctx, "ClearAllChatSessions",
		"user_id", userID.String(),
	)

	if err := s.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
		logger.ErrorContext(ctx, "ClearAllChatSessions failed",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return err
	}

	logger.InfoContext(ctx, "ClearAllChatSessions completed",
		"user_id", userID.String(),
	)

	return nil
}

func (s *AIServiceImpl) SendMessage(ctx context.Context, userID uuid.UUID, req *dto.SendAIChatMessageRequest) (*dto.AIChatMessageResponse, error) {
	startTime := time.Now()

	logger.InfoContext(ctx, "SendMessage started",
		"user_id", userID.String(),
		"session_id", req.SessionID.String(),
		"message_length", len(req.Message),
		"lang", req.Lang,
	)

	// Get language from request, default to Thai
	lang := req.Lang
	if lang == "" {
		lang = "th"
	}

	session, err := s.sessionRepo.GetByID(ctx, req.SessionID)
	if err != nil {
		logger.WarnContext(ctx, "SendMessage - session not found",
			"user_id", userID.String(),
			"session_id", req.SessionID.String(),
			"error", err.Error(),
		)
		return nil, errors.New("session not found")
	}

	if session.UserID != userID {
		logger.WarnContext(ctx, "SendMessage - unauthorized access",
			"user_id", userID.String(),
			"session_id", req.SessionID.String(),
			"session_owner", session.UserID.String(),
		)
		return nil, errors.New("unauthorized")
	}

	// Create user message
	userMessage := &models.AIChatMessage{
		SessionID: req.SessionID,
		Role:      models.MessageRoleUser,
		Content:   req.Message,
		CreatedAt: time.Now(),
	}

	if err := s.messageRepo.Create(ctx, userMessage); err != nil {
		logger.ErrorContext(ctx, "SendMessage - user message creation failed",
			"user_id", userID.String(),
			"session_id", req.SessionID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	// Get recent messages for context
	recentMessages, err := s.messageRepo.GetRecentBySessionID(ctx, req.SessionID, 10)
	if err != nil {
		logger.ErrorContext(ctx, "SendMessage - get recent messages failed",
			"user_id", userID.String(),
			"session_id", req.SessionID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	logger.InfoContext(ctx, "SendMessage - context loaded",
		"user_id", userID.String(),
		"session_id", req.SessionID.String(),
		"context_messages", len(recentMessages),
	)

	// Build conversation history
	var chatHistory []openai.ChatMessage
	for _, msg := range recentMessages {
		chatHistory = append(chatHistory, openai.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Get additional search results if query seems like a new question
	var sources []models.MessageSource
	searchResponse, err := s.googleSearch.SearchAll(ctx, req.Message, 1, 3)
	if err == nil && searchResponse != nil && len(searchResponse.Items) > 0 {
		var searchContext string
		for _, r := range searchResponse.Items {
			sources = append(sources, models.MessageSource{
				Title:   r.Title,
				URL:     r.Link,
				Snippet: r.Snippet,
			})
			searchContext += "Title: " + r.Title + "\nContent: " + r.Snippet + "\n"
		}
		// Add search context to the message (with language-appropriate label)
		if lang == "en" {
			req.Message += "\n\nAdditional reference information:\n" + searchContext
		} else {
			req.Message += "\n\nข้อมูลอ้างอิงเพิ่มเติม:\n" + searchContext
		}
	}

	// Generate AI response with language
	aiResponse, err := s.aiClient.ContinueChat(ctx, chatHistory, req.Message, lang)
	if err != nil {
		logger.ErrorContext(ctx, "SendMessage - OpenAI failed",
			"user_id", userID.String(),
			"session_id", req.SessionID.String(),
			"lang", lang,
			"error", err.Error(),
			"response_time_ms", time.Since(startTime).Milliseconds(),
		)
		return nil, err
	}

	// Extract content from response
	responseContent := ""
	if len(aiResponse.Choices) > 0 {
		responseContent = aiResponse.Choices[0].Message.Content
	}

	// Create assistant message
	sourcesJSON, _ := json.Marshal(sources)
	assistantMessage := &models.AIChatMessage{
		SessionID: req.SessionID,
		Role:      models.MessageRoleAssistant,
		Content:   responseContent,
		Sources:   datatypes.JSON(sourcesJSON),
		CreatedAt: time.Now(),
	}

	if err := s.messageRepo.Create(ctx, assistantMessage); err != nil {
		logger.ErrorContext(ctx, "SendMessage - assistant message creation failed",
			"user_id", userID.String(),
			"session_id", req.SessionID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	// Update session
	session.UpdatedAt = time.Now()
	_ = s.sessionRepo.Update(ctx, session.ID, session)

	logger.InfoContext(ctx, "SendMessage completed",
		"user_id", userID.String(),
		"session_id", req.SessionID.String(),
		"response_length", len(responseContent),
		"response_time_ms", time.Since(startTime).Milliseconds(),
	)

	return dto.AIChatMessageToResponse(assistantMessage), nil
}

func (s *AIServiceImpl) SendMessageStream(ctx context.Context, userID uuid.UUID, req *dto.SendAIChatMessageRequest, writer io.Writer) error {
	// For streaming, we would need to implement SSE or WebSocket
	// For now, fallback to regular response
	response, err := s.SendMessage(ctx, userID, req)
	if err != nil {
		return err
	}

	// Write response as JSON
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func (s *AIServiceImpl) AISearchStream(ctx context.Context, userID uuid.UUID, req *dto.AISearchRequest, writer io.Writer) error {
	// For streaming, we would need to implement SSE or WebSocket
	// For now, fallback to regular response
	response, err := s.AISearch(ctx, userID, req)
	if err != nil {
		return err
	}

	// Write response as JSON
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
