# Arquitetura SaaS Multitenant WhatsApp

## Visão geral
- `core-go`: cérebro do produto (auth, tenants, sessões, API pública, logs, webhooks, orquestração)
- `whatsapp-engine`: motor especializado Baileys (conexão, QR, reconexão, envio/recebimento)
- Contrato desacoplado: HTTP interno para comandos + Redis Pub/Sub para eventos

## Decisões técnicas
1. Go + Gin no Core:
- bom equilíbrio entre performance, middlewares e produtividade
- ecossistema maduro para APIs corporativas

2. GORM + PostgreSQL:
- agilidade com modelagem inicial e evolução incremental
- migrations SQL explícitas para governança de schema

3. Redis event bus:
- simples para início, fácil migrar para NATS/Rabbit depois
- desacopla eventos de sessão/mensagem do request síncrono

4. Node + Fastify + TypeScript no Engine:
- HTTP interno leve e rápido
- camadas isoladas (`http -> service -> provider`) e Baileys encapsulado

## Fluxo principal
1. Core cria sessão no Engine via HTTP interno
2. Engine publica `session.created/session.qr.updated/...` no Redis
3. Core consome eventos e atualiza estado persistido (`whatsapp_sessions`, `message_logs`)
4. Core dispara webhooks externos (estrutura pronta para retries)

## Isolamento multitenant
- Todas as entidades relevantes têm `tenant_id`
- JWT inclui `tenant_id`
- Controllers protegidos resolvem tenant no contexto da requisição

## Escalabilidade
- Core e Engine escalam horizontalmente de forma independente
- Estado de sessão do Baileys persistido por sessão
- Event bus permite desacoplamento para workers dedicados no futuro
