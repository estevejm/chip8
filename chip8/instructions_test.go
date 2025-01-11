package chip8

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddRegister(t *testing.T) {
	c := NewChip8(slog.Default())
	c.registers[0] = 0x07
	c.registers[1] = 0x03

	i := addRegister{x: 0, y: 1}
	i.Execute(c)

	assert.Equal(t, byte(0x0a), c.registers[0])
	assert.Equal(t, byte(0x03), c.registers[1])
	assert.Equal(t, byte(0x00), c.registers[0xF])
}

func TestAddRegisterCarry(t *testing.T) {
	c := NewChip8(slog.Default())
	c.registers[0] = 0xff
	c.registers[1] = 0x11

	i := addRegister{x: 0, y: 1}
	i.Execute(c)

	assert.Equal(t, byte(0x10), c.registers[0])
	assert.Equal(t, byte(0x01), c.registers[0xF])
}

func TestSubRegister(t *testing.T) {
	c := NewChip8(slog.Default())
	c.registers[0] = 0x07
	c.registers[1] = 0x03

	i := subRegister{x: 0, y: 1}
	i.Execute(c)

	assert.Equal(t, byte(0x04), c.registers[0])
	assert.Equal(t, byte(0x03), c.registers[1])
	assert.Equal(t, byte(0x00), c.registers[0xF])
}

func TestSubRegisterBorrow(t *testing.T) {
	c := NewChip8(slog.Default())
	c.registers[0] = 0x00
	c.registers[1] = 0x01

	i := subRegister{x: 0, y: 1}
	i.Execute(c)

	assert.Equal(t, byte(0xff), c.registers[0])
	assert.Equal(t, byte(0x01), c.registers[0xF])
}

func TestReverseSubRegister(t *testing.T) {
	c := NewChip8(slog.Default())
	c.registers[0] = 0x03
	c.registers[1] = 0x07

	i := reverseSubRegister{x: 0, y: 1}
	i.Execute(c)

	assert.Equal(t, byte(0x04), c.registers[0])
	assert.Equal(t, byte(0x07), c.registers[1])
	assert.Equal(t, byte(0x00), c.registers[0xF])
}

func TestReverseSubRegisterBorrow(t *testing.T) {
	c := NewChip8(slog.Default())
	c.registers[0] = 0x01
	c.registers[1] = 0x00

	i := reverseSubRegister{x: 0, y: 1}
	i.Execute(c)

	assert.Equal(t, byte(0xff), c.registers[0])
	assert.Equal(t, byte(0x01), c.registers[0xF])
}
