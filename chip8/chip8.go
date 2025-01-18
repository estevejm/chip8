package chip8

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	memoryLocations  = 0x1000
	programStart     = 0x200
	instructionBytes = 2
	registerCount    = 16
	flagRegister     = 0xF
	stackLevels      = 16
	screenWidth      = 64
	screenHeight     = 32
	bytePixels       = 8
	timerRateHz      = 60
)

var (
	pixelColorOn  = color.White
	pixelColorOff = color.Black
)

func hexdump8(b uint8) string {
	return fmt.Sprintf("%02x", b)
}
func hexdump16(b uint16) string {
	return fmt.Sprintf("%04x", b)
}

type Memory [memoryLocations]uint8

func (m *Memory) Write(start uint16, bytes []byte) {
	// TODO: ensure data copied in bounds -> check data len
	copy(m[start:], bytes)
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

type Registers [registerCount]uint8

func (r Registers) String() string {
	var sb strings.Builder

	for i, b := range r {
		sb.WriteString(fmt.Sprintf("%x:%s ", i, hexdump8(b)))
	}

	return strings.TrimSpace(sb.String())
}

type Stack [stackLevels]uint16

func (s Stack) String() string {
	var sb strings.Builder

	for i, b := range s {
		sb.WriteString(fmt.Sprintf("%x:%s ", i, hexdump16(b)))
	}

	return strings.TrimSpace(sb.String())
}

type Screen [screenHeight][screenWidth]bool

type Chip8 struct {
	log            *slog.Logger
	memory         Memory
	registers      Registers
	stack          Stack
	stackPointer   uint8
	programCounter uint16
	index          uint16
	delayTimer     *Timer
	soundTimer     *Timer
	input          *Input
	screen         Screen
	sound          *Sound
}

func NewChip8(tps uint, log *slog.Logger) *Chip8 {
	memory := Memory{}
	memory.Write(fontStartMemoryAddress, font[:])

	return &Chip8{
		log:            log,
		memory:         memory,
		registers:      Registers{},
		stack:          Stack{},
		stackPointer:   0,
		programCounter: programStart,
		index:          0,
		delayTimer:     NewTimer(tps, timerRateHz),
		soundTimer:     NewTimer(tps, timerRateHz),
		input:          NewInput(),
		screen:         Screen{},
		sound:          NewSound(),
	}
}

func (c *Chip8) LoadROMFromBytes(b []byte) error {
	return c.LoadROM(bytes.NewReader(b))
}

func (c *Chip8) LoadROMFromPath(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}

	return c.LoadROM(f)
}

func (c *Chip8) LoadROM(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	c.memory.Write(programStart, data)
	c.log.Info("ROM loaded", slog.Int("bytes", len(data)))
	//println(c.memory.String())
	return nil
}

func (c *Chip8) Update() error {
	wait := c.input.Detect()
	c.log.Info("input:", slog.String("V", c.input.String()))

	if !wait {
		opcode := c.fetch()
		c.log.Info("fetch:", slog.String("PC", hexdump16(c.programCounter)), slog.String("opcode", hexdump16(opcode)))

		c.incrementProgramCounter()

		instruction, ok := c.decode(opcode)
		if !ok {
			return fmt.Errorf("invalid opcode: %s", hexdump16(opcode))
		}
		c.log.Info("decode: " + instruction.String())

		instruction.Execute(c)
	}

	c.delayTimer.Update()
	c.soundTimer.Update()

	c.outputSound()

	c.log.Info(
		"execute:",
		slog.String("PC", hexdump16(c.programCounter)),
		slog.String("I", hexdump16(c.index)),
		slog.Any("V", c.registers),
		slog.Int("SP", int(c.stackPointer)),
		slog.Any("S", c.stack),
		slog.Any("DT", c.delayTimer),
		slog.Any("ST", c.soundTimer),
	)

	c.log.Info(
		"bench:",
		slog.Int("TPS", int(ebiten.ActualTPS())),
		slog.Int("FPS", int(ebiten.ActualFPS())),
	)

	return nil
}

func (c *Chip8) outputSound() {
	if c.soundTimer.GetValue() == 0 {
		c.sound.Pause()
	} else {
		c.sound.Play()
	}
}

func (c *Chip8) fetch() uint16 {
	// big-endian
	return uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])
}

func (c *Chip8) incrementProgramCounter() {
	// TODO: handle PC > 4096 / 0x1000 (12 bits). 2 options: PC overflow (error) or wrap (modulo)
	// TODO: also check PC < 521 / 0x200 (program start)
	c.programCounter += instructionBytes
}

func (c *Chip8) decode(opcode uint16) (Instruction, bool) {
	return decode(opcode)
}

func (c *Chip8) Draw(screen *ebiten.Image) {
	// TODO: write pixels only if there are changes
	pixelColor := func(pixelOn bool) color.Color {
		if pixelOn {
			return pixelColorOn
		}
		return pixelColorOff
	}

	for i := range c.screen {
		for j := range c.screen[i] {
			screen.Set(j, i, pixelColor(c.screen[i][j]))
		}
	}
}

func (c *Chip8) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return screenWidth, screenHeight
}
