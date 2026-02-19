export type MoveHandler = (deltaX: number, deltaY: number) => void;

export function setupTouchpad(element: HTMLElement, onMove: MoveHandler): void {
  let lastX = 0;
  let lastY = 0;
  let tracking = false;

  element.addEventListener(
    "touchstart",
    (e) => {
      e.preventDefault();
      const touch = e.touches[0];
      lastX = touch.clientX;
      lastY = touch.clientY;
      tracking = true;
    },
    { passive: false },
  );

  element.addEventListener(
    "touchmove",
    (e) => {
      e.preventDefault();
      if (!tracking) return;
      const touch = e.touches[0];
      const deltaX = touch.clientX - lastX;
      const deltaY = touch.clientY - lastY;
      lastX = touch.clientX;
      lastY = touch.clientY;
      if (deltaX !== 0 || deltaY !== 0) {
        onMove(deltaX, deltaY);
      }
    },
    { passive: false },
  );

  const endTracking = () => {
    tracking = false;
  };

  element.addEventListener("touchend", endTracking);
  element.addEventListener("touchcancel", endTracking);
}
