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

type Defaults struct {
	Mode     Mode
	CodeBase uint64
	CodeSize uint64
}

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

type archConstructor func() Arch

type Arch interface {
	/* Return per-Architecture defaults */
	Defaults() Defaults

	/* Return the Architecture type */
	Type() Type

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

func New(arch string) (Arch, error) {
	arch = strings.Trim(arch, "\r\n\t ")
	arch = strings.ToLower(arch)

	if newArch, found := archMap[arch]; !found {
		return nil, fmt.Errorf("Unsupported architecture: %s", arch)
	} else {
		return newArch(), nil
	}
}

var archMap = map[string]archConstructor{
	"arm": armConstructor,
}
