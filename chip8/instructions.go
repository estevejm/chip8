package chip8

import (
	"crypto/rand"
	"fmt"
)

type Instruction interface {
	fmt.Stringer
	Execute(c *Chip8)
}

// NoOperation do nothing
func NoOperation() Instruction {
	return &noOperation{}
}

type noOperation struct{}

func (i noOperation) String() string {
	return "NOP"
}

func (i noOperation) Execute(_ *Chip8) {

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
	c.screen.Clear()
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
	c.fetcher.SetCounter(c.stack[c.stackPointer])
}

// Jump 1NNN: Jump to address NNN
func Jump(nnn uint16) Instruction {
	return &jump{nnn: nnn}
}

type jump struct {
	nnn uint16
}

func (i jump) String() string {
	return fmt.Sprintf("JUMP 0x%04x", i.nnn)
}

func (i jump) Execute(c *Chip8) {
	c.fetcher.SetCounter(i.nnn)
}

// Call 2NNN: Execute subroutine starting at address NNN
func Call(nnn uint16) Instruction {
	return &call{
		nnn: nnn & 0xFFF,
	}
}

type call struct {
	nnn uint16
}

func (i call) String() string {
	return fmt.Sprintf("CALL 0x%04x", i.nnn)
}

func (i call) Execute(c *Chip8) {
	// TODO: check stack overflow
	c.stack[c.stackPointer] = c.fetcher.GetCounter()
	c.stackPointer++
	c.fetcher.SetCounter(i.nnn)
}

// SkipEqual 3XNN: Skip the following instruction if the value of register VX equals NN
func SkipEqual(x, nn uint8) Instruction {
	return &skipEqual{x: x, nn: nn}
}

type skipEqual struct {
	x, nn uint8
}

func (i skipEqual) String() string {
	return fmt.Sprintf("SKE V%x,%x", i.x, i.nn)
}

func (i skipEqual) Execute(c *Chip8) {
	if c.registers[i.x] == i.nn {
		c.fetcher.Skip()
	}
}

// SkipNotEqual 4XNN: Skip the following instruction if the value of register VX is not equal to NN
func SkipNotEqual(x, nn uint8) Instruction {
	return &skipNotEqual{x: x, nn: nn}
}

type skipNotEqual struct {
	x, nn uint8
}

func (i skipNotEqual) String() string {
	return fmt.Sprintf("SKNE V%x,%x", i.x, i.nn)
}

func (i skipNotEqual) Execute(c *Chip8) {
	if c.registers[i.x] != i.nn {
		c.fetcher.Skip()
	}
}

// SkipEqualRegister 5XY0: Skip the following instruction
// if the value of register VX is equal to the value of register VY
func SkipEqualRegister(x, y uint8) Instruction {
	return &skipEqualRegister{x: x, y: y}
}

type skipEqualRegister struct {
	x, y uint8
}

func (i skipEqualRegister) String() string {
	return fmt.Sprintf("SKE V%x,V%x", i.x, i.y)
}

func (i skipEqualRegister) Execute(c *Chip8) {
	if c.registers[i.x] == c.registers[i.y] {
		c.fetcher.Skip()
	}
}

// Load 6XNN: Store number NN in register VX
func Load(x, nn uint8) Instruction {
	return &load{x: x, nn: nn}
}

type load struct {
	x, nn uint8
}

func (i load) String() string {
	return fmt.Sprintf("LOAD V%x,%x", i.x, i.nn)
}

func (i load) Execute(c *Chip8) {
	c.registers[i.x] = i.nn
}

// Add 7XNN: Add the value NN to register VX
func Add(x, nn uint8) Instruction {
	return &add{x: x, nn: nn}
}

type add struct {
	x, nn uint8
}

func (i add) String() string {
	return fmt.Sprintf("ADD V%x,%x", i.x, i.nn)
}

func (i add) Execute(c *Chip8) {
	c.registers[i.x] += i.nn
}

// LoadRegister 8XY0: Store the value of register VY in register VX
func LoadRegister(x, y uint8) Instruction {
	return &loadRegister{x: x, y: y}
}

type loadRegister struct {
	x, y uint8
}

func (i loadRegister) String() string {
	return fmt.Sprintf("LOAD V%x,V%x", i.x, i.y)
}

func (i loadRegister) Execute(c *Chip8) {
	c.registers[i.x] = c.registers[i.y]
}

// Or 8XY1: Set VX to VX OR VY
func Or(x, y uint8) Instruction {
	return &or{x: x, y: y}
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
func And(x, y uint8) Instruction {
	return &and{x: x, y: y}
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
func Xor(x, y uint8) Instruction {
	return &xor{x: x, y: y}
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
func AddRegister(x, y uint8) Instruction {
	return &addRegister{x: x, y: y}
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
func SubRegister(x, y uint8) Instruction {
	return &subRegister{x: x, y: y}
}

type subRegister struct {
	x, y uint8
}

func (i subRegister) String() string {
	return fmt.Sprintf("SUB V%x,V%x", i.x, i.y)
}

func (i subRegister) Execute(c *Chip8) {
	borrow := c.registers[i.y] > c.registers[i.x]
	c.registers[i.x] -= c.registers[i.y]
	if borrow {
		c.registers[flagRegister] = 0
	} else {
		c.registers[flagRegister] = 1
	}
}

// ShiftRight 8XY6: Store the value of register VY shifted right one bit in register VX
// Set register VF to the least significant bit prior to the shift
// VY is unchanged
func ShiftRight(x, y uint8) Instruction {
	return &shiftRight{x: x, y: y}
}

type shiftRight struct {
	x, y uint8
}

func (i shiftRight) String() string {
	return fmt.Sprintf("SHR V%x,V%x", i.x, i.y)
}

func (i shiftRight) Execute(c *Chip8) {
	out := c.registers[i.y] & 1
	c.registers[i.x] = c.registers[i.y] >> 1
	c.registers[flagRegister] = out
}

// ReverseSubRegister 8XY7: Set register VX to the value of VY minus VX
// Set VF to 00 if a borrow occurs
// Set VF to 01 if a borrow does not occur
func ReverseSubRegister(x, y uint8) Instruction {
	return &reverseSubRegister{x: x, y: y}
}

type reverseSubRegister struct {
	x, y uint8
}

func (i reverseSubRegister) String() string {
	return fmt.Sprintf("RSB V%x,V%x", i.x, i.y)
}

func (i reverseSubRegister) Execute(c *Chip8) {
	borrow := c.registers[i.x] > c.registers[i.y]
	c.registers[i.x] = c.registers[i.y] - c.registers[i.x]
	if borrow {
		c.registers[flagRegister] = 0
	} else {
		c.registers[flagRegister] = 1
	}
}

// ShiftLeft 8XYE: Store the value of register VY shifted left one bit in register VX
// Set register VF to the most significant bit prior to the shift
// VY is unchanged
func ShiftLeft(x, y uint8) Instruction {
	return &shiftLeft{x: x, y: y}
}

type shiftLeft struct {
	x, y uint8
}

func (i shiftLeft) String() string {
	return fmt.Sprintf("SHL V%x,V%x", i.x, i.y)
}

func (i shiftLeft) Execute(c *Chip8) {
	out := c.registers[i.y] >> 7
	c.registers[i.x] = c.registers[i.y] << 1
	c.registers[flagRegister] = out
}

// SkipNotEqualRegister 9XY0: Skip the following instruction
// if the value of register VX is not equal to the value of register VY
func SkipNotEqualRegister(x, y uint8) Instruction {
	return &skipNotEqualRegister{x: x, y: y}
}

type skipNotEqualRegister struct {
	x, y uint8
}

func (i skipNotEqualRegister) String() string {
	return fmt.Sprintf("SKNE V%x,V%x", i.x, i.y)
}

func (i skipNotEqualRegister) Execute(c *Chip8) {
	if c.registers[i.x] != c.registers[i.y] {
		c.fetcher.Skip()
	}
}

// LoadIndex ANNN: Store memory address NNN in register I
func LoadIndex(nnn uint16) Instruction {
	return &loadIndex{nnn: nnn}
}

type loadIndex struct {
	nnn uint16
}

func (i loadIndex) String() string {
	return fmt.Sprintf("LOAD I,0x%04x", i.nnn)
}

func (i loadIndex) Execute(c *Chip8) {
	c.index = i.nnn
}

// JumpRegister0 BNNN: Jump to address NNN + V0
func JumpRegister0(nnn uint16) Instruction {
	return &jumpRegister0{nnn: nnn}
}

type jumpRegister0 struct {
	nnn uint16
}

func (i jumpRegister0) String() string {
	return fmt.Sprintf("JUMP 0x%04x+V0", i.nnn)
}

func (i jumpRegister0) Execute(c *Chip8) {
	c.fetcher.SetCounter(i.nnn + uint16(c.registers[0]))
}

// Random CXNN: Set VX to a random number with a mask of NN
func Random(x, nn uint8) Instruction {
	return &random{x: x, nn: nn}
}

type random struct {
	x, nn uint8
}

func (i random) String() string {
	return fmt.Sprintf("RAND V%x,0x%04x", i.x, i.nn)
}

func (i random) Execute(c *Chip8) {
	b := make([]byte, 1)
	rand.Read(b)
	c.registers[i.x] = b[0] & i.nn
}

// DrawSprite DXYN: Draw a sprite at position VX, VY with N bytes of sprite data starting at the address stored in I
// Set VF to 01 if any set pixels are changed to unset, and 00 otherwise
func DrawSprite(x, y, n uint8) Instruction {
	return &drawSprite{x: x, y: y, n: n}
}

type drawSprite struct {
	x, y, n uint8
}

func (i drawSprite) String() string {
	return fmt.Sprintf("DRAW V%x,V%x,%x", i.x, i.y, i.n)
}

func (i drawSprite) Execute(c *Chip8) {
	width, height := c.screen.Layout()

	// always start drawing in boundary
	vx := int(c.registers[i.x]) % width
	vy := int(c.registers[i.y]) % height

	sprite := c.memory[c.index : c.index+uint16(i.n)]
	vf := uint8(0)
	for i, b := range sprite {
		for j := 0; j < bytePixels; j++ {
			pixelX := vx + j
			pixelY := vy + i

			// check clipping
			if pixelX >= width || pixelY >= height {
				continue
			}

			screenPixelIsSet := c.screen[pixelY][pixelX]
			spritePixelIsSet := (b>>(bytePixels-1-j))&1 == 1
			if screenPixelIsSet && spritePixelIsSet {
				vf = uint8(1) // collision detected
			}

			// draw using XOR, boolean != should be equivalent
			c.screen[pixelY][pixelX] = screenPixelIsSet != spritePixelIsSet
		}
	}

	c.registers[flagRegister] = vf
}

// SkipPressed EX9E: Skip the following instruction if the key corresponding to the hex value
// currently stored in register VX is pressed
func SkipPressed(x uint8) Instruction {
	return &skipPressed{x: x}
}

type skipPressed struct {
	x uint8
}

func (i skipPressed) String() string {
	return fmt.Sprintf("SKP V%x", i.x)
}

func (i skipPressed) Execute(c *Chip8) {
	v := c.registers[i.x]
	if c.input.keys[v] {
		c.fetcher.Skip()
	}
}

// SkipNotPressed EXA1: Skip the following instruction if the key corresponding to the hex value
// currently stored in register VX is not pressed
func SkipNotPressed(x uint8) Instruction {
	return &skipNotPressed{x: x}
}

type skipNotPressed struct {
	x uint8
}

func (i skipNotPressed) String() string {
	return fmt.Sprintf("SKNP V%x", i.x)
}

func (i skipNotPressed) Execute(c *Chip8) {
	v := c.registers[i.x]
	if !c.input.keys[v] {
		c.fetcher.Skip()
	}
}

// LoadRegisterDelayTimer FX07: Store the current value of the delay timer in register VX
func LoadRegisterDelayTimer(x uint8) Instruction {
	return &loadRegisterDelayTimer{x: x}
}

type loadRegisterDelayTimer struct {
	x uint8
}

func (i loadRegisterDelayTimer) String() string {
	return fmt.Sprintf("LOAD V%x,DT", i.x)
}

func (i loadRegisterDelayTimer) Execute(c *Chip8) {
	c.registers[i.x] = c.delayTimer.GetValue()
}

// WaitKey FX0A: Wait for a keypress and store the result in register VX
func WaitKey(x uint8) Instruction {
	return &waitKey{x: x}
}

type waitKey struct {
	x uint8
}

func (i waitKey) String() string {
	return fmt.Sprintf("LOAD V%x,K", i.x)
}

func (i waitKey) Execute(c *Chip8) {
	c.input.Wait(func(key uint8) {
		c.registers[i.x] = key
	})
}

// LoadDelayTimerRegister FX15: Set the delay timer to the value of register VX
func LoadDelayTimerRegister(x uint8) Instruction {
	return &loadDelayTimerRegister{x: x}
}

type loadDelayTimerRegister struct {
	x uint8
}

func (i loadDelayTimerRegister) String() string {
	return fmt.Sprintf("LOAD DT,V%x", i.x)
}

func (i loadDelayTimerRegister) Execute(c *Chip8) {
	c.delayTimer.SetValue(c.registers[i.x])
}

// LoadSoundTimerRegister FX15: Set the delay timer to the value of register VX
func LoadSoundTimerRegister(x uint8) Instruction {
	return &loadSoundTimerRegister{x: x}
}

type loadSoundTimerRegister struct {
	x uint8
}

func (i loadSoundTimerRegister) String() string {
	return fmt.Sprintf("LOAD ST,V%x", i.x)
}

func (i loadSoundTimerRegister) Execute(c *Chip8) {
	c.sound.SetTimerValue(c.registers[i.x])
}

// AddIndex FX1E: Add the value stored in register VX to register I
func AddIndex(x uint8) Instruction {
	return &addIndex{x: x}
}

type addIndex struct {
	x uint8
}

func (i addIndex) String() string {
	return fmt.Sprintf("ADD I,V%x", i.x)
}

func (i addIndex) Execute(c *Chip8) {
	c.index += uint16(c.registers[i.x])
}

// LoadDigitIndex FX29: Set I to the memory address of the sprite data
// corresponding to the hexadecimal digit stored in register VX
func LoadDigitIndex(x uint8) Instruction {
	return &loadDigitIndex{x: x}
}

type loadDigitIndex struct {
	x uint8
}

func (i loadDigitIndex) String() string {
	return fmt.Sprintf("LOAD I,V%x", i.x)
}

func (i loadDigitIndex) Execute(c *Chip8) {
	c.index = fontStartMemoryAddress + uint16(c.registers[i.x])
}

// BCD FX33: Store the binary-coded decimal equivalent of the value stored in register VX
// at addresses I, I+1, and I + 2
func BCD(x uint8) Instruction {
	return &bcd{x: x}
}

type bcd struct {
	x uint8
}

func (i bcd) String() string {
	return fmt.Sprintf("BCD V%x", i.x)
}

func (i bcd) Execute(c *Chip8) {
	v := c.registers[i.x]
	c.memory[c.index] = v / 100
	c.memory[c.index+1] = v % 100 / 10
	c.memory[c.index+2] = v % 10
}

// Write FX55: Store the values of registers V0 to VX inclusive in memory starting at address I
// I is set to I + X + 1 after operation
func Write(x uint8) Instruction {
	return &write{x: x}
}

type write struct {
	x uint8
}

func (i write) String() string {
	return fmt.Sprintf("WRITE V0-V%x", i.x)
}

func (i write) Execute(c *Chip8) {
	high := uint16(i.x + 1)
	copy(c.memory[c.index:c.index+high], c.registers[:high])
	c.index += high
}

// Read FX65: Fill registers V0 to VX inclusive with the values stored in memory starting at address I
// I is set to I + X + 1 after operation
func Read(x uint8) Instruction {
	return &read{x: x}
}

type read struct {
	x uint8
}

func (i read) String() string {
	return fmt.Sprintf("READ V0-V%x", i.x)
}

func (i read) Execute(c *Chip8) {
	high := uint16(i.x + 1)
	copy(c.registers[:high], c.memory[c.index:c.index+high])
	c.index += high
}
