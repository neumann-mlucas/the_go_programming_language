package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

func main() {
	dec := xml.NewDecoder(os.Stdin)
	root := &Element{}
	root, err := RDParse(dec, root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Println(root)
}

type Node interface{} // CharData or *Element

type CharData string

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func (n *Element) String() string {
	b := &bytes.Buffer{}
	visit(n, b, 0)
	return b.String()
}

func visit(n Node, w io.Writer, depth int) {
	switch n := n.(type) {
	case *Element:
		fmt.Fprintf(w, "%*s%s %s\n", depth*2, "", n.Type.Local, n.Attr)
		for _, c := range n.Children {
			visit(c, w, depth+1)
		}
	case CharData:
		fmt.Fprintf(w, "%*s%q\n", depth*2, "", n)
	default:
		panic(fmt.Sprintf("got %T", n))
	}
}

func RDParse(dec *xml.Decoder, tree *Element) (*Element, error) {
	tok, err := dec.Token()
	if err == io.EOF {
		return tree, nil
	} else if err != nil {
		return nil, err
	}

	switch tok := tok.(type) {
	case xml.StartElement:
		children := &Element{tok.Name, tok.Attr, []Node{}}
		children, _ = RDParse(dec, children) // push
		if err == nil {
			tree.Children = append(tree.Children, children)
		}
	case xml.EndElement:
		return tree, nil // pop
	case xml.CharData:
		children := CharData(tok)
		tree.Children = append(tree.Children, children)
	}
	return RDParse(dec, tree) // go to the next token
}
