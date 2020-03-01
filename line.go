package pen

import (
	"strings"
)

/*
line represents a single line from the parsed input.
*/
type line string

/*
String is a Stringer method for line, which returns the
string form of line.
*/
func (l line) string() string {
	return string(l)
}

/*
IsZero returns a boolean value indicative of whether
the receiver instance of line is determined to be zero.
*/
func (l line) isZero() bool {
	return len(l) == 0
}

/*
Len returns the integer-based length of the receiver
instance of line.
*/
func (l line) len() int {
	return len(l)
}

/*
TrimLeadingSpace returns a copy of the receiver instance
of line that lacks any leading whitespace that was present
before.
*/
func (l line) trimLeadingSpace() line {
	return line(strings.TrimLeft(l.string(), ` `))
}

/*
IsNumbersOnly returns a boolean value indicative of whether
the receiver instance of line contains ONLY digits. In IANA
PEN terms, this indicates the start of a new OID entry.
*/
func (l line) isNumbersOnly() bool {
	for ch := range l {
		if '0' <= l[ch] && l[ch] <= '9' {
			continue
		}
		return false
	}
	return true
}

func ul(str string) (ul string) {
	ul = str[:] + "\n"
	for i := 0; i < len(str); i++ {
		ul += `-`
	}
	return
}
