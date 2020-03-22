package emulator

import (
	"fmt"
	"os"
)

type ModRM struct {
	Mod      uint8
	Opecode  uint8
	RegIndex uint8
	Rm       uint8

	Sib    uint8
	Disp8  int8
	Disp32 uint32
}

func (emu *Emulator) parseModrm() *ModRM {
	modrm := &ModRM{}

	code := emu.getCode8(0)
	modrm.Mod = uint8((code & 0xC0) >> 6)
	modrm.Opecode = uint8((code & 0x38) >> 3)
	modrm.RegIndex = modrm.Opecode
	modrm.Rm = uint8(code & 0x07)

	emu.Eip += 1

	if modrm.Mod != 3 && modrm.Rm == 4 {
		modrm.Sib = uint8(emu.getCode8(0))
		emu.Eip += 1
	}

	if (modrm.Mod == 0 && modrm.Rm == 5) || modrm.Mod == 2 {
		modrm.Disp32 = uint32(emu.getSignCode32(0))
		emu.Eip += 4
	} else if modrm.Mod == 1 {
		modrm.Disp8 = int8(emu.getSignCode8(0))
		emu.Eip += 1
	}

	return modrm
}

func (emu *Emulator) setRm32(modrm *ModRM, value uint32) {
	if modrm.Mod == 3 {
		emu.setRegister32(int(modrm.Rm), value)
	} else {
		address := emu.calcMemoryAddress(modrm)
		emu.setMemory32(address, value)
	}
}

func (emu *Emulator) getRm32(modrm *ModRM) uint32 {
	if modrm.Mod == 3 {
		return emu.getRegister32(int(modrm.Rm))
	} else {
		address := emu.calcMemoryAddress(modrm)
		return emu.getMemory32(address)
	}
}

func (emu *Emulator) setR32(modrm *ModRM, value uint32) {
	emu.setRegister32(int(modrm.RegIndex), value)
}

func (emu *Emulator) getR32(modrm *ModRM) uint32 {
	return emu.getRegister32(int(modrm.RegIndex))
}

func (emu *Emulator) calcMemoryAddress(modrm *ModRM) uint32 {
	if modrm.Mod == 0 {
		if modrm.Rm == 4 {
			fmt.Println("not implemented ModRM mod = 0, rm = 4")
			os.Exit(0)
		} else if modrm.Rm == 5 {
			return modrm.Disp32
		} else {
			return emu.getRegister32(int(modrm.Rm))
		}
	} else if modrm.Mod == 1 {
		if modrm.Rm == 4 {
			fmt.Println("not implemented ModRM mod = 1, rm = 4")
			os.Exit(0)
		} else {
			return emu.getRegister32(int(modrm.Rm)) + uint32(modrm.Disp8)
		}
	} else if modrm.Mod == 2 {
		if modrm.Rm == 4 {
			fmt.Println("not implemented ModRM mod = 2, rm = 4")
			os.Exit(0)
		} else {
			return emu.getRegister32(int(modrm.Rm)) + modrm.Disp32
		}
	} else {
		fmt.Println("not implemented ModRM mod = 3")
		os.Exit(0)
	}

	// unreachable
	return 0
}
