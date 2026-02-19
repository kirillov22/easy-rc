import { Message, MouseButton } from "./generated/messages.js";

export function encodeMove(deltaX: number, deltaY: number): Uint8Array {
  const msg = Message.create({
    move: {
      timestamp: Date.now(),
      moveX: Math.round(deltaX),
      moveY: Math.round(deltaY),
    },
  });
  return Message.encode(msg).finish();
}

export function encodeClick(button: "left" | "middle" | "right"): Uint8Array {
  const buttonMap = {
    left: MouseButton.MOUSE_BUTTON_LEFT,
    middle: MouseButton.MOUSE_BUTTON_MIDDLE,
    right: MouseButton.MOUSE_BUTTON_RIGHT,
  } as const;

  const msg = Message.create({
    click: {
      timestamp: Date.now(),
      mouseButton: buttonMap[button],
    },
  });
  return Message.encode(msg).finish();
}

export function encodePing(): Uint8Array {
  const msg = Message.create({
    ping: { timestamp: Date.now() },
  });
  return Message.encode(msg).finish();
}

export function decode(data: ArrayBuffer): Message {
  return Message.decode(new Uint8Array(data));
}
