export type ClickHandler = (button: "left" | "middle" | "right") => void;

export function setupButtons(container: HTMLElement, onClick: ClickHandler): void {
  container.addEventListener("pointerdown", (e) => {
    const target = e.target as HTMLElement;
    const button = target.dataset.button as "left" | "middle" | "right" | undefined;
    if (button) {
      onClick(button);
    }
  });
}
