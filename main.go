// goparse parses Go source code.  You may independently enable parser
// tracing, which yields a trace of parsed productions; and AST
// printing, which gives a lower-level view of the resultant AST.  Both
// are ideal for programmer consumption and debugging, aiding in
// understanding the structure of Go programs as well as the behavior of
// the built-in parser library.
//
// If a positional argument is specified, it will be used as the input
// file.  If no positional arguments are specified, the input will be
// read from standard in.
//
// There are other options as well that enable other parser
// functionality.
//
//	usage: goparse [-h] [options] [file]
//	  -all-errors=false: report all errors (not just the first 10
//	                     on different lines)
//	  -ast-print=false: print AST with ast.Fprint and no field filter
//	  -declaration-errors=false: report declaration errors
//	  -imports-only=false: stop parsing after import declarations
//	  -parse-comments=false: parse comments and add them to AST
//	  -trace=false: print a trace of parsed productions
//
// If both tracing and AST printing are enabled, the trace will precede
// AST output.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"go/ast"
	"go/parser"
	"go/token"
)

var mode parser.Mode
var astPrint *bool

func init() {
	importsOnly       := flag.Bool("imports-only",       false,
		"stop parsing after import declarations")
	parseComments     := flag.Bool("parse-comments",     false,
		"parse comments and add them to AST")
	trace             := flag.Bool("trace",              false,
		"print a trace of parsed productions")
	declarationErrors := flag.Bool("declaration-errors", false,
		"report declaration errors")
	allErrors         := flag.Bool("all-errors",         false,
		"report all errors (not just the first 10 on different lines)")
	astPrint           = flag.Bool("ast-print",          false,
		"print AST with ast.Fprint and no field filter")
	flag.Parse()
	if (*importsOnly) {
		mode |= parser.ImportsOnly
	}
	if (*parseComments) {
		mode |= parser.ParseComments
	}
	if (*trace) {
		mode |= parser.Trace
	}
	if (*declarationErrors) {
		mode |= parser.DeclarationErrors
	}
	if (*allErrors) {
		mode |= parser.AllErrors
	}
}

func main() {
	var filename string
	var src io.Reader
	var err error
	var astfile *ast.File

	fargs := flag.Args()
	if len(fargs) == 0 {
		filename = "-"
		src = os.Stdin
	} else if len(fargs) == 1 {
		filename = fargs[0]
		src, err = os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("usage: %s [-h] [options] [file]\n", path.Base(os.Args[0]))
		os.Exit(2)
	}

	fset := token.NewFileSet()
	astfile, err = parser.ParseFile(fset, filename, src, mode)
	if err != nil {
		log.Fatal(err)
	}

	if *astPrint {
		err = ast.Fprint(os.Stdout, fset, astfile, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
