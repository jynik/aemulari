package main

import (
	"fmt"
	"os"

	"../internal/cmdline"
	"./ui"
)

const version = "v0.1.0"

var usageText string = "" +
	"aemulari - Terminal GUI frontend for aemulari debugger (" + version + ")\n" +
	"Usage: %s [options]\n" +
	"\n" +
	"Options:\n" +
	cmdline.FlagStr_arch +
	cmdline.FlagStr_regs +
	cmdline.FlagStr_mem +
	cmdline.FlagStr_breakpoint +
	cmdline.FlagStr_help +
	cmdline.Details_arch +
	cmdline.Details_mem +
	cmdline.Notes +
	"  - Available GUI commands can be viewed by running the \"help\" command.\n" +
	"\n" +
	"Examples:\n" +
	"  Run myprogram.bin with memory at 0x48000 initialized with the contents\n" +
	"  of a mydata.bin file.\n"+
	"    aemulari -m code:0x10000:0x1000:rx:./myprogram.bin \\\n" +
	"      -m mydata:0x48000:0x200:rwx:./mydata.bin\n" +
	"\n" +
	"  Execute myprogram.bin with the breakpoints set. Note that \n" +
	"  breakpoints can also be set from within the GUI.\n" +
	"    aemulari -m code:0x48000000:0x4000:rx:./myprogram.bin \\\n" +
	"      -m mydata:0x80000000:0x4000:rw::./mydata.bin\n" +
	"\n"

func main() {
	supportedFlags := cmdline.SupportedFlags{
		cmdline.Flag_arch,
		cmdline.Flag_reg,
		cmdline.Flag_mem,
		cmdline.Flag_instrcount,
		cmdline.Flag_breakpoint,
		cmdline.Flag_printRegs,
		cmdline.Flag_hexdump,
	}

	_, arch, dbg := cmdline.Parse(supportedFlags, usageText)

	if gui, err := ui.Create(arch, dbg); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	} else {
		gui.Run()
		gui.Close()
		dbg.Close()
	}
}
