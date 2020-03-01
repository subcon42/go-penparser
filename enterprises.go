package pen

import (
	"bufio"
	"bytes"
	"encoding/asn1"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*
OID, IRI & ASN.1 Prefix Variables
*/
var enterpriseOID string
var enterpriseIRI string
var enterpriseASN1 string = `{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) <--X-->}`

/*
Enterprises represents the entire (parsed) branch of
the IANA Private Enterprise Number (PEN) List. This
type also contains prefix OID and IRI values, as well
as parsing-related data that may be useful to admin
responsible for handling regular refreshes.
*/
type Enterprises struct {
	Nodes     []Node
	ParseTime time.Duration
	URI       *url.URL // Source URI
	Title,
	Section,
	LastUpdated string
}

/*
append will attempt to append unique instance of Node (n)
to the receiver instance of *Enterprises' Nodes field.

Non-unique instances of Node (so declared by a matching OID)
will be silently discarded if an append is attempted.
*/
func (e *Enterprises) append(n Node) bool {
	if exists, _ := e.oidExists(n.Decimal); exists {
		return false
	}
	e.Nodes = append(e.Nodes, n)
	return true
}

/*
oidExists returns a boolean value indicative of whether the
provided node decimal or OID sequence exists within a given
node within the Enterprises receiver object.

This method accepts the following search terms:

 * Single node number as an integer (int)
 * Single node number (string)
 * Stringified OID (1.3.6.1.4.1.<num>)
 * []int (raw integer slices)
 * asn1.ObjectIdentifier

If a match is found, the first return value is the boolean
indicator.  The second return value indicates the storage
index number of the Node in question as reported by the
receiver instance of Enterprises.

If a match is found, the second t
*/
func (e Enterprises) oidExists(dec interface{}) (bool, int) {
	for el := range e.Nodes {
		switch tv := dec.(type) {
		case asn1.ObjectIdentifier:
			return e.oidExists([]int(tv))
		case string:
			if x, err := strconv.Atoi(tv); err == nil {
				return e.oidExists(x)
			}
			if x := strings.Split(tv, `.`); len(x) >= 1 {
				return e.oidExists(x[len(x)-1])
			}
		case int:
			if tv < 0 {
				return false, -1
			}

			if e.Nodes[el].Decimal == tv {
				return true, el
			}
		case []int:
			if len(tv) <= 1 {

				return false, -1
			}

			// Don't bother running another loop if the
			// OID prefix is bogus to begin with ...
			if asn1.ObjectIdentifier(tv[:len(tv)-1]).String() != enterpriseOID {
				return false, -1
			}
			return e.oidExists(tv[len(tv)-1])
		}
	}
	return false, -1
}

/*
FindByOID will conduct a bonafide ASN.1 ObjectIdentifier match between
the provided value and each parsed value found within the Enterprises
receiver instance.
*/
func (e Enterprises) FindByOID(oid interface{}) (Node, bool) {
	if exists, idx := e.oidExists(oid); exists {
		return e.Nodes[idx], exists
	}
	return Node{}, false
}

/*
FindByIRI will conduct a caseless string comparison of all IRI values
observed during a looped search, and the provided IRI value (iri).
*/
func (e Enterprises) FindByIRI(iri string) (Node, bool) {
	iri = strings.ToLower(iri)

	for i := 0; i < e.Count(); i++ {
		n := e.Nodes[i]
		target := strings.ToLower(n.IRI())
		if iri == target {
			return n, true
		}
	}
	return Node{}, false
}

/*
FindByEmail performs a caseless match between the provided email
address (email) and each discovered email address within the OID
index. For search convenience, it is unnecessary to replace the
ampersand (`&`) with the so-called "Commercial At Sign" (`@`),
as this is done under-the-hood per email search request.

If found an instance of Node is returned along with an affirmative
boolean value; else an empty node and a negative boolean value.
*/
func (e Enterprises) FindByEmail(email string) (Node, bool) {
	for i := 0; i < e.Count(); i++ {
		for em := 0; em < len(e.Nodes[i].Email); em++ {
			email = strings.ReplaceAll(email, `&`, `@`)
			target := strings.ReplaceAll(e.Nodes[i].Email[em], `&`, `@`)

			if strings.ToLower(email) == strings.ToLower(target) {
				return e.Nodes[i], true
			}
		}
	}
	return Node{}, false
}

/*
FindByContact will conduct a caseless name-based match between
each Contact name found within the Enterprises receiver instance
and the provided name input argument.
*/
func (e Enterprises) FindByContact(name string) (Node, bool) {
	for i := 0; i < e.Count(); i++ {
		name = strings.ReplaceAll(name, ` `, ``)
		target := strings.ReplaceAll(e.Nodes[i].Contact, ` `, ``)

		if strings.ToLower(name) == strings.ToLower(target) {
			return e.Nodes[i], true
		}
	}
	return Node{}, false
}

func (e *Enterprises) setLastUpdated(lu line) bool {
	if lu.len() <= 1 {
		return false
	}

	lus := strings.Split(lu.string()[1:lu.len()-1], ` `)
	if len(lus) == 0 {
		return false
	}

	e.LastUpdated = lus[len(lus)-1]
	return true
}

func (e *Enterprises) setSection(sec line) bool {
	if sec.len() <= 1 {
		return false
	}

	e.Section = sec.string()[0 : sec.len()-1]
	return true
}

func (e *Enterprises) setPrefix(pfx line) bool {
	if pfx.len() <= 7 {
		return false
	}

	npfx := line(pfx[8:pfx.len()])

	pfxs := strings.Split(npfx.string(), ` `)
	if len(pfxs) >= 2 {
		if len(pfxs[0])|len(pfxs[1]) <= 2 {
			return false
		}

		// Set our global vars with IRI and OID info
		enterpriseIRI = `/` + strings.ReplaceAll(pfxs[0], `.`, `/`)
		enterpriseOID = pfxs[1][1 : len(pfxs[1])-1]

		return true
	}

	return false
}

func (e *Enterprises) setURI(uri line) bool {
	if uri.len() <= 1 {
		return false
	}

	var err error
	f := strings.Split(uri.string(), ` `)
	e.URI, err = url.Parse(f[len(f)-1])
	if err != nil {
		return false
	}

	return true
}

/*
Count returns the number of Node instances present within
the receiver instance of Enterprises.
*/
func (e Enterprises) Count() int {
	return len(e.Nodes)
}

/*
DumpHeader is a stringer method for the Enterprises receiver
instance header-information.
*/
func (e Enterprises) DumpHeader() (head string) {
	head += fmt.Sprintf("\n%s\n\n", ul(e.Section+`:`))

	head += fmt.Sprintf("Parse Info:\n")

	if e.URI != nil {
		head += fmt.Sprintf("  URI: %s\n", e.URI.String())
	}

	head += fmt.Sprintf("  Entries: %d\n", e.Count())
	head += fmt.Sprintf("  Duration: %d sec.\n", e.ParseTime/time.Second)
	head += fmt.Sprintf("  Last Updated: %s\n\n", e.LastUpdated)

	head += fmt.Sprintf("Prefix Info:\n")
	head += fmt.Sprintf("  OID: %s\n", enterpriseOID)
	head += fmt.Sprintf("  IRI: %s\n", enterpriseIRI)
	head += fmt.Sprintf("  ASN: %s\n", strings.Replace(enterpriseASN1, ` <--X-->`, ``, 1))

	return
}

func (e *Enterprises) setHeader(l line, ct int) (bool, error) {

	switch ct - 1 {
	case 1:
		e.Title = l.string() // no special processing needed
	case 3:
		if ok := e.setLastUpdated(l); !ok {
			return false, errors.New("Unable to set LastUpdated header value")
		}
	case 5:
		if ok := e.setSection(l); !ok {
			return false, errors.New("Unable to set Section header value")
		}
	case 7:
		if ok := e.setPrefix(l); !ok {
			return false, errors.New("Unable to set Prefix header value")
		}
	case 9:
		if ok := e.setURI(l); !ok {
			return false, errors.New("Unable to set URI header value")
		}
	}
	return true, nil
}

/*
New parses the file specified via input argument as the complete
IANA Private Enterprise Numbers List.

If at any point parsing encounters an error, it is returned alongside
a likely nil instance of the *Enterprises type. Else, a fully-populated
instance of *Enterprises shall be returned alongside a nil error.

Note that you must download the IANA Private Enterprise Numbers List
yourself (this package will not do that part for you).
*/
func New(file string) (ents *Enterprises, err error) {

	var startParse int64 = time.Now().UnixNano()

	penBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	ents = new(Enterprises)
	scan := bufio.NewScanner(bytes.NewReader(penBytes))

	ct := 0
	for scan.Scan() {
		ct++
		L := line(scan.Text())
		if L.isZero() {
			continue
		}

		// Lines 0 - 10 are for header info
		if ct <= 10 {
			if _, err := ents.setHeader(L, ct); err != nil {
				return nil, err
			}
		}

		// Any line that is wholly numerical indicates
		// the start of a new entry ...
		if L.isNumbersOnly() {
			if n, err := parseNode(scan, L); err == nil {
				_ = ents.append(n) // duplicates silently ignored ...
			} else {
				return nil, err
			}
		}
	}

	doneParsed := time.Now().UnixNano()
	ents.ParseTime = time.Duration(doneParsed - startParse)

	return
}
