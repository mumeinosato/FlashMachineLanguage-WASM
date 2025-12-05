package main

import (
	"backend/emulator"
	"fmt"
)

func main() {
	cpu := emulator.NewCPU()
	fmt.Printf("Completed cpu initialization\n")
	fmt.Printf("RAX=%d, RBX=%d, RCX=%d, RDX=%d\n\n", cpu.GetRegister(emulator.RAX), cpu.GetRegister(emulator.RBX), cpu.GetRegister(emulator.RCX), cpu.GetRegister(emulator.RDX))

	var hexInput string
	fmt.Print("Enter machine code (hex): ")
	fmt.Scanln(&hexInput)

	code, err := emulator.ParseHexString(hexInput)
	if err != nil {
		fmt.Printf("Error parsing hex input: %v\n", err)
		return
	}

	decoder := emulator.NewDecoder(code)
	for decoder.HasMore() {
		inst, err := decoder.DecodeNext()
		if err != nil {
			fmt.Printf("Decode Error: %v\n", err)
			break
		}

		if err := cpu.Execute(inst); err != nil {
			fmt.Printf("Execute Error: %v\n", err)
			break
		}
	}

	fmt.Printf("\n=== Execution Result ===\n")
	fmt.Printf("RAX=%d, RBX=%d, RCX=%d, RDX=%d\n", cpu.GetRegister(emulator.RAX), cpu.GetRegister(emulator.RBX), cpu.GetRegister(emulator.RCX), cpu.GetRegister(emulator.RDX))
	fmt.Printf("Final result (int32): %d\n", cpu.GetResult())
}
