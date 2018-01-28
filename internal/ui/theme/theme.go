package theme

import (
	"fmt"
	"regexp"
	"strings"
)

// These functions return the provided string with VT100 coloring
// sequences prepended or ammended
type Theme interface {
	/**************************************************************************
	 * General purpose
	 *************************************************************************/

	// Colorize str if value != prev
	ColorIfBytesDiffer(str string, value, prev byte) string

	// Colorize str if str != prev
	ColorIfStringsDiffer(str, prev string) string

	// Colorize a code address (e.g., memory view, disassembly view)
	ColorAddress(fmtspec string, addr uint64) string

	/**************************************************************************
	 * Disassembly View
	 *************************************************************************/

	// Return colorized armed breakpoint symbol
	ArmedBreakpointSymbol() string

	// Return colorized disabled breakpoint symbol
	DisabledBreakpointSymbol() string

	// Return colorized current instruction symbol
	CurrentInstructionSymbol() string

	// Colorize (highlight) andisassembled instruction that has been modified
	ColorModifiedInstruction(line string) string

	// Color an instruction opcode
	ColorOpcode(opcode string) string

	// Color an instruction mnemonic
	ColorMnemonic(mnemonic string) string

	// Color an instruction's operands
	ColorOperands(operands string) string

	/**************************************************************************
	 * Commands View
	 *************************************************************************/

	// Return colorized successful command symbol
	CmdSuccessSymbol() string

	// Return colorized unsuccessful command symbol
	CmdFailureSymbol() string

	// Return command prompt symbol(s)
	CmdPrompt() string

	/**************************************************************************
	 * Console View
	 *************************************************************************/

	// Colorize error message
	ErrorMessage(err error) string
}

func New(name string, regNames *regexp.Regexp) (Theme, error) {
	c := strings.Trim(strings.ToLower(name), " \t\r\n\x00")
	if c == "" || c == "none" {
		return NoTheme{}, nil
	} else if c == "default" || c == "jon" || c == "jynik" {
		return CreateDefaultTheme(regNames), nil
	}

	return NoTheme{}, fmt.Errorf("\"%s\" is not a valid colorscheme.", name)
}

// Colorize the foreground (text)
func colorizeFg(color uint8, s string) string {
	return fmt.Sprintf("\x1b[38;5;%dm%s\x1b[0m", color, s)
}
