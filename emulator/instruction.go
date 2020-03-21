package emulator

type instructionFunc func(*Emulator)

var instructions [256]instructionFunc

func initInstructions() {
	for i := 0; i < 8; i++ {
		instructions[0xB8+i] = moveR32Imm32
	}
	instructions[0xE9] = nearJump
	instructions[0xEB] = shortJump
}

func moveR32Imm32(emu *Emulator) {
	reg := emu.getCode8(0) - 0xB8
	value := emu.getCode32(1)
	emu.Registers[reg] = value
	emu.Eip += 5
}

func shortJump(emu *Emulator) {
	diff := emu.getSignCode8(1)
	emu.Eip += uint32(diff + 2)
}

func nearJump(emu *Emulator) {
	diff := emu.getSignCode32(1)
	emu.Eip += uint32(diff + 5)
}
