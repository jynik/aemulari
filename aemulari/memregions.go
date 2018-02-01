package aemulari

import (
	"bytes"
	"fmt"
)

type MemRegions map[string]MemRegion

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
	return !m.IsZero()
}

// Create a MemRegion from specification `s` and add it to `regions`
//
// This function implements a portion of the Value interface.
//
// TODO: Make a Set() wrapper with a more appropriately named fn?
func (regions *MemRegions) Set(s string) error {
	var m MemRegion
	if err := m.Set(s); err != nil {
		return err
	}

	return regions.Add(m)
}

func (regions *MemRegions) Add(m MemRegion) error {
	if regions.Contains(m.name) {
		return fmt.Errorf("A region named \"%s\" has already been created.", m.name)
	}

	(*regions)[m.name] = m
	return nil
}

func (regions MemRegions) Get(s string) (MemRegion, error) {
	m := regions[s]
	if m.IsZero() {
		return m, fmt.Errorf("Memory region \"%s\" does not exist.", s)
	}
	return m, nil
}

func (regions *MemRegions) Remove(name string) {
	delete(*regions, name)
}
