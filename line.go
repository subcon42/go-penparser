package pen

import (
	"strings"
)

/*
Line represents a single line from the parsed input.
*/
type Line string

/*
Bytes returns the byte form of Line.
*/
func (l Line) Bytes() []byte {
	return []byte(l)
}

/*
String is a Stringer method for Line, which returns the
string form of Line.
*/
func (l Line) String() string {
	return string(l)
}

/*
IsZero returns a boolean value indicative of whether
the receiver instance of Line is determined to be zero.
*/
func (l Line) IsZero() bool {
	return len(l) == 0
}

/*
Len returns the integer-based length of the receiver
instance of Line.
*/
func (l Line) Len() int {
	return len(l)
}

/*
StartsWith returns a boolean value indicative of whether
the receiver instance of Line begins with the provided
prefix (s).
*/
func (l Line) StartsWith(s string) bool {
	return strings.HasPrefix(l.String(), s)
}

/*
Contains returns a boolean value indicative of whether
the receiver instance of Line contains the provided
prefix (s).
*/
func (l Line) Contains(s string) bool {
	return strings.Contains(l.String(), s)
}

/*
TrimLeadingSpace returns a copy of the receiver instance
of Line that lacks any leading whitespace that was present
before.
*/
func (l Line) TrimLeadingSpace() Line {
	return Line(strings.TrimLeft(l.String(), ` `))
}

/*
IsNumbersOnly returns a boolean value indicative of whether
the receiver instance of Line contains ONLY digits. In IANA
PEN terms, this indicates the start of a new OID entry.
*/
func (l Line) IsNumbersOnly() bool {
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
