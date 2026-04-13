export class InternalSessionController {
    service;
    constructor(service) {
        this.service = service;
    }
    create = async (req, reply) => {
        const out = await this.service.create(req.body);
        reply.code(201).send(out);
    };
    start = async (req, reply) => {
        const out = await this.service.start(req.params.sessionId);
        reply.send(out);
    };
    status = async (req, reply) => {
        reply.send(await this.service.status(req.params.sessionId));
    };
    qr = async (req, reply) => {
        reply.send(await this.service.qr(req.params.sessionId));
    };
    reconnect = async (req, reply) => {
        reply.send(await this.service.reconnect(req.params.sessionId));
    };
    disconnect = async (req, reply) => {
        reply.send(await this.service.disconnect(req.params.sessionId));
    };
    remove = async (req, reply) => {
        reply.send(await this.service.remove(req.params.sessionId));
    };
    sendText = async (req, reply) => {
        reply.send(await this.service.sendText(req.params.sessionId, req.body));
    };
    sendImage = async (req, reply) => {
        reply.send(await this.service.sendImage(req.params.sessionId, req.body));
    };
    sendDocument = async (req, reply) => {
        reply.send(await this.service.sendDocument(req.params.sessionId, req.body));
    };
    sendAudio = async (req, reply) => {
        reply.send(await this.service.sendAudio(req.params.sessionId, req.body));
    };
    sendButtons = async (req, reply) => {
        reply.send(await this.service.sendButtons(req.params.sessionId, req.body));
    };
}
