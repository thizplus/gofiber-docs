package google

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
)

const (
	placesNearbyURL     = "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	placesTextSearchURL = "https://maps.googleapis.com/maps/api/place/textsearch/json"
	placeDetailsURL     = "https://maps.googleapis.com/maps/api/place/details/json"
	placePhotoURL       = "https://maps.googleapis.com/maps/api/place/photo"
)

// PlacesClient handles Google Places API
type PlacesClient struct {
	*GoogleClient
}

// NewPlacesClient creates a new Places client
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

// TextSearchRequest represents text search parameters
type TextSearchRequest struct {
	Query    string
	Language string
	Region   string // e.g., "th" for Thailand
}

// NearbySearchResponse represents Places API nearby search response
type NearbySearchResponse struct {
	Status           string   `json:"status"`
	Results          []Place  `json:"results"`
	NextPageToken    string   `json:"next_page_token,omitempty"`
	HTMLAttributions []string `json:"html_attributions"`
	ErrorMessage     string   `json:"error_message,omitempty"`
}

// Place represents a place result
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
	PriceLevel     int    `json:"price_level,omitempty"`
	BusinessStatus string `json:"business_status,omitempty"`
}

// Geometry contains location data
type Geometry struct {
	Location struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"location"`
}

// Photo contains photo reference data
type Photo struct {
	PhotoReference   string   `json:"photo_reference"`
	Height           int      `json:"height"`
	Width            int      `json:"width"`
	HTMLAttributions []string `json:"html_attributions"`
}

// PlaceDetailsRequest represents place details request
type PlaceDetailsRequest struct {
	PlaceID  string
	Fields   []string
	Language string
}

// PlaceDetailsResponse represents place details response
type PlaceDetailsResponse struct {
	Status       string       `json:"status"`
	Result       PlaceDetails `json:"result"`
	ErrorMessage string       `json:"error_message,omitempty"`
}

// PlaceDetails contains detailed place information
type PlaceDetails struct {
	PlaceID              string             `json:"place_id"`
	Name                 string             `json:"name"`
	FormattedAddress     string             `json:"formatted_address"`
	FormattedPhoneNumber string             `json:"formatted_phone_number,omitempty"`
	InternationalPhone   string             `json:"international_phone_number,omitempty"`
	Website              string             `json:"website,omitempty"`
	Rating               float64            `json:"rating"`
	UserRatingsTotal     int                `json:"user_ratings_total"`
	PriceLevel           int                `json:"price_level,omitempty"`
	Types                []string           `json:"types"`
	Geometry             Geometry           `json:"geometry"`
	Photos               []Photo            `json:"photos,omitempty"`
	Reviews              []Review           `json:"reviews,omitempty"`
	OpeningHours         *OpeningHoursDetail `json:"opening_hours,omitempty"`
	URL                  string             `json:"url"`
}

// Review represents a place review
type Review struct {
	AuthorName              string `json:"author_name"`
	AuthorURL               string `json:"author_url,omitempty"`
	ProfilePhotoURL         string `json:"profile_photo_url,omitempty"`
	Rating                  int    `json:"rating"`
	Text                    string `json:"text"`
	Time                    int64  `json:"time"`
	RelativeTimeDescription string `json:"relative_time_description"`
}

// OpeningHoursDetail contains detailed opening hours
type OpeningHoursDetail struct {
	OpenNow     bool     `json:"open_now"`
	WeekdayText []string `json:"weekday_text"`
}

// TextSearch searches for places using text query (like Google Maps search)
func (c *PlacesClient) TextSearch(ctx context.Context, req *TextSearchRequest) (*NearbySearchResponse, error) {
	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("query", req.Query)

	if req.Language != "" {
		params.Set("language", req.Language)
	} else {
		params.Set("language", "th")
	}

	if req.Region != "" {
		params.Set("region", req.Region)
	}

	searchURL := fmt.Sprintf("%s?%s", placesTextSearchURL, params.Encode())

	var result NearbySearchResponse
	if err := c.doRequest(ctx, searchURL, &result); err != nil {
		return nil, err
	}

	if result.Status != "OK" && result.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("Places API error: %s - %s", result.Status, result.ErrorMessage)
	}

	return &result, nil
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

// CalculateDistance calculates distance between two points using Haversine formula
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000 // Earth's radius in meters

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
