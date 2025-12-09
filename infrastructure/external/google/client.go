package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

// GetAPIKey returns the API key
func (c *GoogleClient) GetAPIKey() string {
	return c.apiKey
}
