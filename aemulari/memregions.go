package aemulari

import (
	"bytes"
	"fmt"
)

// A set MemRegion structures that may be accessed by name
type MemRegions map[string]MemRegion

// Returns a newline-delimited list of memory region strings
func (regions MemRegions) String() string {
	var buf bytes.Buffer

	for _, r := range regions {
		buf.WriteString(r.String() + "\n")
	}

	return buf.String()
}

// Return true if `regions` contains a region named `r`
func (regions MemRegions) Contains(name string) bool {
	m := regions[name]
	return m.name == ""
}

// Create a MemRegion from specification `s` and add it to `regions`
//
// This function implements a portion of the Value interface.
//
// TODO: Make a Set() wrapper with a more appropriately named fn?
func (regions *MemRegions) Set(s string) error {
	if m, err := NewMemRegion(s); err != nil {
		return err
	} else {
		return regions.Add(m)
	}
}

// Add a memory region to the mapping
func (regions *MemRegions) Add(m MemRegion) error {
	if regions.Contains(m.name) {
		return fmt.Errorf("A region named \"%s\" has already been created.", m.name)
	}

	(*regions)[m.name] = m
	return nil
}

// Retrieve the memory region named `name`
func (regions MemRegions) Get(name string) (MemRegion, error) {
	m := regions[name]
	valid, _ := m.IsValid()
	if !valid {
		return m, fmt.Errorf("Memory region \"%s\" does not exist.", name)
	}
	return m, nil
}

// Remove the memory region named `name`
func (regions *MemRegions) Remove(name string) {
	delete(*regions, name)
}
