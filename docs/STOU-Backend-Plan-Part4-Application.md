# STOU Smart Tour - Backend Development Plan
# Part 4: Application Layer (Services Implementation)

---

## Table of Contents - All Parts

| Part | หัวข้อ | สถานะ |
|------|--------|-------|
| Part 1 | Project Overview & Foundation | ✅ Done |
| Part 2 | Domain Layer (Models, DTOs, Interfaces) | ✅ Done |
| Part 3 | Infrastructure Layer (External APIs, Cache) | ✅ Done |
| **Part 4** | **Application Layer (Services Implementation)** | 📍 Current |
| Part 5 | Interface Layer (Handlers, Routes, Middleware) | ⏳ Pending |

---

## 1. Service Implementations Overview

### 1.1 โครงสร้าง Application Layer

```
application/
└── serviceimpl/
    ├── user_service_impl.go       # ✅ มีอยู่แล้ว (ปรับเพิ่ม student_id)
    ├── task_service_impl.go       # ✅ มีอยู่แล้ว
    ├── file_service_impl.go       # ✅ มีอยู่แล้ว
    ├── job_service_impl.go        # ✅ มีอยู่แล้ว
    ├── search_service_impl.go     # 🆕 NEW
    ├── ai_service_impl.go         # 🆕 NEW
    ├── folder_service_impl.go     # 🆕 NEW
    ├── favorite_service_impl.go   # 🆕 NEW
    ├── translate_service_impl.go  # 🆕 NEW
    └── qrcode_service_impl.go     # 🆕 NEW
```

---

## 2. Search Service Implementation

```go
// application/serviceimpl/search_service_impl.go

package serviceimpl

import (
    "context"
    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/dto"
    "github.com/your-org/stou-smart-tour/domain/models"
    "github.com/your-org/stou-smart-tour/domain/repositories"
    "github.com/your-org/stou-smart-tour/domain/services"
    "github.com/your-org/stou-smart-tour/infrastructure/cache"
    "github.com/your-org/stou-smart-tour/infrastructure/external/google"
    "github.com/your-org/stou-smart-tour/infrastructure/redis"
)

type searchService struct {
    searchClient   *google.SearchClient
    placesClient   *google.PlacesClient
    historyRepo    repositories.SearchHistoryRepository
    redis          *redis.RedisClient
}

func NewSearchService(
    searchClient *google.SearchClient,
    placesClient *google.PlacesClient,
    historyRepo repositories.SearchHistoryRepository,
    redis *redis.RedisClient,
) services.SearchService {
    return &searchService{
        searchClient: searchClient,
        placesClient: placesClient,
        historyRepo:  historyRepo,
        redis:        redis,
    }
}

// Search performs Google Custom Search
func (s *searchService) Search(ctx context.Context, req *dto.SearchRequest) (*dto.SearchResponse, error) {
    // Set defaults
    if req.Page <= 0 {
        req.Page = 1
    }
    if req.PerPage <= 0 {
        req.PerPage = 10
    }
    if req.PerPage > 10 {
        req.PerPage = 10 // Google API limit
    }
    if req.Type == "" {
        req.Type = "all"
    }

    // Check cache first
    cacheKey := cache.SearchKey(req.Query, req.Type, req.Page)
    var cachedResult dto.SearchResponse
    if err := s.redis.Get(ctx, cacheKey, &cachedResult); err == nil {
        return &cachedResult, nil
    }

    // Perform search based on type
    var searchResults *google.SearchResponse
    var err error

    switch req.Type {
    case "image":
        searchResults, err = s.searchClient.SearchImages(ctx, req.Query, req.Page, req.PerPage)
    case "video":
        // For video, we modify query to search YouTube-related content
        searchResults, err = s.searchClient.SearchAll(ctx, req.Query+" site:youtube.com", req.Page, req.PerPage)
    default: // "all", "website"
        searchResults, err = s.searchClient.SearchAll(ctx, req.Query, req.Page, req.PerPage)
    }

    if err != nil {
        return nil, fmt.Errorf("search failed: %w", err)
    }

    // Convert to response DTO
    response := s.convertSearchResults(searchResults, req)

    // Cache the result
    _ = s.redis.Set(ctx, cacheKey, response, cache.TTLSearch)

    return response, nil
}

// convertSearchResults converts Google API response to our DTO
func (s *searchService) convertSearchResults(results *google.SearchResponse, req *dto.SearchRequest) *dto.SearchResponse {
    items := make([]dto.SearchResult, 0, len(results.Items))

    for _, item := range results.Items {
        result := dto.SearchResult{
            ID:      generateResultID(item.Link),
            Type:    s.determineResultType(item, req.Type),
            Title:   item.Title,
            URL:     item.Link,
            Snippet: item.Snippet,
            Source:  item.DisplayLink,
        }

        // Get thumbnail
        if len(item.Pagemap.CSEThumbnail) > 0 {
            result.ThumbnailURL = item.Pagemap.CSEThumbnail[0].Src
        } else if len(item.Pagemap.CSEImage) > 0 {
            result.ThumbnailURL = item.Pagemap.CSEImage[0].Src
        }

        // For image search, include image info
        if item.Image != nil {
            result.Image = &dto.ImageInfo{
                Width:  item.Image.Width,
                Height: item.Image.Height,
            }
            result.ThumbnailURL = item.Image.ThumbnailLink
        }

        // For video results (YouTube)
        if strings.Contains(item.Link, "youtube.com") || strings.Contains(item.Link, "youtu.be") {
            result.Type = "video"
            result.Video = s.extractVideoInfo(item)
        }

        items = append(items, result)
    }

    // Parse total results
    total, _ := strconv.ParseInt(results.SearchInformation.TotalResults, 10, 64)

    return &dto.SearchResponse{
        Results: items,
        Meta: dto.SearchMeta{
            Query:      req.Query,
            SearchType: req.Type,
            Page:       req.Page,
            PerPage:    req.PerPage,
            Total:      total,
            TotalPages: dto.CalculateTotalPages(total, req.PerPage),
        },
    }
}

// determineResultType determines the type of search result
func (s *searchService) determineResultType(item google.SearchItem, requestType string) string {
    if requestType == "image" {
        return "image"
    }

    url := strings.ToLower(item.Link)
    if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
        return "video"
    }

    return "website"
}

// extractVideoInfo extracts YouTube video information from search result
func (s *searchService) extractVideoInfo(item google.SearchItem) *dto.VideoInfo {
    info := &dto.VideoInfo{}

    // Extract video ID from URL
    if strings.Contains(item.Link, "youtube.com/watch") {
        parts := strings.Split(item.Link, "v=")
        if len(parts) > 1 {
            videoID := strings.Split(parts[1], "&")[0]
            info.VideoID = videoID
        }
    } else if strings.Contains(item.Link, "youtu.be/") {
        parts := strings.Split(item.Link, "youtu.be/")
        if len(parts) > 1 {
            info.VideoID = strings.Split(parts[1], "?")[0]
        }
    }

    // Extract channel from metatags
    if len(item.Pagemap.Metatags) > 0 {
        meta := item.Pagemap.Metatags[0]
        if channel, ok := meta["og:site_name"]; ok {
            info.Channel = channel
        }
    }

    return info
}

// SearchPlaces searches for nearby places using Google Places API
func (s *searchService) SearchPlaces(ctx context.Context, req *dto.PlacesSearchRequest) (*dto.PlacesResponse, error) {
    // Set defaults
    if req.Radius <= 0 {
        req.Radius = 5000 // 5km default
    }

    // Check cache
    cacheKey := cache.NearbyPlacesKey(req.Lat, req.Lng, req.Radius, req.Type, req.Keyword)
    var cachedResult dto.PlacesResponse
    if err := s.redis.Get(ctx, cacheKey, &cachedResult); err == nil {
        return &cachedResult, nil
    }

    // Search places
    searchReq := &google.NearbySearchRequest{
        Lat:      req.Lat,
        Lng:      req.Lng,
        Radius:   req.Radius,
        Type:     req.Type,
        Keyword:  req.Keyword,
        Language: "th",
    }

    results, err := s.placesClient.NearbySearch(ctx, searchReq)
    if err != nil {
        return nil, fmt.Errorf("places search failed: %w", err)
    }

    // Convert to response
    places := make([]dto.PlaceResult, 0, len(results.Results))
    for _, place := range results.Results {
        placeResult := dto.PlaceResult{
            PlaceID:     place.PlaceID,
            Name:        place.Name,
            Address:     place.Vicinity,
            Lat:         place.Geometry.Location.Lat,
            Lng:         place.Geometry.Location.Lng,
            Rating:      place.Rating,
            ReviewCount: place.UserRatingsTotal,
            Types:       place.Types,
            Distance:    google.CalculateDistance(req.Lat, req.Lng, place.Geometry.Location.Lat, place.Geometry.Location.Lng),
        }

        // Get photo URL
        if len(place.Photos) > 0 {
            placeResult.PhotoURL = s.placesClient.GetPhotoURL(place.Photos[0].PhotoReference, 400)
        }

        // Set open status
        if place.OpeningHours != nil {
            isOpen := place.OpeningHours.OpenNow
            placeResult.IsOpen = &isOpen
        }

        places = append(places, placeResult)
    }

    response := &dto.PlacesResponse{Places: places}

    // Cache result
    _ = s.redis.Set(ctx, cacheKey, response, cache.TTLNearbyPlaces)

    return response, nil
}

// GetPlaceDetail gets detailed information about a place
func (s *searchService) GetPlaceDetail(ctx context.Context, placeID string) (*dto.PlaceDetailResponse, error) {
    // Check cache
    cacheKey := cache.PlaceDetailsKey(placeID)
    var cachedResult dto.PlaceDetailResponse
    if err := s.redis.Get(ctx, cacheKey, &cachedResult); err == nil {
        return &cachedResult, nil
    }

    // Get place details
    result, err := s.placesClient.GetPlaceDetails(ctx, &google.PlaceDetailsRequest{
        PlaceID:  placeID,
        Language: "th",
    })
    if err != nil {
        return nil, fmt.Errorf("get place details failed: %w", err)
    }

    place := result.Result

    // Convert photos
    photos := make([]dto.PhotoInfo, 0, len(place.Photos))
    for _, photo := range place.Photos {
        photos = append(photos, dto.PhotoInfo{
            URL: s.placesClient.GetPhotoURL(photo.PhotoReference, 800),
        })
    }

    // Convert reviews
    reviews := make([]dto.ReviewInfo, 0, len(place.Reviews))
    for _, review := range place.Reviews {
        reviews = append(reviews, dto.ReviewInfo{
            Author: review.AuthorName,
            Rating: review.Rating,
            Text:   review.Text,
            Time:   time.Unix(review.Time, 0),
        })
    }

    // Build response
    response := &dto.PlaceDetailResponse{
        PlaceID:      place.PlaceID,
        Name:         place.Name,
        Address:      place.FormattedAddress,
        Phone:        place.FormattedPhoneNumber,
        Website:      place.Website,
        Rating:       place.Rating,
        ReviewCount:  place.UserRatingsTotal,
        PriceLevel:   place.PriceLevel,
        Photos:       photos,
        Reviews:      reviews,
        Lat:          place.Geometry.Location.Lat,
        Lng:          place.Geometry.Location.Lng,
        Types:        place.Types,
    }

    // Opening hours
    if place.OpeningHours != nil {
        response.OpeningHours = &dto.OpeningHours{
            WeekdayText: place.OpeningHours.WeekdayText,
            IsOpenNow:   place.OpeningHours.OpenNow,
        }
    }

    // Cache result
    _ = s.redis.Set(ctx, cacheKey, response, cache.TTLPlaceDetails)

    return response, nil
}

// SaveSearchHistory saves a search query to history
func (s *searchService) SaveSearchHistory(ctx context.Context, userID uuid.UUID, query, searchType string, resultCount int) error {
    history := &models.SearchHistory{
        UserID:      userID,
        Query:       query,
        SearchType:  searchType,
        ResultCount: resultCount,
    }
    return s.historyRepo.Create(ctx, history)
}

// GetSearchHistory gets user's search history
func (s *searchService) GetSearchHistory(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.SearchHistoryResponse, error) {
    if limit <= 0 {
        limit = dto.DefaultPerPage
    }

    histories, err := s.historyRepo.GetByUserID(ctx, userID, offset, limit)
    if err != nil {
        return nil, err
    }

    total, err := s.historyRepo.CountByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }

    return &dto.SearchHistoryResponse{
        History: dto.SearchHistoriesToSearchHistoryItems(histories),
        Meta: dto.PaginationMeta{
            Total:      total,
            Page:       (offset / limit) + 1,
            PerPage:    limit,
            TotalPages: dto.CalculateTotalPages(total, limit),
        },
    }, nil
}

// ClearSearchHistory clears user's search history
func (s *searchService) ClearSearchHistory(ctx context.Context, userID uuid.UUID) error {
    return s.historyRepo.DeleteByUserID(ctx, userID)
}

// Helper function to generate result ID
func generateResultID(url string) string {
    hash := md5.Sum([]byte(url))
    return hex.EncodeToString(hash[:8])
}
```

---

## 3. AI Service Implementation

```go
// application/serviceimpl/ai_service_impl.go

package serviceimpl

import (
    "context"
    "encoding/json"
    "fmt"
    "strings"

    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/dto"
    "github.com/your-org/stou-smart-tour/domain/models"
    "github.com/your-org/stou-smart-tour/domain/repositories"
    "github.com/your-org/stou-smart-tour/domain/services"
    "github.com/your-org/stou-smart-tour/infrastructure/cache"
    "github.com/your-org/stou-smart-tour/infrastructure/external/google"
    "github.com/your-org/stou-smart-tour/infrastructure/external/openai"
    "github.com/your-org/stou-smart-tour/infrastructure/redis"
)

type aiService struct {
    aiClient      *openai.AIClient
    searchClient  *google.SearchClient
    youtubeClient *google.YouTubeClient
    chatRepo      repositories.AIChatRepository
    redis         *redis.RedisClient
}

func NewAIService(
    aiClient *openai.AIClient,
    searchClient *google.SearchClient,
    youtubeClient *google.YouTubeClient,
    chatRepo repositories.AIChatRepository,
    redis *redis.RedisClient,
) services.AIService {
    return &aiService{
        aiClient:      aiClient,
        searchClient:  searchClient,
        youtubeClient: youtubeClient,
        chatRepo:      chatRepo,
        redis:         redis,
    }
}

// AISearch performs AI-enhanced search with summary generation
func (s *aiService) AISearch(ctx context.Context, query string, userID *uuid.UUID) (*dto.AISearchResponse, error) {
    // Check cache
    cacheKey := cache.SearchAIKey(query)
    var cachedResult dto.AISearchResponse
    if err := s.redis.Get(ctx, cacheKey, &cachedResult); err == nil {
        return &cachedResult, nil
    }

    // Step 1: Search Google for relevant content
    searchResults, err := s.searchClient.SearchAll(ctx, query, 1, 10)
    if err != nil {
        return nil, fmt.Errorf("search failed: %w", err)
    }

    // Step 2: Prepare context from search results
    searchContext := make([]openai.SearchResultContext, 0, len(searchResults.Items))
    sources := make([]dto.SourceInfo, 0, len(searchResults.Items))

    for _, item := range searchResults.Items {
        searchContext = append(searchContext, openai.SearchResultContext{
            Title:   item.Title,
            Snippet: item.Snippet,
            URL:     item.Link,
        })
        sources = append(sources, dto.SourceInfo{
            Title:   item.Title,
            URL:     item.Link,
            Snippet: item.Snippet,
        })
    }

    // Step 3: Generate AI summary
    aiResponse, err := s.aiClient.GenerateTravelSummary(ctx, query, searchContext)
    if err != nil {
        return nil, fmt.Errorf("AI generation failed: %w", err)
    }

    // Extract AI content
    aiContent := ""
    if len(aiResponse.Choices) > 0 {
        aiContent = aiResponse.Choices[0].Message.Content
    }

    // Step 4: Get related videos
    videos, err := s.GetRelatedVideos(ctx, query, 5)
    if err != nil {
        // Don't fail, just log
        videos = []dto.VideoResult{}
    }

    // Step 5: Generate follow-up suggestions
    suggestions := s.generateFollowUpSuggestions(query)

    // Step 6: Create session if user is logged in
    var sessionID string
    if userID != nil {
        session := &models.AIChatSession{
            UserID:       *userID,
            Title:        truncateString(query, 100),
            InitialQuery: query,
        }
        if err := s.chatRepo.CreateSession(ctx, session); err == nil {
            sessionID = session.ID.String()

            // Save initial messages
            s.saveChatMessages(ctx, session.ID, query, aiContent, sources)
        }
    } else {
        sessionID = uuid.New().String() // Temporary session ID for anonymous users
    }

    // Build response
    response := &dto.AISearchResponse{
        Summary: dto.AISummary{
            Content: aiContent,
            Sources: sources,
        },
        RelatedVideos:       videos,
        SessionID:           sessionID,
        FollowUpSuggestions: suggestions,
    }

    // Cache result
    _ = s.redis.Set(ctx, cacheKey, response, cache.TTLSearchAI)

    return response, nil
}

// Chat continues an AI conversation
func (s *aiService) Chat(ctx context.Context, req *dto.AIChatRequest, userID uuid.UUID) (*dto.AIChatResponse, error) {
    // Parse session ID
    sessionID, err := uuid.Parse(req.SessionID)
    if err != nil {
        return nil, fmt.Errorf("invalid session ID")
    }

    // Get session
    session, err := s.chatRepo.GetSessionByID(ctx, sessionID)
    if err != nil {
        return nil, fmt.Errorf("session not found")
    }

    // Verify ownership
    if session.UserID != userID {
        return nil, fmt.Errorf("unauthorized access to session")
    }

    // Build chat history
    history := make([]openai.ChatMessage, 0, len(session.Messages))
    for _, msg := range session.Messages {
        history = append(history, openai.ChatMessage{
            Role:    msg.Role,
            Content: msg.Content,
        })
    }

    // Generate AI response
    aiResponse, err := s.aiClient.ContinueChat(ctx, history, req.Message)
    if err != nil {
        return nil, fmt.Errorf("AI chat failed: %w", err)
    }

    // Extract response content
    responseContent := ""
    if len(aiResponse.Choices) > 0 {
        responseContent = aiResponse.Choices[0].Message.Content
    }

    // Save messages to database
    userMessage := &models.AIChatMessage{
        SessionID: sessionID,
        Role:      models.MessageRoleUser,
        Content:   req.Message,
    }
    s.chatRepo.CreateMessage(ctx, userMessage)

    assistantMessage := &models.AIChatMessage{
        SessionID: sessionID,
        Role:      models.MessageRoleAssistant,
        Content:   responseContent,
    }
    s.chatRepo.CreateMessage(ctx, assistantMessage)

    // Update session timestamp
    session.UpdatedAt = time.Now()
    s.chatRepo.UpdateSession(ctx, session)

    return &dto.AIChatResponse{
        Message: dto.AIChatMessageResponse{
            Role:    models.MessageRoleAssistant,
            Content: responseContent,
        },
        SessionID: sessionID.String(),
    }, nil
}

// GetRelatedVideos gets YouTube videos related to query
func (s *aiService) GetRelatedVideos(ctx context.Context, query string, limit int) ([]dto.VideoResult, error) {
    if limit <= 0 {
        limit = 5
    }

    // Check cache
    cacheKey := cache.YouTubeKey(query, limit)
    var cachedResult []dto.VideoResult
    if err := s.redis.Get(ctx, cacheKey, &cachedResult); err == nil {
        return cachedResult, nil
    }

    // Search YouTube
    searchResult, err := s.youtubeClient.SearchVideos(ctx, &google.VideoSearchRequest{
        Query:      query + " ท่องเที่ยว",
        MaxResults: limit,
        Order:      "relevance",
        RegionCode: "TH",
    })
    if err != nil {
        return nil, err
    }

    // Get video IDs for additional details
    videoIDs := make([]string, 0, len(searchResult.Items))
    for _, item := range searchResult.Items {
        videoIDs = append(videoIDs, item.ID.VideoID)
    }

    // Get video details (duration, views)
    var detailsMap = make(map[string]google.VideoDetails)
    if len(videoIDs) > 0 {
        details, err := s.youtubeClient.GetVideoDetails(ctx, videoIDs)
        if err == nil {
            for _, d := range details.Items {
                detailsMap[d.ID] = d
            }
        }
    }

    // Build response
    videos := make([]dto.VideoResult, 0, len(searchResult.Items))
    for _, item := range searchResult.Items {
        video := dto.VideoResult{
            ID:        item.ID.VideoID,
            Title:     item.Snippet.Title,
            Thumbnail: item.Snippet.Thumbnails.High.URL,
            Channel:   item.Snippet.ChannelTitle,
            URL:       google.GetVideoURL(item.ID.VideoID),
        }

        // Add details if available
        if details, ok := detailsMap[item.ID.VideoID]; ok {
            video.Duration = google.ParseDuration(details.ContentDetails.Duration)
            if viewCount, err := strconv.ParseInt(details.Statistics.ViewCount, 10, 64); err == nil {
                video.ViewCount = viewCount
            }
        }

        videos = append(videos, video)
    }

    // Cache result
    _ = s.redis.Set(ctx, cacheKey, videos, cache.TTLYouTube)

    return videos, nil
}

// GetChatSession gets a chat session by ID
func (s *aiService) GetChatSession(ctx context.Context, sessionID uuid.UUID) (*models.AIChatSession, error) {
    return s.chatRepo.GetSessionByID(ctx, sessionID)
}

// GetChatHistory gets user's chat history
func (s *aiService) GetChatHistory(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.AIChatSession, error) {
    return s.chatRepo.GetSessionsByUserID(ctx, userID, offset, limit)
}

// DeleteChatSession deletes a chat session
func (s *aiService) DeleteChatSession(ctx context.Context, sessionID, userID uuid.UUID) error {
    session, err := s.chatRepo.GetSessionByID(ctx, sessionID)
    if err != nil {
        return err
    }
    if session.UserID != userID {
        return fmt.Errorf("unauthorized")
    }
    return s.chatRepo.DeleteSession(ctx, sessionID)
}

// Helper: save chat messages
func (s *aiService) saveChatMessages(ctx context.Context, sessionID uuid.UUID, userQuery, aiResponse string, sources []dto.SourceInfo) {
    // Save user message
    userMessage := &models.AIChatMessage{
        SessionID: sessionID,
        Role:      models.MessageRoleUser,
        Content:   userQuery,
    }
    s.chatRepo.CreateMessage(ctx, userMessage)

    // Save assistant message with sources
    sourcesJSON, _ := json.Marshal(sources)
    assistantMessage := &models.AIChatMessage{
        SessionID: sessionID,
        Role:      models.MessageRoleAssistant,
        Content:   aiResponse,
        Sources:   sourcesJSON,
    }
    s.chatRepo.CreateMessage(ctx, assistantMessage)
}

// Helper: generate follow-up suggestions
func (s *aiService) generateFollowUpSuggestions(query string) []string {
    // Simple rule-based suggestions
    suggestions := []string{}

    if strings.Contains(query, "เชียงใหม่") || strings.Contains(query, "chiang mai") {
        suggestions = append(suggestions,
            "ที่พักแนะนำในเชียงใหม่?",
            "อาหารเหนือที่ต้องลอง?",
            "วัดสำคัญในเชียงใหม่?",
        )
    } else if strings.Contains(query, "กรุงเทพ") || strings.Contains(query, "bangkok") {
        suggestions = append(suggestions,
            "ที่เที่ยวกลางคืนกรุงเทพ?",
            "ร้านอาหารแนะนำ?",
            "การเดินทางในกรุงเทพ?",
        )
    } else {
        // Generic suggestions
        suggestions = append(suggestions,
            "ค่าใช้จ่ายโดยประมาณ?",
            "ที่พักแนะนำ?",
            "การเดินทางไปที่นั่น?",
        )
    }

    return suggestions
}

// Helper: truncate string
func truncateString(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen-3] + "..."
}
```

---

## 4. Folder Service Implementation

```go
// application/serviceimpl/folder_service_impl.go

package serviceimpl

import (
    "context"
    "encoding/base64"
    "fmt"

    "github.com/google/uuid"
    "github.com/skip2/go-qrcode"
    "github.com/your-org/stou-smart-tour/domain/dto"
    "github.com/your-org/stou-smart-tour/domain/models"
    "github.com/your-org/stou-smart-tour/domain/repositories"
    "github.com/your-org/stou-smart-tour/domain/services"
    "github.com/your-org/stou-smart-tour/pkg/config"
    "gorm.io/gorm"
)

type folderService struct {
    folderRepo     repositories.FolderRepository
    folderItemRepo repositories.FolderItemRepository
    config         *config.Config
}

func NewFolderService(
    folderRepo repositories.FolderRepository,
    folderItemRepo repositories.FolderItemRepository,
    config *config.Config,
) services.FolderService {
    return &folderService{
        folderRepo:     folderRepo,
        folderItemRepo: folderItemRepo,
        config:         config,
    }
}

// CreateFolder creates a new folder
func (s *folderService) CreateFolder(ctx context.Context, req *dto.CreateFolderRequest, userID uuid.UUID) (*dto.FolderResponse, error) {
    folder := dto.CreateFolderRequestToFolder(req, userID)

    if err := s.folderRepo.Create(ctx, folder); err != nil {
        return nil, fmt.Errorf("create folder failed: %w", err)
    }

    response := dto.FolderToFolderResponse(folder)
    return &response, nil
}

// GetFolder gets a folder with its items
func (s *folderService) GetFolder(ctx context.Context, folderID, userID uuid.UUID) (*dto.FolderDetailResponse, error) {
    folder, err := s.folderRepo.GetByID(ctx, folderID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("folder not found")
        }
        return nil, err
    }

    // Check access permission
    if folder.UserID != userID && !folder.IsPublic {
        return nil, fmt.Errorf("access denied")
    }

    // Get items
    items, err := s.folderItemRepo.GetByFolderID(ctx, folderID, 0, 100)
    if err != nil {
        return nil, err
    }

    total, err := s.folderItemRepo.CountByFolderID(ctx, folderID)
    if err != nil {
        return nil, err
    }

    return &dto.FolderDetailResponse{
        Folder: dto.FolderToFolderResponse(folder),
        Items:  dto.FolderItemsToFolderItemResponses(items),
        Meta: dto.PaginationMeta{
            Total:      total,
            Page:       1,
            PerPage:    100,
            TotalPages: dto.CalculateTotalPages(total, 100),
        },
    }, nil
}

// UpdateFolder updates a folder
func (s *folderService) UpdateFolder(ctx context.Context, folderID uuid.UUID, req *dto.UpdateFolderRequest, userID uuid.UUID) (*dto.FolderResponse, error) {
    folder, err := s.folderRepo.GetByID(ctx, folderID)
    if err != nil {
        return nil, fmt.Errorf("folder not found")
    }

    // Check ownership
    if folder.UserID != userID {
        return nil, fmt.Errorf("access denied")
    }

    // Update fields
    if req.Name != "" {
        folder.Name = req.Name
    }
    if req.Description != "" {
        folder.Description = req.Description
    }
    if req.IsPublic != nil {
        folder.IsPublic = *req.IsPublic
    }

    if err := s.folderRepo.Update(ctx, folder); err != nil {
        return nil, fmt.Errorf("update failed: %w", err)
    }

    response := dto.FolderToFolderResponse(folder)
    return &response, nil
}

// DeleteFolder deletes a folder and all its items
func (s *folderService) DeleteFolder(ctx context.Context, folderID, userID uuid.UUID) error {
    folder, err := s.folderRepo.GetByID(ctx, folderID)
    if err != nil {
        return fmt.Errorf("folder not found")
    }

    // Check ownership
    if folder.UserID != userID {
        return fmt.Errorf("access denied")
    }

    // Delete all items first
    if err := s.folderItemRepo.DeleteByFolderID(ctx, folderID); err != nil {
        return fmt.Errorf("delete items failed: %w", err)
    }

    // Delete folder
    if err := s.folderRepo.Delete(ctx, folderID); err != nil {
        return fmt.Errorf("delete folder failed: %w", err)
    }

    return nil
}

// GetUserFolders gets all folders for a user
func (s *folderService) GetUserFolders(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.FolderListResponse, error) {
    if limit <= 0 {
        limit = dto.DefaultPerPage
    }

    folders, err := s.folderRepo.GetByUserID(ctx, userID, offset, limit)
    if err != nil {
        return nil, err
    }

    total, err := s.folderRepo.CountByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }

    return &dto.FolderListResponse{
        Folders: dto.FoldersToFolderResponses(folders),
        Meta: dto.PaginationMeta{
            Total:      total,
            Page:       (offset / limit) + 1,
            PerPage:    limit,
            TotalPages: dto.CalculateTotalPages(total, limit),
        },
    }, nil
}

// AddItem adds an item to a folder
func (s *folderService) AddItem(ctx context.Context, folderID uuid.UUID, req *dto.AddFolderItemRequest, userID uuid.UUID) (*dto.FolderItemResponse, error) {
    // Check folder ownership
    folder, err := s.folderRepo.GetByID(ctx, folderID)
    if err != nil {
        return nil, fmt.Errorf("folder not found")
    }
    if folder.UserID != userID {
        return nil, fmt.Errorf("access denied")
    }

    // Create item
    item, err := dto.AddFolderItemRequestToFolderItem(req, folderID)
    if err != nil {
        return nil, fmt.Errorf("invalid item data: %w", err)
    }

    if err := s.folderItemRepo.Create(ctx, item); err != nil {
        return nil, fmt.Errorf("create item failed: %w", err)
    }

    // Increment folder item count
    s.folderRepo.IncrementItemCount(ctx, folderID)

    // Update cover image if this is the first image item
    if item.ThumbnailURL != "" && folder.CoverImageURL == "" {
        s.folderRepo.UpdateCoverImage(ctx, folderID, item.ThumbnailURL)
    }

    response := dto.FolderItemToFolderItemResponse(item)
    return &response, nil
}

// RemoveItem removes an item from a folder
func (s *folderService) RemoveItem(ctx context.Context, folderID, itemID, userID uuid.UUID) error {
    // Check folder ownership
    folder, err := s.folderRepo.GetByID(ctx, folderID)
    if err != nil {
        return fmt.Errorf("folder not found")
    }
    if folder.UserID != userID {
        return fmt.Errorf("access denied")
    }

    // Check item belongs to folder
    item, err := s.folderItemRepo.GetByID(ctx, itemID)
    if err != nil {
        return fmt.Errorf("item not found")
    }
    if item.FolderID != folderID {
        return fmt.Errorf("item does not belong to this folder")
    }

    // Delete item
    if err := s.folderItemRepo.Delete(ctx, itemID); err != nil {
        return fmt.Errorf("delete item failed: %w", err)
    }

    // Decrement folder item count
    s.folderRepo.DecrementItemCount(ctx, folderID)

    return nil
}

// GetFolderItems gets paginated items from a folder
func (s *folderService) GetFolderItems(ctx context.Context, folderID uuid.UUID, offset, limit int) ([]dto.FolderItemResponse, int64, error) {
    if limit <= 0 {
        limit = dto.DefaultPerPage
    }

    items, err := s.folderItemRepo.GetByFolderID(ctx, folderID, offset, limit)
    if err != nil {
        return nil, 0, err
    }

    total, err := s.folderItemRepo.CountByFolderID(ctx, folderID)
    if err != nil {
        return nil, 0, err
    }

    return dto.FolderItemsToFolderItemResponses(items), total, nil
}

// GenerateShareLink generates a share link and QR code for a folder
func (s *folderService) GenerateShareLink(ctx context.Context, folderID, userID uuid.UUID) (*dto.FolderShareResponse, error) {
    folder, err := s.folderRepo.GetByID(ctx, folderID)
    if err != nil {
        return nil, fmt.Errorf("folder not found")
    }
    if folder.UserID != userID {
        return nil, fmt.Errorf("access denied")
    }

    // Make folder public if not already
    if !folder.IsPublic {
        folder.IsPublic = true
        s.folderRepo.Update(ctx, folder)
    }

    // Generate share URL
    shareURL := fmt.Sprintf("%s/shared/folder/%s", s.config.AppURL, folderID.String())

    // Generate QR code
    qrCode, err := qrcode.Encode(shareURL, qrcode.Medium, 256)
    if err != nil {
        return nil, fmt.Errorf("generate QR code failed: %w", err)
    }

    return &dto.FolderShareResponse{
        ShareURL: shareURL,
        QRCode:   base64.StdEncoding.EncodeToString(qrCode),
    }, nil
}

// GetPublicFolder gets a public folder without authentication
func (s *folderService) GetPublicFolder(ctx context.Context, folderID uuid.UUID) (*dto.FolderDetailResponse, error) {
    folder, err := s.folderRepo.GetByID(ctx, folderID)
    if err != nil {
        return nil, fmt.Errorf("folder not found")
    }

    if !folder.IsPublic {
        return nil, fmt.Errorf("folder is not public")
    }

    items, err := s.folderItemRepo.GetByFolderID(ctx, folderID, 0, 100)
    if err != nil {
        return nil, err
    }

    return &dto.FolderDetailResponse{
        Folder: dto.FolderToFolderResponse(folder),
        Items:  dto.FolderItemsToFolderItemResponses(items),
    }, nil
}
```

---

## 5. Favorite Service Implementation

```go
// application/serviceimpl/favorite_service_impl.go

package serviceimpl

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/dto"
    "github.com/your-org/stou-smart-tour/domain/repositories"
    "github.com/your-org/stou-smart-tour/domain/services"
    "gorm.io/gorm"
)

type favoriteService struct {
    favoriteRepo repositories.FavoriteRepository
}

func NewFavoriteService(favoriteRepo repositories.FavoriteRepository) services.FavoriteService {
    return &favoriteService{
        favoriteRepo: favoriteRepo,
    }
}

// AddFavorite adds an item to favorites
func (s *favoriteService) AddFavorite(ctx context.Context, req *dto.AddFavoriteRequest, userID uuid.UUID) (*dto.FavoriteResponse, error) {
    // Check if already favorited
    if req.ExternalID != "" {
        exists, _ := s.favoriteRepo.ExistsByUserAndExternal(ctx, userID, req.Type, req.ExternalID)
        if exists {
            return nil, fmt.Errorf("already in favorites")
        }
    }

    // Create favorite
    favorite, err := dto.AddFavoriteRequestToFavorite(req, userID)
    if err != nil {
        return nil, fmt.Errorf("invalid favorite data: %w", err)
    }

    if err := s.favoriteRepo.Create(ctx, favorite); err != nil {
        return nil, fmt.Errorf("create favorite failed: %w", err)
    }

    response := dto.FavoriteToFavoriteResponse(favorite)
    return &response, nil
}

// RemoveFavorite removes an item from favorites
func (s *favoriteService) RemoveFavorite(ctx context.Context, favoriteID, userID uuid.UUID) error {
    favorite, err := s.favoriteRepo.GetByID(ctx, favoriteID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return fmt.Errorf("favorite not found")
        }
        return err
    }

    // Check ownership
    if favorite.UserID != userID {
        return fmt.Errorf("access denied")
    }

    return s.favoriteRepo.Delete(ctx, favoriteID)
}

// GetFavorites gets user's favorites with optional type filter
func (s *favoriteService) GetFavorites(ctx context.Context, userID uuid.UUID, favType string, offset, limit int) (*dto.FavoriteListResponse, error) {
    if limit <= 0 {
        limit = dto.DefaultPerPage
    }

    var favorites []models.Favorite
    var total int64
    var err error

    if favType != "" {
        favorites, err = s.favoriteRepo.GetByUserIDAndType(ctx, userID, favType, offset, limit)
        if err != nil {
            return nil, err
        }
        total, err = s.favoriteRepo.CountByUserIDAndType(ctx, userID, favType)
    } else {
        favorites, err = s.favoriteRepo.GetByUserID(ctx, userID, offset, limit)
        if err != nil {
            return nil, err
        }
        total, err = s.favoriteRepo.CountByUserID(ctx, userID)
    }

    if err != nil {
        return nil, err
    }

    return &dto.FavoriteListResponse{
        Favorites: dto.FavoritesToFavoriteResponses(favorites),
        Meta: dto.PaginationMeta{
            Total:      total,
            Page:       (offset / limit) + 1,
            PerPage:    limit,
            TotalPages: dto.CalculateTotalPages(total, limit),
        },
    }, nil
}

// IsFavorited checks if an item is in favorites
func (s *favoriteService) IsFavorited(ctx context.Context, userID uuid.UUID, favType, externalID string) (*dto.CheckFavoriteResponse, error) {
    favorite, err := s.favoriteRepo.GetByUserAndExternal(ctx, userID, favType, externalID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return &dto.CheckFavoriteResponse{
                IsFavorited: false,
            }, nil
        }
        return nil, err
    }

    return &dto.CheckFavoriteResponse{
        IsFavorited: true,
        FavoriteID:  &favorite.ID,
    }, nil
}
```

---

## 6. Translate Service Implementation

```go
// application/serviceimpl/translate_service_impl.go

package serviceimpl

import (
    "context"
    "fmt"

    "github.com/your-org/stou-smart-tour/domain/dto"
    "github.com/your-org/stou-smart-tour/domain/services"
    "github.com/your-org/stou-smart-tour/infrastructure/cache"
    "github.com/your-org/stou-smart-tour/infrastructure/external/google"
    "github.com/your-org/stou-smart-tour/infrastructure/redis"
)

type translateService struct {
    translateClient *google.TranslateClient
    redis           *redis.RedisClient
}

func NewTranslateService(
    translateClient *google.TranslateClient,
    redis *redis.RedisClient,
) services.TranslateService {
    return &translateService{
        translateClient: translateClient,
        redis:           redis,
    }
}

// Translate translates text to target language
func (s *translateService) Translate(ctx context.Context, req *dto.TranslateRequest) (*dto.TranslateResponse, error) {
    // Check cache
    cacheKey := cache.TranslateKey(req.Text, req.SourceLanguage, req.TargetLanguage)
    var cachedResult dto.TranslateResponse
    if err := s.redis.Get(ctx, cacheKey, &cachedResult); err == nil {
        return &cachedResult, nil
    }

    // Perform translation
    result, err := s.translateClient.Translate(ctx, &google.TranslateRequest{
        Text:           req.Text,
        SourceLanguage: req.SourceLanguage,
        TargetLanguage: req.TargetLanguage,
    })
    if err != nil {
        return nil, fmt.Errorf("translation failed: %w", err)
    }

    response := &dto.TranslateResponse{
        TranslatedText: result.TranslatedText,
        SourceLanguage: req.SourceLanguage,
        TargetLanguage: req.TargetLanguage,
    }

    // If source was auto-detected
    if result.DetectedSourceLanguage != "" {
        response.SourceLanguage = result.DetectedSourceLanguage
    }

    // Cache result
    _ = s.redis.Set(ctx, cacheKey, response, cache.TTLTranslate)

    return response, nil
}

// DetectLanguage detects the language of text
func (s *translateService) DetectLanguage(ctx context.Context, text string) (string, error) {
    result, err := s.translateClient.DetectLanguage(ctx, text)
    if err != nil {
        return "", fmt.Errorf("language detection failed: %w", err)
    }
    return result.Language, nil
}

// GetSupportedLanguages returns list of supported languages
func (s *translateService) GetSupportedLanguages(ctx context.Context) ([]string, error) {
    return google.SupportedLanguages, nil
}
```

---

## 7. QR Code Service Implementation

```go
// application/serviceimpl/qrcode_service_impl.go

package serviceimpl

import (
    "context"
    "encoding/base64"
    "fmt"

    "github.com/skip2/go-qrcode"
    "github.com/your-org/stou-smart-tour/domain/dto"
    "github.com/your-org/stou-smart-tour/domain/services"
)

type qrcodeService struct{}

func NewQRCodeService() services.QRCodeService {
    return &qrcodeService{}
}

// Generate generates a QR code
func (s *qrcodeService) Generate(ctx context.Context, req *dto.QRCodeRequest) (*dto.QRCodeResponse, error) {
    // Set defaults
    size := req.Size
    if size <= 0 {
        size = 256
    }
    if size > 1024 {
        size = 1024
    }

    format := req.Format
    if format == "" {
        format = "png"
    }

    // Generate QR code
    var qrData []byte
    var err error

    switch format {
    case "svg":
        // For SVG, we'd need a different library or custom generation
        // For now, fall back to PNG
        fallthrough
    case "png":
        qrData, err = qrcode.Encode(req.Content, qrcode.Medium, size)
    default:
        return nil, fmt.Errorf("unsupported format: %s", format)
    }

    if err != nil {
        return nil, fmt.Errorf("generate QR code failed: %w", err)
    }

    // Encode to base64
    encoded := base64.StdEncoding.EncodeToString(qrData)

    // Add data URI prefix for PNG
    if format == "png" {
        encoded = "data:image/png;base64," + encoded
    }

    return &dto.QRCodeResponse{
        QRCode:  encoded,
        Content: req.Content,
        Format:  format,
    }, nil
}
```

---

## 8. Update User Service (เพิ่ม student_id)

```go
// application/serviceimpl/user_service_impl.go
// เพิ่มเติมจากเดิม - support student_id

// อัพเดท CreateUser method
func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
    // Check email exists
    existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, fmt.Errorf("email already exists")
    }

    // Check username exists
    existingUser, _ = s.userRepo.GetByUsername(ctx, req.Username)
    if existingUser != nil {
        return nil, fmt.Errorf("username already exists")
    }

    // Check student_id exists (if provided)
    if req.StudentID != "" {
        existingUser, _ = s.userRepo.GetByStudentID(ctx, req.StudentID) // Need to add this method
        if existingUser != nil {
            return nil, fmt.Errorf("student ID already exists")
        }
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("hash password failed: %w", err)
    }

    // Create user
    user := &models.User{
        StudentID: req.StudentID, // 🆕 NEW
        Email:     req.Email,
        Username:  req.Username,
        Password:  string(hashedPassword),
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Role:      "user",
        IsActive:  true,
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("create user failed: %w", err)
    }

    response := dto.UserToUserResponse(user)
    return &response, nil
}

// อัพเดท Register method
func (s *userService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
    // Validate student_id is required for STOU
    if req.StudentID == "" {
        return nil, fmt.Errorf("student ID is required")
    }

    // Create user
    createReq := &dto.CreateUserRequest{
        StudentID: req.StudentID,
        Email:     req.Email,
        Username:  req.Username,
        Password:  req.Password,
        FirstName: req.FirstName,
        LastName:  req.LastName,
    }

    userResp, err := s.CreateUser(ctx, createReq)
    if err != nil {
        return nil, err
    }

    // Generate token
    token, err := s.generateToken(userResp)
    if err != nil {
        return nil, fmt.Errorf("generate token failed: %w", err)
    }

    return &dto.RegisterResponse{
        AccessToken: token,
        ExpiresIn:   3600,
        User:        *userResp,
    }, nil
}
```

---

## 9. Update DI Container

```go
// pkg/di/container.go
// เพิ่ม dependencies ใหม่

package di

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    // ... existing imports ...

    // External API clients
    "github.com/your-org/stou-smart-tour/infrastructure/external/google"
    "github.com/your-org/stou-smart-tour/infrastructure/external/openai"

    // New repositories
    // ...

    // New services
    // ...
)

type Container struct {
    // Existing fields...
    DB        *gorm.DB
    Redis     *redis.RedisClient
    Config    *config.Config

    // Repositories (existing)
    UserRepo    repositories.UserRepository
    TaskRepo    repositories.TaskRepository
    FileRepo    repositories.FileRepository
    JobRepo     repositories.JobRepository

    // Repositories (new)
    FolderRepo        repositories.FolderRepository
    FolderItemRepo    repositories.FolderItemRepository
    FavoriteRepo      repositories.FavoriteRepository
    SearchHistoryRepo repositories.SearchHistoryRepository
    AIChatRepo        repositories.AIChatRepository

    // Services (existing)
    UserService services.UserService
    TaskService services.TaskService
    FileService services.FileService
    JobService  services.JobService

    // Services (new)
    SearchService    services.SearchService
    AIService        services.AIService
    FolderService    services.FolderService
    FavoriteService  services.FavoriteService
    TranslateService services.TranslateService
    QRCodeService    services.QRCodeService

    // External clients (new)
    GoogleSearchClient  *google.SearchClient
    GooglePlacesClient  *google.PlacesClient
    YouTubeClient       *google.YouTubeClient
    GoogleTranslateClient *google.TranslateClient
    OpenAIClient        *openai.AIClient
}

func NewContainer() *Container {
    cfg := config.LoadConfig()
    container := &Container{Config: cfg}

    // Initialize database
    db, err := postgres.NewDatabase(&postgres.DatabaseConfig{
        Host:     cfg.DBHost,
        Port:     cfg.DBPort,
        User:     cfg.DBUser,
        Password: cfg.DBPassword,
        DBName:   cfg.DBName,
        SSLMode:  cfg.DBSSLMode,
    })
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    container.DB = db

    // Run migrations
    if err := postgres.Migrate(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // Initialize Redis
    redisClient, err := redis.NewRedisClient(cfg.RedisURL)
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    container.Redis = redisClient

    // Initialize external clients
    container.GoogleSearchClient = google.NewSearchClient(cfg.GoogleAPIKey, cfg.GoogleSearchEngineID)
    container.GooglePlacesClient = google.NewPlacesClient(cfg.GoogleAPIKey)
    container.YouTubeClient = google.NewYouTubeClient(cfg.GoogleAPIKey)
    container.GoogleTranslateClient = google.NewTranslateClient(cfg.GoogleAPIKey)
    container.OpenAIClient = openai.NewAIClient(cfg.OpenAIAPIKey, cfg.OpenAIModel)

    // Initialize repositories (existing)
    container.UserRepo = postgres.NewUserRepository(db)
    container.TaskRepo = postgres.NewTaskRepository(db)
    container.FileRepo = postgres.NewFileRepository(db)
    container.JobRepo = postgres.NewJobRepository(db)

    // Initialize repositories (new)
    container.FolderRepo = postgres.NewFolderRepository(db)
    container.FolderItemRepo = postgres.NewFolderItemRepository(db)
    container.FavoriteRepo = postgres.NewFavoriteRepository(db)
    container.SearchHistoryRepo = postgres.NewSearchHistoryRepository(db)
    container.AIChatRepo = postgres.NewAIChatRepository(db)

    // Initialize services (existing)
    container.UserService = serviceimpl.NewUserService(container.UserRepo, cfg)
    container.TaskService = serviceimpl.NewTaskService(container.TaskRepo)
    container.FileService = serviceimpl.NewFileService(container.FileRepo, nil) // Add storage if needed
    container.JobService = serviceimpl.NewJobService(container.JobRepo, nil)    // Add scheduler if needed

    // Initialize services (new)
    container.SearchService = serviceimpl.NewSearchService(
        container.GoogleSearchClient,
        container.GooglePlacesClient,
        container.SearchHistoryRepo,
        container.Redis,
    )

    container.AIService = serviceimpl.NewAIService(
        container.OpenAIClient,
        container.GoogleSearchClient,
        container.YouTubeClient,
        container.AIChatRepo,
        container.Redis,
    )

    container.FolderService = serviceimpl.NewFolderService(
        container.FolderRepo,
        container.FolderItemRepo,
        cfg,
    )

    container.FavoriteService = serviceimpl.NewFavoriteService(container.FavoriteRepo)

    container.TranslateService = serviceimpl.NewTranslateService(
        container.GoogleTranslateClient,
        container.Redis,
    )

    container.QRCodeService = serviceimpl.NewQRCodeService()

    // Setup graceful shutdown
    container.setupGracefulShutdown()

    return container
}

func (c *Container) setupGracefulShutdown() {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-quit
        log.Println("Shutting down...")

        // Close Redis connection
        if c.Redis != nil {
            c.Redis.Close()
        }

        // Close database connection
        if c.DB != nil {
            sqlDB, _ := c.DB.DB()
            sqlDB.Close()
        }

        os.Exit(0)
    }()
}
```

---

## 10. Update Config

```go
// pkg/config/config.go
// เพิ่ม config fields ใหม่

package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    // Existing fields
    AppEnv      string
    AppPort     string
    DBHost      string
    DBPort      string
    DBUser      string
    DBPassword  string
    DBName      string
    DBSSLMode   string
    RedisURL    string
    JWTSecret   string

    // New fields
    AppURL              string  // For generating share links
    GoogleAPIKey        string
    GoogleSearchEngineID string
    OpenAIAPIKey        string
    OpenAIModel         string

    // Rate limiting
    RateLimitSearch  int
    RateLimitAI      int
    RateLimitGeneral int
}

func LoadConfig() *Config {
    godotenv.Load()

    return &Config{
        // Existing
        AppEnv:     getEnv("APP_ENV", "development"),
        AppPort:    getEnv("APP_PORT", "8080"),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "stou_smart_tour"),
        DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
        RedisURL:   getEnv("REDIS_URL", "redis://localhost:6379"),
        JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),

        // New
        AppURL:               getEnv("APP_URL", "http://localhost:3000"),
        GoogleAPIKey:         getEnv("GOOGLE_API_KEY", ""),
        GoogleSearchEngineID: getEnv("GOOGLE_SEARCH_ENGINE_ID", ""),
        OpenAIAPIKey:         getEnv("OPENAI_API_KEY", ""),
        OpenAIModel:          getEnv("OPENAI_MODEL", "gpt-4-turbo-preview"),

        // Rate limiting
        RateLimitSearch:  getEnvInt("RATE_LIMIT_SEARCH", 30),
        RateLimitAI:      getEnvInt("RATE_LIMIT_AI", 10),
        RateLimitGeneral: getEnvInt("RATE_LIMIT_GENERAL", 100),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}
```

---

## 11. Summary - Application Layer Files

```
application/
└── serviceimpl/
    ├── user_service_impl.go       # ✏️ UPDATE (add student_id support)
    ├── task_service_impl.go       # ✅ KEEP
    ├── file_service_impl.go       # ✅ KEEP
    ├── job_service_impl.go        # ✅ KEEP
    ├── search_service_impl.go     # 🆕 NEW
    ├── ai_service_impl.go         # 🆕 NEW
    ├── folder_service_impl.go     # 🆕 NEW
    ├── favorite_service_impl.go   # 🆕 NEW
    ├── translate_service_impl.go  # 🆕 NEW
    └── qrcode_service_impl.go     # 🆕 NEW

pkg/
├── config/
│   └── config.go                  # ✏️ UPDATE (add new config fields)
└── di/
    └── container.go               # ✏️ UPDATE (add new dependencies)
```

---

## Next Part

➡️ ไปต่อที่ **Part 5: Interface Layer (Handlers, Routes, Middleware)**
- Search Handler & Routes
- AI Handler & Routes
- Folder Handler & Routes
- Favorite Handler & Routes
- Utility Handler & Routes
- Rate Limit Middleware

---

*Document Version: 1.0*
*Part: 4 of 5*
