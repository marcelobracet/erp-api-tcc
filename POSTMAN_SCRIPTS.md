# Scripts para Dados Fictícios - Postman

## 1. Criar Tenant

**POST** `http://localhost:8080/api/v1/tenants`

```json
{
  "name": "Empresa ABC Ltda",
  "plan": "premium"
}
```

## 2. Registrar Usuário

**POST** `http://localhost:8080/api/v1/users/register`

```json
{
  "tenant_id": "{{tenant_id}}",
  "email": "admin@empresa.com",
  "password": "123456",
  "name": "Administrador",
  "role": "admin"
}
```

## 3. Login

**POST** `http://localhost:8080/api/v1/auth/login`

```json
{
  "email": "admin@empresa.com",
  "password": "123456"
}
```

## 4. Criar Cliente

**POST** `http://localhost:8080/api/v1/clients`

**Headers:**

```
Authorization: Bearer {{access_token}}
```

```json
{
  "tenant_id": "{{tenant_id}}",
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

## 5. Criar Produto

**POST** `http://localhost:8080/api/v1/products`

**Headers:**

```
Authorization: Bearer {{access_token}}
```

```json
{
  "tenant_id": "{{tenant_id}}",
  "name": "Mármore Carrara",
  "description": "Mármore branco de alta qualidade",
  "price": 150.0,
  "stock": 100,
  "sku": "MAR-001",
  "category": "Mármore",
  "image_url": "https://example.com/marmore-carrara.jpg"
}
```

## 6. Criar Orçamento

**POST** `http://localhost:8080/api/v1/quotes`

**Headers:**

```
Authorization: Bearer {{access_token}}
```

```json
{
  "tenant_id": "{{tenant_id}}",
  "client_id": "{{client_id}}",
  "user_id": "{{user_id}}",
  "discount": 10.0,
  "status": "pending",
  "notes": "Orçamento para reforma da cozinha",
  "items": [
    {
      "product_id": "{{product_id}}",
      "quantity": 2,
      "price": 150.0
    }
  ]
}
```

## 7. Atualizar Configurações

**PUT** `http://localhost:8080/api/v1/settings`

**Headers:**

```
Authorization: Bearer {{access_token}}
```

```json
{
  "tenant_id": "{{tenant_id}}",
  "settings": {
    "company_name": "Empresa ABC Ltda",
    "company_email": "contato@empresa.com",
    "company_phone": "(11) 3333-4444",
    "company_address": "Rua Principal, 456",
    "company_city": "São Paulo",
    "company_state": "SP",
    "company_zip": "01234-567",
    "primary_color": "#2196F3",
    "secondary_color": "#FFC107",
    "logo_url": "https://example.com/logo.png"
  }
}
```

## Variáveis do Postman

Configure as seguintes variáveis no Postman:

- `tenant_id`: ID do tenant criado
- `client_id`: ID do cliente criado
- `product_id`: ID do produto criado
- `user_id`: ID do usuário logado
- `access_token`: Token de acesso do login

## Scripts de Teste

### Script para capturar tenant_id

```javascript
if (pm.response.code === 201) {
  const response = pm.response.json();
  pm.environment.set("tenant_id", response.data.id);
}
```

### Script para capturar access_token

```javascript
if (pm.response.code === 200) {
  const response = pm.response.json();
  pm.environment.set("access_token", response.data.access_token);
  pm.environment.set("user_id", response.data.user.id);
}
```

### Script para capturar client_id

```javascript
if (pm.response.code === 201) {
  const response = pm.response.json();
  pm.environment.set("client_id", response.data.id);
}
```

### Script para capturar product_id

```javascript
if (pm.response.code === 201) {
  const response = pm.response.json();
  pm.environment.set("product_id", response.data.id);
}
```

## Ordem de Execução

1. Criar Tenant
2. Registrar Usuário
3. Login
4. Criar Cliente
5. Criar Produto
6. Criar Orçamento
7. Atualizar Configurações

## Dados Adicionais

### Mais Clientes

```json
{
  "tenant_id": "{{tenant_id}}",
  "name": "Maria Santos",
  "email": "maria@email.com",
  "phone": "(11) 88888-8888",
  "document": "987.654.321-00",
  "document_type": "CPF",
  "address": "Av. Paulista, 1000",
  "city": "São Paulo",
  "state": "SP",
  "zip_code": "01310-100"
}
```

### Mais Produtos

```json
{
  "tenant_id": "{{tenant_id}}",
  "name": "Granito Preto",
  "description": "Granito preto absoluto",
  "price": 200.0,
  "stock": 50,
  "sku": "GRA-001",
  "category": "Granito",
  "image_url": "https://example.com/granito-preto.jpg"
}
```

### Mais Orçamentos

```json
{
  "tenant_id": "{{tenant_id}}",
  "client_id": "{{client_id}}",
  "user_id": "{{user_id}}",
  "discount": 5.0,
  "status": "approved",
  "notes": "Orçamento aprovado para banheiro",
  "items": [
    {
      "product_id": "{{product_id}}",
      "quantity": 1,
      "price": 200.0
    }
  ]
}
```
