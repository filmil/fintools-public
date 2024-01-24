package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/filmil/fintools/pkg/cfg"
	"github.com/filmil/fintools/pkg/csv2"
	"github.com/filmil/fintools/pkg/index"
	"github.com/filmil/fintools/pkg/tiller"
)

func main() {
	var (
		inFile, rowsFile, balancesFile string
	)

	flag.StringVar(&inFile, "csv", "", "Input CSV file")
	flag.StringVar(&rowsFile, "rows", "", "Output `rows` file")
	flag.StringVar(&balancesFile, "bal", "", "Output `balances` file")

	flag.Parse()

	if inFile == "" {
		log.Fatalf("flag --csv=... is required")
	}
	if rowsFile == "" {
		log.Fatalf("flag --rows=... is required")
	}
	if balancesFile == "" {
		log.Fatalf("flag --bal=... is required")
	}

	f, err := os.Open(inFile)
	if err != nil {
		log.Fatalf("could not open: %v: %v", inFile, err)
	}

	c, err := csv2.NewCSVData(f)
	if err != nil {
		log.Fatalf("could not parse: %v: %v", inFile, err)
	}

	r, err := csv2.NewReport(c)
	if err != nil {
		log.Fatalf("could not create report: %v", err)
	}

	i := index.New(r)

	// TODO: filmil- load the config from a file instead
	lsx, err := cfg.LoadSchema(strings.NewReader(`
		{
		  "account_id": "manual:some_account",
		  "account_map": [
			{
				"original": "Utilities - Common SD1/Sewer",
				"category": "XX - Category"
			}
		  ]
		}

	`))
	if err != nil {
		log.Fatalf("could not load schema: %v", err)
	}
	cfx := cfg.New(lsx)
	t := tiller.New(i, cfx)

	rf, err := os.Create(rowsFile)
	if err != nil {
		log.Fatalf("could not create rows file: %v: %v", rowsFile, err)
	}
	defer rf.Close()
	if err := t.WriteRows(rf); err != nil {
		log.Fatalf("could not write rows file: %v: %v", rowsFile, err)
	}

	bf, err := os.Create(balancesFile)
	if err != nil {
		log.Fatalf("could not create balances file: %v: %v", balancesFile, err)
	}
	defer bf.Close()
	if err := t.WriteBalances(bf); err != nil {
		log.Fatalf("could not create balances file: %v: %v", balancesFile, err)
	}
}
