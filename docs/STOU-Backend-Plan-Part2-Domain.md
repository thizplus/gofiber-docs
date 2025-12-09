# STOU Smart Tour - Backend Development Plan
# Part 2: Domain Layer (Models, DTOs, Interfaces)

---

## Table of Contents - All Parts

| Part | หัวข้อ | สถานะ |
|------|--------|-------|
| Part 1 | Project Overview & Foundation | ✅ Done |
| **Part 2** | **Domain Layer (Models, DTOs, Interfaces)** | 📍 Current |
| Part 3 | Infrastructure Layer (External APIs, Cache) | ⏳ Pending |
| Part 4 | Application Layer (Services Implementation) | ⏳ Pending |
| Part 5 | Interface Layer (Handlers, Routes, Middleware) | ⏳ Pending |

---

## 1. Models (domain/models/)

### 1.1 Update: user.go (เพิ่ม student_id)

```go
// domain/models/user.go
// อัพเดทจากโครงสร้างเดิม - เพิ่ม StudentID

package models

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    StudentID string    `gorm:"type:varchar(20);uniqueIndex"` // 🆕 NEW: รหัสนักศึกษา
    Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
    Username  string    `gorm:"type:varchar(100);uniqueIndex;not null"`
    Password  string    `gorm:"type:varchar(255);not null"`
    FirstName string    `gorm:"type:varchar(100)"`
    LastName  string    `gorm:"type:varchar(100)"`
    Avatar    string    `gorm:"type:text"`
    Role      string    `gorm:"type:varchar(20);default:'user'"` // user, admin
    IsActive  bool      `gorm:"default:true"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`

    // Relationships
    Folders       []Folder       `gorm:"foreignKey:UserID"`
    Favorites     []Favorite     `gorm:"foreignKey:UserID"`
    SearchHistory []SearchHistory `gorm:"foreignKey:UserID"`
    AIChatSessions []AIChatSession `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
    return "users"
}
```

### 1.2 New: folder.go

```go
// domain/models/folder.go

package models

import (
    "time"
    "github.com/google/uuid"
)

type Folder struct {
    ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID        uuid.UUID `gorm:"type:uuid;not null;index"`
    Name          string    `gorm:"type:varchar(255);not null"`
    Description   string    `gorm:"type:text"`
    CoverImageURL string    `gorm:"type:text"`
    IsPublic      bool      `gorm:"default:false"`
    ItemCount     int       `gorm:"default:0"`
    CreatedAt     time.Time `gorm:"autoCreateTime"`
    UpdatedAt     time.Time `gorm:"autoUpdateTime"`

    // Relationships
    User  User         `gorm:"foreignKey:UserID"`
    Items []FolderItem `gorm:"foreignKey:FolderID"`
}

func (Folder) TableName() string {
    return "folders"
}

// Folder item types
const (
    FolderItemTypePlace   = "place"
    FolderItemTypeWebsite = "website"
    FolderItemTypeImage   = "image"
    FolderItemTypeVideo   = "video"
    FolderItemTypeLink    = "link"
)
```

### 1.3 New: folder_item.go

```go
// domain/models/folder_item.go

package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/datatypes"
)

type FolderItem struct {
    ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    FolderID     uuid.UUID      `gorm:"type:uuid;not null;index"`
    Type         string         `gorm:"type:varchar(50);not null"` // place, website, image, video, link
    Title        string         `gorm:"type:varchar(255);not null"`
    URL          string         `gorm:"type:text;not null"`
    ThumbnailURL string         `gorm:"type:text"`
    Description  string         `gorm:"type:text"`
    Metadata     datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
    SortOrder    int            `gorm:"default:0"`
    CreatedAt    time.Time      `gorm:"autoCreateTime"`

    // Relationships
    Folder Folder `gorm:"foreignKey:FolderID"`
}

func (FolderItem) TableName() string {
    return "folder_items"
}

// Metadata structure for different item types
type PlaceMetadata struct {
    PlaceID     string  `json:"place_id,omitempty"`
    Address     string  `json:"address,omitempty"`
    Rating      float64 `json:"rating,omitempty"`
    ReviewCount int     `json:"review_count,omitempty"`
    Lat         float64 `json:"lat,omitempty"`
    Lng         float64 `json:"lng,omitempty"`
}

type WebsiteMetadata struct {
    Snippet string `json:"snippet,omitempty"`
    Source  string `json:"source,omitempty"`
}

type ImageMetadata struct {
    Width  int    `json:"width,omitempty"`
    Height int    `json:"height,omitempty"`
    Source string `json:"source,omitempty"`
}

type VideoMetadata struct {
    Duration  string `json:"duration,omitempty"`
    Channel   string `json:"channel,omitempty"`
    VideoID   string `json:"video_id,omitempty"`
    ViewCount int64  `json:"view_count,omitempty"`
}
```

### 1.4 New: favorite.go

```go
// domain/models/favorite.go

package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/datatypes"
)

type Favorite struct {
    ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID       uuid.UUID      `gorm:"type:uuid;not null;index"`
    Type         string         `gorm:"type:varchar(50);not null"` // place, website, image, video
    ExternalID   string         `gorm:"type:varchar(255);index"`   // Google Place ID, etc.
    Title        string         `gorm:"type:varchar(255);not null"`
    URL          string         `gorm:"type:text;not null"`
    ThumbnailURL string         `gorm:"type:text"`
    Rating       float64        `gorm:"type:decimal(2,1)"`
    ReviewCount  int            `gorm:"default:0"`
    Address      string         `gorm:"type:text"`
    Metadata     datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
    CreatedAt    time.Time      `gorm:"autoCreateTime"`

    // Relationships
    User User `gorm:"foreignKey:UserID"`
}

func (Favorite) TableName() string {
    return "favorites"
}

// Favorite types
const (
    FavoriteTypePlace   = "place"
    FavoriteTypeWebsite = "website"
    FavoriteTypeImage   = "image"
    FavoriteTypeVideo   = "video"
)

// Unique constraint
func (Favorite) Indexes() []string {
    return []string{
        "CREATE UNIQUE INDEX IF NOT EXISTS idx_favorites_unique ON favorites(user_id, type, external_id)",
    }
}
```

### 1.5 New: search_history.go

```go
// domain/models/search_history.go

package models

import (
    "time"
    "github.com/google/uuid"
)

type SearchHistory struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
    Query       string    `gorm:"type:varchar(500);not null"`
    SearchType  string    `gorm:"type:varchar(50);not null;default:'all'"` // all, website, image, video, map, ai
    ResultCount int       `gorm:"default:0"`
    CreatedAt   time.Time `gorm:"autoCreateTime;index"`

    // Relationships
    User User `gorm:"foreignKey:UserID"`
}

func (SearchHistory) TableName() string {
    return "search_history"
}

// Search types
const (
    SearchTypeAll     = "all"
    SearchTypeWebsite = "website"
    SearchTypeImage   = "image"
    SearchTypeVideo   = "video"
    SearchTypeMap     = "map"
    SearchTypeAI      = "ai"
)
```

### 1.6 New: ai_chat_session.go

```go
// domain/models/ai_chat_session.go

package models

import (
    "time"
    "github.com/google/uuid"
)

type AIChatSession struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
    Title        string    `gorm:"type:varchar(255)"`
    InitialQuery string    `gorm:"type:varchar(500)"`
    CreatedAt    time.Time `gorm:"autoCreateTime"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime"`

    // Relationships
    User     User            `gorm:"foreignKey:UserID"`
    Messages []AIChatMessage `gorm:"foreignKey:SessionID"`
}

func (AIChatSession) TableName() string {
    return "ai_chat_sessions"
}
```

### 1.7 New: ai_chat_message.go

```go
// domain/models/ai_chat_message.go

package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/datatypes"
)

type AIChatMessage struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    SessionID uuid.UUID      `gorm:"type:uuid;not null;index"`
    Role      string         `gorm:"type:varchar(20);not null"` // user, assistant
    Content   string         `gorm:"type:text;not null"`
    Sources   datatypes.JSON `gorm:"type:jsonb;default:'[]'"`
    CreatedAt time.Time      `gorm:"autoCreateTime"`

    // Relationships
    Session AIChatSession `gorm:"foreignKey:SessionID"`
}

func (AIChatMessage) TableName() string {
    return "ai_chat_messages"
}

// Message roles
const (
    MessageRoleUser      = "user"
    MessageRoleAssistant = "assistant"
)

// Source structure
type MessageSource struct {
    Title   string `json:"title"`
    URL     string `json:"url"`
    Snippet string `json:"snippet,omitempty"`
}
```

---

## 2. DTOs (domain/dto/)

### 2.1 Update: user.go (เพิ่ม student_id)

```go
// domain/dto/user.go
// อัพเดทจากเดิม - เพิ่ม StudentID

package dto

import (
    "time"
    "github.com/google/uuid"
)

// ============================================
// Request DTOs
// ============================================

type CreateUserRequest struct {
    StudentID string `json:"student_id" validate:"omitempty,min=8,max=20"`
    Email     string `json:"email" validate:"required,email"`
    Username  string `json:"username" validate:"required,min=3,max=50"`
    Password  string `json:"password" validate:"required,min=8"`
    FirstName string `json:"first_name" validate:"required,min=1,max=100"`
    LastName  string `json:"last_name" validate:"required,min=1,max=100"`
}

type UpdateUserRequest struct {
    StudentID string `json:"student_id" validate:"omitempty,min=8,max=20"`
    FirstName string `json:"first_name" validate:"omitempty,min=1,max=100"`
    LastName  string `json:"last_name" validate:"omitempty,min=1,max=100"`
    Avatar    string `json:"avatar" validate:"omitempty,url"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
    StudentID string `json:"student_id" validate:"required,min=8,max=20"` // Required for STOU
    Email     string `json:"email" validate:"required,email"`
    Username  string `json:"username" validate:"required,min=3,max=50"`
    Password  string `json:"password" validate:"required,min=8"`
    FirstName string `json:"first_name" validate:"required,min=1,max=100"`
    LastName  string `json:"last_name" validate:"required,min=1,max=100"`
}

type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password" validate:"required"`
    NewPassword     string `json:"new_password" validate:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// ============================================
// Response DTOs
// ============================================

type UserResponse struct {
    ID        uuid.UUID `json:"id"`
    StudentID string    `json:"student_id,omitempty"`
    Email     string    `json:"email"`
    Username  string    `json:"username"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Avatar    string    `json:"avatar,omitempty"`
    Role      string    `json:"role"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
    AccessToken  string       `json:"access_token"`
    RefreshToken string       `json:"refresh_token,omitempty"`
    ExpiresIn    int          `json:"expires_in"` // seconds
    User         UserResponse `json:"user"`
}

type RegisterResponse struct {
    AccessToken  string       `json:"access_token"`
    RefreshToken string       `json:"refresh_token,omitempty"`
    ExpiresIn    int          `json:"expires_in"`
    User         UserResponse `json:"user"`
}
```

### 2.2 New: search.go

```go
// domain/dto/search.go

package dto

import (
    "time"
    "github.com/google/uuid"
)

// ============================================
// Request DTOs
// ============================================

type SearchRequest struct {
    Query    string `query:"q" validate:"required,min=1,max=500"`
    Type     string `query:"type" validate:"omitempty,oneof=all website image video"` // default: all
    Page     int    `query:"page" validate:"omitempty,min=1"`                          // default: 1
    PerPage  int    `query:"per_page" validate:"omitempty,min=1,max=50"`               // default: 20
    Location string `query:"location" validate:"omitempty"`                            // lat,lng
    Language string `query:"language" validate:"omitempty,oneof=th en"`                // default: th
}

type AISearchRequest struct {
    Query     string `query:"q" validate:"required,min=1,max=500"`
    SessionID string `query:"session_id" validate:"omitempty,uuid"`
}

type AIChatRequest struct {
    SessionID string `json:"session_id" validate:"required,uuid"`
    Message   string `json:"message" validate:"required,min=1,max=2000"`
    ImageURL  string `json:"image_url" validate:"omitempty,url"`
}

type PlacesSearchRequest struct {
    Lat     float64 `query:"lat" validate:"required,latitude"`
    Lng     float64 `query:"lng" validate:"required,longitude"`
    Radius  int     `query:"radius" validate:"omitempty,min=100,max=50000"` // meters, default: 5000
    Type    string  `query:"type" validate:"omitempty"`                     // restaurant, tourist_attraction, etc.
    Keyword string  `query:"keyword" validate:"omitempty,max=200"`
}

type SearchHistoryRequest struct {
    Page    int `query:"page" validate:"omitempty,min=1"`
    PerPage int `query:"per_page" validate:"omitempty,min=1,max=100"`
}

// ============================================
// Response DTOs
// ============================================

type SearchResponse struct {
    Results []SearchResult `json:"results"`
    Meta    SearchMeta     `json:"meta"`
}

type SearchResult struct {
    ID           string       `json:"id"`
    Type         string       `json:"type"` // website, image, video
    Title        string       `json:"title"`
    URL          string       `json:"url"`
    Snippet      string       `json:"snippet,omitempty"`
    ThumbnailURL string       `json:"thumbnail_url,omitempty"`
    Image        *ImageInfo   `json:"image,omitempty"`
    Video        *VideoInfo   `json:"video,omitempty"`
    Source       string       `json:"source,omitempty"`
}

type ImageInfo struct {
    Width  int `json:"width"`
    Height int `json:"height"`
}

type VideoInfo struct {
    Duration  string `json:"duration"`
    Channel   string `json:"channel"`
    VideoID   string `json:"video_id"`
    ViewCount int64  `json:"view_count,omitempty"`
}

type SearchMeta struct {
    Query      string `json:"query"`
    SearchType string `json:"search_type"`
    Page       int    `json:"page"`
    PerPage    int    `json:"per_page"`
    Total      int64  `json:"total"`
    TotalPages int    `json:"total_pages"`
}

// AI Search Response
type AISearchResponse struct {
    Summary             AISummary     `json:"summary"`
    RelatedVideos       []VideoResult `json:"related_videos"`
    SessionID           string        `json:"session_id"`
    FollowUpSuggestions []string      `json:"follow_up_suggestions"`
}

type AISummary struct {
    Content  string          `json:"content"` // Markdown
    Sections []AISection     `json:"sections,omitempty"`
    Sources  []SourceInfo    `json:"sources"`
}

type AISection struct {
    Title   string `json:"title"`
    Content string `json:"content"`
}

type SourceInfo struct {
    Title   string `json:"title"`
    URL     string `json:"url"`
    Snippet string `json:"snippet,omitempty"`
}

type VideoResult struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    Thumbnail string `json:"thumbnail"`
    Channel   string `json:"channel"`
    Duration  string `json:"duration"`
    URL       string `json:"url"`
    ViewCount int64  `json:"view_count,omitempty"`
}

// AI Chat Response
type AIChatResponse struct {
    Message   AIChatMessageResponse `json:"message"`
    SessionID string                `json:"session_id"`
}

type AIChatMessageResponse struct {
    Role    string       `json:"role"`
    Content string       `json:"content"`
    Sources []SourceInfo `json:"sources,omitempty"`
}

// Places Response
type PlacesResponse struct {
    Places []PlaceResult `json:"places"`
}

type PlaceResult struct {
    PlaceID     string   `json:"place_id"`
    Name        string   `json:"name"`
    Address     string   `json:"address"`
    Lat         float64  `json:"lat"`
    Lng         float64  `json:"lng"`
    Rating      float64  `json:"rating"`
    ReviewCount int      `json:"review_count"`
    Types       []string `json:"types"`
    PhotoURL    string   `json:"photo_url,omitempty"`
    IsOpen      *bool    `json:"is_open,omitempty"`
    Distance    float64  `json:"distance"` // meters
}

type PlaceDetailResponse struct {
    PlaceID      string         `json:"place_id"`
    Name         string         `json:"name"`
    Address      string         `json:"address"`
    Phone        string         `json:"phone,omitempty"`
    Website      string         `json:"website,omitempty"`
    Rating       float64        `json:"rating"`
    ReviewCount  int            `json:"review_count"`
    PriceLevel   int            `json:"price_level,omitempty"`
    OpeningHours *OpeningHours  `json:"opening_hours,omitempty"`
    Photos       []PhotoInfo    `json:"photos"`
    Reviews      []ReviewInfo   `json:"reviews"`
    Lat          float64        `json:"lat"`
    Lng          float64        `json:"lng"`
    Types        []string       `json:"types"`
}

type OpeningHours struct {
    WeekdayText []string `json:"weekday_text"`
    IsOpenNow   bool     `json:"is_open_now"`
}

type PhotoInfo struct {
    URL         string `json:"url"`
    Attribution string `json:"attribution,omitempty"`
}

type ReviewInfo struct {
    Author string    `json:"author"`
    Rating int       `json:"rating"`
    Text   string    `json:"text"`
    Time   time.Time `json:"time"`
}

// Search History Response
type SearchHistoryResponse struct {
    History []SearchHistoryItem `json:"history"`
    Meta    PaginationMeta      `json:"meta"`
}

type SearchHistoryItem struct {
    ID         uuid.UUID `json:"id"`
    Query      string    `json:"query"`
    SearchType string    `json:"search_type"`
    CreatedAt  time.Time `json:"created_at"`
}
```

### 2.3 New: folder.go

```go
// domain/dto/folder.go

package dto

import (
    "time"
    "github.com/google/uuid"
)

// ============================================
// Request DTOs
// ============================================

type CreateFolderRequest struct {
    Name        string `json:"name" validate:"required,min=1,max=255"`
    Description string `json:"description" validate:"omitempty,max=1000"`
    IsPublic    bool   `json:"is_public"`
}

type UpdateFolderRequest struct {
    Name        string `json:"name" validate:"omitempty,min=1,max=255"`
    Description string `json:"description" validate:"omitempty,max=1000"`
    IsPublic    *bool  `json:"is_public"`
}

type AddFolderItemRequest struct {
    Type         string                 `json:"type" validate:"required,oneof=place website image video link"`
    Title        string                 `json:"title" validate:"required,min=1,max=255"`
    URL          string                 `json:"url" validate:"required,url"`
    ThumbnailURL string                 `json:"thumbnail_url" validate:"omitempty,url"`
    Description  string                 `json:"description" validate:"omitempty,max=1000"`
    Metadata     map[string]interface{} `json:"metadata" validate:"omitempty"`
}

type FolderFilterRequest struct {
    Page    int `query:"page" validate:"omitempty,min=1"`
    PerPage int `query:"per_page" validate:"omitempty,min=1,max=100"`
}

// ============================================
// Response DTOs
// ============================================

type FolderResponse struct {
    ID            uuid.UUID `json:"id"`
    Name          string    `json:"name"`
    Description   string    `json:"description,omitempty"`
    CoverImageURL string    `json:"cover_image_url,omitempty"`
    IsPublic      bool      `json:"is_public"`
    ItemCount     int       `json:"item_count"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type FolderListResponse struct {
    Folders []FolderResponse `json:"folders"`
    Meta    PaginationMeta   `json:"meta"`
}

type FolderDetailResponse struct {
    Folder FolderResponse       `json:"folder"`
    Items  []FolderItemResponse `json:"items"`
    Meta   PaginationMeta       `json:"meta"`
}

type FolderItemResponse struct {
    ID           uuid.UUID              `json:"id"`
    Type         string                 `json:"type"`
    Title        string                 `json:"title"`
    URL          string                 `json:"url"`
    ThumbnailURL string                 `json:"thumbnail_url,omitempty"`
    Description  string                 `json:"description,omitempty"`
    Metadata     map[string]interface{} `json:"metadata,omitempty"`
    CreatedAt    time.Time              `json:"created_at"`
}

type FolderShareResponse struct {
    ShareURL string `json:"share_url"`
    QRCode   string `json:"qr_code"` // base64
}
```

### 2.4 New: favorite.go

```go
// domain/dto/favorite.go

package dto

import (
    "time"
    "github.com/google/uuid"
)

// ============================================
// Request DTOs
// ============================================

type AddFavoriteRequest struct {
    Type         string                 `json:"type" validate:"required,oneof=place website image video"`
    ExternalID   string                 `json:"external_id" validate:"omitempty,max=255"`
    Title        string                 `json:"title" validate:"required,min=1,max=255"`
    URL          string                 `json:"url" validate:"required,url"`
    ThumbnailURL string                 `json:"thumbnail_url" validate:"omitempty,url"`
    Rating       float64                `json:"rating" validate:"omitempty,min=0,max=5"`
    Address      string                 `json:"address" validate:"omitempty,max=500"`
    Metadata     map[string]interface{} `json:"metadata" validate:"omitempty"`
}

type FavoriteFilterRequest struct {
    Type    string `query:"type" validate:"omitempty,oneof=place website image video"`
    Page    int    `query:"page" validate:"omitempty,min=1"`
    PerPage int    `query:"per_page" validate:"omitempty,min=1,max=100"`
}

type CheckFavoriteRequest struct {
    Type       string `query:"type" validate:"required,oneof=place website image video"`
    ExternalID string `query:"external_id" validate:"required"`
}

// ============================================
// Response DTOs
// ============================================

type FavoriteResponse struct {
    ID           uuid.UUID              `json:"id"`
    Type         string                 `json:"type"`
    ExternalID   string                 `json:"external_id,omitempty"`
    Title        string                 `json:"title"`
    URL          string                 `json:"url"`
    ThumbnailURL string                 `json:"thumbnail_url,omitempty"`
    Rating       float64                `json:"rating,omitempty"`
    Address      string                 `json:"address,omitempty"`
    Metadata     map[string]interface{} `json:"metadata,omitempty"`
    CreatedAt    time.Time              `json:"created_at"`
}

type FavoriteListResponse struct {
    Favorites []FavoriteResponse `json:"favorites"`
    Meta      PaginationMeta     `json:"meta"`
}

type CheckFavoriteResponse struct {
    IsFavorited bool       `json:"is_favorited"`
    FavoriteID  *uuid.UUID `json:"favorite_id,omitempty"`
}
```

### 2.5 New: utility.go

```go
// domain/dto/utility.go

package dto

// ============================================
// Request DTOs
// ============================================

type TranslateRequest struct {
    Text           string `json:"text" validate:"required,min=1,max=5000"`
    SourceLanguage string `json:"source_language" validate:"omitempty,len=2"` // auto-detect if empty
    TargetLanguage string `json:"target_language" validate:"required,len=2"`  // th, en, etc.
}

type QRCodeRequest struct {
    Content string `json:"content" validate:"required,min=1,max=2000"`
    Size    int    `json:"size" validate:"omitempty,min=64,max=1024"` // default: 256
    Format  string `json:"format" validate:"omitempty,oneof=png svg"` // default: png
}

// ============================================
// Response DTOs
// ============================================

type TranslateResponse struct {
    TranslatedText string `json:"translated_text"`
    SourceLanguage string `json:"source_language"`
    TargetLanguage string `json:"target_language"`
}

type QRCodeResponse struct {
    QRCode  string `json:"qr_code"` // base64 encoded image
    Content string `json:"content"`
    Format  string `json:"format"`
}
```

### 2.6 Update: common.go (เพิ่มเติม)

```go
// domain/dto/common.go
// เพิ่มเติมจากเดิม

package dto

// ============================================
// Common Response DTOs
// ============================================

type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
    Success bool           `json:"success"`
    Message string         `json:"message,omitempty"`
    Data    interface{}    `json:"data,omitempty"`
    Meta    PaginationMeta `json:"meta"`
    Error   string         `json:"error,omitempty"`
}

type PaginationMeta struct {
    Total      int64 `json:"total"`
    Page       int   `json:"page"`
    PerPage    int   `json:"per_page"`
    TotalPages int   `json:"total_pages"`
}

// Helper function to calculate total pages
func CalculateTotalPages(total int64, perPage int) int {
    if perPage <= 0 {
        return 0
    }
    pages := int(total) / perPage
    if int(total)%perPage > 0 {
        pages++
    }
    return pages
}

// Default pagination values
const (
    DefaultPage    = 1
    DefaultPerPage = 20
    MaxPerPage     = 100
)

// ============================================
// Error Codes
// ============================================

const (
    ErrCodeValidation     = "VALIDATION_ERROR"
    ErrCodeUnauthorized   = "UNAUTHORIZED"
    ErrCodeForbidden      = "FORBIDDEN"
    ErrCodeNotFound       = "NOT_FOUND"
    ErrCodeConflict       = "CONFLICT"
    ErrCodeInternalError  = "INTERNAL_ERROR"
    ErrCodeRateLimited    = "RATE_LIMITED"
    ErrCodeExternalAPI    = "EXTERNAL_API_ERROR"
)

type ErrorDetail struct {
    Field   string `json:"field,omitempty"`
    Message string `json:"message"`
}

type ErrorResponse struct {
    Success bool          `json:"success"`
    Error   ErrorInfo     `json:"error"`
}

type ErrorInfo struct {
    Code    string        `json:"code"`
    Message string        `json:"message"`
    Details []ErrorDetail `json:"details,omitempty"`
}
```

### 2.7 New: mappers.go (เพิ่มเติม)

```go
// domain/dto/mappers.go
// เพิ่ม mapper functions สำหรับ models ใหม่

package dto

import (
    "encoding/json"
    "github.com/your-org/stou-smart-tour/domain/models"
)

// ============================================
// Folder Mappers
// ============================================

func FolderToFolderResponse(folder *models.Folder) FolderResponse {
    return FolderResponse{
        ID:            folder.ID,
        Name:          folder.Name,
        Description:   folder.Description,
        CoverImageURL: folder.CoverImageURL,
        IsPublic:      folder.IsPublic,
        ItemCount:     folder.ItemCount,
        CreatedAt:     folder.CreatedAt,
        UpdatedAt:     folder.UpdatedAt,
    }
}

func FoldersToFolderResponses(folders []models.Folder) []FolderResponse {
    responses := make([]FolderResponse, len(folders))
    for i, folder := range folders {
        responses[i] = FolderToFolderResponse(&folder)
    }
    return responses
}

func CreateFolderRequestToFolder(req *CreateFolderRequest, userID uuid.UUID) *models.Folder {
    return &models.Folder{
        UserID:      userID,
        Name:        req.Name,
        Description: req.Description,
        IsPublic:    req.IsPublic,
    }
}

// ============================================
// Folder Item Mappers
// ============================================

func FolderItemToFolderItemResponse(item *models.FolderItem) FolderItemResponse {
    var metadata map[string]interface{}
    if item.Metadata != nil {
        json.Unmarshal(item.Metadata, &metadata)
    }

    return FolderItemResponse{
        ID:           item.ID,
        Type:         item.Type,
        Title:        item.Title,
        URL:          item.URL,
        ThumbnailURL: item.ThumbnailURL,
        Description:  item.Description,
        Metadata:     metadata,
        CreatedAt:    item.CreatedAt,
    }
}

func FolderItemsToFolderItemResponses(items []models.FolderItem) []FolderItemResponse {
    responses := make([]FolderItemResponse, len(items))
    for i, item := range items {
        responses[i] = FolderItemToFolderItemResponse(&item)
    }
    return responses
}

func AddFolderItemRequestToFolderItem(req *AddFolderItemRequest, folderID uuid.UUID) (*models.FolderItem, error) {
    var metadata []byte
    var err error
    if req.Metadata != nil {
        metadata, err = json.Marshal(req.Metadata)
        if err != nil {
            return nil, err
        }
    }

    return &models.FolderItem{
        FolderID:     folderID,
        Type:         req.Type,
        Title:        req.Title,
        URL:          req.URL,
        ThumbnailURL: req.ThumbnailURL,
        Description:  req.Description,
        Metadata:     metadata,
    }, nil
}

// ============================================
// Favorite Mappers
// ============================================

func FavoriteToFavoriteResponse(fav *models.Favorite) FavoriteResponse {
    var metadata map[string]interface{}
    if fav.Metadata != nil {
        json.Unmarshal(fav.Metadata, &metadata)
    }

    return FavoriteResponse{
        ID:           fav.ID,
        Type:         fav.Type,
        ExternalID:   fav.ExternalID,
        Title:        fav.Title,
        URL:          fav.URL,
        ThumbnailURL: fav.ThumbnailURL,
        Rating:       fav.Rating,
        Address:      fav.Address,
        Metadata:     metadata,
        CreatedAt:    fav.CreatedAt,
    }
}

func FavoritesToFavoriteResponses(favorites []models.Favorite) []FavoriteResponse {
    responses := make([]FavoriteResponse, len(favorites))
    for i, fav := range favorites {
        responses[i] = FavoriteToFavoriteResponse(&fav)
    }
    return responses
}

func AddFavoriteRequestToFavorite(req *AddFavoriteRequest, userID uuid.UUID) (*models.Favorite, error) {
    var metadata []byte
    var err error
    if req.Metadata != nil {
        metadata, err = json.Marshal(req.Metadata)
        if err != nil {
            return nil, err
        }
    }

    return &models.Favorite{
        UserID:       userID,
        Type:         req.Type,
        ExternalID:   req.ExternalID,
        Title:        req.Title,
        URL:          req.URL,
        ThumbnailURL: req.ThumbnailURL,
        Rating:       req.Rating,
        Address:      req.Address,
        Metadata:     metadata,
    }, nil
}

// ============================================
// Search History Mappers
// ============================================

func SearchHistoryToSearchHistoryItem(history *models.SearchHistory) SearchHistoryItem {
    return SearchHistoryItem{
        ID:         history.ID,
        Query:      history.Query,
        SearchType: history.SearchType,
        CreatedAt:  history.CreatedAt,
    }
}

func SearchHistoriesToSearchHistoryItems(histories []models.SearchHistory) []SearchHistoryItem {
    items := make([]SearchHistoryItem, len(histories))
    for i, history := range histories {
        items[i] = SearchHistoryToSearchHistoryItem(&history)
    }
    return items
}
```

---

## 3. Repository Interfaces (domain/repositories/)

### 3.1 New: folder_repository.go

```go
// domain/repositories/folder_repository.go

package repositories

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/models"
)

type FolderRepository interface {
    // CRUD
    Create(ctx context.Context, folder *models.Folder) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.Folder, error)
    Update(ctx context.Context, folder *models.Folder) error
    Delete(ctx context.Context, id uuid.UUID) error

    // User-specific queries
    GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.Folder, error)
    CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)

    // Public folders
    GetPublicFolders(ctx context.Context, offset, limit int) ([]models.Folder, error)
    CountPublicFolders(ctx context.Context) (int64, error)

    // Update item count
    IncrementItemCount(ctx context.Context, folderID uuid.UUID) error
    DecrementItemCount(ctx context.Context, folderID uuid.UUID) error

    // Update cover image
    UpdateCoverImage(ctx context.Context, folderID uuid.UUID, imageURL string) error
}
```

### 3.2 New: folder_item_repository.go

```go
// domain/repositories/folder_item_repository.go

package repositories

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/models"
)

type FolderItemRepository interface {
    // CRUD
    Create(ctx context.Context, item *models.FolderItem) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.FolderItem, error)
    Delete(ctx context.Context, id uuid.UUID) error

    // Folder-specific queries
    GetByFolderID(ctx context.Context, folderID uuid.UUID, offset, limit int) ([]models.FolderItem, error)
    CountByFolderID(ctx context.Context, folderID uuid.UUID) (int64, error)

    // Get first item (for cover image)
    GetFirstImageByFolderID(ctx context.Context, folderID uuid.UUID) (*models.FolderItem, error)

    // Bulk operations
    DeleteByFolderID(ctx context.Context, folderID uuid.UUID) error

    // Reorder
    UpdateSortOrder(ctx context.Context, id uuid.UUID, sortOrder int) error
}
```

### 3.3 New: favorite_repository.go

```go
// domain/repositories/favorite_repository.go

package repositories

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/models"
)

type FavoriteRepository interface {
    // CRUD
    Create(ctx context.Context, favorite *models.Favorite) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.Favorite, error)
    Delete(ctx context.Context, id uuid.UUID) error

    // User-specific queries
    GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.Favorite, error)
    GetByUserIDAndType(ctx context.Context, userID uuid.UUID, favType string, offset, limit int) ([]models.Favorite, error)
    CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
    CountByUserIDAndType(ctx context.Context, userID uuid.UUID, favType string) (int64, error)

    // Check existence
    ExistsByUserAndExternal(ctx context.Context, userID uuid.UUID, favType, externalID string) (bool, error)
    GetByUserAndExternal(ctx context.Context, userID uuid.UUID, favType, externalID string) (*models.Favorite, error)
}
```

### 3.4 New: search_history_repository.go

```go
// domain/repositories/search_history_repository.go

package repositories

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/models"
)

type SearchHistoryRepository interface {
    // Create
    Create(ctx context.Context, history *models.SearchHistory) error

    // User-specific queries
    GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.SearchHistory, error)
    CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)

    // Delete
    DeleteByUserID(ctx context.Context, userID uuid.UUID) error
    DeleteOlderThan(ctx context.Context, userID uuid.UUID, days int) error

    // Analytics (optional)
    GetPopularQueries(ctx context.Context, limit int) ([]string, error)
}
```

### 3.5 New: ai_chat_repository.go

```go
// domain/repositories/ai_chat_repository.go

package repositories

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/models"
)

type AIChatRepository interface {
    // Session operations
    CreateSession(ctx context.Context, session *models.AIChatSession) error
    GetSessionByID(ctx context.Context, id uuid.UUID) (*models.AIChatSession, error)
    UpdateSession(ctx context.Context, session *models.AIChatSession) error
    DeleteSession(ctx context.Context, id uuid.UUID) error

    // User sessions
    GetSessionsByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.AIChatSession, error)
    CountSessionsByUserID(ctx context.Context, userID uuid.UUID) (int64, error)

    // Message operations
    CreateMessage(ctx context.Context, message *models.AIChatMessage) error
    GetMessagesBySessionID(ctx context.Context, sessionID uuid.UUID) ([]models.AIChatMessage, error)

    // Cleanup
    DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error
}
```

---

## 4. Service Interfaces (domain/services/)

### 4.1 New: search_service.go

```go
// domain/services/search_service.go

package services

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/dto"
)

type SearchService interface {
    // Google Custom Search
    Search(ctx context.Context, req *dto.SearchRequest) (*dto.SearchResponse, error)

    // Places
    SearchPlaces(ctx context.Context, req *dto.PlacesSearchRequest) (*dto.PlacesResponse, error)
    GetPlaceDetail(ctx context.Context, placeID string) (*dto.PlaceDetailResponse, error)

    // Search History
    SaveSearchHistory(ctx context.Context, userID uuid.UUID, query, searchType string, resultCount int) error
    GetSearchHistory(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.SearchHistoryResponse, error)
    ClearSearchHistory(ctx context.Context, userID uuid.UUID) error
}
```

### 4.2 New: ai_service.go

```go
// domain/services/ai_service.go

package services

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/dto"
)

type AIService interface {
    // AI Search (generates summary from search results)
    AISearch(ctx context.Context, query string, userID *uuid.UUID) (*dto.AISearchResponse, error)

    // AI Chat (continue conversation)
    Chat(ctx context.Context, req *dto.AIChatRequest, userID uuid.UUID) (*dto.AIChatResponse, error)

    // Get related videos
    GetRelatedVideos(ctx context.Context, query string, limit int) ([]dto.VideoResult, error)

    // Session management
    GetChatSession(ctx context.Context, sessionID uuid.UUID) (*models.AIChatSession, error)
    GetChatHistory(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.AIChatSession, error)
    DeleteChatSession(ctx context.Context, sessionID, userID uuid.UUID) error
}
```

### 4.3 New: folder_service.go

```go
// domain/services/folder_service.go

package services

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/dto"
    "github.com/your-org/stou-smart-tour/domain/models"
)

type FolderService interface {
    // Folder CRUD
    CreateFolder(ctx context.Context, req *dto.CreateFolderRequest, userID uuid.UUID) (*dto.FolderResponse, error)
    GetFolder(ctx context.Context, folderID, userID uuid.UUID) (*dto.FolderDetailResponse, error)
    UpdateFolder(ctx context.Context, folderID uuid.UUID, req *dto.UpdateFolderRequest, userID uuid.UUID) (*dto.FolderResponse, error)
    DeleteFolder(ctx context.Context, folderID, userID uuid.UUID) error

    // List folders
    GetUserFolders(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.FolderListResponse, error)

    // Folder items
    AddItem(ctx context.Context, folderID uuid.UUID, req *dto.AddFolderItemRequest, userID uuid.UUID) (*dto.FolderItemResponse, error)
    RemoveItem(ctx context.Context, folderID, itemID, userID uuid.UUID) error
    GetFolderItems(ctx context.Context, folderID uuid.UUID, offset, limit int) ([]dto.FolderItemResponse, int64, error)

    // Sharing
    GenerateShareLink(ctx context.Context, folderID, userID uuid.UUID) (*dto.FolderShareResponse, error)
    GetPublicFolder(ctx context.Context, folderID uuid.UUID) (*dto.FolderDetailResponse, error)
}
```

### 4.4 New: favorite_service.go

```go
// domain/services/favorite_service.go

package services

import (
    "context"
    "github.com/google/uuid"
    "github.com/your-org/stou-smart-tour/domain/dto"
)

type FavoriteService interface {
    // CRUD
    AddFavorite(ctx context.Context, req *dto.AddFavoriteRequest, userID uuid.UUID) (*dto.FavoriteResponse, error)
    RemoveFavorite(ctx context.Context, favoriteID, userID uuid.UUID) error

    // List
    GetFavorites(ctx context.Context, userID uuid.UUID, favType string, offset, limit int) (*dto.FavoriteListResponse, error)

    // Check
    IsFavorited(ctx context.Context, userID uuid.UUID, favType, externalID string) (*dto.CheckFavoriteResponse, error)
}
```

### 4.5 New: translate_service.go

```go
// domain/services/translate_service.go

package services

import (
    "context"
    "github.com/your-org/stou-smart-tour/domain/dto"
)

type TranslateService interface {
    Translate(ctx context.Context, req *dto.TranslateRequest) (*dto.TranslateResponse, error)
    DetectLanguage(ctx context.Context, text string) (string, error)
    GetSupportedLanguages(ctx context.Context) ([]string, error)
}
```

### 4.6 New: qrcode_service.go

```go
// domain/services/qrcode_service.go

package services

import (
    "context"
    "github.com/your-org/stou-smart-tour/domain/dto"
)

type QRCodeService interface {
    Generate(ctx context.Context, req *dto.QRCodeRequest) (*dto.QRCodeResponse, error)
}
```

---

## 5. Summary - Domain Layer Files

```
domain/
├── models/
│   ├── user.go                    # ✏️ UPDATE (เพิ่ม student_id)
│   ├── task.go                    # ✅ KEEP
│   ├── file.go                    # ✅ KEEP
│   ├── job.go                     # ✅ KEEP
│   ├── folder.go                  # 🆕 NEW
│   ├── folder_item.go             # 🆕 NEW
│   ├── favorite.go                # 🆕 NEW
│   ├── search_history.go          # 🆕 NEW
│   ├── ai_chat_session.go         # 🆕 NEW
│   └── ai_chat_message.go         # 🆕 NEW
│
├── repositories/
│   ├── user_repository.go         # ✅ KEEP
│   ├── task_repository.go         # ✅ KEEP
│   ├── file_repository.go         # ✅ KEEP
│   ├── job_repository.go          # ✅ KEEP
│   ├── folder_repository.go       # 🆕 NEW
│   ├── folder_item_repository.go  # 🆕 NEW
│   ├── favorite_repository.go     # 🆕 NEW
│   ├── search_history_repository.go # 🆕 NEW
│   └── ai_chat_repository.go      # 🆕 NEW
│
├── services/
│   ├── user_service.go            # ✅ KEEP
│   ├── task_service.go            # ✅ KEEP
│   ├── file_service.go            # ✅ KEEP
│   ├── job_service.go             # ✅ KEEP
│   ├── search_service.go          # 🆕 NEW
│   ├── ai_service.go              # 🆕 NEW
│   ├── folder_service.go          # 🆕 NEW
│   ├── favorite_service.go        # 🆕 NEW
│   ├── translate_service.go       # 🆕 NEW
│   └── qrcode_service.go          # 🆕 NEW
│
└── dto/
    ├── user.go                    # ✏️ UPDATE
    ├── auth.go                    # ✅ KEEP
    ├── task.go                    # ✅ KEEP
    ├── file.go                    # ✅ KEEP
    ├── job.go                     # ✅ KEEP
    ├── common.go                  # ✏️ UPDATE
    ├── mappers.go                 # ✏️ UPDATE
    ├── search.go                  # 🆕 NEW
    ├── ai.go                      # 🆕 NEW (หรือรวมใน search.go)
    ├── folder.go                  # 🆕 NEW
    ├── favorite.go                # 🆕 NEW
    └── utility.go                 # 🆕 NEW
```

---

## Next Part

➡️ ไปต่อที่ **Part 3: Infrastructure Layer (External APIs, Cache)**
- Google API Clients (Search, Places, YouTube, Translate)
- OpenAI/Anthropic Client
- Repository Implementations
- Cache Layer Enhancement

---

*Document Version: 1.0*
*Part: 2 of 5*
