package chip8

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	sampleRate = 48000
	frequency  = 440
)

// from https://ebitengine.org/en/examples/sinewave.html

// stream is an infinite stream of 440 Hz sine wave.
type stream struct {
	pos int64
}

// Read is io.Reader's Read.
//
// Read fills the data with sine wave samples.
func (s *stream) Read(buf []byte) (int, error) {
	const bytesPerSample = 8

	n := len(buf) / bytesPerSample * bytesPerSample

	const length = sampleRate / frequency
	for i := 0; i < n/bytesPerSample; i++ {
		v := math.Float32bits(float32(math.Sin(2 * math.Pi * float64(s.pos/bytesPerSample+int64(i)) / length)))

		buf[8*i] = byte(v)
		buf[8*i+1] = byte(v >> 8)
		buf[8*i+2] = byte(v >> 16)
		buf[8*i+3] = byte(v >> 24)
		buf[8*i+4] = byte(v)
		buf[8*i+5] = byte(v >> 8)
		buf[8*i+6] = byte(v >> 16)
		buf[8*i+7] = byte(v >> 24)
	}

	s.pos += int64(n)
	s.pos %= length * bytesPerSample

	return n, nil
}

// Close is io.Closer's Close.
func (s *stream) Close() error {
	return nil
}

type Sound struct {
	player *audio.Player
	timer  *Timer
}

func NewSound(timer *Timer) *Sound {
	context := audio.NewContext(sampleRate)

	player, _ := context.NewPlayerF32(&stream{})
	player.SetVolume(0)
	player.Play()

	return &Sound{
		player: player,
		timer:  timer,
	}
}

func (s *Sound) SetTimerValue(value uint8) {
	s.timer.SetValue(value)
}

func (s *Sound) Update() {
	s.timer.Update()
	if s.timer.GetValue() == 0 {
		s.player.SetVolume(0)
	} else {
		s.player.SetVolume(1)
	}
}
