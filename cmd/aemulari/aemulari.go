package main

import (
	"path/filepath"
	"fmt"
	"os"
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
	" + Execution terminates when an exception occurs or a when breakpoint is hit.\n" +
	"\n"

func main() {
	params := cmdline.ArgList{ cmdline.Arg_arch }

	if cmdline.HelpRequested(os.Args) {
		fmt.Printf(usageText, filepath.Base(os.Args[0]))
		os.Exit(0)
	}

	args, err := params.Parse(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}


}
