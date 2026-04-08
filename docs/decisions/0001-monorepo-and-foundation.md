# ADR 0001: Monorepo e Fundacao Tecnica

## Status

Aceita

## Contexto

O projeto precisa evoluir rapido no MVP, mantendo coordenação entre mobile, backend, documentacao e automacao.

## Decisao

Adotar um monorepo com:

- mobile Expo/React Native em `apps/mobile`
- API Go em `apps/api`
- pacotes compartilhados em `packages/*`
- documentacao tecnica em `docs/*`

## Consequencias

### Positivas

- versionamento coordenado entre app e API
- padroes de codigo e CI centralizados
- onboarding mais simples

### Negativas

- pipeline inicial precisa lidar com Node e Go
- fronteiras entre apps devem ser bem mantidas para evitar acoplamento
