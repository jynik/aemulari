package cmdline

import (
	"fmt"
	"strconv"
)

// Flag value requirement, if any
type ValueReqt int
const (
	None = iota	// No value should accompany the flag
	Optional	// An argument is optional with the flag
	Required	// The flag requires a value
)

// Flag occurrence constraint
type Occurrence int

const (
	Multiple = iota		 // Flag may occur multiple times
	Once				 // Flag must occur only once
)

type Arg struct {
	Short      string     // Short form of flag
	Long       string     // Long form of flag
	ValueReqt  ValueReqt  // Is there value? Is it required or optional?
	Occurrence Occurrence // Can this argument occur multiple times or just once?
	ValidValues []string  // Only values in this list are permitted, if non-empty
	Value      string     // Value provided on the command line
}

func (a *Arg) Name() string {
	return a.Long[2:]
}

// Mapping of flag name to list of arguments associated with it
type argMap map[string][]string

func (m *argMap) getString(name, default_val string) string {
	if val, exists := (*m)[name]; exists {
		if len(val) != 1 {
			panic("Bug - using getString() for multi-value option")
		}
		return val[0]
	} else {
		return default_val
	}
}

func (m *argMap) getU64List(name string) ([]uint64, error) {
	ret := []uint64{}

	for _, s := range (*m)[name] {
		if val, err := strconv.ParseUint(s, 0, 64); err != nil {
			return []uint64{}, fmt.Errorf("%s", s)
		} else {
			ret = append(ret, val)
		}
	}

	return ret, nil
}

type SupportedArgs []*Arg

func (s *SupportedArgs) Add(a *Arg) *SupportedArgs {
	*s = append(*s, a)
	return s
}

func (s *SupportedArgs) lookup(flag string) *Arg {
	for _, elt := range *s {
		if flag == elt.Short || flag == elt.Long {
			return elt
		}
	}

	return nil
}

// Parse command line arguments and separate them into their
// associated flags.
func (l *SupportedArgs) parse(args []string) (argMap, error) {
	var ret argMap = make(argMap)
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

		// Enforce any occurrance restrictions
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

