package emulator

import "fmt"

type Decoder struct {
	code []byte
	pos  int
}

type Instruction struct {
	Opcode   byte
	ModRM    byte
	HasModRM bool
	Rex      byte
	HasRex   bool
	Imm32    int32
	HasImm32 bool
	Imm64    int64
	HasImm64 bool
	Length   int
}

func NewDecoder(code []byte) *Decoder {
	return &Decoder{
		code: code,
		pos:  0,
	}
}

func (d *Decoder) ReadByte() (byte, error) {
	if d.pos >= len(d.code) {
		return 0, fmt.Errorf("unexpected end of code")
	}
	b := d.code[d.pos]
	d.pos++
	return b, nil
}

func (d *Decoder) PeekByte() (byte, error) {
	if d.pos >= len(d.code) {
		return 0, fmt.Errorf("unexpected end of code")
	}
	return d.code[d.pos], nil
}

func (d *Decoder) ReadImm32() (int32, error) {
	if d.pos+4 > len(d.code) {
		return 0, fmt.Errorf("unexpected end of code")
	}
	imm := int32(d.code[d.pos]) | int32(d.code[d.pos+1])<<8 | int32(d.code[d.pos+2])<<16 | int32(d.code[d.pos+3])<<24
	d.pos += 4
	return imm, nil
}

func (d *Decoder) ReadImm64() (int64, error) {
	if d.pos+8 > len(d.code) {
		return 0, fmt.Errorf("unexpected end of code")
	}
	imm := int64(d.code[d.pos]) | int64(d.code[d.pos+1])<<8 | int64(d.code[d.pos+2])<<16 |
		int64(d.code[d.pos+3])<<24 | int64(d.code[d.pos+4])<<32 | int64(d.code[d.pos+5])<<40 |
		int64(d.code[d.pos+6])<<48 | int64(d.code[d.pos+7])<<56
	d.pos += 8
	return imm, nil
}

func (d *Decoder) DecodeNext() (*Instruction, error) {
	startPos := d.pos
	inst := &Instruction{}

	b, err := d.PeekByte()
	if err != nil {
		return nil, err
	}
	if b >= 0x48 && b <= 0x4F {
		inst.Rex, _ = d.ReadByte()
		inst.HasRex = true
	}

	inst.Opcode, err = d.ReadByte()
	if err != nil {
		return nil, err
	}

	needsModRM := false
	needsImm32 := false

	switch inst.Opcode {
	case 0x89, 0x8B:
		needsModRM = true
	case 0x01, 0x29, 0x31:
		needsModRM = true
	case 0x05, 0x2D:
		needsImm32 = true
	case 0x81, 0xC7:
		needsModRM = true
		needsImm32 = true
	case 0xB8, 0xB9, 0xBA, 0xBB:
		if inst.HasRex {
			inst.Imm64, err = d.ReadImm64()
			if err != nil {
				return nil, err
			}
			inst.HasImm64 = true
		} else {
			needsImm32 = true
		}
	default:
		return nil, &EmulatorError{PC: startPos, Message: fmt.Sprintf("unknown opcode 0x%02X", inst.Opcode)}
	}

	if needsModRM {
		inst.ModRM, err = d.ReadByte()
		if err != nil {
			return nil, err
		}
		inst.HasModRM = true

		mod := (inst.ModRM >> 6) & 0x03
		if mod != 0x03 {
			return nil, &EmulatorError{PC: startPos, Message: "memory access not supported"}
		}
	}

	if needsImm32 {
		inst.Imm32, err = d.ReadImm32()
		if err != nil {
			return nil, err
		}
		inst.HasImm32 = true
	}

	inst.Length = d.pos - startPos
	return inst, nil
}

func (d *Decoder) HasMore() bool {
	return d.pos < len(d.code)
}

func GetRegFromModRM(modrm byte, isRM bool) (Register, error) {
	var code byte
	if isRM {
		code = modrm & 0x07
	} else {
		code = (modrm >> 3) & 0x07
	}

	switch code {
	case 0:
		return RAX, nil
	case 3:
		return RBX, nil
	case 1:
		return RCX, nil
	case 2:
		return RDX, nil
	default:
		return 0, fmt.Errorf("not allowed register code: %d", code)

	}
}
