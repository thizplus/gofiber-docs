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

// hashString creates MD5 hash of a string
func hashString(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// SearchKey generates cache key for search results
func SearchKey(query, searchType string, page int) string {
	return fmt.Sprintf("%s:%s:%d", PrefixSearch, hashString(query+":"+searchType), page)
}

// SearchAIKey generates cache key for AI search results
func SearchAIKey(query string) string {
	return fmt.Sprintf("%s:%s", PrefixSearchAI, hashString(query))
}

// PlaceKey generates cache key for place basic info
func PlaceKey(placeID string) string {
	return fmt.Sprintf("%s:%s", PrefixPlace, placeID)
}

// PlaceDetailsKey generates cache key for place details
func PlaceDetailsKey(placeID, lang string) string {
	if lang == "" {
		lang = "th"
	}
	return fmt.Sprintf("%s:%s:%s", PrefixPlaceDetails, placeID, lang)
}

// NearbyPlacesKey generates cache key for nearby places search
func NearbyPlacesKey(lat, lng float64, radius int, placeType, keyword, lang string) string {
	if lang == "" {
		lang = "th"
	}
	key := fmt.Sprintf("%f:%f:%d:%s:%s:%s", lat, lng, radius, placeType, keyword, lang)
	return fmt.Sprintf("%s:%s", PrefixNearbyPlaces, hashString(key))
}

// PlaceTextSearchKey generates cache key for text search places
func PlaceTextSearchKey(query, lang string) string {
	if lang == "" {
		lang = "th"
	}
	return fmt.Sprintf("%s:text:%s:%s", PrefixPlace, hashString(query), lang)
}

// YouTubeKey generates cache key for YouTube search results
func YouTubeKey(query string, limit int) string {
	key := fmt.Sprintf("%s:%d", query, limit)
	return fmt.Sprintf("%s:%s", PrefixYouTube, hashString(key))
}

// TranslateKey generates cache key for translation
func TranslateKey(text, sourceLang, targetLang string) string {
	key := fmt.Sprintf("%s:%s:%s", text, sourceLang, targetLang)
	return fmt.Sprintf("%s:%s", PrefixTranslate, hashString(key))
}

// DetectLanguageKey generates cache key for language detection
func DetectLanguageKey(text string) string {
	return fmt.Sprintf("detect:%s", hashString(text))
}

// VideoDetailsKey generates cache key for video details
func VideoDetailsKey(videoID string) string {
	return fmt.Sprintf("%s:details:%s", PrefixYouTube, videoID)
}

// ImageSearchKey generates cache key for image search
func ImageSearchKey(query string, page int) string {
	return fmt.Sprintf("%s:image:%s:%d", PrefixSearch, hashString(query), page)
}
