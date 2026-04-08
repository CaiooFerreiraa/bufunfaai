# Arquitetura Inicial

## Objetivo

Estabelecer uma base para evolucao segura do produto sem acoplamento desnecessario entre mobile, backend e integracoes futuras.

## Monorepo

- `apps/mobile`: cliente Expo/React Native
- `apps/api`: API em Go
- `packages/*`: configuracoes e tipos compartilhados
- `docs/*`: arquitetura, decisoes, onboarding e contratos

## Mobile

- `src/app`: rotas com Expo Router
- `src/features`: organizacao por capacidade do produto
- `src/services`: integracoes com API, storage seguro, biometria e notificacoes
- `src/components/ui`: base visual reutilizavel

## API

- `cmd/api`: ponto de entrada
- `internal/platform`: infraestrutura transversal
- `internal/modules`: modulos por dominio
- `internal/shared`: dto, erros e utilitarios comuns

## Decisoes principais

- Monorepo leve para o MVP
- Monolito modular no backend
- SQL explicito com `pgx` e `sqlc`
- Mobile orientado a features
- Estado hibrido com TanStack Query e Zustand
