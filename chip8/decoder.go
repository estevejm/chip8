package chip8

func decode(opcode uint16) (Instruction, bool) {
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0xFF {
		case 0xE0:
			return ClearScreen(), true
		case 0xEE:
			return Return(), true
		}
	case 0x1000:
		return Jump(opcode), true
	case 0x2000:
		return Call(opcode), true
	case 0x3000:
		return SkipEqual(opcode), true
	case 0x4000:
		return SkipNotEqual(opcode), true
	case 0x5000:
		return SkipEqualRegister(opcode), true
	case 0x6000:
		return Load(opcode), true
	case 0x7000:
		return Add(opcode), true
	case 0x8000:
		switch opcode & 0xF {
		case 0x0:
			return LoadRegister(opcode), true
		case 0x1:
			return Or(opcode), true
		case 0x2:
			return And(opcode), true
		case 0x3:
			return Xor(opcode), true
		case 0x4:
			return AddRegister(opcode), true
		case 0x5:
			return SubRegister(opcode), true
		case 0x6:
			return ShiftRight(opcode), true
		case 0x7:
			return ReverseSubRegister(opcode), true
		case 0xE:
			return ShiftLeft(opcode), true
		}
	case 0x9000:
		return SkipNotEqualRegister(opcode), true
	case 0xA000:
		return LoadIndex(opcode), true
	case 0xB000:
		return JumpRegister0(opcode), true
	case 0xC000:
		return Random(opcode), true
	case 0xD000:
		return DrawSprite(opcode), true
	case 0xE000:
		switch opcode & 0xFF {
		case 0x9E:
			return SkipPressed(opcode), true
		case 0xA1:
			return SkipNotPressed(opcode), true
		}
	case 0xF000:
		switch opcode & 0xFF {
		case 0x07:
			return LoadRegisterDelayTimer(opcode), true
		case 0x0A:
			return WaitKey(opcode), true
		case 0x15:
			return LoadDelayTimerRegister(opcode), true
		case 0x18:
			return LoadSoundTimerRegister(opcode), true
		case 0x1E:
			return AddIndex(opcode), true
		case 0x29:
			return LoadDigitIndex(opcode), true
		case 0x33:
			return BCD(opcode), true
		case 0x55:
			return Write(opcode), true
		case 0x65:
			return Read(opcode), true
		}
	}

	return NoOperation(), false
}
