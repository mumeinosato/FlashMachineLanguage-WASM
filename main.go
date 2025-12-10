//go:build !wasm
// +build !wasm

package main

import (
	"fmt"

	"backend/checker"
	"backend/emulator"
	"backend/genhex"
)

func main() {
	cpu := emulator.NewCPU()
	fmt.Printf("Completed cpu initialization\n")
	fmt.Printf("RAX=%d, RBX=%d, RCX=%d, RDX=%d\n\n",
		cpu.GetRegister(emulator.RAX),
		cpu.GetRegister(emulator.RBX),
		cpu.GetRegister(emulator.RCX),
		cpu.GetRegister(emulator.RDX),
	)

	debugMode := true

	if debugMode {
		for i := 1; i < 5; i++ {
			fmt.Printf("=== Level %d ===\n", i)
			if err := genAndRunLevel(cpu, i); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Println("======================\n")
		}
	} else {
		var hexInput string
		fmt.Print("Enter machine code (hex): ")
		fmt.Scanln(&hexInput)

		if err := runHex(cpu, hexInput); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}
}

func assertEqualInt(expected, actual int) error {
	if expected != actual {
		return fmt.Errorf("assert equal failed: expected=%d, actual=%d", expected, actual)
	}
	return nil
}

func runHex(cpu *emulator.CPU, hex string) error {
	code, err := emulator.ParseHexString(hex)
	if err != nil {
		return fmt.Errorf("parse hex: %w", err)
	}

	decoder := emulator.NewDecoder(code)
	for decoder.HasMore() {
		inst, err := decoder.DecodeNext()
		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}
		if err := cpu.Execute(inst); err != nil {
			return fmt.Errorf("execute: %w", err)
		}
	}

	fmt.Printf("\n=== Execution Result ===\n")
	fmt.Printf("RAX=%d, RBX=%d, RCX=%d, RDX=%d\n",
		cpu.GetRegister(emulator.RAX),
		cpu.GetRegister(emulator.RBX),
		cpu.GetRegister(emulator.RCX),
		cpu.GetRegister(emulator.RDX),
	)
	fmt.Printf("Final result (int32): %d\n", cpu.GetResult())
	fmt.Printf("Final result (int32): %x\n", int32(cpu.GetResult()))
	return nil
}

func genAndRunLevel(cpu *emulator.CPU, level int) error {
	spaceHex, noSpaceHex, err := genhex.GenerateHex(level)
	if err != nil {
		return fmt.Errorf("GenerateHex: %w", err)
	}
	fmt.Println("Generated hex: " + spaceHex)

	check, err := checker.CheckLevel(noSpaceHex)
	if err != nil {
		return fmt.Errorf("CheckLevel: %w", err)
	}
	fmt.Printf("Checker returned: %d\n", check)

	if err := assertEqualInt(level, check); err != nil {
		return err
	}

	return runHex(cpu, noSpaceHex)
}
