# 📚 **Documentação Completa - Payloads da API**

## 🔐 **Autenticação**

### **POST /api/v1/auth/login**

**URL:** `http://localhost:8080/api/v1/auth/login`

**Headers:**

```
Content-Type: application/json
```

**Payload:**

```json
{
  "email": "admin@erp.com",
  "password": "admin123"
}
```

**Response (200):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "8936bd32-a5d6-4360-9b27-e58d828bdb61",
    "email": "admin@erp.com",
    "name": "Admin User",
    "role": "admin",
    "is_active": true,
    "created_at": "2025-08-08T01:46:42.568281Z",
    "updated_at": "2025-08-08T02:00:38.420722Z",
    "last_login_at": "2025-08-08T02:00:38.419659Z"
  }
}
```

### **POST /api/v1/auth/refresh**

**URL:** `http://localhost:8080/api/v1/auth/refresh`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {access_token}
```

**Payload:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## 👥 **Usuários**

### **POST /api/v1/users/register**

**URL:** `http://localhost:8080/api/v1/users/register`

**Headers:**

```
Content-Type: application/json
```

**Payload:**

```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe",
  "role": "user"
}
```

**Response (201):**

```json
{
  "id": "22efff77-e75f-4347-8b87-de5e43fa3f3d",
  "email": "user@example.com",
  "name": "John Doe",
  "role": "user",
  "is_active": true,
  "created_at": "2025-08-08T02:00:42.531Z",
  "updated_at": "2025-08-08T02:00:42.531Z"
}
```

### **GET /api/v1/users/profile**

**URL:** `http://localhost:8080/api/v1/users/profile`

**Headers:**

```
Authorization: Bearer {access_token}
```

**Response (200):**

```json
{
  "id": "8936bd32-a5d6-4360-9b27-e58d828bdb61",
  "email": "admin@erp.com",
  "name": "Admin User",
  "role": "admin",
  "is_active": true,
  "created_at": "2025-08-08T01:46:42.568281Z",
  "updated_at": "2025-08-08T02:00:38.420722Z",
  "last_login_at": "2025-08-08T02:00:38.419659Z"
}
```

---

## 👥 **Clientes**

### **POST /api/v1/clients**

**URL:** `http://localhost:8080/api/v1/clients`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {access_token}
```

**Payload:**

```json
{
  "name": "João Silva",
  "email": "joao@email.com",
  "phone": "(11) 99999-9999",
  "document": "123.456.789-00",
  "document_type": "CPF",
  "address": "Rua das Flores, 123",
  "city": "São Paulo",
  "state": "SP",
  "zip_code": "01234-567"
}
```

**Response (201):**

```json
{
  "id": "dbaad96a-66a4-44c2-a84e-1f9abaf1fb63",
  "name": "João Silva",
  "email": "joao@email.com",
  "phone": "(11) 99999-9999",
  "document": "123.456.789-00",
  "document_type": "CPF",
  "address": "Rua das Flores, 123",
  "city": "São Paulo",
  "state": "SP",
  "zip_code": "01234-567",
  "is_active": true,
  "created_at": "2025-08-08T02:16:24.060033588Z",
  "updated_at": "2025-08-08T02:16:24.060033588Z"
}
```

### **GET /api/v1/clients**

**URL:** `http://localhost:8080/api/v1/clients?limit=10&offset=0`

**Headers:**

```
Authorization: Bearer {access_token}
```

**Response (200):**

```json
{
  "clients": [
    {
      "id": "c3011a54-97b4-453a-8377-f569f4229316",
      "name": "Empresa ABC Ltda",
      "email": "contato@empresaabc.com",
      "phone": "(11) 88888-8888",
      "document": "12.345.678/0001-90",
      "document_type": "CNPJ",
      "address": "",
      "city": "",
      "state": "",
      "zip_code": "",
      "is_active": true,
      "created_at": "2025-08-08T02:16:32Z",
      "updated_at": "2025-08-08T02:16:32Z"
    },
    {
      "id": "dbaad96a-66a4-44c2-a84e-1f9abaf1fb63",
      "name": "João Silva",
      "email": "joao@email.com",
      "phone": "(11) 99999-9999",
      "document": "123.456.789-00",
      "document_type": "CPF",
      "address": "",
      "city": "",
      "state": "",
      "zip_code": "",
      "is_active": true,
      "created_at": "2025-08-08T02:16:24Z",
      "updated_at": "2025-08-08T02:16:24Z"
    }
  ],
  "total": 2,
  "limit": 10,
  "offset": 0
}
```

### **PUT /api/v1/clients/:id**

**URL:** `http://localhost:8080/api/v1/clients/dbaad96a-66a4-44c2-a84e-1f9abaf1fb63`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {access_token}
```

**Payload:**

```json
{
  "name": "João Silva Atualizado",
  "email": "joao.novo@email.com",
  "phone": "(11) 88888-8888"
}
```

**Response (200):**

```json
{
  "id": "dbaad96a-66a4-44c2-a84e-1f9abaf1fb63",
  "name": "João Silva Atualizado",
  "email": "joao.novo@email.com",
  "phone": "(11) 88888-8888",
  "document": "123.456.789-00",
  "document_type": "CPF",
  "address": "",
  "city": "",
  "state": "",
  "zip_code": "",
  "is_active": true,
  "created_at": "2025-08-08T02:16:24.060033588Z",
  "updated_at": "2025-08-08T02:20:15.123456789Z"
}
```

---

## 📦 **Produtos**

### **POST /api/v1/products**

**URL:** `http://localhost:8080/api/v1/products`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {access_token}
```

**Payload:**

```json
{
  "name": "Mármore Branco Carrara",
  "description": "Mármore branco de alta qualidade",
  "type": "Mármore",
  "price": 150.0,
  "unit": "m²"
}
```

**Response (201):**

```json
{
  "id": "7d2982e8-9ad2-4d08-8a7b-386ac98a2ddb",
  "name": "Mármore Branco Carrara",
  "description": "Mármore branco de alta qualidade",
  "type": "Mármore",
  "price": 150,
  "unit": "m²",
  "is_active": true,
  "created_at": "2025-08-08T02:20:30.155928924Z",
  "updated_at": "2025-08-08T02:20:30.155928924Z"
}
```

### **GET /api/v1/products**

**URL:** `http://localhost:8080/api/v1/products?limit=10&offset=0`

**Headers:**

```
Authorization: Bearer {access_token}
```

**Response (200):**

```json
{
  "products": [
    {
      "id": "b8957b54-39c6-42fa-951b-aa1d59e5a88d",
      "name": "Instalação de Pia",
      "description": "Serviço de instalação de pia",
      "type": "Serviço",
      "price": 80,
      "unit": "unit",
      "is_active": true,
      "created_at": "2025-08-08T02:20:47Z",
      "updated_at": "2025-08-08T02:20:47Z"
    },
    {
      "id": "792b9454-3e79-4b48-97f0-10f4408c3e04",
      "name": "Granito Preto Absoluto",
      "description": "Granito preto para bancadas",
      "type": "Granito",
      "price": 200,
      "unit": "m²",
      "is_active": true,
      "created_at": "2025-08-08T02:20:39Z",
      "updated_at": "2025-08-08T02:20:39Z"
    },
    {
      "id": "7d2982e8-9ad2-4d08-8a7b-386ac98a2ddb",
      "name": "Mármore Branco Carrara",
      "description": "Mármore branco de alta qualidade",
      "type": "Mármore",
      "price": 150,
      "unit": "m²",
      "is_active": true,
      "created_at": "2025-08-08T02:20:30Z",
      "updated_at": "2025-08-08T02:20:30Z"
    }
  ],
  "total": 3,
  "limit": 10,
  "offset": 0
}
```

---

## 📋 **Orçamentos**

### **POST /api/v1/quotes**

**URL:** `http://localhost:8080/api/v1/quotes`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {access_token}
```

**Payload:**

```json
{
  "client_id": "dbaad96a-66a4-44c2-a84e-1f9abaf1fb63",
  "date": "2025-01-14T00:00:00Z",
  "valid_until": "2025-02-14T00:00:00Z",
  "notes": "Orçamento para cozinha",
  "items": [
    {
      "product_id": "7d2982e8-9ad2-4d08-8a7b-386ac98a2ddb",
      "quantity": 5.5
    },
    {
      "product_id": "792b9454-3e79-4b48-97f0-10f4408c3e04",
      "quantity": 2.0
    }
  ]
}
```

**Response (201):**

```json
{
  "id": "f1407fb3-7735-4660-9209-400a249b94ff",
  "client_id": "dbaad96a-66a4-44c2-a84e-1f9abaf1fb63",
  "total_value": 1225,
  "status": "Pendente",
  "date": "2025-01-14T00:00:00Z",
  "valid_until": "2025-02-14T00:00:00Z",
  "notes": "Orçamento para cozinha",
  "is_active": true,
  "created_at": "2025-08-08T03:08:54.458169587Z",
  "updated_at": "2025-08-08T03:08:54.458169587Z"
}
```

### **GET /api/v1/quotes**

**URL:** `http://localhost:8080/api/v1/quotes?limit=10&offset=0`

**Headers:**

```
Authorization: Bearer {access_token}
```

**Response (200):**

```json
{
  "quotes": [
    {
      "id": "f1407fb3-7735-4660-9209-400a249b94ff",
      "client_id": "dbaad96a-66a4-44c2-a84e-1f9abaf1fb63",
      "client": {
        "id": "dbaad96a-66a4-44c2-a84e-1f9abaf1fb63",
        "name": "João Silva"
      },
      "total_value": 550,
      "status": "Pendente",
      "date": "2025-01-14",
      "valid_until": "2025-02-14",
      "notes": "",
      "is_active": true,
      "created_at": "2025-08-08T03:08:54Z",
      "updated_at": "2025-08-08T03:08:54Z"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

### **PUT /api/v1/quotes/:id/status**

**URL:** `http://localhost:8080/api/v1/quotes/f1407fb3-7735-4660-9209-400a249b94ff/status`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {access_token}
```

**Payload:**

```json
{
  "status": "Aprovado"
}
```

**Response (200):**

```json
{
  "message": "Quote status updated successfully"
}
```

**Status válidos:** `"Pendente"`, `"Aprovado"`, `"Rejeitado"`, `"Cancelado"`

---

## ⚙️ **Configurações**

### **GET /api/v1/settings**

**URL:** `http://localhost:8080/api/v1/settings`

**Headers:**

```
Authorization: Bearer {access_token}
```

**Response (200):**

```json
{
  "id": "2eaaa018-a831-4cb4-8442-ba86b73e730f",
  "trade_name": "Marmoraria Exemplo",
  "legal_name": "Marmoraria Exemplo LTI",
  "cnpj": "00.000.000/0000-00",
  "phone": "(00) 0000-0000",
  "email": "contato@marmorariaexemplo.com",
  "street": "Rua Exemplo",
  "number": "123",
  "neighborhood": "Centro",
  "city": "São Paulo",
  "state": "SP",
  "zip_code": "00000-000",
  "primary_color": "#1976d2",
  "secondary_color": "#9c27b0",
  "logo_url": "/api/placeholder/200/8C",
  "created_at": "2025-08-08T03:08:45.056442847Z",
  "updated_at": "2025-08-08T03:08:45.056442847Z"
}
```

### **PUT /api/v1/settings**

**URL:** `http://localhost:8080/api/v1/settings`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {access_token}
```

**Payload:**

```json
{
  "trade_name": "Marmoraria Exemplo Atualizada",
  "legal_name": "Marmoraria Exemplo LTI Atualizada",
  "cnpj": "11.111.111/0001-11",
  "phone": "(11) 1111-1111",
  "email": "contato@marmorariaexemplo.com.br",
  "street": "Rua Atualizada",
  "number": "456",
  "neighborhood": "Vila Nova",
  "city": "São Paulo",
  "state": "SP",
  "zip_code": "11111-111",
  "primary_color": "#2196f3",
  "secondary_color": "#e91e63",
  "logo_url": "/api/placeholder/200/FF"
}
```

**Response (200):**

```json
{
  "id": "2eaaa018-a831-4cb4-8442-ba86b73e730f",
  "trade_name": "Marmoraria Exemplo Atualizada",
  "legal_name": "Marmoraria Exemplo LTI Atualizada",
  "cnpj": "11.111.111/0001-11",
  "phone": "(11) 1111-1111",
  "email": "contato@marmorariaexemplo.com.br",
  "street": "Rua Atualizada",
  "number": "456",
  "neighborhood": "Vila Nova",
  "city": "São Paulo",
  "state": "SP",
  "zip_code": "11111-111",
  "primary_color": "#2196f3",
  "secondary_color": "#e91e63",
  "logo_url": "/api/placeholder/200/FF",
  "created_at": "2025-08-08T03:08:45.056442847Z",
  "updated_at": "2025-08-08T03:15:30.123456789Z"
}
```

---

## 🔧 **Configuração do Frontend**

### **Base URL:**

```
http://localhost:8080
```

### **Headers Padrão:**

```javascript
const headers = {
  "Content-Type": "application/json",
  Authorization: `Bearer ${accessToken}`,
};
```

### **Exemplo de Cliente JavaScript:**

```javascript
class ErpApiClient {
  constructor(baseURL = "http://localhost:8080") {
    this.baseURL = baseURL;
  }

  async request(endpoint, options = {}) {
    const token = localStorage.getItem("access_token");

    const config = {
      headers: {
        "Content-Type": "application/json",
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
      ...options,
    };

    const response = await fetch(`${this.baseURL}${endpoint}`, config);

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || "Request failed");
    }

    return response.json();
  }

  // Autenticação
  async login(email, password) {
    return this.request("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  }

  // Clientes
  async createClient(clientData) {
    return this.request("/api/v1/clients", {
      method: "POST",
      body: JSON.stringify(clientData),
    });
  }

  async getClients(limit = 10, offset = 0) {
    return this.request(`/api/v1/clients?limit=${limit}&offset=${offset}`);
  }

  // Produtos
  async createProduct(productData) {
    return this.request("/api/v1/products", {
      method: "POST",
      body: JSON.stringify(productData),
    });
  }

  async getProducts(limit = 10, offset = 0) {
    return this.request(`/api/v1/products?limit=${limit}&offset=${offset}`);
  }

  // Orçamentos
  async createQuote(quoteData) {
    return this.request("/api/v1/quotes", {
      method: "POST",
      body: JSON.stringify(quoteData),
    });
  }

  async getQuotes(limit = 10, offset = 0) {
    return this.request(`/api/v1/quotes?limit=${limit}&offset=${offset}`);
  }

  async updateQuoteStatus(quoteId, status) {
    return this.request(`/api/v1/quotes/${quoteId}/status`, {
      method: "PUT",
      body: JSON.stringify({ status }),
    });
  }

  // Configurações
  async getSettings() {
    return this.request("/api/v1/settings");
  }

  async updateSettings(settingsData) {
    return this.request("/api/v1/settings", {
      method: "PUT",
      body: JSON.stringify(settingsData),
    });
  }
}
```

### **Exemplo de Uso:**

```javascript
const api = new ErpApiClient();

// Login
const auth = await api.login("admin@erp.com", "admin123");
localStorage.setItem("access_token", auth.access_token);

// Criar cliente
const client = await api.createClient({
  name: "João Silva",
  email: "joao@email.com",
  phone: "(11) 99999-9999",
  document: "123.456.789-00",
  document_type: "CPF",
});

// Criar produto
const product = await api.createProduct({
  name: "Mármore Branco Carrara",
  description: "Mármore branco de alta qualidade",
  type: "Mármore",
  price: 150.0,
  unit: "m²",
});

// Criar orçamento
const quote = await api.createQuote({
  client_id: client.id,
  date: "2025-01-14T00:00:00Z",
  valid_until: "2025-02-14T00:00:00Z",
  items: [
    {
      product_id: product.id,
      quantity: 5.5,
    },
  ],
});
```

---

## 🚨 **Códigos de Erro**

### **400 Bad Request:**

```json
{
  "error": "Invalid request body",
  "details": "Field validation failed"
}
```

### **401 Unauthorized:**

```json
{
  "error": "Invalid token"
}
```

### **404 Not Found:**

```json
{
  "error": "Client not found"
}
```

### **409 Conflict:**

```json
{
  "error": "Client already exists"
}
```

### **500 Internal Server Error:**

```json
{
  "error": "Database connection failed"
}
```

---

## 📝 **Notas Importantes**

1. **Autenticação:** Sempre inclua o token JWT no header `Authorization`
2. **Content-Type:** Sempre use `application/json` para POST/PUT
3. **IDs:** Todos os IDs são UUIDs
4. **Datas:** Use formato ISO 8601 (`2025-01-14T00:00:00Z`)
5. **Paginação:** Use `limit` e `offset` para listagens
6. **Status de Orçamentos:** `"Pendente"`, `"Aprovado"`, `"Rejeitado"`, `"Cancelado"`
7. **Tipos de Produtos:** `"Mármore"`, `"Granito"`, `"Serviço"`
8. **Tipos de Documentos:** `"CPF"`, `"CNPJ"`
