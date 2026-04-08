import { webcrypto as nodeWebcrypto } from "crypto";

if (!globalThis.crypto?.getRandomValues) {
  if (nodeWebcrypto?.getRandomValues) {
    globalThis.crypto = nodeWebcrypto;
  } else if ((globalThis as any).crypto && (globalThis as any).crypto.getRandomValues) {
    // already has it
  } else {
    // Fallback to simple Math-based shim
    (globalThis as any).crypto = {
      getRandomValues(buffer: Uint8Array) {
        for (let i = 0; i < buffer.length; i += 1) {
          buffer[i] = Math.floor(Math.random() * 256);
        }
        return buffer;
      }
    };
  }
}
