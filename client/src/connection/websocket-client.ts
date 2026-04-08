import type { Connection, ConnectionState, StateListener, MessageListener, Environment } from "./types.js";
import { browserEnvironment } from "./types.js";
export type { ConnectionState, StateListener } from "./types.js";

const INITIAL_RECONNECT_DELAY_MS = 2000;
const MAX_RECONNECT_DELAY_MS = 30000;

export class WebSocketClient implements Connection {
  private ws: WebSocket | null = null;
  private readonly stateListeners: StateListener[] = [];
  private messageListeners: MessageListener[] = [];
  private shouldReconnect = true;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private reconnectDelay = INITIAL_RECONNECT_DELAY_MS;
  private wasHidden = false;
  private unsubscribeVisibility: (() => void) | null = null;

  constructor(
    private readonly url: string,
    private readonly env: Environment = browserEnvironment,
  ) {}

  connect(): void {
    this.shouldReconnect = true;
    this.reconnectDelay = INITIAL_RECONNECT_DELAY_MS;
    this.unsubscribeVisibility = this.env.onVisibilityChange(() => {
      if (!this.env.isVisible()) {
        this.wasHidden = true;
        return;
      }
      if (this.shouldReconnect && this.wasHidden) {
        this.wasHidden = false;
        this.clearReconnectTimer();
        this.reconnectDelay = INITIAL_RECONNECT_DELAY_MS;
        this.open();
      }
    });
    this.open();
  }

  reconnect(): void {
    this.clearReconnectTimer();
    this.reconnectDelay = INITIAL_RECONNECT_DELAY_MS;
    this.ws?.close();
  }

  private clearReconnectTimer(): void {
    if (this.reconnectTimer) {
      this.env.clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
  }

  send(data: Uint8Array<ArrayBuffer>): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(data);
    }
  }

  onMessage(handler: MessageListener): () => void {
    this.messageListeners.push(handler);
    return () => {
      this.messageListeners = this.messageListeners.filter(h => h !== handler);
    };
  }

  onStateChange(listener: StateListener): void {
    this.stateListeners.push(listener);
  }

  get connected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  private open(): void {
    if (this.ws) {
      this.ws.onclose = null;
      this.ws.onerror = null;
      this.ws.close();
    }
    this.notify("connecting");
    const ws = this.env.createWebSocket(this.url);
    this.ws = ws;
    ws.binaryType = "arraybuffer";

    ws.onopen = () => {
      this.reconnectDelay = INITIAL_RECONNECT_DELAY_MS;
      this.notify("connected");
    };

    ws.onmessage = (event) => {
      if (event.data instanceof ArrayBuffer) {
        for (const listener of this.messageListeners) {
          listener(event.data);
        }
      }
    };

    ws.onclose = () => {
      if (this.ws !== ws) return;
      this.notify("disconnected");
      this.ws = null;
      if (this.shouldReconnect && this.env.isVisible()) {
        this.reconnectTimer = this.env.setTimeout(() => this.open(), this.reconnectDelay);
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, MAX_RECONNECT_DELAY_MS);
      }
    };

    ws.onerror = (event) => {
      console.error("WebSocket error:", event);
      ws.close();
    };
  }

  private notify(state: ConnectionState): void {
    for (const listener of this.stateListeners) {
      listener(state);
    }
  }
}

export function buildWsUrl(): string {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return `${protocol}//${window.location.host}/ws`;
}
