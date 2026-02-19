export type ConnectionState = "connecting" | "connected" | "disconnected";
export type StateListener = (state: ConnectionState) => void;

const INITIAL_RECONNECT_DELAY_MS = 2000;
const MAX_RECONNECT_DELAY_MS = 30000;

export class WebSocketClient {
  private ws: WebSocket | null = null;
  private readonly stateListeners: StateListener[] = [];
  private messageListeners: ((data: ArrayBuffer) => void)[] = [];
  private shouldReconnect = true;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private reconnectDelay = INITIAL_RECONNECT_DELAY_MS;

  constructor(private readonly url: string) {}

  connect(): void {
    this.shouldReconnect = true;
    this.reconnectDelay = INITIAL_RECONNECT_DELAY_MS;
    this.open();
  }

  disconnect(): void {
    this.shouldReconnect = false;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    this.ws?.close();
  }

  send(data: Uint8Array): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(data);
    }
  }

  onMessage(handler: (data: ArrayBuffer) => void): () => void {
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
    this.ws?.close();
    this.notify("connecting");
    const ws = new WebSocket(this.url);
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
      this.notify("disconnected");
      this.ws = null;
      if (this.shouldReconnect) {
        this.reconnectTimer = setTimeout(() => this.open(), this.reconnectDelay);
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, MAX_RECONNECT_DELAY_MS);
      }
    };

    ws.onerror = (event) => {
      console.error("WebSocket error:", event);
      ws.close();
    };

    this.ws = ws;
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
