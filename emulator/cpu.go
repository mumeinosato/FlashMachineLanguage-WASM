package emulator

import "fmt"

type Register int

const (
	RAX Register = iota
	RBX
	RCX
	RDX
)

type CPU struct {
	registers [4]int64
	pc        int
}

func (r Register) String() string {
	switch r {
	case RAX:
		return "RAX"
	case RBX:
		return "RBX"
	case RCX:
		return "RCX"
	case RDX:
		return "RDX"
	default:
		return "unknown"
	}
}

func NewCPU() *CPU {
	return &CPU{
		registers: [4]int64{0, 0, 0, 0},
		pc:        0,
	}
}

func (cpu *CPU) GetRegister(reg Register) int64 {
	return cpu.registers[reg]
}

func (cpu *CPU) SetRegister(reg Register, value int64) {
	cpu.registers[reg] = value
}

func (cpu *CPU) AddOverflowCheck(a, b int64) (int64, bool) {
	result := a + b
	overflow := (a > 0 && b > 0 && result < 0) || (a < 0 && b < 0 && result > 0)
	return result, overflow
}

func (cpu *CPU) SubOverflowCheck(a, b int64) (int64, bool) {
	result := a - b
	overflow := (a > 0 && b < 0 && result < 0) || (a < 0 && b > 0 && result > 0)
	return result, overflow
}

func (cpu *CPU) GetResult() int32 {
	rax := cpu.GetRegister(RAX)
	if rax > int64(0x7FFFFFFF) || rax < int64(-0x80000000) {
		return -1
	}
	return int32(rax)
}

func (cpu *CPU) Execute(inst *Instruction) error {
	switch inst.Opcode {
	case 0xB8, 0xB9, 0xBA, 0xBB:
		regCode := inst.Opcode - 0xB8
		reg, _ := GetRegFromModRM(regCode<<3, false)
		if inst.HasImm64 {
			cpu.SetRegister(reg, inst.Imm64)
		} else {
			cpu.SetRegister(reg, int64(inst.Imm32))
		}

	case 0x89:
		src, err := GetRegFromModRM(inst.ModRM, false)
		if err != nil {
			return err
		}
		dst, err := GetRegFromModRM(inst.ModRM, true)
		if err != nil {
			return err
		}
		cpu.SetRegister(dst, cpu.GetRegister(src))

	case 0x8B:
		src, err := GetRegFromModRM(inst.ModRM, true)
		if err != nil {
			return err
		}
		dst, err := GetRegFromModRM(inst.ModRM, false)
		if err != nil {
			return err
		}
		cpu.SetRegister(dst, cpu.GetRegister(src))

	case 0x01:
		src, err := GetRegFromModRM(inst.ModRM, false)
		if err != nil {
			return err
		}
		dst, err := GetRegFromModRM(inst.ModRM, true)
		if err != nil {
			return err
		}
		result, overflow := cpu.AddOverflowCheck(cpu.GetRegister(dst), cpu.GetRegister(src))
		if overflow {
			return &EmulatorError{PC: cpu.pc, Message: "overflow detected in ADD"}
		}
		cpu.SetRegister(dst, result)

	case 0x29:
		src, err := GetRegFromModRM(inst.ModRM, false)
		if err != nil {
			return err
		}
		dst, err := GetRegFromModRM(inst.ModRM, true)
		if err != nil {
			return err
		}
		result, overflow := cpu.SubOverflowCheck(cpu.GetRegister(dst), cpu.GetRegister(src))
		if overflow {
			return &EmulatorError{PC: cpu.pc, Message: "overflow detected in SUB"}
		}
		cpu.SetRegister(dst, result)

	case 0x81:
		subOpcode := (inst.ModRM >> 3) & 0x07
		dst, err := GetRegFromModRM(inst.ModRM, true)
		if err != nil {
			return err
		}
		switch subOpcode {
		case 0:
			result, overflow := cpu.AddOverflowCheck(cpu.GetRegister(dst), int64(inst.Imm32))
			if overflow {
				return &EmulatorError{PC: cpu.pc, Message: "overflow detected in ADD"}
			}
			cpu.SetRegister(dst, result)
		case 5:
			result, overflow := cpu.SubOverflowCheck(cpu.GetRegister(dst), int64(inst.Imm32))
			if overflow {
				return &EmulatorError{PC: cpu.pc, Message: "overflow detected in SUB"}
			}
			cpu.SetRegister(dst, result)
		default:
			return fmt.Errorf("unsupported 0x81 subopcode: %d", subOpcode)
		}

	case 0xC7:
		subOpcode := (inst.ModRM >> 3) & 0x07
		if subOpcode != 0 {
			return fmt.Errorf("unsupported 0xC7 subopcode: %d", subOpcode)
		}
		dst, err := GetRegFromModRM(inst.ModRM, true)
		if err != nil {
			return err
		}
		cpu.SetRegister(dst, int64(inst.Imm32))

	case 0x05:
		result, overflow := cpu.AddOverflowCheck(cpu.GetRegister(RAX), int64(inst.Imm32))
		if overflow {
			return &EmulatorError{PC: cpu.pc, Message: "overflow detected in ADD"}
		}
		cpu.SetRegister(RAX, result)

	case 0x2D:
		result, overflow := cpu.SubOverflowCheck(cpu.GetRegister(RAX), int64(inst.Imm32))
		if overflow {
			return &EmulatorError{PC: cpu.pc, Message: "overflow detected in SUB"}
		}
		cpu.SetRegister(RAX, result)

	case 0x31:
		src, err := GetRegFromModRM(inst.ModRM, false)
		if err != nil {
			return err
		}
		dst, err := GetRegFromModRM(inst.ModRM, true)
		if err != nil {
			return err
		}
		cpu.SetRegister(dst, cpu.GetRegister(dst)^cpu.GetRegister(src))

	default:
		return fmt.Errorf("unknown opcode: 0x%02X", inst.Opcode)
	}

	cpu.pc += inst.Length
	return nil

}
