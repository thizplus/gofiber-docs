package google

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	translateURL = "https://translation.googleapis.com/language/translate/v2"
	detectURL    = "https://translation.googleapis.com/language/translate/v2/detect"
)

// TranslateClient handles Google Translate API
type TranslateClient struct {
	*GoogleClient
}

// NewTranslateClient creates a new Translate client
func NewTranslateClient(apiKey string) *TranslateClient {
	return &TranslateClient{
		GoogleClient: NewGoogleClient(apiKey),
	}
}

// TranslateRequest represents translation request
type TranslateRequest struct {
	Text           string
	SourceLanguage string // Optional, auto-detect if empty
	TargetLanguage string
}

// TranslateResponse represents translation API response
type TranslateResponse struct {
	Data TranslateData `json:"data"`
}

// TranslateData contains translations
type TranslateData struct {
	Translations []Translation `json:"translations"`
}

// Translation represents a single translation
type Translation struct {
	TranslatedText         string `json:"translatedText"`
	DetectedSourceLanguage string `json:"detectedSourceLanguage,omitempty"`
}

// DetectResponse represents language detection response
type DetectResponse struct {
	Data DetectData `json:"data"`
}

// DetectData contains detections
type DetectData struct {
	Detections [][]Detection `json:"detections"`
}

// Detection represents a single detection
type Detection struct {
	Language   string  `json:"language"`
	IsReliable bool    `json:"isReliable"`
	Confidence float64 `json:"confidence"`
}

// Translate translates text to target language
func (c *TranslateClient) Translate(ctx context.Context, req *TranslateRequest) (*Translation, error) {
	body := map[string]interface{}{
		"q":      req.Text,
		"target": req.TargetLanguage,
		"format": "text",
	}

	if req.SourceLanguage != "" {
		body["source"] = req.SourceLanguage
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	requestURL := fmt.Sprintf("%s?key=%s", translateURL, c.apiKey)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var result TranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Data.Translations) == 0 {
		return nil, fmt.Errorf("no translation returned")
	}

	return &result.Data.Translations[0], nil
}

// DetectLanguage detects the language of text
func (c *TranslateClient) DetectLanguage(ctx context.Context, text string) (*Detection, error) {
	body := map[string]interface{}{
		"q": text,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	requestURL := fmt.Sprintf("%s?key=%s", detectURL, c.apiKey)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var result DetectResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Data.Detections) == 0 || len(result.Data.Detections[0]) == 0 {
		return nil, fmt.Errorf("no detection returned")
	}

	return &result.Data.Detections[0][0], nil
}

// SupportedLanguages - commonly used languages
var SupportedLanguages = []string{
	"th", // Thai
	"en", // English
	"zh", // Chinese
	"ja", // Japanese
	"ko", // Korean
	"vi", // Vietnamese
	"ms", // Malay
	"id", // Indonesian
	"fr", // French
	"de", // German
	"es", // Spanish
	"ru", // Russian
}
