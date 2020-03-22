package emulator

const (
	CarryFlag    = 1
	ZeroFlag     = 1 << 6
	SignFlag     = 1 << 7
	OverflowFlag = 1 << 11
)

func (emu *Emulator) getCode8(index int) uint32 {
	return uint32(emu.Memory[emu.Eip+uint32(index)])
}

func (emu *Emulator) getSignCode8(index int) int32 {
	return int32(int8(emu.Memory[emu.Eip+uint32(index)]))
}

func (emu *Emulator) getCode32(index int) uint32 {
	var ret uint32
	for i := 0; i < 4; i++ {
		ret |= emu.getCode8(index+i) << (i * 8)
	}

	return ret
}

func (emu *Emulator) getSignCode32(index int) int32 {
	return int32(emu.getCode32(index))
}

func (emu *Emulator) getRegister8(index int) uint8 {
	if index < 4 {
		return uint8(emu.Registers[index] & 0xff)
	} else {
		return uint8((emu.Registers[index-4] >> 6) & 0xff)
	}
}

func (emu *Emulator) setRegister8(index int, value uint8) {
	if index < 4 {
		r := emu.Registers[index] & 0xffffff00
		emu.Registers[index] = r | uint32(value)
	} else {
		r := emu.Registers[index-4] & 0xffff00ff
		emu.Registers[index-4] = r | uint32(value<<8)
	}
}

func (emu *Emulator) getRegister32(index int) uint32 {
	return emu.Registers[index]
}

func (emu *Emulator) setRegister32(index int, value uint32) {
	emu.Registers[index] = value
}

func (emu *Emulator) setMemory8(address, value uint32) {
	emu.Memory[address] = byte(value & 0xFF)
}

func (emu *Emulator) setMemory32(address, value uint32) {
	for i := 0; i < 4; i++ {
		emu.setMemory8(address+uint32(i), value>>(i*8))
	}
}

func (emu *Emulator) getMemory8(address uint32) uint32 {
	return uint32(emu.Memory[address])
}

func (emu *Emulator) getMemory32(address uint32) uint32 {
	var ret uint32
	for i := 0; i < 4; i++ {
		ret |= emu.getMemory8(address+uint32(i)) << (8 * i)
	}
	return ret
}
func (emu *Emulator) push32(value uint32) {
	address := emu.getRegister32(ESP) - 4
	emu.setRegister32(ESP, address)
	emu.setMemory32(address, value)
}

func (emu *Emulator) pop32() uint32 {
	address := emu.getRegister32(ESP)
	ret := emu.getMemory32(address)
	emu.setRegister32(ESP, address+4)
	return ret
}

func (emu *Emulator) updateEflagsSub(v1, v2 uint32, result uint64) {
	sign1 := v1 >> 31
	sign2 := v2 >> 31
	signr := (result >> 31) & 1

	emu.setCarry((result >> 32) != 0)
	emu.setZero(result == 0)
	emu.setSign(signr != 0)
	emu.setOverflow(sign1 != sign2 && uint64(sign1) != signr)
}

func (emu *Emulator) setCarry(isCarry bool) {
	if isCarry {
		emu.Eflags |= CarryFlag
	} else {
		emu.Eflags &^= CarryFlag
	}
}

func (emu *Emulator) setZero(isZero bool) {
	if isZero {
		emu.Eflags |= ZeroFlag
	} else {
		emu.Eflags &^= ZeroFlag
	}
}

func (emu *Emulator) setSign(isSign bool) {
	if isSign {
		emu.Eflags |= SignFlag
	} else {
		emu.Eflags &^= SignFlag
	}
}

func (emu *Emulator) setOverflow(isOverflow bool) {
	if isOverflow {
		emu.Eflags |= OverflowFlag
	} else {
		emu.Eflags &^= OverflowFlag
	}
}

func (emu *Emulator) isCarry() bool {
	return emu.Eflags&CarryFlag != 0
}

func (emu *Emulator) isZero() bool {
	return emu.Eflags&ZeroFlag != 0
}

func (emu *Emulator) isSign() bool {
	return emu.Eflags&SignFlag != 0
}

func (emu *Emulator) isOverflow() bool {
	return emu.Eflags&OverflowFlag != 0
}
