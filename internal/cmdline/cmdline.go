package cmdline

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/op/go-logging"

	"../arch"
	dbg "../debugger"
)

type cmdlineArgs struct {

	// Raw arguments
	arch string // Target architecture

	verbosity string // Log verbosity

	regDefs appender // Initial register values

	// Parsed arguments
	mem  dbg.MemRegions // Memory regions to configure
	regs arch.RegisterMap
}

var args cmdlineArgs
var log = logging.MustGetLogger("")

func handleVerbosity(cfg *dbg.Config) error {

	args.verbosity = strings.ToLower(strings.Trim(args.verbosity, "\r\n\t "))
	if args.verbosity == "critical" {
		logging.SetLevel(logging.CRITICAL, "")
	} else if args.verbosity == "error" {
		logging.SetLevel(logging.ERROR, "")
	} else if args.verbosity == "warning" {
		logging.SetLevel(logging.WARNING, "")
	} else if args.verbosity == "notice" {
		logging.SetLevel(logging.NOTICE, "")
	} else if args.verbosity == "info" {
		logging.SetLevel(logging.INFO, "")
	} else if args.verbosity == "debug" || args.verbosity == "verbose" {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		return fmt.Errorf("Invalid verbosity level: %s", args.verbosity)
	}

	return nil
}

func handleArch(cfg *dbg.Config) error {
	var err error
	cfg.Arch, err = arch.New(args.arch)
	return err
}

func handleMem(cfg *dbg.Config) error {
	var err error = nil

	haveCodeRegion := args.mem.Contains("code")
	if !haveCodeRegion {
		return errors.New("A memory mapped region named \"code\" must be provided.")
	}

	cfg.Mem = args.mem

	return err
}

func handleRegDefs(cfg *dbg.Config) error {
	var err error
	cfg.RegDefs, err = cfg.Arch.ParseRegisters(args.regDefs)
	return err
}

func InitCommonFlags() {
	flag.StringVar(&args.verbosity, "v", "warning", "Logging verbosity.")
	flag.StringVar(&args.arch, "a", "arm", "Target architecture.")
	flag.Var(&args.mem, "m", "Mapped Memory regions. Specify in the form: "+dbg.MemRegionUsage())
	flag.Var(&args.regDefs, "r", "Set initial register value.")
}

func Parse() (dbg.Config, error) {
	var cfg dbg.Config

	args.mem = make(dbg.MemRegions)

	flag.Parse()

	if err := handleVerbosity(&cfg); err != nil {
		return cfg, err
	}
	if err := handleArch(&cfg); err != nil {
		return cfg, err
	}
	if err := handleMem(&cfg); err != nil {
		return cfg, err
	}
	if err := handleRegDefs(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
