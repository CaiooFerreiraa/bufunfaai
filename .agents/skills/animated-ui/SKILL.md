---
name: animated-ui
description: Criacao e implementacao de interfaces React animadas de alto impacto usando React Bits, Aceternity UI e Magic UI, com foco em composicao memoravel, performance, acessibilidade e responsividade. Use quando Codex precisar desenhar ou codar UIs animadas para landing pages, dashboards, portfólios, heroes, social proof, backgrounds interativos ou componentes de destaque que nao possam parecer genericos.
---

# Animated UI

## Objetivo

Criar interfaces animadas memoraveis e fluidas, nao UIs genericas com motion decorativo. Toda escolha deve ter papel claro na hierarquia visual da tela.

## Fluxo

1. Entender o objetivo da tela, secao ou componente.
2. Escolher a biblioteca principal pela matriz de decisao.
3. Instalar ou copiar apenas os componentes necessarios.
4. Customizar via props e composicao; evitar editar o core do componente sem necessidade real.
5. Validar performance, acessibilidade, responsividade e `prefers-reduced-motion`.

## Bibliotecas e Papeis

### React Bits

Usar para `statement pieces`, fundos animados, profundidade, efeitos 3D, particulas e animacoes de texto criativas.

Preferir quando precisar de:

- backgrounds animados
- texto com reveal marcante
- efeitos que facam o usuario parar e olhar
- secoes hero com presenca visual forte

### Aceternity UI

Usar para heroes, secoes de marketing, cards com hover avancado, pricing, testimonials e navegacoes com polish visual alto.

Preferir quando precisar de:

- landing pages SaaS
- heroes impactantes
- cards com tilt, glow ou borda animada
- secoes promocionais completas

### Magic UI

Usar para elementos funcionais com camada visual extra, especialmente em produtos reais que ja usam design system estruturado.

Preferir quando precisar de:

- marquees e animated lists
- contadores e KPIs animados
- dock, toolbar e navegacao criativa
- componentes intermediarios entre utilitario e visual

## Matriz de Decisao

- Hero section impactante: `Aceternity UI`
- Background animado ou profundidade: `React Bits`
- Componente integrado ao design system: `Magic UI`
- Texto com animacao criativa: `React Bits`
- Cards de features com hover avancado: `Aceternity UI`
- Marquee, social proof ou lista animada: `Magic UI`
- Efeito 3D ou particulas: `React Bits`
- Pricing ou testimonials: `Aceternity UI`
- Dock ou navegacao flutuante: `Magic UI`
- Elemento unico de destaque: `React Bits`

## Regras de Composicao

- Usar uma biblioteca por camada visual.
- `React Bits` para fundo e efeito.
- `Aceternity UI` para secao principal.
- `Magic UI` para reforcos funcionais e detalhes.
- Nao misturar dois efeitos pesados no mesmo viewport.
- Manter um orcamento de no maximo `3` ou `4` componentes animados simultaneos por pagina.

## Performance

- Nunca animar propriedades que causem reflow como `width`, `height`, `top` ou `left`.
- Preferir `transform` e `opacity`.
- Usar `will-change: transform` com parcimonia.
- Lazy-load de componentes pesados como `Three.js` com `React.lazy` e `Suspense`.
- Em Framer Motion, usar `layout` apenas quando realmente necessario.

## Acessibilidade

- Todo componente animado precisa de alternativa estatica quando `prefers-reduced-motion` estiver ativo.
- Elementos decorativos animados devem usar `aria-hidden="true"`.
- Nunca depender de animacao para transmitir informacao critica.
- Garantir fallback de hover para touch.

## Responsividade

- Validar comportamento mobile antes de levar um efeito para producao.
- Simplificar ou desabilitar 3D em telas pequenas quando necessario.
- Testar densidade, legibilidade e performance em viewports menores.

## Combinacoes Recomendadas

### Landing Page SaaS

- Fundo: `React Bits`
- Hero: `Aceternity UI`
- Features: `Aceternity UI`
- Social proof: `Magic UI`
- CTA destacado: `React Bits`

### Portfolio

- Fundo: `React Bits`
- Intro: `React Bits`
- Projetos: `Aceternity UI`
- Skills: `Magic UI`
- Contato: `Aceternity UI`

### Dashboard

- Sidebar ou dock: `Magic UI`
- KPIs: `Magic UI`
- Onboarding: `Aceternity UI`
- Empty states: `React Bits`

## Convencoes de Codigo

- Usar TypeScript sempre.
- Extrair animacoes repetidas para `/components/animated/`.
- Nomear variantes semanticamente, como `entering`, `visible` e `exiting`.
- Centralizar duracoes e easings com CSS custom properties ou tokens equivalentes.

## Implementacao Padrao

1. Escolher a biblioteca certa.
2. Copiar ou instalar somente o necessario.
3. Adaptar ao design system existente.
4. Orquestrar motion sem excesso.
5. Testar reducao de movimento, mobile e performance.

## Checklist Final

- A tela tem um elemento de destaque claro.
- As animacoes reforcam a hierarquia, nao a poluem.
- Nao ha dois efeitos pesados competindo no mesmo viewport.
- `prefers-reduced-motion` foi respeitado.
- Hover tem fallback em touch.
- A pagina continua legivel e rapida em mobile.
