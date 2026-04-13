import { env } from "../../config/env.js";
async function internalAuth(req, reply) {
    const key = req.headers["x-internal-api-key"];
    if (key !== env.internalApiKey) {
        reply.code(401).send({ error: "unauthorized" });
    }
}
const internalOpts = { preHandler: internalAuth };
export async function registerInternalRoutes(app, controller, bus) {
    app.get("/healthz", async () => ({ status: "ok" }));
    app.get("/readyz", async () => ({ status: (await bus.health()) ? "ready" : "not_ready" }));
    app.post("/internal/v1/sessions", internalOpts, controller.create);
    app.post("/internal/v1/sessions/:sessionId/start", internalOpts, controller.start);
    app.get("/internal/v1/sessions/:sessionId/status", internalOpts, controller.status);
    app.get("/internal/v1/sessions/:sessionId/qr", internalOpts, controller.qr);
    app.post("/internal/v1/sessions/:sessionId/reconnect", internalOpts, controller.reconnect);
    app.post("/internal/v1/sessions/:sessionId/disconnect", internalOpts, controller.disconnect);
    app.post("/internal/v1/sessions/:sessionId/remove", internalOpts, controller.remove);
    app.post("/internal/v1/sessions/:sessionId/messages/text", internalOpts, controller.sendText);
    app.post("/internal/v1/sessions/:sessionId/messages/image", internalOpts, controller.sendImage);
    app.post("/internal/v1/sessions/:sessionId/messages/document", internalOpts, controller.sendDocument);
    app.post("/internal/v1/sessions/:sessionId/messages/audio", internalOpts, controller.sendAudio);
}
