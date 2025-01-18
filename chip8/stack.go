package chip8

import (
	"fmt"
	"strings"
)

const stackLevels = 16

type Stack struct {
	data    [stackLevels]uint16
	pointer uint8
}

func NewStack() *Stack {
	return &Stack{
		data:    [stackLevels]uint16{},
		pointer: 0,
	}
}

func (s *Stack) Push(value uint16) {
	// TODO: check stack overflow
	s.data[s.pointer] = value
	s.pointer++
}

func (s *Stack) Pop() uint16 {
	// TODO: check stack pointer won't be < 0
	s.pointer--
	return s.data[s.pointer]
}

func (s *Stack) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("P:%d ", s.pointer))

	for i, b := range s.data[:s.pointer] {
		sb.WriteString(fmt.Sprintf("%x:%s ", i, hexdump16(b)))
	}

	return strings.TrimSpace(sb.String())
}
