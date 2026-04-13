import dotenv from "dotenv";
import { z } from "zod";
dotenv.config();
const schema = z.object({
    NODE_ENV: z.enum(["development", "test", "production"]).default("development"),
    ENGINE_PORT: z.string().default("8090"),
    REDIS_URL: z.string().default("redis://redis:6379"),
    INTERNAL_API_KEY: z.string().min(12),
    WA_EVENTS_CHANNEL: z.string().default("wa.events.v1"),
    SESSION_STORE_PATH: z.string().default("./data/sessions")
});
const parsed = schema.parse(process.env);
export const env = {
    nodeEnv: parsed.NODE_ENV,
    port: Number(parsed.ENGINE_PORT),
    redisUrl: parsed.REDIS_URL,
    internalApiKey: parsed.INTERNAL_API_KEY,
    eventsChannel: parsed.WA_EVENTS_CHANNEL,
    sessionStorePath: parsed.SESSION_STORE_PATH
};
