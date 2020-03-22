package emulator

import (
	"fmt"
	"os"
)

type instructionFunc func(*Emulator)

var instructions [256]instructionFunc

func initInstructions() {
	instructions[0x01] = addRm32R32
	instructions[0x3B] = cmpR32Rm32
	for i := 0; i < 8; i++ {
		instructions[0x50+i] = pushR32
	}
	for i := 0; i < 8; i++ {
		instructions[0x58+i] = popR32
	}
	instructions[0x68] = pushImm32
	instructions[0x6A] = pushImm8
	instructions[0x70] = jo
	instructions[0x71] = jno
	instructions[0x72] = jc
	instructions[0x73] = jnc
	instructions[0x74] = jz
	instructions[0x75] = jnz
	instructions[0x78] = js
	instructions[0x79] = jns
	instructions[0x7C] = jl
	instructions[0x7E] = jle
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
	result := uint64(rm32) - uint64(imm8)
	emu.updateEflagsSub(rm32, uint32(imm8), result)
}

func code83(emu *Emulator) {
	emu.Eip += 1
	modrm := emu.parseModrm()

	switch modrm.Opecode {
	case 0:
		addRm32Imm8(emu, modrm)
	case 5:
		subRm32Imm8(emu, modrm)
	case 7:
		cmpRm32Imm8(emu, modrm)
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

func cmpR32Rm32(emu *Emulator) {
	emu.Eip += 1
	modrm := emu.parseModrm()
	r32 := emu.getR32(modrm)
	rm32 := emu.getRm32(modrm)
	result := uint64(r32) - uint64(rm32)
	emu.updateEflagsSub(r32, rm32, result)
}

func cmpRm32Imm8(emu *Emulator, modrm *ModRM) {
	rm32 := emu.getRm32(modrm)
	imm8 := emu.getSignCode8(0)
	emu.Eip += 1
	result := uint64(rm32) - uint64(imm8)
	emu.updateEflagsSub(rm32, uint32(imm8), result)
}

func jc(emu *Emulator) {
	var diff int32
	if emu.isCarry() {
		diff = emu.getSignCode8(1)
	} else {
		diff = 0
	}
	emu.Eip += uint32(diff + 2)
}

func jnc(emu *Emulator) {
	var diff int32
	if emu.isCarry() {
		diff = 0
	} else {
		diff = emu.getSignCode8(1)
	}
	emu.Eip += uint32(diff + 2)
}

func jz(emu *Emulator) {
	var diff int32
	if emu.isZero() {
		diff = emu.getSignCode8(1)
	} else {
		diff = 0
	}
	emu.Eip += uint32(diff + 2)
}

func jnz(emu *Emulator) {
	var diff int32
	if emu.isZero() {
		diff = 0
	} else {
		diff = emu.getSignCode8(1)
	}
	emu.Eip += uint32(diff + 2)
}

func js(emu *Emulator) {
	var diff int32
	if emu.isSign() {
		diff = emu.getSignCode8(1)
	} else {
		diff = 0
	}
	emu.Eip += uint32(diff + 2)
}

func jns(emu *Emulator) {
	var diff int32
	if emu.isSign() {
		diff = 0
	} else {
		diff = emu.getSignCode8(1)
	}
	emu.Eip += uint32(diff + 2)
}

func jo(emu *Emulator) {
	var diff int32
	if emu.isOverflow() {
		diff = emu.getSignCode8(1)
	} else {
		diff = 0
	}
	emu.Eip += uint32(diff + 2)
}

func jno(emu *Emulator) {
	var diff int32
	if emu.isOverflow() {
		diff = 0
	} else {
		diff = emu.getSignCode8(1)
	}
	emu.Eip += uint32(diff + 2)
}

func jl(emu *Emulator) {
	var diff int32
	if emu.isSign() != emu.isOverflow() {
		diff = emu.getSignCode8(1)
	} else {
		diff = 0
	}
	emu.Eip += uint32(diff + 2)
}

func jle(emu *Emulator) {
	var diff int32
	if emu.isZero() || (emu.isSign() != emu.isOverflow()) {
		diff = emu.getSignCode8(1)
	} else {
		diff = 0
	}
	emu.Eip += uint32(diff + 2)
}

func shortJump(emu *Emulator) {
	diff := emu.getSignCode8(1)
	emu.Eip += uint32(diff + 2)
}

func nearJump(emu *Emulator) {
	diff := emu.getSignCode32(1)
	emu.Eip += uint32(diff + 5)
}
