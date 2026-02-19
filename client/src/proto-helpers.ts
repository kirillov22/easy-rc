import * as $protobuf from "protobufjs/minimal";
import { Message, MouseButton } from "./generated/messages.js";

const moveMsg = Message.create({ move: { timestamp: 0, moveX: 0, moveY: 0 } });
const moveWriter = $protobuf.Writer.create();

export function encodeMove(deltaX: number, deltaY: number): Uint8Array {
  moveMsg.move!.timestamp = Date.now();
  moveMsg.move!.moveX = Math.round(deltaX);
  moveMsg.move!.moveY = Math.round(deltaY);
  moveWriter.reset();
  return Message.encode(moveMsg, moveWriter).finish();
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
