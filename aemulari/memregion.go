package aemulari

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type MemRegion struct {
	name string

	base uint64
	size uint64

	perms Permissions

	inputFile string

	outputFile string
}

func (r MemRegion) Name() string {
	return r.name
}

func (r MemRegion) Region() (uint64, uint64) {
	return r.base, r.size
}

func (r MemRegion) End() uint64 {
	return r.base + r.size
}

func (r MemRegion) Permissions() Permissions {
	return r.perms
}

func (r MemRegion) String() string {
	perm := r.perms.String()
	end := r.base + r.size - 1

	ret := fmt.Sprintf("%s [0x%08x-0x%08x] {%s} in:\"%s\" out:\"%s\"",
		r.name, r.base, end, perm, r.inputFile, r.outputFile)

	return ret
}

func MemRegionUsage() string {
	return "<name>:<addr>:<size>[:input file[:output file]]"
}

func (r MemRegion) LoadInputData() ([]byte, error) {
	return ioutil.ReadFile(r.inputFile)
}

func (r MemRegion) HasOutputFile() bool {
	return r.outputFile != ""
}

func (r MemRegion) WriteFile(data []byte) error {
	return ioutil.WriteFile(r.outputFile, data, 0644)
}

func (r *MemRegion) Set(s string) error {
	var err error
	var fields []string

	if fields = strings.Split(s, ":"); len(fields) < 3 {
		return errors.New("MemRegion requires at least 3 fields.")
	}

	if r.name = fields[0]; len(r.name) == 0 {
		return errors.New("Memory region name cannot be blank.")
	}

	if r.base, err = strconv.ParseUint(fields[1], 0, 64); err != nil {
		return fmt.Errorf("Invalid memory region base address: %s", fields[1])
	}

	if r.size, err = strconv.ParseUint(fields[2], 0, 64); err != nil || r.size == 0 {
		return fmt.Errorf("Invalid memory region size: %s", fields[2])
	}

	if (^uint64(0) - r.size) < r.base {
		return fmt.Errorf("0x%08x:0x%08x exceeds address space limits.", r.base, r.size)
	}

	if len(fields) > 3 {
		if err := r.perms.Set(fields[3]); err != nil {
			return err
		}
	}

	if len(fields) > 4 {
		r.inputFile = fields[4]
	}

	if len(fields) > 5 {
		r.outputFile = fields[5]
	}

	return nil
}

func (r MemRegion) IsZero() bool {
	return r.size == 0
}
