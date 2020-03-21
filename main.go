package main

import (
	"fmt"
	"os"
)

type instructionFunc func(*Emulator)

var instructions [256]instructionFunc

const MEMORY_SIZE = 1024 * 1024
const REGISTERS_COUNT = 8

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

func NewEmulator(eip int, esp int) *Emulator {
	emu := &Emulator{
		Registers: make([]uint32, REGISTERS_COUNT),
		Eflags:    0,
		Memory:    make([]byte, MEMORY_SIZE),
		Eip:       uint32(eip),
	}

	emu.Registers[ESP] = uint32(esp)

	return emu
}

func (emu *Emulator) dumpRegisters() {
	for i, reg := range emu.Registers {
		fmt.Printf("%s = %08x\n", RegistersName[i], reg)
	}

	fmt.Printf("EIP = %08x\n", emu.Eip)
}

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

func initInstructions() {
	for i := 0; i < 8; i++ {
		instructions[0xB8+i] = moveR32Imm32
	}
	instructions[0xEB] = shortJump
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: px86 filename")
		os.Exit(1)
	}

	emu := NewEmulator(0x0000, 0x7c00)

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("%sファイルが開けません\n", os.Args[1])
	}
	defer f.Close()

	f.Read(emu.Memory)

	initInstructions()

	for emu.Eip < MEMORY_SIZE {
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

	emu.dumpRegisters()
}
