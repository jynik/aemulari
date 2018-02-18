package aemulari

// Contains information describing a processor exception
type Exception struct {
	intno uint32 // Interrupt/Exception number
	pc    uint64 // Address at which exception occurred
	desc  string // Printable string describing the exception
}

// Returns true if the Exception object contains information
// about a processor exception that occurred.
func (e *Exception) Occurred() bool {
	return e.desc != ""
}

// Return a string describing a processor exception, if one occurred.
func (e *Exception) String() string {
	return e.desc
}
