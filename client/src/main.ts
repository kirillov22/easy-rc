import "./style.css";
import { WebSocketClient, buildWsUrl } from "./connection/websocket-client.js";
import { Heartbeat } from "./connection/heartbeat.js";
import { bootstrapApp } from "./app.js";

function getElement(id: string): HTMLElement {
  const el = document.getElementById(id);
  if (!el) throw new Error(`Missing element: #${id}`);
  return el;
}

const client = new WebSocketClient(buildWsUrl());

bootstrapApp({
  elements: {
    statusDot: getElement("status-dot"),
    statusText: getElement("status-text"),
    touchpad: getElement("touchpad"),
    buttons: getElement("buttons"),
    accelSelect: getElement("accel-select") as HTMLSelectElement,
  },
  connection: client,
  heartbeat: new Heartbeat(client),
});
