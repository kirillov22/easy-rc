import type { ConnectionState } from "../connection/types.js";

export function setupStatusIndicator(
  dot: HTMLElement,
  text: HTMLElement,
): (state: ConnectionState) => void {
  const labels: Record<ConnectionState, string> = {
    connecting: "Connecting...",
    connected: "Connected",
    disconnected: "Disconnected",
  };

  return (state) => {
    dot.classList.toggle("connected", state === "connected");
    text.textContent = labels[state];
  };
}
