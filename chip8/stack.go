package chip8

import (
	"fmt"
	"strings"
)

const stackLevels = 16

type Stack [stackLevels]uint16

func (s Stack) String() string {
	var sb strings.Builder

	for i, b := range s {
		sb.WriteString(fmt.Sprintf("%x:%s ", i, hexdump16(b)))
	}

	return strings.TrimSpace(sb.String())
}
