/*
Package pen parses the user-retrieved IANA Private Enteprise
Numbers (PEN) file by path/filename.

Credits

Jesse Coretta (2/29/20)

Advisory

You must download the PEN file yourself using your preferred HTTP
client.  The New() method takes the path of that downloaded file.

Keep in mind, this is a very rough and unofficial draft; subject
to change without notice!

Usage

Basic usage is described as follows:

  func main() {

        // Update this to reference your freshly downloaded
        // IANA PEN file (see ents.URI() in header).
        var file string = `/tmp/pen`

        // Create our *Enterprises object based on
        // data parsed via file
        ents, err := New(file)
        if err != nil {
                fmt.Println(err)
                return
        }

        // Print our header (useful data)
        fmt.Println(ents.DumpHeader())

        //////////////////////////////////////////////////////////
        // Begin our basic demo ...

        fmt.Println(`Results`)

        // Conduct a search for a Node. Note: you have more than
        // one means of searching for Node instances ...
        myNode, ok := ents.FindByOID(asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 54399}) // Search by asn1.ObjectIdentifier ...
        //myNode, ok := ents.FindByIRI(`/iso/org/dod/internet/private/enterprise/54399`) // or by string Internationalized Resource Identifier (IRI) path notation ...
        //myNode, ok := ents.FindByOID([]int{1, 3, 6, 1, 4, 1, 54399})                   // or by []int cast of asn1.ObjectIdentifier ...
        //myNode, ok := ents.FindByOID(`1.3.6.1.4.1.54399`)                              // or by stringer of asn1.ObjectIdentifier ...
        //myNode, ok := ents.FindByOID(54399)                                            // or by leaf-node decimal (far-right digit)
        if !ok {
                return // no match!
        }

        // Print our retrieved Node
        fmt.Println(myNode.DumpNode())

        // Alt. search options
        // Find By Email
        //myNode, ok = ents.FindByEmail(`subcon.co.42&gmail.com`)
        //fmt.Printf("Found by email: %t\n", ok)

        // Find By Contact
        //myNode, ok = ents.FindByContact(`Jesse Coretta`)
        //fmt.Printf("Found by contact name: %t\n", ok)

  }

*/
package pen
