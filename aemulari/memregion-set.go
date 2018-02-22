package aemulari

import (
	"fmt"
)

// A set of MemRegion objects
type MemRegionSet map[string]MemRegion

// Create a new MemRegionSet and initialize it with MemRegions described
// by each entry of the `regionSpecs` slice.  See NewMemRegion() for the syntax
// of these entries. Each region must have a unique name.
func NewMemRegionSet(regionSpecs []string) (MemRegionSet, error) {
	ret := make(MemRegionSet)
	for _, s := range regionSpecs {
		err := ret.AddNewMemRegion(s)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}

// Create a MemRegion from the specification string `s` and add it to the set.
// An error is returned if a region with the same name already exists.
func (regions *MemRegionSet) AddNewMemRegion(s string) error {
	if m, err := NewMemRegion(s); err != nil {
		return err
	} else {
		return regions.Add(m)
	}
}

// Add an existing MemRegion to the set.
// An error is returned if a region with the same name already exists.
func (regions *MemRegionSet) Add(m MemRegion) error {
	if regions.Contains(m.name) {
		return fmt.Errorf("A region named \"%s\" has already been created.", m.name)
	}

	(*regions)[m.name] = m
	return nil
}

// Return true if `regions` contains a region named `name`
func (regions MemRegionSet) Contains(name string) bool {
	_, exists := regions[name]
	return exists
}

// Retrieve the memory region named `name`
func (regions MemRegionSet) Get(name string) (MemRegion, error) {
	m := regions[name]
	valid, _ := m.IsValid()
	if !valid {
		return m, fmt.Errorf("No such memory region: %s", name)
	}
	return m, nil
}

// Remove the MemRegion named `name`.
// Returns an error if no such MemRegion exists in the set.
func (regions *MemRegionSet) Remove(name string) error {
	if regions.Contains(name) {
		delete(*regions, name)
		return nil
	}
	return fmt.Errorf("No such memory region: %s", name)
}
