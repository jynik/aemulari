package aemulari

import "fmt"

// A Breakpoint may be used to halt execution when it reaches a specific address.
type Breakpoint struct {
	ID      int					// Numeric identifier for the Breakpoint
	Address uint64				// Address where Breakpoint is placed
	count   uint				// Number of times the breakpoint's been hit
	state   breakpointState		// Current state of the breakpoint
}

// A list of Breakpoint objects
type BreakpointList []Breakpoint

type breakpointState int

const (
	breakpointInvalid = iota
	breakpointInactive
	breakpointArmed
	breakpointTriggered
	breakpointMax
)

func newBreakpoint(id int, addr uint64) Breakpoint {
	var b Breakpoint

	b.Address = addr
	b.ID = id
	b.Reset()

	return b
}

// Reset the Breakpoint hit count and re-enable the Breakpoint
func (b *Breakpoint) Reset() {
	b.count = 0
	b.Enable()
}

// Disable the breakpoint.
// Execution will no longer halt at the associated address.
func (b *Breakpoint) Disable() {
	b.state = breakpointInactive
}

// Enable the breakpoint.
// Execution will halt at the associated address.
func (b *Breakpoint) Enable() {
	b.state = breakpointArmed
}

// Return true if the Breakpoint is enabled, and false otherwise.
func (b *Breakpoint) Enabled() bool {
	return b.state != breakpointInactive
}

// Register a potential breakpoint hit
func (b *Breakpoint) hit(addr uint64) bool {
	if addr != b.Address {
		return false
	}

	if b.state <= breakpointInvalid || b.state >= breakpointMax {
		log.Warning(fmt.Sprintf("BP is in valid state (%d)", b.state))
		return false
	}

	b.count++

	if b.state == breakpointArmed {
		b.state = breakpointTriggered
		return true
	}

	return false
}

// Return a string representation of the breakpoint
func (b Breakpoint) String() string {
	// FIXME address format string should be arch-dependent
	return fmt.Sprintf("Breakpoint %2d: 0x%08x, Hit count = %d", b.ID, b.Address, b.count)
}

// Returns true if any of the breakpoints in the provided list are enabled
func (bpl BreakpointList) Enabled() bool {
	for _, b := range bpl {
		if b.Enabled() {
			return true
		}
	}
	return false
}
