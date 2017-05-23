package cmdline

type appender []string

func (a *appender) String() string {
	var ret = ""
	for _, s := range *a {
		ret += s
	}
	return ret
}

func (a *appender) Set(s string) error {
	*a = append(*a, s)
	return nil
}
