package aemulari

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Not actually backed by map because
type RegisterMap struct {
	regList []*RegisterDef          // Sorted list of primary register names
	regMap  map[string]*RegisterDef // Random access
}

// Register definitions should be added in the desired display order.
func (r *RegisterMap) add(names []string, reg *RegisterDef) {
	if r.regMap == nil {
		r.regMap = make(map[string]*RegisterDef)
	}

	r.regList = append(r.regList, reg)

	for _, name := range names {
		r.regMap[name] = reg
	}
}

func (rm *RegisterMap) register(name string) (*RegisterDef, error) {
	if reg, found := rm.regMap[name]; found {
		return reg, nil
	}

	return nil, fmt.Errorf("\"%s\" is not a valid register name.", name)
}

func (rm *RegisterMap) registers() []*RegisterDef {
	count := len(rm.regList)
	regs := make([]*RegisterDef, count, count)
	copy(regs, rm.regList)
	return regs
}

func (rm *RegisterMap) ParseRegister(s string) (RegisterValue, error) {
	var rv RegisterValue
	fields := strings.Split(strings.ToLower(strings.Trim(s, "\r\n\t ")), "=")

	if len(fields) != 2 {
		return rv, fmt.Errorf("\"%s\" is not a valid register assignment.", s)
	}

	val, err := strconv.ParseUint(fields[1], 0, 64)
	if err != nil {
		return rv, err
	}

	if reg, err := rm.register(fields[0]); err == nil {
		rv.Value = reg.mask & val
		rv.Reg = reg
		return rv, nil
	} else {
		return rv, err
	}
}

func (rm *RegisterMap) ParseRegisters(strs []string) ([]RegisterValue, error) {
	var ret []RegisterValue

	for _, str := range strs {
		if rv, err := rm.ParseRegister(str); err != nil {
			return ret, err
		} else {
			ret = append(ret, rv)
		}
	}

	return ret, nil
}

func (rm *RegisterMap) RegisterRegexp() *regexp.Regexp {
	restr := "(^|[^[:alpha:]])("

	for name, _ := range rm.regMap {
		restr += name + "|"
	}

	return regexp.MustCompile(restr[:len(restr)-1] + ")")
}
