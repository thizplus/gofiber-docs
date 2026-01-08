package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gofiber-template/pkg/logger"
)

const (
	openaiChatURL = "https://api.openai.com/v1/chat/completions"
)

// AIClient handles OpenAI API
type AIClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewAIClient creates a new OpenAI client
func NewAIClient(apiKey, model string) *AIClient {
	if model == "" {
		model = "gpt-4-turbo-preview"
	}
	return &AIClient{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Role    string `json:"role"` // system, user, assistant
	Content string `json:"content"`
}

// ChatRequest represents the request to OpenAI API
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// ChatResponse represents the response from OpenAI API
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a single choice in the response
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// Usage contains token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Chat sends a chat completion request
func (c *AIClient) Chat(ctx context.Context, messages []ChatMessage, maxTokens int, temperature float64) (*ChatResponse, error) {
	startTime := time.Now()

	if maxTokens == 0 {
		maxTokens = 2000
	}
	if temperature == 0 {
		temperature = 0.7
	}

	logger.InfoContext(ctx, "OpenAI Chat request started",
		"model", c.model,
		"message_count", len(messages),
		"max_tokens", maxTokens,
	)

	reqBody := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		logger.ErrorContext(ctx, "OpenAI request marshal failed",
			"error", err.Error(),
		)
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", openaiChatURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.ErrorContext(ctx, "OpenAI request creation failed",
			"error", err.Error(),
		)
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.ErrorContext(ctx, "OpenAI HTTP request failed",
			"error", err.Error(),
			"model", c.model,
			"response_time_ms", time.Since(startTime).Milliseconds(),
		)
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorContext(ctx, "OpenAI response read failed",
			"error", err.Error(),
			"status_code", resp.StatusCode,
		)
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logger.ErrorContext(ctx, "OpenAI API error response",
			"status_code", resp.StatusCode,
			"response_body", string(body),
			"model", c.model,
			"response_time_ms", time.Since(startTime).Milliseconds(),
		)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		logger.ErrorContext(ctx, "OpenAI response parse failed",
			"error", err.Error(),
			"response_body", string(body),
		)
		return nil, fmt.Errorf("decode response: %w", err)
	}

	logger.InfoContext(ctx, "OpenAI Chat request completed",
		"model", c.model,
		"prompt_tokens", result.Usage.PromptTokens,
		"completion_tokens", result.Usage.CompletionTokens,
		"total_tokens", result.Usage.TotalTokens,
		"response_time_ms", time.Since(startTime).Milliseconds(),
	)

	return &result, nil
}

// SearchResultContext represents context from search results
type SearchResultContext struct {
	Title   string
	Snippet string
	URL     string
}

// GetSystemPromptByLang returns system prompt based on language
func GetSystemPromptByLang(lang string) string {
	if lang == "en" {
		return `You are a travel information assistant for students at Sukhothai Thammathirat Open University (STOU).
Summarize information from the provided sources concisely, clearly, and usefully.

Rules for responding:
1. Respond in English
2. Format as Markdown
3. Summarize with main headings and bullet points
4. Cite the source of information
5. If there is price information, opening hours, or important details, include them
6. Suggest 2-3 relevant follow-up questions`
	}
	return `คุณเป็นผู้ช่วยค้นหาข้อมูลท่องเที่ยวสำหรับนักศึกษามหาวิทยาลัยสุโขทัยธรรมาธิราช (มสธ.)
ให้สรุปข้อมูลจาก sources ที่ได้รับอย่างกระชับ ชัดเจน และเป็นประโยชน์

กฎในการตอบ:
1. ตอบเป็นภาษาไทย
2. จัดรูปแบบเป็น Markdown
3. สรุปเป็นหัวข้อหลักๆ พร้อม bullet points
4. ระบุ source ที่มาของข้อมูล
5. หากมีข้อมูลราคา เวลาเปิด-ปิด หรือข้อมูลสำคัญ ให้ระบุด้วย
6. เสนอคำถาม follow-up ที่เกี่ยวข้อง 2-3 ข้อ`
}

// GetChatSystemPromptByLang returns chat system prompt based on language
func GetChatSystemPromptByLang(lang string) string {
	if lang == "en" {
		return `You are a travel information assistant for STOU students.
Answer travel-related questions in a friendly manner and provide useful information.
Respond in English and use Markdown format.`
	}
	return `คุณเป็นผู้ช่วยค้นหาข้อมูลท่องเที่ยวสำหรับนักศึกษา มสธ.
ตอบคำถามเกี่ยวกับการท่องเที่ยวอย่างเป็นมิตรและให้ข้อมูลที่เป็นประโยชน์
ตอบเป็นภาษาไทยและใช้ Markdown format`
}

// GenerateTravelSummary generates a travel summary from search results
func (c *AIClient) GenerateTravelSummary(ctx context.Context, query string, searchResults []SearchResultContext, lang string) (*ChatResponse, error) {
	logger.InfoContext(ctx, "GenerateTravelSummary started",
		"query", query,
		"lang", lang,
		"search_results_count", len(searchResults),
	)

	systemPrompt := GetSystemPromptByLang(lang)

	// Build user prompt with search results
	var userPrompt string
	if lang == "en" {
		userPrompt = fmt.Sprintf("Search query: %s\n\nInformation from various sources:\n", query)
		for i, result := range searchResults {
			userPrompt += fmt.Sprintf("\n[Source %d: %s]\n%s\nURL: %s\n",
				i+1, result.Title, result.Snippet, result.URL)
		}
		userPrompt += "\nPlease summarize the above information systematically."
	} else {
		userPrompt = fmt.Sprintf("คำค้นหา: %s\n\nข้อมูลจากแหล่งต่างๆ:\n", query)
		for i, result := range searchResults {
			userPrompt += fmt.Sprintf("\n[Source %d: %s]\n%s\nURL: %s\n",
				i+1, result.Title, result.Snippet, result.URL)
		}
		userPrompt += "\nกรุณาสรุปข้อมูลข้างต้นอย่างเป็นระบบ"
	}

	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	response, err := c.Chat(ctx, messages, 2000, 0.7)
	if err != nil {
		logger.ErrorContext(ctx, "GenerateTravelSummary failed",
			"query", query,
			"lang", lang,
			"error", err.Error(),
		)
		return nil, err
	}

	logger.InfoContext(ctx, "GenerateTravelSummary completed",
		"query", query,
		"lang", lang,
	)

	return response, nil
}

// ContinueChat continues an existing chat conversation
func (c *AIClient) ContinueChat(ctx context.Context, history []ChatMessage, newMessage string, lang string) (*ChatResponse, error) {
	logger.InfoContext(ctx, "ContinueChat started",
		"lang", lang,
		"history_count", len(history),
		"new_message_length", len(newMessage),
	)

	systemPrompt := GetChatSystemPromptByLang(lang)

	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
	}
	messages = append(messages, history...)
	messages = append(messages, ChatMessage{Role: "user", Content: newMessage})

	response, err := c.Chat(ctx, messages, 1500, 0.7)
	if err != nil {
		logger.ErrorContext(ctx, "ContinueChat failed",
			"lang", lang,
			"history_count", len(history),
			"error", err.Error(),
		)
		return nil, err
	}

	logger.InfoContext(ctx, "ContinueChat completed",
		"lang", lang,
	)

	return response, nil
}
