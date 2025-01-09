package chip8

import (
	"fmt"
)

type Instruction interface {
	fmt.Stringer
	Execute(c *Chip8)
}

// ClearScreen 00E0: Clear the screen
func ClearScreen() Instruction {
	return &clearScreen{}
}

type clearScreen struct{}

func (i clearScreen) String() string {
	return "CLR"
}

func (i clearScreen) Execute(c *Chip8) {
	for y := range c.screen {
		for x := range c.screen[y] {
			c.screen[y][x] = false
		}
	}
}

// Jump 1NNN: Jump to address NNN
func Jump(opcode uint16) Instruction {
	return &jump{
		n: opcode & 0xFFF,
	}
}

type jump struct {
	n uint16
}

func (i jump) String() string {
	return fmt.Sprintf("JUMP 0x%04x", i.n)
}

func (i jump) Execute(c *Chip8) {
	c.programCounter = i.n
}

// SkipEqual 3XNN: Skip the following instruction if the value of register VX equals NN
func SkipEqual(opcode uint16) Instruction {
	return &skipEqual{
		x: uint8(opcode>>8) & 0xF,
		n: uint8(opcode & 0xFF),
	}
}

type skipEqual struct {
	x, n uint8
}

func (i skipEqual) String() string {
	return fmt.Sprintf("SKE V%x,%d", i.x, i.n)
}

func (i skipEqual) Execute(c *Chip8) {
	if c.registers[i.x] == i.n {
		c.incrementProgramCounter()
	}
}

// SkipNotEqual 4XNN: Skip the following instruction if the value of register VX is not equal to NN
func SkipNotEqual(opcode uint16) Instruction {
	return &skipNotEqual{
		x: uint8(opcode>>8) & 0xF,
		n: uint8(opcode & 0xFF),
	}
}

type skipNotEqual struct {
	x, n uint8
}

func (i skipNotEqual) String() string {
	return fmt.Sprintf("SKNE V%x,%d", i.x, i.n)
}

func (i skipNotEqual) Execute(c *Chip8) {
	if c.registers[i.x] != i.n {
		c.incrementProgramCounter()
	}
}

// SkipEqualRegister 5XY0: Skip the following instruction if the value of register VX is equal to the value of register VY
func SkipEqualRegister(opcode uint16) Instruction {
	return &skipEqualRegister{
		x: uint8(opcode>>8) & 0xF,
		y: uint8(opcode>>4) & 0xF,
	}
}

type skipEqualRegister struct {
	x, y uint8
}

func (i skipEqualRegister) String() string {
	return fmt.Sprintf("SKRE V%x,V%x", i.x, i.y)
}

func (i skipEqualRegister) Execute(c *Chip8) {
	if c.registers[i.x] == c.registers[i.y] {
		c.incrementProgramCounter()
	}
}

// Load 6XNN: Store number NN in register VX
func Load(opcode uint16) Instruction {
	return &load{
		x: uint8(opcode>>8) & 0xF,
		n: uint8(opcode & 0xFF),
	}
}

type load struct {
	x, n uint8
}

func (i load) String() string {
	return fmt.Sprintf("LOAD V%x,%d", i.x, i.n)
}

func (i load) Execute(c *Chip8) {
	c.registers[i.x] = i.n
}

// Add 7XNN: Add the value NN to register VX
func Add(opcode uint16) Instruction {
	return &add{
		x: uint8(opcode>>8) & 0xF,
		n: uint8(opcode & 0xFF),
	}
}

type add struct {
	x, n uint8
}

func (i add) String() string {
	return fmt.Sprintf("ADD V%x,%d", i.x, i.n)
}

func (i add) Execute(c *Chip8) {
	c.registers[i.x] += i.n
}

// SkipNotEqualRegister 9XY0: Skip the following instruction if the value of register VX is not equal to the value of register VY
func SkipNotEqualRegister(opcode uint16) Instruction {
	return &skipNotEqualRegister{
		x: uint8(opcode>>8) & 0xF,
		y: uint8(opcode>>4) & 0xF,
	}
}

type skipNotEqualRegister struct {
	x, y uint8
}

func (i skipNotEqualRegister) String() string {
	return fmt.Sprintf("SKRNE V%x,V%x", i.x, i.y)
}

func (i skipNotEqualRegister) Execute(c *Chip8) {
	if c.registers[i.x] != c.registers[i.y] {
		c.incrementProgramCounter()
	}
}

// LoadIndex ANNN: Store memory address NNN in register I
func LoadIndex(opcode uint16) Instruction {
	return &loadIndex{
		n: opcode & 0xFFF,
	}
}

type loadIndex struct {
	n uint16
}

func (i loadIndex) String() string {
	return fmt.Sprintf("LOADI 0x%04x", i.n)
}

func (i loadIndex) Execute(c *Chip8) {
	c.index = i.n
}

// DrawSprite DXYN: Draw a sprite at position VX, VY with N bytes of sprite data starting at the address stored in I
// Set VF to 01 if any set pixels are changed to unset, and 00 otherwise
func DrawSprite(opcode uint16) Instruction {
	return &drawSprite{
		x: uint8(opcode>>8) & 0xF,
		y: uint8(opcode>>4) & 0xF,
		n: uint8(opcode & 0xF),
	}
}

type drawSprite struct {
	x, y, n uint8
}

func (i drawSprite) String() string {
	return fmt.Sprintf("DRAW V%x,V%x,%d", i.x, i.y, i.n)
}

func (i drawSprite) Execute(c *Chip8) {
	// always start drawing in boundary
	vx := int(c.registers[i.x]) % screenWidth
	vy := int(c.registers[i.y]) % screenHeight

	sprite := c.memory[c.index : c.index+uint16(i.n)]
	vf := byte(0)
	for i, b := range sprite {
		for j := 0; j < bytePixels; j++ {
			pixelX := vx + j
			pixelY := vy + i

			// check clipping
			if pixelX >= screenWidth || pixelY >= screenHeight {
				continue
			}

			screenPixelIsSet := c.screen[pixelY][pixelX]
			spritePixelIsSet := (b>>(bytePixels-1-j))&1 == 1
			if screenPixelIsSet && spritePixelIsSet {
				vf = byte(1) // collision detected
			}

			// draw using XOR, boolean != should be equivalent
			c.screen[pixelY][pixelX] = screenPixelIsSet != spritePixelIsSet
		}
	}

	c.registers[0xF] = vf
}
