package main

import (
	"chip8/chip8"
	"flag"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth = 1024
	defaultTPS  = 700
	defaultRom  = "roms/1-chip8-logo.ch8"
)

var (
	tps int
	rom string
)

func main() {
	log := slog.Default()
	log.Info("CHIP-8 starting...")

	flag.IntVar(&tps, "tps", defaultTPS, "ticks per second (clock Hz)")
	flag.StringVar(&rom, "rom", defaultRom, "rom path")
	flag.Parse()

	log.Info("CHIP-8 ready", slog.Int("tps", tps))

	ebiten.SetWindowSize(screenWidth, screenWidth/2)
	ebiten.SetWindowTitle("CHIP-8")
	ebiten.SetTPS(tps)

	emulator := chip8.NewChip8(log)
	if err := emulator.LoadROMFromPath(rom); err != nil {
		log.Error("Failed to load ROM file")
		os.Exit(1)
	}

	if err := ebiten.RunGame(emulator); err != nil {
		log.Error(err.Error())
		os.Exit(2)
	}

	log.Info("CHIP-8 stopping...")
	os.Exit(0)
}
