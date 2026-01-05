package serviceimpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/cache"
	"gofiber-template/infrastructure/external/google"
	"gofiber-template/infrastructure/external/openai"
)

type SearchServiceImpl struct {
	searchHistoryRepo    repositories.SearchHistoryRepository
	placeAIContentRepo   repositories.PlaceAIContentRepository
	googleSearch         *google.SearchClient
	googlePlaces         *google.PlacesClient
	googleYouTube        *google.YouTubeClient
	openaiClient         *openai.AIClient
	redisClient          *redis.Client
}

func NewSearchService(
	searchHistoryRepo repositories.SearchHistoryRepository,
	placeAIContentRepo repositories.PlaceAIContentRepository,
	googleSearch *google.SearchClient,
	googlePlaces *google.PlacesClient,
	googleYouTube *google.YouTubeClient,
	openaiClient *openai.AIClient,
	redisClient *redis.Client,
) services.SearchService {
	return &SearchServiceImpl{
		searchHistoryRepo:  searchHistoryRepo,
		placeAIContentRepo: placeAIContentRepo,
		googleSearch:       googleSearch,
		googlePlaces:       googlePlaces,
		googleYouTube:      googleYouTube,
		openaiClient:       openaiClient,
		redisClient:        redisClient,
	}
}

func (s *SearchServiceImpl) Search(ctx context.Context, userID uuid.UUID, req *dto.SearchRequest) (*dto.SearchResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.Type == "" {
		req.Type = "all"
	}

	var results []dto.SearchResult
	var totalCount int64

	switch req.Type {
	case "website":
		webResp, err := s.SearchWebsites(ctx, userID, req)
		if err != nil {
			return nil, err
		}
		for _, r := range webResp.Results {
			results = append(results, dto.SearchResult{
				Type:    "website",
				Title:   r.Title,
				URL:     r.URL,
				Snippet: r.Snippet,
				Source:  r.DisplayLink,
			})
		}
		totalCount = webResp.TotalCount
	case "image":
		imgReq := &dto.ImageSearchRequest{
			Query:    req.Query,
			Page:     req.Page,
			PageSize: req.PageSize,
		}
		imgResp, err := s.SearchImages(ctx, userID, imgReq)
		if err != nil {
			return nil, err
		}
		for _, r := range imgResp.Results {
			results = append(results, dto.SearchResult{
				Type:         "image",
				Title:        r.Title,
				URL:          r.URL,
				ThumbnailURL: r.ThumbnailURL,
				Source:       r.Source,
			})
		}
		totalCount = imgResp.TotalCount
	case "video":
		vidReq := &dto.VideoSearchRequest{
			Query:    req.Query,
			Page:     req.Page,
			PageSize: req.PageSize,
		}
		vidResp, err := s.SearchVideos(ctx, userID, vidReq)
		if err != nil {
			return nil, err
		}
		for _, r := range vidResp.Results {
			results = append(results, dto.SearchResult{
				Type:         "video",
				Title:        r.Title,
				URL:          "https://www.youtube.com/watch?v=" + r.VideoID,
				Snippet:      r.Description,
				ThumbnailURL: r.ThumbnailURL,
				Source:       r.ChannelTitle,
				PublishedAt:  r.PublishedAt,
			})
		}
		totalCount = vidResp.TotalCount
	default:
		// All - search multiple types in parallel
		// Search places
		placeReq := &dto.PlaceSearchRequest{
			Query:    req.Query,
			Page:     1,
			PageSize: 8,
		}
		placeResp, _ := s.SearchPlaces(ctx, userID, placeReq)
		if placeResp != nil {
			for _, r := range placeResp.Results {
				results = append(results, dto.SearchResult{
					Type:         "place",
					PlaceID:      r.PlaceID,
					Title:        r.Name,
					Snippet:      r.Address,
					ThumbnailURL: r.PhotoURL,
					Lat:          r.Lat,
					Lng:          r.Lng,
					Rating:       r.Rating,
					ReviewCount:  r.ReviewCount,
					Types:        r.Types,
				})
			}
		}

		// Search videos
		vidReq := &dto.VideoSearchRequest{
			Query:    req.Query,
			Page:     1,
			PageSize: 4,
		}
		vidResp, _ := s.SearchVideos(ctx, userID, vidReq)
		if vidResp != nil {
			for _, r := range vidResp.Results {
				results = append(results, dto.SearchResult{
					Type:         "video",
					VideoID:      r.VideoID,
					Title:        r.Title,
					Snippet:      r.Description,
					ThumbnailURL: r.ThumbnailURL,
					Source:       r.ChannelTitle,
					PublishedAt:  r.PublishedAt,
					Duration:     r.Duration,
					ViewCount:    r.ViewCount,
				})
			}
		}

		// Search websites
		webResp, err := s.SearchWebsites(ctx, userID, req)
		if err != nil {
			return nil, err
		}
		for _, r := range webResp.Results {
			results = append(results, dto.SearchResult{
				Type:    "website",
				Title:   r.Title,
				URL:     r.URL,
				Snippet: r.Snippet,
				Source:  r.DisplayLink,
			})
		}
		totalCount = webResp.TotalCount
	}

	return &dto.SearchResponse{
		Query:      req.Query,
		Type:       req.Type,
		Results:    results,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

func (s *SearchServiceImpl) SearchWebsites(ctx context.Context, userID uuid.UUID, req *dto.SearchRequest) (*dto.WebsiteSearchResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// Check cache first
	cacheKey := cache.SearchKey(req.Query, "website", req.Page)
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedResult dto.WebsiteSearchResponse
		if json.Unmarshal([]byte(cached), &cachedResult) == nil {
			// Don't save history for cache hits - only first search counts
			return &cachedResult, nil
		}
	}

	// Cache miss - call API
	searchResponse, err := s.googleSearch.SearchAll(ctx, req.Query, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var websiteResults []dto.WebsiteResult
	for _, r := range searchResponse.Items {
		websiteResults = append(websiteResults, dto.WebsiteResult{
			Title:       r.Title,
			URL:         r.Link,
			Snippet:     r.Snippet,
			DisplayLink: r.DisplayLink,
		})
	}

	response := &dto.WebsiteSearchResponse{
		Query:      req.Query,
		Results:    websiteResults,
		TotalCount: int64(len(websiteResults)),
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	// Store in cache
	if jsonData, err := json.Marshal(response); err == nil {
		s.redisClient.Set(ctx, cacheKey, jsonData, cache.TTLSearch)
	}

	// Save search history
	s.saveSearchHistory(ctx, userID, req.Query, models.SearchTypeWebsite, len(websiteResults))

	return response, nil
}

func (s *SearchServiceImpl) SearchImages(ctx context.Context, userID uuid.UUID, req *dto.ImageSearchRequest) (*dto.ImageSearchResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// Check cache first
	cacheKey := cache.ImageSearchKey(req.Query, req.Page)
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedResult dto.ImageSearchResponse
		if json.Unmarshal([]byte(cached), &cachedResult) == nil {
			// Don't save history for cache hits - only first search counts
			return &cachedResult, nil
		}
	}

	// Cache miss - call API
	searchResponse, err := s.googleSearch.SearchImages(ctx, req.Query, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var imageResults []dto.ImageResult
	for _, r := range searchResponse.Items {
		thumbnailURL := ""
		width := 0
		height := 0
		contextLink := ""
		if r.Image != nil {
			thumbnailURL = r.Image.ThumbnailLink
			width = r.Image.Width
			height = r.Image.Height
			contextLink = r.Image.ContextLink
		}
		imageResults = append(imageResults, dto.ImageResult{
			Title:        r.Title,
			URL:          r.Link,
			ThumbnailURL: thumbnailURL,
			Width:        width,
			Height:       height,
			Source:       r.DisplayLink,
			ContextLink:  contextLink,
		})
	}

	response := &dto.ImageSearchResponse{
		Query:      req.Query,
		Results:    imageResults,
		TotalCount: int64(len(imageResults)),
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	// Store in cache (6 hours for images)
	if jsonData, err := json.Marshal(response); err == nil {
		s.redisClient.Set(ctx, cacheKey, jsonData, cache.TTLYouTube)
	}

	// Save search history
	s.saveSearchHistory(ctx, userID, req.Query, models.SearchTypeImage, len(imageResults))

	return response, nil
}

func (s *SearchServiceImpl) SearchVideos(ctx context.Context, userID uuid.UUID, req *dto.VideoSearchRequest) (*dto.VideoSearchResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.Order == "" {
		req.Order = "relevance"
	}

	// Check cache first
	cacheKey := cache.YouTubeKey(req.Query, req.PageSize)
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedResult dto.VideoSearchResponse
		if json.Unmarshal([]byte(cached), &cachedResult) == nil {
			// Don't save history for cache hits - only first search counts
			return &cachedResult, nil
		}
	}

	// Cache miss - call API
	searchReq := &google.VideoSearchRequest{
		Query:      req.Query,
		MaxResults: req.PageSize,
		Order:      req.Order,
	}

	searchResponse, err := s.googleYouTube.SearchVideos(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	// Collect video IDs for details lookup
	var videoIDs []string
	for _, item := range searchResponse.Items {
		videoIDs = append(videoIDs, item.ID.VideoID)
	}

	// Get video details (duration, view count)
	var detailsMap = make(map[string]google.VideoDetails)
	if len(videoIDs) > 0 {
		detailsResponse, err := s.googleYouTube.GetVideoDetails(ctx, videoIDs)
		if err == nil && detailsResponse != nil {
			for _, d := range detailsResponse.Items {
				detailsMap[d.ID] = d
			}
		}
	}

	var videoResults []dto.VideoResult
	for _, item := range searchResponse.Items {
		videoID := item.ID.VideoID
		duration := ""
		viewCount := ""
		if details, ok := detailsMap[videoID]; ok {
			duration = google.ParseDuration(details.ContentDetails.Duration)
			viewCount = details.Statistics.ViewCount
		}

		thumbnailURL := ""
		if item.Snippet.Thumbnails.High.URL != "" {
			thumbnailURL = item.Snippet.Thumbnails.High.URL
		} else if item.Snippet.Thumbnails.Medium.URL != "" {
			thumbnailURL = item.Snippet.Thumbnails.Medium.URL
		} else {
			thumbnailURL = item.Snippet.Thumbnails.Default.URL
		}

		viewCountInt, _ := strconv.ParseInt(viewCount, 10, 64)
		videoResults = append(videoResults, dto.VideoResult{
			VideoID:      videoID,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			ThumbnailURL: thumbnailURL,
			ChannelTitle: item.Snippet.ChannelTitle,
			PublishedAt:  item.Snippet.PublishedAt,
			Duration:     duration,
			ViewCount:    viewCountInt,
		})
	}

	response := &dto.VideoSearchResponse{
		Query:      req.Query,
		Results:    videoResults,
		TotalCount: int64(len(videoResults)),
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	// Store in cache (6 hours)
	if jsonData, err := json.Marshal(response); err == nil {
		s.redisClient.Set(ctx, cacheKey, jsonData, cache.TTLYouTube)
	}

	// Save search history
	s.saveSearchHistory(ctx, userID, req.Query, models.SearchTypeVideo, len(videoResults))

	return response, nil
}

func (s *SearchServiceImpl) GetVideoDetails(ctx context.Context, videoID string) (*dto.VideoResult, error) {
	// Check cache first
	cacheKey := cache.VideoDetailsKey(videoID)
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedResult dto.VideoResult
		if json.Unmarshal([]byte(cached), &cachedResult) == nil {
			return &cachedResult, nil
		}
	}

	// Cache miss - Get video details (duration, statistics)
	detailsResponse, err := s.googleYouTube.GetVideoDetails(ctx, []string{videoID})
	if err != nil {
		return nil, err
	}

	if len(detailsResponse.Items) == 0 {
		return nil, errors.New("video not found")
	}

	details := detailsResponse.Items[0]

	// Search for video to get snippet info (title, description, etc.)
	searchReq := &google.VideoSearchRequest{
		Query:      "video:" + videoID,
		MaxResults: 1,
	}
	searchResponse, _ := s.googleYouTube.SearchVideos(ctx, searchReq)

	var title, description, thumbnailURL, channelTitle, publishedAt string
	if searchResponse != nil && len(searchResponse.Items) > 0 {
		snippet := searchResponse.Items[0].Snippet
		title = snippet.Title
		description = snippet.Description
		channelTitle = snippet.ChannelTitle
		publishedAt = snippet.PublishedAt
		if snippet.Thumbnails.High.URL != "" {
			thumbnailURL = snippet.Thumbnails.High.URL
		} else {
			thumbnailURL = snippet.Thumbnails.Default.URL
		}
	}

	viewCountInt, _ := strconv.ParseInt(details.Statistics.ViewCount, 10, 64)
	likeCountInt, _ := strconv.ParseInt(details.Statistics.LikeCount, 10, 64)

	result := &dto.VideoResult{
		VideoID:      videoID,
		Title:        title,
		Description:  description,
		ThumbnailURL: thumbnailURL,
		ChannelTitle: channelTitle,
		PublishedAt:  publishedAt,
		Duration:     google.ParseDuration(details.ContentDetails.Duration),
		ViewCount:    viewCountInt,
		LikeCount:    likeCountInt,
	}

	// Store in cache (24 hours)
	if jsonData, err := json.Marshal(result); err == nil {
		s.redisClient.Set(ctx, cacheKey, jsonData, cache.TTLPlaceDetails)
	}

	return result, nil
}

func (s *SearchServiceImpl) SearchPlaces(ctx context.Context, userID uuid.UUID, req *dto.PlaceSearchRequest) (*dto.PlaceSearchResponse, error) {
	// Determine search type: Text Search (no lat/lng) or Nearby Search (with lat/lng)
	useTextSearch := req.Lat == 0 && req.Lng == 0

	var cacheKey string
	if useTextSearch {
		cacheKey = cache.PlaceTextSearchKey(req.Query)
	} else {
		if req.Radius == 0 {
			req.Radius = 5000
		}
		cacheKey = cache.NearbyPlacesKey(req.Lat, req.Lng, req.Radius, req.PlaceType, req.Query)
	}

	// Check cache first
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedResult dto.PlaceSearchResponse
		if json.Unmarshal([]byte(cached), &cachedResult) == nil {
			// Don't save history for cache hits - only first search counts
			return &cachedResult, nil
		}
	}

	// Cache miss - call API
	var searchResponse *google.NearbySearchResponse
	var err error

	if useTextSearch {
		// Text Search - search by query text only (like Google Maps search)
		textReq := &google.TextSearchRequest{
			Query:    req.Query,
			Language: "th",
			Region:   "th",
		}
		searchResponse, err = s.googlePlaces.TextSearch(ctx, textReq)
	} else {
		// Nearby Search - search by location
		nearbyReq := &google.NearbySearchRequest{
			Lat:     req.Lat,
			Lng:     req.Lng,
			Radius:  req.Radius,
			Type:    req.PlaceType,
			Keyword: req.Query,
		}
		searchResponse, err = s.googlePlaces.NearbySearch(ctx, nearbyReq)
	}

	if err != nil {
		return nil, err
	}

	var placeResults []dto.PlaceResult
	for _, place := range searchResponse.Results {
		photoURL := ""
		if len(place.Photos) > 0 {
			photoURL = s.googlePlaces.GetPhotoURL(place.Photos[0].PhotoReference, 400)
		}

		var isOpen *bool
		if place.OpeningHours != nil {
			isOpen = &place.OpeningHours.OpenNow
		}

		// Use FormattedAddress for text search, Vicinity for nearby search
		address := place.Vicinity
		if address == "" {
			address = place.FormattedAddress
		}

		placeResult := dto.PlaceResult{
			PlaceID:     place.PlaceID,
			Name:        place.Name,
			Address:     address,
			Lat:         place.Geometry.Location.Lat,
			Lng:         place.Geometry.Location.Lng,
			Rating:      place.Rating,
			ReviewCount: place.UserRatingsTotal,
			PriceLevel:  place.PriceLevel,
			Types:       place.Types,
			PhotoURL:    photoURL,
			IsOpen:      isOpen,
		}

		// Calculate distance if user location provided
		if req.Lat != 0 && req.Lng != 0 {
			distance := google.CalculateDistance(req.Lat, req.Lng, place.Geometry.Location.Lat, place.Geometry.Location.Lng)
			placeResult.Distance = distance
			placeResult.DistanceText = formatDistance(distance)
		}

		placeResults = append(placeResults, placeResult)
	}

	response := &dto.PlaceSearchResponse{
		Query:      req.Query,
		Results:    placeResults,
		TotalCount: int64(len(placeResults)),
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	// Store in cache
	if jsonData, err := json.Marshal(response); err == nil {
		s.redisClient.Set(ctx, cacheKey, jsonData, cache.TTLNearbyPlaces)
	}

	// Save search history
	s.saveSearchHistory(ctx, userID, req.Query, models.SearchTypeMap, len(placeResults))

	return response, nil
}

func (s *SearchServiceImpl) GetPlaceDetails(ctx context.Context, placeID string, userLat, userLng float64) (*dto.PlaceDetailResponse, error) {
	// Check cache first (without distance - distance calculated per user)
	cacheKey := cache.PlaceDetailsKey(placeID)
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedResult dto.PlaceDetailResponse
		if json.Unmarshal([]byte(cached), &cachedResult) == nil {
			// Calculate distance for this user
			if userLat != 0 && userLng != 0 {
				distance := google.CalculateDistance(userLat, userLng, cachedResult.Lat, cachedResult.Lng)
				cachedResult.Distance = distance
				cachedResult.DistanceText = formatDistance(distance)
			}
			return &cachedResult, nil
		}
	}

	// Cache miss - call API
	detailsReq := &google.PlaceDetailsRequest{
		PlaceID: placeID,
	}

	detailsResponse, err := s.googlePlaces.GetPlaceDetails(ctx, detailsReq)
	if err != nil {
		return nil, err
	}

	result := detailsResponse.Result

	var openingHours []string
	if result.OpeningHours != nil {
		openingHours = result.OpeningHours.WeekdayText
	}

	response := &dto.PlaceDetailResponse{
		PlaceID:          result.PlaceID,
		Name:             result.Name,
		FormattedAddress: result.FormattedAddress,
		Lat:              result.Geometry.Location.Lat,
		Lng:              result.Geometry.Location.Lng,
		Rating:           result.Rating,
		ReviewCount:      result.UserRatingsTotal,
		PriceLevel:       result.PriceLevel,
		Types:            result.Types,
		Phone:            result.FormattedPhoneNumber,
		Website:          result.Website,
		GoogleMapsURL:    result.URL,
		OpeningHours:     openingHours,
	}

	// Convert reviews
	for _, r := range result.Reviews {
		response.Reviews = append(response.Reviews, dto.PlaceReview{
			Author:   r.AuthorName,
			Rating:   r.Rating,
			Text:     r.Text,
			Time:     r.RelativeTimeDescription,
			PhotoURL: r.ProfilePhotoURL,
		})
	}

	// Convert photos
	for _, p := range result.Photos {
		photoURL := s.googlePlaces.GetPhotoURL(p.PhotoReference, 800)
		response.Photos = append(response.Photos, dto.PlacePhoto{
			URL:    photoURL,
			Width:  p.Width,
			Height: p.Height,
		})
	}

	// Store in cache (24 hours - place details rarely change)
	if jsonData, err := json.Marshal(response); err == nil {
		s.redisClient.Set(ctx, cacheKey, jsonData, cache.TTLPlaceDetails)
	}

	// Calculate distance if user location provided (after caching base data)
	if userLat != 0 && userLng != 0 {
		distance := google.CalculateDistance(userLat, userLng, result.Geometry.Location.Lat, result.Geometry.Location.Lng)
		response.Distance = distance
		response.DistanceText = formatDistance(distance)
	}

	return response, nil
}

func (s *SearchServiceImpl) SearchNearbyPlaces(ctx context.Context, req *dto.NearbyPlacesRequest) (*dto.PlaceSearchResponse, error) {
	if req.Radius == 0 {
		req.Radius = 5000
	}

	// Check cache first
	cacheKey := cache.NearbyPlacesKey(req.Lat, req.Lng, req.Radius, req.PlaceType, req.Keyword)
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedResult dto.PlaceSearchResponse
		if json.Unmarshal([]byte(cached), &cachedResult) == nil {
			return &cachedResult, nil
		}
	}

	// Cache miss - call API
	searchReq := &google.NearbySearchRequest{
		Lat:     req.Lat,
		Lng:     req.Lng,
		Radius:  req.Radius,
		Type:    req.PlaceType,
		Keyword: req.Keyword,
	}

	searchResponse, err := s.googlePlaces.NearbySearch(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	var placeResults []dto.PlaceResult
	for _, place := range searchResponse.Results {
		photoURL := ""
		if len(place.Photos) > 0 {
			photoURL = s.googlePlaces.GetPhotoURL(place.Photos[0].PhotoReference, 400)
		}

		var isOpen *bool
		if place.OpeningHours != nil {
			isOpen = &place.OpeningHours.OpenNow
		}

		distance := google.CalculateDistance(req.Lat, req.Lng, place.Geometry.Location.Lat, place.Geometry.Location.Lng)
		placeResults = append(placeResults, dto.PlaceResult{
			PlaceID:      place.PlaceID,
			Name:         place.Name,
			Address:      place.Vicinity,
			Lat:          place.Geometry.Location.Lat,
			Lng:          place.Geometry.Location.Lng,
			Rating:       place.Rating,
			ReviewCount:  place.UserRatingsTotal,
			PriceLevel:   place.PriceLevel,
			Types:        place.Types,
			PhotoURL:     photoURL,
			IsOpen:       isOpen,
			Distance:     distance,
			DistanceText: formatDistance(distance),
		})
	}

	response := &dto.PlaceSearchResponse{
		Query:      req.Keyword,
		Results:    placeResults,
		TotalCount: int64(len(placeResults)),
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	// Store in cache
	if jsonData, err := json.Marshal(response); err == nil {
		s.redisClient.Set(ctx, cacheKey, jsonData, cache.TTLNearbyPlaces)
	}

	return response, nil
}

func (s *SearchServiceImpl) GetSearchHistory(ctx context.Context, userID uuid.UUID, req *dto.GetSearchHistoryRequest) (*dto.SearchHistoryListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	var histories []*models.SearchHistory
	var total int64
	var err error

	if req.SearchType != "" {
		histories, err = s.searchHistoryRepo.GetByUserIDAndType(ctx, userID, req.SearchType, offset, req.PageSize)
		if err != nil {
			return nil, err
		}
		total, err = s.searchHistoryRepo.CountByUserIDAndType(ctx, userID, req.SearchType)
	} else {
		histories, err = s.searchHistoryRepo.GetByUserID(ctx, userID, offset, req.PageSize)
		if err != nil {
			return nil, err
		}
		total, err = s.searchHistoryRepo.CountByUserID(ctx, userID)
	}

	if err != nil {
		return nil, err
	}

	var historyResponses []dto.SearchHistoryResponse
	for _, h := range histories {
		historyResponses = append(historyResponses, dto.SearchHistoryResponse{
			ID:          h.ID,
			Query:       h.Query,
			SearchType:  h.SearchType,
			ResultCount: h.ResultCount,
			CreatedAt:   h.CreatedAt,
		})
	}

	return &dto.SearchHistoryListResponse{
		Histories: historyResponses,
		Meta: dto.PaginationMeta{
			Total:  total,
			Offset: offset,
			Limit:  req.PageSize,
		},
	}, nil
}

func (s *SearchServiceImpl) ClearSearchHistory(ctx context.Context, userID uuid.UUID, req *dto.ClearSearchHistoryRequest) error {
	if req.SearchType != "" {
		return s.searchHistoryRepo.DeleteByUserIDAndType(ctx, userID, req.SearchType)
	}
	return s.searchHistoryRepo.DeleteByUserID(ctx, userID)
}

func (s *SearchServiceImpl) DeleteSearchHistoryItem(ctx context.Context, userID uuid.UUID, historyID uuid.UUID) error {
	history, err := s.searchHistoryRepo.GetByID(ctx, historyID)
	if err != nil {
		return errors.New("history not found")
	}

	if history.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.searchHistoryRepo.Delete(ctx, historyID)
}

func (s *SearchServiceImpl) saveSearchHistory(ctx context.Context, userID uuid.UUID, query, searchType string, resultCount int) {
	if userID == uuid.Nil {
		return
	}

	history := &models.SearchHistory{
		UserID:      userID,
		Query:       query,
		SearchType:  searchType,
		ResultCount: resultCount,
	}

	_ = s.searchHistoryRepo.Create(ctx, history)
}

func formatDistance(meters float64) string {
	if meters < 1000 {
		return fmt.Sprintf("%.0f m", meters)
	}
	return fmt.Sprintf("%.1f km", meters/1000)
}

// In-memory map to track generating places (simple approach)
var generatingPlaces = make(map[string]bool)
var generatingMutex = &sync.Mutex{}

// GetPlaceDetailsEnhanced returns place details with AI-generated content
// Returns immediately - AI content generates in background if not cached
func (s *SearchServiceImpl) GetPlaceDetailsEnhanced(ctx context.Context, placeID string, userLat, userLng float64, includeAI bool) (*dto.PlaceDetailEnhancedResponse, error) {
	// 1. Get basic place details first
	basicDetails, err := s.GetPlaceDetails(ctx, placeID, userLat, userLng)
	if err != nil {
		return nil, err
	}

	// Build enhanced response from basic details
	response := &dto.PlaceDetailEnhancedResponse{
		PlaceID:          basicDetails.PlaceID,
		Name:             basicDetails.Name,
		FormattedAddress: basicDetails.FormattedAddress,
		Lat:              basicDetails.Lat,
		Lng:              basicDetails.Lng,
		Rating:           basicDetails.Rating,
		ReviewCount:      basicDetails.ReviewCount,
		PriceLevel:       basicDetails.PriceLevel,
		Types:            basicDetails.Types,
		Phone:            basicDetails.Phone,
		Website:          basicDetails.Website,
		GoogleMapsURL:    basicDetails.GoogleMapsURL,
		OpeningHours:     basicDetails.OpeningHours,
		Reviews:          basicDetails.Reviews,
		Photos:           basicDetails.Photos,
		Distance:         basicDetails.Distance,
		DistanceText:     basicDetails.DistanceText,
		AIStatus:         "unavailable",
	}

	// If AI content not requested, return basic response
	if !includeAI {
		return response, nil
	}

	// 2. Check if AI content exists in database
	aiContent, err := s.placeAIContentRepo.GetByPlaceID(ctx, placeID)
	if err == nil && aiContent != nil {
		// Found in database - use cached content
		response.AIStatus = "ready"
		response.AIOverview = s.mapAIContentToOverview(aiContent)
		response.GuideInfo = s.mapAIContentToGuideInfo(aiContent)
		response.RelatedVideos = s.mapAIContentToVideos(aiContent)
		return response, nil
	}

	// 3. Check if already generating
	generatingMutex.Lock()
	isGenerating := generatingPlaces[placeID]
	if !isGenerating {
		// Mark as generating
		generatingPlaces[placeID] = true
	}
	generatingMutex.Unlock()

	if isGenerating {
		// Already generating - return with generating status
		response.AIStatus = "generating"
		return response, nil
	}

	// 4. Start background generation
	response.AIStatus = "generating"
	go s.generateAIContentBackground(placeID, basicDetails)

	return response, nil
}

// generateAIContentBackground generates AI content in background
func (s *SearchServiceImpl) generateAIContentBackground(placeID string, basicDetails *dto.PlaceDetailResponse) {
	// Create a new context for background operation
	ctx := context.Background()

	defer func() {
		// Remove from generating map when done
		generatingMutex.Lock()
		delete(generatingPlaces, placeID)
		generatingMutex.Unlock()
	}()

	// Generate AI content
	aiContent, err := s.generateAIContent(ctx, basicDetails)
	if err != nil {
		fmt.Printf("Background: Failed to generate AI content for place %s: %v\n", placeID, err)
		return
	}

	// Save to database
	if err := s.placeAIContentRepo.Upsert(ctx, aiContent); err != nil {
		fmt.Printf("Background: Failed to save AI content for place %s: %v\n", placeID, err)
		return
	}

	fmt.Printf("Background: Successfully generated AI content for place %s\n", placeID)
}

// generateAIContent generates AI content for a place
func (s *SearchServiceImpl) generateAIContent(ctx context.Context, place *dto.PlaceDetailResponse) (*models.PlaceAIContent, error) {
	// Generate AI overview using OpenAI
	aiOverview, err := s.generateAIOverview(ctx, place)
	if err != nil {
		return nil, fmt.Errorf("generate AI overview: %w", err)
	}

	// Generate guide info using OpenAI
	guideInfo, err := s.generateGuideInfo(ctx, place)
	if err != nil {
		// Don't fail, just skip guide info
		fmt.Printf("Failed to generate guide info: %v\n", err)
	}

	// Get related videos from YouTube
	videos, err := s.getRelatedVideos(ctx, place.Name)
	if err != nil {
		// Don't fail, just skip videos
		fmt.Printf("Failed to get related videos: %v\n", err)
	}

	// Marshal JSON fields
	highlightsJSON, _ := json.Marshal(aiOverview.Highlights)
	tipsJSON, _ := json.Marshal(aiOverview.Tips)
	quickFactsJSON, _ := json.Marshal(guideInfo.QuickFacts)
	talkingPointsJSON, _ := json.Marshal(guideInfo.TalkingPoints)
	commonQuestionsJSON, _ := json.Marshal(guideInfo.CommonQuestions)
	videosJSON, _ := json.Marshal(videos)

	// Create content record
	content := &models.PlaceAIContent{
		PlaceID:         place.PlaceID,
		PlaceName:       place.Name,
		Summary:         aiOverview.Summary,
		History:         aiOverview.History,
		Highlights:      highlightsJSON,
		BestTimeToVisit: aiOverview.BestTimeToVisit,
		Tips:            tipsJSON,
		QuickFacts:      quickFactsJSON,
		TalkingPoints:   talkingPointsJSON,
		CommonQuestions: commonQuestionsJSON,
		RelatedVideos:   videosJSON,
		Language:        "th",
		GeneratedAt:     time.Now(),
		ExpiresAt:       time.Now().AddDate(0, 1, 0), // 1 month expiry
	}

	return content, nil
}

// generateAIOverview generates AI overview using OpenAI
func (s *SearchServiceImpl) generateAIOverview(ctx context.Context, place *dto.PlaceDetailResponse) (*dto.AIPlaceOverview, error) {
	prompt := fmt.Sprintf(`คุณเป็นมัคคุเทศก์ผู้เชี่ยวชาญด้านการท่องเที่ยวไทย กรุณาสร้างข้อมูลที่ละเอียดและมีประโยชน์เกี่ยวกับสถานที่นี้:

📍 ชื่อสถานที่: %s
📍 ที่ตั้ง: %s
📍 ประเภท: %v
⭐ คะแนน: %.1f (%d รีวิว)
🌐 พิกัด: %.6f, %.6f

กรุณาสร้างข้อมูลในรูปแบบ JSON ดังนี้:
{
    "summary": "ภาพรวมของสถานที่ที่ครอบคลุมและน่าสนใจ อธิบายว่าสถานที่นี้คืออะไร มีความสำคัญอย่างไร ทำไมนักท่องเที่ยวควรมาเยี่ยมชม (5-7 ประโยค ประมาณ 150-200 คำ)",
    "history": "ประวัติความเป็นมาที่ละเอียด รวมถึงปีที่ก่อตั้ง/สร้าง ผู้ก่อตั้ง เหตุการณ์สำคัญ และวิวัฒนาการตลอดประวัติศาสตร์ (2-3 ย่อหน้า ประมาณ 200-300 คำ)",
    "highlights": [
        "จุดเด่นที่ 1 - อธิบายสั้นๆ ว่าทำไมถึงพิเศษ",
        "จุดเด่นที่ 2 - สิ่งที่น่าสนใจเฉพาะของสถานที่นี้",
        "จุดเด่นที่ 3 - กิจกรรมหรือประสบการณ์ที่ห้ามพลาด",
        "จุดเด่นที่ 4 - สถาปัตยกรรม/ศิลปะ/ธรรมชาติที่โดดเด่น",
        "จุดเด่นที่ 5 - สิ่งที่ทำให้แตกต่างจากที่อื่น"
    ],
    "bestTimeToVisit": "เวลาที่เหมาะสมในการเยี่ยมชม รวมถึงฤดูกาล ช่วงเวลาของวัน และเหตุผล (2-3 ประโยค)",
    "tips": [
        "เคล็ดลับที่ 1 - การเตรียมตัวก่อนมา",
        "เคล็ดลับที่ 2 - สิ่งที่ควรรู้เกี่ยวกับการแต่งกาย/มารยาท",
        "เคล็ดลับที่ 3 - จุดถ่ายรูปที่ดีที่สุด",
        "เคล็ดลับที่ 4 - ร้านอาหาร/ที่พักใกล้เคียง",
        "เคล็ดลับที่ 5 - การเดินทางและที่จอดรถ",
        "เคล็ดลับที่ 6 - ค่าใช้จ่ายและเวลาที่ควรใช้"
    ]
}

⚠️ กฎสำคัญ:
- ตอบเป็นภาษาไทยเท่านั้น
- ข้อมูลต้องถูกต้องตามความเป็นจริง ถ้าไม่แน่ใจให้ระบุว่า "ควรตรวจสอบข้อมูลเพิ่มเติม"
- เนื้อหาต้องละเอียดและเป็นประโยชน์สำหรับมัคคุเทศก์
- ตอบเฉพาะ JSON เท่านั้น ไม่ต้องมีข้อความอื่น`, place.Name, place.FormattedAddress, place.Types, place.Rating, place.ReviewCount, place.Lat, place.Lng)

	messages := []openai.ChatMessage{
		{Role: "system", Content: "คุณเป็นมัคคุเทศก์ผู้เชี่ยวชาญด้านการท่องเที่ยวไทยที่มีประสบการณ์มากกว่า 20 ปี คุณมีความรู้ลึกซึ้งเกี่ยวกับประวัติศาสตร์ วัฒนธรรม และสถานที่ท่องเที่ยวทั่วประเทศไทย ให้ข้อมูลที่ถูกต้อง ละเอียด และเป็นประโยชน์สำหรับการนำเที่ยว"},
		{Role: "user", Content: prompt},
	}

	response, err := s.openaiClient.Chat(ctx, messages, 3000, 0.7)
	if err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}

	// Parse JSON response
	var overview dto.AIPlaceOverview
	content := response.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &overview); err != nil {
		// Try to extract JSON from response
		if start := findJSONStart(content); start >= 0 {
			if end := findJSONEnd(content, start); end > start {
				if err := json.Unmarshal([]byte(content[start:end+1]), &overview); err != nil {
					return nil, fmt.Errorf("parse AI response: %w", err)
				}
			}
		} else {
			return nil, fmt.Errorf("parse AI response: %w", err)
		}
	}

	overview.GeneratedAt = time.Now().Format(time.RFC3339)
	return &overview, nil
}

// generateGuideInfo generates guide info using OpenAI
func (s *SearchServiceImpl) generateGuideInfo(ctx context.Context, place *dto.PlaceDetailResponse) (*dto.PlaceGuideInfo, error) {
	prompt := fmt.Sprintf(`คุณเป็นมัคคุเทศก์มืออาชีพ กรุณาสร้างข้อมูลที่เป็นประโยชน์สำหรับการนำเที่ยวที่:

📍 สถานที่: %s
📍 ที่ตั้ง: %s
📍 ประเภท: %v

กรุณาสร้างข้อมูลในรูปแบบ JSON:
{
    "quickFacts": [
        "ข้อเท็จจริงที่ 1 - ข้อมูลตัวเลขหรือสถิติที่น่าสนใจ (เช่น พื้นที่ ปีที่สร้าง จำนวนผู้เข้าชม)",
        "ข้อเท็จจริงที่ 2 - ความพิเศษหรือสถิติที่โดดเด่น (เช่น ใหญ่ที่สุด เก่าที่สุด แห่งแรก)",
        "ข้อเท็จจริงที่ 3 - ข้อมูลที่นักท่องเที่ยวมักไม่รู้",
        "ข้อเท็จจริงที่ 4 - ความเชื่อมโยงกับประวัติศาสตร์หรือบุคคลสำคัญ",
        "ข้อเท็จจริงที่ 5 - ข้อมูลที่ทำให้สถานที่นี้มีเอกลักษณ์"
    ],
    "talkingPoints": [
        "ประเด็นที่ 1 - เรื่องราวที่น่าสนใจสำหรับเล่าให้นักท่องเที่ยวฟัง (2-3 ประโยค)",
        "ประเด็นที่ 2 - ตำนานหรือเรื่องเล่าที่เกี่ยวข้อง",
        "ประเด็นที่ 3 - ความสำคัญทางวัฒนธรรม/ศาสนา/ประวัติศาสตร์",
        "ประเด็นที่ 4 - เหตุการณ์พิเศษหรือเทศกาลที่จัดขึ้น",
        "ประเด็นที่ 5 - การเปรียบเทียบกับสถานที่อื่นที่คล้ายกัน"
    ],
    "commonQuestions": [
        {"question": "คำถามที่ 1 - คำถามเกี่ยวกับประวัติ/ที่มา", "answer": "คำตอบที่ละเอียดและถูกต้อง (3-4 ประโยค)"},
        {"question": "คำถามที่ 2 - คำถามเกี่ยวกับการเข้าชม/ค่าใช้จ่าย", "answer": "คำตอบที่ละเอียดพร้อมข้อมูลที่เป็นประโยชน์"},
        {"question": "คำถามที่ 3 - คำถามเกี่ยวกับสิ่งที่น่าสนใจ", "answer": "คำตอบที่ช่วยให้นักท่องเที่ยวได้รับประสบการณ์ที่ดี"},
        {"question": "คำถามที่ 4 - คำถามเกี่ยวกับข้อห้ามหรือมารยาท", "answer": "คำตอบที่ช่วยให้ปฏิบัติตัวได้ถูกต้อง"},
        {"question": "คำถามที่ 5 - คำถามอื่นที่นักท่องเที่ยวมักถาม", "answer": "คำตอบที่ครบถ้วนและเป็นประโยชน์"}
    ]
}

⚠️ กฎสำคัญ:
- ตอบเป็นภาษาไทยเท่านั้น
- ข้อมูลต้องถูกต้องและเป็นประโยชน์สำหรับมัคคุเทศก์จริงๆ
- คำตอบใน commonQuestions ต้องละเอียดพอที่จะตอบนักท่องเที่ยวได้
- ตอบเฉพาะ JSON เท่านั้น ไม่ต้องมีข้อความอื่น`, place.Name, place.FormattedAddress, place.Types)

	messages := []openai.ChatMessage{
		{Role: "system", Content: "คุณเป็นมัคคุเทศก์ผู้เชี่ยวชาญที่มีประสบการณ์นำเที่ยวมากกว่า 20 ปี คุณรู้วิธีเล่าเรื่องให้น่าสนใจและรู้คำถามที่นักท่องเที่ยวมักถาม ให้ข้อมูลที่ละเอียดและเป็นประโยชน์จริง"},
		{Role: "user", Content: prompt},
	}

	response, err := s.openaiClient.Chat(ctx, messages, 2500, 0.7)
	if err != nil {
		return &dto.PlaceGuideInfo{}, err
	}

	if len(response.Choices) == 0 {
		return &dto.PlaceGuideInfo{}, errors.New("no response from OpenAI")
	}

	var guideInfo dto.PlaceGuideInfo
	content := response.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &guideInfo); err != nil {
		if start := findJSONStart(content); start >= 0 {
			if end := findJSONEnd(content, start); end > start {
				if err := json.Unmarshal([]byte(content[start:end+1]), &guideInfo); err != nil {
					return &dto.PlaceGuideInfo{}, nil
				}
			}
		}
	}

	return &guideInfo, nil
}

// getRelatedVideos gets related YouTube videos
func (s *SearchServiceImpl) getRelatedVideos(ctx context.Context, placeName string) ([]dto.RelatedVideo, error) {
	searchReq := &google.VideoSearchRequest{
		Query:      placeName + " ท่องเที่ยว",
		MaxResults: 5,
		Order:      "relevance",
	}

	searchResponse, err := s.googleYouTube.SearchVideos(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	// Get video IDs for details
	var videoIDs []string
	for _, item := range searchResponse.Items {
		videoIDs = append(videoIDs, item.ID.VideoID)
	}

	// Get video details
	var detailsMap = make(map[string]google.VideoDetails)
	if len(videoIDs) > 0 {
		detailsResponse, err := s.googleYouTube.GetVideoDetails(ctx, videoIDs)
		if err == nil && detailsResponse != nil {
			for _, d := range detailsResponse.Items {
				detailsMap[d.ID] = d
			}
		}
	}

	var videos []dto.RelatedVideo
	for _, item := range searchResponse.Items {
		videoID := item.ID.VideoID
		duration := ""
		var viewCount int64
		if details, ok := detailsMap[videoID]; ok {
			duration = google.ParseDuration(details.ContentDetails.Duration)
			viewCount, _ = strconv.ParseInt(details.Statistics.ViewCount, 10, 64)
		}

		thumbnailURL := item.Snippet.Thumbnails.High.URL
		if thumbnailURL == "" {
			thumbnailURL = item.Snippet.Thumbnails.Default.URL
		}

		videos = append(videos, dto.RelatedVideo{
			VideoID:      videoID,
			Title:        item.Snippet.Title,
			ThumbnailURL: thumbnailURL,
			ChannelTitle: item.Snippet.ChannelTitle,
			Duration:     duration,
			ViewCount:    viewCount,
		})
	}

	return videos, nil
}

// Helper functions for mapping database content to DTOs
func (s *SearchServiceImpl) mapAIContentToOverview(content *models.PlaceAIContent) *dto.AIPlaceOverview {
	var highlights []string
	var tips []string
	_ = json.Unmarshal(content.Highlights, &highlights)
	_ = json.Unmarshal(content.Tips, &tips)

	return &dto.AIPlaceOverview{
		Summary:         content.Summary,
		History:         content.History,
		Highlights:      highlights,
		BestTimeToVisit: content.BestTimeToVisit,
		Tips:            tips,
		GeneratedAt:     content.GeneratedAt.Format(time.RFC3339),
	}
}

func (s *SearchServiceImpl) mapAIContentToGuideInfo(content *models.PlaceAIContent) *dto.PlaceGuideInfo {
	var quickFacts []string
	var talkingPoints []string
	var commonQuestions []dto.PlaceFAQ
	_ = json.Unmarshal(content.QuickFacts, &quickFacts)
	_ = json.Unmarshal(content.TalkingPoints, &talkingPoints)
	_ = json.Unmarshal(content.CommonQuestions, &commonQuestions)

	return &dto.PlaceGuideInfo{
		QuickFacts:      quickFacts,
		TalkingPoints:   talkingPoints,
		CommonQuestions: commonQuestions,
	}
}

func (s *SearchServiceImpl) mapAIContentToVideos(content *models.PlaceAIContent) []dto.RelatedVideo {
	var videos []dto.RelatedVideo
	_ = json.Unmarshal(content.RelatedVideos, &videos)
	return videos
}

// Helper to find JSON start in string
func findJSONStart(s string) int {
	for i, c := range s {
		if c == '{' {
			return i
		}
	}
	return -1
}

// Helper to find JSON end in string
func findJSONEnd(s string, start int) int {
	depth := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}
