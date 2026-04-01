import { describe, it, expect } from "vitest";
import { Heartbeat } from "../heartbeat.js";
import type { Connection, Clock, MessageListener } from "../types.js";
import { Message } from "../../generated/messages.js";

class FakeConnection implements Connection {
  connected = true;
  sent: Uint8Array[] = [];
  private messageHandlers: MessageListener[] = [];

  send(data: Uint8Array) { this.sent.push(data); }

  onMessage(handler: MessageListener): () => void {
    this.messageHandlers.push(handler);
    return () => {
      this.messageHandlers = this.messageHandlers.filter((h) => h !== handler);
    };
  }

  onStateChange() {}
  connect() {}
  reconnect() {}

  simulateMessage(data: ArrayBuffer) {
    for (const h of this.messageHandlers) h(data);
  }
}

class FakeClock implements Clock {
  current = 0;
  private callbacks = new Map<number, { fn: () => void; interval: number; next: number }>();
  private nextId = 1;

  now() { return this.current; }

  setInterval(fn: () => void, ms: number) {
    const id = this.nextId++;
    this.callbacks.set(id, { fn, interval: ms, next: this.current + ms });
    return id as unknown as ReturnType<typeof setInterval>;
  }

  clearInterval(id: ReturnType<typeof setInterval>) {
    this.callbacks.delete(id as unknown as number);
  }

  advance(ms: number) {
    this.current += ms;
    for (const [, cb] of this.callbacks) {
      while (cb.next <= this.current) {
        cb.fn();
        cb.next += cb.interval;
      }
    }
  }
}

function encodePong(): ArrayBuffer {
  const msg = Message.create({ pong: { timestamp: Date.now() } });
  const bytes = Message.encode(msg).finish();
  return bytes.buffer.slice(bytes.byteOffset, bytes.byteOffset + bytes.byteLength) as ArrayBuffer;
}

describe("Heartbeat", () => {
  it("sends pings at regular intervals when connected", () => {
    const conn = new FakeConnection();
    const clock = new FakeClock();
    const heartbeat = new Heartbeat(conn, clock);

    heartbeat.start(() => {});
    clock.advance(5000);

    expect(conn.sent.length).toBe(1);

    clock.advance(5000);
    expect(conn.sent.length).toBe(2);
  });

  it("does not send pings when disconnected", () => {
    const conn = new FakeConnection();
    conn.connected = false;
    const clock = new FakeClock();
    const heartbeat = new Heartbeat(conn, clock);

    heartbeat.start(() => {});
    clock.advance(5000);

    expect(conn.sent.length).toBe(0);
  });

  it("triggers timeout when no pong received", () => {
    const conn = new FakeConnection();
    const clock = new FakeClock();
    const heartbeat = new Heartbeat(conn, clock);
    let timedOut = false;

    heartbeat.start(() => { timedOut = true; });

    clock.advance(12000);
    expect(timedOut).toBe(true);
  });

  it("does not timeout when pong is received in time", () => {
    const conn = new FakeConnection();
    const clock = new FakeClock();
    const heartbeat = new Heartbeat(conn, clock);
    let timedOut = false;

    heartbeat.start(() => { timedOut = true; });

    clock.advance(4000);
    conn.simulateMessage(encodePong());
    clock.advance(8000);

    expect(timedOut).toBe(false);
  });

  it("stops sending pings after stop()", () => {
    const conn = new FakeConnection();
    const clock = new FakeClock();
    const heartbeat = new Heartbeat(conn, clock);

    heartbeat.start(() => {});
    clock.advance(5000);
    expect(conn.sent.length).toBe(1);

    heartbeat.stop();
    clock.advance(10000);
    expect(conn.sent.length).toBe(1);
  });

  it("cleans up previous timers when start() is called twice", () => {
    const conn = new FakeConnection();
    const clock = new FakeClock();
    const heartbeat = new Heartbeat(conn, clock);
    let firstCallbackCalled = false;

    heartbeat.start(() => { firstCallbackCalled = true; });
    heartbeat.start(() => {});

    clock.advance(12000);
    expect(firstCallbackCalled).toBe(false);
  });
});
