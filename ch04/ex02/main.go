package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"os"
)

func main() {
	var useSHA384 = flag.Bool("SHA384", false, "use SHA384")
	var useSHA512 = flag.Bool("SHA512", false, "use SHA512")
	flag.Parse()

	b := bytes.Buffer{}
	b.ReadFrom(os.Stdin)

	if *useSHA512 {
		checksum := sha512.Sum512(b.Bytes())
		fmt.Printf("SHA512: %x", checksum)
	} else if *useSHA384 {
		checksum := sha512.Sum384(b.Bytes())
		fmt.Printf("SHA384: %x", checksum)
	} else {
		checksum := sha256.Sum256(b.Bytes())
		fmt.Printf("SHA256: %x", checksum)
	}
}
