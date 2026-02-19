export type ConnectionState = "connecting" | "connected" | "disconnected";
export type StateListener = (state: ConnectionState) => void;

const RECONNECT_DELAY_MS = 2000;

export class WebSocketClient {
  private ws: WebSocket | null = null;
  private listeners: StateListener[] = [];
  private shouldReconnect = true;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;

  constructor(private url: string) {}

  connect(): void {
    this.shouldReconnect = true;
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

  onMessage(handler: (data: ArrayBuffer) => void): void {
    this._onMessage = handler;
  }

  onStateChange(listener: StateListener): void {
    this.listeners.push(listener);
  }

  get connected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  private _onMessage: ((data: ArrayBuffer) => void) | null = null;

  private open(): void {
    this.notify("connecting");
    const ws = new WebSocket(this.url);
    ws.binaryType = "arraybuffer";

    ws.onopen = () => this.notify("connected");

    ws.onmessage = (event) => {
      if (event.data instanceof ArrayBuffer) {
        this._onMessage?.(event.data);
      }
    };

    ws.onclose = () => {
      this.notify("disconnected");
      this.ws = null;
      if (this.shouldReconnect) {
        this.reconnectTimer = setTimeout(() => this.open(), RECONNECT_DELAY_MS);
      }
    };

    ws.onerror = () => ws.close();

    this.ws = ws;
  }

  private notify(state: ConnectionState): void {
    for (const listener of this.listeners) {
      listener(state);
    }
  }
}

export function buildWsUrl(): string {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return `${protocol}//${window.location.host}/ws`;
}
