# Matriz de Secrets

## Tier 1

- `DATABASE_URL`
- `REDIS_URL`
- `JWT_ACCESS_SECRET`
- `APP_ENCRYPTION_KEY`
- `OPENAI_API_KEY`
- `OPENFINANCE_CLIENT_ID`
- `OPENFINANCE_CLIENT_SECRET`
- `OPENFINANCE_MTLS_CERT_PATH`
- `OPENFINANCE_MTLS_KEY_PATH`

## Tier 2

- `SENTRY_DSN`
- `DATADOG_API_KEY`
- `EXPORTS_STORAGE_BUCKET`
- `EXPORTS_STORAGE_BASE_URL`

## Regras

- nenhum secret em repositório
- separar valores por ambiente
- usar variaveis sensiveis na Vercel para segredos criticos
- rotacao trimestral para segredos normais
- rotacao imediata em incidente ou exposicao
- certificados Open Finance ficam fora do codigo e fora do bundle do app
