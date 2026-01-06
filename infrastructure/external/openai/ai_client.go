package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
	if maxTokens == 0 {
		maxTokens = 2000
	}
	if temperature == 0 {
		temperature = 0.7
	}

	reqBody := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", openaiChatURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

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

	return c.Chat(ctx, messages, 2000, 0.7)
}

// ContinueChat continues an existing chat conversation
func (c *AIClient) ContinueChat(ctx context.Context, history []ChatMessage, newMessage string, lang string) (*ChatResponse, error) {
	systemPrompt := GetChatSystemPromptByLang(lang)

	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
	}
	messages = append(messages, history...)
	messages = append(messages, ChatMessage{Role: "user", Content: newMessage})

	return c.Chat(ctx, messages, 1500, 0.7)
}
