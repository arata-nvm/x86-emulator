package emulator

import (
	"fmt"
	"os"
)

type instructionFunc func(*Emulator)

var instructions [256]instructionFunc

func initInstructions() {
	instructions[0x01] = addRm32R32
	for i := 0; i < 8; i++ {
		instructions[0x50+i] = pushR32
	}
	for i := 0; i < 8; i++ {
		instructions[0x58+i] = popR32
	}
	instructions[0x68] = pushImm32
	instructions[0x6A] = pushImm8
	instructions[0x83] = code83
	instructions[0x89] = movRm32R32
	instructions[0x8B] = movR32Rm32
	for i := 0; i < 8; i++ {
		instructions[0xB8+i] = movR32Imm32
	}
	instructions[0xC3] = ret
	instructions[0xC7] = movRm32Imm32
	instructions[0xC9] = leave
	instructions[0xE8] = callRel32
	instructions[0xE9] = nearJump
	instructions[0xEB] = shortJump
	instructions[0xFF] = codeFF
}

func movR32Imm32(emu *Emulator) {
	reg := emu.getCode8(0) - 0xB8
	value := emu.getCode32(1)
	emu.Registers[reg] = value
	emu.Eip += 5
}

func movRm32Imm32(emu *Emulator) {
	emu.Eip += 1
	modrm := emu.parseModrm()
	value := emu.getCode32(0)
	emu.Eip += 4
	emu.setRm32(modrm, value)
}

func movRm32R32(emu *Emulator) {
	emu.Eip += 1
	modrm := emu.parseModrm()
	r32 := emu.getR32(modrm)
	emu.setRm32(modrm, r32)
}

func movR32Rm32(emu *Emulator) {
	emu.Eip += 1
	modrm := emu.parseModrm()
	rm32 := emu.getRm32(modrm)
	emu.setR32(modrm, rm32)
}

func addRm32R32(emu *Emulator) {
	emu.Eip += 1
	modrm := emu.parseModrm()
	r32 := emu.getR32(modrm)
	rm32 := emu.getRm32(modrm)
	emu.setRm32(modrm, rm32+r32)
}

func addRm32Imm8(emu *Emulator, modrm *ModRM) {
	rm32 := emu.getRm32(modrm)
	imm8 := emu.getSignCode8(0)
	emu.Eip += 1
	emu.setRm32(modrm, rm32+uint32(imm8))
}

func subRm32Imm8(emu *Emulator, modrm *ModRM) {
	rm32 := emu.getRm32(modrm)
	imm8 := emu.getSignCode8(0)
	emu.Eip += 1
	emu.setRm32(modrm, rm32-uint32(imm8))
}

func code83(emu *Emulator) {
	emu.Eip += 1
	modrm := emu.parseModrm()

	switch modrm.Opecode {
	case 0:
		addRm32Imm8(emu, modrm)
	case 5:
		subRm32Imm8(emu, modrm)
	default:
		fmt.Printf("not implemented: 83 /%d\n", modrm.Opecode)
		os.Exit(1)
	}
}

func incRm32(emu *Emulator, modrm *ModRM) {
	value := emu.getRm32(modrm)
	emu.setRm32(modrm, value+1)
}

func codeFF(emu *Emulator) {
	emu.Eip += 1
	modrm := emu.parseModrm()

	switch modrm.Opecode {
	case 0:
		incRm32(emu, modrm)
	default:
		fmt.Printf("not implemented: FF /%d\n", modrm.Opecode)
		os.Exit(1)
	}
}

func pushR32(emu *Emulator) {
	reg := emu.getCode8(0) - 0x50
	emu.push32(emu.getRegister32(int(reg)))
	emu.Eip += 1
}

func pushImm32(emu *Emulator) {
	value := emu.getCode32(1)
	emu.push32(value)
	emu.Eip += 5
}

func pushImm8(emu *Emulator) {
	value := emu.getCode8(1)
	emu.push32(value)
	emu.Eip += 2
}

func popR32(emu *Emulator) {
	reg := emu.getCode8(0) - 0x58
	emu.setRegister32(int(reg), emu.pop32())
	emu.Eip += 1
}

func callRel32(emu *Emulator) {
	diff := emu.getSignCode32(1)
	emu.push32(emu.Eip + 5)
	emu.Eip += uint32(diff + 5)
}

func ret(emu *Emulator) {
	emu.Eip = emu.pop32()
}

func leave(emu *Emulator) {
	ebp := emu.getRegister32(EBP)
	emu.setRegister32(ESP, ebp)
	emu.setRegister32(EBP, emu.pop32())
	emu.Eip += 1
}

func shortJump(emu *Emulator) {
	diff := emu.getSignCode8(1)
	emu.Eip += uint32(diff + 2)
}

func nearJump(emu *Emulator) {
	diff := emu.getSignCode32(1)
	emu.Eip += uint32(diff + 5)
}
