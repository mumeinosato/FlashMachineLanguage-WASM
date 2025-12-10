//go:build wasm
// +build wasm

package main

import (
	"backend/emulator"
	"fmt"
	"syscall/js"
)

func main() {
	js.Global().Set("RunCode", js.FuncOf(run))
	js.Global().Set("GenHex", js.FuncOf(genMachineLanguage))

	select {}
}

func run(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{
			"error": "hex string required",
		}
	}

	hexInput := args[0].String()

	cpu := emulator.NewCPU()
	code, err := emulator.ParseHexString(hexInput)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("error parsing hex input: %v", err),
		}
	}

	decoder := emulator.NewDecoder(code)
	for decoder.HasMore() {
		inst, err := decoder.DecodeNext()
		if err != nil {
			return map[string]interface{}{
				"error": fmt.Sprintf("decode error: %v", err),
			}
		}

		if err := cpu.Execute(inst); err != nil {
			return map[string]interface{}{
				"error": fmt.Sprintf("execute error: %v", err),
			}
		}
	}

	return fmt.Sprintf("%x", int32(cpu.GetResult()))
}

func genMachineLanguage(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{"error": "int argument required"}
	}

	level := args[0].Int()

	if level < 1 || level > 4 {
		return map[string]interface{}{"error": "invalid level"}
	}

	spaceHex, noSpaceHex, err := genhex.GenerateHex(level)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("error generating hex: %v", err)}
	}

	return []interface{}{spaceHex, noSpaceHex}
}
