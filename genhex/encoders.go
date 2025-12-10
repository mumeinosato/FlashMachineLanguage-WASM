package genhex

func encMovRegImm(reg int, imm int32) []byte {
	op := byte(0xB8 + reg)
	out := []byte{op}
	out = append(out, u32Bytes(imm)...)
	return out
}

func encMovRegReg(dst, src int) []byte {
	modrm := byte(0xC0 | (src << 3) | dst)
	return []byte{0x48, 0x89, modrm}
}

func encAddRegImm(reg int, imm int32) []byte {
	modrm := byte(0xC0 | reg)
	out := []byte{0x48, 0x81, modrm}
	out = append(out, u32Bytes(imm)...)
	return out
}

func encSubRegImm(reg int, imm int32) []byte {
	modrm := byte(0xC0 | (5 << 3) | reg)
	out := []byte{0x48, 0x81, modrm}
	out = append(out, u32Bytes(imm)...)
	return out
}

func encAddRegReg(dst, src int) []byte {
	modrm := byte(0xC0 | (src << 3) | dst)
	return []byte{0x48, 0x01, modrm}
}

func encSubRegReg(dst, src int) []byte {
	modrm := byte(0xC0 | (src << 3) | dst)
	return []byte{0x48, 0x29, modrm}
}
