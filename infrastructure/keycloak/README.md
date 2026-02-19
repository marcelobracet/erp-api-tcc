# Keycloak realm import (dev)

Qualquer arquivo `.json` colocado nesta pasta será montado em `/opt/keycloak/data/import` dentro do container e importado automaticamente na inicialização (`--import-realm`).

Sugestão:
- Exporte seu realm do Keycloak (Realm settings → Action → Partial export / Export) e salve aqui como `realm-erp.json`.

Depois reinicie o container:
- `docker compose restart keycloak`
