package dto

// ==================== Translation DTOs ====================

type TranslateRequest struct {
	Text       string `json:"text" validate:"required,min=1,max=5000"`
	SourceLang string `json:"sourceLang" validate:"omitempty,len=2"`
	TargetLang string `json:"targetLang" validate:"required,len=2"`
}

type TranslateResponse struct {
	OriginalText   string `json:"originalText"`
	TranslatedText string `json:"translatedText"`
	SourceLang     string `json:"sourceLang"`
	TargetLang     string `json:"targetLang"`
	DetectedLang   string `json:"detectedLang,omitempty"`
}

type DetectLanguageRequest struct {
	Text string `json:"text" validate:"required,min=1,max=1000"`
}

type DetectLanguageResponse struct {
	Text       string `json:"text"`
	Language   string `json:"language"`
	Confidence float64 `json:"confidence"`
}

// ==================== QR Code DTOs ====================

type GenerateQRRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
	Size    int    `json:"size" validate:"omitempty,min=100,max=1000"`
	Format  string `json:"format" validate:"omitempty,oneof=png svg"`
}

type GenerateQRResponse struct {
	Content  string `json:"content"`
	QRCodeURL string `json:"qrCodeUrl"`
	Size     int    `json:"size"`
	Format   string `json:"format"`
}

// ==================== Distance Calculation DTOs ====================

type CalculateDistanceRequest struct {
	OriginLat      float64 `json:"originLat" validate:"required,latitude"`
	OriginLng      float64 `json:"originLng" validate:"required,longitude"`
	DestinationLat float64 `json:"destinationLat" validate:"required,latitude"`
	DestinationLng float64 `json:"destinationLng" validate:"required,longitude"`
}

type CalculateDistanceResponse struct {
	DistanceMeters float64 `json:"distanceMeters"`
	DistanceKm     float64 `json:"distanceKm"`
	DistanceText   string  `json:"distanceText"`
}

// ==================== Nearby Places DTOs ====================

type NearbyPlacesRequest struct {
	Lat       float64 `json:"lat" query:"lat" validate:"required,latitude"`
	Lng       float64 `json:"lng" query:"lng" validate:"required,longitude"`
	Radius    int     `json:"radius" query:"radius" validate:"omitempty,min=100,max=50000"`
	PlaceType string  `json:"type" query:"type" validate:"omitempty"`
	Keyword   string  `json:"keyword" query:"keyword" validate:"omitempty,max=200"`
	Page      int     `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize  int     `json:"pageSize" query:"pageSize" validate:"omitempty,min=1,max=20"`
}

// ==================== Location DTOs ====================

type LocationRequest struct {
	Lat float64 `json:"lat" query:"lat" validate:"required,latitude"`
	Lng float64 `json:"lng" query:"lng" validate:"required,longitude"`
}

type GeocodeRequest struct {
	Address string `json:"address" query:"address" validate:"required,min=1,max=500"`
}

type GeocodeResponse struct {
	Address       string  `json:"address"`
	FormattedAddr string  `json:"formattedAddress"`
	Lat           float64 `json:"lat"`
	Lng           float64 `json:"lng"`
	PlaceID       string  `json:"placeId,omitempty"`
}

type ReverseGeocodeResponse struct {
	Lat           float64 `json:"lat"`
	Lng           float64 `json:"lng"`
	FormattedAddr string  `json:"formattedAddress"`
	Components    AddressComponents `json:"components"`
}

type AddressComponents struct {
	Street       string `json:"street,omitempty"`
	Subdistrict  string `json:"subdistrict,omitempty"`
	District     string `json:"district,omitempty"`
	Province     string `json:"province,omitempty"`
	Country      string `json:"country,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
}

// ==================== Health Check DTOs ====================

type HealthCheckResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
}
