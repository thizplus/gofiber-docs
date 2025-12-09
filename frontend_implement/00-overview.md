# STOU Smart Tour - API Documentation Overview

## Project Information
- **Project Name**: STOU Smart Tour Backend
- **Framework**: GoFiber
- **Base URL**: `http://localhost:8080/api/v1`
- **Version**: 1.0.0

---

## API Structure

```
/api/v1
├── /auth              # Authentication (Register, Login)
├── /users             # User Management
├── /search            # Search APIs
│   ├── /websites      # Website Search
│   ├── /images        # Image Search
│   ├── /videos        # Video Search (YouTube)
│   ├── /places        # Places Search (Google Maps)
│   ├── /nearby        # Nearby Places
│   └── /history       # Search History
├── /ai                # AI APIs
│   ├── /search        # AI Search (Quick)
│   └── /chat          # AI Chat Sessions
├── /folders           # Folder Management
│   └── /items         # Folder Items
├── /favorites         # Favorites Management
└── /utils             # Utility Tools
    ├── /translate     # Translation
    ├── /detect-language
    ├── /qrcode        # QR Code Generator
    └── /distance      # Distance Calculator
```

---

## Authentication

### Bearer Token
ทุก endpoint ที่ต้อง login ให้ส่ง header:
```
Authorization: Bearer <token>
```

### Token Expiration
- Token หมดอายุใน 24 ชั่วโมง
- ถ้า token หมดอายุจะได้ response:
```json
{
  "success": false,
  "message": "Token expired",
  "error": "Unauthorized"
}
```

---

## Response Format

### Success Response
```typescript
interface APIResponse<T> {
  success: true;
  message: string;
  data: T;
}
```

### Error Response
```typescript
interface APIErrorResponse {
  success: false;
  message: string;
  error: string;
}
```

### Paginated Response
```typescript
interface PaginatedResponse<T> {
  success: true;
  message: string;
  data: T;
  meta: {
    total: number;
    offset: number;
    limit: number;
  };
}
```

---

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - สำเร็จ |
| 201 | Created - สร้างสำเร็จ |
| 400 | Bad Request - ข้อมูลไม่ถูกต้อง |
| 401 | Unauthorized - ไม่ได้ login หรือ token หมดอายุ |
| 403 | Forbidden - ไม่มีสิทธิ์เข้าถึง |
| 404 | Not Found - ไม่พบข้อมูล |
| 429 | Too Many Requests - เกิน rate limit |
| 500 | Internal Server Error - ข้อผิดพลาดภายในระบบ |

---

## Documentation Files

| File | Description |
|------|-------------|
| [01-authentication.md](./01-authentication.md) | Register, Login, Profile |
| [02-search.md](./02-search.md) | Website, Image, Video Search |
| [03-places.md](./03-places.md) | Google Places Search |
| [04-ai.md](./04-ai.md) | AI Search & Chat |
| [05-folders.md](./05-folders.md) | Folder Management |
| [06-favorites.md](./06-favorites.md) | Favorites Management |
| [07-utility.md](./07-utility.md) | Translation, QR Code, Distance |

---

## Complete API Endpoints

### Authentication (Public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | ลงทะเบียน |
| POST | `/auth/login` | เข้าสู่ระบบ |

### Users (Protected)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users/profile` | ดูโปรไฟล์ |
| PUT | `/users/profile` | แก้ไขโปรไฟล์ |
| DELETE | `/users/profile` | ลบบัญชี |

### Search (Public/Optional Auth)
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/search` | Optional | ค้นหาทั่วไป |
| GET | `/search/websites` | Optional | ค้นหาเว็บไซต์ |
| GET | `/search/images` | Optional | ค้นหารูปภาพ |
| GET | `/search/videos` | Optional | ค้นหาวิดีโอ |
| GET | `/search/videos/:videoId` | No | รายละเอียดวิดีโอ |
| GET | `/search/places` | Optional | ค้นหาสถานที่ |
| GET | `/search/places/:placeId` | No | รายละเอียดสถานที่ |
| GET | `/search/nearby` | No | สถานที่ใกล้เคียง |
| GET | `/search/history` | Yes | ประวัติการค้นหา |
| DELETE | `/search/history` | Yes | ลบประวัติทั้งหมด |
| DELETE | `/search/history/:id` | Yes | ลบประวัติรายการเดียว |

### AI (Public/Protected)
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/ai/search` | Optional | AI Search |
| POST | `/ai/chat` | Yes | สร้าง Chat Session |
| GET | `/ai/chat` | Yes | รายการ Sessions |
| GET | `/ai/chat/:sessionId` | Yes | รายละเอียด Session |
| POST | `/ai/chat/:sessionId/messages` | Yes | ส่งข้อความ |
| DELETE | `/ai/chat/:sessionId` | Yes | ลบ Session |
| DELETE | `/ai/chat` | Yes | ลบทั้งหมด |

### Folders (Protected)
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/folders` | Yes | สร้างโฟลเดอร์ |
| GET | `/folders` | Yes | รายการโฟลเดอร์ |
| GET | `/folders/:id` | Yes | รายละเอียดโฟลเดอร์ |
| PUT | `/folders/:id` | Yes | แก้ไขโฟลเดอร์ |
| DELETE | `/folders/:id` | Yes | ลบโฟลเดอร์ |
| POST | `/folders/:id/share` | Yes | แชร์โฟลเดอร์ |
| GET | `/folders/public/:id` | No | ดูโฟลเดอร์สาธารณะ |
| POST | `/folders/:id/items` | Yes | เพิ่มไอเทม |
| GET | `/folders/:id/items` | Yes | รายการไอเทม |
| PUT | `/folders/:id/items/reorder` | Yes | เรียงลำดับ |
| PUT | `/folders/items/:itemId` | Yes | แก้ไขไอเทม |
| DELETE | `/folders/items/:itemId` | Yes | ลบไอเทม |

### Favorites (Protected)
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/favorites` | Yes | เพิ่มรายการโปรด |
| GET | `/favorites` | Yes | รายการโปรดทั้งหมด |
| DELETE | `/favorites/:id` | Yes | ลบรายการโปรด |
| GET | `/favorites/check` | Yes | ตรวจสอบสถานะ |
| POST | `/favorites/toggle` | Yes | สลับสถานะ |

### Utility (Public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/utils/translate` | แปลภาษา |
| POST | `/utils/detect-language` | ตรวจจับภาษา |
| POST | `/utils/qrcode` | สร้าง QR Code |
| GET | `/utils/distance` | คำนวณระยะทาง |

---

## TypeScript Types File

สร้างไฟล์ `types/api.ts` รวม types ทั้งหมด:

```typescript
// types/api.ts

// ==================== Common ====================
export interface APIResponse<T> {
  success: boolean;
  message: string;
  data?: T;
  error?: string;
}

export interface PaginationMeta {
  total: number;
  offset: number;
  limit: number;
}

// ==================== Auth ====================
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface User {
  id: string;
  email: string;
  username: string;
  firstName: string;
  lastName: string;
  avatar: string;
  role: 'user' | 'admin';
  isActive: boolean;
  studentId?: string;
  createdAt: string;
  updatedAt: string;
}

// ==================== Search ====================
export interface WebsiteResult {
  title: string;
  url: string;
  snippet: string;
  displayLink: string;
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

export interface PlaceResult {
  placeId: string;
  name: string;
  address: string;
  lat: number;
  lng: number;
  rating: number;
  reviewCount: number;
  types: string[];
  photoUrl?: string;
  isOpen?: boolean;
  distance?: number;
  distanceText?: string;
}

// ==================== AI ====================
export interface AISearchResponse {
  query: string;
  summary: string;
  sources: MessageSource[];
  keywords?: string[];
}

export interface MessageSource {
  title: string;
  url: string;
  snippet?: string;
}

export interface ChatSession {
  id: string;
  title: string;
  initialQuery: string;
  messages?: ChatMessage[];
  createdAt: string;
  updatedAt: string;
}

export interface ChatMessage {
  id: string;
  sessionId: string;
  role: 'user' | 'assistant';
  content: string;
  sources?: MessageSource[];
  createdAt: string;
}

// ==================== Folders ====================
export type FolderItemType = 'place' | 'website' | 'image' | 'video' | 'link';

export interface Folder {
  id: string;
  name: string;
  description: string;
  coverImageUrl?: string;
  isPublic: boolean;
  itemCount: number;
  items?: FolderItem[];
  createdAt: string;
  updatedAt: string;
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

// ==================== Favorites ====================
export type FavoriteType = 'place' | 'website' | 'image' | 'video';

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

// ==================== Utility ====================
export interface TranslateResponse {
  originalText: string;
  translatedText: string;
  sourceLang: string;
  targetLang: string;
  detectedLang?: string;
}

export interface QRCodeResponse {
  content: string;
  qrCodeUrl: string;
  size: number;
  format: string;
}

export interface DistanceResponse {
  distanceMeters: number;
  distanceKm: number;
  distanceText: string;
}
```

---

## Caching Information

| API | Cache Duration |
|-----|----------------|
| Website Search | 1 hour |
| Image Search | 1 hour |
| YouTube Search | 6 hours |
| Places Search | 1 hour |
| Place Details | 24 hours |
| AI Search | 6 hours |
| Translation | 7 days |
| AI Chat | No cache |

---

## Environment Variables (Backend)

```env
# App
APP_NAME=STOU Smart Tour
APP_PORT=8080
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=mtu_docs

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your_jwt_secret

# Google APIs
GOOGLE_API_KEY=your_google_api_key
GOOGLE_SEARCH_ENGINE_ID=your_search_engine_id

# OpenAI
OPENAI_API_KEY=your_openai_api_key
OPENAI_MODEL=gpt-4o-mini
```

---

## Getting Started (Frontend)

### 1. Create API Client
```typescript
// lib/api.ts
import axios from 'axios';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token interceptor
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;
```

### 2. Create Auth Context
```typescript
// contexts/AuthContext.tsx
import { createContext, useContext, useState, useEffect } from 'react';
import api from '@/lib/api';
import { User } from '@/types/api';

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const savedToken = localStorage.getItem('token');
    if (savedToken) {
      setToken(savedToken);
      fetchProfile();
    } else {
      setIsLoading(false);
    }
  }, []);

  const fetchProfile = async () => {
    try {
      const res = await api.get('/users/profile');
      setUser(res.data.data);
    } catch (error) {
      localStorage.removeItem('token');
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (email: string, password: string) => {
    const res = await api.post('/auth/login', { email, password });
    const { token, user } = res.data.data;
    localStorage.setItem('token', token);
    setToken(token);
    setUser(user);
  };

  const logout = () => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, token, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
```

### 3. Example API Hooks
```typescript
// hooks/useSearch.ts
import { useState } from 'react';
import api from '@/lib/api';
import { PlaceResult } from '@/types/api';

export const usePlaceSearch = () => {
  const [results, setResults] = useState<PlaceResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const search = async (query: string) => {
    setLoading(true);
    setError(null);
    try {
      const res = await api.get('/search/places', { params: { q: query } });
      setResults(res.data.data.results);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Search failed');
    } finally {
      setLoading(false);
    }
  };

  return { results, loading, error, search };
};
```

---

## Contact & Support
- GitHub: [Project Repository]
- Documentation: This folder
- API Testing: Use Postman or Insomnia
