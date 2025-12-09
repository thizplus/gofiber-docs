# Part 7: Utility APIs (เครื่องมือช่วยเหลือ)

## Overview
API เครื่องมือช่วยเหลือต่างๆ เช่น แปลภาษา, สร้าง QR Code, คำนวณระยะทาง

## Base URL
```
/api/v1/utils
```

## Authentication
ไม่ต้อง (Public) - ทุก endpoint ใช้งานได้โดยไม่ต้อง login

---

## 7.1 Translate (แปลภาษา)

### Endpoint
```
POST /api/v1/utils/translate
```

### Request Body
```typescript
interface TranslateRequest {
  text: string;        // required, 1-5000 chars - ข้อความที่ต้องการแปล
  sourceLang?: string; // optional, 2 chars - ภาษาต้นทาง (auto-detect ถ้าไม่ใส่)
  targetLang: string;  // required, 2 chars - ภาษาปลายทาง
}
```

### Language Codes
| Code | Language |
|------|----------|
| th | ไทย |
| en | English |
| ja | 日本語 |
| ko | 한국어 |
| zh | 中文 |
| vi | Tiếng Việt |
| de | Deutsch |
| fr | Français |
| es | Español |

### Example Request
```json
{
  "text": "Hello, welcome to Bangkok!",
  "targetLang": "th"
}
```

### Example Request (with source language)
```json
{
  "text": "สวัสดีครับ",
  "sourceLang": "th",
  "targetLang": "en"
}
```

### Response
```typescript
interface TranslateResponse {
  originalText: string;    // ข้อความต้นฉบับ
  translatedText: string;  // ข้อความที่แปลแล้ว
  sourceLang: string;      // ภาษาต้นทาง
  targetLang: string;      // ภาษาปลายทาง
  detectedLang?: string;   // ภาษาที่ตรวจพบ (ถ้า auto-detect)
}
```

### Example Response
```json
{
  "success": true,
  "message": "Translation completed",
  "data": {
    "originalText": "Hello, welcome to Bangkok!",
    "translatedText": "สวัสดี ยินดีต้อนรับสู่กรุงเทพฯ!",
    "sourceLang": "en",
    "targetLang": "th",
    "detectedLang": "en"
  }
}
```

---

## 7.2 Detect Language (ตรวจจับภาษา)

### Endpoint
```
POST /api/v1/utils/detect-language
```

### Request Body
```typescript
interface DetectLanguageRequest {
  text: string;  // required, 1-1000 chars
}
```

### Example Request
```json
{
  "text": "こんにちは、バンコクへようこそ！"
}
```

### Response
```typescript
interface DetectLanguageResponse {
  text: string;        // ข้อความที่ตรวจ
  language: string;    // รหัสภาษา (e.g., "ja")
  confidence: number;  // ความมั่นใจ (0-1)
}
```

### Example Response
```json
{
  "success": true,
  "message": "Language detected",
  "data": {
    "text": "こんにちは、バンコクへようこそ！",
    "language": "ja",
    "confidence": 0.98
  }
}
```

---

## 7.3 Generate QR Code (สร้าง QR Code)

### Endpoint
```
POST /api/v1/utils/qrcode
```

### Request Body
```typescript
interface GenerateQRRequest {
  content: string;  // required, 1-2000 chars - ข้อมูลที่ต้องการเข้ารหัส
  size?: number;    // optional, 100-1000, default 200 - ขนาด (pixels)
  format?: string;  // optional, "png" | "svg", default "png"
}
```

### Example Request
```json
{
  "content": "https://www.google.com/maps/place/?q=place_id:ChIJ5Wl37g6Z4jARiP4itarBPDQ",
  "size": 300,
  "format": "png"
}
```

### Response
```typescript
interface GenerateQRResponse {
  content: string;   // ข้อมูลที่เข้ารหัส
  qrCodeUrl: string; // URL รูป QR Code (base64 หรือ URL)
  size: number;
  format: string;
}
```

### Example Response
```json
{
  "success": true,
  "message": "QR code generated",
  "data": {
    "content": "https://www.google.com/maps/place/?q=place_id:ChIJ5Wl37g6Z4jARiP4itarBPDQ",
    "qrCodeUrl": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
    "size": 300,
    "format": "png"
  }
}
```

### Use Cases
- QR Code สำหรับแชร์สถานที่
- QR Code สำหรับ link เว็บไซต์
- QR Code สำหรับข้อความ

---

## 7.4 Calculate Distance (คำนวณระยะทาง)

### Endpoint
```
GET /api/v1/utils/distance
```

### Query Parameters
```typescript
interface CalculateDistanceRequest {
  originLat: number;       // required - latitude ต้นทาง
  originLng: number;       // required - longitude ต้นทาง
  destinationLat: number;  // required - latitude ปลายทาง
  destinationLng: number;  // required - longitude ปลายทาง
}
```

### Example Request
```
GET /api/v1/utils/distance?originLat=13.7563&originLng=100.5018&destinationLat=13.7516&destinationLng=100.4927
```

### Response
```typescript
interface CalculateDistanceResponse {
  distanceMeters: number;  // ระยะทาง (เมตร)
  distanceKm: number;      // ระยะทาง (กิโลเมตร)
  distanceText: string;    // ระยะทาง (text format)
}
```

### Example Response
```json
{
  "success": true,
  "message": "Distance calculated",
  "data": {
    "distanceMeters": 1234.56,
    "distanceKm": 1.23,
    "distanceText": "1.2 km"
  }
}
```

### Note
- ใช้ Haversine formula ในการคำนวณ
- เป็นระยะทางเส้นตรง (as the crow flies) ไม่ใช่ระยะทางถนน

---

## TypeScript Types สำหรับ Frontend

```typescript
// types/utility.ts

// Translation
export interface TranslateRequest {
  text: string;
  sourceLang?: string;
  targetLang: string;
}

export interface TranslateResponse {
  originalText: string;
  translatedText: string;
  sourceLang: string;
  targetLang: string;
  detectedLang?: string;
}

export interface DetectLanguageRequest {
  text: string;
}

export interface DetectLanguageResponse {
  text: string;
  language: string;
  confidence: number;
}

// QR Code
export interface GenerateQRRequest {
  content: string;
  size?: number;
  format?: 'png' | 'svg';
}

export interface GenerateQRResponse {
  content: string;
  qrCodeUrl: string;
  size: number;
  format: string;
}

// Distance
export interface CalculateDistanceRequest {
  originLat: number;
  originLng: number;
  destinationLat: number;
  destinationLng: number;
}

export interface CalculateDistanceResponse {
  distanceMeters: number;
  distanceKm: number;
  distanceText: string;
}

// Language codes
export type LanguageCode =
  | 'th' | 'en' | 'ja' | 'ko' | 'zh'
  | 'vi' | 'de' | 'fr' | 'es' | 'pt'
  | 'it' | 'ru' | 'ar' | 'hi' | 'id';

export const LANGUAGE_NAMES: Record<LanguageCode, string> = {
  th: 'ไทย',
  en: 'English',
  ja: '日本語',
  ko: '한국어',
  zh: '中文',
  vi: 'Tiếng Việt',
  de: 'Deutsch',
  fr: 'Français',
  es: 'Español',
  pt: 'Português',
  it: 'Italiano',
  ru: 'Русский',
  ar: 'العربية',
  hi: 'हिन्दी',
  id: 'Bahasa Indonesia'
};

// Helper functions
export function formatDistance(meters: number): string {
  if (meters < 1000) {
    return `${Math.round(meters)} m`;
  }
  return `${(meters / 1000).toFixed(1)} km`;
}

export function calculateDistanceFromCoords(
  lat1: number, lng1: number,
  lat2: number, lng2: number
): number {
  const R = 6371000; // Earth's radius in meters
  const dLat = (lat2 - lat1) * Math.PI / 180;
  const dLng = (lng2 - lng1) * Math.PI / 180;
  const a =
    Math.sin(dLat/2) * Math.sin(dLat/2) +
    Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
    Math.sin(dLng/2) * Math.sin(dLng/2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));
  return R * c;
}
```

---

## API Routes Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/utils/translate` | No | แปลภาษา |
| POST | `/api/v1/utils/detect-language` | No | ตรวจจับภาษา |
| POST | `/api/v1/utils/qrcode` | No | สร้าง QR Code |
| GET | `/api/v1/utils/distance` | No | คำนวณระยะทาง |

---

## Notes

### Translation Caching
- การแปลภาษา cache 7 วัน ใน Redis
- ตรวจจับภาษา cache ตามข้อความ

### QR Code
- รองรับ format PNG และ SVG
- ขนาดสูงสุด 1000x1000 pixels
- Return เป็น base64 data URL

### Distance Calculation
- ใช้ Haversine formula
- คำนวณระยะทางเส้นตรง
- สำหรับระยะทางถนนจริง ต้องใช้ Google Directions API

### Rate Limiting
- Translation: 100 requests/hour
- QR Code: 200 requests/hour
- Distance: ไม่จำกัด (local calculation)

### Use Case Examples

**1. Translation in Search Results**
```typescript
// แปลคำอธิบายสถานที่เป็นภาษาผู้ใช้
const translateDescription = async (text: string, userLang: string) => {
  const res = await api.post('/utils/translate', {
    text,
    targetLang: userLang
  });
  return res.data.translatedText;
};
```

**2. Share Place with QR Code**
```typescript
// สร้าง QR Code สำหรับแชร์สถานที่
const generatePlaceQR = async (placeId: string) => {
  const url = `https://yourapp.com/place/${placeId}`;
  const res = await api.post('/utils/qrcode', {
    content: url,
    size: 300
  });
  return res.data.qrCodeUrl;
};
```

**3. Sort Places by Distance**
```typescript
// เรียงสถานที่ตามระยะทางจากตำแหน่งปัจจุบัน
const sortByDistance = (places: Place[], userLat: number, userLng: number) => {
  return places
    .map(place => ({
      ...place,
      distance: calculateDistanceFromCoords(
        userLat, userLng,
        place.lat, place.lng
      )
    }))
    .sort((a, b) => a.distance - b.distance);
};
```
