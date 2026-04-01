import type { Connection, Clock } from "./types.js";
import { systemClock } from "./types.js";
import { encodePing, decode } from "../proto-helpers.js";

const PING_INTERVAL_MS = 5000;
const PONG_TIMEOUT_MS = 10000;
const TIMEOUT_CHECK_MS = 2000;

export class Heartbeat {
  private pingTimer: ReturnType<typeof setInterval> | null = null;
  private lastPong: number;
  private timeoutTimer: ReturnType<typeof setInterval> | null = null;
  private onTimeout: (() => void) | null = null;
  private unsubscribe: (() => void) | null = null;

  constructor(
    private readonly client: Connection,
    private readonly clock: Clock = systemClock,
  ) {
    this.lastPong = this.clock.now();
  }

  start(onTimeout: () => void): void {
    this.stop();
    this.onTimeout = onTimeout;
    this.lastPong = this.clock.now();

    this.unsubscribe = this.client.onMessage((data) => {
      const msg = decode(data);
      if (msg.pong) {
        this.lastPong = this.clock.now();
      }
    });

    this.pingTimer = this.clock.setInterval(() => {
      if (this.client.connected) {
        this.client.send(encodePing());
      }
    }, PING_INTERVAL_MS);

    this.timeoutTimer = this.clock.setInterval(() => {
      if (this.client.connected && this.clock.now() - this.lastPong > PONG_TIMEOUT_MS) {
        this.onTimeout?.();
      }
    }, TIMEOUT_CHECK_MS);
  }

  stop(): void {
    if (this.pingTimer) {
      this.clock.clearInterval(this.pingTimer);
      this.pingTimer = null;
    }
    if (this.timeoutTimer) {
      this.clock.clearInterval(this.timeoutTimer);
      this.timeoutTimer = null;
    }
    this.unsubscribe?.();
    this.unsubscribe = null;
  }
}
