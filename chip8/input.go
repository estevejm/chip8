package chip8

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const keyCount = 16

// 1 2 3 C     |\   1 2 3 4
// 4 5 6 D  ---- \  Q W E R
// 7 8 9 E  ---- /  A S D F
// A 0 B F     |/   Z X C V
var keyMap = [keyCount]ebiten.Key{
	ebiten.KeyX,
	ebiten.Key1,
	ebiten.Key2,
	ebiten.Key3,
	ebiten.KeyQ,
	ebiten.KeyW,
	ebiten.KeyE,
	ebiten.KeyA,
	ebiten.KeyS,
	ebiten.KeyD,
	ebiten.KeyZ,
	ebiten.KeyC,
	ebiten.Key4,
	ebiten.KeyR,
	ebiten.KeyF,
	ebiten.KeyV,
}

type Input struct {
	keys         [keyCount]bool // true if pressed
	waitCallback func(uint8)
}

func NewInput() *Input {
	return &Input{
		keys:         [keyCount]bool{},
		waitCallback: nil,
	}
}

// Detect return true waiting for key press + release, false otherwise
func (i *Input) Detect() bool {
	buf := make([]ebiten.Key, 0, keyCount)
	if i.isWaiting() {
		detected := inpututil.AppendJustReleasedKeys(buf)
		if len(detected) == 0 {
			return true
		}

		index := slices.Index(keyMap[:], detected[0]) // TODO: use map instead
		i.waitCallback(uint8(index))
		i.waitCallback = nil

		return false
	}

	detected := inpututil.AppendPressedKeys(buf)
	for j, key := range keyMap {
		i.keys[j] = slices.Contains(detected, key) // TODO: use map instead
	}

	return false
}

func (i *Input) Wait(callback func(uint8)) {
	i.waitCallback = callback
}

func (i *Input) isWaiting() bool {
	return i.waitCallback != nil
}

func (i *Input) String() string {
	var sb strings.Builder

	for j, b := range i.keys {
		sb.WriteString(fmt.Sprintf("%x:%s ", j, boolString(b)))
	}

	sb.WriteString(fmt.Sprintf("wait:%s", boolString(i.isWaiting())))

	return strings.TrimSpace(sb.String())
}

func boolString(pressed bool) string {
	if pressed {
		return "Y"
	}

	return "N"
}
