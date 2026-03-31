export type MoveHandler = (deltaX: number, deltaY: number) => void;
export type AccelLevel = "off" | "small" | "medium" | "large";

const ACCEL_CONFIGS: Record<AccelLevel, { threshold: number; gain: number } | null> = {
  off: null,
  small:  { threshold: 8, gain: 0.3 },
  medium: { threshold: 6, gain: 0.6 },
  large:  { threshold: 4, gain: 0.9 },
};

function applyAcceleration(dx: number, dy: number, level: AccelLevel): [number, number] {
  const config = ACCEL_CONFIGS[level];
  if (!config) return [dx, dy];
  const speed = Math.hypot(dx, dy);
  const multiplier = 1 + (speed / config.threshold) * config.gain;
  return [dx * multiplier, dy * multiplier];
}

function loadAccelLevel(): AccelLevel {
  const stored = localStorage.getItem("accel") as AccelLevel | null;
  return stored && stored in ACCEL_CONFIGS ? stored : "medium";
}

export function setupTouchpad(element: HTMLElement, onMove: MoveHandler) {
  let lastX = 0;
  let lastY = 0;
  let tracking = false;
  let accelLevel = loadAccelLevel();

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
      localStorage.setItem("accel", level);
    },
    get accelLevel() { return accelLevel; },
  };
}
