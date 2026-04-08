const nodeCrypto = require("crypto");
const { webcrypto } = nodeCrypto;

if (!nodeCrypto.getRandomValues && webcrypto?.getRandomValues) {
  nodeCrypto.getRandomValues = webcrypto.getRandomValues.bind(webcrypto);
}

if (!globalThis.crypto?.getRandomValues && webcrypto?.getRandomValues) {
  globalThis.crypto = webcrypto;
} else if (!globalThis.crypto) {
  globalThis.crypto = {
    getRandomValues(buffer) {
      for (let i = 0; i < buffer.length; i += 1) {
        buffer[i] = Math.floor(Math.random() * 256);
      }
      return buffer;
    }
  };
}
