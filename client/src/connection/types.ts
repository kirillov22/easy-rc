export type ConnectionState = "connecting" | "connected" | "disconnected";
export type StateListener = (state: ConnectionState) => void;
export type MessageListener = (data: ArrayBuffer) => void;

export interface Connection {
  readonly connected: boolean;
  send(data: Uint8Array): void;
  onMessage(handler: MessageListener): () => void;
  onStateChange(listener: StateListener): void;
  connect(): void;
  reconnect(): void;
}

export interface Clock {
  now(): number;
  setInterval(fn: () => void, ms: number): ReturnType<typeof setInterval>;
  clearInterval(id: ReturnType<typeof setInterval>): void;
}

export const systemClock: Clock = {
  now: () => Date.now(),
  setInterval: (fn, ms) => setInterval(fn, ms),
  clearInterval: (id) => clearInterval(id),
};

export interface Environment {
  isVisible(): boolean;
  onVisibilityChange(handler: () => void): () => void;
  createWebSocket(url: string): WebSocket;
  setTimeout(fn: () => void, ms: number): ReturnType<typeof setTimeout>;
  clearTimeout(id: ReturnType<typeof setTimeout>): void;
}

export const browserEnvironment: Environment = {
  isVisible: () => document.visibilityState === "visible",
  onVisibilityChange(handler) {
    document.addEventListener("visibilitychange", handler);
    return () => document.removeEventListener("visibilitychange", handler);
  },
  createWebSocket: (url) => new WebSocket(url),
  setTimeout: (fn, ms) => setTimeout(fn, ms),
  clearTimeout: (id) => clearTimeout(id),
};
