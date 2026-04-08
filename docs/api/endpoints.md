# Endpoints Iniciais

## Healthcheck

- `GET /health`

## Auth base

- `POST /v1/auth/demo-login`

## Open Finance

- `GET /v1/open-finance/institutions`
- `GET /v1/open-finance/institutions/:id`
- `POST /v1/open-finance/consents`
- `GET /v1/open-finance/consents/:id`
- `POST /v1/open-finance/consents/:id/authorize`
- `GET /v1/open-finance/callback`
- `POST /v1/open-finance/callback`
- `POST /v1/open-finance/consents/:id/revoke`
- `GET /v1/open-finance/connections`
- `GET /v1/open-finance/connections/:id`
- `DELETE /v1/open-finance/connections/:id`
- `POST /v1/open-finance/connections/:id/sync`
- `GET /v1/open-finance/connections/:id/sync-status`
