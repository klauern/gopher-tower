import "@testing-library/jest-dom";
import { cleanup } from "@testing-library/react";
import { afterEach, beforeEach, vi } from "vitest";

// Automatically cleanup after each test
afterEach(() => {
  cleanup();
});

// Mock fetch globally
global.fetch = vi.fn();

// Mock IntersectionObserver
class MockIntersectionObserver implements IntersectionObserver {
  readonly root: Element | null = null;
  readonly rootMargin: string = "0px";
  readonly thresholds: ReadonlyArray<number> = [0];

  observe = vi.fn();
  unobserve = vi.fn();
  disconnect = vi.fn();
  takeRecords = vi.fn().mockReturnValue([]);
}

global.IntersectionObserver =
  MockIntersectionObserver as unknown as typeof IntersectionObserver;

// Mock the Radix UI Select components
vi.mock("@radix-ui/react-select", () => {
  return import("../test/mocks/select");
});

// Mock next/navigation
vi.mock("next/navigation", () => ({
  useRouter: () => ({
    push: vi.fn(),
    back: vi.fn(),
    forward: vi.fn(),
  }),
  usePathname: vi.fn(),
}));

// Reset all mocks before each test
beforeEach(() => {
  vi.clearAllMocks();
});
