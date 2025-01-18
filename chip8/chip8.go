package chip8

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const instructionBytes = 2

type Chip8 struct {
	log            *slog.Logger
	memory         Memory
	registers      Registers
	stack          Stack
	stackPointer   uint8
	programCounter uint16
	index          uint16
	delayTimer     *Timer
	input          *Input
	screen         Screen
	sound          *Sound
}

func NewChip8(tps uint, log *slog.Logger) *Chip8 {
	return &Chip8{
		log:            log,
		memory:         Memory{},
		registers:      Registers{},
		stack:          Stack{},
		stackPointer:   0,
		programCounter: programStartMemoryAddress,
		index:          0,
		delayTimer:     NewTimer(tps),
		input:          NewInput(),
		screen:         Screen{},
		sound:          NewSound(NewTimer(tps)),
	}
}

func (c *Chip8) LoadFont() error {
	_, err := c.memory.Write(fontStartMemoryAddress, bytes.NewReader(font[:]))
	return err
}

func (c *Chip8) LoadROM(r io.Reader) error {
	n, err := c.memory.Write(programStartMemoryAddress, r)
	if err != nil {
		return err
	}

	c.log.Info("ROM loaded", slog.Int("bytes", n))
	//println(c.memory.String())
	return nil
}

func (c *Chip8) Update() error {
	wait := c.input.Detect()
	c.log.Info("input  :", slog.String("V", c.input.String()))

	if !wait {
		if err := c.Cycle(); err != nil {
			return err
		}
	}

	c.delayTimer.Update()
	c.sound.Update()

	c.log.Info(
		"execute:",
		slog.Int("TPS", int(math.Round(ebiten.ActualTPS()))),
		slog.String("PC", hexdump16(c.programCounter)),
		slog.String("I", hexdump16(c.index)),
		slog.Any("V", c.registers),
		slog.Int("SP", int(c.stackPointer)),
		slog.Any("S", c.stack),
		slog.Any("DT", c.delayTimer),
		slog.Any("ST", c.sound.timer),
	)

	return nil
}

func (c *Chip8) Cycle() error {
	opcode := c.fetch()
	c.log.Info("fetch  :", slog.String("PC", hexdump16(c.programCounter)), slog.String("opcode", hexdump16(opcode)))

	c.incrementProgramCounter()

	instruction, ok := c.decode(opcode)
	if !ok {
		return fmt.Errorf("invalid opcode: %s", hexdump16(opcode))
	}
	c.log.Info("decode :" + instruction.String())

	instruction.Execute(c)

	return nil
}

func (c *Chip8) fetch() uint16 {
	return c.memory.ReadWord(c.programCounter)
}

func (c *Chip8) incrementProgramCounter() {
	// TODO: handle PC > 4096 / 0x1000 (12 bits). 2 options: PC overflow (error) or wrap (modulo)
	// TODO: also check PC < 521 / 0x200 (program start)
	c.programCounter += instructionBytes
}

func (c *Chip8) decode(opcode uint16) (Instruction, bool) {
	return decode(opcode)
}

func (c *Chip8) Draw(image *ebiten.Image) {
	c.screen.Draw(image)
	c.log.Info("draw   :", slog.Int("FPS", int(math.Round(ebiten.ActualFPS()))))
}

func (c *Chip8) Layout(_, _ int) (w, h int) {
	return c.screen.Layout()
}
