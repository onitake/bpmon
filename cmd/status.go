package cmd

type status int

const (
	StatusOK = iota
	StatusNOK
	StatusUnknown
)

func (s status) String() string {
	var out string
	switch s {
	case StatusOK:
		out = "ok"
	case StatusNOK:
		out = "not ok"
	case StatusUnknown:
		out = "unknown"
	}
	return out
}

func (s status) toInt() int {
	return int(s)
}

func (s status) Colorize(in string) string {
	var out string
	switch s {
	case StatusOK:
		out = "\x1b[32;1m" + in + "\x1b[0m"
	case StatusNOK:
		out = "\x1b[31;1m" + in + "\x1b[0m"
	case StatusUnknown:
		out = "\x1b[35;1m" + in + "\x1b[0m"
	}
	return out
}

func boolAsStatus(ok bool) status {
	if ok {
		return StatusOK
	} else {
		return StatusNOK
	}
}
