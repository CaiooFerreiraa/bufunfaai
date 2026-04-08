---
name: nextjs-bun
description: Desenvolvimento e evolucao de aplicacoes web modernas com Next.js App Router, Bun como runtime padrao, TypeScript strict, shadcn/ui, Tailwind CSS, Prisma, Neon e autenticacao com NextAuth.js. Use quando Codex precisar criar, refatorar ou expandir um projeto Next.js moderno, especialmente quando a regra for usar a versao estavel mais recente do Next.js, executar tudo com Bun e manter arquitetura limpa entre app, domain, infrastructure, lib e actions.
---

# Next.js + Bun

## Fluxo

1. Como primeira acao de cada tarefa, tentar ler `.agent/context.md`.
2. Se `.agent/context.md` nao existir, criar `.agent/context.md` com o template padrao de contexto.
3. Antes de iniciar um projeto novo ou atualizar dependencias principais, verificar a versao estavel atual do `next` e do `next-auth`.
4. Usar `Bun` como runtime padrao para instalar dependencias, executar scripts e rodar arquivos TypeScript.
5. Se a tarefa envolver UI, ler `.agent/design.md` antes de qualquer alteracao visual.
6. Se `.gitignore` existir, conferir se ele cobre os itens obrigatorios; se nao existir, criar antes do primeiro commit.
7. Sempre que a alteracao for relevante, sugerir atualizacao de `.agent/context.md`.
8. Sempre que houver mudanca de banco ou decisao visual relevante, registrar isso em `.agent/context.md`.

## Regra de Versao

- Sempre usar a versao estavel mais recente do `Next.js`.
- Sempre usar `Bun` como runtime do projeto.
- Nunca fixar uma versao antiga de `Next.js` sem justificativa explicita do usuario.
- Antes de iniciar um projeto ou upgrade relevante, verificar a versao atual com `bunx npm view next version`.
- Se o projeto usar `NextAuth.js`, verificar a versao atual com `bunx npm view next-auth version`.
- Se houver incompatibilidade entre a versao mais nova e o ecossistema do projeto, explicar o trade-off antes de fugir do latest.

Referencia verificada em 2026-04-08:

- `next`: `16.2.2`
- `next-auth`: `4.24.13`

## Stack Base

- Framework: `Next.js` com `App Router`
- Runtime: `Bun`
- Linguagem: `TypeScript` em strict mode
- UI: `shadcn/ui` + `Tailwind CSS`
- Auth: `NextAuth.js`
- ORM: `Prisma`
- Banco: `Neon`

## Bun como Runtime

- Nunca usar `node`, `npm`, `npx` ou `ts-node` diretamente no projeto.
- Usar `bun install` no lugar de `npm install`.
- Usar `bun add` no lugar de `npm install <pacote>`.
- Usar `bunx` no lugar de `npx`.
- Usar `bun run` no lugar de `npm run`.
- Usar `bun <arquivo>.ts` para executar TypeScript diretamente quando fizer sentido.
- Nao versionar `package-lock.json` nem `yarn.lock`.
- Versionar `bun.lockb`; ele e o lockfile oficial do projeto.

## Estrutura Base

```text
src/
  app/
  components/
  domain/
    entities/
    repositories/
    use-cases/
  infrastructure/
    database/
    repositories/
  lib/
  actions/
```

- `app/`: rotas, composicao de tela e orquestracao do App Router
- `components/`: componentes reutilizaveis e UI
- `domain/`: entidades, contratos e regras puras
- `infrastructure/`: Prisma, auth adapters e integracoes externas
- `lib/`: helpers e configuracoes compartilhadas
- `actions/`: Server Actions como ponte entre UI e use-cases

## Regras de Arquitetura

- Manter `domain` sem dependencia de `Next.js`, `Prisma`, `NextAuth.js` ou infraestrutura externa.
- Fazer `use-cases` dependerem de contratos definidos no dominio.
- Implementar repositorios concretos em `infrastructure`.
- Fazer `Server Actions` chamarem use-cases; nunca acessar banco diretamente nelas.
- Manter logica de autorizacao nos use-cases, nao espalhada em componentes de UI.
- Preferir `Server Components`; usar `Client Components` apenas quando houver necessidade real de estado, efeitos ou APIs do navegador.

## Codigo e TypeScript

- Tipificar explicitamente variaveis, parametros e retornos quando a declaracao nao for obvia.
- Evitar `any`.
- Usar `unknown` apenas quando houver narrowing posterior.
- Usar arquivos em `kebab-case` e componentes em `PascalCase`.
- Usar `cn()` do `shadcn/ui` para merge de classes Tailwind.
- Validar entradas com `zod` nas bordas da aplicacao.
- Validar variaveis de ambiente em `env.ts`.

## UI e UX

- Ler `.agent/design.md` antes de qualquer alteracao visual.
- Respeitar design system, tokens, componentes e convencoes visuais ja definidos.
- Nao inventar cores, componentes ou estilos sem antes verificar o que ja existe.
- Usar design minimalista: espacamento generoso, paleta neutra e tipografia clara.
- Mapear estados `default`, `hover`, `focus`, `active`, `disabled`, `loading`, `error`, `empty` e `success`.
- Garantir contraste WCAG AA, foco visivel e labels semanticos.
- Tratar mobile-first como padrao.
- Manter microtransicoes na faixa de `150ms` a `300ms`.

## Cursores

- Declarar explicitamente o cursor em elementos interativos.
- Aplicar `cursor-pointer` em botoes, links, icones clicaveis, tags clicaveis e qualquer elemento com `onClick`.
- Aplicar `cursor-text` em inputs e textareas.
- Aplicar `cursor-not-allowed` em elementos desabilitados junto com sinal visual coerente.
- Aplicar `cursor-grab` e `cursor-grabbing` em drag and drop.
- Nao confiar no cursor padrao do navegador como regra geral.

## Auth

- Usar `NextAuth.js` para autenticacao e sessao.
- Configurar auth em `src/lib/auth.ts` ou `auth.ts` na raiz, conforme a convencao da versao usada.
- Proteger rotas com `middleware.ts` quando necessario.
- Configurar adapters de auth na camada `infrastructure`.
- Nunca expor dados sensiveis no client.

## Prisma e Neon

- Centralizar schema em `prisma/schema.prisma`.
- Criar migrations com `bunx prisma migrate dev`.
- Instanciar Prisma uma unica vez em `src/infrastructure/database/prisma.ts`.
- Nao importar `PrismaClient` diretamente nos use-cases.
- Fazer mapeamento explicito de modelos Prisma para entidades de dominio.
- Usar `DATABASE_URL` validada no boot.
- Preferir URL Neon com `?pgbouncer=true&connect_timeout=10`.
- Exigir `?sslmode=require` em producao.

## .gitignore

Ao criar ou revisar `.gitignore`, garantir estes itens:

```gitignore
# Dependencies
node_modules/
.pnp
.pnp.js

# Build outputs
.next/
out/
dist/
build/

# Environment variables
.env
.env.*
!.env.example

# Agent context
.agent/

# Editors e sistemas operacionais
.DS_Store
Thumbs.db
*.pem
.vscode/
.idea/

# Debug
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Vercel
.vercel

# TypeScript
*.tsbuildinfo
next-env.d.ts
```

Regras obrigatorias:

- Nunca versionar `.env`; apenas `.env.example` sem segredos.
- Sempre ignorar `.agent/`.
- Sempre ignorar `.next/`, `node_modules/`, `dist/` e outros artefatos de build.
- Nunca ignorar `bun.lockb`; ele deve ser versionado.

## Regras de Terminal

- Nao usar `&&`.
- Nao usar `||` como fallback em cadeia.
- Preferir um comando por linha.
- Nao usar `sudo` sem avisar o usuario.
- Nao usar remocao destrutiva sem confirmacao explicita.
- Preferir caminhos absolutos em scripts quando isso reduzir ambiguidade.
- Nao assumir o sistema operacional se isso nao estiver claro no contexto.

## Browser

- Quando for necessario abrir uma URL no ambiente local, usar `Firefox`.
- Nao usar o browser padrao do sistema como atalho.

## Implementacao Padrao

Quando criar ou alterar uma feature:

1. Ler contexto e design relevantes.
2. Confirmar versoes atuais de `next` e, se aplicavel, `next-auth`.
3. Verificar `.gitignore`.
4. Definir a separacao entre UI, actions, domain e infrastructure.
5. Criar ou ajustar componentes, actions e use-cases na camada correta.
6. Validar entradas com `zod`.
7. Integrar `Prisma` e `NextAuth.js` apenas nas camadas apropriadas.
8. Subir um servidor local para visualizacao.
9. Registrar decisoes importantes em `.agent/context.md`.

## Checklist Final

- `.agent/context.md` foi lido ou criado.
- `Next.js` esta alinhado com a versao estavel mais recente.
- `Bun` foi usado como runtime padrao.
- App Router foi mantido.
- Server Actions nao acessam DB diretamente.
- `domain` continua puro.
- Entradas externas foram validadas com `zod`.
- `.gitignore` esta correto e `bun.lockb` nao foi ignorado.
- Em alteracoes visuais, `.agent/design.md` foi lido.
- Existe servidor local para visualizacao.
