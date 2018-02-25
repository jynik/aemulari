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

	ret := fmt.Sprintf("%-12s [0x%08x-0x%08x] {%3s}  in:\"%s\"  out:\"%s\"",
		r.name, r.base, end, perm, r.inputFile, r.outputFile)

	return ret
}

// Returns true if the memory region has an input initialization file,
// and false otherwise
func (r MemRegion) HasInputFile() bool {
	return r.inputFile != ""
}

// Craft an error message that includes e
func (r *MemRegion) loadError(e error) ([]byte, error) {
	return []byte{}, fmt.Errorf(
		"Failed to load data for memory region \"%s\" - %s",
		r.name, e.Error())
}

// Load data that should be used to initialize a memory region
func (r *MemRegion) LoadInputData() ([]byte, error) {
	var data []byte = make([]byte, r.size)

	if r.HasInputFile() {

		f, err := os.Open(r.inputFile)
		if err != nil {
			return r.loadError(err)
		}

		// TODO Allow this to be specified in the memory specifier syntax
		offset := int64(0)

		_, err = f.Seek(offset, 0)
		if err != nil {
			return r.loadError(err)
		}

		_, err = f.Read(data)
		if err != nil {
			return r.loadError(err)
		}
		return data, err
	}

	return data, nil
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

// Returns true if MemRegion is configured correctly.
// Otherwise, false and an error is returned.
func (r MemRegion) IsValid() (bool, error) {
	var err error

	if len(r.name) == 0 {
		return false, errors.New("MemRegion name cannot be blank.")
	} else if r.name == "code" && len(r.inputFile) == 0 {
		return false, errors.New("The \"code\" memory region requires an input file.")
	}

	if (^uint64(0) - r.size) < r.base {
		err = fmt.Errorf("MemRegion \"%s\" exceeds address space limits: "+
			"base=0x%08x, size=0x%08x.", r.name, r.base, r.size)
		return false, err
	}

	if len(r.inputFile) > 0 {
		if _, err = os.Stat(r.inputFile); os.IsNotExist(err) {
			return false,
				fmt.Errorf("Cannot open input file for memory region \"%s\" - %s",
					r.name, err.Error())
		}
	}

	if r.name == "code" && !r.perms.Exec {
		return false, errors.New("The \"code\" memory region must be executable.")
	}

	return true, nil
}

// Create a memory region based upon the specification string `s`.
// The syntax of this MemRegion specification string is:
//	<name>:<addr>:<size>:[permissions]:[input file]:[output_file]
func NewMemRegion(s string) (region MemRegion, err error) {
	var fields []string

	if fields = strings.Split(s, ":"); len(fields) < 3 {
		err = errors.New("MemRegion requires at least 3 fields.")
		return
	}

	region.name = fields[0]

	if region.base, err = strconv.ParseUint(fields[1], 0, 64); err != nil {
		err = fmt.Errorf("Invalid memory region base address: %s", fields[1])
		return
	}

	if region.size, err = strconv.ParseUint(fields[2], 0, 64); err != nil || region.size == 0 {
		err = fmt.Errorf("Invalid memory region size: %s", fields[2])
		return
	}

	if len(fields) > 3 {
		if err = region.perms.Set(fields[3]); err != nil {
			return
		}
	} else {
		// Default to something very permissive and easy to work with
		region.perms.Set("rwx")
	}

	if len(fields) > 4 {
		region.inputFile = fields[4]
	}

	if len(fields) > 5 {
		region.outputFile = fields[5]
	}

	_, err = region.IsValid()
	return
}
