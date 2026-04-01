export type AccelLevel = "off" | "low" | "medium" | "high";

export const ACCEL_CONFIGS: Record<AccelLevel, { threshold: number; gain: number } | null> = {
  off: null,
  low:  { threshold: 8, gain: 0.3 },
  medium: { threshold: 6, gain: 0.6 },
  high:  { threshold: 4, gain: 0.9 },
};

export function applyAcceleration(dx: number, dy: number, level: AccelLevel): [number, number] {
  const config = ACCEL_CONFIGS[level];
  if (!config) return [dx, dy];
  const speed = Math.hypot(dx, dy);
  const multiplier = 1 + (speed / config.threshold) * config.gain;
  return [dx * multiplier, dy * multiplier];
}

export function isValidAccelLevel(value: string): value is AccelLevel {
  return value in ACCEL_CONFIGS;
}

export interface Storage {
  get(key: string): string | null;
  set(key: string, value: string): void;
}

export const localStorageAdapter: Storage = {
  get: (key) => localStorage.getItem(key),
  set: (key, value) => localStorage.setItem(key, value),
};
