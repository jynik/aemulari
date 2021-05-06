package cmdline

var Flag_arch *Flag = &Flag{
	Short:      "-a",
	Long:       "--arch",
	Occurrence: Once,
	ValueReqt:  Required,
}

var Flag_reg *Flag = &Flag{
	Short:      "-r",
	Long:       "--reg",
	Occurrence: Multiple,
	ValueReqt:  Required,
}

var Flag_mem *Flag = &Flag{
	Short:      "-m",
	Long:       "--mem",
	Occurrence: Multiple,
	ValueReqt:  Required,
}

var Flag_instrcount *Flag = &Flag{
	Short:      "-n",
	Long:       "--instr-count",
	Occurrence: Once,
	ValueReqt:  Required,
}

var Flag_breakpoint *Flag = &Flag{
	Short:      "-b",
	Long:       "--break",
	Occurrence: Multiple,
	ValueReqt:  Required,
}

var Flag_printRegs *Flag = &Flag{
	Short:       "-R",
	Long:        "--print-regs",
	Occurrence:  Once,
	ValueReqt:   Optional,
	ValidValues: []string{"pretty", "list"},
}

var Flag_hexdump *Flag = &Flag{
	Short:      "-d",
	Long:       "--hexdump",
	Occurrence: Multiple,
	ValueReqt:  Required,
}

var Flag_ghidra *Flag = &Flag{
	Short:		"-G",
	Long:		"--ghidra",
	Occurrence: Once,
	ValueReqt:	None,
}
