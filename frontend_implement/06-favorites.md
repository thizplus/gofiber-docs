# Part 6: Favorites APIs (รายการโปรด)

## Overview
ระบบจัดการรายการโปรดสำหรับบันทึกสถานที่ เว็บ รูป หรือวิดีโอที่ชื่นชอบ

## Base URL
```
/api/v1/favorites
```

## Authentication
**Required** - ทุก endpoint ต้อง login

---

## 6.1 Add Favorite (เพิ่มรายการโปรด)

### Endpoint
```
POST /api/v1/favorites
```

### Request Body
```typescript
interface AddFavoriteRequest {
  type: 'place' | 'website' | 'image' | 'video';
  externalId?: string;    // optional - Google Place ID, YouTube Video ID
  title: string;          // required, 1-255 chars
  url: string;            // required, URL format, max 2000 chars
  thumbnailUrl?: string;  // optional, URL format
  rating?: number;        // optional, 0-5
  reviewCount?: number;   // optional
  address?: string;       // optional, max 500 chars
  metadata?: Record<string, any>;  // optional - ข้อมูลเพิ่มเติม
}
```

### Example Request (Place)
```json
{
  "type": "place",
  "externalId": "ChIJ5Wl37g6Z4jARiP4itarBPDQ",
  "title": "วัดพระศรีรัตนศาสดาราม",
  "url": "https://maps.google.com/?cid=123456",
  "thumbnailUrl": "https://example.com/photo.jpg",
  "rating": 4.7,
  "reviewCount": 40683,
  "address": "ถ. หน้าพระลาน กรุงเทพมหานคร",
  "metadata": {
    "lat": 13.7516435,
    "lng": 100.4927041,
    "types": ["tourist_attraction", "place_of_worship"]
  }
}
```

### Example Request (Video)
```json
{
  "type": "video",
  "externalId": "dQw4w9WgXcQ",
  "title": "10 ที่เที่ยวกรุงเทพ 2024",
  "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "thumbnailUrl": "https://img.youtube.com/vi/dQw4w9WgXcQ/hqdefault.jpg",
  "metadata": {
    "channelTitle": "Travel Thailand",
    "viewCount": 150000,
    "duration": "PT15M30S"
  }
}
```

### Example Request (Website)
```json
{
  "type": "website",
  "title": "มหาวิทยาลัยสุโขทัยธรรมาธิราช",
  "url": "https://www.stou.ac.th/",
  "thumbnailUrl": "https://www.stou.ac.th/logo.png",
  "metadata": {
    "displayLink": "www.stou.ac.th",
    "snippet": "มหาวิทยาลัยเปิดแห่งแรกของประเทศไทย"
  }
}
```

### Example Request (Image)
```json
{
  "type": "image",
  "title": "วัดอรุณราชวราราม",
  "url": "https://example.com/images/wat-arun-full.jpg",
  "thumbnailUrl": "https://example.com/images/wat-arun-thumb.jpg",
  "metadata": {
    "width": 1920,
    "height": 1080,
    "source": "thailand-photos.com"
  }
}
```

### Response
```typescript
interface FavoriteResponse {
  id: string;              // UUID
  type: string;
  externalId?: string;
  title: string;
  url: string;
  thumbnailUrl?: string;
  rating?: number;
  reviewCount?: number;
  address?: string;
  metadata?: Record<string, any>;
  createdAt: string;
}
```

### Example Response
```json
{
  "success": true,
  "message": "Favorite added successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "place",
    "externalId": "ChIJ5Wl37g6Z4jARiP4itarBPDQ",
    "title": "วัดพระศรีรัตนศาสดาราม",
    "url": "https://maps.google.com/?cid=123456",
    "thumbnailUrl": "https://example.com/photo.jpg",
    "rating": 4.7,
    "reviewCount": 40683,
    "address": "ถ. หน้าพระลาน กรุงเทพมหานคร",
    "metadata": {
      "lat": 13.7516435,
      "lng": 100.4927041,
      "types": ["tourist_attraction", "place_of_worship"]
    },
    "createdAt": "2024-01-15T10:30:00Z"
  }
}
```

---

## 6.2 Get Favorites (รายการโปรดทั้งหมด)

### Endpoint
```
GET /api/v1/favorites
```

### Query Parameters
```typescript
interface GetFavoritesRequest {
  type?: 'place' | 'website' | 'image' | 'video';  // optional - filter by type
  page?: number;     // optional, default 1
  pageSize?: number; // optional, default 20, max 50
}
```

### Example Requests
```
GET /api/v1/favorites                     # ทั้งหมด
GET /api/v1/favorites?type=place          # เฉพาะสถานที่
GET /api/v1/favorites?type=video&page=2   # เฉพาะวิดีโอ หน้าที่ 2
```

### Response
```typescript
interface FavoriteListResponse {
  favorites: FavoriteResponse[];
  meta: PaginationMeta;
}
```

### Example Response
```json
{
  "success": true,
  "message": "Favorites retrieved",
  "data": {
    "favorites": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "type": "place",
        "externalId": "ChIJ5Wl37g6Z4jARiP4itarBPDQ",
        "title": "วัดพระศรีรัตนศาสดาราม",
        "url": "https://maps.google.com/?cid=123456",
        "thumbnailUrl": "https://example.com/photo.jpg",
        "rating": 4.7,
        "reviewCount": 40683,
        "address": "กรุงเทพมหานคร",
        "createdAt": "2024-01-15T10:30:00Z"
      },
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "type": "video",
        "externalId": "dQw4w9WgXcQ",
        "title": "10 ที่เที่ยวกรุงเทพ",
        "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
        "thumbnailUrl": "https://img.youtube.com/vi/dQw4w9WgXcQ/hqdefault.jpg",
        "createdAt": "2024-01-14T15:00:00Z"
      }
    ],
    "meta": {
      "total": 25,
      "offset": 0,
      "limit": 20
    }
  }
}
```

---

## 6.3 Remove Favorite (ลบรายการโปรด)

### Endpoint
```
DELETE /api/v1/favorites/:id
```

### Path Parameters
- `id`: UUID ของรายการโปรด

### Example Request
```
DELETE /api/v1/favorites/550e8400-e29b-41d4-a716-446655440000
```

### Response
```json
{
  "success": true,
  "message": "Favorite removed successfully"
}
```

---

## 6.4 Check Favorite (ตรวจสอบว่าเป็นรายการโปรดหรือไม่)

### Endpoint
```
GET /api/v1/favorites/check
```

### Query Parameters
```typescript
interface CheckFavoriteRequest {
  type: 'place' | 'website' | 'image' | 'video';  // required
  url?: string;        // required if no externalId
  externalId?: string; // required if no url
}
```

### Example Requests
```
GET /api/v1/favorites/check?type=place&externalId=ChIJ5Wl37g6Z4jARiP4itarBPDQ
GET /api/v1/favorites/check?type=website&url=https://www.stou.ac.th/
```

### Response
```typescript
interface CheckFavoriteResponse {
  isFavorite: boolean;
  favoriteId?: string;  // UUID ถ้าเป็นรายการโปรด
}
```

### Example Response (เป็นรายการโปรด)
```json
{
  "success": true,
  "message": "Favorite status checked",
  "data": {
    "isFavorite": true,
    "favoriteId": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### Example Response (ไม่เป็นรายการโปรด)
```json
{
  "success": true,
  "message": "Favorite status checked",
  "data": {
    "isFavorite": false
  }
}
```

---

## 6.5 Toggle Favorite (สลับสถานะรายการโปรด)

### Endpoint
```
POST /api/v1/favorites/toggle
```

### Description
ถ้ายังไม่เป็นรายการโปรด จะเพิ่ม
ถ้าเป็นรายการโปรดแล้ว จะลบออก

### Request Body
Same as `AddFavoriteRequest`

### Example Request
```json
{
  "type": "place",
  "externalId": "ChIJ5Wl37g6Z4jARiP4itarBPDQ",
  "title": "วัดพระศรีรัตนศาสดาราม",
  "url": "https://maps.google.com/?cid=123456",
  "thumbnailUrl": "https://example.com/photo.jpg",
  "rating": 4.7
}
```

### Response (เพิ่ม)
```json
{
  "success": true,
  "message": "Favorite added",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "place",
    "externalId": "ChIJ5Wl37g6Z4jARiP4itarBPDQ",
    "title": "วัดพระศรีรัตนศาสดาราม",
    ...
  }
}
```

### Response (ลบ)
```json
{
  "success": true,
  "message": "Favorite removed"
}
```

---

## TypeScript Types สำหรับ Frontend

```typescript
// types/favorites.ts

export type FavoriteType = 'place' | 'website' | 'image' | 'video';

export interface AddFavoriteRequest {
  type: FavoriteType;
  externalId?: string;
  title: string;
  url: string;
  thumbnailUrl?: string;
  rating?: number;
  reviewCount?: number;
  address?: string;
  metadata?: Record<string, any>;
}

export interface Favorite {
  id: string;
  type: FavoriteType;
  externalId?: string;
  title: string;
  url: string;
  thumbnailUrl?: string;
  rating?: number;
  reviewCount?: number;
  address?: string;
  metadata?: Record<string, any>;
  createdAt: string;
}

export interface FavoriteListResponse {
  favorites: Favorite[];
  meta: PaginationMeta;
}

export interface CheckFavoriteRequest {
  type: FavoriteType;
  url?: string;
  externalId?: string;
}

export interface CheckFavoriteResponse {
  isFavorite: boolean;
  favoriteId?: string;
}

export interface PaginationMeta {
  total: number;
  offset: number;
  limit: number;
}

// Helper: สร้าง favorite request จาก search result
export function createPlaceFavoriteRequest(place: PlaceResult): AddFavoriteRequest {
  return {
    type: 'place',
    externalId: place.placeId,
    title: place.name,
    url: `https://www.google.com/maps/place/?q=place_id:${place.placeId}`,
    thumbnailUrl: place.photoUrl,
    rating: place.rating,
    reviewCount: place.reviewCount,
    address: place.address,
    metadata: {
      lat: place.lat,
      lng: place.lng,
      types: place.types
    }
  };
}

export function createVideoFavoriteRequest(video: VideoResult): AddFavoriteRequest {
  return {
    type: 'video',
    externalId: video.videoId,
    title: video.title,
    url: `https://www.youtube.com/watch?v=${video.videoId}`,
    thumbnailUrl: video.thumbnailUrl,
    metadata: {
      channelTitle: video.channelTitle,
      viewCount: video.viewCount,
      duration: video.duration
    }
  };
}
```

---

## API Routes Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/favorites` | Yes | เพิ่มรายการโปรด |
| GET | `/api/v1/favorites` | Yes | รายการโปรดทั้งหมด |
| DELETE | `/api/v1/favorites/:id` | Yes | ลบรายการโปรด |
| GET | `/api/v1/favorites/check` | Yes | ตรวจสอบสถานะ |
| POST | `/api/v1/favorites/toggle` | Yes | สลับสถานะ |

---

## Notes

### Favorites vs Folders
- **Favorites**: รายการเดี่ยว สำหรับบันทึกอย่างรวดเร็ว
- **Folders**: กลุ่มรายการ สำหรับจัดหมวดหมู่

### externalId Usage
- **Place**: ใช้ `placeId` จาก Google Places
- **Video**: ใช้ `videoId` จาก YouTube
- **Website/Image**: ไม่จำเป็นต้องใส่

### Check Before Adding
แนะนำให้ใช้ `/favorites/check` ก่อนเพิ่ม เพื่อแสดงสถานะปุ่มให้ถูกต้อง

### Toggle Pattern
ใช้ `/favorites/toggle` สำหรับปุ่มหัวใจ/ดาว - ถ้ากดแล้วเป็นรายการโปรด กดอีกทีจะลบออก

### UI Example
```tsx
// ตัวอย่างการใช้งานใน React
const FavoriteButton = ({ item, type }: Props) => {
  const [isFavorite, setIsFavorite] = useState(false);
  const [favoriteId, setFavoriteId] = useState<string | null>(null);

  // Check status on mount
  useEffect(() => {
    checkFavoriteStatus();
  }, []);

  const checkFavoriteStatus = async () => {
    const res = await api.get('/favorites/check', {
      params: { type, externalId: item.placeId || item.videoId }
    });
    setIsFavorite(res.data.isFavorite);
    setFavoriteId(res.data.favoriteId);
  };

  const handleToggle = async () => {
    await api.post('/favorites/toggle', {
      type,
      externalId: item.placeId || item.videoId,
      title: item.name || item.title,
      url: item.url,
      thumbnailUrl: item.thumbnailUrl || item.photoUrl
    });
    setIsFavorite(!isFavorite);
  };

  return (
    <button onClick={handleToggle}>
      {isFavorite ? '❤️' : '🤍'}
    </button>
  );
};
```
