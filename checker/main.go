package checker

import "encoding/hex"

func CheckLevel(codeHex string) (int, error) {
	code, err := hex.DecodeString(codeHex)
	if err != nil {
		return 0, err
	}

	maxImmSize := 0
	calcCount := 0

	for i := 0; i < len(code); {
		op := code[i]

		if op >= 0xB8 && op <= 0xBF {
			if i+5 > len(code) {
				return 0, nil
			}
			maxImmSize = max(maxImmSize, immClass(code[i+1:i+5]))
			i += 5
			continue
		}

		if op == 0x48 {
			if i+1 >= len(code) {
				break
			}
			next := code[i+1]

			if next == 0x81 {
				if i+7 > len(code) {
					break
				}
				maxImmSize = max(maxImmSize, immClass(code[i+3:i+7]))
				calcCount++
				i += 7
				continue
			}

			if next == 0xB8 {
				if i+10 > len(code) {
					break
				}
				maxImmSize = 4
				i += 10
				continue
			}
		}

		if op == 0x05 || op == 0x2D {
			if i+5 > len(code) {
				break
			}
			maxImmSize = max(maxImmSize, immClass(code[i+1:i+5]))
			calcCount++
			i += 5
			continue
		}

		if op == 0x01 || op == 0x29 {
			calcCount++
			i += 2
			continue
		}

		i++
	}

	switch {
	case maxImmSize == 1 && calcCount >= 1:
		return 1, nil
	case maxImmSize == 2:
		return 2, nil
	case maxImmSize == 4 && calcCount == 1:
		return 3, nil
	case maxImmSize == 4 && calcCount >= 2:
		return 4, nil
	}
	return 0, nil
}

func immClass(b []byte) int {
	v := int(int32(int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24))
	if v >= -128 && v <= 127 {
		return 1
	}
	if v >= -32768 && v <= 32767 {
		return 2
	}
	return 4
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
