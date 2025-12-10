package genhex

import (
	"bytes"
	"errors"
	"fmt"
	rand2 "math/rand"
	"time"
)

func GenerateHex(level int) (spaceHex string, noSpaceHex string, err error) {
	rnd := rand2.New(rand2.NewSource(time.Now().UnixNano()))

	var bytesOut []byte
	switch level {
	case 1:
		bytesOut = genLevel1(rnd)
	case 2:
		bytesOut = genLevel2(rnd)
	case 3:
		bytesOut = genLevel3(rnd)
	case 4:
		bytesOut = genLevel4(rnd)
	default:
		return "", "", errors.New("unsupported level")
	}

	var buf bytes.Buffer
	var noSpaceBuf bytes.Buffer
	for i, b := range bytesOut {
		if i > 0 {
			buf.WriteRune(' ')
		}
		fmt.Fprintf(&buf, "%02x", b)
		fmt.Fprintf(&noSpaceBuf, "%02x", b)
	}

	return buf.String(), noSpaceBuf.String(), nil
}
