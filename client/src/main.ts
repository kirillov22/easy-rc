import "./style.css";
import { WebSocketClient, buildWsUrl } from "./connection/websocket-client.js";
import { Heartbeat } from "./connection/heartbeat.js";
import { setupTouchpad } from "./input/touchpad.js";
import { setupButtons } from "./input/buttons.js";
import { setupStatusIndicator } from "./ui/status-indicator.js";
import { encodeMove, encodeClick } from "./proto-helpers.js";

const client = new WebSocketClient(buildWsUrl());
const heartbeat = new Heartbeat(client);

const updateStatus = setupStatusIndicator(
  document.getElementById("status-dot")!,
  document.getElementById("status-text")!,
);

client.onStateChange((state) => {
  updateStatus(state);
  // if (state === "connected") {
  //   heartbeat.start(() => client.disconnect());
  // } else if (state === "disconnected") {
  //   heartbeat.stop();
  // }
});

setupTouchpad(document.getElementById("touchpad")!, (dx, dy) => {
  client.send(encodeMove(dx, dy));
});

setupButtons(document.getElementById("buttons")!, (button) => {
  client.send(encodeClick(button));
});

client.connect();
