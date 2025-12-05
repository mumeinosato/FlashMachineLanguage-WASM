package emulator

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseHexString(hexStr string) ([]byte, error) {
	hexStr = strings.ReplaceAll(hexStr, " ", "")

	if len(hexStr)%2 != 0 {
		return nil, fmt.Errorf("hex string length must be even")
	}

	code := make([]byte, len(hexStr)/2)
	for i := 0; i < len(hexStr); i += 2 {
		b, err := strconv.ParseUint(hexStr[i:i+2], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid hex character at position %d: %v", i, err)
		}
		code[i/2] = byte(b)
	}
	return code, nil
}
