---
name: abacatepay
description: Integracao com a API da AbacatePay para clientes, cupons, cobrancas, PIX QRCode, saques, loja e webhooks, incluindo autenticacao por Bearer token, Dev Mode e transicao para producao. Use quando Codex precisar implementar servicos, SDKs internos, endpoints de backend, webhooks ou fluxos de pagamento usando AbacatePay sem expor credenciais no frontend.
---

# AbacatePay

## Principios

- Tratar AbacatePay como integracao de backend.
- Nunca expor token no frontend ou em repositorio publico.
- Usar `Authorization: Bearer <token>` em todas as chamadas.
- Priorizar `Dev Mode` para validar integracao, webhooks e simulacoes antes de producao.

## Fluxo

1. Implementar a integracao no backend, nunca direto no client.
2. Centralizar base URL, token, timeouts e tratamento de erro em um cliente HTTP dedicado.
3. Modelar tipos e validacoes para requests e responses.
4. Usar Dev Mode para testar criacao de cobrancas, QR Codes e webhooks.
5. Preparar transicao para producao trocando token, desligando Dev Mode e revalidando eventos.

## Autenticacao

- Metodo: `Bearer Token`
- Header obrigatorio:

```http
Authorization: Bearer {SEU_TOKEN_AQUI}
```

- Nunca logar o token em texto puro.
- Armazenar token apenas em variavel de ambiente ou secret manager.

## Recursos Principais

### Clientes

Usar para criar e listar clientes associados a cobrancas.

Endpoints principais:

- `POST /v1/customer/create`
- `GET /v1/customer/list`

Boas praticas:

- Reaproveitar `customerId` quando o cliente ja existir.
- Validar `name`, `cellphone`, `email` e `taxId` antes do envio.

### Cupons

Usar para descontos percentuais ou fixos.

Endpoints principais:

- `POST /v1/coupon/create`
- `GET /v1/coupon/list`

Cuidados:

- Tratar `discountKind` como enum.
- Documentar claramente se o desconto esta em percentual ou centavos.

### Cobrancas

Usar para gerar URL de pagamento.

Endpoints principais:

- `POST /v1/billing/create`
- `GET /v1/billing/get?id=...`
- `GET /v1/billing/list`

Boas praticas:

- Validar `frequency`, `methods`, `products`, `returnUrl` e `completionUrl`.
- Tratar `customerId` e `customer` como campos mutuamente exclusivos.
- Trabalhar valores monetarios em centavos.

### PIX QRCode

Usar para gerar `brCode` e QR Code em base64 para exibir no app ou sistema proprio.

Endpoints principais:

- `POST /v1/pixQrCode/create`
- `GET /v1/pixQrCode/check?id=...`
- `POST /v1/pixQrCode/simulate-payment?id=...`

Status importantes:

- `PENDING`
- `PAID`
- `EXPIRED`

Boas praticas:

- Exibir tempo de expiracao claramente.
- Fazer polling com criterio ou preferir webhook quando disponivel.
- Usar simulacao apenas em Dev Mode.

### Saques

Usar para transferir saldo da conta para chave PIX.

Endpoints principais:

- `POST /v1/withdraw/create`
- `GET /v1/withdraw/get?id=...`
- `GET /v1/withdraw/list`

Cuidados:

- Validar `amount` em centavos.
- Validar formato da `pixKey`.

### Loja

Usar `GET /v1/store/get` para obter dados da loja vinculada ao token.

## Webhooks

Eventos principais:

- `billing.paid`
- `pix.paid`
- `pix.expired`
- `withdraw.paid`

Boas praticas obrigatorias:

- Validar assinatura do cabecalho.
- Implementar retries idempotentes.
- Persistir evento recebido antes de processar efeitos colaterais quando necessario.
- Retornar HTTP `200` ao receber evento valido.
- Nao assumir que os eventos chegam em ordem perfeita.

## Modelagem e Validacao

- Criar DTOs e schemas para requests e responses.
- Usar enums para status, frequencia, metodos e tipos de desconto.
- Tratar `error` e `data` explicitamente no contrato da API.
- Nunca confiar em campos opcionais sem validacao.

## Erros e Seguranca

- Mapear erros da AbacatePay para erros internos tipados.
- Registrar contexto suficiente para diagnostico sem vazar credenciais.
- Nunca repassar mensagens sensiveis da API bruta direto ao client.
- Implementar timeouts e tratamento de indisponibilidade.

## Transicao para Producao

1. Desativar Dev Mode.
2. Completar cadastro exigido pela plataforma.
3. Trocar token de desenvolvimento por token de producao.
4. Revalidar criacao de cobrancas e recebimento de webhooks.

## Implementacao Padrao

1. Criar cliente HTTP ou service wrapper da AbacatePay.
2. Modelar tipos e validadores.
3. Implementar operacoes por dominio: clientes, cobrancas, PIX, saques e webhooks.
4. Testar tudo em Dev Mode.
5. Tratar observabilidade, retries e idempotencia.

## Checklist Final

- O token esta apenas no backend.
- O fluxo foi testado em Dev Mode.
- Valores monetarios usam centavos.
- Webhooks validam assinatura e sao idempotentes.
- Erros foram mapeados para contratos internos claros.
- Nao ha segredos expostos em logs ou respostas publicas.
