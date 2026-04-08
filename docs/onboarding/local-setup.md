# Setup Local

## Requisitos

- Node.js 22+
- pnpm 10+
- Go 1.24+
- Docker Desktop

## Mobile

1. Copie `apps/mobile/.env.example` para `apps/mobile/.env`
2. Rode `pnpm install`
3. Rode `pnpm --filter mobile dev`

## API

1. Copie `apps/api/.env.example` para `apps/api/.env`
2. Se estiver usando Neon, ajuste `DATABASE_URL` no `.env` e suba apenas o Redis local
3. Se estiver usando Postgres local, suba dependencias com `docker compose -f apps/api/docker-compose.yml up -d postgres redis`
4. Rode `make -C apps/api dev`

## Validacao

- Mobile deve abrir no Expo
- API deve responder `GET /health`
