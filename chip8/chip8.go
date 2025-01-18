package chip8

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	fontStartMemoryAddress    = 0x50
	programStartMemoryAddress = 0x200
)

type Chip8 struct {
	log        *slog.Logger
	clockHz    uint
	memory     Memory
	index      uint16
	registers  Registers
	stack      *Stack
	fetcher    *Fetcher
	delayTimer *Timer
	input      *Input
	screen     *Screen
	sound      *Sound
}

func NewChip8(tps uint, log *slog.Logger) *Chip8 {
	return &Chip8{
		log:        log,
		clockHz:    tps,
		memory:     Memory{},
		index:      0,
		registers:  Registers{},
		stack:      NewStack(),
		fetcher:    NewFetcher(programStartMemoryAddress),
		delayTimer: NewTimer(tps),
		input:      NewInput(),
		screen:     NewScreen(),
		sound:      NewSound(NewTimer(tps)),
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

func (c *Chip8) Run() error {
	ebiten.SetWindowSize(screenWidth*screenMultiplier, screenHeight*screenMultiplier)
	ebiten.SetWindowTitle("CHIP-8")
	ebiten.SetTPS(int(c.clockHz))

	return ebiten.RunGame(c)
}

func (c *Chip8) Update() error {
	// TODO: wait using channel + select?
	wait := c.input.Detect()
	c.log.Info("input  :", slog.Any("keys", c.input))

	if !wait {
		if err := c.Cycle(); err != nil {
			return err
		}
	}

	c.delayTimer.Update()
	c.sound.Update()

	c.log.Info(
		"execute:",
		slog.String("PC", hexdump16(c.fetcher.counter)),
		slog.String("I", hexdump16(c.index)),
		slog.Any("V", c.registers),
		slog.Any("S", c.stack),
		slog.Any("DT", c.delayTimer),
		slog.Any("ST", c.sound.timer),
		slog.Int("TPS", int(math.Round(ebiten.ActualTPS()))),
		slog.Int("FPS", int(math.Round(ebiten.ActualFPS()))),
	)

	return nil
}

func (c *Chip8) Cycle() error {
	pc := c.fetcher.counter
	opcode := c.fetcher.Fetch(c)
	c.log.Info("fetch  :", slog.String("PC", hexdump16(pc)), slog.String("opcode", hexdump16(opcode)))

	instruction, ok := decode(opcode)
	if !ok {
		return fmt.Errorf("invalid opcode: %s", hexdump16(opcode))
	}
	c.log.Info("decode : " + instruction.String())

	instruction.Execute(c)

	return nil
}

func (c *Chip8) Draw(image *ebiten.Image) {
	c.screen.Draw(image)
}

func (c *Chip8) Layout(_, _ int) (w, h int) {
	return c.screen.Layout()
}
