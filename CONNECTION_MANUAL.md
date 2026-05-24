# Manual de Conexão WhatsApp

## 📱 Como Conectar uma Conta WhatsApp

Este guia mostra como autenticar uma conta WhatsApp na plataforma Baileys API e começar a enviar mensagens.

---

## 🔄 Fluxo de Conexão

```
Criar Sessão → Gerar QR Code → Escanear com WhatsApp → Autenticado ✅ → Enviar Mensagens
```

---

## ✅ Pré-requisitos

- ✔️ Conta de usuário criada
- ✔️ Tenant criado
- ✔️ Access Token JWT válido
- ✔️ WhatsApp instalado no celular
- ✔️ Conexão de internet estável

---

## 📋 Passo a Passo

### **1️⃣ Criar uma Sessão WhatsApp**

Uma sessão representa uma conexão com uma conta WhatsApp.

**Endpoint:**
```
POST /api/v1/sessions
```

**Headers:**
```
Authorization: Bearer {ACCESS_TOKEN}
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Meu WhatsApp Principal"
}
```

**cURL:**
```bash
curl -X POST https://zap-api.wesenderbrasil.com.br/api/v1/sessions \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Meu WhatsApp Principal"
  }'
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "tenantId": "tenant-uuid",
  "name": "Meu WhatsApp Principal",
  "status": "created",
  "createdAt": "2026-05-22T19:30:00Z"
}
```

**⚠️ Salve o `id` da sessão para os próximos passos!**

---

### **2️⃣ Iniciar a Sessão e Gerar QR Code**

Inicie a sessão para gerar o QR Code.

**Endpoint:**
```
POST /api/v1/sessions/{sessionId}/start
```

**Headers:**
```
Authorization: Bearer {ACCESS_TOKEN}
X-Engine-Session-ID: {OPTIONAL - engine_session_id para priorizar com API Key}
```

**Opção A: Com sessionId no caminho**
```bash
curl -X POST https://zap-api.wesenderbrasil.com.br/api/v1/sessions/550e8400-e29b-41d4-a716-446655440000/start \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN"
```

**Opção B: Com engine_session_id no header (prioritário)**
```bash
curl -X POST https://zap-api.wesenderbrasil.com.br/api/v1/sessions/any-value/start \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN" \
  -H "X-Engine-Session-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "qr_code_pending",
  "qrCode": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAHQAAAB0CAYAAABr...",
  "expiresIn": 60
}
```

**💡 O QR Code é uma imagem PNG em base64 que expira em ~60 segundos**

---

### **3️⃣ Exibir o QR Code**

Você pode exibir o QR Code de várias formas:

#### **Opção A: Em HTML**
```html
<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAHQAAAB0CAYAAABr..." />
```

#### **Opção B: Em Terminal**
```bash
# Salve a imagem base64 em um arquivo e abra
echo "iVBORw0KGgoAAAANSUhEUgAAAHQAAAB0CAYAAABr..." | base64 -d > qrcode.png
open qrcode.png  # macOS
xdg-open qrcode.png  # Linux
start qrcode.png  # Windows
```

#### **Opção C: QR Code Scanner Online**
Se tiver a base64, pode usar um decodificador online para visualizar.

---

### **4️⃣ Escanear o QR Code com WhatsApp**

**No seu celular:**

1. 📱 Abra o **WhatsApp**
2. Toque em **Configurações** (ou ⚙️)
3. Vá para **Aparelhos conectados** (ou "Linked devices")
4. Toque em **Conectar um aparelho** (ou "Link a device")
5. 📸 **Escaneie o QR Code** com a câmera do seu celular

**Aguarde a autenticação completar (~5-10 segundos)**

---

### **5️⃣ Verificar Status de Autenticação**

Verifique se a sessão foi autenticada com sucesso.

**Endpoint:**
```
GET /api/v1/sessions/{sessionId}/status
```

**cURL:**
```bash
curl -X GET https://zap-api.wesenderbrasil.com.br/api/v1/sessions/550e8400-e29b-41d4-a716-446655440000/status \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN"
```

**Response (200 OK - Autenticado):**
```json
{
  "status": "authenticated",
  "phoneNumber": "5511999999999",
  "isOnline": true,
  "lastActivity": "2026-05-22T19:35:00Z"
}
```

**Response (200 OK - Pendente):**
```json
{
  "status": "qr_code_pending",
  "isOnline": false
}
```

**✅ Se `status: "authenticated"` aparecer, você está conectado!**

---

### **6️⃣ Enviar uma Mensagem de Teste**

Teste a conexão enviando uma mensagem.

**Endpoint:**
```
POST /api/v1/messages/text
```

**Request Body:**
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "to": "5511987654321",
  "text": "Olá! Testando a API Baileys 🎉"
}
```

**cURL:**
```bash
curl -X POST https://zap-api.wesenderbrasil.com.br/api/v1/messages/text \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "to": "5511987654321",
    "text": "Olá! Testando a API Baileys 🎉"
  }'
```

**Response (201 Created):**
```json
{
  "id": "message-uuid",
  "sessionId": "session-uuid",
  "to": "5511987654321",
  "text": "Olá! Testando a API Baileys 🎉",
  "type": "text",
  "status": "sent",
  "timestamp": "2026-05-22T19:36:00Z"
}
```

**✅ Mensagem enviada com sucesso!**

---

## 🔐 Métodos de Autenticação

A API Baileys suporta **3 métodos de autenticação** diferentes dependendo da rota:

### **1️⃣ JWT (Access Token)**
Usado na maioria das rotas protegidas. O token expira automaticamente.

**Header:**
```
Authorization: Bearer SEU_ACCESS_TOKEN
```

**Exemplo:**
```bash
curl -X GET https://zap-api.wesenderbrasil.com.br/api/v1/sessions \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

### **2️⃣ API Key**
Sem expiração. Use para integrações de longa duração.

**Headers aceitos:**
```
X-API-Key: sua-api-key-aqui
```
ou
```
api-key: sua-api-key-aqui
```

**Exemplo:**
```bash
curl -X GET https://zap-api.wesenderbrasil.com.br/api/v1/sessions \
  -H "X-API-Key: 065e8658-9f0d-44a4-a1f0-12f599019ebd"
```

### **3️⃣ INTERNAL_API_KEY (Para Criar Sessões)**
Chave interna para criar novas sessões. Não requer usuário autenticado, mas precisa do `tenant_id`.

**Headers aceitos:**
```
X-Internal-Key: INTERNAL_API_KEY_VALUE
```
ou
```
api-key: INTERNAL_API_KEY_VALUE
```

**Body (obrigatório fornecer tenant_id):**
```json
{
  "name": "Nome da Sessão",
  "tenant_id": "uuid-do-tenant"
}
```

**Exemplo:**
```bash
curl -X POST https://zap-api.wesenderbrasil.com.br/api/v1/sessions \
  -H "X-Internal-Key: changeme123456789012345" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "WhatsApp Business",
    "tenant_id": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

**⚠️ Importante:** Quando criando com `INTERNAL_API_KEY`, você **deve** fornecer o `tenant_id` no body!

### **4️⃣ Engine Session ID (Para Iniciar Sessão)**
Use o `engine_session_id` para iniciar uma sessão sem expiração de token.

**Header:**
```
Authorization: Bearer SEU_ACCESS_TOKEN
X-Engine-Session-ID: ID_DA_SESSAO
```

**Exemplo:**
```bash
curl -X POST https://zap-api.wesenderbrasil.com.br/api/v1/sessions/550e8400-e29b-41d4-a716-446655440000/start \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN" \
  -H "X-Engine-Session-ID: 550e8400-e29b-41d4-a716-446655440000"
```

---

## 📍 Quais rotas usam qual autenticação?

| Rota | Método | Autenticação | Header |
|------|--------|--------------|--------|
| POST /sessions | **Criar** | INTERNAL_API_KEY | `X-Internal-Key` ou `api-key` |
| POST /sessions/:id/start | **Iniciar** | JWT + Engine Session | `Authorization` + `X-Engine-Session-ID` |
| GET /sessions | **Listar** | JWT ou INTERNAL_API_KEY | `Authorization` ou `X-Internal-Key` |
| GET /sessions/:id | **Obter** | JWT | `Authorization` |
| POST /messages/text | **Enviar** | JWT | `Authorization` |

---

## 🔍 Obter QR Code Later

Se o QR Code expirou, você pode obter um novo:

```bash
curl -X GET https://zap-api.wesenderbrasil.com.br/api/v1/sessions/550e8400-e29b-41d4-a716-446655440000/qr \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN"
```

---

## 🔄 Gerenciar Sessões

### **Listar Todas as Sessões**

**Com JWT:**
```bash
curl -X GET https://zap-api.wesenderbrasil.com.br/api/v1/sessions \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN"
```

**Com INTERNAL_API_KEY:**
```bash
curl -X GET https://zap-api.wesenderbrasil.com.br/api/v1/sessions \
  -H "X-Internal-Key: changeme123456789012345"
```

### **Obter Detalhes de uma Sessão**
```bash
curl -X GET https://zap-api.wesenderbrasil.com.br/api/v1/sessions/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN"
```

### **Reconectar uma Sessão Desconectada**
```bash
curl -X POST https://zap-api.wesenderbrasil.com.br/api/v1/sessions/550e8400-e29b-41d4-a716-446655440000/reconnect \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN"
```

### **Desconectar uma Sessão**
```bash
curl -X POST https://zap-api.wesenderbrasil.com.br/api/v1/sessions/550e8400-e29b-41d4-a716-446655440000/disconnect \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN"
```

### **Remover uma Sessão**
```bash
curl -X DELETE https://zap-api.wesenderbrasil.com.br/api/v1/sessions/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN"
```

---

## ⚠️ Troubleshooting

### **QR Code expirou**
- ❌ Erro: "QR Code expired"
- ✅ Solução: Chame o endpoint `/start` novamente para gerar um novo QR Code

### **Sessão não autentica**
- ❌ Erro: `status: "qr_code_pending"` após vários minutos
- ✅ Solução:
  1. Verifique se tem internet estável no celular
  2. Tente escanear novamente
  3. Reinicie o WhatsApp no celular
  4. Crie uma nova sessão

### **Mensagem não é enviada**
- ❌ Erro: `"session not authenticated"`
- ✅ Solução: Verifique o status da sessão com `/status`

### **Sessão desconecta**
- ❌ Status: `"disconnected"`
- ✅ Solução: Use o endpoint `/reconnect` para reconectar

### **WhatsApp pede verificação 2FA**
- ℹ️ Se sua conta tem 2FA ativo, o QR Code pode não aparecer
- ✅ Solução: Desative 2FA temporariamente ou escaneie em 60 segundos

---

## 📊 Estatísticas

| Métrica | Valor |
|---------|-------|
| Tempo para gerar QR | ~2 segundos |
| Validade do QR | ~60 segundos |
| Tempo para autenticar | ~5-10 segundos |
| Tempo para enviar mensagem | ~1-2 segundos |
| Sessões simultâneas | Ilimitado |

---

## 🔐 Segurança

- ✅ A senha do WhatsApp **nunca é pedida**
- ✅ Usa apenas o QR Code da API oficial do WhatsApp
- ✅ A sessão é armazenada com segurança
- ✅ Tokens JWT expiram automaticamente
- ✅ Sempre use HTTPS em produção

---

## 📞 Tipos de Mensagens Suportadas

Além de texto, você pode enviar:

- 📝 **Texto**: `/messages/text`
- 🖼️ **Imagem**: `/messages/image`
- 📄 **Documento**: `/messages/document`
- 🎵 **Áudio/Voz**: `/messages/audio`
- 🔘 **Botões**: `/messages/buttons`
- 🎠 **Carrossel**: `/messages/carousel`

Veja o [API_MANUAL.md](API_MANUAL.md) para mais detalhes!

---

## 🚀 Próximos Passos

1. ✅ Criar e autenticar sessão
2. 📤 Enviar mensagens
3. 📨 Configurar webhooks para receber mensagens
4. 📊 Monitorar logs e status
5. 🔄 Escalar para múltiplas sessões

---

**Versão**: 1.0.0  
**Última atualização**: 2026-05-22  
**Suporte**: support@wesenderbrasil.com.br
