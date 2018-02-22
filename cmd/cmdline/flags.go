package cmdline

import (
	"fmt"
)

// Flag value requirement, if any
type ValueReqt int
const (
	None = iota	// No value should accompany the flag
	Optional	// An argument is optional with the flag
	Required	// The flag requires a value
)

// Flag occurrence constraints
type Occurrence int

const (
	Multiple = iota		 // Flag may occur multiple times
	Once				 // Flag must occur only once
)

// Command line flag attributes and constraints
type Flag struct {
	Short      string     // Short form of flag
	Long       string     // Long form of flag
	ValueReqt  ValueReqt  // Is there value? Is it required or optional?
	Occurrence Occurrence // Can this argument occur multiple times or just once?
	ValidValues []string  // Only values in this list are permitted, if non-empty
}

// Return a command line flag's name, which is derived from it's long form
func (f *Flag) Name() string {
	return f.Long[2:]
}

// This type is used to maintain a list of flags supported by the program
type SupportedFlags []*Flag

// Add a Flag to the list to the list of supported options
func (s *SupportedFlags) Add(f *Flag) *SupportedFlags{
	*s = append(*s, f)
	return s
}

// Retrieve a Flag by its short or long form
func (s *SupportedFlags) lookup(flag string) *Flag {
	for _, elt := range *s {
		if flag == elt.Short || flag == elt.Long {
			return elt
		}
	}

	return nil
}

// Parse command line arguments and separate them into their
// associated flags, reporting any misuse as an error.
//
// Returns a ArgMap on success.  This maps flag names to the
// associated argument values.
func (l *SupportedFlags) parse(args []string) (ArgMap, error) {
	var ret ArgMap = make(ArgMap)
	var err error
	var arg string

	for i := 0; i < len(args); i++ {
		flag := l.lookup(args[i])
		if flag == nil {
			return ret, fmt.Errorf("Invalid option provided: %s", args[i])
		}

		flagName := fmt.Sprintf("%s/%s", flag.Short, flag.Long)

		if i+1 < len(args) {
			arg = args[i+1]
		} else {
			arg = ""
		}

		// Test our argument value requirements
		switch flag.ValueReqt {
		case None:
			if len(arg) != 0 && l.lookup(arg) != nil {
				return ret, fmt.Errorf("The %s option does not take an argument.", flagName)
			}
		case Required:
			if len(arg) == 0 || l.lookup(arg) != nil {
				return ret, fmt.Errorf("The %s option requires an argument.", flagName)
			}
			i += 1
		case Optional:
			if len(arg) != 0 && l.lookup(arg) != nil {
				// Zap this, it's not really an arg
				arg = ""
			} else {
				i += 1
			}
		}

		// Enforce value whitelist
		if len(arg) != 0 && len(flag.ValidValues) != 0 {
			valid := false
			for _, val := range flag.ValidValues {
				if arg == val {
					valid = true
					break
				}
			}
			if !valid {
				return ret, fmt.Errorf("Invalid value for %s option: %s", flagName, arg)
			}
		}

		// Enforce any occurrence restrictions
		values, alreadyPresent := ret[flag.Name()]
		if alreadyPresent {
			if flag.Occurrence == Once {
				return ret, fmt.Errorf("The %s option may not be specified more than once.", flagName)
			}

			values = append(values, arg)
			ret[flag.Name()] = values
		} else {
			ret[flag.Name()] = []string{ arg }
		}
	}

	return ret, err
}
