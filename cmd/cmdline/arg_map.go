package cmdline

import (
	"errors"
	"strconv"
)

// Mapping of flag names to associated arguments
type ArgMap map[string][]string

// Retrieve the string argument provided supplied with a specific flag, or a
// default value if the flag was not specified by the user.
//
// If a flag was provided multiple times, this only returns the first value.
func (m *ArgMap) GetString(name, default_val string) string {
	if val, exists := (*m)[name]; !exists || len(val) == 0 || len(val[0]) == 0 {
		return default_val
	} else {
		return val[0]
	}
}

// Get values passed by multiple usages of a flag as string slice
func (m *ArgMap) GetStrings(name string) []string {
	return (*m)[name]
}


// Retrieve a flag's values as a list of uint64 values
// Upon encountering an invalid value, that value is returned as an error
func (m *ArgMap) GetU64List(name string) ([]uint64, error) {
	ret := []uint64{}

	for _, s := range (*m)[name] {
		if val, err := strconv.ParseUint(s, 0, 64); err != nil {
			return []uint64{}, errors.New(s)
		} else {
			ret = append(ret, val)
		}
	}

	return ret, nil
}

// Retrieve a value provided to a flag  as an int64.
// If the flag was specified multiple times, only the first value is returned.
// Upon encountering an invalid value, that value is returned as an error
func (m *ArgMap) GetInt64(name string) (int64, error) {
	if s, exists := (*m)[name]; exists {
		if val, err := strconv.ParseInt(s[0], 0, 64); err != nil {
			return -1, errors.New(s[0])
		} else {
			return val, nil
		}
	} else {
		panic("Bug: GetInt64() called with invalid name: " + name)
	}
}

// Returns true if a flag was specified by a user
func (m *ArgMap) Contains(name string) bool {
	_, present := (*m)[name]
	return present
}

// Remove an entry from an argument map
func (m *ArgMap) remove(name string) {
	delete(*m, name)
}

