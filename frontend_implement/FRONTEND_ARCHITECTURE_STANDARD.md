# Frontend Architecture Standard

## Overview

โครงสร้างนี้เป็น **Feature-Based Architecture** สำหรับ Next.js 15+ พร้อม TypeScript, React Query, และ Zustand

---

## 1. Folder Structure

```
project/
├── app/                          # Next.js App Router (pages & layouts)
├── src/
│   ├── features/                 # Feature modules (self-contained)
│   ├── shared/                   # Shared utilities & components
│   ├── services/                 # API service layer
│   └── providers/                # Global React providers
├── public/                       # Static assets
└── tsconfig.json                 # Path aliases
```

### Path Aliases (tsconfig.json)

```json
{
  "compilerOptions": {
    "paths": {
      "@/features/*": ["./src/features/*"],
      "@/shared/*": ["./src/shared/*"],
      "@/services/*": ["./src/services/*"],
      "@/types/*": ["./src/shared/types/*"],
      "@/hooks/*": ["./src/shared/hooks/*"],
      "@/components/*": ["./src/shared/components/*"],
      "@/lib/*": ["./src/shared/lib/*"],
      "@/config/*": ["./src/shared/config/*"]
    }
  }
}
```

---

## 2. Feature Module Structure

แต่ละ feature เป็น **self-contained module** มีโครงสร้างดังนี้:

```
src/features/[feature-name]/
├── components/              # UI components เฉพาะ feature
│   ├── Component.tsx
│   └── index.ts            # Barrel export
├── hooks/                   # React Query hooks + business logic
│   ├── useFeature.ts
│   └── index.ts
├── stores/                  # Zustand stores (ถ้าต้องการ)
│   ├── featureStore.ts
│   └── actions/            # แยก actions ตาม concern
├── utils/                   # Utilities เฉพาะ feature
├── types.ts                 # Types เฉพาะ feature (optional)
└── index.ts                 # PUBLIC API - export เฉพาะที่ต้องการ expose
```

### ตัวอย่าง Feature

```
src/features/posts/
├── components/
│   ├── PostCard.tsx
│   ├── PostFeed.tsx
│   ├── CreatePostForm.tsx
│   └── index.ts
├── hooks/
│   ├── usePosts.ts          # Query: list, detail
│   ├── usePostMutations.ts  # Mutations: create, update, delete
│   ├── useVotes.ts          # Vote logic
│   └── index.ts
├── stores/
│   └── optimisticPostStore.ts
├── utils/
│   └── voteCalculations.ts
└── index.ts
```

---

## 3. Shared Layer

สิ่งที่ใช้ร่วมกันทุก feature:

```
src/shared/
├── components/
│   ├── common/              # EmptyState, LoadingState, ErrorState
│   ├── layouts/             # AppLayout, ChatLayout
│   ├── media/               # MediaDisplay, VideoPlayer
│   ├── navigation/          # Sidebar, NavBar
│   └── ui/                  # Button, Card, Dialog (Radix wrappers)
│
├── hooks/
│   ├── useHydration.ts      # SSR hydration check
│   ├── useMobile.ts         # Mobile detection
│   └── useAuthGuard.ts      # Auth protection
│
├── lib/
│   ├── api/
│   │   ├── http-client.ts   # Axios instances + interceptors
│   │   └── constants/
│   │       └── api.ts       # All API endpoints
│   ├── storage/             # IndexedDB, localStorage
│   ├── upload/              # Upload utilities
│   └── utils/               # General utilities (cn, date, etc.)
│
├── types/
│   ├── common.ts            # ApiResponse, Pagination, Enums
│   ├── models.ts            # User, Post, Comment, etc.
│   ├── request.ts           # API request types
│   ├── response.ts          # API response types
│   └── index.ts             # Barrel export
│
└── config/
    ├── constants.ts         # PAGINATION, FORM_LIMITS, UI, etc.
    ├── validation.ts        # Validation rules
    └── index.ts
```

---

## 4. Services Layer

Services จัดกลุ่มตาม **Microservice** ไม่ใช่ตาม feature:

```
src/services/
├── index.ts                 # Main exports
├── auth/                    # Auth Service (port 8088)
│   ├── auth.service.ts
│   ├── user.service.ts
│   └── index.ts
├── backend/                 # Backend Service (port 8080)
│   ├── post.service.ts
│   ├── comment.service.ts
│   ├── vote.service.ts
│   └── index.ts
├── upload/                  # Upload Service (port 8090)
│   ├── upload.service.ts
│   └── index.ts
└── chat/                    # Chat Service
    ├── chat.service.ts
    └── index.ts
```

### Service Pattern

```typescript
// src/services/backend/post.service.ts
import { backendClient } from '@/lib/api/http-client';
import { POST_API } from '@/lib/api/constants/api';
import type { Post, CreatePostRequest, ApiResponse } from '@/types';

export const postService = {
  // Read
  list: async (params?): Promise<ApiResponse<Post[]>> => {
    const { data } = await backendClient.get(POST_API.LIST, { params });
    return data;
  },

  getById: async (id: string): Promise<ApiResponse<Post>> => {
    const { data } = await backendClient.get(POST_API.DETAIL(id));
    return data;
  },

  // Write
  create: async (payload: CreatePostRequest): Promise<ApiResponse<Post>> => {
    const { data } = await backendClient.post(POST_API.CREATE, payload);
    return data;
  },

  update: async (id: string, payload): Promise<ApiResponse<Post>> => {
    const { data } = await backendClient.put(POST_API.UPDATE(id), payload);
    return data;
  },

  delete: async (id: string): Promise<ApiResponse<null>> => {
    const { data } = await backendClient.delete(POST_API.DELETE(id));
    return data;
  },
};
```

---

## 5. HTTP Client

```typescript
// src/shared/lib/api/http-client.ts
import axios from 'axios';
import { useAuthStore } from '@/features/auth';

const getToken = (): string | null => {
  return useAuthStore.getState().token;
};

// Auth Service Client
export const authClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_AUTH_API_URL,
  headers: { 'Content-Type': 'application/json' },
});

// Backend Service Client
export const backendClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_BACKEND_API_URL,
  headers: { 'Content-Type': 'application/json' },
});

// Request Interceptor - Add token
backendClient.interceptors.request.use((config) => {
  const token = getToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response Interceptor - Handle 401
backendClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().clearAuth();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

---

## 6. API Endpoints Configuration

```typescript
// src/shared/lib/api/constants/api.ts

// Base URLs
export const AUTH_BASE_URL = process.env.NEXT_PUBLIC_AUTH_API_URL;
export const BACKEND_BASE_URL = process.env.NEXT_PUBLIC_BACKEND_API_URL;
export const UPLOAD_BASE_URL = process.env.NEXT_PUBLIC_UPLOAD_API_URL;

// Auth Endpoints
export const AUTH_API = {
  LOGIN: '/auth/login',
  REGISTER: '/auth/register',
  LOGOUT: '/auth/logout',
  REFRESH: '/auth/refresh',
  GOOGLE_CALLBACK: '/auth/google/callback',
};

// Post Endpoints
export const POST_API = {
  LIST: '/posts',
  CREATE: '/posts',
  DETAIL: (id: string) => `/posts/${id}`,
  UPDATE: (id: string) => `/posts/${id}`,
  DELETE: (id: string) => `/posts/${id}`,
  BY_AUTHOR: (userId: string) => `/users/${userId}/posts`,
  BY_TAG: (tagName: string) => `/tags/${tagName}/posts`,
};

// Comment Endpoints
export const COMMENT_API = {
  LIST: (postId: string) => `/posts/${postId}/comments`,
  CREATE: (postId: string) => `/posts/${postId}/comments`,
  UPDATE: (id: string) => `/comments/${id}`,
  DELETE: (id: string) => `/comments/${id}`,
};

// ... other endpoints
```

---

## 7. Types Organization

### common.ts - Shared Types

```typescript
// src/shared/types/common.ts

// Standard API Response
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  error?: string;
}

// Pagination
export interface CursorPaginationMeta {
  hasMore: boolean;
  nextCursor?: string;
  prevCursor?: string;
  total?: number;
}

export interface CursorPaginatedResponse<T> {
  success: boolean;
  data: {
    items: T[];
    pagination: CursorPaginationMeta;
  };
}

// Enums
export enum VoteType {
  UP = 'up',
  DOWN = 'down',
}

export enum NotificationType {
  LIKE = 'like',
  COMMENT = 'comment',
  FOLLOW = 'follow',
  MENTION = 'mention',
}
```

### models.ts - Data Models

```typescript
// src/shared/types/models.ts

export interface User {
  id: string;
  username: string;
  displayName: string;
  email?: string;
  avatar: string | null;
  bio?: string | null;
  followersCount?: number;
  followingCount?: number;
  createdAt: string;
  updatedAt?: string;
}

export interface Post {
  id: string;
  title: string;
  content: string;
  author: User;
  tags: string[];
  media: Media[];
  upvotes: number;
  downvotes: number;
  commentCount: number;
  userVote?: VoteType | null;
  isSaved?: boolean;
  createdAt: string;
  updatedAt?: string;
}

export interface Comment {
  id: string;
  content: string;
  author: User;
  postId: string;
  parentId?: string | null;
  upvotes: number;
  downvotes: number;
  userVote?: VoteType | null;
  createdAt: string;
}

export interface Media {
  id: string;
  type: 'image' | 'video';
  url: string;
  thumbnailUrl?: string;
  width?: number;
  height?: number;
}
```

### request.ts - API Request Types

```typescript
// src/shared/types/request.ts

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  displayName?: string;
}

export interface CreatePostRequest {
  title: string;
  content: string;
  tags?: string[];
  mediaIds?: string[];
}

export interface UpdatePostRequest {
  title?: string;
  content?: string;
  tags?: string[];
}

export interface CreateCommentRequest {
  content: string;
  parentId?: string;
}

export interface VoteRequest {
  targetType: 'post' | 'comment';
  targetId: string;
  voteType: VoteType;
}
```

### response.ts - API Response Types

```typescript
// src/shared/types/response.ts

export interface LoginResponse {
  token: string;
  user: User;
}

export interface CreatePostResponse {
  post: Post;
}

export interface ListPostsResponse {
  posts: Post[];
  pagination: CursorPaginationMeta;
}
```

---

## 8. React Query Hooks

### Query Keys Pattern

```typescript
// src/features/posts/hooks/queryKeys.ts
export const postKeys = {
  all: ['posts'] as const,
  lists: () => [...postKeys.all, 'list'] as const,
  list: (params?: object) => [...postKeys.lists(), params] as const,
  details: () => [...postKeys.all, 'detail'] as const,
  detail: (id: string) => [...postKeys.details(), id] as const,
  byAuthor: (userId: string) => [...postKeys.all, 'author', userId] as const,
};
```

### Query Hook

```typescript
// src/features/posts/hooks/usePosts.ts
import { useQuery, useInfiniteQuery } from '@tanstack/react-query';
import { postService } from '@/services';
import { postKeys } from './queryKeys';

// Single item query
export function usePost(id: string) {
  return useQuery({
    queryKey: postKeys.detail(id),
    queryFn: () => postService.getById(id),
    enabled: !!id,
  });
}

// List query
export function usePosts(params?: GetPostsParams) {
  return useQuery({
    queryKey: postKeys.list(params),
    queryFn: () => postService.list(params),
  });
}

// Infinite query (pagination)
export function useInfinitePosts(params?: GetPostsParams) {
  return useInfiniteQuery({
    queryKey: postKeys.list(params),
    queryFn: ({ pageParam }) =>
      postService.list({ ...params, cursor: pageParam }),
    getNextPageParam: (lastPage) =>
      lastPage.data.pagination.hasMore
        ? lastPage.data.pagination.nextCursor
        : undefined,
    initialPageParam: undefined,
  });
}
```

### Mutation Hook

```typescript
// src/features/posts/hooks/usePostMutations.ts
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { postService } from '@/services';
import { postKeys } from './queryKeys';
import { toast } from 'sonner';

export function useCreatePost() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: postService.create,
    onSuccess: (response) => {
      // Invalidate list to refetch
      queryClient.invalidateQueries({ queryKey: postKeys.lists() });
      toast.success('Post created!');
    },
    onError: (error) => {
      toast.error('Failed to create post');
    },
  });
}

export function useDeletePost() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: postService.delete,
    onMutate: async (postId) => {
      // Optimistic update
      await queryClient.cancelQueries({ queryKey: postKeys.lists() });

      const previousPosts = queryClient.getQueryData(postKeys.lists());

      queryClient.setQueryData(postKeys.lists(), (old: any) => ({
        ...old,
        data: {
          ...old.data,
          items: old.data.items.filter((p: Post) => p.id !== postId),
        },
      }));

      return { previousPosts };
    },
    onError: (err, postId, context) => {
      // Rollback on error
      queryClient.setQueryData(postKeys.lists(), context?.previousPosts);
      toast.error('Failed to delete post');
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: postKeys.lists() });
    },
  });
}
```

---

## 9. Zustand Store

### Store Pattern

```typescript
// src/features/auth/stores/authStore.ts
import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import type { User } from '@/types';

interface AuthState {
  // State
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
  _hasHydrated: boolean;

  // Actions
  setAuth: (token: string, user: User) => void;
  clearAuth: () => void;
  setUser: (user: User) => void;
  setHasHydrated: (state: boolean) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      // Initial state
      token: null,
      user: null,
      isAuthenticated: false,
      _hasHydrated: false,

      // Actions
      setAuth: (token, user) => set({
        token,
        user,
        isAuthenticated: true,
      }),

      clearAuth: () => set({
        token: null,
        user: null,
        isAuthenticated: false,
      }),

      setUser: (user) => set({ user }),

      setHasHydrated: (state) => set({ _hasHydrated: state }),
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => localStorage),
      onRehydrateStorage: () => (state) => {
        state?.setHasHydrated(true);
      },
    }
  )
);

// Selector hooks (prevent unnecessary re-renders)
export const useUser = () => useAuthStore((state) => state.user);
export const useToken = () => useAuthStore((state) => state.token);
export const useIsAuthenticated = () => useAuthStore((state) => state.isAuthenticated);
export const useHasHydrated = () => useAuthStore((state) => state._hasHydrated);
```

### Complex Store with Actions Split

```typescript
// src/features/chat/stores/chatStore.ts
import { create } from 'zustand';
import { createConversationActions } from './actions/conversationActions';
import { createMessageActions } from './actions/messageActions';

interface ChatState {
  conversations: Conversation[];
  activeConversationId: string | null;
  messagesByConversation: Record<string, ChatMessage[]>;
  // ... more state
}

interface ChatActions {
  // Conversation actions
  setConversations: (conversations: Conversation[]) => void;
  addConversation: (conversation: Conversation) => void;
  setActiveConversation: (id: string | null) => void;

  // Message actions
  addMessage: (conversationId: string, message: ChatMessage) => void;
  updateMessage: (conversationId: string, messageId: string, updates: Partial<ChatMessage>) => void;
}

export const useChatStore = create<ChatState & ChatActions>()((set, get) => ({
  // Initial state
  conversations: [],
  activeConversationId: null,
  messagesByConversation: {},

  // Spread actions from separate files
  ...createConversationActions(set, get),
  ...createMessageActions(set, get),
}));
```

---

## 10. Constants & Configuration

```typescript
// src/shared/config/constants.ts

export const PAGINATION = {
  DEFAULT_LIMIT: 20,
  MESSAGE_LIMIT: 50,
  COMMENT_LIMIT: 10,
};

export const FORM_LIMITS = {
  POST: {
    TITLE_MAX: 300,
    CONTENT_MAX: 10000,
    TAGS_MAX: 5,
  },
  COMMENT: {
    CONTENT_MAX: 2000,
  },
  USER: {
    USERNAME_MIN: 3,
    USERNAME_MAX: 30,
    BIO_MAX: 500,
  },
};

export const UI = {
  TOAST_DURATION: 3000,
  DEBOUNCE_MS: 300,
  ANIMATION_DURATION: 200,
};

export const WEBSOCKET = {
  MAX_RECONNECT_ATTEMPTS: 5,
  RECONNECT_INTERVAL: 3000,
  PING_INTERVAL: 30000,
};

export const ROUTES = {
  HOME: '/',
  LOGIN: '/auth/login',
  REGISTER: '/auth/register',
  PROFILE: (username: string) => `/u/${username}`,
  POST: (id: string) => `/post/${id}`,
  SETTINGS: '/settings',
};
```

---

## 11. Data Flow

```
User Action (Click, Submit)
    ↓
React Component
    ↓
Custom Hook (usePosts, useCreatePost)
    ├─ useQuery (READ)
    └─ useMutation (WRITE)
    ↓
Service Layer (postService.create)
    ↓
HTTP Client (backendClient)
    ├─ Add Authorization header
    └─ Handle errors
    ↓
API Microservice
    ↓
Response
    ↓
React Query Cache Update
    ↓
Component Re-render
```

---

## 12. Import Patterns

```typescript
// Features - ใช้ public API จาก index.ts
import { usePost, useCreatePost, PostCard } from '@/features/posts';
import { useAuth, useAuthStore } from '@/features/auth';

// Shared Components
import { Button, Card, Dialog } from '@/components/ui';
import { AppLayout } from '@/components/layouts';
import { EmptyState, LoadingState } from '@/components/common';

// Shared Hooks
import { useHydration, useMobile } from '@/hooks';

// Utilities
import { cn } from '@/lib/utils';
import { PAGINATION, FORM_LIMITS } from '@/config';

// Types
import type { Post, User, Comment } from '@/types';
import type { CreatePostRequest } from '@/types/request';

// Services (ใช้ใน hooks เท่านั้น)
import { postService } from '@/services';
```

---

## 13. Best Practices

### Feature Isolation
- แต่ละ feature พัฒนาและ test แยกกันได้
- Export เฉพาะสิ่งที่ต้องการ expose ผ่าน `index.ts`
- ลด cross-feature dependencies

### Service Abstraction
- Components ไม่เรียก HTTP โดยตรง
- ใช้ service layer เป็น abstraction
- Services จัดกลุ่มตาม microservice

### Type Safety
- Type ทุก API (request & response)
- ใช้ Enum สำหรับ fixed values
- Optional fields สำหรับ backward compatibility

### Optimistic Updates
- ใช้ `onMutate` ใน React Query
- Rollback เมื่อเกิด error
- IndexedDB สำหรับ persist ระหว่าง upload

### Cache Invalidation
- ใช้ structured query keys
- Selective invalidation
- `onSuccess` callback

### Error Handling
- Centralized ใน HTTP interceptors
- Feature-level handling ใน components
- Toast notifications สำหรับ user feedback

---

## 14. Checklist สำหรับ Feature ใหม่

```
[ ] สร้าง folder structure ตาม template
[ ] สร้าง types (ถ้าต้องการ types เฉพาะ feature)
[ ] สร้าง service methods ใน src/services/
[ ] สร้าง query keys
[ ] สร้าง hooks (query + mutation)
[ ] สร้าง components
[ ] สร้าง store (ถ้าต้องการ local state)
[ ] สร้าง index.ts สำหรับ public API
[ ] เพิ่ม routes ใน app/
```

---

## สรุป

Architecture นี้ให้:
- **Modularity** - Features แยกกันชัดเจน
- **Scalability** - เพิ่ม feature ใหม่ง่าย
- **Maintainability** - Separation of concerns ชัดเจน
- **Type Safety** - TypeScript ครบถ้วน
- **Performance** - React Query caching, optimistic updates
- **DX** - Path aliases, barrel exports, clear patterns
