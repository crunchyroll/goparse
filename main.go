// goparse parses Go source code with the parser.Trace option set,
// yielding a trace of parsed productions.  This is ideal for human
// consumption and debugging, aiding in the understanding of the
// structure of Go programs as well as the behavior of the library
// parser.
//
// If a command line argument is specified, it will be used as the input
// file.  If no arguments are specified, the input will be read from
// standard in.
//
//	usage: goparse [file]
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"go/parser"
	"go/token"
)

func main() {
	var filename string
	var src io.Reader
	var err error

	if len(os.Args) == 1 {
		filename = "-"
		src = os.Stdin
	} else if len(os.Args) == 2 {
		filename = os.Args[1]
		src, err = os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("usage: %s [file]\n", path.Base(os.Args[0]))
		os.Exit(2)
	}

	fset := token.NewFileSet()
	mode := parser.Trace
	_, err = parser.ParseFile(fset, filename, src, mode)
	if err != nil {
		log.Fatal(err)
	}
}
