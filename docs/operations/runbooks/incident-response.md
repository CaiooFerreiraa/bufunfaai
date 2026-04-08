# Runbook de Incidente

## Severidades

- `SEV1`: indisponibilidade geral, vazamento potencial, falha sistemica de auth
- `SEV2`: falha relevante em consentimento, callback ou sync
- `SEV3`: degradacao parcial em relatorios ou analytics
- `SEV4`: bug localizado sem risco operacional direto

## Fluxo

1. Detectar por alerta, dashboard ou suporte.
2. Classificar severidade e abrir canal de incidente.
3. Conter impacto com WAF, feature flag ou rollback.
4. Preservar evidencias e logs relevantes.
5. Restaurar servico.
6. Registrar RCA e acao corretiva.

## Contencao rapida

- rotacionar segredo exposto
- bloquear rota no WAF
- reduzir trafego em endpoint critico
- desligar integracao externa por feature flag
- restaurar deploy anterior
