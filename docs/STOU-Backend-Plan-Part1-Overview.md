# STOU Smart Tour - Backend Development Plan
# Part 1: Project Overview & Foundation

---

## Table of Contents - All Parts

| Part | หัวข้อ | สถานะ |
|------|--------|-------|
| **Part 1** | **Project Overview & Foundation** | 📍 Current |
| Part 2 | Domain Layer (Models, DTOs, Interfaces) | ⏳ Pending |
| Part 3 | Infrastructure Layer (External APIs, Cache) | ⏳ Pending |
| Part 4 | Application Layer (Services Implementation) | ⏳ Pending |
| Part 5 | Interface Layer (Handlers, Routes, Middleware) | ⏳ Pending |

---

## 1. Project Overview

### 1.1 สรุปโปรเจค STOU Smart Tour

```
┌─────────────────────────────────────────────────────────────────┐
│                    STOU Smart Tour                               │
│         ระบบค้นหาข้อมูลท่องเที่ยวสำหรับนักศึกษา มสธ.               │
└─────────────────────────────────────────────────────────────────┘

Features หลัก:
├── 🔍 Search System (Google Custom Search API)
│   ├── All - ค้นหาทุกประเภท
│   ├── Website - ค้นหาเว็บไซต์
│   ├── Image - ค้นหารูปภาพ
│   ├── Video - ค้นหาวิดีโอ
│   └── Map - แสดงบน Map
│
├── 🤖 AI Mode (OpenAI/Anthropic API)
│   ├── AI Summary - สรุปข้อมูลจาก AI
│   ├── Chat - ถามตอบกับ AI
│   └── Related Videos - วิดีโอที่เกี่ยวข้อง
│
├── 📁 Folder System
│   ├── Create/Edit/Delete folders
│   ├── Save items to folders
│   └── Share folders
│
├── ❤️ Favorites System
│   ├── Add/Remove favorites
│   └── Favorites list
│
├── 🔐 Authentication
│   ├── Register (student ID)
│   ├── Login/Logout
│   └── Profile management
│
└── 🛠️ Utilities
    ├── Translation (Google Translate)
    └── QR Code Generator
```

### 1.2 Tech Stack ที่ใช้

```
Backend:
├── Framework: Go Fiber v2.52.0
├── Architecture: Clean Architecture
├── Database: PostgreSQL + GORM
├── Cache: Redis
├── Auth: JWT
└── Existing Features: User, Task, File, Job (ใช้เป็น reference)

External APIs:
├── Google Custom Search API
├── Google Places API
├── Google Maps API
├── Google Translate API
├── YouTube Data API
└── OpenAI API / Anthropic API
```

---

## 2. การ Map กับโครงสร้างเดิม

### 2.1 โครงสร้าง Clean Architecture ปัจจุบัน

```
gofiber-docs/
├── cmd/api/main.go                    # Entry point
├── domain/                            # Domain layer
│   ├── models/                        # ✅ มี user.go, task.go, file.go, job.go
│   ├── repositories/                  # ✅ มี interfaces
│   ├── services/                      # ✅ มี interfaces
│   └── dto/                           # ✅ มี request/response DTOs
├── application/serviceimpl/           # ✅ Service implementations
├── infrastructure/                    # ✅ มี postgres, redis, storage
│   ├── postgres/
│   ├── redis/
│   └── storage/
├── interfaces/api/                    # ✅ มี handlers, middleware, routes
│   ├── handlers/
│   ├── middleware/
│   └── routes/
└── pkg/                               # ✅ มี config, di, utils
```

### 2.2 สิ่งที่ต้องเพิ่มใหม่

```
gofiber-docs/
├── domain/
│   ├── models/
│   │   ├── user.go                    # ✅ มีอยู่แล้ว (ปรับเพิ่ม student_id)
│   │   ├── folder.go                  # 🆕 NEW
│   │   ├── folder_item.go             # 🆕 NEW
│   │   ├── favorite.go                # 🆕 NEW
│   │   ├── search_history.go          # 🆕 NEW
│   │   ├── ai_chat_session.go         # 🆕 NEW
│   │   └── ai_chat_message.go         # 🆕 NEW
│   │
│   ├── repositories/
│   │   ├── folder_repository.go       # 🆕 NEW
│   │   ├── folder_item_repository.go  # 🆕 NEW
│   │   ├── favorite_repository.go     # 🆕 NEW
│   │   ├── search_history_repository.go # 🆕 NEW
│   │   └── ai_chat_repository.go      # 🆕 NEW
│   │
│   ├── services/
│   │   ├── search_service.go          # 🆕 NEW
│   │   ├── ai_service.go              # 🆕 NEW
│   │   ├── folder_service.go          # 🆕 NEW
│   │   ├── favorite_service.go        # 🆕 NEW
│   │   ├── translate_service.go       # 🆕 NEW
│   │   └── qrcode_service.go          # 🆕 NEW
│   │
│   └── dto/
│       ├── search.go                  # 🆕 NEW
│       ├── ai.go                      # 🆕 NEW
│       ├── folder.go                  # 🆕 NEW
│       ├── favorite.go                # 🆕 NEW
│       └── utility.go                 # 🆕 NEW
│
├── application/serviceimpl/
│   ├── search_service_impl.go         # 🆕 NEW
│   ├── ai_service_impl.go             # 🆕 NEW
│   ├── folder_service_impl.go         # 🆕 NEW
│   ├── favorite_service_impl.go       # 🆕 NEW
│   ├── translate_service_impl.go      # 🆕 NEW
│   └── qrcode_service_impl.go         # 🆕 NEW
│
├── infrastructure/
│   ├── postgres/
│   │   ├── folder_repository_impl.go      # 🆕 NEW
│   │   ├── folder_item_repository_impl.go # 🆕 NEW
│   │   ├── favorite_repository_impl.go    # 🆕 NEW
│   │   ├── search_history_repository_impl.go # 🆕 NEW
│   │   └── ai_chat_repository_impl.go     # 🆕 NEW
│   │
│   ├── external/                      # 🆕 NEW FOLDER
│   │   ├── google/
│   │   │   ├── search_client.go       # 🆕 Google Custom Search
│   │   │   ├── places_client.go       # 🆕 Google Places
│   │   │   ├── youtube_client.go      # 🆕 YouTube Data
│   │   │   └── translate_client.go    # 🆕 Google Translate
│   │   │
│   │   └── openai/
│   │       └── ai_client.go           # 🆕 OpenAI/Anthropic
│   │
│   └── cache/                         # 🆕 Enhanced cache
│       └── cache_keys.go              # 🆕 Cache key patterns
│
├── interfaces/api/
│   ├── handlers/
│   │   ├── search_handler.go          # 🆕 NEW
│   │   ├── ai_handler.go              # 🆕 NEW
│   │   ├── folder_handler.go          # 🆕 NEW
│   │   ├── favorite_handler.go        # 🆕 NEW
│   │   └── utility_handler.go         # 🆕 NEW
│   │
│   ├── middleware/
│   │   └── rate_limit_middleware.go   # 🆕 NEW
│   │
│   └── routes/
│       ├── search_routes.go           # 🆕 NEW
│       ├── ai_routes.go               # 🆕 NEW
│       ├── folder_routes.go           # 🆕 NEW
│       ├── favorite_routes.go         # 🆕 NEW
│       └── utility_routes.go          # 🆕 NEW
│
└── pkg/
    └── utils/
        └── qrcode.go                  # 🆕 NEW
```

---

## 3. Development Phases

### Phase 1: Foundation (สัปดาห์ที่ 1)

```
┌─────────────────────────────────────────────────────────────────┐
│                     Phase 1: Foundation                          │
└─────────────────────────────────────────────────────────────────┘

Tasks:
├── 1.1 Update User Model
│   ├── เพิ่ม student_id field
│   └── Update DTOs และ mappers
│
├── 1.2 Setup External API Clients
│   ├── Google Custom Search client
│   ├── Google Places client
│   ├── YouTube Data client
│   └── OpenAI/Anthropic client
│
├── 1.3 Setup Cache Layer
│   ├── Define cache keys
│   └── Implement cache patterns
│
└── 1.4 Update DI Container
    └── Register new dependencies

Files to modify:
├── domain/models/user.go
├── domain/dto/user.go
├── infrastructure/external/ (new folder)
├── pkg/di/container.go
└── pkg/config/config.go
```

### Phase 2: Core Features (สัปดาห์ที่ 2-3)

```
┌─────────────────────────────────────────────────────────────────┐
│                   Phase 2: Core Features                         │
└─────────────────────────────────────────────────────────────────┘

Tasks:
├── 2.1 Search Feature
│   ├── Search service (Google API integration)
│   ├── Search handler (all, website, image, video)
│   ├── Search history repository
│   └── Caching for search results
│
├── 2.2 Folder Feature
│   ├── Folder model & repository
│   ├── Folder item model & repository
│   ├── Folder service
│   └── Folder handler
│
└── 2.3 Favorites Feature
    ├── Favorite model & repository
    ├── Favorite service
    └── Favorite handler

New files:
├── domain/models/folder.go, folder_item.go, favorite.go
├── domain/repositories/folder_repository.go, favorite_repository.go
├── domain/services/search_service.go, folder_service.go, favorite_service.go
├── domain/dto/search.go, folder.go, favorite.go
├── application/serviceimpl/search_service_impl.go, folder_service_impl.go
├── infrastructure/postgres/folder_repository_impl.go, favorite_repository_impl.go
└── interfaces/api/handlers/search_handler.go, folder_handler.go, favorite_handler.go
```

### Phase 3: Advanced Features (สัปดาห์ที่ 4-5)

```
┌─────────────────────────────────────────────────────────────────┐
│                  Phase 3: Advanced Features                      │
└─────────────────────────────────────────────────────────────────┘

Tasks:
├── 3.1 AI Mode Feature
│   ├── AI service (OpenAI integration)
│   ├── AI chat session & messages models
│   ├── AI chat repository
│   ├── AI handler
│   └── YouTube video integration
│
├── 3.2 Places Feature
│   ├── Google Places integration
│   ├── Nearby places search
│   └── Place details
│
└── 3.3 Map Integration
    └── Location-based search

New files:
├── domain/models/ai_chat_session.go, ai_chat_message.go
├── domain/repositories/ai_chat_repository.go
├── domain/services/ai_service.go, place_service.go
├── domain/dto/ai.go, place.go
├── application/serviceimpl/ai_service_impl.go
├── infrastructure/postgres/ai_chat_repository_impl.go
├── infrastructure/external/openai/ai_client.go
└── interfaces/api/handlers/ai_handler.go
```

### Phase 4: Utilities & Polish (สัปดาห์ที่ 6)

```
┌─────────────────────────────────────────────────────────────────┐
│                  Phase 4: Utilities & Polish                     │
└─────────────────────────────────────────────────────────────────┘

Tasks:
├── 4.1 Translation Feature
│   ├── Google Translate integration
│   └── Translation handler
│
├── 4.2 QR Code Feature
│   ├── QR code generator utility
│   └── QR code handler
│
├── 4.3 Rate Limiting
│   ├── Rate limit middleware
│   └── Per-endpoint configuration
│
└── 4.4 Testing & Documentation
    ├── Unit tests
    ├── Integration tests
    └── API documentation (Swagger)

New files:
├── domain/services/translate_service.go, qrcode_service.go
├── domain/dto/utility.go
├── application/serviceimpl/translate_service_impl.go, qrcode_service_impl.go
├── interfaces/api/middleware/rate_limit_middleware.go
├── interfaces/api/handlers/utility_handler.go
└── pkg/utils/qrcode.go
```

---

## 4. API Endpoints Summary

### 4.1 Authentication (มีอยู่แล้ว - ปรับเพิ่ม)

```
POST   /api/v1/auth/register       # ✅ มีอยู่ (เพิ่ม student_id)
POST   /api/v1/auth/login          # ✅ มีอยู่
POST   /api/v1/auth/refresh        # 🆕 NEW
POST   /api/v1/auth/logout         # 🆕 NEW
GET    /api/v1/auth/me             # ✅ มีอยู่ (users/profile)
PUT    /api/v1/auth/me             # ✅ มีอยู่ (users/profile)
```

### 4.2 Search (ใหม่ทั้งหมด)

```
GET    /api/v1/search              # 🆕 Search with Google API
       ?q={query}
       &type={all|website|image|video}
       &page={page}
       &per_page={limit}

GET    /api/v1/search/ai           # 🆕 AI Mode search
       ?q={query}

POST   /api/v1/search/ai/chat      # 🆕 AI Chat (protected)
       body: { session_id, message, image_url }

GET    /api/v1/search/places       # 🆕 Nearby places
       ?lat={lat}&lng={lng}&radius={radius}

GET    /api/v1/search/places/:id   # 🆕 Place details

GET    /api/v1/search/history      # 🆕 Search history (protected)
```

### 4.3 Folders (ใหม่ทั้งหมด)

```
GET    /api/v1/folders             # 🆕 List user's folders
POST   /api/v1/folders             # 🆕 Create folder
GET    /api/v1/folders/:id         # 🆕 Get folder with items
PUT    /api/v1/folders/:id         # 🆕 Update folder
DELETE /api/v1/folders/:id         # 🆕 Delete folder
POST   /api/v1/folders/:id/items   # 🆕 Add item to folder
DELETE /api/v1/folders/:id/items/:itemId  # 🆕 Remove item
POST   /api/v1/folders/:id/share   # 🆕 Generate share link
```

### 4.4 Favorites (ใหม่ทั้งหมด)

```
GET    /api/v1/favorites           # 🆕 List favorites
POST   /api/v1/favorites           # 🆕 Add to favorites
DELETE /api/v1/favorites/:id       # 🆕 Remove from favorites
GET    /api/v1/favorites/check     # 🆕 Check if favorited
       ?type={type}&external_id={id}
```

### 4.5 Utilities (ใหม่ทั้งหมด)

```
POST   /api/v1/translate           # 🆕 Translate text
POST   /api/v1/qrcode              # 🆕 Generate QR code
```

---

## 5. Environment Variables ที่ต้องเพิ่ม

```bash
# .env (เพิ่มเติมจากที่มีอยู่)

# ============================================
# Existing (keep as is)
# ============================================
APP_ENV=development
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=stou_smart_tour
DB_SSL_MODE=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your_jwt_secret

# ============================================
# NEW - Google APIs
# ============================================
GOOGLE_API_KEY=your_google_api_key
GOOGLE_SEARCH_ENGINE_ID=your_search_engine_id

# ============================================
# NEW - OpenAI
# ============================================
OPENAI_API_KEY=your_openai_api_key
OPENAI_MODEL=gpt-4-turbo

# ============================================
# NEW - Rate Limiting
# ============================================
RATE_LIMIT_SEARCH=30          # requests per minute
RATE_LIMIT_AI=10              # requests per minute
RATE_LIMIT_GENERAL=100        # requests per minute

# ============================================
# NEW - Cache TTL (seconds)
# ============================================
CACHE_TTL_SEARCH=3600         # 1 hour
CACHE_TTL_AI=21600            # 6 hours
CACHE_TTL_PLACE=86400         # 24 hours
```

---

## 6. Database Migration Plan

### 6.1 Tables ที่ต้องสร้างใหม่

```sql
-- Run after existing migrations

-- 1. Update users table (add student_id)
ALTER TABLE users ADD COLUMN IF NOT EXISTS student_id VARCHAR(20) UNIQUE;

-- 2. Create folders table
CREATE TABLE folders (...);

-- 3. Create folder_items table
CREATE TABLE folder_items (...);

-- 4. Create favorites table
CREATE TABLE favorites (...);

-- 5. Create search_history table
CREATE TABLE search_history (...);

-- 6. Create ai_chat_sessions table
CREATE TABLE ai_chat_sessions (...);

-- 7. Create ai_chat_messages table
CREATE TABLE ai_chat_messages (...);
```

### 6.2 Migration Files

```
infrastructure/postgres/migrations/
├── 000001_create_users.up.sql        # ✅ มีอยู่
├── 000001_create_users.down.sql
├── 000002_create_tasks.up.sql        # ✅ มีอยู่
├── 000002_create_tasks.down.sql
├── 000003_create_files.up.sql        # ✅ มีอยู่
├── 000003_create_files.down.sql
├── 000004_create_jobs.up.sql         # ✅ มีอยู่
├── 000004_create_jobs.down.sql
├── 000005_add_student_id.up.sql      # 🆕 NEW
├── 000005_add_student_id.down.sql
├── 000006_create_folders.up.sql      # 🆕 NEW
├── 000006_create_folders.down.sql
├── 000007_create_folder_items.up.sql # 🆕 NEW
├── 000007_create_folder_items.down.sql
├── 000008_create_favorites.up.sql    # 🆕 NEW
├── 000008_create_favorites.down.sql
├── 000009_create_search_history.up.sql    # 🆕 NEW
├── 000009_create_search_history.down.sql
├── 000010_create_ai_chat_sessions.up.sql  # 🆕 NEW
├── 000010_create_ai_chat_sessions.down.sql
├── 000011_create_ai_chat_messages.up.sql  # 🆕 NEW
└── 000011_create_ai_chat_messages.down.sql
```

---

## 7. Dependencies ที่ต้องเพิ่ม

```go
// go.mod - เพิ่มเติม

require (
    // Existing dependencies...

    // NEW - QR Code
    github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e

    // NEW - Rate Limiting
    github.com/gofiber/fiber/v2/middleware/limiter

    // NEW - HTTP Client for external APIs
    // (ใช้ net/http มาตรฐานได้)
)
```

---

## 8. ลำดับการพัฒนา (Step by Step)

```
Week 1:
├── Day 1-2: Setup external API clients
├── Day 3-4: Update User model & DTOs
└── Day 5: Update DI Container & Config

Week 2:
├── Day 1-2: Search feature (domain + infrastructure)
├── Day 3-4: Search feature (application + interfaces)
└── Day 5: Testing & bug fixes

Week 3:
├── Day 1-2: Folder feature (domain + infrastructure)
├── Day 3-4: Folder feature (application + interfaces)
└── Day 5: Favorites feature

Week 4:
├── Day 1-2: AI Mode - OpenAI integration
├── Day 3-4: AI Mode - Chat feature
└── Day 5: YouTube integration

Week 5:
├── Day 1-2: Places feature (Google Places)
├── Day 3-4: Map integration
└── Day 5: Location-based search

Week 6:
├── Day 1-2: Translation & QR Code
├── Day 3-4: Rate limiting & Security
└── Day 5: Testing & Documentation
```

---

## Next Part

➡️ ไปต่อที่ **Part 2: Domain Layer (Models, DTOs, Interfaces)**
- รายละเอียด Models ทั้งหมด
- รายละเอียด DTOs ทั้งหมด
- รายละเอียด Repository Interfaces
- รายละเอียด Service Interfaces

---

*Document Version: 1.0*
*Part: 1 of 5*
