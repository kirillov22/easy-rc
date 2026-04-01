// @vitest-environment jsdom
import { describe, it, expect } from "vitest";
import { bootstrapApp, type AppDeps } from "../app.js";
import { Heartbeat } from "../connection/heartbeat.js";
import type { Connection, MessageListener, StateListener } from "../connection/types.js";
import type { Storage } from "../input/acceleration.js";

class FakeConnection implements Connection {
  connected = false;
  sent: Uint8Array[] = [];
  private messageHandlers: MessageListener[] = [];
  private stateListeners: StateListener[] = [];
  connectCalled = false;

  send(data: Uint8Array) { this.sent.push(data); }

  onMessage(handler: MessageListener): () => void {
    this.messageHandlers.push(handler);
    return () => {
      this.messageHandlers = this.messageHandlers.filter((h) => h !== handler);
    };
  }

  onStateChange(listener: StateListener) {
    this.stateListeners.push(listener);
  }

  connect() { this.connectCalled = true; }
  reconnect() {}

  simulateState(state: "connecting" | "connected" | "disconnected") {
    if (state === "connected") this.connected = true;
    if (state === "disconnected") this.connected = false;
    for (const l of this.stateListeners) l(state);
  }
}

function createFakeStorage(initial: Record<string, string> = {}): Storage {
  const data = new Map(Object.entries(initial));
  return {
    get: (key) => data.get(key) ?? null,
    set: (key, value) => data.set(key, value),
  };
}

function createElements() {
  return {
    statusDot: document.createElement("span"),
    statusText: document.createElement("span"),
    touchpad: document.createElement("div"),
    buttons: document.createElement("div"),
    accelSelect: document.createElement("select") as HTMLSelectElement,
  };
}

function setup(overrides: Partial<AppDeps> = {}) {
  const conn = new FakeConnection();
  const storage = createFakeStorage({ accel: "medium" });
  const elements = createElements();
  const heartbeat = new Heartbeat(conn);

  const deps: AppDeps = {
    elements,
    connection: conn,
    heartbeat,
    storage,
    ...overrides,
  };

  bootstrapApp(deps);
  return { conn, storage, elements, heartbeat };
}

describe("bootstrapApp", () => {
  it("calls connect on the connection", () => {
    const { conn } = setup();
    expect(conn.connectCalled).toBe(true);
  });

  it("updates status text on state change", () => {
    const { conn, elements } = setup();
    conn.simulateState("connected");
    expect(elements.statusText.textContent).toBe("Connected");

    conn.simulateState("disconnected");
    expect(elements.statusText.textContent).toBe("Disconnected");
  });

  it("toggles status dot class on state change", () => {
    const { conn, elements } = setup();
    conn.simulateState("connected");
    expect(elements.statusDot.classList.contains("connected")).toBe(true);

    conn.simulateState("disconnected");
    expect(elements.statusDot.classList.contains("connected")).toBe(false);
  });

  it("sends encoded click when button is pressed", () => {
    const { conn, elements } = setup();

    const button = document.createElement("button");
    button.dataset.button = "left";
    elements.buttons.appendChild(button);

    button.dispatchEvent(new PointerEvent("pointerdown", { bubbles: true }));

    expect(conn.sent.length).toBe(1);
    expect(conn.sent[0].length).toBeGreaterThan(0);
  });

  it("loads accel level from storage", () => {
    const storage = createFakeStorage({ accel: "large" });
    const conn = new FakeConnection();
    const elements = createElements();
    for (const val of ["off", "small", "medium", "large"]) {
      const opt = document.createElement("option");
      opt.value = val;
      elements.accelSelect.appendChild(opt);
    }
    const heartbeat = new Heartbeat(conn);

    bootstrapApp({ elements, connection: conn, heartbeat, storage });

    expect(elements.accelSelect.value).toBe("large");
  });
});
