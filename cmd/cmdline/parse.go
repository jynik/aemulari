package cmdline

import (
	"fmt"
	"os"
	"path/filepath"

	ae "../../aemulari"
)


// Scan arguments for help requests
func helpRequested() bool {
	for _, arg := range(os.Args) {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

// Parse command line arguments (argv) based upon list of supported arguments.
// Configures and returns an Architecture and Debugger on success.
// Prints errors to stderr and exits the program on failure.
func Parse(supported SupportedArgs, usage string) (*ae.Architecture, *ae.Debugger) {
	var dbgCfg ae.DebuggerConfig

	if helpRequested() || len(os.Args) == 1 {
		fmt.Printf(usage, filepath.Base(os.Args[0]))
		os.Exit(0)
	}

	args, err := supported.parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Aggregate and convert breakpoints to uint64 addresses
	breakpoints, err := args.getU64List("break")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid breakpoint address encountered: %s\n", err)
		os.Exit(1)
	}

	// Aggregate memory regions
	dbgCfg.Mem, err = ae.NewMemRegionSet(args["mem"])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Determine which architecture we're emulating
	arch, err := ae.NewArchitecture(args.getString("arch", "arm"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Parse user-provided initial register values for the configured architecture
	dbgCfg.Regs, err = arch.ParseRegisters(args["regs"])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Create the debugger and set any initial breakpoints
	dbg, err := ae.NewDebugger(arch, dbgCfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, b := range breakpoints {
		dbg.SetBreakpoint(b)
	}

	return &arch, dbg
}
