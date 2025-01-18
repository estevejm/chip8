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

type Screen struct {
	buffer *ebiten.Image
}

func NewScreen() *Screen {
	return &Screen{
		buffer: ebiten.NewImage(screenWidth, screenHeight),
	}
}

func (s *Screen) Layout() (w, h int) {
	return screenWidth, screenHeight
}

func (s *Screen) Clear() {
	s.buffer.Fill(pixelColorOff)
}

func (s *Screen) Get(x, y int) bool {
	// TODO: reduce calls https://ebitengine.org/en/documents/performancetips.html#Don't_call_(*Image).At_too_much
	return pixelOn(s.buffer.At(x, y))
}

func (s *Screen) Set(x, y int, on bool) {
	s.buffer.Set(x, y, pixelColor(on))
}

func (s *Screen) Draw(image *ebiten.Image) {
	// TODO: Use BlendXor, if we still can detect collisions https://ebitengine.org/en/examples/blend.html
	image.DrawImage(s.buffer, nil)
}

func pixelColor(on bool) color.Color {
	if on {
		return pixelColorOn
	}
	return pixelColorOff
}

func pixelOn(c color.Color) bool {
	r1, g1, b1, a1 := c.RGBA()
	r2, g2, b2, a2 := pixelColorOn.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}
