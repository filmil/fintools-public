// Package main contains a test parser for ofx files.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aclindsa/ofxgo"
)

var inFile = flag.String("infile", "", "The input ofx file to parse")

func main() {
	flag.Parse()
	if *inFile == "" {
		fmt.Fprintf(os.Stderr, "flag --infile=... is required\n")
		os.Exit(-1)
	}
	f, err := os.Open(*inFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "while opening file: %v: %v\n", *inFile, err)
		os.Exit(-2)
	}
	r, err := ofxgo.ParseResponse(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse contents from: %v: %v\n", *inFile, err)
		os.Exit(-3)
	}
	fmt.Printf("%v\n", r)
}
