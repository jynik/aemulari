package common

import (
	"flag"
	"strings"

	ae "../../aemulari"
)

type cmdlineArgs struct {

	// Raw arguments
	arch string // Target architecture

	regDefaults appender // Initial register values

	// Parsed arguments
	mem ae.MemRegions // Memory regions to configure
}

var args cmdlineArgs

func handleArchitecture(cfg *ae.DebuggerConfig) (ae.Architecture, error) {
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

	cfg.Regs, err = ret.ParseRegisters(args.regDefaults)
	return ret, err
}

func InitCommonFlags() {
	memRegionUsage := "<name>:<addr>:<size>:[permissions]:[input file]:[output file]"

	flag.StringVar(&args.arch, "a", "arm", "Target architecture.")
	flag.Var(&args.mem, "m", "Mapped Memory regions. Specify in the form: "+memRegionUsage)
	flag.Var(&args.regDefaults, "r", "Set initial register value.")
}

func Parse() (ae.Architecture, ae.DebuggerConfig, error) {
	var cfg ae.DebuggerConfig
	var arch ae.Architecture
	var err error

	args.mem = make(ae.MemRegions)

	flag.Parse()

	if arch, err = handleArchitecture(&cfg); err != nil {
		return arch, cfg, err
	}

	cfg.Mem = args.mem
	return arch, cfg, nil
}
