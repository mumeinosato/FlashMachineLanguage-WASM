package genhex

import (
	"encoding/binary"
	rand2 "math/rand"
)

var regMap = map[string]int{
	"rax": 0,
	"rbx": 3,
	"rcx": 1,
	"rdx": 2,
}

func u32Bytes(v int32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(v))
	return b
}

func u64Bytes(v int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v))
	return b
}

func randInt8(rnd *rand2.Rand) int8 {
	return int8(rnd.Intn(0x100) - 0x80)
}

func randInt16(rnd *rand2.Rand) int16 {
	return int16(rnd.Intn(0x10000) - 0x8000)
}

func randInt32(rnd *rand2.Rand) int32 {
	u := rnd.Int31()
	if rnd.Intn(2) == 0 {
		return int32(u)
	}
	return -int32(u)
}
