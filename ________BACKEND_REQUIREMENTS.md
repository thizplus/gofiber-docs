# Backend Requirements - STOU Smart Tour

## สถานะปัจจุบัน (Current State)

### โครงสร้าง Project
- **Framework:** Go Fiber v2.52
- **Database:** PostgreSQL with GORM
- **Cache:** Redis
- **Storage:** Bunny CDN
- **Architecture:** Clean Architecture (Domain-Driven Design)

### External API Integrations
- **Google Search API** - Web search
- **Google Places API** - Place details
- **Google YouTube API** - Video search
- **Google Translate API** - Translation
- **OpenAI API** - AI Chat (GPT-4)

### API Endpoints ที่มีอยู่แล้ว

#### Authentication
| Method | Endpoint | Auth | สถานะ |
|--------|----------|------|-------|
| POST | `/api/v1/auth/register` | Public | ✅ ใช้งานได้ |
| POST | `/api/v1/auth/login` | Public | ✅ ใช้งานได้ |

#### User Management
| Method | Endpoint | Auth | สถานะ |
|--------|----------|------|-------|
| GET | `/api/v1/users/profile` | Protected | ✅ ใช้งานได้ |
| PUT | `/api/v1/users/profile` | Protected | ✅ ใช้งานได้ |
| DELETE | `/api/v1/users/profile` | Protected | ✅ ใช้งานได้ |
| GET | `/api/v1/users/` | Admin | ✅ ใช้งานได้ |

#### Search
| Method | Endpoint | Auth | สถานะ |
|--------|----------|------|-------|
| GET | `/api/v1/search/` | OptionalAuth | ✅ ใช้งานได้ |
| GET | `/api/v1/search/websites` | OptionalAuth | ✅ ใช้งานได้ |
| GET | `/api/v1/search/images` | OptionalAuth | ✅ ใช้งานได้ |
| GET | `/api/v1/search/videos` | OptionalAuth | ✅ ใช้งานได้ |
| GET | `/api/v1/search/videos/:videoId` | Public | ✅ ใช้งานได้ |
| GET | `/api/v1/search/places` | OptionalAuth | ✅ ใช้งานได้ |
| GET | `/api/v1/search/places/:placeId` | Public | ✅ ใช้งานได้ |
| GET | `/api/v1/search/nearby` | Public | ✅ ใช้งานได้ |
| GET | `/api/v1/search/history` | Protected | ✅ ใช้งานได้ |
| DELETE | `/api/v1/search/history` | Protected | ✅ ใช้งานได้ |
| DELETE | `/api/v1/search/history/:id` | Protected | ✅ ใช้งานได้ |

#### AI
| Method | Endpoint | Auth | สถานะ |
|--------|----------|------|-------|
| GET | `/api/v1/ai/search` | OptionalAuth | ✅ ใช้งานได้ |
| POST | `/api/v1/ai/chat/` | Protected | ✅ ใช้งานได้ |
| GET | `/api/v1/ai/chat/` | Protected | ✅ ใช้งานได้ |
| GET | `/api/v1/ai/chat/:sessionId` | Protected | ✅ ใช้งานได้ |
| POST | `/api/v1/ai/chat/:sessionId/messages` | Protected | ✅ ใช้งานได้ |
| DELETE | `/api/v1/ai/chat/` | Protected | ✅ ใช้งานได้ |
| DELETE | `/api/v1/ai/chat/:sessionId` | Protected | ✅ ใช้งานได้ |

#### Folders
| Method | Endpoint | Auth | สถานะ |
|--------|----------|------|-------|
| GET | `/api/v1/folders/public/:id` | Public | ✅ ใช้งานได้ |
| POST | `/api/v1/folders/` | Protected | ✅ ใช้งานได้ |
| GET | `/api/v1/folders/` | Protected | ✅ ใช้งานได้ |
| GET | `/api/v1/folders/:id` | Protected | ✅ ใช้งานได้ |
| PUT | `/api/v1/folders/:id` | Protected | ✅ ใช้งานได้ |
| DELETE | `/api/v1/folders/:id` | Protected | ✅ ใช้งานได้ |
| POST | `/api/v1/folders/:id/share` | Protected | ✅ ใช้งานได้ |
| POST | `/api/v1/folders/:id/items` | Protected | ✅ ใช้งานได้ |
| GET | `/api/v1/folders/:id/items` | Protected | ✅ ใช้งานได้ |
| PUT | `/api/v1/folders/:id/items/reorder` | Protected | ✅ ใช้งานได้ |
| GET | `/api/v1/folders/items/check` | Protected | ✅ ใช้งานได้ |
| PUT | `/api/v1/folders/items/:itemId` | Protected | ✅ ใช้งานได้ |
| DELETE | `/api/v1/folders/items/:itemId` | Protected | ✅ ใช้งานได้ |

#### Favorites
| Method | Endpoint | Auth | สถานะ |
|--------|----------|------|-------|
| POST | `/api/v1/favorites/` | Protected | ✅ ใช้งานได้ |
| GET | `/api/v1/favorites/` | Protected | ✅ ใช้งานได้ |
| DELETE | `/api/v1/favorites/:id` | Protected | ✅ ใช้งานได้ |
| GET | `/api/v1/favorites/check` | Protected | ✅ ใช้งานได้ |
| POST | `/api/v1/favorites/toggle` | Protected | ✅ ใช้งานได้ |

#### Files
| Method | Endpoint | Auth | สถานะ |
|--------|----------|------|-------|
| POST | `/api/v1/files/upload` | Protected | ✅ ใช้งานได้ |
| GET | `/api/v1/files/` | Admin | ✅ ใช้งานได้ |
| GET | `/api/v1/files/my` | Protected | ✅ ใช้งานได้ |
| GET | `/api/v1/files/:id` | Protected | ✅ ใช้งานได้ |
| DELETE | `/api/v1/files/:id` | Owner | ✅ ใช้งานได้ |

#### Utilities
| Method | Endpoint | Auth | สถานะ |
|--------|----------|------|-------|
| GET | `/api/v1/utils/config` | Public | ✅ ใช้งานได้ |
| POST | `/api/v1/utils/translate` | Public | ✅ ใช้งานได้ |
| POST | `/api/v1/utils/detect-language` | Public | ✅ ใช้งานได้ |
| POST | `/api/v1/utils/qrcode` | Public | ✅ ใช้งานได้ |
| GET | `/api/v1/utils/distance` | Public | ✅ ใช้งานได้ |

### Database Models ที่มีอยู่
- `User` - ผู้ใช้งาน
- `Folder` - โฟลเดอร์
- `FolderItem` - รายการในโฟลเดอร์
- `Favorite` - รายการโปรด
- `SearchHistory` - ประวัติการค้นหา
- `AIChatSession` - เซสชัน AI Chat
- `AIChatMessage` - ข้อความใน AI Chat
- `File` - ไฟล์ที่อัปโหลด
- `Task` - งานที่กำหนดเวลา
- `Job` - งาน scheduled

---

## งานที่ต้องทำ (Requirements from requirement.txt)

### 1. OAuth Authentication (Gmail, LINE, Facebook)
**Priority: HIGH**

**สิ่งที่ต้องทำ:**
- [ ] เพิ่ม OAuth2 flow สำหรับ Google
- [ ] เพิ่ม OAuth2 flow สำหรับ LINE
- [ ] เพิ่ม OAuth2 flow สำหรับ Facebook
- [ ] Link OAuth account กับ existing user (by email)
- [ ] เพิ่มฟิลด์ `oauth_provider` และ `oauth_id` ใน User model

**Endpoints ที่ต้องสร้าง:**
```
POST /api/v1/auth/google
POST /api/v1/auth/line
POST /api/v1/auth/facebook
POST /api/v1/auth/link-oauth  (optional - link existing account)
```

**ไฟล์ที่ต้องสร้าง/แก้ไข:**
- `domain/models/user.go` - เพิ่ม OAuth fields
- สร้าง `infrastructure/external/oauth/google.go`
- สร้าง `infrastructure/external/oauth/line.go`
- สร้าง `infrastructure/external/oauth/facebook.go`
- `interfaces/api/handlers/auth_handler.go` - เพิ่ม OAuth handlers
- `interfaces/api/routes/auth_routes.go` - เพิ่ม routes
- `pkg/config/config.go` - เพิ่ม OAuth config

**User Model Update:**
```go
type User struct {
    // ... existing fields
    OAuthProvider string  `gorm:"type:varchar(20)"` // google, line, facebook
    OAuthID       string  `gorm:"type:varchar(255);index"`
    // ... rest
}
```

**Config ที่ต้องเพิ่ม:**
```go
type OAuthConfig struct {
    Google   GoogleOAuthConfig
    LINE     LINEOAuthConfig
    Facebook FacebookOAuthConfig
}

type GoogleOAuthConfig struct {
    ClientID     string
    ClientSecret string
    RedirectURL  string
}
// ... LINE, Facebook similar
```

**Environment Variables:**
```env
# Google OAuth
GOOGLE_OAUTH_CLIENT_ID=
GOOGLE_OAUTH_CLIENT_SECRET=
GOOGLE_OAUTH_REDIRECT_URL=

# LINE OAuth
LINE_CHANNEL_ID=
LINE_CHANNEL_SECRET=
LINE_REDIRECT_URL=

# Facebook OAuth
FACEBOOK_APP_ID=
FACEBOOK_APP_SECRET=
FACEBOOK_REDIRECT_URL=
```

---

### 2. AI Place Summary Endpoint
**Priority: HIGH**

**สิ่งที่ต้องทำ:**
- [ ] สร้าง endpoint ที่รับ Place ID และส่งคืน AI summary
- [ ] ใช้ข้อมูลจาก Google Places + Reviews ในการสร้าง summary
- [ ] Cache AI summary ใน Redis (เพื่อประหยัด API calls)
- [ ] รองรับภาษาไทยและอังกฤษ

**Endpoint ที่ต้องสร้าง:**
```
GET /api/v1/places/:placeId/ai-summary?lang=th
```

**ไฟล์ที่ต้องแก้ไข:**
- สร้าง handler ใน `interfaces/api/handlers/search_handler.go` หรือแยกไฟล์ใหม่
- `interfaces/api/routes/search_routes.go` - เพิ่ม route
- `application/serviceimpl/ai_service_impl.go` - เพิ่ม method

**Response Format:**
```json
{
  "success": true,
  "data": {
    "placeId": "ChIJX49VWsWZ4jARwEYBM0MDXnE",
    "summary": "วัดพระศรีรัตนศาสดาราม หรือ วัดพระแก้ว เป็นวัดสำคัญที่สุดในประเทศไทย...",
    "history": "สร้างขึ้นในปี พ.ศ. 2325 ในรัชกาลที่ 1...",
    "highlights": ["พระแก้วมรกต", "จิตรกรรมฝาผนัง", "พระบรมมหาราชวัง"],
    "tips": ["ควรแต่งกายสุภาพ", "เปิดทุกวัน 8:30-15:30"],
    "language": "th",
    "cachedAt": "2024-01-15T10:30:00Z"
  }
}
```

---

### 3. Search Suggestions/Keywords Endpoint
**Priority: MEDIUM**

**สิ่งที่ต้องทำ:**
- [ ] สร้าง endpoint สำหรับ search suggestions
- [ ] ดึงข้อมูลจาก popular searches
- [ ] ดึงข้อมูลจาก user's recent searches
- [ ] Cache popular keywords

**Endpoint ที่ต้องสร้าง:**
```
GET /api/v1/search/suggestions?q={query}&limit=10
```

**ไฟล์ที่ต้องแก้ไข:**
- `interfaces/api/handlers/search_handler.go` - เพิ่ม handler
- `interfaces/api/routes/search_routes.go` - เพิ่ม route
- `application/serviceimpl/search_service_impl.go` - เพิ่ม method
- `domain/repositories/search_history_repository.go` - เพิ่ม method

**Response Format:**
```json
{
  "success": true,
  "data": {
    "suggestions": [
      {"text": "วัดพระแก้ว", "type": "popular"},
      {"text": "วัดอรุณ", "type": "recent"},
      {"text": "พระบรมมหาราชวัง", "type": "popular"}
    ]
  }
}
```

**Database Query:**
- นับ search history ที่มี query คล้ายกัน
- เรียงตามความถี่
- Filter by prefix matching

---

### 4. Cloudflare R2 Storage (Optional - แทน Bunny CDN)
**Priority: LOW**

**หมายเหตุ:** ปัจจุบันใช้ Bunny Storage อยู่แล้วและทำงานได้ดี การเปลี่ยนไป Cloudflare R2 เป็น optional

**ถ้าต้องการเปลี่ยน:**
- [ ] สร้าง R2 storage client
- [ ] แก้ไข file upload service
- [ ] Migrate existing files (optional)

**ไฟล์ที่ต้องสร้าง:**
- `infrastructure/storage/r2_storage.go`

**Config ที่ต้องเพิ่ม:**
```env
CLOUDFLARE_ACCOUNT_ID=
CLOUDFLARE_R2_ACCESS_KEY_ID=
CLOUDFLARE_R2_SECRET_ACCESS_KEY=
CLOUDFLARE_R2_BUCKET_NAME=
CLOUDFLARE_R2_PUBLIC_URL=
```

---

### 5. Place History/Background Data (Optional Enhancement)
**Priority: LOW**

**สิ่งที่ต้องทำ:**
- [ ] เพิ่ม endpoint สำหรับดึงข้อมูลประวัติสถานที่
- [ ] อาจใช้ Wikipedia API หรือ AI generate
- [ ] Cache ข้อมูลใน database

**Endpoint ที่อาจสร้าง:**
```
GET /api/v1/places/:placeId/history
```

**หมายเหตุ:** สามารถรวมเข้ากับ AI Summary endpoint ได้

---

## Database Migrations ที่ต้องทำ

### 1. User Model - OAuth Support
```sql
ALTER TABLE users ADD COLUMN oauth_provider VARCHAR(20);
ALTER TABLE users ADD COLUMN oauth_id VARCHAR(255);
CREATE INDEX idx_users_oauth ON users(oauth_provider, oauth_id);
```

### 2. Popular Searches Table (Optional)
```sql
CREATE TABLE popular_searches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query VARCHAR(255) NOT NULL,
    search_count INT DEFAULT 1,
    last_searched_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(query)
);

CREATE INDEX idx_popular_searches_count ON popular_searches(search_count DESC);
```

---

## สรุป Priority

| Priority | งาน | ความยาก | Dependencies |
|----------|-----|---------|--------------|
| **HIGH** | 1. OAuth Login (Google, LINE, Facebook) | Hard | OAuth credentials |
| **HIGH** | 2. AI Place Summary | Medium | OpenAI API |
| **MEDIUM** | 3. Search Suggestions | Easy | None |
| **LOW** | 4. Cloudflare R2 | Medium | R2 credentials |
| **LOW** | 5. Place History | Medium | Wikipedia API / AI |

---

## Environment Variables ที่ต้องเพิ่ม

```env
# OAuth - Google
GOOGLE_OAUTH_CLIENT_ID=
GOOGLE_OAUTH_CLIENT_SECRET=
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:3000/api/auth/callback/google

# OAuth - LINE
LINE_CHANNEL_ID=
LINE_CHANNEL_SECRET=
LINE_REDIRECT_URL=http://localhost:3000/api/auth/callback/line

# OAuth - Facebook
FACEBOOK_APP_ID=
FACEBOOK_APP_SECRET=
FACEBOOK_REDIRECT_URL=http://localhost:3000/api/auth/callback/facebook

# Cloudflare R2 (Optional)
CLOUDFLARE_ACCOUNT_ID=
CLOUDFLARE_R2_ACCESS_KEY_ID=
CLOUDFLARE_R2_SECRET_ACCESS_KEY=
CLOUDFLARE_R2_BUCKET_NAME=
CLOUDFLARE_R2_PUBLIC_URL=
```

---

## Recommended Implementation Order

### Phase 1 (ทำได้เลย - ไม่ต้อง setup เพิ่ม):
1. **Search Suggestions** - ใช้ข้อมูล search history ที่มีอยู่
2. **AI Place Summary** - ใช้ OpenAI ที่ setup ไว้แล้ว

### Phase 2 (ต้อง setup credentials):
3. **Google OAuth** - ต้องสร้าง Google Cloud Console project
4. **LINE OAuth** - ต้องสร้าง LINE Developers account
5. **Facebook OAuth** - ต้องสร้าง Meta for Developers app

### Phase 3 (Optional):
6. **Cloudflare R2** - ถ้าต้องการย้ายจาก Bunny CDN
7. **Place History** - ถ้าต้องการข้อมูลประวัติละเอียด

---

## API Design Guidelines

1. **Consistency:** ใช้ response format เดียวกันทั้ง project
   ```json
   {
     "success": true/false,
     "message": "...",
     "data": { ... },
     "error": "..." // only when success: false
   }
   ```

2. **Caching:** ใช้ Redis cache สำหรับข้อมูลที่ query บ่อย
   - AI summaries: 24 hours
   - Search suggestions: 1 hour
   - Place details: 6 hours

3. **Rate Limiting:** ใช้ rate limit ที่มีอยู่แล้ว
   - Search: 30 req/min
   - AI: 10 req/min
   - General: 100 req/min

4. **Error Handling:** Return meaningful error messages
