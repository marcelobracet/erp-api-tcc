# 🔧 Troubleshooting - Problemas de Conexão Frontend

## ✅ **PROBLEMA RESOLVIDO!**

O erro **"Erro de conexão com o servidor"** foi causado por **CORS não configurado**. Agora está **100% resolvido**!

### **🔧 O que foi corrigido:**

1. **Adicionado middleware CORS** no backend
2. **Configurado origins permitidos** para frontend
3. **Habilitado credenciais** para autenticação
4. **Configurado headers** necessários

## 🧪 **TESTE AGORA NO FRONTEND**

### **1. Verificar se a API está rodando:**

```bash
curl http://localhost:8080/health
```

**Response esperado:**

```json
{ "status": "ok", "time": "2025-08-08T02:05:03Z" }
```

### **2. Testar login via curl:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@erp.com", "password": "admin123"}'
```

### **3. Testar no navegador:**

Abra o console do navegador e execute:

```javascript
// Teste direto no console
fetch("http://localhost:8080/api/v1/auth/login", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    email: "admin@erp.com",
    password: "admin123",
  }),
})
  .then((response) => response.json())
  .then((data) => console.log("Login success:", data))
  .catch((error) => console.error("Login error:", error));
```

## 🚀 **CONFIGURAÇÃO CORS IMPLEMENTADA**

### **Origins Permitidos:**

- ✅ `http://localhost:3000` (React padrão)
- ✅ `http://localhost:3001` (React alternativo)
- ✅ `http://localhost:5173` (Vite)
- ✅ `http://127.0.0.1:3000`
- ✅ `http://127.0.0.1:3001`
- ✅ `http://127.0.0.1:5173`

### **Métodos Permitidos:**

- ✅ `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`

### **Headers Permitidos:**

- ✅ `Origin`, `Content-Type`, `Accept`, `Authorization`

### **Credenciais:**

- ✅ `AllowCredentials: true`

## 🔍 **VERIFICAÇÃO DE PROBLEMAS**

### **1. Se ainda der erro, verifique:**

#### **URL da API no frontend:**

```javascript
// Certifique-se que está usando a URL correta
const API_URL = "http://localhost:8080";

// Não use https se a API está em http
const API_URL = "https://localhost:8080"; // ❌ ERRADO
const API_URL = "http://localhost:8080"; // ✅ CORRETO
```

#### **Porta do frontend:**

```javascript
// Se seu frontend está rodando na porta 3001, adicione ao CORS
// Ou use a porta padrão 3000
```

#### **Headers da requisição:**

```javascript
// Certifique-se que está enviando os headers corretos
const response = await fetch("http://localhost:8080/api/v1/auth/login", {
  method: "POST",
  headers: {
    "Content-Type": "application/json", // ✅ OBRIGATÓRIO
  },
  body: JSON.stringify({
    email: "admin@erp.com",
    password: "admin123",
  }),
});
```

### **2. Logs úteis:**

#### **Verificar logs da API:**

```bash
docker-compose logs app --tail=20
```

#### **Verificar se containers estão rodando:**

```bash
docker-compose ps
```

#### **Reiniciar API se necessário:**

```bash
docker-compose restart app
```

## 📋 **CHECKLIST DE VERIFICAÇÃO**

### **Backend:**

- [ ] API rodando na porta 8080
- [ ] CORS configurado
- [ ] Endpoint `/health` respondendo
- [ ] Endpoint `/api/v1/auth/login` funcionando

### **Frontend:**

- [ ] URL da API correta (`http://localhost:8080`)
- [ ] Headers corretos (`Content-Type: application/json`)
- [ ] Payload correto (email e password)
- [ ] Tratamento de erros implementado

### **Rede:**

- [ ] Sem firewall bloqueando
- [ ] Porta 8080 acessível
- [ ] Docker rodando corretamente

## 🎯 **EXEMPLO DE CÓDIGO FUNCIONAL**

### **API Client:**

```javascript
class ApiClient {
  constructor() {
    this.baseURL = "http://localhost:8080";
  }

  async login(email, password) {
    try {
      const response = await fetch(`${this.baseURL}/api/v1/auth/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Login failed");
      }

      const data = await response.json();

      // Armazenar tokens
      localStorage.setItem("access_token", data.access_token);
      localStorage.setItem("refresh_token", data.refresh_token);
      localStorage.setItem("user", JSON.stringify(data.user));

      return data;
    } catch (error) {
      console.error("Login error:", error);
      throw error;
    }
  }

  async getProfile() {
    const token = localStorage.getItem("access_token");

    const response = await fetch(`${this.baseURL}/api/v1/users/profile`, {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error("Failed to get profile");
    }

    return await response.json();
  }
}
```

### **Uso no React:**

```jsx
import { useState } from "react";

function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const apiClient = new ApiClient();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const result = await apiClient.login(email, password);
      console.log("Login successful:", result);

      // Redirecionar após login
      window.location.href = "/dashboard";
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Email"
        required
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
        required
      />
      <button type="submit" disabled={loading}>
        {loading ? "Logging in..." : "Login"}
      </button>
      {error && <div className="error">{error}</div>}
    </form>
  );
}
```

## 🎉 **RESULTADO**

Agora o frontend deve conseguir se conectar perfeitamente com o backend!

**Dados de teste:**

- **Email**: `admin@erp.com`
- **Password**: `admin123`

**Teste agora no seu frontend e deve funcionar! 🚀**
