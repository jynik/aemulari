package main

import (
	"flag"
	"fmt"
	"os"

	"../../internal/cmdline"
	"../../internal/debugger"
	mylog "../../internal/log"
	"../../internal/ui"
	"github.com/op/go-logging"
)

type Flags struct {
	count    int64               // Instruction count
	showRegs bool                // Show register values after execution
	dumpMem  debugger.MemRegions // Memory regions to display after execution
}

func initFlags(f *Flags) {
	f.dumpMem = make(debugger.MemRegions)

	flag.Int64Var(&f.count, "n", -1, "Execute only the specified number of instructions.")
	flag.BoolVar(&f.showRegs, "R", false, "Show register values after execution.")
	flag.Var(&f.dumpMem, "M", "Show specified memory regoin after execution.")
}

var log = logging.MustGetLogger("")

func main() {
	var ret int = 0
	var flags Flags

	mylog.Init()

	cmdline.InitCommonFlags()
	initFlags(&flags)

	cfg, err := cmdline.Parse()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	dbg, err := debugger.New(cfg)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	exception, err := dbg.Step(flags.count)
	if err != nil {
		log.Error(err)
	}

	if exception.Occurred() {
		fmt.Println("Execution halted due to exception: " + exception.String())
	}

	// Print final register state, if requested
	if flags.showRegs {
		rvs, err := dbg.ReadRegAll()
		if err != nil {
			log.Error(err)
		} else {
			fmt.Println()
			ui.PrintRegisters(rvs)
		}
	}

	// Print final state of memory regions, if requested
	for _, m := range flags.dumpMem {
		base, size := m.Region()
		data, err := dbg.ReadMem(base, size)
		if err != nil {
			log.Error(err)
			break
		} else {
			ui.PrintMemory(m.Name(), base, data)
		}
	}

	// Ensure output files are written
	dbg.Close()
	os.Exit(ret)
}
