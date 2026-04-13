Crie uma **base de arquitetura híbrida para um SaaS multitenant de WhatsApp**, usando:

- **Go** como backend principal
- **Node.js + TypeScript + Baileys** como engine de sessões WhatsApp

Quero uma base de projeto **realista, modular, escalável e pronta para evoluir para produção**, evitando estrutura de demo ou código improvisado.

## Objetivo do sistema
O sistema será uma plataforma SaaS multitenant para gerenciamento de múltiplas conexões WhatsApp por cliente (tenant), com API própria, autenticação, sessões independentes, envio de mensagens, recebimento de eventos e webhooks.

A arquitetura deve separar claramente:

1. **Core SaaS / API principal em Go**
2. **Engine de WhatsApp em Node.js com Baileys**
3. **Comunicação interna limpa entre os serviços**

## Diretriz principal de arquitetura
Quero que o **Go seja o cérebro do produto** e o **Node/Baileys seja o motor especializado de mensageria**.

### Go deve cuidar de:
- autenticação
- tenants
- usuários
- API pública
- API keys
- billing-ready design
- planos e limites futuros
- orquestração das sessões
- persistência principal
- webhooks
- auditoria
- observabilidade
- controle multitenant
- autorização
- rate limiting
- filas e jobs do core

### Node + Baileys deve cuidar de:
- criar sessão WhatsApp
- iniciar conexão
- gerar QR code
- reconectar
- desconectar
- manter auth state
- enviar mensagens
- receber mensagens
- emitir eventos técnicos da sessão
- encapsular toda complexidade do Baileys

## Princípios obrigatórios
Quero que o projeto siga estes princípios:

- separação clara de responsabilidades
- baixo acoplamento entre Go e Node
- comunicação por contrato bem definido
- preparado para múltiplos tenants e múltiplas sessões por tenant
- código limpo e profissional
- fácil de operar com Docker
- fácil de evoluir para produção
- preparado para observabilidade, filas e scale-out
- sem gambiarra
- sem arquivo monolítico gigante
- sem acoplar regra de negócio ao Baileys

## Stack desejada

### Core API
Use:
- **Go**
- framework HTTP: **Fiber**, **Gin** ou **Echo** (escolha e justifique)
- **PostgreSQL**
- **Redis**
- ORM ou query layer idiomática em Go, como **GORM**, **sqlc**, **bun** ou equivalente
- **JWT**
- **Swagger/OpenAPI**
- logs estruturados
- config por `.env`
- Docker / Docker Compose

### Engine WhatsApp
Use:
- **Node.js**
- **TypeScript**
- **Baileys**
- **Fastify** ou framework HTTP leve, se necessário
- Redis se fizer sentido
- persistência adequada para auth state e metadados da sessão

## Comunicação entre Go e Node
Quero uma estratégia clara e profissional.

Você pode propor um destes modelos, com justificativa:
- HTTP interno + eventos assíncronos
- Redis Pub/Sub
- NATS
- RabbitMQ

Minha preferência inicial:
- **HTTP interno para comandos**
- **event bus assíncrono para eventos**

Exemplo:
- Go chama o engine Node para criar/iniciar/enviar
- Node publica eventos de sessão e mensagem
- Go consome os eventos e atualiza estado / dispara webhooks

## Multitenancy
O sistema deve ser verdadeiramente multitenant.

Cada tenant deve ter:
- id
- nome
- slug ou identificador
- status
- plano
- limites futuros
- timestamps
- configuração opcional

Cada tenant pode ter:
- usuários
- sessões WhatsApp
- webhooks
- API keys
- logs/eventos
- configurações

O isolamento deve ser pelo menos **lógico por tenant**, com campo `tenant_id` em todas as entidades relevantes.

## Entidades mínimas do Core em Go
Modele ao menos:

- **Tenant**
- **User**
- **TenantUser**
- **ApiKey**
- **WhatsAppSession**
- **WebhookEndpoint**
- **WebhookDelivery**
- **MessageLog**
- **AuditLog**
- **Plan** ou estrutura preparada para plano/limites futuros

Sugira também entidades extras se fizer sentido.

## WhatsAppSession
A entidade de sessão deve contemplar algo como:

- id
- tenant_id
- external_engine_id ou engine_session_id
- nome
- status
- phone
- push_name
- qr_code ou ponteiro temporário
- last_seen_at
- connected_at
- disconnected_at
- failure_reason
- metadata json
- created_at
- updated_at

Status sugeridos:
- `created`
- `starting`
- `qr_pending`
- `connected`
- `disconnected`
- `failed`
- `reconnecting`

## Comandos que o Go deve enviar ao engine Node
Defina contratos para operações como:

- criar sessão
- iniciar sessão
- obter QR code
- consultar status
- reconectar sessão
- desconectar sessão
- remover sessão
- enviar texto
- enviar imagem
- enviar documento
- enviar áudio
- enviar lista
- enviar botão simples, se suportado
- healthcheck do engine

## Eventos que o Node deve emitir para o Go
Defina eventos como:

- session.created
- session.starting
- session.qr.updated
- session.connected
- session.disconnected
- session.reconnecting
- session.failed
- message.received
- message.sent
- message.delivery_update
- engine.error

Quero contratos claros, versionáveis e bem definidos.

## Requisitos da engine Node/Baileys
A engine deve:

- suportar múltiplas sessões simultâneas
- suportar múltiplos tenants
- persistir auth state de forma segura e organizada
- restaurar sessões no startup
- tratar reconexão
- isolar a lógica do Baileys em módulos/serviços limpos
- não misturar HTTP com regra de sessão
- ter camada própria de adapter/provider para Baileys
- expor API interna mínima e limpa
- publicar eventos de forma estruturada

Quero que o código trate o Baileys como uma dependência encapsulada.

## API do Core em Go
Quero endpoints iniciais para:

### Auth
- login
- refresh token
- me

### Tenants
- criar tenant
- listar tenants
- detalhar tenant
- atualizar tenant

### Users
- criar usuário
- listar usuários do tenant
- associar usuário ao tenant

### API Keys
- criar api key
- listar api keys
- revogar api key

### Sessions
- criar sessão
- iniciar sessão
- listar sessões
- detalhar sessão
- buscar QR code
- consultar status
- reconectar
- desconectar
- remover

### Messages
- enviar texto
- enviar imagem
- enviar documento
- enviar áudio
- listar logs de mensagens

### Webhooks
- criar webhook
- listar webhooks
- atualizar webhook
- remover webhook
- listar tentativas de entrega

## Webhooks
O Core em Go deve ser responsável por webhooks externos.

Requisitos:
- registrar múltiplos webhooks por tenant
- assinar payloads com secret opcional
- retry futuro preparado
- logar tentativas
- suportar eventos configuráveis
- desacoplar emissão de webhook do request principal

## Segurança
Quero:
- hash de senha
- JWT access + refresh token
- autorização por tenant
- contexto do tenant em cada request
- suporte futuro a RBAC
- API key por tenant para integrações
- validação de payloads
- tratamento centralizado de erros
- logs de auditoria para ações sensíveis

## Banco de dados
Quero modelagem inicial no PostgreSQL com:
- migrations
- seed inicial
- índices
- relacionamentos corretos
- campos de auditoria
- suporte a soft delete onde fizer sentido

## Redis
Use Redis para o que fizer sentido, por exemplo:
- cache transitório de QR
- pub/sub de eventos
- locks
- rate limiting
- filas futuras

## Observabilidade
A base deve nascer preparada para:
- logs estruturados
- correlation id / request id
- health endpoints
- readiness/liveness
- tracing futuro
- métricas futuras
- separação de logs do core e do engine

## Estrutura de pastas
Quero uma estrutura profissional.

### Core Go
Sugira algo como:
- `cmd/`
- `internal/`
- `pkg/`
- `configs/`
- `migrations/`
- `docs/`

ou equivalente idiomático em Go.

### Engine Node
Sugira algo como:
- `src/modules`
- `src/core`
- `src/infra`
- `src/providers`
- `src/events`
- `src/http`
- `src/config`

ou equivalente.

## Docker
Quero:
- `docker-compose.yml`
- serviços separados para:
  - api-core-go
  - whatsapp-engine-node
  - postgres
  - redis
- volumes necessários
- ambiente local fácil de subir
- healthchecks

## Entrega esperada
Quero que você entregue:

1. **visão arquitetural**
2. **decisões técnicas justificadas**
3. **árvore de pastas**
4. **contrato entre Go e Node**
5. **modelagem inicial do banco**
6. **código base dos dois serviços**
7. **docker compose**
8. **.env.example**
9. **migrations iniciais**
10. **endpoints básicos funcionando**
11. **bootstrap/restauração de sessões no engine**
12. **exemplos de requests/responses**
13. **guia curto para rodar localmente**

## Regras de implementação
- não criar código fake que não compila
- não simplificar demais o multitenancy
- não acoplar diretamente controllers HTTP com Baileys
- não enfiar tudo em um único serviço
- evitar abstrações inúteis, mas manter organização séria
- escrever código pensando em produto real
- deixar pontos claros para escalar horizontalmente
- se houver tradeoff entre rapidez e arquitetura correta, prefira arquitetura correta

## Importante sobre recursos WhatsApp
Considere que:
- texto, mídia e sessão são prioridade
- botões simples e listas podem existir como capacidades opcionais
- recursos mais frágeis do Baileys devem ficar isolados atrás de uma interface de capability
- não trate interativos avançados como garantidos

## Resultado desejado
Quero um **starter kit de produto SaaS**, não um exemplo acadêmico.

Se precisar escolher prioridades, priorize nesta ordem:
1. arquitetura correta
2. multitenancy
3. engine Baileys bem encapsulada
4. contratos entre serviços
5. facilidade de manutenção
6. operação local com Docker

Ao final, explique também:
- por que essa arquitetura híbrida é superior a colocar tudo em Node neste cenário
- quais são os pontos de atenção operacionais
- como evoluir depois para filas, billing, limites por plano e observabilidade
