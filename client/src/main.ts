import "./style.css";
import { WebSocketClient, buildWsUrl } from "./connection/websocket-client.js";
import { Heartbeat } from "./connection/heartbeat.js";
import { setupTouchpad, type AccelLevel } from "./input/touchpad.js";
import { setupButtons } from "./input/buttons.js";
import { setupStatusIndicator } from "./ui/status-indicator.js";
import { encodeMove, encodeClick } from "./proto-helpers.js";

function getElement(id: string): HTMLElement {
  const el = document.getElementById(id);
  if (!el) throw new Error(`Missing element: #${id}`);
  return el;
}

const client = new WebSocketClient(buildWsUrl());
const heartbeat = new Heartbeat(client);

const updateStatus = setupStatusIndicator(
  getElement("status-dot"),
  getElement("status-text"),
);

client.onStateChange((state) => {
  updateStatus(state);
  if (state === "connected") {
    heartbeat.start(() => client.reconnect());
  } else if (state === "disconnected") {
    heartbeat.stop();
  }
});

const touchpad = setupTouchpad(getElement("touchpad"), (dx, dy) => {
  client.send(encodeMove(dx, dy));
});

const accelSelect = getElement("accel-select") as HTMLSelectElement;
accelSelect.value = touchpad.accelLevel;
accelSelect.addEventListener("change", () => {
  touchpad.setAccelLevel(accelSelect.value as AccelLevel);
});

setupButtons(getElement("buttons"), (button) => {
  client.send(encodeClick(button));
});

client.connect();
