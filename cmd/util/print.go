package util

import (
	"encoding/hex"
	"strings"

	ae "../../aemulari"
)

const LineSeparator string = strings.Repeat("-", 80)

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

func PrintHexDump(header string, addr uint64, data []byte) {
	if len(header) != 0 {
		fmt.Println(header)
		fmt.Println(LineSeparator)
	}

	fmt.Println(hex.Dump(data))
}