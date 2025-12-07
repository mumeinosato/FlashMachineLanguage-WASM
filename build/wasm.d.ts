declare global {
  interface Window {
    /**
     * Run WASM code with hex string input
     * @param hexInput - Hexadecimal machine code string
     * @returns Result as hex string or error object
     */
    RunCode(hexInput: string): string | { error: string };
  }
}

export {};
