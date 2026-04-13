import { createClient, type RedisClientType } from "redis";
import { env } from "../../config/env.js";
import { logger } from "../../config/logger.js";
import type { EventEnvelope } from "../../core/types.js";

export class EventBus {
  private readonly redis: RedisClientType;

  constructor() {
    this.redis = createClient({ url: env.redisUrl });
    this.redis.on("error", (err) => logger.error({ err }, "redis error"));
  }

  private async connectIfNeeded(): Promise<void> {
    if (!this.redis.isOpen) {
      await this.redis.connect();
    }
  }

  async publish(event: EventEnvelope): Promise<void> {
    await this.connectIfNeeded();
    await this.redis.publish(env.eventsChannel, JSON.stringify(event));
  }

  async health(): Promise<boolean> {
    await this.connectIfNeeded();
    const res = await this.redis.ping();
    return res === "PONG";
  }

  async close(): Promise<void> {
    if (this.redis.isOpen) {
      await this.redis.quit();
    }
  }
}
