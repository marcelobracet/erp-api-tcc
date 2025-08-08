# 🚀 Frontend Integration Summary - ERP API

## 📋 **Arquivos Criados**

### **1. Documentação Completa**

- ✅ **`API_DOCUMENTATION.md`** - Documentação completa da API
- ✅ **`frontend-integration-example.js`** - Exemplo prático de integração
- ✅ **`frontend-env-example.env`** - Configuração de ambiente
- ✅ **`FRONTEND_INTEGRATION_SUMMARY.md`** - Este resumo

## 🔧 **Configuração Rápida**

### **1. Copiar arquivos para o frontend:**

```bash
# Copiar para seu projeto frontend
cp API_DOCUMENTATION.md /path/to/your/frontend/
cp frontend-integration-example.js /path/to/your/frontend/src/services/
cp frontend-env-example.env /path/to/your/frontend/.env
```

### **2. Configurar ambiente:**

```bash
# No seu projeto frontend
cp frontend-env-example.env .env
```

### **3. Instalar dependências (se necessário):**

```bash
npm install react-router-dom axios
```

## 📡 **Endpoints Principais**

### **🔐 Autenticação**

```
POST /api/v1/auth/login
POST /api/v1/auth/refresh
```

### **👥 Usuários**

```
POST /api/v1/users/register
GET /api/v1/users/profile
GET /api/v1/users/:id
PUT /api/v1/users/:id
DELETE /api/v1/users/:id
GET /api/v1/users
GET /api/v1/users/count
```

### **🏥 Health Check**

```
GET /health
```

## 🚀 **Exemplo de Uso Rápido**

### **1. Login:**

```javascript
import apiClient from "./services/frontend-integration-example.js";

// Login
const login = async (email, password) => {
  try {
    const result = await apiClient.login(email, password);
    console.log("Login successful:", result.user);
    return result;
  } catch (error) {
    console.error("Login failed:", error.message);
    throw error;
  }
};
```

### **2. Verificar autenticação:**

```javascript
// Verificar se está logado
if (apiClient.isAuthenticated()) {
  console.log("User is logged in");
} else {
  console.log("User needs to login");
}

// Verificar role
if (apiClient.hasRole("admin")) {
  console.log("User is admin");
}
```

### **3. Fazer requisições autenticadas:**

```javascript
// Obter perfil do usuário
const profile = await apiClient.getUserProfile();

// Listar usuários (admin)
const users = await apiClient.listUsers(10, 0);

// Atualizar usuário
const updatedUser = await apiClient.updateUser(userId, userData);
```

## 🔐 **Configurações de Segurança**

### **1. Tokens JWT:**

- **Access Token**: 24 horas
- **Refresh Token**: 7 dias
- **Armazenamento**: localStorage (dev) / httpOnly cookies (prod)

### **2. Roles:**

- **admin**: Acesso total
- **manager**: Acesso limitado
- **user**: Acesso básico

### **3. Headers obrigatórios:**

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

## 🌐 **URLs por Ambiente**

### **Development:**

```
API: http://localhost:8080
Frontend: http://localhost:3000
```

### **QA:**

```
API: https://qa-api.seudominio.com
Frontend: https://qa.seudominio.com
```

### **Production:**

```
API: https://api.seudominio.com
Frontend: https://app.seudominio.com
```

## 📊 **Dados de Teste**

### **Usuário Admin:**

```json
{
  "email": "admin@erp.com",
  "password": "admin123",
  "name": "Admin User",
  "role": "admin"
}
```

### **Usuário Normal:**

```json
{
  "email": "user@erp.com",
  "password": "user123",
  "name": "Normal User",
  "role": "user"
}
```

## 🔧 **Troubleshooting**

### **1. Erro 401 (Unauthorized):**

```javascript
// Token expirou, será renovado automaticamente
// Se falhar, usuário será redirecionado para login
```

### **2. Erro 403 (Forbidden):**

```javascript
// Usuário não tem permissão
// Verificar role do usuário
if (!apiClient.hasRole("admin")) {
  showError("Access denied");
}
```

### **3. Erro de CORS:**

```javascript
// Verificar se a URL da API está correta
// Verificar se o backend está rodando
```

### **4. Token não encontrado:**

```javascript
// Verificar se o login foi feito corretamente
// Verificar localStorage
console.log(localStorage.getItem("erp_access_token"));
```

## 📝 **Checklist de Implementação**

### **Frontend:**

- [ ] Copiar arquivos de integração
- [ ] Configurar variáveis de ambiente
- [ ] Implementar tela de login
- [ ] Implementar proteção de rotas
- [ ] Implementar logout
- [ ] Testar todos os endpoints
- [ ] Configurar tratamento de erros
- [ ] Implementar loading states

### **Backend:**

- [ ] Verificar se API está rodando
- [ ] Testar endpoints com curl
- [ ] Configurar CORS para frontend
- [ ] Verificar logs de erro

## 🧪 **Testes Rápidos**

### **1. Testar API:**

```bash
# Health check
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@erp.com", "password": "admin123"}'

# Registrar usuário
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@erp.com", "password": "test123", "name": "Test User", "role": "user"}'
```

### **2. Testar Frontend:**

```javascript
// No console do navegador
const apiClient = new ErpApiClient();

// Testar login
apiClient
  .login("admin@erp.com", "admin123")
  .then((result) => console.log("Login success:", result))
  .catch((error) => console.error("Login failed:", error));

// Testar perfil
apiClient
  .getUserProfile()
  .then((profile) => console.log("Profile:", profile))
  .catch((error) => console.error("Profile failed:", error));
```

## 📞 **Suporte**

### **Problemas Comuns:**

1. **API não responde**: Verificar se Docker está rodando
2. **CORS error**: Verificar configuração de CORS no backend
3. **Token inválido**: Fazer logout e login novamente
4. **Permissão negada**: Verificar role do usuário

### **Logs Úteis:**

```bash
# Ver logs da API
docker-compose logs app

# Ver logs do banco
docker-compose logs postgres

# Verificar status dos containers
docker-compose ps
```

---

## 🎉 **PRONTO PARA USAR!**

Agora você tem:

- ✅ **Documentação completa** da API
- ✅ **Exemplo prático** de integração
- ✅ **Configuração de ambiente** pronta
- ✅ **Dados de teste** funcionando
- ✅ **Troubleshooting** detalhado

**Basta copiar os arquivos para seu projeto frontend e começar a usar! 🚀**
