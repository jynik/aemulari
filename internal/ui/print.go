package ui

import (
	"encoding/hex"
	"fmt"
	"strings"

	"../arch"
)

var linesep string = strings.Repeat("-", 80)

func PrintRegisters(rvs []arch.RegisterValue) {
	var rv arch.RegisterValue
	var i int

	fmt.Println(" Registers\n" + linesep)
	for i, rv = range rvs {
		fmt.Printf("%s    ", &rv)
		if (i+1)%3 == 0 {
			fmt.Println()
		}
	}

	if (i+1)%3 != 0 {
		fmt.Println()
	}
	fmt.Println()
}

func PrintMemory(name string, addr uint64, data []byte) {
	fmt.Printf(" Memory at 0x%08x (%s)\n%s\n%s\n", addr, name, linesep, hex.Dump(data))
}
