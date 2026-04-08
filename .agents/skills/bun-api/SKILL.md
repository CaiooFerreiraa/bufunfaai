---
name: bun-api
description: Desenvolvimento e evolucao de APIs modernas com Bun, TypeScript strict, arquitetura limpa, TDD, Zod, Prisma e Neon. Use quando Codex precisar criar endpoints, organizar camadas, implementar casos de uso, validar entradas, integrar banco PostgreSQL via Prisma, escrever testes com Bun Test ou refatorar uma API para manter isolamento entre domain, application e infrastructure.
---

# Bun API

## Fluxo

1. Ler `.agent/context.md` antes de qualquer tarefa para respeitar decisoes, convencoes e historico do projeto.
2. Confirmar o objetivo da alteracao e em qual camada ela deve entrar antes de escrever codigo.
3. Definir ou ajustar os testes primeiro seguindo `Red -> Green -> Refactor`.
4. Implementar a solucao respeitando a regra de dependencia: camadas internas nao dependem das externas.
5. Validar entradas na borda com Zod, mapear erros para respostas HTTP padronizadas e manter o dominio puro.
6. Sempre que a mudanca for relevante, sugerir atualizacao de `.agent/context.md`.
7. Sempre que houver mudanca de banco ou alteracao visual significativa, registrar a decisao no historico de `.agent/context.md`.
8. Antes do primeiro commit de qualquer projeto, verificar se o `.gitignore` esta correto.

## Stack Obrigatoria

- Runtime: `Bun`
- Linguagem: `TypeScript` em modo strict
- Testes: `bun test`
- Validacao: `Zod`
- ORM: `Prisma`
- Banco: `Neon` com PostgreSQL serverless

## Estrutura Base

Organizar a API desta forma:

```text
src/
  domain/
    entities/
    repositories/
    use-cases/
    errors/
  application/
    dtos/
    services/
  infrastructure/
    repositories/
    database/
    http/
      controllers/
      middlewares/
      routes/
  shared/
```

Organizar testes espelhando `src/`:

```text
tests/
  unit/
    domain/
    application/
  integration/
    http/
```

## Regras de Arquitetura

- Manter `domain` sem dependencias de Bun, Prisma, banco, HTTP ou qualquer framework.
- Definir contratos de repositorio em `src/domain/repositories`.
- Fazer `use-cases` orquestrarem entidades e contratos, nunca detalhes de infraestrutura.
- Usar `application` para DTOs, coordenacao e servicos de aplicacao.
- Implementar adaptadores concretos em `infrastructure`.
- Manter controllers finos: receber, validar, delegar e responder.
- Nao colocar regra de negocio em controller, middleware, rota ou repositorio Prisma.
- Mapear modelos Prisma para entidades de dominio explicitamente.
- Nao importar `PrismaClient` diretamente em use-cases.

## TDD

- Escrever o teste antes da implementacao.
- Criar testes unitarios isolados para cada use-case com mocks dos repositorios.
- Criar testes de integracao para controllers e rotas.
- Nomear testes como `should [comportamento esperado] when [condicao]`.
- Cobrir comportamento, erros esperados e casos de borda de entidades e use-cases.
- Rodar `bun test` como runner padrao.

## Validacao e DTOs

- Validar toda entrada externa na borda com Zod.
- Centralizar schemas nos DTOs de `application`.
- Nao deixar dados nao validados entrarem no dominio.
- Retornar `400` para erro de validacao com mensagem descritiva.
- Validar variaveis de ambiente no boot via `env.ts` com Zod.

## TypeScript

- Tipificar explicitamente variaveis, parametros, callbacks e retornos quando a declaracao nao for obvia.
- Evitar `any`.
- Usar `unknown` apenas quando o tipo realmente for indeterminado e fizer narrowing depois.
- Preferir interfaces para contratos e `type` para composicoes.
- Usar imports absolutos via `@/` quando o projeto estiver configurado assim.
- Usar arquivos em `kebab-case` e classes em `PascalCase`.

## Codigo Limpo

- Manter funcoes pequenas com responsabilidade unica.
- Usar nomes autoexplicativos.
- Preferir composicao a heranca.
- Usar `early return` para reduzir aninhamento.
- Extrair logica complexa para funcoes nomeadas.
- Evitar duplicacao com bom senso.
- Nao lancar strings cruas; usar erros tipados.

## Erros

- Fazer erros de dominio herdarem de `DomainError`.
- Tratar erros inesperados em middleware global.
- Nao expor stack trace em producao.
- Padronizar resposta HTTP como `{ data, error, meta }`.
- Usar status codes semanticos.
- Produzir log estruturado para falhas inesperadas.

## Prisma e Neon

- Centralizar schema em `prisma/schema.prisma`.
- Criar migrations versionadas com `prisma migrate dev`.
- Instanciar o client uma unica vez em `src/infrastructure/database/prisma.ts`.
- Implementar repositorios Prisma em `src/infrastructure/repositories/prisma-*.ts`.
- Usar `DATABASE_URL` validada no boot.
- Preferir URL Neon com `?pgbouncer=true&connect_timeout=10`.
- Em producao, exigir `?sslmode=require`.

## Implementacao Padrao

Quando a tarefa envolver um novo fluxo de negocio:

1. Definir ou ajustar a entidade e os erros de dominio, se necessario.
2. Definir o contrato de repositorio no `domain`.
3. Escrever testes do use-case.
4. Implementar o use-case.
5. Criar DTOs e schemas Zod na `application`.
6. Implementar controller, rota e middleware necessarios.
7. Implementar adaptador Prisma ou outro repositorio concreto em `infrastructure`.
8. Adicionar testes de integracao para a borda HTTP.
9. Rodar testes e fazer refactor sem quebrar o isolamento entre camadas.

## Proibicoes

- Nao misturar regra de negocio com HTTP.
- Nao acessar Prisma direto de `domain` ou `application`.
- Nao usar modelos Prisma como entidades de dominio.
- Nao confiar em payload sem Zod.
- Nao deixar tipos implicitos em trechos importantes.
- Nao ignorar `.agent/context.md`.
- Nao alterar banco sem refletir a decisao no historico de `.agent/context.md`.

## Checklist Final

- `.agent/context.md` foi lido e respeitado.
- A mudanca ficou na camada certa.
- O teste veio antes da implementacao.
- Entradas externas foram validadas com Zod.
- Use-cases dependem de contratos, nao de Prisma ou HTTP.
- Erros estao tipados e mapeados corretamente.
- Respostas HTTP seguem `{ data, error, meta }`.
- Tipos estao explicitos onde precisam estar.
- Se houve mudanca relevante, foi sugerida atualizacao de `.agent/context.md`.
- Se houve mudanca de banco, isso foi registrado no historico de decisoes.
