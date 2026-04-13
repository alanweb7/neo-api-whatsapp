import { z } from "zod";
const createSchema = z.object({ tenant_id: z.string().uuid(), name: z.string().min(2).max(120) });
const sendTextSchema = z.object({ to: z.string().min(8), text: z.string().min(1) });
const sendMediaSchema = z.object({ to: z.string().min(8), media_url: z.string().url(), caption: z.string().optional(), file_name: z.string().optional() });
const sendButtonsSchema = z.object({
    jid: z.string().min(8),
    text: z.string().min(1).max(1024),
    footer: z.string().max(60).optional(),
    fallback_text: z.string().max(1024).optional(),
    buttons: z.array(z.object({
        type: z.literal("quick_reply"),
        displayText: z.string().min(1).max(40),
        id: z.string().min(1).max(128)
    })).min(1).max(3)
});
export class SessionService {
    adapter;
    constructor(adapter) {
        this.adapter = adapter;
    }
    async create(payload) {
        const input = createSchema.parse(payload);
        const session = await this.adapter.createSession(input.tenant_id, input.name);
        return { session_id: session.sessionId, status: session.status, qr_code: session.qrCode };
    }
    async start(sessionId) {
        const session = await this.adapter.startSession(sessionId);
        return { session_id: session.sessionId, status: session.status };
    }
    async status(sessionId) {
        return this.adapter.getStatus(sessionId);
    }
    async qr(sessionId) {
        return this.adapter.getStatus(sessionId);
    }
    async reconnect(sessionId) {
        await this.adapter.startSession(sessionId);
        return { reconnecting: true };
    }
    async disconnect(sessionId) {
        await this.adapter.disconnect(sessionId);
        return { disconnected: true };
    }
    async remove(sessionId) {
        await this.adapter.remove(sessionId);
        return { removed: true };
    }
    async sendText(sessionId, payload) {
        const parsed = sendTextSchema.parse(payload);
        return this.adapter.sendText(sessionId, parsed);
    }
    async sendImage(sessionId, payload) {
        const parsed = sendMediaSchema.parse(payload);
        return this.adapter.sendImage(sessionId, parsed);
    }
    async sendDocument(sessionId, payload) {
        const parsed = sendMediaSchema.parse(payload);
        return this.adapter.sendDocument(sessionId, parsed);
    }
    async sendAudio(sessionId, payload) {
        const parsed = sendMediaSchema.parse(payload);
        return this.adapter.sendAudio(sessionId, parsed);
    }
    async sendButtons(sessionId, payload) {
        const parsed = sendButtonsSchema.parse(payload);
        return this.adapter.sendButtons(sessionId, parsed);
    }
    async bootstrap() {
        await this.adapter.bootstrap();
    }
}
