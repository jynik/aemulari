package theme

import "fmt"

type NoTheme struct{}

func (n NoTheme) ColorIfBytesDiffer(str string, value, prev byte) string {
	return str
}

func (n NoTheme) ColorIfStringsDiffer(str, prev string) string {
	return str
}

func (n NoTheme) ColorAddress(fmtspec string, addr uint64) string {
	return fmt.Sprintf(fmtspec, addr)
}

func (n NoTheme) ColorOpcode(opcode string) string {
	return opcode
}

func (n NoTheme) ColorMnemonic(mnemonic string) string {
	return mnemonic
}

func (n NoTheme) ColorOperands(operands string) string {
	return operands
}

func (n NoTheme) ArmedBreakpointSymbol() string {
	return "@"
}

func (n NoTheme) DisabledBreakpointSymbol() string {
	return "O"
}

func (n NoTheme) CurrentInstructionSymbol() string {
	return ">"
}

func (n NoTheme) ColorModifiedInstruction(line string) string {
	return line
}

func (n NoTheme) CmdPrompt() string {
	return "> "
}

func (n NoTheme) CmdSuccessSymbol() string {
	return "+"
}

func (n NoTheme) CmdFailureSymbol() string {
	return "-"
}

func (n NoTheme) ErrorMessage(e error) string {
	return e.Error()
}
