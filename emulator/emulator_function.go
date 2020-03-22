package emulator

func (emu *Emulator) getCode8(index int) uint32 {
	return uint32(emu.Memory[emu.Eip+uint32(index)])
}

func (emu *Emulator) getSignCode8(index int) int32 {
	return int32(emu.Memory[emu.Eip+uint32(index)])
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
