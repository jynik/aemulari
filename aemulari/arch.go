package aemulari

import (
	"fmt"
	"regexp"
)

type Endianness int

const (
	BigEndian = iota
	LittleEndian
)

// Processor type ID
type processorType struct {
	uc int	// Unicorn ID for processor type
	cs int  // Capstone ID for processor type
}

// Processor mode ID
type processorMode struct {
	uc int	// Unicorn ID for the initial mode
	cs uint // Capstone ID for for the initial mode
}

// Processor exception
type exception struct {
	intno uint32 // Interrupt/Exception number
	pc    uint64 // Address at which exception occurred
	desc  string // Printable string describing the exception
}

// Tests if exception metadata indicates that an exception has been recorded
func (e *exceptionInfo) Occurred() bool {
	return e.desc != ""
}

type archConstructor func(mode string) (Architecture, error)

type archBase struct {
	processor processorType
	mode	  processorMode
	maxInstrLen uint
	RegisterMap
}

type Architecture interface {

	// Adjust, if necessary (e.g., based upon mode or alignment), and return
	// the initial PC value.
	initialPC(pc uint64) uint64

	// Adjust current PC, if necessary.  This allows architecture-specific
	// information (e.g., current mode denoted by status register) to be
	// considered before passing the PC the emulator when (re)starting it.
	currentPC(pc uint64, regs []RegisterValue) uint64

	// Determine current data endianess.
	// The `regs` parameter should contain the current state of registers.
	// Returns BigEndian or LittleEndian
	endianness(regs []RegisterValue) Endianness

	// Create an Exception with a descriptive String() output
	// FIXME this looks outdated
	//
	//	intno	Exception/Interrupt number
	//	regs	Current state of registers
	//	instr	Instruction that generated exception
	//
	exception(intno uint32, regs []RegisterValue, instr []byte) exceptionInfo

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
	// FIXME - not needed?
	//register(name string) (*RegisterDef, error)

	/*
	 * Return a Regular Expression for matching register names and aliases
	 */
	RegisterRegexp() *regexp.Regexp

	/* Get all register definitions */
	Registers() []*RegisterDef
}

// Obtain an implementation of the Architecture interface for
// the specified architecture type (`arch`). The `initialMode`
// may be used to set the initial mode of the processor, or can
// be left empty to use the default.
func NewArchitecture(arch, initialMode string) (Architechture, error) {
	if ret, found := archMap[arch]; !found {
		return nil, fmt.Errorf("Unsupported architecture: %s", arch)
	} else {
		return ret(initialMode)
	}
}

var archMap = map[string]archConstructor{
	"arm": armConstructor,
}
