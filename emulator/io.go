package emulator

import "C"
import (
	"bufio"
	"fmt"
	"os"
)

func ioIn8(address uint16) uint8 {
	switch address {
	case 0x03f8:
		return getchar()
	default:
		return 0
	}
}

func getchar() byte {
	reader := bufio.NewReader(os.Stdin)
	ch, _ := reader.ReadByte()
	return ch
}

func ioOut8(address uint16, value uint8) {
	switch address {
	case 0x03f8:
		putchar(value)
	}
}

func putchar(ch byte) {
	fmt.Printf("%c", ch)
}
