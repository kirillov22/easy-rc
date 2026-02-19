export type MoveHandler = (deltaX: number, deltaY: number) => void;

const MAX_STEP = 10;

export function setupTouchpad(element: HTMLElement, onMove: MoveHandler): void {
  let lastX = 0;
  let lastY = 0;
  let tracking = false;
  let pendingDX = 0;
  let pendingDY = 0;
  let rafId = 0;

  function flush() {
    rafId = 0;
    const dx = pendingDX;
    const dy = pendingDY;
    pendingDX = 0;
    pendingDY = 0;
    if (dx === 0 && dy === 0) return;

    const maxDelta = Math.max(Math.abs(dx), Math.abs(dy));
    const steps = Math.max(1, Math.ceil(maxDelta / MAX_STEP));
    const stepDX = Math.round(dx / steps);
    const stepDY = Math.round(dy / steps);

    for (let i = 0; i < steps - 1; i++) {
      onMove(stepDX, stepDY);
    }
    onMove(dx - stepDX * (steps - 1), dy - stepDY * (steps - 1));
  }

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
    pendingDX += touch.clientX - lastX;
    pendingDY += touch.clientY - lastY;
    lastX = touch.clientX;
    lastY = touch.clientY;
    if (!rafId) {
      rafId = requestAnimationFrame(flush);
    }
  }, { passive: false });

  const endTracking = () => {
    tracking = false;
    if (rafId) {
      cancelAnimationFrame(rafId);
      rafId = 0;
    }
    flush();
  };

  element.addEventListener("touchend", endTracking);
  element.addEventListener("touchcancel", endTracking);
}
