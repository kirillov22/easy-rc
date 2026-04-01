import type { Connection } from "./connection/types.js";
import { Heartbeat } from "./connection/heartbeat.js";
import { setupTouchpad } from "./input/touchpad.js";
import type { AccelLevel, Storage } from "./input/acceleration.js";
import { setupButtons } from "./input/buttons.js";
import { setupStatusIndicator } from "./ui/status-indicator.js";
import { encodeMove, encodeClick } from "./proto-helpers.js";

export interface AppElements {
  statusDot: HTMLElement;
  statusText: HTMLElement;
  touchpad: HTMLElement;
  buttons: HTMLElement;
  accelSelect: HTMLSelectElement;
}

export interface AppDeps {
  elements: AppElements;
  connection: Connection;
  heartbeat: Heartbeat;
  storage?: Storage;
}

export function bootstrapApp(deps: AppDeps): void {
  const { elements, connection, heartbeat } = deps;

  const updateStatus = setupStatusIndicator(elements.statusDot, elements.statusText);

  connection.onStateChange((state) => {
    updateStatus(state);
    if (state === "connected") {
      heartbeat.start(() => connection.reconnect());
    } else if (state === "disconnected") {
      heartbeat.stop();
    }
  });

  const touchpad = setupTouchpad(
    elements.touchpad,
    (dx, dy) => connection.send(encodeMove(dx, dy)),
    deps.storage,
  );

  elements.accelSelect.value = touchpad.accelLevel;
  elements.accelSelect.addEventListener("change", () => {
    touchpad.setAccelLevel(elements.accelSelect.value as AccelLevel);
  });

  setupButtons(elements.buttons, (button) => {
    connection.send(encodeClick(button));
  });

  connection.connect();
}
