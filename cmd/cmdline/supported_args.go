package cmdline

var Arg_arch *Arg = &Arg{
	Short:      "-a",
	Long:       "--arch",
	Occurrence: Once,
	ValueReqt:       Required,
}

var Arg_reg *Arg = &Arg{
	Short:      "-r",
	Long:       "--reg",
	Occurrence: Multiple,
	ValueReqt:       Required,
}

var Arg_mem *Arg = &Arg{
	Short:      "-m",
	Long:       "--mem",
	Occurrence: Multiple,
	ValueReqt:		Required,
}

var Arg_breakpoint *Arg = &Arg{
	Short: "-b",
	Long: "--break",
	Occurrence: Multiple,
	ValueReqt: Required,
}

var Arg_printRegs *Arg = &Arg{
	Short: "-R",
	Long: "--print-regs",
	Occurrence: Once,
	ValueReqt: Optional,
	ValidValues: []string{ "pretty", "list" },
}

var Arg_hexdump *Arg = &Arg{
	Short: "-d",
	Long: "--hexdump",
	Occurrence: Once,
	ValueReqt: Required,
}
