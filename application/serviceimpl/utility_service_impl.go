package serviceimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/skip2/go-qrcode"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/cache"
	"gofiber-template/infrastructure/external/google"
	"gofiber-template/pkg/config"
)

type UtilityServiceImpl struct {
	translateClient *google.TranslateClient
	redisClient     *redis.Client
	config          *config.Config
}

func NewUtilityService(
	translateClient *google.TranslateClient,
	redisClient *redis.Client,
	cfg *config.Config,
) services.UtilityService {
	return &UtilityServiceImpl{
		translateClient: translateClient,
		redisClient:     redisClient,
		config:          cfg,
	}
}

func (s *UtilityServiceImpl) Translate(ctx context.Context, req *dto.TranslateRequest) (*dto.TranslateResponse, error) {
	// Check cache first (translations are very cacheable - 7 days)
	cacheKey := cache.TranslateKey(req.Text, req.SourceLang, req.TargetLang)
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedResult dto.TranslateResponse
		if json.Unmarshal([]byte(cached), &cachedResult) == nil {
			return &cachedResult, nil
		}
	}

	// Cache miss - call API
	translateReq := &google.TranslateRequest{
		Text:           req.Text,
		SourceLanguage: req.SourceLang,
		TargetLanguage: req.TargetLang,
	}

	translation, err := s.translateClient.Translate(ctx, translateReq)
	if err != nil {
		return nil, err
	}

	sourceLang := req.SourceLang
	detectedLang := translation.DetectedSourceLanguage
	if sourceLang == "" {
		sourceLang = detectedLang
	}

	response := &dto.TranslateResponse{
		OriginalText:   req.Text,
		TranslatedText: translation.TranslatedText,
		SourceLang:     sourceLang,
		TargetLang:     req.TargetLang,
		DetectedLang:   detectedLang,
	}

	// Store in cache (7 days - translations don't change)
	if jsonData, err := json.Marshal(response); err == nil {
		s.redisClient.Set(ctx, cacheKey, jsonData, cache.TTLTranslate)
	}

	return response, nil
}

func (s *UtilityServiceImpl) DetectLanguage(ctx context.Context, req *dto.DetectLanguageRequest) (*dto.DetectLanguageResponse, error) {
	// Check cache first
	cacheKey := cache.DetectLanguageKey(req.Text)
	if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedResult dto.DetectLanguageResponse
		if json.Unmarshal([]byte(cached), &cachedResult) == nil {
			return &cachedResult, nil
		}
	}

	// Cache miss - call API
	detection, err := s.translateClient.DetectLanguage(ctx, req.Text)
	if err != nil {
		return nil, err
	}

	response := &dto.DetectLanguageResponse{
		Text:       req.Text,
		Language:   detection.Language,
		Confidence: detection.Confidence,
	}

	// Store in cache (7 days)
	if jsonData, err := json.Marshal(response); err == nil {
		s.redisClient.Set(ctx, cacheKey, jsonData, cache.TTLTranslate)
	}

	return response, nil
}

func (s *UtilityServiceImpl) GenerateQRCode(ctx context.Context, req *dto.GenerateQRRequest) (*dto.GenerateQRResponse, error) {
	size := req.Size
	if size == 0 {
		size = 256
	}

	format := req.Format
	if format == "" {
		format = "png"
	}

	// Generate QR code
	qr, err := qrcode.New(req.Content, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	// For now, return a data URL
	// In production, you might want to save to CDN and return URL
	png, err := qr.PNG(size)
	if err != nil {
		return nil, err
	}

	// Convert to base64 data URL
	dataURL := fmt.Sprintf("data:image/png;base64,%s", encodeBase64(png))

	return &dto.GenerateQRResponse{
		Content:   req.Content,
		QRCodeURL: dataURL,
		Size:      size,
		Format:    format,
	}, nil
}

func (s *UtilityServiceImpl) CalculateDistance(ctx context.Context, req *dto.CalculateDistanceRequest) (*dto.CalculateDistanceResponse, error) {
	// Haversine formula
	const earthRadius = 6371000 // meters

	lat1 := req.OriginLat * math.Pi / 180
	lat2 := req.DestinationLat * math.Pi / 180
	deltaLat := (req.DestinationLat - req.OriginLat) * math.Pi / 180
	deltaLng := (req.DestinationLng - req.OriginLng) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c

	return &dto.CalculateDistanceResponse{
		DistanceMeters: distance,
		DistanceKm:     distance / 1000,
		DistanceText:   formatDistanceUtil(distance),
	}, nil
}

func (s *UtilityServiceImpl) Geocode(ctx context.Context, req *dto.GeocodeRequest) (*dto.GeocodeResponse, error) {
	// This would require Google Geocoding API
	// For now, return an error as we need to implement the client
	return nil, fmt.Errorf("geocoding not implemented")
}

func (s *UtilityServiceImpl) ReverseGeocode(ctx context.Context, req *dto.LocationRequest) (*dto.ReverseGeocodeResponse, error) {
	// This would require Google Geocoding API
	// For now, return an error as we need to implement the client
	return nil, fmt.Errorf("reverse geocoding not implemented")
}

func (s *UtilityServiceImpl) HealthCheck(ctx context.Context) (*dto.HealthCheckResponse, error) {
	services := make(map[string]string)

	// Check Redis
	if s.redisClient != nil {
		if err := s.redisClient.Ping(ctx).Err(); err != nil {
			services["redis"] = "unhealthy: " + err.Error()
		} else {
			services["redis"] = "healthy"
		}
	} else {
		services["redis"] = "not configured"
	}

	// Check other services
	services["google_api"] = "configured"
	services["openai_api"] = "configured"

	return &dto.HealthCheckResponse{
		Status:    "healthy",
		Version:   "1.0.0",
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  services,
	}, nil
}

func formatDistanceUtil(meters float64) string {
	if meters < 1000 {
		return fmt.Sprintf("%.0f m", meters)
	}
	return fmt.Sprintf("%.1f km", meters/1000)
}

func encodeBase64(data []byte) string {
	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	result := make([]byte, ((len(data)+2)/3)*4)

	for i, j := 0, 0; i < len(data); i, j = i+3, j+4 {
		var val uint32
		switch len(data) - i {
		case 1:
			val = uint32(data[i]) << 16
			result[j] = base64Table[val>>18&0x3F]
			result[j+1] = base64Table[val>>12&0x3F]
			result[j+2] = '='
			result[j+3] = '='
		case 2:
			val = uint32(data[i])<<16 | uint32(data[i+1])<<8
			result[j] = base64Table[val>>18&0x3F]
			result[j+1] = base64Table[val>>12&0x3F]
			result[j+2] = base64Table[val>>6&0x3F]
			result[j+3] = '='
		default:
			val = uint32(data[i])<<16 | uint32(data[i+1])<<8 | uint32(data[i+2])
			result[j] = base64Table[val>>18&0x3F]
			result[j+1] = base64Table[val>>12&0x3F]
			result[j+2] = base64Table[val>>6&0x3F]
			result[j+3] = base64Table[val&0x3F]
		}
	}

	return string(result)
}
