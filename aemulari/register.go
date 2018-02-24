package aemulari

import (
	"fmt"
	"strings"
)

// Named flag bit in a register
type registerFlag struct {
	name string // Flag (short) name
	desc string // Flag description
	mask uint64 // Flag bit mask
	lsb  uint   // Least significant bit position
	fmt  string // Format string for flag representation
}

// Definition of register and its attributes
type registerAttr struct {
	name  string         // Primary name
	mask  uint64         // Read/Write mask
	fmt   string         // Format to use for string representation
	uc    int            // Unicorn register identifier
	pc    bool           // This register is the program counter
	flags []registerFlag // Named flag bits for this register
}

// Information about a processor's register, including its name, current value, and flags.
type Register struct {
	attr  *registerAttr
	Value uint64
}

// Return the name of a Register.
func (r *Register) Name() string {
	return r.attr.name
}

// Return a string that includes a Register's name and current value.
func (r *Register) String() string {
	return fmt.Sprintf("%-6s"+r.attr.fmt, r.attr.name, r.Value)
}

// Lookup a register's flag by name
func (r *Register) getFlagByName(name string) (*registerFlag, error) {
	targetFlag := strings.TrimSpace(strings.ToLower(name))

	for _, f := range r.attr.flags {
		if f.name == targetFlag {
			return &f, nil
		}
	}

	return nil, fmt.Errorf("No such flag in register %s: %s", r.attr.name, name)
}

// Extract the flag bit values from a register value
func (r *Register) getFlagValue(f *registerFlag) uint64 {
	return (r.Value & f.mask) >> f.lsb
}

// Extract the flag bit values from a register value,
// specifying the flag bits by name
func (r *Register) getFlagValueByName(name string) (uint64, error) {
	if flag, err := r.getFlagByName(name); err != nil {
		return ^uint64(0), err
	} else {
		return r.getFlagValue(flag), nil
	}
}

// Update a register value by setting flag bits to the specified state
func (r *Register) setFlagValue(f *registerFlag, value uint64) {
	r.Value &= ^f.mask
	r.Value |= ((value << f.lsb) & f.mask)
}

// Update a register value by setting flag bits to the provided state,
// with the flag bit(s) specified by name
func (r *Register) setFlagValueByName(name string, value uint64) error {
	if flag, err := r.getFlagByName(name); err != nil {
		return err
	} else {
		r.setFlagValue(flag, value)
		return nil
	}
}

// Return a list of strings that include a Register's flag bits and their associated values.
func (r *Register) FlagStrings() []string {
	var ret []string
	for _, flag := range r.attr.flags {
		ret = append(ret, fmt.Sprintf("%-4s"+flag.fmt, flag.name, r.getFlagValue(&flag)))
	}
	return ret
}
