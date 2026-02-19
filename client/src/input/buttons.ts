export type ClickHandler = (button: "left" | "middle" | "right") => void;

type ButtonName = "left" | "middle" | "right";
const VALID_BUTTONS: Set<string> = new Set(["left", "middle", "right"]);

export function setupButtons(container: HTMLElement, onClick: ClickHandler): void {
  container.addEventListener("pointerdown", (e) => {
    const target = e.target as HTMLElement;
    const raw = target.dataset.button;
    if (raw && VALID_BUTTONS.has(raw)) {
      onClick(raw as ButtonName);
    }
  });
}
