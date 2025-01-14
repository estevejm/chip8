package chip8

import (
	"fmt"
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

// Keys true if pressed
type Keys [keyCount]bool

func (k *Keys) detectInput() {
	for i, key := range keyMap {
		k[i] = inpututil.KeyPressDuration(key) > 0
	}
}

func (k *Keys) String() string {
	var sb strings.Builder

	for i, b := range k {
		sb.WriteString(fmt.Sprintf("%x:%t ", i, b))
	}

	return strings.TrimSpace(sb.String())
}
