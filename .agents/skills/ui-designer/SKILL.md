---
name: ui-designer
description: Direcao de design e implementacao de interfaces modernas, minimalistas e memoraveis, com identidade visual forte, hierarquia clara, composicao intencional e movimento com proposito. Use quando Codex precisar criar, redesenhar ou refinar telas, landing pages, dashboards, formularios, fluxos, design systems ou componentes visuais, especialmente quando o resultado nao puder ter cara de template generico ou de UI gerada por IA.
---

# UI Designer

## Fluxo

1. Ler `.agent/context.md` antes de qualquer tarefa para respeitar a identidade visual, os tokens e as decisoes de design do projeto.
2. Extrair do contexto as restricoes de tipografia, paleta, espacamento, componentes base, acessibilidade e motion.
3. Responder internamente antes de implementar: qual problema a interface resolve, quem a usa, qual emocao deve causar, qual diferencial visual deve torna-la memoravel e qual direcao estetica deve guiar tudo.
4. Escolher uma direcao estetica principal e executar com consistencia total em tipografia, cor, espacamento, composicao, imagem, icones e movimento.
5. Implementar ou orientar com base nessa direcao, sem cair em layout generico ou em repeticao de padroes tipicos de UI gerada por IA.
6. Validar contraste, foco, estados, responsividade e clareza da hierarquia antes de encerrar.

## Escolher a Direcao

Escolher uma destas direcoes e comprometer a interface inteira com ela:

- `minimalista refinado`
- `editorial/magazine`
- `luxo sobrio`
- `brutalista`
- `organico/natural`
- `futurista frio`
- `retro-moderno`
- `industrial`
- `art deco geometrico`
- `soft pastel`

Misturar direcoes so quando existir justificativa real do produto. Se o projeto ja tiver uma linguagem visual definida, evoluir essa linguagem em vez de reinventa-la do zero.

## Evitar Cara de IA

Recusar padroes genericos sempre que eles surgirem como primeira ideia:

- Nao usar cards brancos com borda arredondada e sombra generica como solucao padrao.
- Nao usar gradiente roxo/azul em fundo branco ou escuro por inercia.
- Nao montar hero centralizado com titulo, subtitulo e botao empilhados sem tensao visual.
- Nao repetir grid de tres features simetrico com icone, titulo e texto sem identidade.
- Nao montar navbar padrao de logo, links e CTA sem linguagem propria.
- Nao aplicar o mesmo espacamento em tudo; criar ritmo.
- Nao depender da paleta tipica de UI generica com `#6366f1`, `#3b82f6` e cinzas comuns.
- Nao usar `fade-in` identico em todos os elementos ao mesmo tempo.
- Nao repetir a mesma estrutura em todas as telas quando o produto pede variacao e hierarquia.

Se o resultado parecer template gratuito, refazer.

## Tipografia

- Nunca escolher `Inter`, `Roboto`, `Arial` ou stack de sistema como fonte principal, exceto quando o contexto do projeto exigir explicitamente outra coisa.
- Parear uma fonte de display marcante com uma fonte de corpo refinada.
- Construir hierarquia com tamanho, peso, largura, tracking e espacamento antes de depender de cor.
- Preferir familias com personalidade, como `Playfair Display`, `DM Serif Display`, `Syne`, `Fraunces`, `Instrument Serif`, `Cabinet Grotesk`, `Unbounded`, `Cormorant` e `Bebas Neue`.
- Respeitar o contexto do projeto: se a base atual ainda usa fontes de sistema, propor a evolucao sem quebrar a experiencia existente.

## Cor e Tema

- Definir no maximo tres cores base e um acento principal.
- Construir uma paleta intencional; evitar cores timidas que nao comunicam nada.
- Em dark mode, usar profundidade real com camadas como `#0a0a0a`, `#111` e `#1a1a1a`; nunca depender de preto puro sem contexto.
- Em light mode, evitar `#ffffff` puro como fundo dominante sem temperatura ou atmosfera.
- Declarar variaveis CSS ou tokens desde o inicio para garantir consistencia.
- Respeitar os tokens do projeto e evitar valores magicos fora do sistema.

## Composicao e Layout

- Questionar o layout obvio antes de implementa-lo.
- Usar assimetria, sobreposicao, cortes, diagonais ou blocos de densidade diferente quando isso fortalecer a narrativa visual.
- Escolher entre espaco negativo generoso ou densidade controlada; evitar o meio-termo acidental.
- Quebrar a grid com intencao quando isso gerar foco ou identidade.
- Manter alinhamento preciso; desalinhamento nao intencional destroi credibilidade.
- Guiar o olho do usuario por uma trilha clara de leitura e acao.
- Em produtos financeiros, priorizar legibilidade, contraste e leitura rapida de dados densos.

## Movimento e Feedback

- Animar com proposito, nunca por decoracao vazia.
- Preferir reveals escalonados e transicoes coordenadas em vez de animar tudo igual.
- Criar hover, press e focus states com mudanca perceptivel de forma, deslocamento, escala, opacidade ou profundidade.
- Dar feedback visual imediato para cada acao relevante.
- Usar toast para notificacoes do sistema com o usuario.
- Aplicar `cursor: pointer` sempre que o componente representar uma acao clicavel.
- Manter transicoes consistentes em toda a interface.

## Fundos e Atmosfera

- Evitar fundos solidos sem intencao.
- Usar gradient meshes, ruido sutil, texturas leves, padroes geometricos ou camadas com blur quando isso reforcar a atmosfera.
- Escolher entre sombras dramaticas ou ausencia quase total de sombras; evitar a sombra generica do meio-termo.
- Tratar divisores e bordas como parte da linguagem visual, nao como `border-gray-200` por padrao.

## Componentes

- Dar personalidade aos botoes por forma, peso, contraste, ritmo de hover e presenca.
- Tratar cards como unidades visuais intencionais, nao como caixas neutras por habito.
- Fazer inputs comunicarem o sistema de design pelo foco, pela tipografia e pelos estados.
- Usar apenas icones de `Lucide React`.
- Manter tamanho e `stroke` consistentes para os icones em toda a interface.
- Tratar badges, tags e chips como parte da linguagem visual.
- Aplicar tratamento intencional a imagens com recortes, blend, filtros ou composicao contextual quando necessario.

## Acessibilidade Visual

- Garantir contraste minimo WCAG AA em todos os textos.
- Manter foco visivel e coerente com o sistema de design.
- Verificar se a hierarquia continua funcionando em escala de cinza.
- Garantir alvos interativos com pelo menos `44px`.
- Nao depender exclusivamente de cor para comunicar estado.

## Implementar no Projeto

- Comecar pelo contexto real do produto, nao por referencias aleatorias.
- Preservar padroes existentes quando trabalhar dentro de um design system ja estabelecido.
- Evoluir tokens e componentes sem quebrar consistencia visual entre telas.
- Em React, React Native ou web, manter coerencia entre estados, motion e espacamento.
- Se a tarefa for exploratoria, entregar primeiro uma direcao curta com conceito, paleta, tipografia, composicao e interacoes.
- Se a tarefa pedir implementacao, codar direto com intencao visual clara e sem boilerplate estetico.

## Checklist Final

- A interface resolve o problema certo para o usuario certo.
- Existe uma emocao dominante perceptivel.
- A direcao estetica escolhida aparece de forma consistente.
- O layout foge do padrao generico.
- Tipografia, cor e espacamento constroem hierarquia clara.
- Estados interativos, toasts e foco estao bem resolvidos.
- Os icones usam somente `Lucide React`.
- Elementos clicaveis usam `cursor: pointer`.
- A interface funciona em desktop e mobile.
- O resultado parece produto real, nao demo de template.
