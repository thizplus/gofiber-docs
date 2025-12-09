# Part 2: Search APIs (Website, Image, Video)

## Overview
ระบบค้นหาต่างๆ รวมถึง Website, รูปภาพ, และวิดีโอ YouTube

## Base URL
```
/api/v1/search
```

---

## 2.1 General Search (ค้นหาทั่วไป)

### Endpoint
```
GET /api/v1/search
```

### Authentication
Optional (ถ้า login จะบันทึกประวัติการค้นหา)

### Query Parameters
```typescript
interface SearchRequest {
  q: string;         // required, 1-500 chars - คำค้นหา
  type?: string;     // optional: "all" | "website" | "image" | "video" | "map" | "ai"
  page?: number;     // optional, min 1 - หน้าที่ต้องการ
  pageSize?: number; // optional, 1-50 - จำนวนผลลัพธ์ต่อหน้า
  lang?: string;     // optional, 2 chars (e.g., "th", "en")
}
```

### Example Request
```
GET /api/v1/search?q=ที่เที่ยวกรุงเทพ&type=all&page=1&pageSize=10
```

### Response
```typescript
interface SearchResponse {
  query: string;
  type: string;
  results: SearchResult[];
  totalCount: number;
  page: number;
  pageSize: number;
}

interface SearchResult {
  type: string;           // "website" | "image" | "video" | "place"
  title: string;
  url: string;
  snippet?: string;
  thumbnailUrl?: string;
  source?: string;
  publishedAt?: string;
  rating?: number;
  reviewCount?: number;
}
```

---

## 2.2 Website Search (ค้นหาเว็บไซต์)

### Endpoint
```
GET /api/v1/search/websites
```

### Authentication
Optional

### Query Parameters
```typescript
interface SearchRequest {
  q: string;         // required - คำค้นหา
  page?: number;     // optional, default 1
  pageSize?: number; // optional, default 10, max 50
}
```

### Example Request
```
GET /api/v1/search/websites?q=มหาวิทยาลัยสุโขทัยธรรมาธิราช&page=1&pageSize=10
```

### Response
```typescript
interface WebsiteSearchResponse {
  query: string;
  results: WebsiteResult[];
  totalCount: number;
  page: number;
  pageSize: number;
}

interface WebsiteResult {
  title: string;       // หัวข้อเว็บ
  url: string;         // URL ของเว็บ
  snippet: string;     // ข้อความสรุป
  displayLink: string; // ชื่อโดเมน
  formattedAt?: string; // วันที่ format แล้ว
}
```

### Example Response
```json
{
  "success": true,
  "message": "Website search completed",
  "data": {
    "query": "มหาวิทยาลัยสุโขทัยธรรมาธิราช",
    "results": [
      {
        "title": "มหาวิทยาลัยสุโขทัยธรรมาธิราช - STOU",
        "url": "https://www.stou.ac.th/",
        "snippet": "มหาวิทยาลัยสุโขทัยธรรมาธิราช เป็นมหาวิทยาลัยเปิดแห่งแรกของประเทศไทย...",
        "displayLink": "www.stou.ac.th"
      }
    ],
    "totalCount": 1000,
    "page": 1,
    "pageSize": 10
  }
}
```

---

## 2.3 Image Search (ค้นหารูปภาพ)

### Endpoint
```
GET /api/v1/search/images
```

### Authentication
Optional

### Query Parameters
```typescript
interface ImageSearchRequest {
  q: string;         // required - คำค้นหา
  page?: number;     // optional, default 1
  pageSize?: number; // optional, default 10, max 50
  size?: string;     // optional: "small" | "medium" | "large"
}
```

### Example Request
```
GET /api/v1/search/images?q=วัดพระแก้ว&page=1&pageSize=10
```

### Response
```typescript
interface ImageSearchResponse {
  query: string;
  results: ImageResult[];
  totalCount: number;
  page: number;
  pageSize: number;
}

interface ImageResult {
  title: string;        // ชื่อรูปภาพ
  url: string;          // URL รูปภาพขนาดเต็ม
  thumbnailUrl: string; // URL รูปภาพ thumbnail
  width: number;        // ความกว้าง (pixels)
  height: number;       // ความสูง (pixels)
  source: string;       // แหล่งที่มา
  contextLink: string;  // URL หน้าเว็บที่รูปภาพอยู่
}
```

### Example Response
```json
{
  "success": true,
  "message": "Image search completed",
  "data": {
    "query": "วัดพระแก้ว",
    "results": [
      {
        "title": "วัดพระศรีรัตนศาสดาราม",
        "url": "https://example.com/images/wat-phra-kaew-full.jpg",
        "thumbnailUrl": "https://example.com/images/wat-phra-kaew-thumb.jpg",
        "width": 1920,
        "height": 1080,
        "source": "thailand-tourism.com",
        "contextLink": "https://thailand-tourism.com/wat-phra-kaew"
      }
    ],
    "totalCount": 500,
    "page": 1,
    "pageSize": 10
  }
}
```

---

## 2.4 Video Search (ค้นหาวิดีโอ YouTube)

### Endpoint
```
GET /api/v1/search/videos
```

### Authentication
Optional

### Query Parameters
```typescript
interface VideoSearchRequest {
  q: string;         // required - คำค้นหา
  page?: number;     // optional, default 1
  pageSize?: number; // optional, default 10, max 50
  order?: string;    // optional: "relevance" | "date" | "viewCount" | "rating"
}
```

### Example Request
```
GET /api/v1/search/videos?q=อาหารกรุงเทพ&page=1&pageSize=10&order=viewCount
```

### Response
```typescript
interface VideoSearchResponse {
  query: string;
  results: VideoResult[];
  totalCount: number;
  page: number;
  pageSize: number;
}

interface VideoResult {
  videoId: string;       // YouTube Video ID
  title: string;         // ชื่อวิดีโอ
  description: string;   // คำอธิบาย
  thumbnailUrl: string;  // URL รูป thumbnail
  channelTitle: string;  // ชื่อช่อง
  publishedAt: string;   // วันที่เผยแพร่
  duration?: string;     // ความยาว (e.g., "PT10M30S")
  viewCount?: number;    // จำนวนการรับชม
  likeCount?: number;    // จำนวนไลค์
}
```

### Example Response
```json
{
  "success": true,
  "message": "Video search completed",
  "data": {
    "query": "อาหารกรุงเทพ",
    "results": [
      {
        "videoId": "dQw4w9WgXcQ",
        "title": "10 ร้านอาหารกรุงเทพต้องไป",
        "description": "รวม 10 ร้านอาหารที่ดีที่สุดในกรุงเทพมหานคร...",
        "thumbnailUrl": "https://img.youtube.com/vi/dQw4w9WgXcQ/hqdefault.jpg",
        "channelTitle": "Food Explorer",
        "publishedAt": "2024-01-10T12:00:00Z",
        "duration": "PT15M30S",
        "viewCount": 150000,
        "likeCount": 5000
      }
    ],
    "totalCount": 200,
    "page": 1,
    "pageSize": 10
  }
}
```

---

## 2.5 Get Video Details (รายละเอียดวิดีโอ)

### Endpoint
```
GET /api/v1/search/videos/:videoId
```

### Authentication
ไม่ต้อง (Public)

### Path Parameters
- `videoId`: YouTube Video ID

### Example Request
```
GET /api/v1/search/videos/dQw4w9WgXcQ
```

### Response
```json
{
  "success": true,
  "message": "Video details retrieved",
  "data": {
    "videoId": "dQw4w9WgXcQ",
    "title": "10 ร้านอาหารกรุงเทพต้องไป",
    "description": "รวม 10 ร้านอาหารที่ดีที่สุดในกรุงเทพมหานคร...",
    "thumbnailUrl": "https://img.youtube.com/vi/dQw4w9WgXcQ/maxresdefault.jpg",
    "channelTitle": "Food Explorer",
    "publishedAt": "2024-01-10T12:00:00Z",
    "duration": "PT15M30S",
    "viewCount": 150000,
    "likeCount": 5000
  }
}
```

---

## 2.6 Search History (ประวัติการค้นหา)

### Get Search History
```
GET /api/v1/search/history
```

### Authentication
Required (Bearer Token)

### Query Parameters
```typescript
interface GetSearchHistoryRequest {
  type?: string;     // optional: "all" | "website" | "image" | "video" | "map" | "ai"
  page?: number;     // optional, default 1
  pageSize?: number; // optional, default 20, max 50
}
```

### Response
```typescript
interface SearchHistoryListResponse {
  histories: SearchHistoryResponse[];
  meta: PaginationMeta;
}

interface SearchHistoryResponse {
  id: string;         // UUID
  query: string;      // คำค้นหา
  searchType: string; // ประเภทการค้นหา
  resultCount: number; // จำนวนผลลัพธ์
  createdAt: string;  // วันที่ค้นหา
}

interface PaginationMeta {
  total: number;
  offset: number;
  limit: number;
}
```

### Example Response
```json
{
  "success": true,
  "message": "Search history retrieved",
  "data": {
    "histories": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "query": "ที่เที่ยวกรุงเทพ",
        "searchType": "website",
        "resultCount": 100,
        "createdAt": "2024-01-15T10:30:00Z"
      }
    ],
    "meta": {
      "total": 50,
      "offset": 0,
      "limit": 20
    }
  }
}
```

### Clear All Search History
```
DELETE /api/v1/search/history
```

### Delete Single History Item
```
DELETE /api/v1/search/history/:id
```

---

## TypeScript Types สำหรับ Frontend

```typescript
// types/search.ts

export interface SearchRequest {
  q: string;
  type?: 'all' | 'website' | 'image' | 'video' | 'map' | 'ai';
  page?: number;
  pageSize?: number;
  lang?: string;
}

export interface WebsiteResult {
  title: string;
  url: string;
  snippet: string;
  displayLink: string;
  formattedAt?: string;
}

export interface ImageResult {
  title: string;
  url: string;
  thumbnailUrl: string;
  width: number;
  height: number;
  source: string;
  contextLink: string;
}

export interface VideoResult {
  videoId: string;
  title: string;
  description: string;
  thumbnailUrl: string;
  channelTitle: string;
  publishedAt: string;
  duration?: string;
  viewCount?: number;
  likeCount?: number;
}

export interface WebsiteSearchResponse {
  query: string;
  results: WebsiteResult[];
  totalCount: number;
  page: number;
  pageSize: number;
}

export interface ImageSearchResponse {
  query: string;
  results: ImageResult[];
  totalCount: number;
  page: number;
  pageSize: number;
}

export interface VideoSearchResponse {
  query: string;
  results: VideoResult[];
  totalCount: number;
  page: number;
  pageSize: number;
}

export interface SearchHistory {
  id: string;
  query: string;
  searchType: string;
  resultCount: number;
  createdAt: string;
}

export interface PaginationMeta {
  total: number;
  offset: number;
  limit: number;
}
```

---

## API Routes Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/search` | Optional | ค้นหาทั่วไป |
| GET | `/api/v1/search/websites` | Optional | ค้นหาเว็บไซต์ |
| GET | `/api/v1/search/images` | Optional | ค้นหารูปภาพ |
| GET | `/api/v1/search/videos` | Optional | ค้นหาวิดีโอ YouTube |
| GET | `/api/v1/search/videos/:videoId` | No | รายละเอียดวิดีโอ |
| GET | `/api/v1/search/history` | Yes | ดูประวัติการค้นหา |
| DELETE | `/api/v1/search/history` | Yes | ลบประวัติทั้งหมด |
| DELETE | `/api/v1/search/history/:id` | Yes | ลบประวัติรายการเดียว |

---

## Notes
- การค้นหาทุกประเภทมี Redis caching เพื่อลด API calls
- ถ้า login แล้วจะบันทึกประวัติการค้นหาอัตโนมัติ
- ผลลัพธ์ถูก cache ไว้ 1 ชั่วโมง (Website/Image), 6 ชั่วโมง (YouTube)
- Video ID สามารถใช้สร้าง URL YouTube: `https://www.youtube.com/watch?v=${videoId}`
- Thumbnail YouTube URL pattern: `https://img.youtube.com/vi/${videoId}/hqdefault.jpg`
