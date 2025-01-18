package chip8

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 64
	screenHeight = 32
	bytePixels   = 8
)

var (
	pixelColorOn  = color.White
	pixelColorOff = color.Black
)

// TODO: store frame buffer in byte array instead of bidimensional boolean array
type Screen [screenHeight][screenWidth]bool

func (s *Screen) Layout() (w, h int) {
	return screenWidth, screenHeight
}

func (s *Screen) Clear() {
	for y := range s {
		for x := range s[y] {
			s[y][x] = false
		}
	}
}

func (s *Screen) Draw(image *ebiten.Image) {
	// TODO: write pixels only if there are changes
	for i := range s {
		for j := range s[i] {
			image.Set(j, i, pixelColor(s[i][j]))
		}
	}
}

func pixelColor(pixelOn bool) color.Color {
	if pixelOn {
		return pixelColorOn
	}
	return pixelColorOff
}