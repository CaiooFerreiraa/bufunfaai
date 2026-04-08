# Mobile Release

## Arquivos base

- `apps/mobile/eas.json`
- `apps/mobile/app.json`
- `apps/mobile/.env.example`

## Passos

1. Definir `ios.bundleIdentifier` e `android.package` finais.
2. Configurar projeto EAS com `eas init`.
3. Revisar permissões, push notifications e biometria em device real.
4. Rodar `pnpm --filter mobile eas:build:preview`.
5. Rodar `pnpm --filter mobile eas:build:production`.
6. Submeter com `pnpm --filter mobile eas:submit:production`.

## Checklist minimo

- DSN do Sentry por ambiente
- API URL de staging e production
- deep links configurados
- splash e icones finais
- credenciais Apple e Google configuradas
- push token validado em producao
