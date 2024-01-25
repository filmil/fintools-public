package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/filmil/fintools/pkg/cfg"
	"github.com/filmil/fintools/pkg/csv2"
	"github.com/filmil/fintools/pkg/index"
	"github.com/filmil/fintools/pkg/tiller"
)

func require(flagName string, value string) {
	if value == "" {
		log.Fatalf("flag --%v=... is required", flagName)
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Converts Buildium CSV reports into Tiller CSV reports.

Use to convert from Buildium into Tiller.  All flag parameters are required. Sample
config file is below.

{
  "account_id": "manual:some_account",
  "account_map": [
    {
      "original": "Utilities - Common SD1/Sewer",
      "category": "XX - Category",
      "id": "manual:foo"
    }
  ]
}

Args:
`)

		flag.PrintDefaults()
	}
}

func main() {
	var (
		// Obvious variable names are obvious.
		inFile, rowsFile, balancesFile, cfgFile string
	)

	flag.StringVar(&inFile, "csv", "", "Input CSV file from Buildium")
	flag.StringVar(&rowsFile, "rows", "", "Output Tiller `Transactions` CSV file")
	flag.StringVar(&balancesFile, "bal", "", "Output Tiller `Balances` CSV file")
	flag.StringVar(&cfgFile, "cfg", "", "Configuration file, for automated account categorization, JSON format")

	flag.Parse()

	require("csv", inFile)
	require("rows", rowsFile)
	require("bal", balancesFile)
	require("cfg", cfgFile)

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

	cr, err := os.Open(cfgFile)
	if err != nil {
		log.Fatalf("could not read config: %v: %v", cfgFile, err)
	}

	lsx, err := cfg.LoadSchema(cr)
	if err != nil {
		log.Fatalf("could not load schema: %v", err)
	}
	cfx := cfg.New(lsx)
	t := tiller.New(i, cfx)

	if err := writeFile(rowsFile, t.WriteRows); err != nil {
		log.Fatalf("%v", err)
	}
	if err := writeFile(balancesFile, t.WriteBalances); err != nil {
		log.Fatalf("%v", err)
	}
}

// Write file using the writer function provided.
func writeFile(fn string, writeFn func(w io.Writer) error) error {
	rf, err := os.Create(fn)
	if err != nil {
		return fmt.Errorf("could not create rows file: %v: %w", fn, err)
	}
	defer rf.Close()
	if err := writeFn(rf); err != nil {
		return fmt.Errorf("could not write rows file: %v: %w", fn, err)
	}
	return nil
}
