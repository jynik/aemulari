package common

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/op/go-logging"

	ae "../../aemulari"
)

type cmdlineArgs struct {

	// Raw arguments
	arch string // Target architecture

	verbosity string // Log verbosity

	regDefs appender // Initial register values

	// Parsed arguments
	mem  ae.MemRegions // Memory regions to configure
	regs ae.RegisterMap
}

var args cmdlineArgs
var log = logging.MustGetLogger("")

func handleVerbosity(cfg *ae.Config) error {
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

func handleArchitecture(cfg *ae.Config) (ae.Architecture, error) {
	fields := strings.Split(args.arch, ":")

	arch := fields[0]
	mode := ""

	if len(fields) >= 2 {
		mode = fields[1]
	}

	ret, err := ae.NewArchitecture(arch, mode)
	if err != nil {
		return ret, err
	}

	cfg.RegDefs, err = ret.ParseRegisters(args.regDefs)
	return ret, err
}

func handleMem(cfg *ae.Config) error {
	var err error = nil

	haveCodeRegion := args.mem.Contains("code")
	if !haveCodeRegion {
		return errors.New("A memory mapped region named \"code\" must be provided.")
	}

	cfg.Mem = args.mem

	return err
}

func InitCommonFlags() {
	memRegionUsage := "<name>:<addr>:<size>:<permissions>[:input file[:output file]]"

	flag.StringVar(&args.verbosity, "v", "warning", "Logging verbosity.")
	flag.StringVar(&args.arch, "a", "arm", "Target architecture.")
	flag.Var(&args.mem, "m", "Mapped Memory regions. Specify in the form: "+MemRegionUsage())
	flag.Var(&args.regDefs, "r", "Set initial register value.")
}

func Parse() (ae.Architecture, ae.Config, error) {
	var cfg ae.Config
	var arch ae.Architecture
	var err error

	args.mem = make(ae.MemRegions)

	flag.Parse()

	if arch, err = handleArchitecture(&cfg); err != nil {
		return arch, cfg, err
	}

	if err := handleVerbosity(&cfg); err != nil {
		return arch, cfg, err
	}

	if err = handleMem(&cfg); err != nil {
		return arch, cfg, err
	}

	return arch, cfg, nil
}
