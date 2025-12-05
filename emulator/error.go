package emulator

import "fmt"

type EmulatorError struct {
	PC      int
	Message string
}

func (e *EmulatorError) Error() string {
	return fmt.Sprintf("Error at PC:%d %s", e.PC, e.Message)
}
