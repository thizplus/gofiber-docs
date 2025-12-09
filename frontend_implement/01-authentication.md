# Part 1: Authentication (Register / Login)

## Overview
ระบบ Authentication สำหรับการลงทะเบียนและเข้าสู่ระบบ

## Base URL
```
/api/v1/auth
```

---

## 1.1 Register (ลงทะเบียน)

### Endpoint
```
POST /api/v1/auth/register
```

### Authentication
ไม่ต้อง (Public)

### Request Body
```typescript
interface RegisterRequest {
  email: string;      // required, email format, max 255
  username: string;   // required, 3-20 chars, alphanumeric only
  password: string;   // required, 8-72 chars
  firstName: string;  // required, 1-50 chars
  lastName: string;   // required, 1-50 chars
}
```

### Example Request
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "SecurePass123",
  "firstName": "John",
  "lastName": "Doe"
}
```

### Response
```typescript
interface RegisterResponse {
  success: boolean;
  message: string;
  data: {
    token: string;
    user: UserResponse;
  };
}

interface UserResponse {
  id: string;          // UUID
  email: string;
  username: string;
  firstName: string;
  lastName: string;
  avatar: string;
  role: string;        // "user" | "admin"
  isActive: boolean;
  studentId?: string;  // optional
  createdAt: string;   // ISO 8601
  updatedAt: string;   // ISO 8601
}
```

### Example Response (Success)
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "username": "johndoe",
      "firstName": "John",
      "lastName": "Doe",
      "avatar": "",
      "role": "user",
      "isActive": true,
      "createdAt": "2024-01-15T10:30:00Z",
      "updatedAt": "2024-01-15T10:30:00Z"
    }
  }
}
```

### Error Responses
```json
// Email already exists
{
  "success": false,
  "message": "Email already exists",
  "error": "Bad Request"
}

// Validation error
{
  "success": false,
  "message": "Validation failed",
  "error": "password: must be at least 8 characters"
}
```

---

## 1.2 Login (เข้าสู่ระบบ)

### Endpoint
```
POST /api/v1/auth/login
```

### Authentication
ไม่ต้อง (Public)

### Request Body
```typescript
interface LoginRequest {
  email: string;     // required, email format, max 255
  password: string;  // required, min 1 char
}
```

### Example Request
```json
{
  "email": "user@example.com",
  "password": "SecurePass123"
}
```

### Response
```typescript
interface LoginResponse {
  success: boolean;
  message: string;
  data: {
    token: string;
    user: UserResponse;
  };
}
```

### Example Response (Success)
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "username": "johndoe",
      "firstName": "John",
      "lastName": "Doe",
      "avatar": "",
      "role": "user",
      "isActive": true,
      "createdAt": "2024-01-15T10:30:00Z",
      "updatedAt": "2024-01-15T10:30:00Z"
    }
  }
}
```

### Error Responses
```json
// Invalid credentials
{
  "success": false,
  "message": "Invalid email or password",
  "error": "Unauthorized"
}
```

---

## 1.3 Get Profile (ดูโปรไฟล์)

### Endpoint
```
GET /api/v1/users/profile
```

### Authentication
Required (Bearer Token)

### Headers
```
Authorization: Bearer <token>
```

### Response
```typescript
interface APIResponse<UserResponse> {
  success: boolean;
  message: string;
  data: UserResponse;
}
```

### Example Response
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "johndoe",
    "firstName": "John",
    "lastName": "Doe",
    "avatar": "https://example.com/avatar.jpg",
    "role": "user",
    "isActive": true,
    "studentId": "12345678901",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  }
}
```

---

## 1.4 Update Profile (แก้ไขโปรไฟล์)

### Endpoint
```
PUT /api/v1/users/profile
```

### Authentication
Required (Bearer Token)

### Request Body
```typescript
interface UpdateUserRequest {
  firstName?: string;  // optional, 1-50 chars
  lastName?: string;   // optional, 1-50 chars
  avatar?: string;     // optional, URL format, max 500 chars
}
```

### Example Request
```json
{
  "firstName": "Johnny",
  "lastName": "Doe",
  "avatar": "https://example.com/new-avatar.jpg"
}
```

### Response
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "johndoe",
    "firstName": "Johnny",
    "lastName": "Doe",
    "avatar": "https://example.com/new-avatar.jpg",
    "role": "user",
    "isActive": true,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T11:00:00Z"
  }
}
```

---

## 1.5 Delete Account (ลบบัญชี)

### Endpoint
```
DELETE /api/v1/users/profile
```

### Authentication
Required (Bearer Token)

### Response
```json
{
  "success": true,
  "message": "Account deleted successfully"
}
```

---

## TypeScript Types สำหรับ Frontend

```typescript
// types/auth.ts

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface LoginRequest {
  email: string;
  password: string;
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

export interface AuthResponse {
  token: string;
  user: User;
}

export interface UpdateProfileRequest {
  firstName?: string;
  lastName?: string;
  avatar?: string;
}

// API Response wrapper
export interface APIResponse<T> {
  success: boolean;
  message: string;
  data?: T;
  error?: string;
}
```

---

## API Routes Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/auth/register` | No | ลงทะเบียนผู้ใช้ใหม่ |
| POST | `/api/v1/auth/login` | No | เข้าสู่ระบบ |
| GET | `/api/v1/users/profile` | Yes | ดูโปรไฟล์ |
| PUT | `/api/v1/users/profile` | Yes | แก้ไขโปรไฟล์ |
| DELETE | `/api/v1/users/profile` | Yes | ลบบัญชี |

---

## Notes
- Token ใช้รูปแบบ JWT (JSON Web Token)
- Token ต้องส่งใน Header: `Authorization: Bearer <token>`
- Token มีอายุ 24 ชั่วโมง
- รหัสผ่านต้องมีความยาวอย่างน้อย 8 ตัวอักษร
- Username ต้องเป็น alphanumeric เท่านั้น (a-z, A-Z, 0-9)
