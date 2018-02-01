package arch

import (
	"fmt"
	"regexp"
	"strings"
)

type Endianness int

const (
	BigEndian = iota
	LittleEndian
)

type Type struct {
	Uc int
	Cs int
}

type Mode struct {
	Uc int
	Cs uint
}

type Exception struct {
	intno uint32 // Interrupt/Exception number
	pc    uint64 // Address at which exception ocurred
	desc  string // Printable string describing the exception
}

func (e *Exception) String() string {
	return e.desc
}

func (e *Exception) Occurred() bool {
	return e.desc != ""
}

// args form a single string in the form: arch[:mode]
type archConstructor func(args string) (Arch, error)

type Arch interface {
	// Return the processor's architecture type
	Type() Type

	// Return the initial processor mode
	InitialMode() Mode

	// Adjust, if neccessary (e.g., based upon mode or alignment), and return
	// the initial PC value.  For example, we'll need to set the LSB for ARM
	// THUMB mode. Report errors if a provided PC value is invalid.
	InitialPC(pc uint64) (uint64, error)

	// Adjust current PC, if neccessary.  This allows architecture-specifics
	// (e.g., mode denoted in status reguster) to be considered before
	// (re)starting the emulator. Returns error on invalid value.
	CurrentPC(pc uint64, regs []RegisterValue) (uint64, error)

	/* Return the max size of an instruction, in bytes */
	MaxInstrLen() uint

	/* Determine current data endianess.
	 * The `regs` parameter should contain the current state of registers.
	 * Returns BigEndian or LittleEndian
	 */
	Endianness(regs []RegisterValue) Endianness

	/* Create an Exception with a descriptive String() output
	 *
	 * Parameters:
	 *	intno	Exception/Interrupt number
	 *	regs	Current state of registers
	 *	instr	Instruction that generated exception
	 *
	 * Returns a string describing the exception
	 */
	Exception(intno uint32, regs []RegisterValue, instr []byte) Exception

	/* Parse a string and return a RegisterValue.
	 * Expected form: <reg name>=<value>
	 */
	ParseRegister(s string) (RegisterValue, error)

	/* Performs ParseRegister() on each entry in the provided slice, and
	 * returns the associated []RegisterValue. */
	ParseRegisters(s []string) ([]RegisterValue, error)

	/* Look up a register definition by name or alias
	 * Returns a pointer to a register definition on success,
	 * or a non-nil error on failure.
	 */
	Register(name string) (*RegisterDef, error)

	/*
	 * Return a Regular Expression for matching register names and aliases
	 */
	RegisterRegexp() *regexp.Regexp

	/* Get all register definitions */
	Registers() []*RegisterDef
}

// Create a new Arch instance matching the architecture and mode
// specified by "arch[:mode]"
func New(args string) (Arch, error) {
	var arch, mode string

	args = strings.Trim(args, "\r\n\t ")
	args = strings.ToLower(args)

	argv := strings.Split(args, ":")
	arch = argv[0]

	if len(argv) > 1 {
		mode = argv[1]
	} else {
		mode = ""
	}

	if newArch, found := archMap[arch]; !found {
		return nil, fmt.Errorf("Unsupported architecture: %s", arch)
	} else {
		return newArch(mode)
	}
}

var archMap = map[string]archConstructor{
	"arm": armConstructor,
}
