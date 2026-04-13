import { buildApp } from "./app.js";
import { env } from "./config/env.js";
import { logger } from "./config/logger.js";
const run = async () => {
    const { app, eventBus, sessionService } = await buildApp();
    await sessionService.bootstrap();
    const close = async () => {
        logger.info("shutting down engine");
        await app.close();
        await eventBus.close();
        process.exit(0);
    };
    process.on("SIGINT", close);
    process.on("SIGTERM", close);
    await app.listen({ host: "0.0.0.0", port: env.port });
    logger.info({ port: env.port }, "whatsapp engine started");
};
run().catch((err) => {
    logger.fatal({ err }, "engine failed to start");
    process.exit(1);
});
