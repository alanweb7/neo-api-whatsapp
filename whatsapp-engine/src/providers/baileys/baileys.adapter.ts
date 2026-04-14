import makeWASocket, {
  DisconnectReason,
  fetchLatestBaileysVersion,
  generateWAMessageFromContent,
  proto,
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

  async sendButtons(sessionId: string, payload: SendButtonsPayload): Promise<{ message_id: string; mode: "native_flow" | "fallback_text" | "legacy_buttons" }> {
    const s = this.mustConnected(sessionId);
    const jid = this.normalizeJid(payload.jid);

    try {
      const userJid = s.socket.user?.id;
      if (!userJid) {
        throw new Error("socket user not available");
      }

      const allQuickReply = payload.buttons.every((button) => button.type === "quick_reply");
      if (allQuickReply && payload.buttons.length <= 3) {
        try {
          const legacyMessageId = await this.sendLegacyQuickReplyButtons(s, jid, userJid, payload);
          await this.publish("message.sent", s.tenantId, s.sessionId, {
            id: legacyMessageId,
            to: jid,
            type: "buttons_legacy"
          });
          return { message_id: legacyMessageId, mode: "legacy_buttons" };
        } catch (legacyErr) {
          logger.warn(
            { err: legacyErr, sessionId, jid },
            "legacy buttonsMessage failed, trying native flow"
          );
        }
      }

      const nativeMessageId = await this.sendNativeFlowButtons(s, jid, userJid, payload);
      await this.publish("message.sent", s.tenantId, s.sessionId, { id: nativeMessageId, to: jid, type: "buttons_native_flow" });
      return { message_id: nativeMessageId, mode: "native_flow" };
    } catch (err) {
      logger.warn({ err, sessionId, jid }, "native flow buttons failed, sending fallback text");
      const fallbackText = payload.fallback_text ?? this.buildButtonsFallbackText(payload.text, payload.buttons);
      const fallbackRes = await s.socket.sendMessage(jid, { text: fallbackText });
      const fallbackId = this.extractMessageId(fallbackRes);
      await this.publish("engine.error", s.tenantId, s.sessionId, {
        category: "interactive_native_flow_failed",
        message: "native flow buttons failed; fallback text sent"
      });
      await this.publish("message.sent", s.tenantId, s.sessionId, {
        id: fallbackId,
        to: jid,
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

  private normalizeJid(raw: string): string {
    if (raw.includes("@")) return raw;
    return this.toJid(raw);
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
    const options = buttons.map((b, idx) => `${idx + 1} - ${b.displayText}`).join("\n");
    return `${text}\n\n${options}`;
  }

  private async sendLegacyQuickReplyButtons(
    session: RuntimeSession & { socket: WASocket },
    jid: string,
    userJid: string,
    payload: SendButtonsPayload
  ): Promise<string> {
    const buttonsMessage = proto.Message.ButtonsMessage.create({
      contentText: payload.text,
      footerText: payload.footer,
      headerType: proto.Message.ButtonsMessage.HeaderType.EMPTY,
      buttons: payload.buttons.map((button) => ({
        buttonId: button.id,
        buttonText: { displayText: button.displayText },
        type: proto.Message.ButtonsMessage.Button.Type.RESPONSE
      }))
    });

    const waMessage = generateWAMessageFromContent(
      jid,
      { buttonsMessage } as any,
      { userJid }
    );

    await session.socket.relayMessage(jid, waMessage.message as any, {
      messageId: waMessage.key.id ?? undefined
    });
    return waMessage.key.id ?? "unknown";
  }

  private async sendNativeFlowButtons(
    session: RuntimeSession & { socket: WASocket },
    jid: string,
    userJid: string,
    payload: SendButtonsPayload
  ): Promise<string> {
    const nativeFlowButtons = payload.buttons.map((btn) => {
      if (btn.type === "cta_url") {
        return {
          name: "cta_url",
          buttonParamsJson: JSON.stringify({
            display_text: btn.displayText,
            url: btn.url,
            merchant_url: btn.url
          })
        };
      }
      if (btn.type === "cta_call") {
        return {
          name: "cta_call",
          buttonParamsJson: JSON.stringify({
            display_text: btn.displayText,
            phone_number: btn.phoneNumber
          })
        };
      }
      if (btn.type === "cta_copy") {
        return {
          name: "cta_copy",
          buttonParamsJson: JSON.stringify({
            display_text: btn.displayText,
            copy_code: btn.copyCode
          })
        };
      }

      return {
        name: "quick_reply",
        buttonParamsJson: JSON.stringify({
          display_text: btn.displayText,
          id: btn.id
        })
      };
    });

    const interactiveMessage: any = {
      body: { text: payload.text || " " },
      nativeFlowMessage: {
        buttons: nativeFlowButtons,
        messageVersion: 1
      }
    };
    if (payload.footer) {
      interactiveMessage.footer = { text: payload.footer };
    }

    const directMessage = generateWAMessageFromContent(
      jid,
      {
        interactiveMessage: proto.Message.InteractiveMessage.create(interactiveMessage)
      } as any,
      { userJid }
    );

    try {
      await session.socket.relayMessage(jid, directMessage.message as any, {
        messageId: directMessage.key.id ?? undefined
      });
      return directMessage.key.id ?? "unknown";
    } catch (directErr) {
      logger.warn({ err: directErr, jid }, "direct interactiveMessage send failed, retrying viewOnce wrapper");
    }

    const wrappedMessage = generateWAMessageFromContent(
      jid,
      {
        viewOnceMessage: {
          message: {
            interactiveMessage: proto.Message.InteractiveMessage.create(interactiveMessage)
          }
        }
      } as any,
      { userJid }
    );

    await session.socket.relayMessage(jid, wrappedMessage.message as any, {
      messageId: wrappedMessage.key.id ?? undefined
    });
    return wrappedMessage.key.id ?? "unknown";
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
