package main

import (
	"errors"
	"fmt"
	"os"

	ae "../../aemulari"
	"../cmdline"
	"../util"
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

func print_registers(args cmdline.FlagMap, dbg *ae.Debugger) {
	if args.Contains("print-regs") {
		regs, err := dbg.ReadRegAll()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Println()
			switch args.GetString("print-regs", "pretty") {
			case "list":
				for _, r := range regs {
					fmt.Println(r.String())
				}
				fmt.Println()
			default:
				util.PrettyPrintRegisters(regs)
			}
		}
	}
}

func print_hexdumps(args cmdline.FlagMap, dbg *ae.Debugger) {
	return
}

func step(args cmdline.FlagMap, dbg *ae.Debugger) (ae.Exception, error) {
	var ex ae.Exception

	count, err := args.GetInt64("instr-count")
	if err != nil {
		return ex, errors.New("Invalid instruction count: " + err.Error())
	} else if count < 1 {
		return ex, errors.New("Instruction count must be >= 1")
	}

	return dbg.Step(count)
}

func main() {
	var exception ae.Exception
	var err error

	supportedFlags := cmdline.SupportedFlags{
		cmdline.Flag_arch,
		cmdline.Flag_reg,
		cmdline.Flag_mem,
		cmdline.Flag_instrcount,
		cmdline.Flag_breakpoint,
		cmdline.Flag_printRegs,
		cmdline.Flag_hexdump,
	}

	// Fetch an initialized debugger and any unhandled args.
	args, _, dbg := cmdline.Parse(supportedFlags, usageText)

	// Execute our program
	if args.Contains("instr-count") {
		exception, err = step(args, dbg)
	} else {
		exception, err = dbg.Continue()
	}

	if err == nil {
		if exception.Occurred() {
			fmt.Println("Execution terminated due to exception: %s\n", exception.String())
		}

		// Output information requested by cmdline args
		print_registers(args, dbg)
		print_hexdumps(args, dbg)
	} else {
		fmt.Fprintln(os.Stderr, err)
	}

	dbg.Close()
}
