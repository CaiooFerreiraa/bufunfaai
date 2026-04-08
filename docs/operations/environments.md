# Mapa de Ambientes

## Development

- projeto Vercel opcional ou execucao local
- branch Neon isolada
- Redis de desenvolvimento
- Open Finance sandbox e provider `mock`
- telemetria com amostragem baixa

## Staging

- projeto Vercel dedicado
- branch Neon de staging
- Redis de staging
- segredos proprios
- validacao E2E de auth, consentimento, sync e export

## Production

- projeto Vercel dedicado
- branch raiz Neon de producao
- Redis de producao
- WAF mais restritivo
- Sentry e Datadog completos
- runbooks e on-call ativos

## Tags obrigatorias

- `env`
- `service`
- `version`
- `region` quando aplicavel
