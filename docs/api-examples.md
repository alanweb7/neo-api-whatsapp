# Exemplos de API

## Auth
### Login
`POST /api/v1/auth/login`
```json
{
  "email": "admin@example.com",
  "password": "ChangeMe123!",
  "tenant_id": "<TENANT_UUID>"
}
```

### Refresh
`POST /api/v1/auth/refresh`
```json
{
  "refresh_token": "<JWT_REFRESH>",
  "tenant_id": "<TENANT_UUID>"
}
```

## Sessions
### Criar sessão
`POST /api/v1/sessions`
```json
{ "name": "Atendimento 01" }
```

### Iniciar sessão
`POST /api/v1/sessions/{sessionId}/start`

### Buscar QR
`GET /api/v1/sessions/{sessionId}/qr`

## Messages
### Texto
`POST /api/v1/messages/text`
```json
{
  "session_id": "<SESSION_UUID>",
  "to": "5511999999999",
  "text": "Olá!"
}
```

### Imagem
`POST /api/v1/messages/image`
```json
{
  "session_id": "<SESSION_UUID>",
  "to": "5511999999999",
  "media_url": "https://example.com/image.jpg",
  "caption": "Legenda"
}
```
