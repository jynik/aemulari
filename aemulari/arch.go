package aemulari

import (
	"fmt"
	"regexp"
)

// Processor type ID
type processorType struct {
	uc int // Unicorn ID for processor type
	cs int // Capstone ID for processor type
}

// Processor mode ID
type processorMode struct {
	uc int  // Unicorn ID for the initial mode
	cs uint // Capstone ID for for the initial mode
}

type archConstructor func(mode string) (Architecture, error)

type Architecture interface {
	// Return the architecture's processor type ID
	id() processorType

	// Return the processor's initial mode
	initialMode() processorMode

	// Adjust, if necessary (e.g., based upon mode or alignment), and return
	// the initial PC value.
	initialPC(pc uint64) uint64

	// Get the maximum length of an instruction
	maxInstructionSize() uint

	// Adjust current PC, if necessary.  This allows architecture-specific
	// information (e.g., current mode denoted by status register) to be
	// considered before passing the PC the emulator when (re)starting it.
	currentPC(pc uint64, regs []RegisterValue) uint64

	// Determine current data endianess.
	// The `regs` parameter should contain the current state of registers.
	// Returns BigEndian or LittleEndian
	endianness(regs []RegisterValue) Endianness

	// Create an Exception with a descriptive String() output
	//
	//	intno	Exception/Interrupt number
	//	regs	Current state of registers
	//	instr	Instruction that generated exception
	//
	exception(intno uint32, regs []RegisterValue, instr []byte) Exception

	// Parse a string and return a RegisterValue.
	// Expected form: <reg name>=<value>
	ParseRegister(s string) (RegisterValue, error)

	// Performs ParseRegister() on each entry in the provided slice, and
	// returns the associated []RegisterValue.
	ParseRegisters(s []string) ([]RegisterValue, error)

	// Look up a register definition by name or alias
	// Returns a pointer to a register definition on success,
	// or a non-nil error on failure.
	//
	register(name string) (*RegisterDef, error)

	// Return a Regular Expression for matching register names and aliases
	RegisterRegexp() *regexp.Regexp

	// Get all register definitions
	Registers() []*RegisterDef
}

// Obtain an implementation of the Architecture interface for
// the specified architecture type (`arch`). The `initialMode`
// may be used to set the initial mode of the processor, or can
// be left empty to use the default.
func NewArchitecture(arch, initialMode string) (Architecture, error) {
	if ret, found := archMap[arch]; !found {
		return nil, fmt.Errorf("Unsupported architecture: %s", arch)
	} else {
		return ret(initialMode)
	}
}

var archMap = map[string]archConstructor{
	"arm": armConstructor,
}
