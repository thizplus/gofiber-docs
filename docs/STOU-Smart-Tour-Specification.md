# STOU Smart Tour - Project Specification

## Overview
ระบบค้นหาข้อมูลท่องเที่ยวสำหรับนักศึกษา มหาวิทยาลัยสุโขทัยธรรมาธิราช (มสธ.)  
ใช้ Google Search API เป็น backend หลักในการค้นหาข้อมูล

---

## Tech Stack
- **Frontend:** Next.js (React + TypeScript)
- **Backend:** Go Fiber
- **Search API:** Google Custom Search API
- **Database:** PostgreSQL (สำหรับ user data, folders, favorites)

---

## Pages & Features

### 1. Authentication Pages

#### 1.1 Login Page (`/login`)
- Login form (email/password หรือ student ID)
- Social login options (optional)
- Link ไป Register page
- "ลืมรหัสผ่าน" link

#### 1.2 Register Page (`/register`)
- Registration form สำหรับนักศึกษา
- Fields: รหัสนักศึกษา, ชื่อ-นามสกุล, email, password

---

### 2. Main Pages

#### 2.1 Home Page (`/`)
- **Hero Section:** แสดงภาพและข้อความแนะนำ
- **Search Bar:** ช่องค้นหาหลัก พร้อม placeholder "พิมพ์ชื่อจังหวัด/สถานที่"
- **Filter Tabs:** AI Mode | All | Website | Image | Video | Map
- **Featured Destinations:** แสดงสถานที่ท่องเที่ยวยอดนิยม
- **Quick Links:** ลิงก์ไปยังสถานที่สำคัญ เช่น ตลาดจตุจักร, สวนลุมพินี, วัดโพธิ์

#### 2.2 Search Results Page (`/search`)
- **Search Bar:** คงที่ด้านบน พร้อม keyword ที่ค้นหา
- **Filter Tabs:** AI Mode | All | Website | Image | Video | Map
- **Results Display:**
  - **Card View:** แสดงผลแบบ grid สำหรับรูปภาพ/สถานที่
  - **List View:** แสดงผลแบบรายการสำหรับเว็บไซต์
- **Result Card Components:**
  - Thumbnail image
  - ชื่อสถานที่
  - Rating (ดาว) + จำนวน reviews
  - ระยะทาง (กม.)
  - ประเภท (ร้านอาหาร, สถานที่ท่องเที่ยว, etc.)
  - ปุ่ม Favorite (heart icon)
  - ปุ่ม Share
  - ปุ่ม Save to Folder

#### 2.3 AI Mode Page (`/search?mode=ai`)
- **AI-Generated Content:** 
  - สรุปข้อมูลจาก AI เกี่ยวกับหัวข้อที่ค้นหา
  - แสดงเป็น structured content (หัวข้อ, bullet points)
  - มี source links (clickable)
- **Related Videos:** แสดง YouTube videos ที่เกี่ยวข้อง
- **Chat Input:** ช่องสำหรับถามคำถามเพิ่มเติม
  - รองรับ text input
  - รองรับ image upload
  - รองรับ voice input (microphone icon)

#### 2.4 Place Detail Page (`/place/[id]`)
- รายละเอียดสถานที่
- รูปภาพ gallery
- Rating & Reviews
- Map location
- สถานที่ใกล้เคียง
- ร้านอาหารใกล้เคียง

---

### 3. User Features

#### 3.1 My Folder Page (`/my-folder`)
- **Folder List:** แสดง folders ที่ user สร้าง
- **Folder Features:**
  - สร้าง folder ใหม่
  - ตั้งชื่อ folder (เช่น "กรุงเทพฯ", "เชียงใหม่")
  - เก็บได้: ไฟล์, ภาพ, วิดีโอ, ลิงก์
  - แชร์ folder ได้
- **Folder Detail:** แสดงรายการที่เก็บไว้ใน folder

#### 3.2 Profile Page (`/profile`)
- ข้อมูลส่วนตัว
- แก้ไขข้อมูล
- ประวัติการค้นหา
- Settings

---

### 4. Utility Features

#### 4.1 Translation (แปลภาษา)
- แปลข้อมูลที่ค้นหาเป็นภาษาอื่น
- รองรับหลายภาษา (EN, TH, etc.)

#### 4.2 QR Code Generator (แปลง QR Code)
- สร้าง QR Code จาก URL หรือข้อมูลที่เลือก
- แชร์สถานที่ผ่าน QR Code

#### 4.3 Language Switcher
- เปลี่ยนภาษา interface (TH/EN)

---

## Navigation Structure

```
Header Navigation:
├── Logo (STOU Smart Tour)
├── Home
├── My Folder (ต้อง login)
├── Project (optional)
├── Contact
├── AI Mode
├── แปลภาษา
├── แปลง QR Code
├── Profile (ต้อง login)
└── Sign In / Sign Out
```

---

## Components Structure

```
components/
├── layout/
│   ├── Header.tsx
│   ├── Footer.tsx
│   ├── Sidebar.tsx
│   └── MainLayout.tsx
├── search/
│   ├── SearchBar.tsx
│   ├── FilterTabs.tsx
│   ├── SearchResults.tsx
│   ├── ResultCard.tsx
│   ├── ResultList.tsx
│   └── MapView.tsx
├── ai/
│   ├── AIContent.tsx
│   ├── AIChat.tsx
│   ├── RelatedVideos.tsx
│   └── SourceLinks.tsx
├── folder/
│   ├── FolderList.tsx
│   ├── FolderCard.tsx
│   ├── FolderDetail.tsx
│   └── SaveToFolderModal.tsx
├── place/
│   ├── PlaceCard.tsx
│   ├── PlaceDetail.tsx
│   ├── NearbyPlaces.tsx
│   └── NearbyRestaurants.tsx
├── common/
│   ├── Button.tsx
│   ├── Modal.tsx
│   ├── Rating.tsx
│   ├── FavoriteButton.tsx
│   ├── ShareButton.tsx
│   └── QRCodeGenerator.tsx
└── auth/
    ├── LoginForm.tsx
    ├── RegisterForm.tsx
    └── AuthGuard.tsx
```

---

## API Endpoints (Go Fiber Backend)

### Authentication
```
POST   /api/auth/login          - Login
POST   /api/auth/register       - Register
POST   /api/auth/logout         - Logout
GET    /api/auth/me             - Get current user
```

### Search
```
GET    /api/search              - Search (Google API wrapper)
       ?q={query}
       &type={all|website|image|video}
       &location={lat,lng}
       
GET    /api/search/ai           - AI Mode search
       ?q={query}
       
GET    /api/search/places       - Search places nearby
       ?lat={lat}&lng={lng}&radius={km}
```

### Folders
```
GET    /api/folders             - Get user's folders
POST   /api/folders             - Create folder
GET    /api/folders/:id         - Get folder detail
PUT    /api/folders/:id         - Update folder
DELETE /api/folders/:id         - Delete folder
POST   /api/folders/:id/items   - Add item to folder
DELETE /api/folders/:id/items/:itemId - Remove item
```

### Favorites
```
GET    /api/favorites           - Get user's favorites
POST   /api/favorites           - Add to favorites
DELETE /api/favorites/:id       - Remove from favorites
```

### Utilities
```
POST   /api/translate           - Translate text
POST   /api/qrcode              - Generate QR Code
```

---

## Database Schema (PostgreSQL)

### Tables

```sql
-- Users
users (
  id UUID PRIMARY KEY,
  student_id VARCHAR(20) UNIQUE,
  email VARCHAR(255) UNIQUE,
  password_hash VARCHAR(255),
  name VARCHAR(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
)

-- Folders
folders (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  name VARCHAR(255),
  description TEXT,
  is_public BOOLEAN DEFAULT false,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
)

-- Folder Items
folder_items (
  id UUID PRIMARY KEY,
  folder_id UUID REFERENCES folders(id),
  type VARCHAR(50), -- 'place', 'website', 'image', 'video', 'link'
  title VARCHAR(255),
  url TEXT,
  thumbnail_url TEXT,
  metadata JSONB,
  created_at TIMESTAMP
)

-- Favorites
favorites (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  type VARCHAR(50),
  external_id VARCHAR(255), -- Google Place ID, etc.
  title VARCHAR(255),
  url TEXT,
  thumbnail_url TEXT,
  rating DECIMAL(2,1),
  metadata JSONB,
  created_at TIMESTAMP
)

-- Search History
search_history (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  query VARCHAR(255),
  search_type VARCHAR(50),
  created_at TIMESTAMP
)
```

---

## Key Features Summary

| Feature | Description | Priority |
|---------|-------------|----------|
| Search (All) | ค้นหาทุกประเภทจาก Google | High |
| Search (Website) | ค้นหาเฉพาะเว็บไซต์ | High |
| Search (Image) | ค้นหารูปภาพ | High |
| Search (Video) | ค้นหาวิดีโอ | High |
| Search (Map) | แสดงผลบน Map | Medium |
| AI Mode | ค้นหาพร้อมสรุปจาก AI | High |
| My Folder | เก็บข้อมูลที่สนใจ | High |
| Favorites | กดหัวใจเก็บสถานที่ | Medium |
| Share | แชร์ข้อมูล | Medium |
| QR Code | สร้าง QR Code | Low |
| Translation | แปลภาษา | Low |
| Nearby Places | สถานที่ใกล้เคียง | Medium |

---

## UI/UX Notes

### Color Scheme
- Primary: Blue (#3B82F6 หรือใกล้เคียง)
- Background: Light blue gradient / Sky theme
- Accent: Green (for highlights, links)
- Text: Dark gray / Black

### Design Elements
- Rounded corners (cards, buttons, search bar)
- Soft shadows
- Clean, modern look
- Sky/cloud background theme
- TripAdvisor-style badges (Travelers' Choice 2025)

### Responsive
- Desktop: Full navigation, grid view
- Tablet: Collapsible menu, adjusted grid
- Mobile: Hamburger menu, single column

---

## Third-Party Integrations

1. **Google Custom Search API** - ค้นหาข้อมูล
2. **Google Places API** - ข้อมูลสถานที่, rating, reviews
3. **Google Maps API** - แสดง Map
4. **Google Translate API** - แปลภาษา (optional)
5. **OpenAI / Anthropic API** - AI Mode summarization
6. **YouTube Data API** - ดึงวิดีโอที่เกี่ยวข้อง

---

## Development Phases

### Phase 1: Foundation
- [ ] Setup Next.js project
- [ ] Setup Go Fiber backend
- [ ] Authentication system
- [ ] Basic search functionality

### Phase 2: Core Features
- [ ] Search results pages (All, Website, Image, Video)
- [ ] Filter tabs
- [ ] Result cards & list views
- [ ] My Folder feature

### Phase 3: Advanced Features
- [ ] AI Mode integration
- [ ] Map view
- [ ] Nearby places
- [ ] Favorites system

### Phase 4: Polish
- [ ] QR Code generator
- [ ] Translation feature
- [ ] Share functionality
- [ ] Performance optimization
- [ ] Mobile responsive

---

## Notes
- ระบบนี้เป็นระบบสำหรับนักศึกษา มสธ. โดยเฉพาะ
- ควรมี rate limiting สำหรับ Google API calls
- ควรมี caching สำหรับ frequent searches
- พิจารณา offline support สำหรับข้อมูลที่เคยค้นหา
