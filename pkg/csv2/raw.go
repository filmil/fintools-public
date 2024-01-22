package csv2

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/big"
	"strings"
)

const NotFound = -1

var NoData = ""

// CSVData contains the parsed CSV values.
type CSVData struct {
	data [][]string
}

func (d CSVData) String() string {
	var ret []string
	for n, l := range d.data {
		ret = append(ret, fmt.Sprintf("%010d: %v", n, strings.Join(l, ",")))
	}
	return strings.Join(ret, "\n")
}

// Finds the first cell beginning row that contains `s`.  Returns NotFound if none found.
func (d CSVData) FindFirst(s string) int {
	return d.FindFrom(0, s)
}

// Gets a full row of values.
func (d CSVData) Row(i int) []string {
	if i >= len(d.data) || i < 0 {
		// No such row.
		return nil
	}
	return d.data[i]
}

// Cell returns the value of data at row and column, or NoData if not present or empty.
func (d CSVData) Cell(r, c int) string {
	if r < 0 || r >= len(d.data) {
		return NoData
	}
	l := &d.data[r]
	if c < 0 || c >= len(*l) {
		return NoData
	}
	return d.data[r][c]
}

// FindFrom is finding s, starting from row i.  Returns NotFound if none found.
func (d CSVData) FindFrom(f int, s string) int {
	return d.FindFromTo(f, len(d.data), s)
}

func (d CSVData) norm(f, t int) (int, int) {
	if f < 0 {
		f = 0
	}
	if t > len(d.data) {
		t = len(d.data)
	}
	if f > t {
		f, t = t, f
	}
	return f, t
}

func (d CSVData) FindFromTo(f, t int, s string) int {
	f, t = d.norm(f, t)
	for a := f; a < t && a < len(d.data); a++ {
		if len(d.data[a]) == 0 {
			continue
		}
		if strings.Contains(d.data[a][0], s) {
			return a
		}
	}
	return NotFound
}

// Calls `fun` for each row from f (inclusive) to t (exclusive).
func (d CSVData) ForEach(f, t int, fun func([]string) error) error {
	f, t = d.norm(f, t)
	for a := f; a < t; a++ {
		if err := fun(d.data[a]); err != nil {
			return fmt.Errorf("ForEeach: row: %v: %w", a, err)
		}
	}
	return nil
}

// / NewCSVData parses the CSV from the supplied reader.
func NewCSVData(r io.Reader) (*CSVData, error) {
	var ret CSVData

	cr := csv.NewReader(r)

	// Set variable number of fields per record.
	cr.FieldsPerRecord = -1

	for run := true; run; {
		r, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				run = false
				continue
			}
			return nil, fmt.Errorf("parse error: %w", err)
		}
		ret.data = append(ret.data, r)
	}
	return &ret, nil
}

func USDToBFloat(s string) (*big.Float, error) {
	if s == "" {
		s = "$0.00"
	}
	s = s[1:] // Skip $

	bf, _, err := big.ParseFloat(s, 10, 6, big.ToNearestEven)
	if err != nil {
		return nil, fmt.Errorf("could not parse: %q: %w", s, err)
	}
	return bf, nil
}

type Transaction struct {
	Date, Type, Unit, Name, Description string
	Debit, Credit, Balance              big.Float
}

// ParseTransaction extracts one transaction or errors out.
func ParseTransaction(l []string) (Transaction, error) {
	var t Transaction
	t.Date = l[0]
	t.Type = l[1]
	t.Unit = l[2]
	t.Name = l[3]
	t.Description = l[4]
	dStr, err := USDToBFloat(l[5])
	if err != nil {
		return t, fmt.Errorf("dStr: %w", err)
	}
	t.Debit = *dStr
	cStr, err := USDToBFloat(l[6])
	if err != nil {
		return t, fmt.Errorf("cStr: %w", err)
	}
	t.Credit = *cStr
	bStr, err := USDToBFloat(l[7])
	if err != nil {
		return t, fmt.Errorf("bStr: %w", err)
	}
	t.Balance = *bStr
	return t, nil
}

type AccountType int

const (
	AssetType AccountType = 1
	LiabilityType
)

type Account struct {
	Name         string
	Type         AccountType
	BeginBalance big.Float
	EndBalance   big.Float
	Transactions []Transaction
}

// Report is a structured report from an account.
type Report struct {
	Name     string
	Accounts []Account
}

func NewReport(from *CSVData) (*Report, error) {
	const (
		AccountNameCol = 0
		BalanceCol     = 7
	)

	var ret Report
	ret.Name = from.Cell(2, 0)

	// Process assets accounts
	f := from.FindFirst("Assets")
	if f == NotFound {
		return nil, fmt.Errorf("no assets found")
	}
	l := from.FindFirst("Total for Assets")
	if l == NotFound {
		return nil, fmt.Errorf("no assets total found")
	}

	log.Printf("f: %v; l:%v", f, l)

	i := f
	for {
		// Find the first account starting from 'i'
		p := from.FindFromTo(i, l, "Previous balance")
		log.Printf("p: %v", p)
		if p == NotFound {
			log.Println("no accounts?")
			break
		}
		aName := from.Cell(p-1, AccountNameCol) // One previous is the account name. One more is the beginning.
		tfor := fmt.Sprintf("Total for %v", aName)

		// Beginning balance as string.
		bBalStr := from.Cell(p, BalanceCol)

		r := from.FindFromTo(p, l, tfor) // End account
		if r == NotFound {
			return nil, fmt.Errorf("not found: %q", tfor)
		}

		var ac Account

		eBalStr := from.Cell(r, BalanceCol)
		acEB, err := USDToBFloat(eBalStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse: %q", eBalStr)
		}
		ac.EndBalance = *acEB

		ac.Name = aName
		ac.Type = AssetType
		acBB, err := USDToBFloat(bBalStr)
		ac.BeginBalance = *acBB
		if err != nil {
			return nil, fmt.Errorf("could not parse: %q", bBalStr)
		}
		err = from.ForEach(p+1, r, func(l []string) error {
			t, err := ParseTransaction(l)
			if err != nil {
				return fmt.Errorf("could not parse: %v", l)
			}
			ac.Transactions = append(ac.Transactions, t)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("while parsing: %v: %w", aName, err)
		}

		ret.Accounts = append(ret.Accounts, ac)

		i = r
	}

	return &ret, nil
}
