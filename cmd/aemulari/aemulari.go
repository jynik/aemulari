package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/op/go-logging"

	ae "../../aemulari"
	"../common"
)

type Flags struct {
	count    int64         // Instruction count
	showRegs bool          // Show register values after execution
	dumpMem  ae.MemRegions // Memory regions to display after execution
}

func initFlags(f *Flags) {
	f.dumpMem = make(ae.MemRegions)

	flag.Int64Var(&f.count, "n", -1, "Execute only the specified number of instructions.")
	flag.BoolVar(&f.showRegs, "R", false, "Show register values after execution.")
	flag.Var(&f.dumpMem, "M", "Show specified memory regoin after execution.")
}

var log = logging.MustGetLogger("")

var linesep string = strings.Repeat("-", 80)

func PrintRegisters(regs []ae.Register) {
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

func PrintMemory(name string, addr uint64, data []byte) {
	fmt.Printf(" Memory at 0x%08x (%s)\n%s\n%s\n", addr, name, linesep, hex.Dump(data))
}

func main() {
	var ret int = 0
	var flags Flags
	var exception ae.Exception

	common.InitLogging()

	common.InitCommonFlags()
	initFlags(&flags)

	arch, cfg, err := common.Parse()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	dbg, err := ae.NewDebugger(arch, cfg)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if flags.count < 1 {
		exception, err = dbg.Continue()
	} else {
		exception, err = dbg.Step(flags.count)
	}
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
			PrintRegisters(rvs)
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
			PrintMemory(m.Name(), base, data)
		}
	}

	// Ensure output files are written
	dbg.Close()
	os.Exit(ret)
}
