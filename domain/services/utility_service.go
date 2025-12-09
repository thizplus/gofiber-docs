package services

import (
	"context"

	"gofiber-template/domain/dto"
)

type UtilityService interface {
	// Translation
	Translate(ctx context.Context, req *dto.TranslateRequest) (*dto.TranslateResponse, error)
	DetectLanguage(ctx context.Context, req *dto.DetectLanguageRequest) (*dto.DetectLanguageResponse, error)

	// QR Code generation
	GenerateQRCode(ctx context.Context, req *dto.GenerateQRRequest) (*dto.GenerateQRResponse, error)

	// Distance calculation
	CalculateDistance(ctx context.Context, req *dto.CalculateDistanceRequest) (*dto.CalculateDistanceResponse, error)

	// Geocoding
	Geocode(ctx context.Context, req *dto.GeocodeRequest) (*dto.GeocodeResponse, error)
	ReverseGeocode(ctx context.Context, req *dto.LocationRequest) (*dto.ReverseGeocodeResponse, error)

	// Health check
	HealthCheck(ctx context.Context) (*dto.HealthCheckResponse, error)
}
