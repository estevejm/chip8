package chip8

import (
	"fmt"
)

type Timer struct {
	rate      uint
	tps       uint
	count     float64
	increment float64
	value     uint8
}

// NewTimer create timer that counts down at a given hertz rate
func NewTimer(tps uint, rate uint) *Timer {
	timer := &Timer{rate: rate}
	timer.SetTPS(tps)
	return timer
}

func (t *Timer) SetTPS(tps uint) {
	t.tps = tps
	t.increment = float64(t.rate) / float64(t.tps)
}

func (t *Timer) GetValue() uint8 {
	return t.value
}

func (t *Timer) SetValue(value uint8) {
	t.value = value
	t.count = 0
}

func (t *Timer) Update() {
	if t.value == 0 {
		return
	}

	for t.count >= 1 && t.value > 0 {
		t.value--
		t.count--
	}

	t.count += t.increment
}

func (t *Timer) String() string {
	return fmt.Sprintf("%02x", t.value)
}
