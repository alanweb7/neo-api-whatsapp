import Fastify from "fastify";
import helmet from "@fastify/helmet";
import cors from "@fastify/cors";
import { logger } from "./config/logger.js";
import { EventBus } from "./infra/redis/event-bus.js";
import { SessionRegistryStore } from "./infra/storage/session-registry.store.js";
import { BaileysAdapter } from "./providers/baileys/baileys.adapter.js";
import { SessionService } from "./modules/sessions/session.service.js";
import { InternalSessionController } from "./http/controllers/internal-session.controller.js";
import { registerInternalRoutes } from "./http/routes/internal.routes.js";
export async function buildApp() {
    const app = Fastify({ logger: false });
    await app.register(helmet);
    await app.register(cors, { origin: false });
    const eventBus = new EventBus();
    const store = new SessionRegistryStore();
    await store.init();
    const adapter = new BaileysAdapter(eventBus, store);
    const sessionService = new SessionService(adapter);
    const controller = new InternalSessionController(sessionService);
    await registerInternalRoutes(app, controller, eventBus);
    app.setErrorHandler((error, _req, reply) => {
        logger.error({ err: error }, "request failed");
        reply.status(500).send({ error: "internal error", details: error.message });
    });
    return { app, eventBus, sessionService };
}
