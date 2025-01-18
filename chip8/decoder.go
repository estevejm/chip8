package chip8

func decode(encoded uint16) (Instruction, bool) {
	op := encoded & 0xF000

	// 0x0X00
	x := uint8(encoded>>8) & 0xF
	// 0x00Y0
	y := uint8(encoded>>4) & 0xF
	// 0x00NN
	n := uint8(encoded & 0xF)
	// 0x00NN
	nn := uint8(encoded & 0xFF)
	// 0x0NNN
	nnn := encoded & 0xFFF

	switch op {
	case 0x0000:
		switch nn {
		case 0xE0:
			return ClearScreen(), true
		case 0xEE:
			return Return(), true
		}
	case 0x1000:
		return Jump(nnn), true
	case 0x2000:
		return Call(nnn), true
	case 0x3000:
		return SkipEqual(x, nn), true
	case 0x4000:
		return SkipNotEqual(x, nn), true
	case 0x5000:
		return SkipEqualRegister(x, y), true
	case 0x6000:
		return Load(x, nn), true
	case 0x7000:
		return Add(x, nn), true
	case 0x8000:
		switch n {
		case 0x0:
			return LoadRegister(x, y), true
		case 0x1:
			return Or(x, y), true
		case 0x2:
			return And(x, y), true
		case 0x3:
			return Xor(x, y), true
		case 0x4:
			return AddRegister(x, y), true
		case 0x5:
			return SubRegister(x, y), true
		case 0x6:
			return ShiftRight(x, y), true
		case 0x7:
			return ReverseSubRegister(x, y), true
		case 0xE:
			return ShiftLeft(x, y), true
		}
	case 0x9000:
		return SkipNotEqualRegister(x, y), true
	case 0xA000:
		return LoadIndex(nnn), true
	case 0xB000:
		return JumpRegister0(nnn), true
	case 0xC000:
		return Random(x, nn), true
	case 0xD000:
		return DrawSprite(x, y, n), true
	case 0xE000:
		switch nn {
		case 0x9E:
			return SkipPressed(x), true
		case 0xA1:
			return SkipNotPressed(x), true
		}
	case 0xF000:
		switch nn {
		case 0x07:
			return LoadRegisterDelayTimer(x), true
		case 0x0A:
			return WaitKey(x), true
		case 0x15:
			return LoadDelayTimerRegister(x), true
		case 0x18:
			return LoadSoundTimerRegister(x), true
		case 0x1E:
			return AddIndex(x), true
		case 0x29:
			return LoadDigitIndex(x), true
		case 0x33:
			return BCD(x), true
		case 0x55:
			return Write(x), true
		case 0x65:
			return Read(x), true
		}
	}

	return NoOperation(), false
}
