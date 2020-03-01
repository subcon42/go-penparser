package pen

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

/*
Node represents each distinct occurrence of a given Node
entry within the IANA Private Enterprises Numbers List.
The number of fields directly matches the number of lines
found for each Node (and consequently how many bufio.Scan
iterations are conducted per Node during parsing).

As shown on the legend/diagram, IANA describes each field
and the leading lengths of subfields:

 Decimal
 | Organization
 | | Contact
 | | | Email
 | | | |

... which would equate to ...

  nodeNum             // Entry Line 1 :: zero (0) leading spaces
  __Organization      // Entry Line 2 :: two (2) leading spaces
  ____Contact         // Entry Line 3 :: four (4) leading spaces
  ______Email         // Entry Line 4 :: six (6) leading spaces

*/
type Node struct {
	Email []string
	Contact,
	Organization string
	Decimal int // aka node number
}

/*
OID returns the stringified ASN.1 Object Identifier (OID)
value of the receiver Node.
*/
func (n Node) OID() string {
	return fmt.Sprintf("%s.%d", enterpriseOID, n.Decimal)
}

/*
IRI returns the stringified Internationalized Resource
Identifier (IRI) value of the receiver Node.
*/
func (n Node) IRI() string {
	return fmt.Sprintf("%s/%d", enterpriseIRI, n.Decimal)
}

/*
ASN returns the stringified ASN.1 Path Notation value
of the receiver Node.
*/
func (n Node) ASN() string {
	return fmt.Sprintf("%s", strings.Replace(enterpriseASN1, `<--X-->`, fmt.Sprint(n.Decimal), 1))
}

/*
Emails returns the stringified email address(es) found for
the receiver node. An optional boolean value of true will
cause indented email address output (helpful for multi-
valued entries).
*/
func (n Node) Emails(split ...bool) (e string) {
	if len(n.Email) > 0 {
		if len(split) == 0 {
			split = []bool{false}
		}
		if split[0] {
			for i := 0; i < len(n.Email); i++ {
				e += fmt.Sprintf("  - %s\n", n.Email[i])
			}
		} else {
			e += fmt.Sprintf("%s\n", strings.Join(n.Email, `,`))
		}
	}
	return
}

/*
DumpNode is a stringer method for the receiver instance of
Node.
*/
func (n Node) DumpNode() (entry string) {
	entry += fmt.Sprintf("\n## %s (%d)\n", n.Organization, n.Decimal)
	entry += fmt.Sprintf("  Contact: %s\n", n.Contact)

	entry += fmt.Sprintf("  OID: %s\n", n.OID())
	entry += fmt.Sprintf("  IRI: %s\n", n.IRI())
	entry += fmt.Sprintf("  ASN: %s\n", n.ASN())

	entry += "  Email:\n" + n.Emails(true)

	return
}

func parseNode(scan *bufio.Scanner, l line) (n Node, err error) {
	n.Decimal, err = strconv.Atoi(l.string())
	if err != nil {
		return Node{}, err
	}

	for i := 1; i < 4; i++ {
		if scan.Scan() {
			next := line(scan.Text()).trimLeadingSpace()
			switch i {
			case 1:
				n.Organization = next.string()
			case 2:
				n.Contact = next.string()
			case 3:
				em := strings.Split(strings.ReplaceAll(next.string(), ` `, ``), `,`)
				n.Email = make([]string, 0, len(em))
				for el := range em {
					n.Email = append(n.Email, em[el])
				}
			}
		}
	}
	return
}
