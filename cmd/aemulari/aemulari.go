package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	ae "../../aemulari.v0"
	"../internal/cmdline"
	"../internal/util"
)

const version = "v0.1.0"

type hexdumpRequest struct {
	name         string
	addr, length uint64
}

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
	"\n" +
	"Examples:\n" +
	"  Run myprogram.bin and then print the state of registers upon termination\n" +
	"  caused by the breakpoint at address 0x10214 or an unexpected exception.\n" +
	"    aemulari -m code:0x10000:0x1000:rx:./myprogram.bin -b 0x10214 -R\n" +
	"\n" +
	"  Execute myprogram.bin until it terminates and then write the contents of\n" +
	"  the \"mydata\" region to a file named mydata.bin.\n" +
	"    aemulari -m code:0x48000000:0x4000:rx:./myprogram.bin \\\n" +
	"      -m mydata:0x80000000:0x4000:rw::./mydata.bin\n" +
	"\n"

// Output the final states of registers, if requested  to do so
func print_registers(args cmdline.ArgMap, dbg *ae.Debugger) {
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

func parseHexdumpRequests(args cmdline.ArgMap, dbg *ae.Debugger) ([]hexdumpRequest, error) {
	var requests []hexdumpRequest
	for _, region := range args.GetStrings("hexdump") {
		if dbg.IsMapped(region) {
			requests = append(requests, hexdumpRequest{name: region})
		} else {
			fields := strings.Split(region, ":")
			if len(fields) != 2 {
				return requests, errors.New("Invalid memory region specification: " + region)
			}
			addr, err := strconv.ParseUint(fields[0], 0, 64)
			if err != nil {
				return requests, fmt.Errorf("Invalid memory region address: %s", addr)
			}

			length, err := strconv.ParseUint(fields[1], 0, 64)
			if err != nil {
				return requests, fmt.Errorf("Invalid memory region address: %s", addr)
			}

			requests = append(requests, hexdumpRequest{addr: addr, length: length})
		}
	}

	return requests, nil
}

// Print hex dumps of memory regions, if asked to do so.
func print_hexdumps(regions []hexdumpRequest, dbg *ae.Debugger) {
	var failures []string
	var header string
	var addr uint64
	var data []byte
	var err error

	for _, r := range regions {
		name := ""
		if len(r.name) != 0 {
			name = "(" + r.name + ")"
			addr, data, err = dbg.ReadMemRegion(r.name)
			if err != nil {
				f := fmt.Sprintf("Failed to read memory region named \"%s\": %s",
					r.name, err.Error())
				failures = append(failures, f)
				continue
			}
		} else {
			addr = r.addr
			data, err = dbg.ReadMem(addr, r.length)
			if err != nil {
				f := fmt.Sprintf("Failed to read %d bytes of memory at 0x%08x: %s",
					r.length, r.addr, err.Error())
				failures = append(failures, f)
				continue
			}
		}

		header = fmt.Sprintf(" Memory Region at 0x%08x %s", addr, name)
		util.PrintHexDump(header, addr, data)
	}

	// Print errors at the end
	for _, f := range failures {
		fmt.Fprintln(os.Stderr, f)
	}
}

// Step `instr-count` instructions
func step(args cmdline.ArgMap, dbg *ae.Debugger) (ae.Exception, error) {
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

	// Finish remaining argument parsing tasks
	hexdumpRequests, err := parseHexdumpRequests(args, dbg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		goto cleanup
	}

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
		print_hexdumps(hexdumpRequests, dbg)
	} else {
		fmt.Fprintln(os.Stderr, err)
	}

cleanup:
	dbg.Close()
}
