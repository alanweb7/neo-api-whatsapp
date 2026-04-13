import makeWASocket, {
  DisconnectReason,
  fetchLatestBaileysVersion,
  useMultiFileAuthState,
  type WASocket
} from "@whiskeysockets/baileys";
import { Boom } from "@hapi/boom";
import type { EventEnvelope, SendButtonsPayload, SendMediaPayload, SendTextPayload } from "../../core/types.js";
import { logger } from "../../config/logger.js";
import type { EventBus } from "../../infra/redis/event-bus.js";
import type { SessionRegistryStore } from "../../infra/storage/session-registry.store.js";

type RuntimeSession = {
  sessionId: string;
  tenantId: string;
  name: string;
  status: string;
  socket?: WASocket;
  qrCode?: string;
};

export class BaileysAdapter {
  private readonly sessions = new Map<string, RuntimeSession>();

  constructor(private readonly eventBus: EventBus, private readonly store: SessionRegistryStore) {}

  async createSession(tenantId: string, name: string): Promise<RuntimeSession> {
    const sessionId = crypto.randomUUID();
    const now = new Date().toISOString();
    const runtime: RuntimeSession = { sessionId, tenantId, name, status: "created" };
    this.sessions.set(sessionId, runtime);
    await this.store.upsert({ sessionId, tenantId, name, status: "created", createdAt: now, updatedAt: now });

    await this.publish("session.created", tenantId, sessionId, { name });
    return runtime;
  }

  async startSession(sessionId: string): Promise<RuntimeSession> {
    const session = this.sessions.get(sessionId);
    if (!session) throw new Error("session not found");

    const authDir = this.store.authPath(sessionId);
    const { state, saveCreds } = await useMultiFileAuthState(authDir);
    const { version } = await fetchLatestBaileysVersion();

    const sock = makeWASocket({
      auth: state,
      version,
      printQRInTerminal: false,
      syncFullHistory: false
    });

    sock.ev.on("creds.update", saveCreds);
    sock.ev.on("connection.update", async (update) => {
      const { connection, lastDisconnect, qr } = update;

      if (qr) {
        session.qrCode = qr;
        session.status = "qr_pending";
        await this.publish("session.qr.updated", session.tenantId, session.sessionId, { qr_code: qr });
        await this.persistStatus(session, "qr_pending");
      }

      if (connection === "open") {
        session.status = "connected";
        await this.publish("session.connected", session.tenantId, session.sessionId, {});
        await this.persistStatus(session, "connected");
      }

      if (connection === "close") {
        const code = (lastDisconnect?.error as Boom | undefined)?.output?.statusCode;
        const shouldReconnect = code !== DisconnectReason.loggedOut;
        session.status = shouldReconnect ? "reconnecting" : "disconnected";

        await this.publish(
          shouldReconnect ? "session.reconnecting" : "session.disconnected",
          session.tenantId,
          session.sessionId,
          { reason_code: code ?? null }
        );
        await this.persistStatus(session, session.status);

        if (shouldReconnect) {
          await this.startSession(sessionId);
        }
      }
    });

    sock.ev.on("messages.upsert", async (m) => {
      for (const msg of m.messages) {
        await this.publish("message.received", session.tenantId, session.sessionId, {
          id: msg.key.id,
          from: msg.key.remoteJid,
          type: m.type,
          message: msg.message
        });
      }
    });

    session.socket = sock;
    session.status = "starting";
    await this.publish("session.starting", session.tenantId, session.sessionId, {});
    await this.persistStatus(session, "starting");

    return session;
  }

  async getStatus(sessionId: string): Promise<{ status: string; qr_code?: string }> {
    const s = this.sessions.get(sessionId);
    if (!s) throw new Error("session not found");
    return { status: s.status, qr_code: s.qrCode };
  }

  async disconnect(sessionId: string): Promise<void> {
    const s = this.sessions.get(sessionId);
    if (!s) throw new Error("session not found");
    if (s.socket) s.socket.end(new Error("disconnect requested"));
    s.status = "disconnected";
    await this.persistStatus(s, "disconnected");
    await this.publish("session.disconnected", s.tenantId, s.sessionId, { reason: "manual" });
  }

  async remove(sessionId: string): Promise<void> {
    const s = this.sessions.get(sessionId);
    if (!s) return;
    if (s.socket) s.socket.end(new Error("session removed"));
    this.sessions.delete(sessionId);
    await this.store.remove(sessionId);
  }

  async sendText(sessionId: string, payload: SendTextPayload): Promise<{ message_id: string }> {
    const s = this.mustConnected(sessionId);
    const jid = this.toJid(payload.to);
    const res = await s.socket.sendMessage(jid, { text: payload.text });
    const messageId = this.extractMessageId(res);
    await this.publish("message.sent", s.tenantId, s.sessionId, { id: messageId, to: payload.to, type: "text" });
    return { message_id: messageId };
  }

  async sendImage(sessionId: string, payload: SendMediaPayload): Promise<{ message_id: string }> {
    const s = this.mustConnected(sessionId);
    const jid = this.toJid(payload.to);
    const res = await s.socket.sendMessage(jid, { image: { url: payload.media_url }, caption: payload.caption });
    const messageId = this.extractMessageId(res);
    await this.publish("message.sent", s.tenantId, s.sessionId, { id: messageId, to: payload.to, type: "image" });
    return { message_id: messageId };
  }

  async sendDocument(sessionId: string, payload: SendMediaPayload): Promise<{ message_id: string }> {
    const s = this.mustConnected(sessionId);
    const jid = this.toJid(payload.to);
    const res = await s.socket.sendMessage(jid, {
      document: { url: payload.media_url },
      fileName: payload.file_name ?? "document",
      mimetype: "application/octet-stream",
      caption: payload.caption
    });
    const messageId = this.extractMessageId(res);
    await this.publish("message.sent", s.tenantId, s.sessionId, { id: messageId, to: payload.to, type: "document" });
    return { message_id: messageId };
  }

  async sendAudio(sessionId: string, payload: SendMediaPayload): Promise<{ message_id: string }> {
    const s = this.mustConnected(sessionId);
    const jid = this.toJid(payload.to);
    const res = await s.socket.sendMessage(jid, { audio: { url: payload.media_url }, ptt: false });
    const messageId = this.extractMessageId(res);
    await this.publish("message.sent", s.tenantId, s.sessionId, { id: messageId, to: payload.to, type: "audio" });
    return { message_id: messageId };
  }

  async sendButtons(sessionId: string, payload: SendButtonsPayload): Promise<{ message_id: string; mode: "buttons" | "fallback_text" }> {
    const s = this.mustConnected(sessionId);
    const jid = this.toJid(payload.to);
    const buttons = payload.buttons.map((b) => ({
      buttonId: b.id,
      buttonText: { displayText: b.title },
      type: 1
    }));

    try {
      const interactivePayload: any = {
        text: payload.text,
        footer: payload.footer,
        buttons,
        headerType: 1
      };
      const res = await s.socket.sendMessage(jid, interactivePayload);
      const messageId = this.extractMessageId(res);
      await this.publish("message.sent", s.tenantId, s.sessionId, { id: messageId, to: payload.to, type: "buttons" });
      return { message_id: messageId, mode: "buttons" };
    } catch (err) {
      logger.warn({ err, sessionId, to: payload.to }, "interactive buttons not supported, sending fallback text");
      const fallbackText = payload.fallback_text ?? this.buildButtonsFallbackText(payload.text, payload.buttons);
      const fallbackRes = await s.socket.sendMessage(jid, { text: fallbackText });
      const fallbackId = this.extractMessageId(fallbackRes);
      await this.publish("engine.error", s.tenantId, s.sessionId, {
        category: "interactive_buttons_unsupported",
        message: "buttons payload failed; fallback text sent"
      });
      await this.publish("message.sent", s.tenantId, s.sessionId, {
        id: fallbackId,
        to: payload.to,
        type: "text_fallback",
        original_type: "buttons"
      });
      return { message_id: fallbackId, mode: "fallback_text" };
    }
  }

  async bootstrap(): Promise<void> {
    const records = await this.store.list();
    for (const r of records) {
      this.sessions.set(r.sessionId, {
        sessionId: r.sessionId,
        tenantId: r.tenantId,
        name: r.name,
        status: r.status
      });
      if (r.status !== "disconnected" && r.status !== "failed") {
        try {
          await this.startSession(r.sessionId);
        } catch (err) {
          logger.error({ err, sessionId: r.sessionId }, "failed to restore session");
          await this.publish("session.failed", r.tenantId, r.sessionId, { reason: "bootstrap_failed" });
        }
      }
    }
  }

  private mustConnected(sessionId: string): RuntimeSession & { socket: WASocket } {
    const s = this.sessions.get(sessionId);
    if (!s || !s.socket) throw new Error("session not found or not started");
    return s as RuntimeSession & { socket: WASocket };
  }

  private toJid(raw: string): string {
    const normalized = raw.replace(/\D/g, "");
    return `${normalized}@s.whatsapp.net`;
  }

  private async persistStatus(session: RuntimeSession, status: string): Promise<void> {
    const records = await this.store.list();
    const found = records.find((r) => r.sessionId === session.sessionId);
    if (!found) return;
    found.status = status;
    found.updatedAt = new Date().toISOString();
    await this.store.upsert(found);
  }

  private extractMessageId(res: unknown): string {
    const id = (res as { key?: { id?: string } } | undefined)?.key?.id;
    return id ?? "unknown";
  }

  private buildButtonsFallbackText(text: string, buttons: SendButtonsPayload["buttons"]): string {
    const options = buttons.map((b, idx) => `${idx + 1} - ${b.title}`).join("\n");
    return `${text}\n\n${options}`;
  }

  private async publish(
    type: EventEnvelope["type"],
    tenantId: string,
    sessionId: string,
    payload: Record<string, unknown>
  ): Promise<void> {
    await this.eventBus.publish({
      version: "v1",
      type,
      tenant_id: tenantId,
      session_id: sessionId,
      timestamp: new Date().toISOString(),
      payload
    });
  }
}
