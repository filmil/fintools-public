package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/filmil/fintools-public/pkg/buildium"
	"github.com/filmil/fintools-public/pkg/cfg"
	"github.com/filmil/fintools-public/pkg/csv2"
	"github.com/filmil/fintools-public/pkg/tiller"
)

func require(flagName string, value string) {
	if value == "" {
		log.Fatalf("flag --%v=... is required", flagName)
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Converts Buildium CSV Income Statement Detailed into Tiller CSV reports.

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
		inFile, rowsFile, cfgFile string
	)

	flag.StringVar(&inFile, "csv", "", "Input CSV file from Buildium")
	flag.StringVar(&rowsFile, "rows", "", "Output Tiller `Transactions` CSV file")
	flag.StringVar(&cfgFile, "cfg", "", "Configuration file, for automated account categorization, JSON format")

	flag.Parse()

	require("csv", inFile)
	require("rows", rowsFile)
	require("cfg", cfgFile)

	f, err := os.Open(inFile)
	if err != nil {
		log.Fatalf("could not open: %v: %v", inFile, err)
	}

	c, err := csv2.NewCSVData(f)
	if err != nil {
		log.Fatalf("could not parse: %v: %v", inFile, err)
	}

	cr, err := os.Open(cfgFile)
	if err != nil {
		log.Fatalf("could not read config: %v: %v", cfgFile, err)
	}

	lsx, err := cfg.LoadSchema(cr)
	if err != nil {
		log.Fatalf("could not load schema: %v", err)
	}
	cf := cfg.New(lsx)

	tillerRows, err := ConvertTiller(c, cf)
	if err != nil {
		log.Fatalf("could not convert tiller: %v", err)
	}

	rf, err := os.Create(rowsFile)
	defer rf.Close()
	if err != nil {
		log.Fatalf("could not create rows file: %v: %v", rowsFile, err)
	}
	defer rf.Close()
	if err := WriteWriter(rf, tillerRows); err != nil {
		log.Fatalf("%v", err)
	}
}

func ConvertTiller(c *csv2.CSVData, cf *cfg.Instance) ([][]string, error) {
	rowCount := 0
	var tillerRows [][]string
	c.ForEach(0, c.Size(), func(row []string) error {
		var rowStr []string
		if rowCount == 0 {
			rowStr = tiller.ColumnNames
		} else {
			// We got a row.  Convert it.
			var tillerRow tiller.Row
			rowStr = tillerRow.AsCSVRow()
			amount := row[buildium.IndexAmount]
			// Invert income and expense
			if amount[0] == '-' {
				amount = amount[1:]
			} else {
				amount = "-" + amount
			}
			rowStr[tiller.IndexAmount] = amount
			rowStr[tiller.IndexDate] = row[buildium.IndexEntryDate]
			rowStr[tiller.IndexTags] = "Tax"
			rowStr[tiller.IndexDescription] = row[buildium.IndexPayeeName]
			rowStr[tiller.IndexAccountID] = cf.GetID()
			rowStr[tiller.IndexAccountName] = cf.GetID()

			buildiumAccount := row[buildium.IndexAccountID]
			if buildiumAccount == "" {
				return fmt.Errorf("row %d no account: %+v", rowCount+1, rowStr)
			}
			tillerCategory := cf.GetCat(buildiumAccount)
			if tillerCategory == "" {
				return fmt.Errorf("row %d no category for account: %+v", rowCount+1, rowStr)
			}
			rowStr[tiller.IndexCategory] = tillerCategory
		}
		tillerRows = append(tillerRows, rowStr)
		rowCount++
		return nil
	})
	return tillerRows, nil
}

func WriteWriter(rf io.Writer, tillerRows [][]string) error {
	if err := writeFile(rf, func(w io.Writer) error {
		cw := csv.NewWriter(rf)
		defer cw.Flush()
		for _, row := range tillerRows {
			if err := cw.Write(row); err != nil {
				return fmt.Errorf("could not write row: %+v", row)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Write file using the writer function provided.
func writeFile(w io.Writer, writeFn func(w io.Writer) error) error {
	if err := writeFn(w); err != nil {
		return fmt.Errorf("could not write rows file: %w", err)
	}
	return nil
}
