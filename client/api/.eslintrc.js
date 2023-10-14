module.exports = {
  parserOptions: {
    ecmaVersion: 13,
    sourceType: 'module'
  },
  globals: {
    Blob: 'readonly',
    clearInterval: 'readonly',
    clearTimeout: 'readonly',
    console: 'readonly',
    crypto: 'readonly',
    File: 'readonly',
    FormData: 'readonly',
    Promise: 'readonly',
    setInterval: 'readonly',
    setTimeout: 'readonly',
    TextDecoder: 'readonly',
    Uint8Array: 'readonly',
    WebSocket: 'readonly',
    window: 'readonly'
  },
  extends: [
    "eslint:recommended",
    "prettier",
  ],
};
