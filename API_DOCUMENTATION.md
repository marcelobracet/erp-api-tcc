# 🚀 API Documentation - ERP System

## 📋 Base URL

```
Development: http://localhost:8080
QA: https://qa-api.seudominio.com
Production: https://api.seudominio.com
```

## 🔐 Authentication

### **JWT Token Structure**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 86400
}
```

### **Token Claims**

```json
{
  "user_id": "uuid-do-usuario",
  "email": "user@example.com",
  "role": "admin|user|manager",
  "exp": 1234567890,
  "iat": 1234567890
}
```

## 📡 Endpoints

### **1. Authentication**

#### **POST /api/v1/auth/login**

**Login do usuário**

**Request:**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 86400,
  "user": {
    "id": "8936bd32-a5d6-4360-9b27-e58d828bdb61",
    "email": "user@example.com",
    "name": "User Name",
    "role": "admin",
    "is_active": true,
    "created_at": "2025-08-08T01:46:42.568281417Z",
    "updated_at": "2025-08-08T01:46:42.568281417Z",
    "last_login_at": "2025-08-08T01:46:42.568281417Z"
  }
}
```

**Error Response (401):**

```json
{
  "error": "Invalid credentials"
}
```

#### **POST /api/v1/auth/refresh**

**Renovar access token**

**Request:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 86400
}
```

### **2. Users**

#### **POST /api/v1/users/register**

**Registrar novo usuário**

**Request:**

```json
{
  "email": "newuser@example.com",
  "password": "password123",
  "name": "New User",
  "role": "user"
}
```

**Response (201):**

```json
{
  "id": "8936bd32-a5d6-4360-9b27-e58d828bdb61",
  "email": "newuser@example.com",
  "name": "New User",
  "role": "user",
  "is_active": true,
  "created_at": "2025-08-08T01:46:42.568281417Z",
  "updated_at": "2025-08-08T01:46:42.568281417Z"
}
```

**Error Response (400):**

```json
{
  "error": "Validation failed",
  "details": {
    "email": "Email is required",
    "password": "Password must be at least 6 characters"
  }
}
```

#### **GET /api/v1/users/profile**

**Obter perfil do usuário logado**

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response (200):**

```json
{
  "id": "8936bd32-a5d6-4360-9b27-e58d828bdb61",
  "email": "user@example.com",
  "name": "User Name",
  "role": "admin",
  "is_active": true,
  "created_at": "2025-08-08T01:46:42.568281417Z",
  "updated_at": "2025-08-08T01:46:42.568281417Z",
  "last_login_at": "2025-08-08T01:46:42.568281417Z"
}
```

#### **GET /api/v1/users/:id**

**Obter usuário por ID (Admin only)**

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response (200):**

```json
{
  "id": "8936bd32-a5d6-4360-9b27-e58d828bdb61",
  "email": "user@example.com",
  "name": "User Name",
  "role": "admin",
  "is_active": true,
  "created_at": "2025-08-08T01:46:42.568281417Z",
  "updated_at": "2025-08-08T01:46:42.568281417Z",
  "last_login_at": "2025-08-08T01:46:42.568281417Z"
}
```

#### **PUT /api/v1/users/:id**

**Atualizar usuário (Admin only)**

**Headers:**

```
Authorization: Bearer <access_token>
```

**Request:**

```json
{
  "email": "updated@example.com",
  "name": "Updated Name",
  "role": "manager",
  "is_active": true
}
```

**Response (200):**

```json
{
  "id": "8936bd32-a5d6-4360-9b27-e58d828bdb61",
  "email": "updated@example.com",
  "name": "Updated Name",
  "role": "manager",
  "is_active": true,
  "created_at": "2025-08-08T01:46:42.568281417Z",
  "updated_at": "2025-08-08T01:47:42.568281417Z",
  "last_login_at": "2025-08-08T01:46:42.568281417Z"
}
```

#### **DELETE /api/v1/users/:id**

**Deletar usuário (Admin only)**

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response (204):**

```
No Content
```

#### **GET /api/v1/users**

**Listar usuários (Admin only)**

**Headers:**

```
Authorization: Bearer <access_token>
```

**Query Parameters:**

```
?limit=10&offset=0
```

**Response (200):**

```json
{
  "users": [
    {
      "id": "8936bd32-a5d6-4360-9b27-e58d828bdb61",
      "email": "user@example.com",
      "name": "User Name",
      "role": "admin",
      "is_active": true,
      "created_at": "2025-08-08T01:46:42.568281417Z",
      "updated_at": "2025-08-08T01:46:42.568281417Z",
      "last_login_at": "2025-08-08T01:46:42.568281417Z"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

#### **GET /api/v1/users/count**

**Contar usuários (Admin only)**

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response (200):**

```json
{
  "count": 5
}
```

### **3. Health Check**

#### **GET /health**

**Verificar status da API**

**Response (200):**

```json
{
  "status": "ok",
  "time": "2025-08-08T01:46:42.568281417Z"
}
```

## 🔧 Error Responses

### **400 Bad Request**

```json
{
  "error": "Validation failed",
  "details": {
    "field": "error message"
  }
}
```

### **401 Unauthorized**

```json
{
  "error": "Invalid credentials"
}
```

### **403 Forbidden**

```json
{
  "error": "Insufficient permissions"
}
```

### **404 Not Found**

```json
{
  "error": "User not found"
}
```

### **500 Internal Server Error**

```json
{
  "error": "Internal server error"
}
```

## 📝 Frontend Integration Guide

### **1. Authentication Flow**

#### **Login Process:**

```javascript
// 1. Login
const login = async (email, password) => {
  try {
    const response = await fetch("http://localhost:8080/api/v1/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    });

    if (response.ok) {
      const data = await response.json();

      // Store tokens
      localStorage.setItem("access_token", data.access_token);
      localStorage.setItem("refresh_token", data.refresh_token);
      localStorage.setItem("user", JSON.stringify(data.user));

      return data;
    } else {
      throw new Error("Login failed");
    }
  } catch (error) {
    console.error("Login error:", error);
    throw error;
  }
};
```

#### **Token Management:**

```javascript
// 2. API Client with automatic token refresh
class ApiClient {
  constructor(baseURL) {
    this.baseURL = baseURL;
  }

  async request(endpoint, options = {}) {
    const token = localStorage.getItem("access_token");

    const config = {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
    };

    try {
      const response = await fetch(`${this.baseURL}${endpoint}`, config);

      if (response.status === 401) {
        // Token expired, try to refresh
        const refreshed = await this.refreshToken();
        if (refreshed) {
          // Retry request with new token
          config.headers.Authorization = `Bearer ${localStorage.getItem(
            "access_token"
          )}`;
          return await fetch(`${this.baseURL}${endpoint}`, config);
        }
      }

      return response;
    } catch (error) {
      console.error("API request error:", error);
      throw error;
    }
  }

  async refreshToken() {
    try {
      const refreshToken = localStorage.getItem("refresh_token");
      const response = await fetch(`${this.baseURL}/api/v1/auth/refresh`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (response.ok) {
        const data = await response.json();
        localStorage.setItem("access_token", data.access_token);
        localStorage.setItem("refresh_token", data.refresh_token);
        return true;
      }

      // Refresh failed, redirect to login
      this.logout();
      return false;
    } catch (error) {
      this.logout();
      return false;
    }
  }

  logout() {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user");
    // Redirect to login page
    window.location.href = "/login";
  }
}
```

#### **Protected Routes:**

```javascript
// 3. Route Protection
const ProtectedRoute = ({ children, requiredRole }) => {
  const user = JSON.parse(localStorage.getItem("user") || "{}");
  const token = localStorage.getItem("access_token");

  if (!token) {
    return <Navigate to="/login" />;
  }

  if (requiredRole && user.role !== requiredRole) {
    return <Navigate to="/unauthorized" />;
  }

  return children;
};
```

### **2. API Usage Examples**

#### **User Registration:**

```javascript
const registerUser = async (userData) => {
  const response = await apiClient.request("/api/v1/users/register", {
    method: "POST",
    body: JSON.stringify(userData),
  });

  return response.json();
};
```

#### **Get User Profile:**

```javascript
const getUserProfile = async () => {
  const response = await apiClient.request("/api/v1/users/profile");
  return response.json();
};
```

#### **List Users (Admin):**

```javascript
const listUsers = async (limit = 10, offset = 0) => {
  const response = await apiClient.request(
    `/api/v1/users?limit=${limit}&offset=${offset}`
  );
  return response.json();
};
```

### **3. Environment Configuration**

#### **Development:**

```javascript
// .env.development
REACT_APP_API_URL=http://localhost:8080
REACT_APP_ENVIRONMENT=development
```

#### **QA:**

```javascript
// .env.qa
REACT_APP_API_URL=https://qa-api.seudominio.com
REACT_APP_ENVIRONMENT=qa
```

#### **Production:**

```javascript
// .env.production
REACT_APP_API_URL=https://api.seudominio.com
REACT_APP_ENVIRONMENT=production
```

### **4. Error Handling**

```javascript
const handleApiError = (error) => {
  if (error.status === 401) {
    // Unauthorized - redirect to login
    window.location.href = "/login";
  } else if (error.status === 403) {
    // Forbidden - show access denied
    showNotification("Access denied", "error");
  } else if (error.status === 404) {
    // Not found
    showNotification("Resource not found", "error");
  } else {
    // Generic error
    showNotification("An error occurred", "error");
  }
};
```

## 🔒 Security Considerations

### **1. Token Storage**

- Store tokens in `localStorage` for development
- Use `httpOnly` cookies in production
- Implement token rotation

### **2. CORS Configuration**

```javascript
// Backend CORS configuration
app.use(
  cors({
    origin: ["http://localhost:3000", "https://your-frontend-domain.com"],
    credentials: true,
  })
);
```

### **3. Input Validation**

- Always validate input on frontend
- Use proper validation libraries
- Sanitize user input

### **4. Error Messages**

- Don't expose sensitive information in error messages
- Log errors on backend
- Show user-friendly messages on frontend

## 📊 Testing Examples

### **cURL Commands**

#### **Login:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@erp.com", "password": "admin123"}'
```

#### **Register:**

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123", "name": "Test User", "role": "user"}'
```

#### **Get Profile:**

```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### **Postman Collection**

```json
{
  "info": {
    "name": "ERP API",
    "description": "ERP System API endpoints"
  },
  "item": [
    {
      "name": "Authentication",
      "item": [
        {
          "name": "Login",
          "request": {
            "method": "POST",
            "url": "{{base_url}}/api/v1/auth/login",
            "body": {
              "mode": "raw",
              "raw": "{\n  \"email\": \"admin@erp.com\",\n  \"password\": \"admin123\"\n}",
              "options": {
                "raw": {
                  "language": "json"
                }
              }
            }
          }
        }
      ]
    }
  ]
}
```

## 🚀 Deployment Checklist

### **Frontend:**

- [ ] Configure environment variables
- [ ] Set up API client with base URL
- [ ] Implement authentication flow
- [ ] Add error handling
- [ ] Test all endpoints
- [ ] Configure CORS

### **Backend:**

- [ ] Deploy to QA environment
- [ ] Test all endpoints
- [ ] Configure CORS for frontend domain
- [ ] Set up monitoring
- [ ] Deploy to production

---

**📞 Support:** For any questions or issues, contact the development team.
