package emulator

import (
	"fmt"
	"os"
)

const RegistersCount = 8

var RegistersName = []string{"EAX", "ECX", "EDX", "EBX", "ESP", "EBP", "ESI", "EDI"}

const (
	EAX = iota
	ECX
	EDX
	EBX
	ESP
	EBP
	ESI
	EDI
)

type Emulator struct {
	Registers []uint32
	Eflags    uint32
	Memory    []byte
	Eip       uint32
}

func NewEmulator(size, eip int, esp int) *Emulator {
	emu := &Emulator{
		Registers: make([]uint32, RegistersCount),
		Eflags:    0,
		Memory:    make([]byte, size),
		Eip:       uint32(eip),
	}

	emu.Registers[ESP] = uint32(esp)

	return emu
}

func (emu *Emulator) DumpRegisters() {
	for i, reg := range emu.Registers {
		fmt.Printf("%s = %08x\n", RegistersName[i], reg)
	}

	fmt.Printf("EIP = %08x\n", emu.Eip)
}

func (emu *Emulator) ReadBinary(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%sファイルが開けません\n", filename)
		os.Exit(1)
	}
	defer f.Close()

	f.Read(emu.Memory[0x7c00:])
}

func (emu *Emulator) Execute() {
	initInstructions()

	for emu.Eip < uint32(len(emu.Memory)) {
		code := emu.getCode8(0)
		fmt.Printf("EIP = %X, Code = %02X\n", emu.Eip, code)

		if instructions[code] == nil {
			fmt.Printf("Not Implemented: %x\n", code)
			break
		}

		instructions[code](emu)

		if emu.Eip == 0x00 {
			fmt.Printf("\n\nend of program.\n\n")
			break
		}
	}

}
