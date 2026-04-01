import { describe, it, expect } from "vitest";
import { WebSocketClient } from "../websocket-client.js";
import type { Environment, ConnectionState } from "../types.js";

class FakeWebSocket {
  binaryType = "";
  readyState = 0; // CONNECTING
  onopen: (() => void) | null = null;
  onclose: (() => void) | null = null;
  onmessage: ((event: { data: unknown }) => void) | null = null;
  onerror: ((event: unknown) => void) | null = null;
  sent: unknown[] = [];
  closed = false;

  static OPEN = 1;
  static CONNECTING = 0;
  static CLOSING = 2;
  static CLOSED = 3;

  constructor(public url: string) {}

  send(data: unknown) { this.sent.push(data); }

  close() {
    this.closed = true;
    this.readyState = FakeWebSocket.CLOSED;
    this.onclose?.();
  }

  simulateOpen() {
    this.readyState = FakeWebSocket.OPEN;
    this.onopen?.();
  }

  simulateMessage(data: ArrayBuffer) {
    this.onmessage?.({ data });
  }
}

class FakeEnvironment implements Environment {
  visible = true;
  private visibilityHandlers: (() => void)[] = [];
  websockets: FakeWebSocket[] = [];
  private pendingTimeouts = new Map<number, { fn: () => void; ms: number }>();
  private nextTimeoutId = 1;

  isVisible() { return this.visible; }

  onVisibilityChange(handler: () => void): () => void {
    this.visibilityHandlers.push(handler);
    return () => {
      this.visibilityHandlers = this.visibilityHandlers.filter((h) => h !== handler);
    };
  }

  createWebSocket(url: string): WebSocket {
    const ws = new FakeWebSocket(url);
    this.websockets.push(ws);
    return ws as unknown as WebSocket;
  }

  setTimeout(fn: () => void, ms: number): ReturnType<typeof setTimeout> {
    const id = this.nextTimeoutId++;
    this.pendingTimeouts.set(id, { fn, ms });
    return id as unknown as ReturnType<typeof setTimeout>;
  }

  clearTimeout(id: ReturnType<typeof setTimeout>) {
    this.pendingTimeouts.delete(id as unknown as number);
  }

  triggerVisibilityChange() {
    for (const h of this.visibilityHandlers) h();
  }

  flushTimeouts() {
    for (const [id, { fn }] of this.pendingTimeouts) {
      this.pendingTimeouts.delete(id);
      fn();
    }
  }

  get latestWebSocket(): FakeWebSocket {
    return this.websockets[this.websockets.length - 1];
  }
}

describe("WebSocketClient", () => {
  it("creates a websocket and notifies connecting on connect()", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);
    const states: ConnectionState[] = [];
    client.onStateChange((s) => states.push(s));

    client.connect();

    expect(env.websockets.length).toBe(1);
    expect(env.latestWebSocket.url).toBe("ws://test/ws");
    expect(states).toEqual(["connecting"]);
  });

  it("notifies connected when websocket opens", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);
    const states: ConnectionState[] = [];
    client.onStateChange((s) => states.push(s));

    client.connect();
    env.latestWebSocket.simulateOpen();

    expect(states).toEqual(["connecting", "connected"]);
  });

  it("notifies disconnected when websocket closes", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);
    const states: ConnectionState[] = [];
    client.onStateChange((s) => states.push(s));

    client.connect();
    env.latestWebSocket.simulateOpen();
    env.latestWebSocket.close();

    expect(states).toEqual(["connecting", "connected", "disconnected"]);
  });

  it("schedules reconnect after close when visible", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);

    client.connect();
    env.latestWebSocket.simulateOpen();
    env.latestWebSocket.close();

    expect(env.websockets.length).toBe(1);

    env.flushTimeouts();
    expect(env.websockets.length).toBe(2);
  });

  it("does not schedule reconnect when tab is hidden", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);

    client.connect();
    env.latestWebSocket.simulateOpen();

    env.visible = false;
    env.latestWebSocket.close();

    env.flushTimeouts();
    expect(env.websockets.length).toBe(1);
  });

  it("reconnects on visibility change when disconnected", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);

    client.connect();
    env.latestWebSocket.simulateOpen();

    env.visible = false;
    env.triggerVisibilityChange();
    env.latestWebSocket.close();
    expect(env.websockets.length).toBe(1);

    env.visible = true;
    env.triggerVisibilityChange();
    expect(env.websockets.length).toBe(2);
  });

  it("reconnects on visibility change even when ws appears connected (stale connection)", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);

    client.connect();
    env.latestWebSocket.simulateOpen();

    // Page goes hidden (browser backgrounded), then comes back
    // The ws still has readyState OPEN (stale) because onclose never fired
    env.visible = false;
    env.triggerVisibilityChange();
    expect(env.websockets.length).toBe(1);

    env.visible = true;
    env.triggerVisibilityChange();
    // Should reconnect even though old ws appeared connected
    expect(env.websockets.length).toBe(2);
  });

  it("does not reconnect on visibility change if page was never hidden", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);

    client.connect();
    env.latestWebSocket.simulateOpen();
    expect(env.websockets.length).toBe(1);

    // Visibility change fires while visible (e.g. page load event)
    // Page was never hidden, so should NOT disrupt the active connection
    env.triggerVisibilityChange();
    expect(env.websockets.length).toBe(1);
  });

  it("does not fire stale onclose after open() replaces the connection", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);
    const states: ConnectionState[] = [];
    client.onStateChange((s) => states.push(s));

    client.connect();
    const firstWs = env.latestWebSocket;
    firstWs.simulateOpen();

    // Simulate hidden→visible to trigger reconnection
    env.visible = false;
    env.triggerVisibilityChange();
    env.visible = true;
    env.triggerVisibilityChange();

    const secondWs = env.latestWebSocket;
    secondWs.simulateOpen();

    // Old ws's onclose was detached, so closing it should not affect state
    expect(firstWs.onclose).toBeNull();
    expect(secondWs).not.toBe(firstWs);
  });

  it("reconnects after simulated browser close cycle (hidden → ws dies → visible)", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);
    const states: ConnectionState[] = [];
    client.onStateChange((s) => states.push(s));

    client.connect();
    env.latestWebSocket.simulateOpen();
    expect(states).toEqual(["connecting", "connected"]);

    // Browser goes to background
    env.visible = false;
    env.triggerVisibilityChange();

    // Connection drops while hidden
    env.latestWebSocket.close();
    expect(states).toEqual(["connecting", "connected", "disconnected"]);
    // No reconnect timer since hidden
    expect(env.websockets.length).toBe(1);

    // User reopens browser
    env.visible = true;
    env.triggerVisibilityChange();

    // Should have created a new websocket
    expect(env.websockets.length).toBe(2);
    expect(states).toEqual(["connecting", "connected", "disconnected", "connecting"]);

    // New connection succeeds
    env.latestWebSocket.simulateOpen();
    expect(states).toEqual(["connecting", "connected", "disconnected", "connecting", "connected"]);
  });

  it("delivers messages to listeners", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);
    const received: ArrayBuffer[] = [];

    client.onMessage((data) => received.push(data));
    client.connect();
    env.latestWebSocket.simulateOpen();

    const buf = new ArrayBuffer(4);
    env.latestWebSocket.simulateMessage(buf);

    expect(received.length).toBe(1);
    expect(received[0]).toBe(buf);
  });

  it("unsubscribes message listeners", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);
    const received: ArrayBuffer[] = [];

    const unsub = client.onMessage((data) => received.push(data));
    client.connect();
    env.latestWebSocket.simulateOpen();

    unsub();
    env.latestWebSocket.simulateMessage(new ArrayBuffer(4));

    expect(received.length).toBe(0);
  });

  it("reconnect() closes existing connection and resets delay", () => {
    const env = new FakeEnvironment();
    const client = new WebSocketClient("ws://test/ws", env);
    const states: ConnectionState[] = [];
    client.onStateChange((s) => states.push(s));

    client.connect();
    env.latestWebSocket.simulateOpen();
    client.reconnect();

    expect(env.latestWebSocket.closed).toBe(true);
    expect(states).toContain("disconnected");
  });
});
