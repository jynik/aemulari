package aemulari

import (
	"fmt"
)

// Register flag bits
type Flag struct {
	name string // Flag (short) name
	desc string // Flag description
	mask uint64 // Flag bit mask
	lsb  uint   // Least significant bit position
	fmt  string // Format string for flag representation
}

// Definition of register
type RegisterDef struct {
	name  string // Primary name
	mask  uint64 // Read/Write mask
	fmt   string // Format to use for string representation
	uc    int    // Unicorn register identifier
	pc    bool   // This register is the program couner
	Flags []Flag // Named flag bits for this register
}

type RegisterValue struct {
	Reg   *RegisterDef
	Value uint64
}

// Get the flag value, provided a register value
func (f *Flag) Get(reg uint64) uint64 {
	return (reg & f.mask) >> f.lsb
}

// Get a string representation of a flag value, including its name
func (f *Flag) GetString(reg uint64) string {
	return fmt.Sprintf(f.fmt, f.Get(reg))
}

// Similar to GetString, but with the flag short name prefixed
func (f *Flag) GetNamedString(reg uint64) string {
	return fmt.Sprintf("%-4s"+f.fmt, f.name, f.Get(reg))
}

func (f *Flag) Name() string {
	return f.name
}

/* Set the specified flag value, provided the current register state.
 * Return the updated register value. */
func (f *Flag) Set(reg uint64, flag uint64) uint64 {
	return reg | ((flag << f.lsb) & f.mask)
}

func (r *RegisterValue) ValueString() string {
	return fmt.Sprintf(r.Reg.fmt, r.Value)
}

func (r *RegisterValue) Name() string {
	return r.Reg.Name()
}

func (r *RegisterValue) String() string {
	return fmt.Sprintf("%-6s"+r.Reg.fmt, r.Reg.name, r.Value)
}

func (r *RegisterDef) Name() string {
	return r.name
}

// Is this register the program counter?
func (r *RegisterDef) IsProgramCounter() bool {
	return r.pc
}
