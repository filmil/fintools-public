package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/filmil/fintools/pkg/csv2"
)

func main() {
	var (
		inFile string
	)

	flag.StringVar(&inFile, "csv", "", "Input CSV file")

	flag.Parse()

	if inFile == "" {
		log.Fatalf("flag --csv=... is required")
	}

	f, err := os.Open(inFile)
	if err != nil {
		log.Fatalf("could not open: %v: %v", inFile, err)
	}

	c, err := csv2.NewCSVData(f)
	if err != nil {
		log.Fatalf("could not parse: %v: %v", inFile, err)
	}

	fmt.Printf("result:\n%v", c)

	r, err := csv2.NewReport(c)
	if err != nil {
		log.Fatalf("could not create report: %v", err)
	}

	fmt.Printf("report:\n%+v", *r)

}
