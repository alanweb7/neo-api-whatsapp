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
  type: "quick_reply" | "cta_url" | "cta_call" | "cta_copy";
  displayText: string;
  id: string;
  url?: string;
  phoneNumber?: string;
  copyCode?: string;
};

export type SendButtonsPayload = {
  jid: string;
  text: string;
  footer?: string;
  buttons: InteractiveButton[];
  fallback_text?: string;
};

export type CarouselCard = {
  title?: string;
  body: string;
  footer?: string;
  image_url: string;
  buttons: InteractiveButton[];
};

export type SendCarouselPayload = {
  jid: string;
  text: string;
  footer?: string;
  cards: CarouselCard[];
  fallback_text?: string;
};
