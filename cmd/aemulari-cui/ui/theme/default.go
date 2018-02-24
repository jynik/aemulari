package theme

import (
	"fmt"
	"regexp"
)

type DefaultTheme struct {
	immediate *regexp.Regexp
	regNames  *regexp.Regexp
}

const differsColor = immediateColor
const errorColor = 9
const addrColor = 30
const opcodeColor = 242
const mnemonicColor = 255
const registerColor = 80
const immediateColor = 48
const cmdSuccessColor = 70
const cmdErrorColor = errorColor
const breakpointColor = 124
const currentInstrColor = 48

func CreateDefaultTheme(regNames *regexp.Regexp) (theme DefaultTheme) {
	theme.regNames = regNames
	theme.immediate = regexp.MustCompile("(#-?(0x)?[[:xdigit:]]+)")
	return theme
}

func (d DefaultTheme) ColorIfBytesDiffer(str string, value, prev byte) string {
	if value != prev {
		return colorizeFg(differsColor, str)
	}
	return str
}

func (d DefaultTheme) ColorIfStringsDiffer(str, prev string) string {
	if str != prev {
		str = colorizeFg(differsColor, str)
	}
	return str
}

func (d DefaultTheme) ColorAddress(fmtspec string, addr uint64) string {
	return colorizeFg(addrColor, fmt.Sprintf(fmtspec, addr))
}

func (d DefaultTheme) ColorOpcode(opcode string) string {
	return colorizeFg(opcodeColor, opcode)
}

func (d DefaultTheme) ColorMnemonic(mnemonic string) string {
	return colorizeFg(mnemonicColor, mnemonic)
}

func (d DefaultTheme) ColorOperands(operands string) string {
	regsRepl := []byte("$1" + colorizeFg(registerColor, "$2"))
	immRepl := []byte(colorizeFg(immediateColor, "$1"))

	// Color all register names
	coloredOperands := d.regNames.ReplaceAll([]byte(operands), regsRepl)

	// Color intermediate values
	return string(d.immediate.ReplaceAll(coloredOperands, immRepl))
}

func (d DefaultTheme) ArmedBreakpointSymbol() string {
	return colorizeFg(breakpointColor, "B")
}

func (d DefaultTheme) DisabledBreakpointSymbol() string {
	return colorizeFg(breakpointColor, "b")
}

func (d DefaultTheme) CurrentInstructionSymbol() string {
	return colorizeFg(currentInstrColor, ">")
}

func (d DefaultTheme) ColorModifiedInstruction(line string) string {
	return colorizeFg(differsColor, line)
}

func (d DefaultTheme) CmdPrompt() string {
	return "> "
}

func (d DefaultTheme) CmdSuccessSymbol() string {
	return colorizeFg(cmdSuccessColor, "+")
}

func (d DefaultTheme) CmdFailureSymbol() string {
	return colorizeFg(errorColor, "x")
}

func (d DefaultTheme) ErrorMessage(e error) string {
	return colorizeFg(errorColor, e.Error())
}
