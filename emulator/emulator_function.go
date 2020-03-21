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
