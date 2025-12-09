# STOU Smart Tour - Backend Development Plan
# Part 3: Infrastructure Layer (External APIs, Cache, Repository Impl)

---

## Table of Contents - All Parts

| Part | หัวข้อ | สถานะ |
|------|--------|-------|
| Part 1 | Project Overview & Foundation | ✅ Done |
| Part 2 | Domain Layer (Models, DTOs, Interfaces) | ✅ Done |
| **Part 3** | **Infrastructure Layer (External APIs, Cache)** | 📍 Current |
| Part 4 | Application Layer (Services Implementation) | ⏳ Pending |
| Part 5 | Interface Layer (Handlers, Routes, Middleware) | ⏳ Pending |

---

## 1. External API Clients (infrastructure/external/)

### 1.1 โครงสร้าง External Folder

```
infrastructure/
└── external/
    ├── google/
    │   ├── client.go              # Base Google client
    │   ├── search_client.go       # Google Custom Search
    │   ├── places_client.go       # Google Places API
    │   ├── youtube_client.go      # YouTube Data API
    │   └── translate_client.go    # Google Translate API
    │
    └── openai/
        └── ai_client.go           # OpenAI/Anthropic client
```

### 1.2 Google Base Client

```go
// infrastructure/external/google/client.go

package google

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type GoogleClient struct {
    apiKey     string
    httpClient *http.Client
}

func NewGoogleClient(apiKey string) *GoogleClient {
    return &GoogleClient{
        apiKey: apiKey,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

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

func (c *GoogleClient) GetAPIKey() string {
    return c.apiKey
}
```

### 1.3 Google Custom Search Client

```go
// infrastructure/external/google/search_client.go

package google

import (
    "context"
    "fmt"
    "net/url"
    "strconv"
)

const (
    customSearchBaseURL = "https://www.googleapis.com/customsearch/v1"
)

type SearchClient struct {
    *GoogleClient
    searchEngineID string
}

func NewSearchClient(apiKey, searchEngineID string) *SearchClient {
    return &SearchClient{
        GoogleClient:   NewGoogleClient(apiKey),
        searchEngineID: searchEngineID,
    }
}

// Search types
const (
    SearchTypeAll     = ""
    SearchTypeImage   = "image"
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
        TotalResults         string  `json:"totalResults"`
        FormattedTotalResults string `json:"formattedTotalResults"`
        SearchTime           float64 `json:"searchTime"`
    } `json:"searchInformation"`
    Items []SearchItem `json:"items"`
}

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
    Image *ImageInfo `json:"image,omitempty"` // Only for image search
}

type ImageInfo struct {
    ContextLink   string `json:"contextLink"`
    Height        int    `json:"height"`
    Width         int    `json:"width"`
    ByteSize      int    `json:"byteSize"`
    ThumbnailLink string `json:"thumbnailLink"`
}

// Search performs a Google Custom Search
func (c *SearchClient) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
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
        return nil, err
    }

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
```

### 1.4 Google Places Client

```go
// infrastructure/external/google/places_client.go

package google

import (
    "context"
    "fmt"
    "net/url"
    "strconv"
    "strings"
)

const (
    placesNearbyURL  = "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
    placeDetailsURL  = "https://maps.googleapis.com/maps/api/place/details/json"
    placePhotoURL    = "https://maps.googleapis.com/maps/api/place/photo"
)

type PlacesClient struct {
    *GoogleClient
}

func NewPlacesClient(apiKey string) *PlacesClient {
    return &PlacesClient{
        GoogleClient: NewGoogleClient(apiKey),
    }
}

// NearbySearchRequest represents nearby search parameters
type NearbySearchRequest struct {
    Lat      float64
    Lng      float64
    Radius   int    // meters
    Type     string // restaurant, tourist_attraction, etc.
    Keyword  string
    Language string
}

// NearbySearchResponse represents Places API nearby search response
type NearbySearchResponse struct {
    Status           string  `json:"status"`
    Results          []Place `json:"results"`
    NextPageToken    string  `json:"next_page_token,omitempty"`
    HTMLAttributions []string `json:"html_attributions"`
    ErrorMessage     string  `json:"error_message,omitempty"`
}

type Place struct {
    PlaceID          string   `json:"place_id"`
    Name             string   `json:"name"`
    Vicinity         string   `json:"vicinity"`
    FormattedAddress string   `json:"formatted_address,omitempty"`
    Geometry         Geometry `json:"geometry"`
    Rating           float64  `json:"rating"`
    UserRatingsTotal int      `json:"user_ratings_total"`
    Types            []string `json:"types"`
    Photos           []Photo  `json:"photos,omitempty"`
    OpeningHours     *struct {
        OpenNow bool `json:"open_now"`
    } `json:"opening_hours,omitempty"`
    PriceLevel   int    `json:"price_level,omitempty"`
    BusinessStatus string `json:"business_status,omitempty"`
}

type Geometry struct {
    Location struct {
        Lat float64 `json:"lat"`
        Lng float64 `json:"lng"`
    } `json:"location"`
}

type Photo struct {
    PhotoReference   string   `json:"photo_reference"`
    Height           int      `json:"height"`
    Width            int      `json:"width"`
    HTMLAttributions []string `json:"html_attributions"`
}

// PlaceDetailsRequest represents place details request
type PlaceDetailsRequest struct {
    PlaceID  string
    Fields   []string // Optional fields to return
    Language string
}

// PlaceDetailsResponse represents place details response
type PlaceDetailsResponse struct {
    Status       string       `json:"status"`
    Result       PlaceDetails `json:"result"`
    ErrorMessage string       `json:"error_message,omitempty"`
}

type PlaceDetails struct {
    PlaceID              string   `json:"place_id"`
    Name                 string   `json:"name"`
    FormattedAddress     string   `json:"formatted_address"`
    FormattedPhoneNumber string   `json:"formatted_phone_number,omitempty"`
    InternationalPhone   string   `json:"international_phone_number,omitempty"`
    Website              string   `json:"website,omitempty"`
    Rating               float64  `json:"rating"`
    UserRatingsTotal     int      `json:"user_ratings_total"`
    PriceLevel           int      `json:"price_level,omitempty"`
    Types                []string `json:"types"`
    Geometry             Geometry `json:"geometry"`
    Photos               []Photo  `json:"photos,omitempty"`
    Reviews              []Review `json:"reviews,omitempty"`
    OpeningHours         *OpeningHoursDetail `json:"opening_hours,omitempty"`
    URL                  string   `json:"url"` // Google Maps URL
}

type Review struct {
    AuthorName              string `json:"author_name"`
    AuthorURL               string `json:"author_url,omitempty"`
    ProfilePhotoURL         string `json:"profile_photo_url,omitempty"`
    Rating                  int    `json:"rating"`
    Text                    string `json:"text"`
    Time                    int64  `json:"time"`
    RelativeTimeDescription string `json:"relative_time_description"`
}

type OpeningHoursDetail struct {
    OpenNow     bool     `json:"open_now"`
    WeekdayText []string `json:"weekday_text"`
}

// NearbySearch searches for places nearby
func (c *PlacesClient) NearbySearch(ctx context.Context, req *NearbySearchRequest) (*NearbySearchResponse, error) {
    params := url.Values{}
    params.Set("key", c.apiKey)
    params.Set("location", fmt.Sprintf("%f,%f", req.Lat, req.Lng))
    params.Set("radius", strconv.Itoa(req.Radius))

    if req.Type != "" {
        params.Set("type", req.Type)
    }
    if req.Keyword != "" {
        params.Set("keyword", req.Keyword)
    }
    if req.Language != "" {
        params.Set("language", req.Language)
    } else {
        params.Set("language", "th")
    }

    searchURL := fmt.Sprintf("%s?%s", placesNearbyURL, params.Encode())

    var result NearbySearchResponse
    if err := c.doRequest(ctx, searchURL, &result); err != nil {
        return nil, err
    }

    if result.Status != "OK" && result.Status != "ZERO_RESULTS" {
        return nil, fmt.Errorf("Places API error: %s - %s", result.Status, result.ErrorMessage)
    }

    return &result, nil
}

// GetPlaceDetails gets detailed information about a place
func (c *PlacesClient) GetPlaceDetails(ctx context.Context, req *PlaceDetailsRequest) (*PlaceDetailsResponse, error) {
    params := url.Values{}
    params.Set("key", c.apiKey)
    params.Set("place_id", req.PlaceID)

    if len(req.Fields) > 0 {
        params.Set("fields", strings.Join(req.Fields, ","))
    } else {
        // Default fields
        params.Set("fields", "place_id,name,formatted_address,formatted_phone_number,website,rating,user_ratings_total,price_level,types,geometry,photos,reviews,opening_hours,url")
    }

    if req.Language != "" {
        params.Set("language", req.Language)
    } else {
        params.Set("language", "th")
    }

    detailsURL := fmt.Sprintf("%s?%s", placeDetailsURL, params.Encode())

    var result PlaceDetailsResponse
    if err := c.doRequest(ctx, detailsURL, &result); err != nil {
        return nil, err
    }

    if result.Status != "OK" {
        return nil, fmt.Errorf("Places API error: %s - %s", result.Status, result.ErrorMessage)
    }

    return &result, nil
}

// GetPhotoURL generates a URL for a place photo
func (c *PlacesClient) GetPhotoURL(photoReference string, maxWidth int) string {
    params := url.Values{}
    params.Set("key", c.apiKey)
    params.Set("photoreference", photoReference)
    params.Set("maxwidth", strconv.Itoa(maxWidth))

    return fmt.Sprintf("%s?%s", placePhotoURL, params.Encode())
}

// CalculateDistance calculates distance between two points (Haversine formula)
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
    const R = 6371000 // Earth's radius in meters

    lat1Rad := lat1 * 3.14159265359 / 180
    lat2Rad := lat2 * 3.14159265359 / 180
    deltaLat := (lat2 - lat1) * 3.14159265359 / 180
    deltaLng := (lng2 - lng1) * 3.14159265359 / 180

    a := sin(deltaLat/2)*sin(deltaLat/2) +
        cos(lat1Rad)*cos(lat2Rad)*sin(deltaLng/2)*sin(deltaLng/2)
    c := 2 * atan2(sqrt(a), sqrt(1-a))

    return R * c
}

func sin(x float64) float64  { return math.Sin(x) }
func cos(x float64) float64  { return math.Cos(x) }
func sqrt(x float64) float64 { return math.Sqrt(x) }
func atan2(y, x float64) float64 { return math.Atan2(y, x) }
```

### 1.5 YouTube Data API Client

```go
// infrastructure/external/google/youtube_client.go

package google

import (
    "context"
    "fmt"
    "net/url"
    "strconv"
)

const (
    youtubeSearchURL = "https://www.googleapis.com/youtube/v3/search"
    youtubeVideosURL = "https://www.googleapis.com/youtube/v3/videos"
)

type YouTubeClient struct {
    *GoogleClient
}

func NewYouTubeClient(apiKey string) *YouTubeClient {
    return &YouTubeClient{
        GoogleClient: NewGoogleClient(apiKey),
    }
}

// VideoSearchRequest represents YouTube search parameters
type VideoSearchRequest struct {
    Query      string
    MaxResults int    // max 50
    Order      string // relevance, date, rating, viewCount
    Language   string
    RegionCode string
}

// VideoSearchResponse represents YouTube search API response
type VideoSearchResponse struct {
    Kind          string `json:"kind"`
    PageInfo      PageInfo `json:"pageInfo"`
    NextPageToken string `json:"nextPageToken,omitempty"`
    Items         []VideoSearchItem `json:"items"`
}

type PageInfo struct {
    TotalResults   int `json:"totalResults"`
    ResultsPerPage int `json:"resultsPerPage"`
}

type VideoSearchItem struct {
    Kind    string `json:"kind"`
    ID      VideoID `json:"id"`
    Snippet VideoSnippet `json:"snippet"`
}

type VideoID struct {
    Kind    string `json:"kind"`
    VideoID string `json:"videoId"`
}

type VideoSnippet struct {
    PublishedAt  string `json:"publishedAt"`
    ChannelID    string `json:"channelId"`
    Title        string `json:"title"`
    Description  string `json:"description"`
    ChannelTitle string `json:"channelTitle"`
    Thumbnails   Thumbnails `json:"thumbnails"`
    LiveBroadcastContent string `json:"liveBroadcastContent"`
}

type Thumbnails struct {
    Default  Thumbnail `json:"default"`
    Medium   Thumbnail `json:"medium"`
    High     Thumbnail `json:"high"`
    Standard Thumbnail `json:"standard,omitempty"`
    Maxres   Thumbnail `json:"maxres,omitempty"`
}

type Thumbnail struct {
    URL    string `json:"url"`
    Width  int    `json:"width"`
    Height int    `json:"height"`
}

// VideoDetailsResponse for video statistics and content details
type VideoDetailsResponse struct {
    Items []VideoDetails `json:"items"`
}

type VideoDetails struct {
    ID             string `json:"id"`
    ContentDetails ContentDetails `json:"contentDetails"`
    Statistics     Statistics `json:"statistics"`
}

type ContentDetails struct {
    Duration string `json:"duration"` // ISO 8601 format (PT1H2M3S)
}

type Statistics struct {
    ViewCount    string `json:"viewCount"`
    LikeCount    string `json:"likeCount"`
    CommentCount string `json:"commentCount"`
}

// SearchVideos searches for videos on YouTube
func (c *YouTubeClient) SearchVideos(ctx context.Context, req *VideoSearchRequest) (*VideoSearchResponse, error) {
    params := url.Values{}
    params.Set("key", c.apiKey)
    params.Set("q", req.Query)
    params.Set("part", "snippet")
    params.Set("type", "video")

    if req.MaxResults > 0 && req.MaxResults <= 50 {
        params.Set("maxResults", strconv.Itoa(req.MaxResults))
    } else {
        params.Set("maxResults", "10")
    }

    if req.Order != "" {
        params.Set("order", req.Order)
    } else {
        params.Set("order", "relevance")
    }

    if req.RegionCode != "" {
        params.Set("regionCode", req.RegionCode)
    } else {
        params.Set("regionCode", "TH")
    }

    if req.Language != "" {
        params.Set("relevanceLanguage", req.Language)
    }

    searchURL := fmt.Sprintf("%s?%s", youtubeSearchURL, params.Encode())

    var result VideoSearchResponse
    if err := c.doRequest(ctx, searchURL, &result); err != nil {
        return nil, err
    }

    return &result, nil
}

// GetVideoDetails gets video statistics and content details
func (c *YouTubeClient) GetVideoDetails(ctx context.Context, videoIDs []string) (*VideoDetailsResponse, error) {
    if len(videoIDs) == 0 {
        return &VideoDetailsResponse{}, nil
    }

    params := url.Values{}
    params.Set("key", c.apiKey)
    params.Set("id", strings.Join(videoIDs, ","))
    params.Set("part", "contentDetails,statistics")

    detailsURL := fmt.Sprintf("%s?%s", youtubeVideosURL, params.Encode())

    var result VideoDetailsResponse
    if err := c.doRequest(ctx, detailsURL, &result); err != nil {
        return nil, err
    }

    return &result, nil
}

// ParseDuration converts ISO 8601 duration to human readable format
// e.g., PT1H2M3S -> "1:02:03"
func ParseDuration(isoDuration string) string {
    // Simple parser for YouTube duration format
    // PT1H2M3S -> 1:02:03
    // PT2M3S -> 2:03
    // PT45S -> 0:45

    duration := isoDuration
    duration = strings.TrimPrefix(duration, "PT")

    hours := 0
    minutes := 0
    seconds := 0

    if idx := strings.Index(duration, "H"); idx != -1 {
        hours, _ = strconv.Atoi(duration[:idx])
        duration = duration[idx+1:]
    }
    if idx := strings.Index(duration, "M"); idx != -1 {
        minutes, _ = strconv.Atoi(duration[:idx])
        duration = duration[idx+1:]
    }
    if idx := strings.Index(duration, "S"); idx != -1 {
        seconds, _ = strconv.Atoi(duration[:idx])
    }

    if hours > 0 {
        return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
    }
    return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// GetVideoURL returns the YouTube video URL
func GetVideoURL(videoID string) string {
    return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
}
```

### 1.6 Google Translate Client

```go
// infrastructure/external/google/translate_client.go

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

type TranslateClient struct {
    *GoogleClient
}

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

type TranslateData struct {
    Translations []Translation `json:"translations"`
}

type Translation struct {
    TranslatedText         string `json:"translatedText"`
    DetectedSourceLanguage string `json:"detectedSourceLanguage,omitempty"`
}

// DetectResponse represents language detection response
type DetectResponse struct {
    Data DetectData `json:"data"`
}

type DetectData struct {
    Detections [][]Detection `json:"detections"`
}

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
```

### 1.7 OpenAI Client

```go
// infrastructure/external/openai/ai_client.go

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

type AIClient struct {
    apiKey     string
    model      string
    httpClient *http.Client
}

func NewAIClient(apiKey, model string) *AIClient {
    if model == "" {
        model = "gpt-4-turbo-preview"
    }
    return &AIClient{
        apiKey: apiKey,
        model:  model,
        httpClient: &http.Client{
            Timeout: 60 * time.Second, // AI responses can take longer
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
    ID      string `json:"id"`
    Object  string `json:"object"`
    Created int64  `json:"created"`
    Model   string `json:"model"`
    Choices []Choice `json:"choices"`
    Usage   Usage    `json:"usage"`
}

type Choice struct {
    Index        int         `json:"index"`
    Message      ChatMessage `json:"message"`
    FinishReason string      `json:"finish_reason"`
}

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

// GenerateTravelSummary generates a travel summary from search results
func (c *AIClient) GenerateTravelSummary(ctx context.Context, query string, searchResults []SearchResultContext) (*ChatResponse, error) {
    systemPrompt := `คุณเป็นผู้ช่วยค้นหาข้อมูลท่องเที่ยวสำหรับนักศึกษามหาวิทยาลัยสุโขทัยธรรมาธิราช (มสธ.)
ให้สรุปข้อมูลจาก sources ที่ได้รับอย่างกระชับ ชัดเจน และเป็นประโยชน์

กฎในการตอบ:
1. ตอบเป็นภาษาไทย
2. จัดรูปแบบเป็น Markdown
3. สรุปเป็นหัวข้อหลักๆ พร้อม bullet points
4. ระบุ source ที่มาของข้อมูล
5. หากมีข้อมูลราคา เวลาเปิด-ปิด หรือข้อมูลสำคัญ ให้ระบุด้วย
6. เสนอคำถาม follow-up ที่เกี่ยวข้อง 2-3 ข้อ`

    // Build user prompt with search results
    userPrompt := fmt.Sprintf("คำค้นหา: %s\n\nข้อมูลจากแหล่งต่างๆ:\n", query)
    for i, result := range searchResults {
        userPrompt += fmt.Sprintf("\n[Source %d: %s]\n%s\nURL: %s\n",
            i+1, result.Title, result.Snippet, result.URL)
    }
    userPrompt += "\nกรุณาสรุปข้อมูลข้างต้นอย่างเป็นระบบ"

    messages := []ChatMessage{
        {Role: "system", Content: systemPrompt},
        {Role: "user", Content: userPrompt},
    }

    return c.Chat(ctx, messages, 2000, 0.7)
}

// ContinueChat continues an existing chat conversation
func (c *AIClient) ContinueChat(ctx context.Context, history []ChatMessage, newMessage string) (*ChatResponse, error) {
    systemPrompt := `คุณเป็นผู้ช่วยค้นหาข้อมูลท่องเที่ยวสำหรับนักศึกษา มสธ.
ตอบคำถามเกี่ยวกับการท่องเที่ยวอย่างเป็นมิตรและให้ข้อมูลที่เป็นประโยชน์
ตอบเป็นภาษาไทยและใช้ Markdown format`

    messages := []ChatMessage{
        {Role: "system", Content: systemPrompt},
    }
    messages = append(messages, history...)
    messages = append(messages, ChatMessage{Role: "user", Content: newMessage})

    return c.Chat(ctx, messages, 1500, 0.7)
}

// SearchResultContext represents context from search results
type SearchResultContext struct {
    Title   string
    Snippet string
    URL     string
}
```

---

## 2. Repository Implementations (infrastructure/postgres/)

### 2.1 Folder Repository Implementation

```go
// infrastructure/postgres/folder_repository_impl.go

package postgres

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/models"
    "github.com/your-org/stou-smart-tour/domain/repositories"
    "gorm.io/gorm"
)

type folderRepository struct {
    db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) repositories.FolderRepository {
    return &folderRepository{db: db}
}

func (r *folderRepository) Create(ctx context.Context, folder *models.Folder) error {
    return r.db.WithContext(ctx).Create(folder).Error
}

func (r *folderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Folder, error) {
    var folder models.Folder
    err := r.db.WithContext(ctx).First(&folder, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &folder, nil
}

func (r *folderRepository) Update(ctx context.Context, folder *models.Folder) error {
    return r.db.WithContext(ctx).Save(folder).Error
}

func (r *folderRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&models.Folder{}, "id = ?", id).Error
}

func (r *folderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.Folder, error) {
    var folders []models.Folder
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("updated_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&folders).Error
    return folders, err
}

func (r *folderRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&models.Folder{}).
        Where("user_id = ?", userID).
        Count(&count).Error
    return count, err
}

func (r *folderRepository) GetPublicFolders(ctx context.Context, offset, limit int) ([]models.Folder, error) {
    var folders []models.Folder
    err := r.db.WithContext(ctx).
        Where("is_public = ?", true).
        Order("updated_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&folders).Error
    return folders, err
}

func (r *folderRepository) CountPublicFolders(ctx context.Context) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&models.Folder{}).
        Where("is_public = ?", true).
        Count(&count).Error
    return count, err
}

func (r *folderRepository) IncrementItemCount(ctx context.Context, folderID uuid.UUID) error {
    return r.db.WithContext(ctx).
        Model(&models.Folder{}).
        Where("id = ?", folderID).
        UpdateColumn("item_count", gorm.Expr("item_count + 1")).Error
}

func (r *folderRepository) DecrementItemCount(ctx context.Context, folderID uuid.UUID) error {
    return r.db.WithContext(ctx).
        Model(&models.Folder{}).
        Where("id = ? AND item_count > 0", folderID).
        UpdateColumn("item_count", gorm.Expr("item_count - 1")).Error
}

func (r *folderRepository) UpdateCoverImage(ctx context.Context, folderID uuid.UUID, imageURL string) error {
    return r.db.WithContext(ctx).
        Model(&models.Folder{}).
        Where("id = ?", folderID).
        Update("cover_image_url", imageURL).Error
}
```

### 2.2 Folder Item Repository Implementation

```go
// infrastructure/postgres/folder_item_repository_impl.go

package postgres

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/models"
    "github.com/your-org/stou-smart-tour/domain/repositories"
    "gorm.io/gorm"
)

type folderItemRepository struct {
    db *gorm.DB
}

func NewFolderItemRepository(db *gorm.DB) repositories.FolderItemRepository {
    return &folderItemRepository{db: db}
}

func (r *folderItemRepository) Create(ctx context.Context, item *models.FolderItem) error {
    return r.db.WithContext(ctx).Create(item).Error
}

func (r *folderItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.FolderItem, error) {
    var item models.FolderItem
    err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &item, nil
}

func (r *folderItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&models.FolderItem{}, "id = ?", id).Error
}

func (r *folderItemRepository) GetByFolderID(ctx context.Context, folderID uuid.UUID, offset, limit int) ([]models.FolderItem, error) {
    var items []models.FolderItem
    err := r.db.WithContext(ctx).
        Where("folder_id = ?", folderID).
        Order("sort_order ASC, created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&items).Error
    return items, err
}

func (r *folderItemRepository) CountByFolderID(ctx context.Context, folderID uuid.UUID) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&models.FolderItem{}).
        Where("folder_id = ?", folderID).
        Count(&count).Error
    return count, err
}

func (r *folderItemRepository) GetFirstImageByFolderID(ctx context.Context, folderID uuid.UUID) (*models.FolderItem, error) {
    var item models.FolderItem
    err := r.db.WithContext(ctx).
        Where("folder_id = ? AND thumbnail_url IS NOT NULL AND thumbnail_url != ''", folderID).
        Order("created_at ASC").
        First(&item).Error
    if err != nil {
        return nil, err
    }
    return &item, nil
}

func (r *folderItemRepository) DeleteByFolderID(ctx context.Context, folderID uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&models.FolderItem{}, "folder_id = ?", folderID).Error
}

func (r *folderItemRepository) UpdateSortOrder(ctx context.Context, id uuid.UUID, sortOrder int) error {
    return r.db.WithContext(ctx).
        Model(&models.FolderItem{}).
        Where("id = ?", id).
        Update("sort_order", sortOrder).Error
}
```

### 2.3 Favorite Repository Implementation

```go
// infrastructure/postgres/favorite_repository_impl.go

package postgres

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/models"
    "github.com/your-org/stou-smart-tour/domain/repositories"
    "gorm.io/gorm"
)

type favoriteRepository struct {
    db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) repositories.FavoriteRepository {
    return &favoriteRepository{db: db}
}

func (r *favoriteRepository) Create(ctx context.Context, favorite *models.Favorite) error {
    return r.db.WithContext(ctx).Create(favorite).Error
}

func (r *favoriteRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Favorite, error) {
    var favorite models.Favorite
    err := r.db.WithContext(ctx).First(&favorite, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &favorite, nil
}

func (r *favoriteRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&models.Favorite{}, "id = ?", id).Error
}

func (r *favoriteRepository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.Favorite, error) {
    var favorites []models.Favorite
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&favorites).Error
    return favorites, err
}

func (r *favoriteRepository) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, favType string, offset, limit int) ([]models.Favorite, error) {
    var favorites []models.Favorite
    err := r.db.WithContext(ctx).
        Where("user_id = ? AND type = ?", userID, favType).
        Order("created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&favorites).Error
    return favorites, err
}

func (r *favoriteRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&models.Favorite{}).
        Where("user_id = ?", userID).
        Count(&count).Error
    return count, err
}

func (r *favoriteRepository) CountByUserIDAndType(ctx context.Context, userID uuid.UUID, favType string) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&models.Favorite{}).
        Where("user_id = ? AND type = ?", userID, favType).
        Count(&count).Error
    return count, err
}

func (r *favoriteRepository) ExistsByUserAndExternal(ctx context.Context, userID uuid.UUID, favType, externalID string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&models.Favorite{}).
        Where("user_id = ? AND type = ? AND external_id = ?", userID, favType, externalID).
        Count(&count).Error
    return count > 0, err
}

func (r *favoriteRepository) GetByUserAndExternal(ctx context.Context, userID uuid.UUID, favType, externalID string) (*models.Favorite, error) {
    var favorite models.Favorite
    err := r.db.WithContext(ctx).
        Where("user_id = ? AND type = ? AND external_id = ?", userID, favType, externalID).
        First(&favorite).Error
    if err != nil {
        return nil, err
    }
    return &favorite, nil
}
```

### 2.4 Search History Repository Implementation

```go
// infrastructure/postgres/search_history_repository_impl.go

package postgres

import (
    "context"
    "time"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/models"
    "github.com/your-org/stou-smart-tour/domain/repositories"
    "gorm.io/gorm"
)

type searchHistoryRepository struct {
    db *gorm.DB
}

func NewSearchHistoryRepository(db *gorm.DB) repositories.SearchHistoryRepository {
    return &searchHistoryRepository{db: db}
}

func (r *searchHistoryRepository) Create(ctx context.Context, history *models.SearchHistory) error {
    return r.db.WithContext(ctx).Create(history).Error
}

func (r *searchHistoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.SearchHistory, error) {
    var histories []models.SearchHistory
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&histories).Error
    return histories, err
}

func (r *searchHistoryRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&models.SearchHistory{}).
        Where("user_id = ?", userID).
        Count(&count).Error
    return count, err
}

func (r *searchHistoryRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
    return r.db.WithContext(ctx).
        Delete(&models.SearchHistory{}, "user_id = ?", userID).Error
}

func (r *searchHistoryRepository) DeleteOlderThan(ctx context.Context, userID uuid.UUID, days int) error {
    cutoff := time.Now().AddDate(0, 0, -days)
    return r.db.WithContext(ctx).
        Delete(&models.SearchHistory{}, "user_id = ? AND created_at < ?", userID, cutoff).Error
}

func (r *searchHistoryRepository) GetPopularQueries(ctx context.Context, limit int) ([]string, error) {
    var results []struct {
        Query string
        Count int64
    }

    err := r.db.WithContext(ctx).
        Model(&models.SearchHistory{}).
        Select("query, COUNT(*) as count").
        Group("query").
        Order("count DESC").
        Limit(limit).
        Scan(&results).Error

    if err != nil {
        return nil, err
    }

    queries := make([]string, len(results))
    for i, r := range results {
        queries[i] = r.Query
    }
    return queries, nil
}
```

### 2.5 AI Chat Repository Implementation

```go
// infrastructure/postgres/ai_chat_repository_impl.go

package postgres

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/models"
    "github.com/your-org/stou-smart-tour/domain/repositories"
    "gorm.io/gorm"
)

type aiChatRepository struct {
    db *gorm.DB
}

func NewAIChatRepository(db *gorm.DB) repositories.AIChatRepository {
    return &aiChatRepository{db: db}
}

// Session operations
func (r *aiChatRepository) CreateSession(ctx context.Context, session *models.AIChatSession) error {
    return r.db.WithContext(ctx).Create(session).Error
}

func (r *aiChatRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (*models.AIChatSession, error) {
    var session models.AIChatSession
    err := r.db.WithContext(ctx).
        Preload("Messages", func(db *gorm.DB) *gorm.DB {
            return db.Order("created_at ASC")
        }).
        First(&session, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &session, nil
}

func (r *aiChatRepository) UpdateSession(ctx context.Context, session *models.AIChatSession) error {
    return r.db.WithContext(ctx).Save(session).Error
}

func (r *aiChatRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
    // Delete messages first (cascade)
    if err := r.db.WithContext(ctx).Delete(&models.AIChatMessage{}, "session_id = ?", id).Error; err != nil {
        return err
    }
    return r.db.WithContext(ctx).Delete(&models.AIChatSession{}, "id = ?", id).Error
}

func (r *aiChatRepository) GetSessionsByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.AIChatSession, error) {
    var sessions []models.AIChatSession
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("updated_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&sessions).Error
    return sessions, err
}

func (r *aiChatRepository) CountSessionsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&models.AIChatSession{}).
        Where("user_id = ?", userID).
        Count(&count).Error
    return count, err
}

// Message operations
func (r *aiChatRepository) CreateMessage(ctx context.Context, message *models.AIChatMessage) error {
    return r.db.WithContext(ctx).Create(message).Error
}

func (r *aiChatRepository) GetMessagesBySessionID(ctx context.Context, sessionID uuid.UUID) ([]models.AIChatMessage, error) {
    var messages []models.AIChatMessage
    err := r.db.WithContext(ctx).
        Where("session_id = ?", sessionID).
        Order("created_at ASC").
        Find(&messages).Error
    return messages, err
}

func (r *aiChatRepository) DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error {
    // Get all session IDs first
    var sessionIDs []uuid.UUID
    if err := r.db.WithContext(ctx).
        Model(&models.AIChatSession{}).
        Where("user_id = ?", userID).
        Pluck("id", &sessionIDs).Error; err != nil {
        return err
    }

    if len(sessionIDs) == 0 {
        return nil
    }

    // Delete all messages
    if err := r.db.WithContext(ctx).Delete(&models.AIChatMessage{}, "session_id IN ?", sessionIDs).Error; err != nil {
        return err
    }

    // Delete all sessions
    return r.db.WithContext(ctx).Delete(&models.AIChatSession{}, "user_id = ?", userID).Error
}
```

---

## 3. Cache Layer Enhancement (infrastructure/cache/)

### 3.1 Cache Keys

```go
// infrastructure/cache/cache_keys.go

package cache

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "time"
)

// Cache key prefixes
const (
    PrefixSearch       = "search"
    PrefixSearchAI     = "search:ai"
    PrefixPlace        = "place"
    PrefixPlaceDetails = "place:details"
    PrefixNearbyPlaces = "places:nearby"
    PrefixYouTube      = "youtube"
    PrefixTranslate    = "translate"
    PrefixUserSession  = "user:session"
)

// Cache TTLs
const (
    TTLSearch       = 1 * time.Hour
    TTLSearchAI     = 6 * time.Hour
    TTLPlace        = 1 * time.Hour
    TTLPlaceDetails = 24 * time.Hour
    TTLNearbyPlaces = 1 * time.Hour
    TTLYouTube      = 6 * time.Hour
    TTLTranslate    = 7 * 24 * time.Hour
    TTLUserSession  = 24 * time.Hour
)

// Hash helper
func hashString(s string) string {
    hash := md5.Sum([]byte(s))
    return hex.EncodeToString(hash[:])
}

// Key generators
func SearchKey(query, searchType string, page int) string {
    return fmt.Sprintf("%s:%s:%d", PrefixSearch, hashString(query+":"+searchType), page)
}

func SearchAIKey(query string) string {
    return fmt.Sprintf("%s:%s", PrefixSearchAI, hashString(query))
}

func PlaceKey(placeID string) string {
    return fmt.Sprintf("%s:%s", PrefixPlace, placeID)
}

func PlaceDetailsKey(placeID string) string {
    return fmt.Sprintf("%s:%s", PrefixPlaceDetails, placeID)
}

func NearbyPlacesKey(lat, lng float64, radius int, placeType, keyword string) string {
    key := fmt.Sprintf("%f:%f:%d:%s:%s", lat, lng, radius, placeType, keyword)
    return fmt.Sprintf("%s:%s", PrefixNearbyPlaces, hashString(key))
}

func YouTubeKey(query string, limit int) string {
    key := fmt.Sprintf("%s:%d", query, limit)
    return fmt.Sprintf("%s:%s", PrefixYouTube, hashString(key))
}

func TranslateKey(text, sourceLang, targetLang string) string {
    key := fmt.Sprintf("%s:%s:%s", text, sourceLang, targetLang)
    return fmt.Sprintf("%s:%s", PrefixTranslate, hashString(key))
}
```

### 3.2 Enhanced Redis Client

```go
// infrastructure/redis/redis.go
// เพิ่มเติมจากเดิม

package redis

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"
)

type RedisClient struct {
    client *redis.Client
}

func NewRedisClient(url string) (*RedisClient, error) {
    opt, err := redis.ParseURL(url)
    if err != nil {
        return nil, err
    }

    client := redis.NewClient(opt)

    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.Ping(ctx).Err(); err != nil {
        return nil, err
    }

    return &RedisClient{client: client}, nil
}

func (r *RedisClient) Close() error {
    return r.client.Close()
}

// Get retrieves a value and unmarshals it
func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
    val, err := r.client.Get(ctx, key).Result()
    if err != nil {
        return err
    }
    return json.Unmarshal([]byte(val), dest)
}

// Set marshals and stores a value with TTL
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return r.client.Set(ctx, key, data, ttl).Err()
}

// Exists checks if a key exists
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
    result, err := r.client.Exists(ctx, key).Result()
    return result > 0, err
}

// Delete removes a key
func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
    return r.client.Del(ctx, keys...).Err()
}

// GetOrSet gets from cache or sets using the loader function
func (r *RedisClient) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, loader func() (interface{}, error)) error {
    // Try to get from cache first
    err := r.Get(ctx, key, dest)
    if err == nil {
        return nil
    }

    // If not found, load and cache
    if err == redis.Nil {
        value, err := loader()
        if err != nil {
            return err
        }

        // Cache the result
        if err := r.Set(ctx, key, value, ttl); err != nil {
            // Log but don't fail
        }

        // Marshal and unmarshal to dest
        data, _ := json.Marshal(value)
        return json.Unmarshal(data, dest)
    }

    return err
}

// Increment increments a counter
func (r *RedisClient) Increment(ctx context.Context, key string) (int64, error) {
    return r.client.Incr(ctx, key).Result()
}

// SetExpire sets expiration on a key
func (r *RedisClient) SetExpire(ctx context.Context, key string, ttl time.Duration) error {
    return r.client.Expire(ctx, key, ttl).Err()
}

// DeletePattern deletes keys matching a pattern
func (r *RedisClient) DeletePattern(ctx context.Context, pattern string) error {
    var cursor uint64
    for {
        var keys []string
        var err error
        keys, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
        if err != nil {
            return err
        }

        if len(keys) > 0 {
            if err := r.client.Del(ctx, keys...).Err(); err != nil {
                return err
            }
        }

        if cursor == 0 {
            break
        }
    }
    return nil
}
```

---

## 4. Database Migrations

### 4.1 Migration Files

```sql
-- migrations/000005_add_student_id.up.sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS student_id VARCHAR(20) UNIQUE;
CREATE INDEX IF NOT EXISTS idx_users_student_id ON users(student_id);

-- migrations/000005_add_student_id.down.sql
DROP INDEX IF EXISTS idx_users_student_id;
ALTER TABLE users DROP COLUMN IF EXISTS student_id;
```

```sql
-- migrations/000006_create_folders.up.sql
CREATE TABLE IF NOT EXISTS folders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    cover_image_url TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    item_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_folders_user_id ON folders(user_id);
CREATE INDEX idx_folders_is_public ON folders(is_public);

-- migrations/000006_create_folders.down.sql
DROP TABLE IF EXISTS folders;
```

```sql
-- migrations/000007_create_folder_items.up.sql
CREATE TABLE IF NOT EXISTS folder_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    folder_id UUID NOT NULL REFERENCES folders(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    thumbnail_url TEXT,
    description TEXT,
    metadata JSONB DEFAULT '{}',
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_folder_items_folder_id ON folder_items(folder_id);
CREATE INDEX idx_folder_items_type ON folder_items(type);

-- migrations/000007_create_folder_items.down.sql
DROP TABLE IF EXISTS folder_items;
```

```sql
-- migrations/000008_create_favorites.up.sql
CREATE TABLE IF NOT EXISTS favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    external_id VARCHAR(255),
    title VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    thumbnail_url TEXT,
    rating DECIMAL(2,1),
    review_count INTEGER DEFAULT 0,
    address TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, type, external_id)
);

CREATE INDEX idx_favorites_user_id ON favorites(user_id);
CREATE INDEX idx_favorites_type ON favorites(type);

-- migrations/000008_create_favorites.down.sql
DROP TABLE IF EXISTS favorites;
```

```sql
-- migrations/000009_create_search_history.up.sql
CREATE TABLE IF NOT EXISTS search_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    query VARCHAR(500) NOT NULL,
    search_type VARCHAR(50) NOT NULL DEFAULT 'all',
    result_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_search_history_user_id ON search_history(user_id);
CREATE INDEX idx_search_history_created_at ON search_history(created_at DESC);

-- migrations/000009_create_search_history.down.sql
DROP TABLE IF EXISTS search_history;
```

```sql
-- migrations/000010_create_ai_chat_sessions.up.sql
CREATE TABLE IF NOT EXISTS ai_chat_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255),
    initial_query VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ai_chat_sessions_user_id ON ai_chat_sessions(user_id);

-- migrations/000010_create_ai_chat_sessions.down.sql
DROP TABLE IF EXISTS ai_chat_sessions;
```

```sql
-- migrations/000011_create_ai_chat_messages.up.sql
CREATE TABLE IF NOT EXISTS ai_chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES ai_chat_sessions(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    sources JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ai_chat_messages_session_id ON ai_chat_messages(session_id);

-- migrations/000011_create_ai_chat_messages.down.sql
DROP TABLE IF EXISTS ai_chat_messages;
```

---

## 5. Update Database Migration Function

```go
// infrastructure/postgres/database.go
// อัพเดท Migrate function

func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        // Existing models
        &models.User{},
        &models.Task{},
        &models.File{},
        &models.Job{},
        // New models
        &models.Folder{},
        &models.FolderItem{},
        &models.Favorite{},
        &models.SearchHistory{},
        &models.AIChatSession{},
        &models.AIChatMessage{},
    )
}
```

---

## 6. Summary - Infrastructure Layer Files

```
infrastructure/
├── postgres/
│   ├── database.go                        # ✏️ UPDATE (add new models to Migrate)
│   ├── user_repository_impl.go            # ✅ KEEP
│   ├── task_repository_impl.go            # ✅ KEEP
│   ├── file_repository_impl.go            # ✅ KEEP
│   ├── job_repository_impl.go             # ✅ KEEP
│   ├── folder_repository_impl.go          # 🆕 NEW
│   ├── folder_item_repository_impl.go     # 🆕 NEW
│   ├── favorite_repository_impl.go        # 🆕 NEW
│   ├── search_history_repository_impl.go  # 🆕 NEW
│   └── ai_chat_repository_impl.go         # 🆕 NEW
│
├── redis/
│   └── redis.go                           # ✏️ UPDATE (enhanced methods)
│
├── storage/
│   └── bunny_storage.go                   # ✅ KEEP
│
├── cache/                                 # 🆕 NEW FOLDER
│   └── cache_keys.go                      # 🆕 NEW
│
├── external/                              # 🆕 NEW FOLDER
│   ├── google/
│   │   ├── client.go                      # 🆕 NEW
│   │   ├── search_client.go               # 🆕 NEW
│   │   ├── places_client.go               # 🆕 NEW
│   │   ├── youtube_client.go              # 🆕 NEW
│   │   └── translate_client.go            # 🆕 NEW
│   │
│   └── openai/
│       └── ai_client.go                   # 🆕 NEW
│
└── websocket/
    └── websocket.go                       # ✅ KEEP
```

---

## Next Part

➡️ ไปต่อที่ **Part 4: Application Layer (Services Implementation)**
- Search Service Implementation
- AI Service Implementation
- Folder Service Implementation
- Favorite Service Implementation
- Utility Services (Translate, QRCode)

---

*Document Version: 1.0*
*Part: 3 of 5*
