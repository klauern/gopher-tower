import "@testing-library/jest-dom";
import { cleanup } from "@testing-library/react";
import { afterEach } from "vitest";

// Configure testing-library
afterEach(cleanup);

// Add any global test setup here
declare global {
  interface Window {
    IS_REACT_ACT_ENVIRONMENT: boolean;
  }
}

window.IS_REACT_ACT_ENVIRONMENT = true;
