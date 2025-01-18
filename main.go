package main

import (
	"chip8/chip8"
	"flag"
	"log/slog"
	"os"
)

const (
	defaultTPS = 1000
	defaultRom = "roms/1-chip8-logo.ch8"
)

var (
	tps int
	rom string
)

func main() {
	flag.IntVar(&tps, "tps", defaultTPS, "ticks per second (clock Hz)")
	flag.StringVar(&rom, "rom", defaultRom, "rom path")
	flag.Parse()

	log := slog.Default()
	log.Info("CHIP-8 starting...", slog.Int("tps", tps))

	emulator := chip8.NewChip8(uint(tps), log)

	if err := emulator.LoadFont(); err != nil {
		log.Error("Failed to load font")
		os.Exit(1)
	}

	romReader, err := os.Open(rom)
	if err != nil {
		log.Error("Failed to open ROM file")
		os.Exit(1)
	}

	if err := emulator.LoadROM(romReader); err != nil {
		log.Error("Failed to load ROM file")
		os.Exit(1)
	}

	if err := emulator.Run(); err != nil {
		log.Error(err.Error())
		os.Exit(2)
	}

	log.Info("CHIP-8 stopping...")
	os.Exit(0)
}
