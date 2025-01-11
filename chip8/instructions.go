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

// Return 00EE: Return from a subroutine
func Return() Instruction {
	return &returnFromSubroutine{}
}

type returnFromSubroutine struct{}

func (i returnFromSubroutine) String() string {
	return "RTS"
}

func (i returnFromSubroutine) Execute(c *Chip8) {
	// TODO: check stack pointer won't be < 0
	c.stackPointer--
	c.programCounter = c.stack[c.stackPointer]
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

// Call 2NNN: Execute subroutine starting at address NNN
func Call(opcode uint16) Instruction {
	return &call{
		n: opcode & 0xFFF,
	}
}

type call struct {
	n uint16
}

func (i call) String() string {
	return fmt.Sprintf("CALL 0x%04x", i.n)
}

func (i call) Execute(c *Chip8) {
	// TODO: check stack overflow
	c.stack[c.stackPointer] = c.programCounter
	c.stackPointer++
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

// SkipEqualRegister 5XY0: Skip the following instruction
// if the value of register VX is equal to the value of register VY
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
	return fmt.Sprintf("SKE V%x,V%x", i.x, i.y)
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

// Move 8XY0: Store the value of register VY in register VX
func Move(opcode uint16) Instruction {
	return &move{
		x: uint8(opcode>>8) & 0xF,
		y: uint8(opcode>>4) & 0xF,
	}
}

type move struct {
	x, y uint8
}

func (i move) String() string {
	return fmt.Sprintf("MOVE V%x,V%x", i.x, i.y)
}

func (i move) Execute(c *Chip8) {
	c.registers[i.x] = c.registers[i.y]
}

// Or 8XY1: Set VX to VX OR VY
func Or(opcode uint16) Instruction {
	return &or{
		x: uint8(opcode>>8) & 0xF,
		y: uint8(opcode>>4) & 0xF,
	}
}

type or struct {
	x, y uint8
}

func (i or) String() string {
	return fmt.Sprintf("OR V%x,V%x", i.x, i.y)
}

func (i or) Execute(c *Chip8) {
	c.registers[i.x] |= c.registers[i.y]
}

// And 8XY2: Set VX to VX AND VY
func And(opcode uint16) Instruction {
	return &and{
		x: uint8(opcode>>8) & 0xF,
		y: uint8(opcode>>4) & 0xF,
	}
}

type and struct {
	x, y uint8
}

func (i and) String() string {
	return fmt.Sprintf("AND V%x,V%x", i.x, i.y)
}

func (i and) Execute(c *Chip8) {
	c.registers[i.x] &= c.registers[i.y]
}

// Xor 8XY3: Set VX to VX XOR VY
func Xor(opcode uint16) Instruction {
	return &xor{
		x: uint8(opcode>>8) & 0xF,
		y: uint8(opcode>>4) & 0xF,
	}
}

type xor struct {
	x, y uint8
}

func (i xor) String() string {
	return fmt.Sprintf("XOR V%x,V%x", i.x, i.y)
}

func (i xor) Execute(c *Chip8) {
	c.registers[i.x] ^= c.registers[i.y]
}

// AddRegister 8XY4: Add the value of register VY to register VX
// Set VF to 01 if a carry occurs
// Set VF to 00 if a carry does not occur
func AddRegister(opcode uint16) Instruction {
	return &addRegister{
		x: uint8(opcode>>8) & 0xF,
		y: uint8(opcode>>4) & 0xF,
	}
}

type addRegister struct {
	x, y uint8
}

func (i addRegister) String() string {
	return fmt.Sprintf("ADD V%x,V%x", i.x, i.y)
}

func (i addRegister) Execute(c *Chip8) {
	c.registers[i.x] += c.registers[i.y]
	if c.registers[i.x] < c.registers[i.y] {
		c.registers[flagRegister] = 1
	} else {
		c.registers[flagRegister] = 0
	}
}

// SubRegister 8XY5: Subtract the value of register VY from register VX
// Set VF to 00 if a borrow occurs
// Set VF to 01 if a borrow does not occur
func SubRegister(opcode uint16) Instruction {
	return &subRegister{
		x: uint8(opcode>>8) & 0xF,
		y: uint8(opcode>>4) & 0xF,
	}
}

type subRegister struct {
	x, y uint8
}

func (i subRegister) String() string {
	return fmt.Sprintf("SUB V%x,V%x", i.x, i.y)
}

func (i subRegister) Execute(c *Chip8) {
	if c.registers[i.y] > c.registers[i.x] {
		c.registers[flagRegister] = 1
	} else {
		c.registers[flagRegister] = 0
	}
	c.registers[i.x] -= c.registers[i.y]
}

// ShiftRight 8XY6: Store the value of register VY shifted right one bit in register VX
// Set register VF to the least significant bit prior to the shift
// VY is unchanged
func ShiftRight(opcode uint16) Instruction {
	return &shiftRight{
		x: uint8(opcode>>8) & 0xF,
		y: uint8(opcode>>4) & 0xF,
	}
}

type shiftRight struct {
	x, y uint8
}

func (i shiftRight) String() string {
	return fmt.Sprintf("SHR V%x,V%x", i.x, i.y)
}

func (i shiftRight) Execute(c *Chip8) {
	c.registers[flagRegister] = c.registers[i.y] & 1
	c.registers[i.x] = c.registers[i.y] >> 1
}

// ReverseSubRegister 8XY7: Set register VX to the value of VY minus VX
// Set VF to 00 if a borrow occurs
// Set VF to 01 if a borrow does not occur
func ReverseSubRegister(opcode uint16) Instruction {
	return &reverseSubRegister{
		x: uint8(opcode>>8) & 0xF,
		y: uint8(opcode>>4) & 0xF,
	}
}

type reverseSubRegister struct {
	x, y uint8
}

func (i reverseSubRegister) String() string {
	return fmt.Sprintf("RSB V%x,V%x", i.x, i.y)
}

func (i reverseSubRegister) Execute(c *Chip8) {
	if c.registers[i.x] > c.registers[i.y] {
		c.registers[flagRegister] = 1
	} else {
		c.registers[flagRegister] = 0
	}
	c.registers[i.x] = c.registers[i.y] - c.registers[i.x]
}

// ShiftLeft 8XYE: Store the value of register VY shifted left one bit in register VX
// Set register VF to the most significant bit prior to the shift
// VY is unchanged
func ShiftLeft(opcode uint16) Instruction {
	return &shiftLeft{
		x: uint8(opcode>>8) & 0xF,
		y: uint8(opcode>>4) & 0xF,
	}
}

type shiftLeft struct {
	x, y uint8
}

func (i shiftLeft) String() string {
	return fmt.Sprintf("SHL V%x,V%x", i.x, i.y)
}

func (i shiftLeft) Execute(c *Chip8) {
	c.registers[flagRegister] = c.registers[i.y] & 0x80
	c.registers[i.x] = c.registers[i.y] << 1
}

// SkipNotEqualRegister 9XY0: Skip the following instruction
// if the value of register VX is not equal to the value of register VY
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
	return fmt.Sprintf("SKNE V%x,V%x", i.x, i.y)
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

	c.registers[flagRegister] = vf
}

// Read FX65: Fill registers V0 to VX inclusive with the values stored in memory starting at address I
// I is set to I + X + 1 after operation
func Read(opcode uint16) Instruction {
	return &read{
		x: (opcode & 0x0F00) >> 8,
	}
}

type read struct {
	x uint16
}

func (i read) String() string {
	return fmt.Sprintf("READ V0-V%x", i.x)
}

func (i read) Execute(c *Chip8) {
	high := i.x + 1
	copy(c.registers[:high], c.memory[c.index:c.index+high])
	c.index += high
}

// Write FX55: Store the values of registers V0 to VX inclusive in memory starting at address I
// I is set to I + X + 1 after operation
func Write(opcode uint16) Instruction {
	return &write{
		x: (opcode & 0x0F00) >> 8,
	}
}

type write struct {
	x uint16
}

func (i write) String() string {
	return fmt.Sprintf("WRITE V0-V%x", i.x)
}

func (i write) Execute(c *Chip8) {
	high := i.x + 1
	copy(c.memory[c.index:c.index+high], c.registers[:high])
	c.index += high
}
