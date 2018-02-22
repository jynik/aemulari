package aemulari

import (
	"fmt"
	"regexp"
	"strings"
)

// Processor type ID
type processorType struct {
	uc int // Unicorn ID for processor type
	cs int // Capstone ID for processor type
}

// Processor mode ID
type processorMode struct {
	uc int // Unicorn ID for the initial mode
	cs int // Capstone ID for for the initial mode
}

type archConstructor func(mode string) (Architecture, error)

// Architechture presents a standard interface for accessing and
// working with architectures-specific properties, such as register
// definitions.
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
	currentPC(pc uint64, regs []Register) uint64

	// Determine current data endianess.
	// The `regs` parameter should contain the current state of registers.
	// Returns BigEndian or LittleEndian
	endianness(regs []Register) Endianness

	// Create an Exception with a descriptive String() output
	//
	//	intno	Exception/Interrupt number
	//	regs	Current state of registers
	//	instr	Instruction that generated exception
	//
	exception(intno uint32, regs []Register, instr []byte) Exception

	// Parse a string and return a Register.
	// Expected form: <reg name>=<value>
	ParseRegister(s string) (Register, error)

	// Performs ParseRegister() on each entry in the provided slice, and
	// returns the associated []Register.
	ParseRegisters(s []string) ([]Register, error)

	// Return a Regular Expression that matches register names and their aliases
	RegisterRegexp() *regexp.Regexp

	// Look up a register definition by name or alias
	// Returns a pointer to a register definition on success,
	// or a non-nil error on failure.
	//
	register(name string) (*registerAttr, error)

	// Retrieve all register definiitions
	registers() []*registerAttr
}

// Obtain an implementation of the Architecture interface for
// the specified processor type, and optionally, initial mode:
//
// Examples:
//		arch, err := NewArchitecture("arm")
//		arch, err := NewArchitecture("arm:arm")
//		arch, err := NewArchitecture("arm:thumb")
func NewArchitecture(arch string) (Architecture, error) {
	var mode string

	fields := strings.Split(arch, ":")
	arch = fields[0]

	if len(fields) > 1 {
		mode = fields[1]
	} else {
		mode = ""
	}

	if ret, found := archMap[arch]; !found {
		return nil, fmt.Errorf("Unsupported architecture: %s", arch)
	} else {
		return ret(mode)
	}
}

var archMap = map[string]archConstructor{
	"arm": armConstructor,
}
