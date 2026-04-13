import type { FastifyReply, FastifyRequest } from "fastify";
import { SessionService } from "../../modules/sessions/session.service.js";

export class InternalSessionController {
  constructor(private readonly service: SessionService) {}

  create = async (req: FastifyRequest, reply: FastifyReply): Promise<void> => {
    const out = await this.service.create(req.body);
    reply.code(201).send(out);
  };

  start = async (req: FastifyRequest, reply: FastifyReply): Promise<void> => {
    const out = await this.service.start((req.params as { sessionId: string }).sessionId);
    reply.send(out);
  };

  status = async (req: FastifyRequest, reply: FastifyReply): Promise<void> => {
    reply.send(await this.service.status((req.params as { sessionId: string }).sessionId));
  };

  qr = async (req: FastifyRequest, reply: FastifyReply): Promise<void> => {
    reply.send(await this.service.qr((req.params as { sessionId: string }).sessionId));
  };

  reconnect = async (req: FastifyRequest, reply: FastifyReply): Promise<void> => {
    reply.send(await this.service.reconnect((req.params as { sessionId: string }).sessionId));
  };

  disconnect = async (req: FastifyRequest, reply: FastifyReply): Promise<void> => {
    reply.send(await this.service.disconnect((req.params as { sessionId: string }).sessionId));
  };

  remove = async (req: FastifyRequest, reply: FastifyReply): Promise<void> => {
    reply.send(await this.service.remove((req.params as { sessionId: string }).sessionId));
  };

  sendText = async (req: FastifyRequest, reply: FastifyReply): Promise<void> => {
    reply.send(await this.service.sendText((req.params as { sessionId: string }).sessionId, req.body));
  };

  sendImage = async (req: FastifyRequest, reply: FastifyReply): Promise<void> => {
    reply.send(await this.service.sendImage((req.params as { sessionId: string }).sessionId, req.body));
  };

  sendDocument = async (req: FastifyRequest, reply: FastifyReply): Promise<void> => {
    reply.send(await this.service.sendDocument((req.params as { sessionId: string }).sessionId, req.body));
  };

  sendAudio = async (req: FastifyRequest, reply: FastifyReply): Promise<void> => {
    reply.send(await this.service.sendAudio((req.params as { sessionId: string }).sessionId, req.body));
  };
}
