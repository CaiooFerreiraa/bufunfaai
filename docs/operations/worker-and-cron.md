# Worker e Cron

## Componentes

- `apps/api/cmd/worker`: executa jobs assíncronos fora do caminho HTTP publico
- `POST /internal/open-finance/reconcile`: endpoint interno protegido por `CRON_SECRET`

## Job atual

- `openfinance-reconcile`: percorre conexoes ativas e executa reconciliacao limitada por `WORKER_BATCH_SIZE`

## Formas de execucao

- worker dedicado: `make -C apps/api worker`
- execucao pontual: `make -C apps/api reconcile`
- cron externo chamando o endpoint interno com `X-Internal-Secret`

## Regras

- nao expor rota `/internal/*` publicamente sem segredo
- manter batches curtos na Vercel
- jobs pesados ou demorados ficam no worker externo
