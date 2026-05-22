# Docker Swarm Deployment Guide

## 📋 Pré-requisitos

- Docker Swarm inicializado: `docker swarm init` (se necessário)
- Traefik rodando na rede `traefik_public`
- Node com role `manager`

## 🚀 Deploy Steps

### 1️⃣ Criar volumes persistentes

```bash
docker volume create baileys_postgres_data
docker volume create baileys_redis_data
docker volume create baileys_wa_sessions
```

### 2️⃣ Atualizar as variáveis SECRET no docker-compose.swarm.yml

Edite o arquivo e substitua os valores `changeme123456789012345` por valores reais:

```yaml
environment:
  JWT_ACCESS_SECRET: seu-jwt-access-secret-super-secreto
  JWT_REFRESH_SECRET: seu-jwt-refresh-secret-super-secreto
  INTERNAL_API_KEY: seu-internal-api-key-super-secreto
```

**Valores mínimos recomendados:**
- `JWT_ACCESS_SECRET`: 32+ caracteres aleatórios
- `JWT_REFRESH_SECRET`: 32+ caracteres aleatórios
- `INTERNAL_API_KEY`: 32+ caracteres aleatórios

Gerar valores seguros:
```bash
openssl rand -base64 32
```

### 5️⃣ Verificar deploy

```bash
# Status dos serviços
docker stack services baileys

# Logs
docker service logs baileys_api-core-go
docker service logs baileys_whatsapp-engine

# Health check
curl https://zap-api.wesenderbrasil.com.br/healthz
curl https://engine.zap-api.wesenderbrasil.com.br/healthz
```

## 📌 URLs Disponíveis

- **API Core**: `https://zap-api.wesenderbrasil.com.br`
- **Engine**: `https://engine.zap-api.wesenderbrasil.com.br`
- **Alternativa API**: `https://api.zap-api.wesenderbrasil.com.br`

## 🔄 Atualizar Deploy

```bash
# Com novas imagens do GHCR
docker service update --image ghcr.io/alanweb7/baileys-core:latest baileys_api-core-go
docker service update --image ghcr.io/alanweb7/baileys-engine:latest baileys_whatsapp-engine
```

## 🛑 Remover Stack

```bash
docker stack rm baileys
```

## 🐳 Deploy via Portainer

1. **Acesse Portainer** → Stacks → Add Stack
2. **Upload** ou copie o conteúdo de `docker-compose.swarm.yml`
3. **Edite as variáveis antes de deploy:**
   - Clique em "Show advanced options"
   - Procure por `environment` seção do serviço `api-core-go`
   - Substitua os valores `changeme123456789012345`:
     - `JWT_ACCESS_SECRET` 
     - `JWT_REFRESH_SECRET`
     - `INTERNAL_API_KEY`
4. **Deploy** o stack

## 🔐 Security Best Practices

- ✅ Use Docker Secrets para valores sensíveis
- ✅ Configure HTTPS com Traefik
- ✅ Limitar resources (CPU/Memory)
- ✅ Use placeholders constraints para distribuir carga
- ✅ Manter imagens atualizadas do GHCR

## 📊 Monitoramento

```bash
# Ver recursos
docker stats

# Ver eventos
docker events --filter type=service

# Health checks
docker service ps baileys_api-core-go
docker service ps baileys_whatsapp-engine
```

## ⚠️ Troubleshooting

### Serviço não inicia
```bash
docker service logs baileys_api-core-go -f
```

### Traefik não roteia
- Verificar se `traefik_public` network existe
- Verificar labels do serviço
- Verificar DNS para `zap-api.wesenderbrasil.com.br`

### Banco de dados inacessível
```bash
# Entrar no container
docker exec -it <container-id> psql -U wa_user -d wa_saas
```
