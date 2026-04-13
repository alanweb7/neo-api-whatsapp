export type EventEnvelope<T = Record<string, unknown>> = {
  version: "v1";
  type:
    | "session.created"
    | "session.starting"
    | "session.qr.updated"
    | "session.connected"
    | "session.disconnected"
    | "session.reconnecting"
    | "session.failed"
    | "message.received"
    | "message.sent"
    | "message.delivery_update"
    | "engine.error";
  tenant_id: string;
  session_id?: string;
  timestamp: string;
  payload: T;
};

export type SessionRecord = {
  sessionId: string;
  tenantId: string;
  name: string;
  status: string;
  createdAt: string;
  updatedAt: string;
};

export type SendTextPayload = {
  to: string;
  text: string;
};

export type SendMediaPayload = {
  to: string;
  media_url: string;
  caption?: string;
  file_name?: string;
};

export type InteractiveButton = {
  type: "quick_reply";
  displayText: string;
  id: string;
};

export type SendButtonsPayload = {
  jid: string;
  text: string;
  footer?: string;
  buttons: InteractiveButton[];
  fallback_text?: string;
};
