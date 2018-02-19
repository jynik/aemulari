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
type ArgMap map[string][]string

func (m *ArgMap) Get(name string) []string {
	if val, exists := (*m)[name]; exists {
		return val
	} else {
		return []string{ "" }
	}
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
func (l *SupportedArgs) parse(args []string) (ArgMap, error) {
	var ret ArgMap = make(ArgMap)
	var err error
	var arg string

	for i := 0; i < len(args); i++ {
		flag := l.lookup(args[i])
		if flag == nil {
			return ret, fmt.Errorf("Invalid flag provided: %s", args[i])
		}

		i++
		if i < len(args) {
			arg = args[i]
		} else {
			arg = ""
		}

		switch flag.ValueReqt {
		case None:
			if len(arg) != 0 && l.lookup(arg) != nil {
				return ret, fmt.Errorf("The %s/%s flag does not take an argument.", flag.Short, flag.Long)
			}
		case Required:
			if len(arg) == 0 || l.lookup(arg) != nil {
				return ret, fmt.Errorf("The %s/%s flag requires an argument.", flag.Short, flag.Long)
			}
		case Optional:
			if len(arg) != 0 && l.lookup(arg) != nil {
				// Zap this, it's not really an arg
				arg = ""
			}
		}


	}

	return ret, err
}

