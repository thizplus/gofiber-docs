package google

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	youtubeSearchURL = "https://www.googleapis.com/youtube/v3/search"
	youtubeVideosURL = "https://www.googleapis.com/youtube/v3/videos"
)

// YouTubeClient handles YouTube Data API
type YouTubeClient struct {
	*GoogleClient
}

// NewYouTubeClient creates a new YouTube client
func NewYouTubeClient(apiKey string) *YouTubeClient {
	return &YouTubeClient{
		GoogleClient: NewGoogleClient(apiKey),
	}
}

// VideoSearchRequest represents YouTube search parameters
type VideoSearchRequest struct {
	Query      string
	MaxResults int
	Order      string // relevance, date, rating, viewCount
	Language   string
	RegionCode string
}

// VideoSearchResponse represents YouTube search API response
type VideoSearchResponse struct {
	Kind          string       `json:"kind"`
	PageInfo      PageInfo     `json:"pageInfo"`
	NextPageToken string       `json:"nextPageToken,omitempty"`
	Items         []VideoSearchItem `json:"items"`
}

// PageInfo contains pagination info
type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

// VideoSearchItem represents a single video search result
type VideoSearchItem struct {
	Kind    string       `json:"kind"`
	ID      VideoID      `json:"id"`
	Snippet VideoSnippet `json:"snippet"`
}

// VideoID contains video identification
type VideoID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

// VideoSnippet contains video metadata
type VideoSnippet struct {
	PublishedAt          string     `json:"publishedAt"`
	ChannelID            string     `json:"channelId"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	ChannelTitle         string     `json:"channelTitle"`
	Thumbnails           Thumbnails `json:"thumbnails"`
	LiveBroadcastContent string     `json:"liveBroadcastContent"`
}

// Thumbnails contains thumbnail images
type Thumbnails struct {
	Default  Thumbnail `json:"default"`
	Medium   Thumbnail `json:"medium"`
	High     Thumbnail `json:"high"`
	Standard Thumbnail `json:"standard,omitempty"`
	Maxres   Thumbnail `json:"maxres,omitempty"`
}

// Thumbnail represents a single thumbnail
type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// VideoDetailsResponse for video statistics and content details
type VideoDetailsResponse struct {
	Items []VideoDetails `json:"items"`
}

// VideoDetails contains detailed video info
type VideoDetails struct {
	ID             string         `json:"id"`
	ContentDetails ContentDetails `json:"contentDetails"`
	Statistics     Statistics     `json:"statistics"`
}

// ContentDetails contains video duration
type ContentDetails struct {
	Duration string `json:"duration"`
}

// Statistics contains video statistics
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
// e.g., PT1H2M3S -> "1:02:03", PT2M3S -> "2:03"
func ParseDuration(isoDuration string) string {
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
