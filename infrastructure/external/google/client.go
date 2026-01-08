package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gofiber-template/pkg/logger"
)

// GoogleClient is the base client for all Google APIs
type GoogleClient struct {
	apiKey     string
	httpClient *http.Client
}

// NewGoogleClient creates a new Google API client
func NewGoogleClient(apiKey string) *GoogleClient {
	return &GoogleClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP GET request and decodes the JSON response
func (c *GoogleClient) doRequest(ctx context.Context, url string, result interface{}) error {
	startTime := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.ErrorContext(ctx, "Google API request creation failed",
			"error", err.Error(),
		)
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.ErrorContext(ctx, "Google API HTTP request failed",
			"error", err.Error(),
			"response_time_ms", time.Since(startTime).Milliseconds(),
		)
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorContext(ctx, "Google API response read failed",
			"error", err.Error(),
			"status_code", resp.StatusCode,
		)
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logger.ErrorContext(ctx, "Google API error response",
			"status_code", resp.StatusCode,
			"response_body", string(body),
			"response_time_ms", time.Since(startTime).Milliseconds(),
		)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, result); err != nil {
		logger.ErrorContext(ctx, "Google API response parse failed",
			"error", err.Error(),
			"response_body", string(body),
		)
		return fmt.Errorf("decode response: %w", err)
	}

	logger.DebugContext(ctx, "Google API request completed",
		"response_time_ms", time.Since(startTime).Milliseconds(),
	)

	return nil
}

// GetAPIKey returns the API key
func (c *GoogleClient) GetAPIKey() string {
	return c.apiKey
}
