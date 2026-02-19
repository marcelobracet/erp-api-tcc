# Keycloak (Realm único) - Guia rápido

Este projeto suporta autenticação local (JWT) e pode ser configurado para validar tokens do **Keycloak** (OIDC/JWKS).

## 1) Conceito (multi-tenant)
- Realm único
- Cada usuário pertence a 1 tenant (por enquanto)
- O `tenant_id` é um **claim obrigatório** no access token

## 2) Keycloak - configuração mínima
1. Crie um Realm (ex: `erp`).
2. Crie um Client para o frontend (ex: `erp-frontend`) com PKCE habilitado.
3. Crie roles (ex: `admin`, `manager`, `sales`, `finance`, `production`, `viewer`).
4. Crie um **User Attribute** `tenant_id` (você pode setar direto no usuário).
5. Crie um **Protocol Mapper** para incluir `tenant_id` no token:
   - Mapper Type: *User Attribute*
   - User Attribute: `tenant_id`
   - Token Claim Name: `tenant_id`
   - Add to access token: ON

Opcional (roles por client): defina `KEYCLOAK_CLIENT_ID` e crie Client Roles no client correspondente.

## 3) API - variáveis de ambiente
- `AUTH_PROVIDER=keycloak`
- `KEYCLOAK_ISSUER=http://localhost:8081/realms/erp`
- `KEYCLOAK_AUDIENCE=erp-frontend` (ou outro audience aceito)
- `KEYCLOAK_CLIENT_ID=erp-frontend` (se for ler client roles em `resource_access`)

### Dica importante (Docker x Issuer)
O middleware valida o token comparando o claim `iss` com `KEYCLOAK_ISSUER`.

- Se você **rodar a API no host** (ex: `go run ./cmd/api`), use `KEYCLOAK_ISSUER=http://localhost:8081/realms/erp`.
- Se você **rodar a API no Docker**, normalmente você quer manter `KEYCLOAK_ISSUER` como `http://localhost:8081/realms/erp` (pra bater com o token emitido no navegador), mas o container não consegue resolver `localhost` do host.

Para esse cenário, você pode configurar também:
- `KEYCLOAK_JWKS_URL=http://keycloak:8080/realms/erp/protocol/openid-connect/certs`

Assim o container busca as chaves no endereço interno do compose (`keycloak:8080`), mas continua validando `iss` como `localhost:8081`.

## 5) Subindo o Keycloak via Docker Compose
O `docker-compose.yml` já tem um serviço `keycloak` para dev.

1. Suba os containers:
   - `docker compose up -d`
2. Acesse o admin:
   - `http://localhost:8081/admin`
   - usuário/senha (dev): `admin` / `admin`
3. Crie o realm `erp` e o client `erp-frontend` (ou importe um realm exportado colocando o JSON em `infrastructure/keycloak/`).

## 4) Migração gradual sugerida
- Curto prazo: manter `/api/v1/auth/login` (JWT local) para dev e smoke tests.
- Produção: frontend faz login via Keycloak e chama a API com `Authorization: Bearer <access_token>`.
- Longo prazo: criação/gestão de usuários via Admin API do Keycloak (convites), mantendo a tabela `users` como espelho mínimo (tenant_id + dados do app).
