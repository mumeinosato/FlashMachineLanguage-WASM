package genhex

import rand2 "math/rand"

const (
	min8  = -128
	max8  = 127
	min16 = -32768
	max16 = 32767
	min32 = -2147483648
	max32 = 2147483647
)

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}

	if v > hi {
		return hi
	}

	return v
}

func clampInt64(v, lo, hi int64) int64 {
	if v < lo {
		return lo
	}

	if v > hi {
		return hi
	}

	return v
}

func randInRange(rnd *rand2.Rand, lo, hi int) int {
	width := hi - lo + 1

	if width <= 0 {
		return lo
	}

	return rnd.Intn(width) + lo
}

func randInt64InRange(rnd *rand2.Rand, lo, hi int64) int64 {
	width := hi - lo + 1

	if width <= 0 {
		return lo
	}

	return lo + rnd.Int63n(width)
}

func applyInitAddSub(out *[]byte, rnd *rand2.Rand, reg int, init int64, bitMin, bitMax int64, cast func(int64) int32) {
	if rnd.Intn(2) == 0 {
		allowedMin := clampInt64(bitMin-init, bitMin, bitMax)
		allowedMax := clampInt64(bitMax-init, bitMin, bitMax)
		imm := randInt64InRange(rnd, allowedMin, allowedMax)
		*out = append(*out, encAddRegImm(reg, cast(imm))...)
	} else {
		allowedMin := clampInt64(init-bitMax, bitMin, bitMax)
		allowedMax := clampInt64(init-bitMin, bitMin, bitMax)
		imm := randInt64InRange(rnd, allowedMin, allowedMax)
		*out = append(*out, encSubRegImm(reg, cast(imm))...)
	}
}

func appendAddImm32(out *[]byte, rnd *rand2.Rand, dst int, regVals map[int]int64) {
	allowedMin := clampInt64(min32-regVals[dst], min32, max32)
	allowedMax := clampInt64(max32-regVals[dst], min32, max32)
	imm64 := randInt64InRange(rnd, allowedMin, allowedMax)
	*out = append(*out, encAddRegImm(dst, int32(imm64))...)
	regVals[dst] += imm64
}

func appendSubImm32(out *[]byte, rnd *rand2.Rand, dst int, regVals map[int]int64) {
	allowedMin := clampInt64(regVals[dst]-max32, min32, max32)
	allowedMax := clampInt64(regVals[dst]-min32, min32, max32)
	imm64 := randInt64InRange(rnd, allowedMin, allowedMax)
	*out = append(*out, encSubRegImm(dst, int32(imm64))...)
	regVals[dst] -= imm64
}

func genLevel1(rnd *rand2.Rand) []byte {
	var out []byte
	val := randInt8(rnd)
	out = append(out, encMovRegImm(regMap["rax"], int32(val))...)

	applyInitAddSub(&out, rnd, regMap["rax"], int64(val), min8, max8, func(i int64) int32 { return int32(int8(i)) })

	return out
}

func genLevel2(rnd *rand2.Rand) []byte {
	var out []byte
	v := randInt16(rnd)
	out = append(out, encMovRegImm(regMap["rax"], int32(v))...)

	applyInitAddSub(&out, rnd, regMap["rax"], int64(v), min16, max16, func(i int64) int32 { return int32(int16(i)) })
	return out
}

func genLevel3(rnd *rand2.Rand) []byte {
	var out []byte
	v := randInt32(rnd)
	out = append(out, encMovRegImm(regMap["rax"], int32(v))...)

	applyInitAddSub(&out, rnd, regMap["rax"], int64(v), min32, max32, func(i int64) int32 { return int32(i) })
	return out
}

func genLevel4(rnd *rand2.Rand) []byte {
	var out []byte
	regs := []string{"rax", "rbx", "rcx", "rdx"}

	regVals := make(map[int]int64)
	for _, r := range regs {
		v := randInt32(rnd)
		reg := regMap[r]
		regVals[reg] = int64(v)
		out = append(out, encMovRegImm(reg, int32(v))...)
	}

	ops := 2 + rnd.Intn(2)
	for i := 0; i < ops; i++ {
		dstName := regs[rnd.Intn(len(regs))]
		dst := regMap[dstName]

		if rnd.Intn(2) == 0 {
			// try register ops first
			type cand struct {
				op  int
				src int
			}
			var cands []cand
			for _, srcName := range regs {
				src := regMap[srcName]
				resAdd := regVals[dst] + regVals[src]
				if resAdd >= min32 && resAdd <= max32 {
					cands = append(cands, cand{0, src})
				}
				resSub := regVals[dst] - regVals[src]
				if resSub >= min32 && resSub <= max32 {
					cands = append(cands, cand{1, src})
				}
			}

			if len(cands) > 0 {
				choice := cands[rnd.Intn(len(cands))]
				if choice.op == 0 {
					out = append(out, encAddRegReg(dst, choice.src)...)
					regVals[dst] = regVals[dst] + regVals[choice.src]
				} else {
					out = append(out, encSubRegReg(dst, choice.src)...)
					regVals[dst] = regVals[dst] - regVals[choice.src]
				}
			} else {
				if rnd.Intn(2) == 0 {
					appendAddImm32(&out, rnd, dst, regVals)
				} else {
					appendSubImm32(&out, rnd, dst, regVals)
				}
			}
		} else {
			if rnd.Intn(2) == 0 {
				appendAddImm32(&out, rnd, dst, regVals)
			} else {
				appendSubImm32(&out, rnd, dst, regVals)
			}
		}
	}
	return out
}
