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
	arch     string // Target architecture
	codeFile string // Filename of input code

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

	haveCodeFile := (args.codeFile != "")
	haveCodeRegion := args.mem.Contains("code")

	if haveCodeFile && haveCodeRegion {
		return errors.New("Both a code file and code memory region were specified.")
	} else if haveCodeFile {
		err = args.mem.Set(fmt.Sprintf("code:0x%08x:0x%08x:rwx:%s",
			cfg.Arch.Defaults().CodeBase,
			cfg.Arch.Defaults().CodeSize,
			args.codeFile))
	} else {
		return errors.New("Either a code file (-f) or code memory region (-m) must be provided.")
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
	flag.StringVar(&args.codeFile, "f", "", "Code file to execute.")
	flag.Var(&args.mem, "m", "Mapped Memory regions.")
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
