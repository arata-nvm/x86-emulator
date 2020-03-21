package main

import (
	"fmt"
	. "github.com/arata-nvm/x86-emulator/emulator"
	"os"
)

const MemorySize = 1024 * 1024

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: px86 filename")
		os.Exit(1)
	}
	emu := NewEmulator(MemorySize, 0x7c00, 0x7c00)
	emu.ReadBinary(os.Args[1])
	emu.Execute()
	emu.DumpRegisters()
}
