package main

import (
	"flag"
	"fmt"
	. "github.com/arata-nvm/x86-emulator/emulator"
	"os"
)

const MemorySize = 1024 * 1024

func main() {
	var (
		quiet = flag.Bool("q", false, "quiet")
	)
	flag.Parse()

	filename := flag.Arg(0)
	if filename == "" {
		fmt.Println("usage: px86 filename")
		os.Exit(1)
	}
	emu := NewEmulator(MemorySize, 0x7c00, 0x7c00)
	emu.ReadBinary(filename)
	emu.Execute(*quiet)
	emu.DumpRegisters()
}
