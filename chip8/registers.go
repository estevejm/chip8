package chip8

import (
	"fmt"
	"strings"
)

const (
	registerCount = 16
	flagRegister  = 0xF
)

type Registers [registerCount]uint8

func (r Registers) String() string {
	var sb strings.Builder

	for i, b := range r {
		sb.WriteString(fmt.Sprintf("%x:%s ", i, hexdump8(b)))
	}

	return strings.TrimSpace(sb.String())
}
