package main

import (
	//"strings"

	//ae "../../aemulari"
	"../cmdline"
)

const version = "v0.1.0"

var usageText string = "" +
	"aemulari - Batch execution of the aemulari debugger (" + version + ")\n" +
	"Usage: %s [options]\n" +
	"\n" +
	"Options:\n" +
	cmdline.FlagStr_arch +
	cmdline.FlagStr_regs +
	cmdline.FlagStr_mem +
	cmdline.FlagStr_breakpoint +
	cmdline.FlagStr_printRegs +
	cmdline.FlagStr_printHexdump +
	cmdline.FlagStr_help +
	cmdline.Details_arch +
	cmdline.Details_mem +
	cmdline.Notes +
	" - Execution terminates when an exception occurs or a when breakpoint is hit.\n" +
	"\n"

func main() {
	supportedArgs := cmdline.SupportedArgs{
		cmdline.Arg_arch,
		cmdline.Arg_reg,
		cmdline.Arg_mem,
		cmdline.Arg_breakpoint,
		cmdline.Arg_printRegs,
		cmdline.Arg_hexdump,
	}

	cmdline.Parse(supportedArgs, usageText)
}
