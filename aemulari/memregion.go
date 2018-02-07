package aemulari

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Configuration of a memory region
type MemRegion struct {
	name string // Name of the region

	base uint64 // Base address of the region
	size uint64 // Size of the region, in bytes

	perms Permissions // Access permission of the region

	// Path to a file containing data that will be used to initially populate
	// the region. An empty string may be used to specify that the region
	// should just be zeroized.
	inputFile string

	// Path to a file to write the contents of a memory region when it is
	// unmapped. An empty string may be used to denote that the data shouldn't
	// be written.
	outputFile string
}

// Returns the name of a region
func (r MemRegion) Name() string {
	return r.name
}

// Returns the base address of a region and its size (bytes)
func (r MemRegion) Region() (uint64, uint64) {
	return r.base, r.size
}

// Returns the address of the byte AFTER the end of the MemRegion
func (r MemRegion) End() uint64 {
	return r.base + r.size
}

// Returns the permissions of the MemRegion
func (r MemRegion) Permissions() Permissions {
	return r.perms
}

// Returns a string representation of the MemRegion
func (r MemRegion) String() string {
	perm := r.perms.String()
	end := r.base + r.size - 1

	ret := fmt.Sprintf("%s [0x%08x-0x%08x] {%s} in:\"%s\" out:\"%s\"",
		r.name, r.base, end, perm, r.inputFile, r.outputFile)

	return ret
}

// Returns true if the memory region has an input initialization file,
// and false otherwise
func (r MemRegion) HasInputFile() bool {
	return r.inputFile != ""
}

// Load data that should be used to initialize a memory region
func (r MemRegion) LoadInputData() ([]byte, error) {
	if r.HasInputFile() {
		return ioutil.ReadFile(r.inputFile)
	}
	return []byte{}, nil
}

// Returns true if the memory region has an output file, and false otherwise
func (r MemRegion) HasOutputFile() bool {
	return r.outputFile != ""
}

// Write the provided data to the MemRegion's output file. If not configured
// for an output file, this function does nothing and returns nil.
func (r MemRegion) WriteFile(data []byte) error {
	if r.HasOutputFile() {
		return ioutil.WriteFile(r.outputFile, data, 0644)
	}
	return nil
}

// Create a memory region as specified by the string `s`, of the form:
//	<name>:<addr>:<size>:[permissions[:input file[:output file]]]
func NewMemRegion(s string) (region MemRegion, err error) {
	var fields []string

	if fields = strings.Split(s, ":"); len(fields) < 3 {
		err = errors.New("MemRegion requires at least 3 fields.")
		return
	}

	if region.name = fields[0]; len(region.name) == 0 {
		err = errors.New("Memory region name cannot be blank.")
		return
	}

	if region.base, err = strconv.ParseUint(fields[1], 0, 64); err != nil {
		err = fmt.Errorf("Invalid memory region base address: %s", fields[1])
		return
	}

	if region.size, err = strconv.ParseUint(fields[2], 0, 64); err != nil || region.size == 0 {
		err = fmt.Errorf("Invalid memory region size: %s", fields[2])
		return
	}

	if (^uint64(0) - region.size) < region.base {
		err = fmt.Errorf("0x%08x:0x%08x exceeds address space limits.", region.base, region.size)
		return
	}

	if len(fields) > 3 {
		if err = region.perms.Set(fields[3]); err != nil {
			return
		}
	}

	if len(fields) > 4 {
		region.inputFile = fields[4]
		if _, err = os.Stat(region.inputFile); os.IsNotExist(err) {
			return
		}
	}

	if len(fields) > 5 {
		region.outputFile = fields[5]
	}

	return region, nil
}
