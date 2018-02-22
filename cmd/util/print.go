package util

import (
	"encoding/hex"
	"fmt"
	"strings"

	ae "../../aemulari"
)

var linesep string = strings.Repeat("-", 80)

func PrettyPrintRegisters(regs []ae.Register) {
	var reg ae.Register
	var i int

	fmt.Println(" Registers\n" + linesep)
	for i, reg = range regs {
		fmt.Printf("%s    ", &reg)
		if (i+1)%3 == 0 {
			fmt.Println()
		}
	}

	if (i+1)%3 != 0 {
		fmt.Println()
	}
	fmt.Println()
}

// TODO: Rewrite this - I'd prefer that the address in the dump
// starts at `addr` % 16 rather than 00000000.
func PrintHexDump(header string, addr uint64, data []byte) {
	if len(header) != 0 {
		fmt.Println(header)
		fmt.Println(linesep)
	}

	fmt.Println(hex.Dump(data))
}
