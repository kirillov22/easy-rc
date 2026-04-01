import { describe, it, expect } from "vitest";
import { applyAcceleration, isValidAccelLevel, type AccelLevel } from "../acceleration.js";

describe("applyAcceleration", () => {
  it("returns raw deltas when level is off", () => {
    expect(applyAcceleration(10, 5, "off")).toEqual([10, 5]);
  });

  it("returns raw deltas for zero movement regardless of level", () => {
    expect(applyAcceleration(0, 0, "high")).toEqual([0, 0]);
  });

  it("amplifies movement proportional to speed", () => {
    const [ax, ay] = applyAcceleration(20, 0, "high");
    expect(ax).toBeGreaterThan(20);
    expect(ay).toBe(0);
  });

  it("applies higher amplification at higher levels", () => {
    const levels: AccelLevel[] = ["low", "medium", "high"];
    const results = levels.map((level) => applyAcceleration(10, 10, level)[0]);
    for (let i = 1; i < results.length; i++) {
      expect(results[i]).toBeGreaterThan(results[i - 1]);
    }
  });

  it("preserves movement direction", () => {
    const [ax, ay] = applyAcceleration(-15, 10, "medium");
    expect(ax).toBeLessThan(0);
    expect(ay).toBeGreaterThan(0);
  });
});

describe("isValidAccelLevel", () => {
  it("accepts valid levels", () => {
    expect(isValidAccelLevel("off")).toBe(true);
    expect(isValidAccelLevel("low")).toBe(true);
    expect(isValidAccelLevel("medium")).toBe(true);
    expect(isValidAccelLevel("high")).toBe(true);
  });

  it("rejects invalid values", () => {
    expect(isValidAccelLevel("huge")).toBe(false);
    expect(isValidAccelLevel("")).toBe(false);
    expect(isValidAccelLevel("33")).toBe(false);
  });
});
