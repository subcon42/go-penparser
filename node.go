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
func (n Node) Emails() (e string) {
	if len(n.Email) > 0 {
		e = fmt.Sprintf("%s", strings.Join(n.Email, `,`))
	}
	return
}

/*
Node returns a map[string]string containing key pieces of information
regarding the receiver Node instance.
*/
func (n Node) Node() map[string]string {
	return map[string]string{
		`Organization`: n.Organization,
		`Contact`: n.Contact,
		`Decimal`: fmt.Sprint(n.Decimal),
		`Emails`: n.Emails(),
		`OID`: n.OID(),
		`IRI`: n.IRI(),
		`ASN`: n.ASN(),
	}
}

func parseNode(scan *bufio.Scanner, l line) (n Node, err error) {
	n.Decimal, err = strconv.Atoi(l.string())
	if err != nil {
		return emptyNode, err
	}

	// we make an assumption here that each node entry is
	// no more and no less than four (4) lines. This has
	// always been the case for this file, so its fair to
	// operate as such ...
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
