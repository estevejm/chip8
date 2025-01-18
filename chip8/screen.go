package chip8

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth      = 64
	screenHeight     = 32
	screenMultiplier = 16
	bytePixels       = 8
)

var (
	pixelColorOn  = color.White
	pixelColorOff = color.Black
)

// TODO: store frame buffer in byte array instead of bidimensional boolean array
type Screen struct {
	buffer [screenHeight][screenWidth]bool
}

func NewScreen() *Screen {
	return &Screen{buffer: [screenHeight][screenWidth]bool{}}
}

func (s *Screen) Layout() (w, h int) {
	return screenWidth, screenHeight
}

func (s *Screen) Clear() {
	for y := range s.buffer {
		for x := range s.buffer[y] {
			s.buffer[y][x] = false
		}
	}
}

func (s *Screen) Draw(image *ebiten.Image) {
	// TODO: write pixels only if there are changes
	for i := range s.buffer {
		for j := range s.buffer[i] {
			image.Set(j, i, pixelColor(s.buffer[i][j]))
		}
	}
}

func (s *Screen) Get(x, y int) bool {
	return s.buffer[y][x]
}

func (s *Screen) Set(x, y int, v bool) {
	s.buffer[y][x] = v
}

func pixelColor(pixelOn bool) color.Color {
	if pixelOn {
		return pixelColorOn
	}
	return pixelColorOff
}
