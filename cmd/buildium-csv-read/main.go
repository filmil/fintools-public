package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/filmil/fintools/pkg/cfg"
	"github.com/filmil/fintools/pkg/csv2"
	"github.com/filmil/fintools/pkg/index"
	"github.com/filmil/fintools/pkg/tiller"
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

	log.Printf("result:\n%v", c)

	r, err := csv2.NewReport(c)
	if err != nil {
		log.Fatalf("could not create report: %v", err)
	}

	log.Printf("report:\n%+v", spew.Sdump(*r))

	i := index.New(r)

	log.Printf("index:\n%+v", spew.Sdump(*i))

	lsx, err := cfg.LoadSchema(strings.NewReader(`
		{
		  "account_id": "manual:some_account",
		  "account_map": [
			{
				"original": "id",
				"category": "cat"
			}
		  ]
		}

	`))
	if err != nil {
		log.Fatalf("could not load schema: %v", err)
	}
	cfx := cfg.New(lsx)
	t := tiller.New(i, cfx)

	log.Printf("tiller export:\n%v", spew.Sdump(*t))

}
