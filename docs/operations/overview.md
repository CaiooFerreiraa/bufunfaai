# Operacao e Ambientes

## Topologia recomendada

```text
Mobile Expo
  -> API publica em Vercel
      -> Neon Postgres
      -> Redis gerenciado
      -> Sentry
      -> Datadog
      -> Object storage para exports
      -> Worker e cron fora da Vercel
```

## Escopo da Vercel

- API HTTP publica
- callbacks sincronicos
- endpoints do app mobile
- WAF e rate limiting de borda
- deploy continuo por ambiente

## Fora da Vercel

- Redis gerenciado
- sync assíncrono de Open Finance
- reconciliacao
- notificacoes pesadas
- geracao de exports demorados

## Ambientes

- `development`: uso local e preview tecnico
- `staging`: homologacao funcional com segredos separados
- `production`: ambiente isolado com observabilidade completa

## Arquivos desta fase

- [Matriz de Secrets](./secrets-matrix.md)
- [Mapa de Ambientes](./environments.md)
- [Worker e Cron](./worker-and-cron.md)
- [Release Mobile](./mobile-release.md)
- [Runbook de Incidente](./runbooks/incident-response.md)
- [Runbook de Backup e Restore](./runbooks/backup-and-restore.md)
