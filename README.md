# BufunfaAI

Base inicial da Fase 0 para um app financeiro com mobile em Expo/React Native e API em Go.

## Estrutura

```text
apps/
  api/
  mobile/
packages/
docs/
.github/workflows/
```

## Stack

- Mobile: Expo, React Native, Expo Router, TypeScript, Zustand, TanStack Query, Zod
- API: Go, Gin, pgx, sqlc, Redis, PostgreSQL, Docker Compose
- Qualidade: pnpm workspaces, ESLint, Prettier, GitHub Actions

## Comecando

### Requisitos

- Node.js 22+
- pnpm 10+
- Go 1.24+
- Docker + Docker Compose

### Instalar dependencias

```bash
pnpm install
```

### Rodar mobile

```bash
pnpm --filter mobile dev
```

### Rodar API localmente

```bash
docker compose -f apps/api/docker-compose.yml up -d postgres redis
make -C apps/api dev
```

## Comandos principais

```bash
pnpm lint
pnpm typecheck
pnpm test
pnpm --filter mobile dev
make -C apps/api test
```

## Convencoes

- Conventional Commits
- SemVer
- Branches: `main`, `develop`, `feature/*`, `fix/*`, `hotfix/*`, `chore/*`, `docs/*`

## Documentacao

- [Arquitetura](./docs/architecture/overview.md)
- [ADR 0001](./docs/decisions/0001-monorepo-and-foundation.md)
- [Onboarding local](./docs/onboarding/local-setup.md)
- [Operacao e Ambientes](./docs/operations/overview.md)
