package chip8

import (
	"io"
	"strings"
)

const memoryLocations = 0x1000

type Memory [memoryLocations]uint8

// ReadWord returns 2 big-endian bytes
func (m *Memory) ReadWord(address uint16) uint16 {
	msb := m.ReadByte(address)
	lsb := m.ReadByte(address + 1)
	return uint16(msb)<<8 | uint16(lsb)
}

func (m *Memory) ReadByte(address uint16) uint8 {
	return m[address]
}

func (m *Memory) Write(address uint16, r io.Reader) (int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}

	// TODO: ensure data copied in bounds -> check data len
	copy(m[address:], data)

	return len(data), nil
}

func (m *Memory) String() string {
	const bytesPerRow = 16
	var sb strings.Builder

	for i, b := range m {
		if i%bytesPerRow == 0 {
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(hexdump16(uint16(i)))
			sb.WriteString(" ")
		}
		sb.WriteString(" ")
		sb.WriteString(hexdump8(b))
	}

	return sb.String()
}
