import { WebSocketClient } from "./websocket-client.js";
import { encodePing, decode } from "../proto-helpers.js";

const PING_INTERVAL_MS = 5000;
const PONG_TIMEOUT_MS = 10000;
const TIMEOUT_CHECK_MS = 2000;

export class Heartbeat {
  private pingTimer: ReturnType<typeof setInterval> | null = null;
  private lastPong = Date.now();
  private timeoutTimer: ReturnType<typeof setInterval> | null = null;
  private onTimeout: (() => void) | null = null;
  private unsubscribe: (() => void) | null = null;

  constructor(private client: WebSocketClient) {}

  start(onTimeout: () => void): void {
    this.stop();
    this.onTimeout = onTimeout;
    this.lastPong = Date.now();

    this.unsubscribe = this.client.onMessage((data) => {
      const msg = decode(data);
      if (msg.pong) {
        this.lastPong = Date.now();
      }
    });

    this.pingTimer = setInterval(() => {
      if (this.client.connected) {
        this.client.send(encodePing());
      }
    }, PING_INTERVAL_MS);

    this.timeoutTimer = setInterval(() => {
      if (this.client.connected && Date.now() - this.lastPong > PONG_TIMEOUT_MS) {
        this.onTimeout?.();
      }
    }, TIMEOUT_CHECK_MS);
  }

  stop(): void {
    if (this.pingTimer) {
      clearInterval(this.pingTimer);
      this.pingTimer = null;
    }
    if (this.timeoutTimer) {
      clearInterval(this.timeoutTimer);
      this.timeoutTimer = null;
    }
    this.unsubscribe?.();
    this.unsubscribe = null;
  }
}
