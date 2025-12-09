# Part 3: Places APIs (Google Maps / Places)

## Overview
ระบบค้นหาสถานที่จาก Google Places API รองรับทั้ง Text Search และ Nearby Search

## Base URL
```
/api/v1/search
```

---

## 3.1 Search Places (ค้นหาสถานที่)

### Endpoint
```
GET /api/v1/search/places
```

### Authentication
Optional (ถ้า login จะบันทึกประวัติ)

### Query Parameters
```typescript
interface PlaceSearchRequest {
  q: string;          // required - คำค้นหา (e.g., "ที่เที่ยว กรุงเทพ")
  lat?: number;       // optional - latitude (ถ้าใส่จะใช้ Nearby Search)
  lng?: number;       // optional - longitude (ถ้าใส่จะใช้ Nearby Search)
  radius?: number;    // optional - รัศมี (meters), 100-50000, default 5000
  type?: string;      // optional - ประเภทสถานที่ (restaurant, tourist_attraction, etc.)
  page?: number;      // optional - หน้า
  pageSize?: number;  // optional - จำนวนต่อหน้า, max 20
}
```

### Search Modes

**1. Text Search (ค้นหาด้วยข้อความ)**
- ใช้เมื่อ: ไม่ส่ง `lat` และ `lng`
- เหมือนกับการค้นหาใน Google Maps
- ค้นหาได้ทุกที่ในโลก

```
GET /api/v1/search/places?q=ที่เที่ยว+กรุงเทพ
```

**2. Nearby Search (ค้นหาใกล้เคียง)**
- ใช้เมื่อ: ส่ง `lat` และ `lng` มาด้วย
- ค้นหาสถานที่รอบๆ ตำแหน่งที่กำหนด

```
GET /api/v1/search/places?q=ร้านอาหาร&lat=13.7563&lng=100.5018&radius=2000
```

### Response
```typescript
interface PlaceSearchResponse {
  query: string;
  results: PlaceResult[];
  totalCount: number;
  page: number;
  pageSize: number;
}

interface PlaceResult {
  placeId: string;        // Google Place ID
  name: string;           // ชื่อสถานที่
  address: string;        // ที่อยู่
  lat: number;            // latitude
  lng: number;            // longitude
  rating: number;         // คะแนน (0-5)
  reviewCount: number;    // จำนวนรีวิว
  priceLevel?: number;    // ระดับราคา (1-4)
  types: string[];        // ประเภทสถานที่
  photoUrl?: string;      // URL รูปภาพ
  isOpen?: boolean;       // เปิดอยู่หรือไม่
  distance?: number;      // ระยะทาง (meters) - เฉพาะ Nearby Search
  distanceText?: string;  // ระยะทาง (text) - e.g., "1.2 km"
}
```

### Example Request (Text Search)
```
GET /api/v1/search/places?q=ที่เที่ยว+กรุงเทพ
```

### Example Response
```json
{
  "success": true,
  "message": "Place search completed",
  "data": {
    "query": "ที่เที่ยว กรุงเทพ",
    "results": [
      {
        "placeId": "ChIJX49VWsWZ4jARwEYBM0MDXnE",
        "name": "เสาชิงช้า",
        "address": "ถ. ดินสอ แขวง บางขุนพรหม เขตพระนคร กรุงเทพมหานคร",
        "lat": 13.7517707,
        "lng": 100.5012761,
        "rating": 4.6,
        "reviewCount": 3526,
        "types": ["tourist_attraction", "point_of_interest", "establishment"],
        "photoUrl": "https://maps.googleapis.com/maps/api/place/photo?...",
        "isOpen": true
      },
      {
        "placeId": "ChIJH4I85c2e4jARUMYLeiqfBd0",
        "name": "ซีไลฟ์ แบงค็อก โอเชียน เวิลด์",
        "address": "ชั้น บี1-บี2 สยามพารากอน กรุงเทพมหานคร",
        "lat": 13.7459351,
        "lng": 100.5352057,
        "rating": 4.5,
        "reviewCount": 27219,
        "types": ["aquarium", "tourist_attraction", "point_of_interest"],
        "photoUrl": "https://maps.googleapis.com/maps/api/place/photo?...",
        "isOpen": true
      }
    ],
    "totalCount": 20,
    "page": 0,
    "pageSize": 0
  }
}
```

---

## 3.2 Get Place Details (รายละเอียดสถานที่)

### Endpoint
```
GET /api/v1/search/places/:placeId
```

### Authentication
ไม่ต้อง (Public)

### Path Parameters
- `placeId`: Google Place ID

### Example Request
```
GET /api/v1/search/places/ChIJX49VWsWZ4jARwEYBM0MDXnE
```

### Response
```typescript
interface PlaceDetailResponse {
  placeId: string;
  name: string;
  formattedAddress: string;
  lat: number;
  lng: number;
  rating: number;
  reviewCount: number;
  priceLevel?: number;
  types: string[];
  phone?: string;              // เบอร์โทร
  website?: string;            // เว็บไซต์
  googleMapsUrl: string;       // URL Google Maps
  openingHours?: string[];     // เวลาเปิด-ปิด
  reviews?: PlaceReview[];     // รีวิว
  photos?: PlacePhoto[];       // รูปภาพ
  distance?: number;
  distanceText?: string;
}

interface PlaceReview {
  author: string;          // ชื่อผู้รีวิว
  rating: number;          // คะแนน (1-5)
  text: string;            // ข้อความรีวิว
  time: string;            // เวลาที่รีวิว
  photoUrl?: string;       // รูปโปรไฟล์ผู้รีวิว
}

interface PlacePhoto {
  url: string;             // URL รูปภาพ
  width: number;
  height: number;
}
```

### Example Response
```json
{
  "success": true,
  "message": "Place details retrieved",
  "data": {
    "placeId": "ChIJX49VWsWZ4jARwEYBM0MDXnE",
    "name": "เสาชิงช้า",
    "formattedAddress": "QG22+PGF ถ. ดินสอ แขวง บางขุนพรหม เขตพระนคร กรุงเทพมหานคร 10200",
    "lat": 13.7517707,
    "lng": 100.5012761,
    "rating": 4.6,
    "reviewCount": 3526,
    "types": ["tourist_attraction", "point_of_interest", "establishment"],
    "phone": "+66 2 225 9999",
    "website": "https://www.tourismthailand.org/",
    "googleMapsUrl": "https://maps.google.com/?cid=...",
    "openingHours": [
      "Monday: Open 24 hours",
      "Tuesday: Open 24 hours",
      "Wednesday: Open 24 hours",
      "Thursday: Open 24 hours",
      "Friday: Open 24 hours",
      "Saturday: Open 24 hours",
      "Sunday: Open 24 hours"
    ],
    "reviews": [
      {
        "author": "John Doe",
        "rating": 5,
        "text": "สถานที่ท่องเที่ยวสำคัญของกรุงเทพ สวยมาก",
        "time": "2024-01-10T12:00:00Z",
        "photoUrl": "https://..."
      }
    ],
    "photos": [
      {
        "url": "https://maps.googleapis.com/maps/api/place/photo?...",
        "width": 4032,
        "height": 3024
      }
    ]
  }
}
```

---

## 3.3 Search Nearby Places (ค้นหาสถานที่ใกล้เคียง)

### Endpoint
```
GET /api/v1/search/nearby
```

### Authentication
ไม่ต้อง (Public)

### Query Parameters
```typescript
interface NearbyPlacesRequest {
  lat: number;        // required - latitude
  lng: number;        // required - longitude
  radius?: number;    // optional - รัศมี (meters), 100-50000, default 5000
  type?: string;      // optional - ประเภทสถานที่
  keyword?: string;   // optional - คำค้นหา
  page?: number;
  pageSize?: number;  // max 20
}
```

### Example Request
```
GET /api/v1/search/nearby?lat=13.7563&lng=100.5018&radius=1000&type=restaurant
```

### Response
Same as PlaceSearchResponse

---

## Place Types (ประเภทสถานที่)

```typescript
type PlaceType =
  // ท่องเที่ยว
  | 'tourist_attraction'
  | 'museum'
  | 'park'
  | 'zoo'
  | 'aquarium'
  | 'amusement_park'

  // อาหาร
  | 'restaurant'
  | 'cafe'
  | 'bar'
  | 'bakery'

  // ที่พัก
  | 'lodging'
  | 'hotel'

  // ช็อปปิ้ง
  | 'shopping_mall'
  | 'store'

  // ขนส่ง
  | 'airport'
  | 'train_station'
  | 'bus_station'
  | 'transit_station'

  // สถานที่สำคัญ
  | 'church'
  | 'hindu_temple'
  | 'mosque'
  | 'place_of_worship'
  | 'city_hall'

  // อื่นๆ
  | 'hospital'
  | 'pharmacy'
  | 'bank'
  | 'atm'
  | 'gas_station';
```

---

## TypeScript Types สำหรับ Frontend

```typescript
// types/places.ts

export interface PlaceSearchRequest {
  q: string;
  lat?: number;
  lng?: number;
  radius?: number;
  type?: string;
  page?: number;
  pageSize?: number;
}

export interface PlaceResult {
  placeId: string;
  name: string;
  address: string;
  lat: number;
  lng: number;
  rating: number;
  reviewCount: number;
  priceLevel?: number;
  types: string[];
  photoUrl?: string;
  isOpen?: boolean;
  distance?: number;
  distanceText?: string;
}

export interface PlaceSearchResponse {
  query: string;
  results: PlaceResult[];
  totalCount: number;
  page: number;
  pageSize: number;
}

export interface PlaceDetail {
  placeId: string;
  name: string;
  formattedAddress: string;
  lat: number;
  lng: number;
  rating: number;
  reviewCount: number;
  priceLevel?: number;
  types: string[];
  phone?: string;
  website?: string;
  googleMapsUrl: string;
  openingHours?: string[];
  reviews?: PlaceReview[];
  photos?: PlacePhoto[];
  distance?: number;
  distanceText?: string;
}

export interface PlaceReview {
  author: string;
  rating: number;
  text: string;
  time: string;
  photoUrl?: string;
}

export interface PlacePhoto {
  url: string;
  width: number;
  height: number;
}

export interface NearbySearchRequest {
  lat: number;
  lng: number;
  radius?: number;
  type?: string;
  keyword?: string;
  page?: number;
  pageSize?: number;
}
```

---

## API Routes Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/search/places` | Optional | ค้นหาสถานที่ (Text/Nearby) |
| GET | `/api/v1/search/places/:placeId` | No | รายละเอียดสถานที่ |
| GET | `/api/v1/search/nearby` | No | ค้นหาสถานที่ใกล้เคียง |

---

## Notes

### Search Modes
- **ไม่ส่ง lat/lng**: ใช้ Text Search - ค้นหาแบบ Google Maps
- **ส่ง lat/lng**: ใช้ Nearby Search - ค้นหารอบตำแหน่งที่กำหนด

### Caching
- Place Search: cache 1 ชั่วโมง
- Place Details: cache 24 ชั่วโมง

### Photo URLs
- Photo URL จาก response สามารถใช้ได้โดยตรง
- Photo URLs มี API key แนบมาแล้ว

### Google Maps URL
- สร้าง URL สำหรับเปิดใน Google Maps:
  ```typescript
  // เปิดสถานที่ใน Google Maps
  const googleMapsUrl = `https://www.google.com/maps/place/?q=place_id:${placeId}`;

  // เปิด directions
  const directionsUrl = `https://www.google.com/maps/dir/?api=1&destination=${lat},${lng}&destination_place_id=${placeId}`;
  ```

### Distance Calculation
- ระยะทางจะคำนวณเมื่อใช้ Nearby Search
- หน่วยเป็น meters
- `distanceText` format: "500 m", "1.2 km"
