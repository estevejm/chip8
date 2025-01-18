package chip8

import "fmt"

func hexdump8(b uint8) string {
	return fmt.Sprintf("%02x", b)
}
func hexdump16(b uint16) string {
	return fmt.Sprintf("%04x", b)
}
