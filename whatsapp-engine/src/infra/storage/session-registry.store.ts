import { promises as fs } from "node:fs";
import path from "node:path";
import { env } from "../../config/env.js";
import type { SessionRecord } from "../../core/types.js";

const registryFile = path.join(env.sessionStorePath, "registry.json");

export class SessionRegistryStore {
  async init(): Promise<void> {
    await fs.mkdir(env.sessionStorePath, { recursive: true });
    try {
      await fs.access(registryFile);
    } catch {
      await fs.writeFile(registryFile, "[]", "utf-8");
    }
  }

  async list(): Promise<SessionRecord[]> {
    await this.init();
    const raw = await fs.readFile(registryFile, "utf-8");
    return JSON.parse(raw) as SessionRecord[];
  }

  async upsert(record: SessionRecord): Promise<void> {
    const records = await this.list();
    const idx = records.findIndex((r) => r.sessionId === record.sessionId);
    if (idx >= 0) records[idx] = record;
    else records.push(record);
    await fs.writeFile(registryFile, JSON.stringify(records, null, 2), "utf-8");
  }

  async remove(sessionId: string): Promise<void> {
    const records = await this.list();
    const filtered = records.filter((r) => r.sessionId !== sessionId);
    await fs.writeFile(registryFile, JSON.stringify(filtered, null, 2), "utf-8");
  }

  authPath(sessionId: string): string {
    return path.join(env.sessionStorePath, sessionId);
  }
}
