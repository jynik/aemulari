package aemulari

import (
	"fmt"
	"sort"
	"strings"
)

// A set of MemRegion objects, identified by case-sensitive name.
type MemRegionSet struct {
	entries map[string]MemRegion
}

// Create and initialize an empty MemRegionSet
//
// This function or NewMemRegionSet must be used before attempting to
// access a MemRegionSet.
func EmptyMemRegionSet() MemRegionSet {
	return MemRegionSet{ entries: map[string]MemRegion{} }
}

// Create a new MemRegionSet and initialize it with MemRegions described
// by each entry of the `regionSpecs` slice.  See NewMemRegion() for the syntax
// of these entries. Each region must have a unique name.
//
// This function or NewEmptyMemRegionSet must be used before attempting to
// access a MemRegionSet.
func NewMemRegionSet(regionSpecs []string) (MemRegionSet, error) {
	ret := MemRegionSet{}
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
	regions.entries[m.name] = m
	return nil
}

// Return true if `regions` contains a region named `name`
func (regions MemRegionSet) Contains(name string) bool {
	_, exists := regions.entries[name]
	return exists
}

// Retrieve the memory region named `name`
func (regions MemRegionSet) Get(name string) (MemRegion, error) {
	m := regions.entries[name]
	valid, _ := m.IsValid()
	if !valid {
		return m, fmt.Errorf("No such memory region: %s", name)
	}
	return m, nil
}

// Retrieve a slice all entries in a MemRegion, sorted by base address
func (regions *MemRegionSet) Entries() []MemRegion {
	ret := []MemRegion{}
	for _, r := range regions.entries {
		ret = append(ret, r)
	}
	sort.Slice(ret, func(i, j int) bool {
		if ret[i].base == ret[j].base {
			return strings.Compare(ret[i].name, ret[j].name) < 0
		} else {
			return ret[i].base < ret[j].base
		}
	})
	return ret
}

// Remove the MemRegion named `name`.
// Returns an error if no such MemRegion exists in the set.
func (regions *MemRegionSet) Remove(name string) error {
	if regions.Contains(name) {
		delete(regions.entries, name)
		return nil
	}
	return fmt.Errorf("No such memory region: %s", name)
}
