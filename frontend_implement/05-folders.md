# Part 5: Folders APIs (จัดการโฟลเดอร์)

## Overview
ระบบจัดการโฟลเดอร์สำหรับเก็บรวบรวมข้อมูลการค้นหา (สถานที่, เว็บ, รูป, วิดีโอ)

## Base URL
```
/api/v1/folders
```

---

## 5.1 Create Folder (สร้างโฟลเดอร์)

### Endpoint
```
POST /api/v1/folders
```

### Authentication
Required (Bearer Token)

### Request Body
```typescript
interface CreateFolderRequest {
  name: string;           // required, 1-255 chars
  description?: string;   // optional, max 1000 chars
  coverImageUrl?: string; // optional, URL format, max 500 chars
  isPublic?: boolean;     // optional, default false
}
```

### Example Request
```json
{
  "name": "ที่เที่ยวกรุงเทพ",
  "description": "รวมสถานที่ท่องเที่ยวในกรุงเทพมหานคร",
  "coverImageUrl": "https://example.com/cover.jpg",
  "isPublic": false
}
```

### Response
```typescript
interface FolderResponse {
  id: string;             // UUID
  name: string;
  description: string;
  coverImageUrl?: string;
  isPublic: boolean;
  itemCount: number;
  createdAt: string;
  updatedAt: string;
}
```

### Example Response
```json
{
  "success": true,
  "message": "Folder created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "ที่เที่ยวกรุงเทพ",
    "description": "รวมสถานที่ท่องเที่ยวในกรุงเทพมหานคร",
    "coverImageUrl": "https://example.com/cover.jpg",
    "isPublic": false,
    "itemCount": 0,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  }
}
```

---

## 5.2 Get Folders (รายการโฟลเดอร์)

### Endpoint
```
GET /api/v1/folders
```

### Authentication
Required (Bearer Token)

### Query Parameters
```typescript
interface GetFoldersRequest {
  page?: number;     // optional, default 1
  pageSize?: number; // optional, default 20, max 50
  isPublic?: boolean; // optional - filter by public/private
}
```

### Example Request
```
GET /api/v1/folders?page=1&pageSize=10
```

### Response
```typescript
interface FolderListResponse {
  folders: FolderResponse[];
  meta: PaginationMeta;
}
```

### Example Response
```json
{
  "success": true,
  "message": "Folders retrieved",
  "data": {
    "folders": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "ที่เที่ยวกรุงเทพ",
        "description": "รวมสถานที่ท่องเที่ยวในกรุงเทพมหานคร",
        "coverImageUrl": "https://example.com/cover.jpg",
        "isPublic": false,
        "itemCount": 5,
        "createdAt": "2024-01-15T10:30:00Z",
        "updatedAt": "2024-01-15T10:30:00Z"
      }
    ],
    "meta": {
      "total": 10,
      "offset": 0,
      "limit": 20
    }
  }
}
```

---

## 5.3 Get Folder Detail (รายละเอียดโฟลเดอร์)

### Endpoint
```
GET /api/v1/folders/:id
```

### Authentication
Required (Bearer Token)

### Path Parameters
- `id`: UUID ของโฟลเดอร์

### Response
```typescript
interface FolderDetailResponse {
  id: string;
  name: string;
  description: string;
  coverImageUrl?: string;
  isPublic: boolean;
  itemCount: number;
  items: FolderItemResponse[];  // รายการไอเทมในโฟลเดอร์
  createdAt: string;
  updatedAt: string;
}
```

### Example Response
```json
{
  "success": true,
  "message": "Folder retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "ที่เที่ยวกรุงเทพ",
    "description": "รวมสถานที่ท่องเที่ยว",
    "isPublic": false,
    "itemCount": 2,
    "items": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "folderId": "550e8400-e29b-41d4-a716-446655440000",
        "type": "place",
        "title": "วัดพระแก้ว",
        "url": "https://maps.google.com/?cid=...",
        "thumbnailUrl": "https://...",
        "description": "วัดสำคัญในกรุงเทพ",
        "metadata": {
          "placeId": "ChIJ...",
          "rating": 4.7,
          "lat": 13.75,
          "lng": 100.49
        },
        "sortOrder": 0,
        "createdAt": "2024-01-15T10:35:00Z"
      }
    ],
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:35:00Z"
  }
}
```

---

## 5.4 Update Folder (แก้ไขโฟลเดอร์)

### Endpoint
```
PUT /api/v1/folders/:id
```

### Authentication
Required (Bearer Token)

### Request Body
```typescript
interface UpdateFolderRequest {
  name?: string;          // optional, 1-255 chars
  description?: string;   // optional, max 1000 chars
  coverImageUrl?: string; // optional, URL format
  isPublic?: boolean;     // optional
}
```

### Example Request
```json
{
  "name": "ที่เที่ยวกรุงเทพ 2024",
  "isPublic": true
}
```

### Response
Same as `FolderResponse`

---

## 5.5 Delete Folder (ลบโฟลเดอร์)

### Endpoint
```
DELETE /api/v1/folders/:id
```

### Authentication
Required (Bearer Token)

### Response
```json
{
  "success": true,
  "message": "Folder deleted successfully"
}
```

---

## 5.6 Share Folder (แชร์โฟลเดอร์)

### Endpoint
```
POST /api/v1/folders/:id/share
```

### Authentication
Required (Bearer Token)

### Request Body
```typescript
interface ShareFolderRequest {
  isPublic: boolean;  // true = public, false = private
}
```

### Response
```typescript
interface FolderShareResponse {
  folderId: string;
  isPublic: boolean;
  shareUrl?: string;  // URL สำหรับแชร์ (ถ้า public)
}
```

### Example Response
```json
{
  "success": true,
  "message": "Folder sharing updated",
  "data": {
    "folderId": "550e8400-e29b-41d4-a716-446655440000",
    "isPublic": true,
    "shareUrl": "https://yourapp.com/folders/public/550e8400-e29b-41d4-a716-446655440000"
  }
}
```

---

## 5.7 Get Public Folder (ดูโฟลเดอร์สาธารณะ)

### Endpoint
```
GET /api/v1/folders/public/:id
```

### Authentication
ไม่ต้อง (Public)

### Path Parameters
- `id`: UUID ของโฟลเดอร์

### Response
Same as `FolderDetailResponse` (แต่เฉพาะโฟลเดอร์ที่เป็น public เท่านั้น)

---

## 5.8 Folder Items (จัดการไอเทมในโฟลเดอร์)

### Add Item to Folder

```
POST /api/v1/folders/:id/items
```

### Authentication
Required (Bearer Token)

### Request Body
```typescript
interface AddFolderItemRequest {
  type: 'place' | 'website' | 'image' | 'video' | 'link';
  title: string;                    // required, 1-255 chars
  url: string;                      // required, URL format, max 2000 chars
  thumbnailUrl?: string;            // optional, URL format
  description?: string;             // optional, max 1000 chars
  metadata?: Record<string, any>;   // optional - ข้อมูลเพิ่มเติม
}
```

### Example Request (Place)
```json
{
  "type": "place",
  "title": "วัดพระศรีรัตนศาสดาราม",
  "url": "https://maps.google.com/?cid=123456",
  "thumbnailUrl": "https://example.com/photo.jpg",
  "description": "วัดพระแก้ว",
  "metadata": {
    "placeId": "ChIJ5Wl37g6Z4jARiP4itarBPDQ",
    "lat": 13.7516435,
    "lng": 100.4927041,
    "rating": 4.7,
    "reviewCount": 40683
  }
}
```

### Example Request (Video)
```json
{
  "type": "video",
  "title": "10 ที่เที่ยวกรุงเทพ",
  "url": "https://www.youtube.com/watch?v=abc123",
  "thumbnailUrl": "https://img.youtube.com/vi/abc123/hqdefault.jpg",
  "description": "วิดีโอแนะนำที่เที่ยว",
  "metadata": {
    "videoId": "abc123",
    "channelTitle": "Travel Channel",
    "viewCount": 100000
  }
}
```

### Response
```typescript
interface FolderItemResponse {
  id: string;
  folderId: string;
  type: string;
  title: string;
  url: string;
  thumbnailUrl?: string;
  description?: string;
  metadata?: Record<string, any>;
  sortOrder: number;
  createdAt: string;
}
```

---

### Get Folder Items

```
GET /api/v1/folders/:id/items
```

### Query Parameters
```typescript
interface GetFolderItemsRequest {
  type?: 'place' | 'website' | 'image' | 'video' | 'link';
  page?: number;
  pageSize?: number;
}
```

### Response
```typescript
interface FolderItemListResponse {
  items: FolderItemResponse[];
  meta: PaginationMeta;
}
```

---

### Update Folder Item

```
PUT /api/v1/folders/items/:itemId
```

### Request Body
```typescript
interface UpdateFolderItemRequest {
  title?: string;       // optional, 1-255 chars
  description?: string; // optional, max 1000 chars
  sortOrder?: number;   // optional, min 0
}
```

---

### Delete Folder Item

```
DELETE /api/v1/folders/items/:itemId
```

---

### Reorder Folder Items

```
PUT /api/v1/folders/:id/items/reorder
```

### Request Body
```typescript
interface ReorderFolderItemsRequest {
  itemOrders: ItemOrder[];
}

interface ItemOrder {
  itemId: string;  // UUID
  sortOrder: number;
}
```

### Example Request
```json
{
  "itemOrders": [
    { "itemId": "item-1-uuid", "sortOrder": 0 },
    { "itemId": "item-2-uuid", "sortOrder": 1 },
    { "itemId": "item-3-uuid", "sortOrder": 2 }
  ]
}
```

---

## TypeScript Types สำหรับ Frontend

```typescript
// types/folders.ts

export type FolderItemType = 'place' | 'website' | 'image' | 'video' | 'link';

export interface CreateFolderRequest {
  name: string;
  description?: string;
  coverImageUrl?: string;
  isPublic?: boolean;
}

export interface UpdateFolderRequest {
  name?: string;
  description?: string;
  coverImageUrl?: string;
  isPublic?: boolean;
}

export interface Folder {
  id: string;
  name: string;
  description: string;
  coverImageUrl?: string;
  isPublic: boolean;
  itemCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface FolderDetail extends Folder {
  items: FolderItem[];
}

export interface FolderItem {
  id: string;
  folderId: string;
  type: FolderItemType;
  title: string;
  url: string;
  thumbnailUrl?: string;
  description?: string;
  metadata?: Record<string, any>;
  sortOrder: number;
  createdAt: string;
}

export interface AddFolderItemRequest {
  type: FolderItemType;
  title: string;
  url: string;
  thumbnailUrl?: string;
  description?: string;
  metadata?: Record<string, any>;
}

export interface UpdateFolderItemRequest {
  title?: string;
  description?: string;
  sortOrder?: number;
}

export interface ReorderItemsRequest {
  itemOrders: {
    itemId: string;
    sortOrder: number;
  }[];
}

export interface ShareFolderRequest {
  isPublic: boolean;
}

export interface ShareFolderResponse {
  folderId: string;
  isPublic: boolean;
  shareUrl?: string;
}

export interface FolderListResponse {
  folders: Folder[];
  meta: PaginationMeta;
}

export interface PaginationMeta {
  total: number;
  offset: number;
  limit: number;
}

// Place metadata
export interface PlaceMetadata {
  placeId: string;
  lat: number;
  lng: number;
  rating?: number;
  reviewCount?: number;
  address?: string;
  types?: string[];
}

// Video metadata
export interface VideoMetadata {
  videoId: string;
  channelTitle?: string;
  duration?: string;
  viewCount?: number;
}
```

---

## API Routes Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/folders` | Yes | สร้างโฟลเดอร์ |
| GET | `/api/v1/folders` | Yes | รายการโฟลเดอร์ |
| GET | `/api/v1/folders/:id` | Yes | รายละเอียดโฟลเดอร์ |
| PUT | `/api/v1/folders/:id` | Yes | แก้ไขโฟลเดอร์ |
| DELETE | `/api/v1/folders/:id` | Yes | ลบโฟลเดอร์ |
| POST | `/api/v1/folders/:id/share` | Yes | แชร์โฟลเดอร์ |
| GET | `/api/v1/folders/public/:id` | No | ดูโฟลเดอร์สาธารณะ |
| POST | `/api/v1/folders/:id/items` | Yes | เพิ่มไอเทม |
| GET | `/api/v1/folders/:id/items` | Yes | รายการไอเทม |
| PUT | `/api/v1/folders/:id/items/reorder` | Yes | เรียงลำดับไอเทม |
| PUT | `/api/v1/folders/items/:itemId` | Yes | แก้ไขไอเทม |
| DELETE | `/api/v1/folders/items/:itemId` | Yes | ลบไอเทม |

---

## Notes

### Use Cases
1. **บันทึกสถานที่**: เมื่อค้นหาสถานที่แล้วต้องการบันทึก
2. **บันทึกวิดีโอ**: เก็บวิดีโอ YouTube ที่สนใจ
3. **สร้างแผนเที่ยว**: รวบรวมสถานที่ ร้านอาหาร โรงแรม ไว้ในโฟลเดอร์เดียว
4. **แชร์กับเพื่อน**: ทำให้โฟลเดอร์เป็น public แล้วแชร์ link

### Metadata Guidelines
- **Place**: เก็บ `placeId`, `lat`, `lng`, `rating`, `reviewCount`
- **Video**: เก็บ `videoId`, `channelTitle`, `viewCount`
- **Website**: เก็บ `displayLink`, `snippet`
- **Image**: เก็บ `width`, `height`, `source`

### sortOrder
- ใช้สำหรับเรียงลำดับไอเทมในโฟลเดอร์
- ค่าเริ่มต้นจะเป็น 0, 1, 2, ... ตามลำดับที่เพิ่ม
- สามารถ reorder ได้ตามต้องการ
