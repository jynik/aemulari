package aemulari

import (
	"bytes"
	"fmt"
	"strings"
)

type Permissions struct {
	Read  bool
	Write bool
	Exec  bool
}

func (p *Permissions) Set(s string) error {
	s = strings.Trim(strings.ToLower(s), " \t\r\n")
	if len(strings.Trim(s, "rwx")) != 0 {
		return fmt.Errorf("Invalid permissions string: %s", s)
	}

	p.Read = strings.Contains(s, "r")
	p.Write = strings.Contains(s, "w")
	p.Exec = strings.Contains(s, "x")

	return nil
}

func (p Permissions) String() string {
	var buf bytes.Buffer

	if p.Read {
		buf.WriteRune('r')
	}

	if p.Write {
		buf.WriteRune('w')
	}

	if p.Exec {
		buf.WriteRune('x')
	}

	if buf.Len() == 0 {
		return "(none)"
	}

	return buf.String()
}
