import { z } from "zod";
import type { SendButtonsPayload, SendMediaPayload, SendTextPayload } from "../../core/types.js";
import { BaileysAdapter } from "../../providers/baileys/baileys.adapter.js";

const createSchema = z.object({ tenant_id: z.string().uuid(), name: z.string().min(2).max(120) });
const sendTextSchema = z.object({ to: z.string().min(8), text: z.string().min(1) });
const sendMediaSchema = z.object({ to: z.string().min(8), media_url: z.string().url(), caption: z.string().optional(), file_name: z.string().optional() });
const sendButtonsSchema = z.object({
  jid: z.string().min(8),
  text: z.string().min(1).max(1024),
  footer: z.string().max(60).optional(),
  fallback_text: z.string().max(1024).optional(),
  buttons: z.array(
    z.object({
      type: z.enum(["quick_reply", "cta_url", "cta_call", "cta_copy"]),
      displayText: z.string().min(1).max(40),
      id: z.string().min(1).max(128),
      url: z.string().url().optional(),
      phoneNumber: z.string().min(3).max(32).optional(),
      copyCode: z.string().min(1).max(500).optional()
    })
  ).min(1).max(3)
}).superRefine((payload, ctx) => {
  const hasQuickReply = payload.buttons.some((b) => b.type === "quick_reply");
  const hasCTA = payload.buttons.some((b) => b.type !== "quick_reply");
  if (hasQuickReply && hasCTA) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: "Do not mix quick_reply with CTA button types in the same payload."
    });
  }
});

export class SessionService {
  constructor(private readonly adapter: BaileysAdapter) {}

  async create(payload: unknown): Promise<{ session_id: string; status: string; qr_code?: string }> {
    const input = createSchema.parse(payload);
    const session = await this.adapter.createSession(input.tenant_id, input.name);
    return { session_id: session.sessionId, status: session.status, qr_code: session.qrCode };
  }

  async start(sessionId: string): Promise<{ session_id: string; status: string }> {
    const session = await this.adapter.startSession(sessionId);
    return { session_id: session.sessionId, status: session.status };
  }

  async status(sessionId: string): Promise<{ status: string; qr_code?: string }> {
    return this.adapter.getStatus(sessionId);
  }

  async qr(sessionId: string): Promise<{ qr_code?: string; status: string }> {
    return this.adapter.getStatus(sessionId);
  }

  async reconnect(sessionId: string): Promise<{ reconnecting: boolean }> {
    await this.adapter.startSession(sessionId);
    return { reconnecting: true };
  }

  async disconnect(sessionId: string): Promise<{ disconnected: boolean }> {
    await this.adapter.disconnect(sessionId);
    return { disconnected: true };
  }

  async remove(sessionId: string): Promise<{ removed: boolean }> {
    await this.adapter.remove(sessionId);
    return { removed: true };
  }

  async sendText(sessionId: string, payload: unknown): Promise<{ message_id: string }> {
    const parsed = sendTextSchema.parse(payload) as SendTextPayload;
    return this.adapter.sendText(sessionId, parsed);
  }

  async sendImage(sessionId: string, payload: unknown): Promise<{ message_id: string }> {
    const parsed = sendMediaSchema.parse(payload) as SendMediaPayload;
    return this.adapter.sendImage(sessionId, parsed);
  }

  async sendDocument(sessionId: string, payload: unknown): Promise<{ message_id: string }> {
    const parsed = sendMediaSchema.parse(payload) as SendMediaPayload;
    return this.adapter.sendDocument(sessionId, parsed);
  }

  async sendAudio(sessionId: string, payload: unknown): Promise<{ message_id: string }> {
    const parsed = sendMediaSchema.parse(payload) as SendMediaPayload;
    return this.adapter.sendAudio(sessionId, parsed);
  }

  async sendButtons(sessionId: string, payload: unknown): Promise<{ message_id: string; mode: "native_flow" | "fallback_text" }> {
    const parsed = sendButtonsSchema.parse(payload) as SendButtonsPayload;
    return this.adapter.sendButtons(sessionId, parsed);
  }

  async bootstrap(): Promise<void> {
    await this.adapter.bootstrap();
  }
}
