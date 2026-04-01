import { applyAcceleration, isValidAccelLevel, localStorageAdapter } from "./acceleration.js";
import type { AccelLevel, Storage } from "./acceleration.js";
export type { AccelLevel } from "./acceleration.js";
export { applyAcceleration } from "./acceleration.js";

export type MoveHandler = (deltaX: number, deltaY: number) => void;

function loadAccelLevel(storage: Storage): AccelLevel {
  const stored = storage.get("accel");
  return stored && isValidAccelLevel(stored) ? stored : "medium";
}

export function setupTouchpad(
  element: HTMLElement,
  onMove: MoveHandler,
  storage: Storage = localStorageAdapter,
) {
  let lastX = 0;
  let lastY = 0;
  let tracking = false;
  let accelLevel = loadAccelLevel(storage);

  element.addEventListener("touchstart", (e) => {
    e.preventDefault();
    const touch = e.touches[0];
    lastX = touch.clientX;
    lastY = touch.clientY;
    tracking = true;
  }, { passive: false });

  element.addEventListener("touchmove", (e) => {
    e.preventDefault();
    if (!tracking) return;
    const touch = e.touches[0];
    let moveX = touch.clientX - lastX;
    let moveY = touch.clientY - lastY;
    lastX = touch.clientX;
    lastY = touch.clientY;
    if (accelLevel !== "off") {
      [moveX, moveY] = applyAcceleration(moveX, moveY, accelLevel);
    }
    onMove(moveX, moveY);
  }, { passive: false });

  const endTracking = () => { tracking = false; };

  element.addEventListener("touchend", endTracking);
  element.addEventListener("touchcancel", endTracking);

  return {
    setAccelLevel(level: AccelLevel) {
      accelLevel = level;
      storage.set("accel", level);
    },
    get accelLevel() { return accelLevel; },
  };
}
