# Runbook de Backup e Restore

## Objetivo

Garantir recuperacao rapida do Neon e preservar continuidade da API.

## Politica minima

- PITR no branch raiz de producao
- snapshot diario
- snapshot semanal com retencao maior
- export logico periodico criptografado fora do banco principal
- teste mensal de restore

## Fluxo de restore

1. Confirmar escopo do incidente e janela de perda aceitavel.
2. Escolher ponto de restauracao no Neon.
3. Restaurar para branch segura.
4. Validar schema, integridade e migrations.
5. Promover rota de leitura ou redirecionar aplicacao conforme plano.
6. Registrar horario real de RPO e RTO.

## Validacoes obrigatorias

- auth funcionando
- consentimentos e conexoes legiveis
- analytics consistente
- migrations atuais aplicadas
- jobs pausados ate reconciliacao
