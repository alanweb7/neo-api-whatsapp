import { createClient } from "redis";
import { env } from "../../config/env.js";
import { logger } from "../../config/logger.js";
export class EventBus {
    redis;
    constructor() {
        this.redis = createClient({ url: env.redisUrl });
        this.redis.on("error", (err) => logger.error({ err }, "redis error"));
    }
    async connectIfNeeded() {
        if (!this.redis.isOpen) {
            await this.redis.connect();
        }
    }
    async publish(event) {
        await this.connectIfNeeded();
        await this.redis.publish(env.eventsChannel, JSON.stringify(event));
    }
    async health() {
        await this.connectIfNeeded();
        const res = await this.redis.ping();
        return res === "PONG";
    }
    async close() {
        if (this.redis.isOpen) {
            await this.redis.quit();
        }
    }
}
