package aemulari

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type registerMap struct {
	regList []*registerAttr          // Sorted by register name
	regMap  map[string]*registerAttr // Random access
}

// Register definitions should be added in the desired display order.
func (r *registerMap) add(names []string, reg *registerAttr) {
	if r.regMap == nil {
		r.regMap = make(map[string]*registerAttr)
	}

	r.regList = append(r.regList, reg)

	for _, name := range names {
		r.regMap[name] = reg
	}
}

func (rm *registerMap) register(name string) (*registerAttr, error) {
	if reg, found := rm.regMap[name]; found {
		return reg, nil
	}

	return nil, fmt.Errorf("\"%s\" is not a valid register name.", name)
}

func (rm *registerMap) registers() []*registerAttr {
	count := len(rm.regList)
	regs := make([]*registerAttr, count, count)
	copy(regs, rm.regList)
	return regs
}

func (rm *registerMap) ParseRegister(s string) (Register, error) {
	var reg Register
	fields := strings.Split(strings.ToLower(strings.Trim(s, "\r\n\t ")), "=")

	if len(fields) != 2 {
		return reg, fmt.Errorf("\"%s\" is not a valid register assignment.", s)
	}

	val, err := strconv.ParseUint(fields[1], 0, 64)
	if err != nil {
		return reg, err
	}

	attr, err := rm.register(fields[0])
	reg.Value = reg.attr.mask & val
	reg.attr = attr
	return reg, err
}

func (rm *registerMap) ParseRegisters(strs []string) ([]Register, error) {
	var ret []Register

	for _, str := range strs {
		if rv, err := rm.ParseRegister(str); err != nil {
			return ret, err
		} else {
			ret = append(ret, rv)
		}
	}

	return ret, nil
}

func (rm *registerMap) RegisterRegexp() *regexp.Regexp {
	restr := "(^|[^[:alpha:]])("

	for name, _ := range rm.regMap {
		restr += name + "|"
	}

	return regexp.MustCompile(restr[:len(restr)-1] + ")")
}
