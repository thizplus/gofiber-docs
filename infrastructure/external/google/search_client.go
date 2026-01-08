package google

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"gofiber-template/pkg/logger"
)

const (
	customSearchBaseURL = "https://www.googleapis.com/customsearch/v1"
)

// SearchClient handles Google Custom Search API
type SearchClient struct {
	*GoogleClient
	searchEngineID string
}

// NewSearchClient creates a new Search client
func NewSearchClient(apiKey, searchEngineID string) *SearchClient {
	return &SearchClient{
		GoogleClient:   NewGoogleClient(apiKey),
		searchEngineID: searchEngineID,
	}
}

// Search types
const (
	SearchTypeAll   = ""
	SearchTypeImage = "image"
)

// SearchRequest represents search parameters
type SearchRequest struct {
	Query      string
	SearchType string // "", "image"
	Start      int    // pagination start (1-based)
	Num        int    // results per page (max 10)
	Language   string // hl parameter
	SafeSearch string // safe parameter: off, medium, high
}

// SearchResponse represents Google Custom Search API response
type SearchResponse struct {
	Kind              string `json:"kind"`
	SearchInformation struct {
		TotalResults          string  `json:"totalResults"`
		FormattedTotalResults string  `json:"formattedTotalResults"`
		SearchTime            float64 `json:"searchTime"`
	} `json:"searchInformation"`
	Items []SearchItem `json:"items"`
}

// SearchItem represents a single search result
type SearchItem struct {
	Kind             string `json:"kind"`
	Title            string `json:"title"`
	HTMLTitle        string `json:"htmlTitle"`
	Link             string `json:"link"`
	DisplayLink      string `json:"displayLink"`
	Snippet          string `json:"snippet"`
	HTMLSnippet      string `json:"htmlSnippet"`
	FormattedURL     string `json:"formattedUrl"`
	HTMLFormattedURL string `json:"htmlFormattedUrl"`
	Pagemap          struct {
		CSEThumbnail []struct {
			Src    string `json:"src"`
			Width  string `json:"width"`
			Height string `json:"height"`
		} `json:"cse_thumbnail"`
		CSEImage []struct {
			Src string `json:"src"`
		} `json:"cse_image"`
		Metatags []map[string]string `json:"metatags"`
	} `json:"pagemap"`
	Image *ImageInfo `json:"image,omitempty"`
}

// ImageInfo contains image-specific information
type ImageInfo struct {
	ContextLink   string `json:"contextLink"`
	Height        int    `json:"height"`
	Width         int    `json:"width"`
	ByteSize      int    `json:"byteSize"`
	ThumbnailLink string `json:"thumbnailLink"`
}

// Search performs a Google Custom Search
func (c *SearchClient) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	startTime := time.Now()

	logger.InfoContext(ctx, "Google Search request started",
		"query", req.Query,
		"search_type", req.SearchType,
		"start", req.Start,
		"num", req.Num,
	)

	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("cx", c.searchEngineID)
	params.Set("q", req.Query)

	if req.SearchType == SearchTypeImage {
		params.Set("searchType", "image")
	}

	if req.Start > 0 {
		params.Set("start", strconv.Itoa(req.Start))
	}

	if req.Num > 0 && req.Num <= 10 {
		params.Set("num", strconv.Itoa(req.Num))
	}

	if req.Language != "" {
		params.Set("hl", req.Language)
	}

	if req.SafeSearch != "" {
		params.Set("safe", req.SafeSearch)
	}

	searchURL := fmt.Sprintf("%s?%s", customSearchBaseURL, params.Encode())

	var result SearchResponse
	if err := c.doRequest(ctx, searchURL, &result); err != nil {
		logger.ErrorContext(ctx, "Google Search request failed",
			"query", req.Query,
			"error", err.Error(),
			"response_time_ms", time.Since(startTime).Milliseconds(),
		)
		return nil, err
	}

	logger.InfoContext(ctx, "Google Search request completed",
		"query", req.Query,
		"results_count", len(result.Items),
		"total_results", result.SearchInformation.TotalResults,
		"search_time_sec", result.SearchInformation.SearchTime,
		"response_time_ms", time.Since(startTime).Milliseconds(),
	)

	return &result, nil
}

// SearchAll performs a general web search
func (c *SearchClient) SearchAll(ctx context.Context, query string, page, perPage int) (*SearchResponse, error) {
	start := (page-1)*perPage + 1
	if perPage > 10 {
		perPage = 10 // Google API limit
	}

	return c.Search(ctx, &SearchRequest{
		Query: query,
		Start: start,
		Num:   perPage,
	})
}

// SearchImages performs an image search
func (c *SearchClient) SearchImages(ctx context.Context, query string, page, perPage int) (*SearchResponse, error) {
	start := (page-1)*perPage + 1
	if perPage > 10 {
		perPage = 10
	}

	return c.Search(ctx, &SearchRequest{
		Query:      query,
		SearchType: SearchTypeImage,
		Start:      start,
		Num:        perPage,
	})
}
