import { webcrypto } from "crypto";
import { defineConfig } from "vitest/config";

if (!globalThis.crypto?.getRandomValues && webcrypto?.getRandomValues) {
  globalThis.crypto = webcrypto;
}

export default defineConfig({
  test: {
    include: ["src/features/refund/helpers.test.ts"],
    environment: "node",
    globals: true,
    setupFiles: "./vitest.setup.ts"
  }
});
