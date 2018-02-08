package aemulari

// Contains information describing a processor exception
type Exception struct {
	intno uint32 // Interrupt/Exception number
	pc    uint64 // Address at which exception occurred
	desc  string // Printable string describing the exception
}

// Tests whether Exception holds information about an exception that occurred
func (e *Exception) Occurred() bool {
	return e.desc != ""
}

// Return a string representation of an Exception
func (e *Exception) String() string {
	return e.desc
}
