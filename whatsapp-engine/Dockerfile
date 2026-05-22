FROM node:22-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY tsconfig.json ./
COPY src ./src
RUN npm run build
RUN ls -la dist/ || echo "dist not found"

FROM node:22-alpine
WORKDIR /app
RUN apk add --no-cache wget
COPY package*.json ./
RUN npm ci --omit=dev
COPY --from=builder /app/dist ./dist
RUN mkdir -p /app/data/sessions
EXPOSE 8090
HEALTHCHECK --interval=10s --timeout=5s --retries=10 --start-period=30s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8090/healthz || exit 1
CMD ["node", "dist/index.js"]
