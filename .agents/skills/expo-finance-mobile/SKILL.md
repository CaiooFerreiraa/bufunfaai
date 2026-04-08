---
name: expo-finance-mobile
description: Desenvolvimento e evolucao de apps mobile de gerenciamento financeiro com React Native, Expo SDK, Expo Router, Clerk, NativeWind v4, Reanimated, Zustand, TanStack Query, AsyncStorage, Zod, Hono.js, Prisma e Neon em monorepo com Turborepo e pnpm. Use quando Codex precisar criar telas, fluxos autenticados, hooks, use-cases, integracoes mobile ou backend compartilhado, sempre com a versao estavel mais recente das dependencias principais.
---

# Expo Finance Mobile

## Fluxo

1. Ler `.agent/context.md` antes de qualquer tarefa.
2. Se `.agent/context.md` nao existir, criar com o template padrao fornecido pelo projeto.
3. Antes de iniciar projeto novo ou upgrade relevante, verificar as versoes estaveis atuais do ecossistema Expo.
4. Manter a separacao entre `app`, `components`, `domain`, `infrastructure`, `hooks`, `lib`, `server` e `packages/shared`.
5. Implementar logica em hooks e use-cases; manter componentes focados em renderizacao.
6. Registrar decisoes relevantes de arquitetura, banco ou visual em `.agent/context.md`.
7. Validar acessibilidade, estados, touch targets, safe areas e performance antes de encerrar.

## Regra de Versao

- Sempre usar a versao estavel mais recente das dependencias principais.
- Nunca fixar versoes antigas sem justificativa explicita do usuario.
- Revalidar as versoes antes de iniciar projeto ou upgrade importante.

Referencia verificada em 2026-04-08:

- `expo`: `55.0.12`
- `expo-router`: `55.0.11`
- `@clerk/clerk-expo`: `2.19.31`
- `react-native-reanimated`: `4.3.0`
- `nativewind`: `4.2.3`
- `hono`: `4.12.12`
- `prisma`: `7.7.0`

## Estrutura Base

```text
apps/
  mobile/
    app/
      (auth)/
      (app)/
      _layout.tsx
    components/
    domain/
      entities/
      repositories/
      use-cases/
    infrastructure/
      repositories/
    hooks/
    lib/
  server/
    src/
      routes/
      domain/
      infrastructure/
        database/prisma.ts
        repositories/
    prisma/schema.prisma
packages/
  shared/
    types/
    utils/
    validators/
```

## Regras de Arquitetura

- `domain` nao conhece Expo, React Native, Prisma, Clerk ou framework HTTP.
- `use-cases` orquestram contratos; nunca acessam infraestrutura diretamente.
- Componentes renderizam; a logica fica em hooks, actions locais e use-cases.
- Nunca chamar API diretamente em componentes.
- Repositorios concretos vivem em `infrastructure`.
- Tipos, utilitarios e validadores compartilhados ficam em `packages/shared`.

## Auth com Clerk

- Separar rotas publicas e protegidas com `(auth)/` e `(app)/`.
- Verificar sessao no `_layout.tsx` raiz com `useAuth()`.
- Manter autorizacao nos use-cases, nao na UI.
- Configurar Clerk em `lib/clerk.ts`.
- Nunca expor tokens em logs.
- Nao persistir tokens sensiveis em `AsyncStorage` sem criptografia apropriada.

## Estado e Dados

- Usar `Zustand` para estado local/global de interface.
- Usar `TanStack Query` para cache, sincronizacao e estados de rede.
- Usar `AsyncStorage` para persistencia nao sensivel do app.
- Validar payloads e parametros com `Zod`.
- Nunca deixar componente assumir regras de dados ou acesso direto a fetcher sem mediacao.

## Banco com Prisma e Neon

- Centralizar schema em `apps/server/prisma/schema.prisma`.
- Criar migrations com `prisma migrate dev`.
- Instanciar `PrismaClient` em singleton em `infrastructure/database/prisma.ts`.
- Nao importar `PrismaClient` diretamente em use-cases.
- Mapear modelos Prisma explicitamente para entidades de dominio.
- Em producao, usar Neon com `?pgbouncer=true&connect_timeout=10&sslmode=require`.
- Validar `DATABASE_URL` com Zod no boot.

## UI e UX Mobile

- Garantir touch targets minimos de `44x44`.
- Usar `SafeAreaView` ou `useSafeAreaInsets()` em toda tela.
- Usar `KeyboardAvoidingView` em telas com inputs.
- Para listas, usar `FlatList` ou `FlashList`; nunca `ScrollView` com `map` para listas reais.
- Adicionar `accessibilityLabel` em icones e botoes sem texto visivel.
- Garantir contraste minimo WCAG AA.
- Usar tokens de design; evitar valores magicos como cores soltas e espacamentos arbitrarios.

Estados obrigatorios:

- `default`
- `pressed`
- `disabled`
- `loading`
- `error`
- `empty`
- `success`

Haptics obrigatorios:

- Acao primaria: `impactAsync(Medium)`
- Sucesso: `notificationAsync(Success)`
- Erro: `notificationAsync(Error)`
- Selecao: `selectionAsync()`
- Acao destrutiva: `notificationAsync(Warning)`

Animacao:

- Usar `React Native Reanimated` para animacoes com 60fps na thread nativa.
- Manter transicoes na faixa de `150ms` a `300ms`.

## TypeScript e Codigo

- Tipificar explicitamente variaveis, parametros e retornos quando a declaracao nao for obvia.
- Evitar `any`.
- Usar `unknown` apenas quando houver narrowing posterior.
- Manter componentes em `PascalCase` e arquivos em `kebab-case`.
- Usar `cn()` do NativeWind para merge de classes.
- Aplicar SRP, `early return` e nomes autoexplicativos.

## Ferramentas e Terminal

- Em monorepo, usar `pnpm` como package manager.
- Para dependencias Expo e nativas, preferir `npx expo install` para preservar compatibilidade com o SDK.
- Nao usar `&&` nem `||`; usar `;` quando precisar encadear comandos.
- Nao usar `sudo` sem avisar.
- Nao usar remocao destrutiva sem confirmacao explicita.

## .gitignore

Garantir pelo menos:

```gitignore
node_modules/
dist/
build/
.expo/
.env
.env.*
!.env.example
.agent/
*.jks
*.p8
*.p12
*.key
*.mobileprovision
web-build/
.DS_Store
*.tsbuildinfo
```

Nunca versionar `.env`, certificados ou `.agent/`.

## Implementacao Padrao

1. Ler contexto do projeto.
2. Verificar versoes estaveis relevantes.
3. Separar o que e UI, hook, caso de uso, repositorio e rota de servidor.
4. Validar entradas com Zod.
5. Integrar auth, estado e networking nas camadas corretas.
6. Aplicar regras de mobile UX, acessibilidade e haptics.
7. Registrar decisoes relevantes em `.agent/context.md`.

## Checklist Final

- `.agent/context.md` foi lido ou criado.
- As dependencias principais respeitam o latest stable.
- `domain` continua puro.
- Componentes nao chamam API diretamente.
- Auth e autorizacao estao nas camadas corretas.
- Safe area, teclado, acessibilidade e touch targets foram tratados.
- Listas usam `FlatList` ou `FlashList`.
- Haptics e estados essenciais foram mapeados.
- Mudancas relevantes foram registradas no contexto.
