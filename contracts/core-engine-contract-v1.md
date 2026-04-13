# Contrato Core <-> Engine (v1)

## Estratégia
- Comandos síncronos: HTTP interno (`X-Internal-API-Key`)
- Eventos assíncronos: Redis Pub/Sub canal `wa.events.v1`
- Envelope versionado: `version = "v1"`

## HTTP interno (Core -> Engine)
Base: `http://whatsapp-engine:8090/internal/v1`

### Criar sessão
`POST /sessions`
```json
{ "tenant_id": "uuid", "name": "Atendimento Principal" }
```
Resposta:
```json
{ "session_id": "uuid", "status": "created", "qr_code": null }
```

### Iniciar sessão
`POST /sessions/:sessionId/start`

### QR da sessão
`GET /sessions/:sessionId/qr`

### Status da sessão
`GET /sessions/:sessionId/status`

### Reconnect / Disconnect / Remove
- `POST /sessions/:sessionId/reconnect`
- `POST /sessions/:sessionId/disconnect`
- `POST /sessions/:sessionId/remove`

### Mensagens
- `POST /sessions/:sessionId/messages/text`
- `POST /sessions/:sessionId/messages/image`
- `POST /sessions/:sessionId/messages/document`
- `POST /sessions/:sessionId/messages/audio`

Payload base mídia:
```json
{
  "to": "5511999999999",
  "media_url": "https://example.com/file.jpg",
  "caption": "opcional",
  "file_name": "opcional"
}
```

## Eventos (Engine -> Core)
Canal Redis: `wa.events.v1`

Envelope:
```json
{
  "version": "v1",
  "type": "session.connected",
  "tenant_id": "uuid",
  "session_id": "uuid",
  "timestamp": "2026-04-12T20:45:00Z",
  "payload": {}
}
```

### Tipos previstos
- `session.created`
- `session.starting`
- `session.qr.updated`
- `session.connected`
- `session.disconnected`
- `session.reconnecting`
- `session.failed`
- `message.received`
- `message.sent`
- `message.delivery_update`
- `engine.error`

## Versionamento
- Alterações compatíveis: adicionar campos opcionais
- Alterações quebrando contrato: novo `version` + novo prefixo de rota (`/internal/v2`)
