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
	stackLevels      = 16
	inputKeys        = 16
	screenWidth      = 64
	screenHeight     = 32
	bytePixels       = 8
)

var (
	pixelColorOn  = color.White
	pixelColorOff = color.Black
)

func hexdump8(b byte) string {
	return fmt.Sprintf("%02x", b)
}
func hexdump16(b uint16) string {
	return fmt.Sprintf("%04x", b)
}

type Memory [memoryLocations]byte

func (m Memory) Hexdump() string {
	const bytesPerRow = 16
	var sb strings.Builder
	sb.Grow(len(m))

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

type Registers [registerCount]byte

type Screen [screenHeight][screenWidth]bool

func (s *Screen) Clear() {
	for i := range s {
		for j := range s[i] {
			s[i][j] = false
		}
	}
}

type Chip8 struct {
	log            *slog.Logger
	memory         Memory
	registers      Registers
	screen         Screen
	programCounter uint16
	index          uint16
}

func NewChip8(log *slog.Logger) *Chip8 {
	return &Chip8{
		log:            log,
		memory:         Memory{},
		registers:      Registers{},
		screen:         Screen{},
		programCounter: programStart,
		index:          0,
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

	// TODO: ensure data copied in bounds -> check data len
	copy(c.memory[programStart:], data)
	c.log.Info("ROM loaded", slog.Int("bytes", len(data)))
	//println(c.memory.Hexdump())
	return nil
}

func (c *Chip8) Update() error {
	opcode := c.fetch()
	c.log.Info("fetch:", slog.String("PC", hexdump16(c.programCounter)), slog.String("opcode", hexdump16(opcode)))

	c.incrementProgramCounter()

	instruction, ok := c.decode(opcode)
	if !ok {
		return fmt.Errorf("invalid opcode: %s", hexdump16(opcode))
	}
	c.log.Info("decode: " + instruction.String())

	instruction.Execute(c)
	c.log.Info(
		"execute:",
		slog.String("PC", hexdump16(c.programCounter)),
		slog.String("I", hexdump16(c.index)),
		slog.Any("V", c.registers),
	)

	return nil
}

func (c *Chip8) fetch() uint16 {
	return uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])
}

func (c *Chip8) incrementProgramCounter() {
	// TODO: handle PC > 4096 / 0x1000 (12 bits). 2 options: PC overflow (error) or wrap (modulo)
	c.programCounter += instructionBytes
}

func (c *Chip8) decode(opcode uint16) (Instruction, bool) {
	switch opcode & 0xF000 {
	case 0x0000:
		return ClearScreen(), true
	case 0x1000:
		return Jump(opcode), true
	case 0x3000:
		return SkipEqual(opcode), true
	case 0x4000:
		return SkipNotEqual(opcode), true
	case 0x5000:
		return SkipEqualRegister(opcode), true
	case 0x6000:
		return Load(opcode), true
	case 0x7000:
		return Add(opcode), true
	case 0x9000:
		return SkipNotEqualRegister(opcode), true
	case 0xA000:
		return LoadIndex(opcode), true
	case 0xD000:
		return DrawSprite(opcode), true
	default:
		return nil, false
	}
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
