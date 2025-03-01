// Package tiller provides tiller output formatting.
package tiller

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"io"
	"math/big"
	"sort"
	"time"

	"github.com/filmil/fintools-public/pkg/cfg"
	"github.com/filmil/fintools-public/pkg/csv2"
	"github.com/filmil/fintools-public/pkg/index"
	"github.com/google/uuid"
)

const (
	IndexDate        = 1
	IndexDescription = 2
	IndexCategory    = 3
	IndexAmount      = 4
	IndexTags        = 5
	IndexAccountName = 6
	IndexAccountID   = 12
)

var (
	ColumnNames []string = []string{
		"T",
		"Date",
		"Description",
		"Category",
		"Amount", // 5
		"Tags",
		"Account", // Account descriptive name
		"Account #",
		"Institution",
		"Month", // 10
		"Week",
		"Transaction ID",
		"Account ID",
		"Check Number",
		"Full Description", // 15
		"Date Added",
		"Categorized Date",
	}
)

type Dater interface {
	DateM() time.Time
}

func SortDater[D Dater](d []D) {
	sort.Slice(d, func(i, j int) bool {
		return d[i].DateM().Before(d[j].DateM())
	})
}

var _ fmt.Stringer = (*DateTime)(nil)

// Only so that we can implement a special String.
type DateTime time.Time

func (d DateTime) String() string {
	return fmt.Sprintf("%v", time.Time(d).Format("1/2/2006"))
}

// GenID generates a stable ID from a given string (account name really)
func GenID(n string) string {
	h := sha256.New()
	h.Write([]byte(n))
	return fmt.Sprintf("manual:%x", h.Sum(nil))
}

// Row represents one row of a tiller structure.
type Row struct {
	Date DateTime
	Empty,
	Description,
	Category string
	Amount big.Float
	Tags   string
	// An account's textual description
	AccountDesc string
	AccountNum  string
	Institution string
	// The first of the month of Date.
	Month time.Time
	// The first of the week of Date.
	Week time.Time
	// A transaction UUID. Each transaction leg has the same UUID.
	TransactionID string
	AccountID,
	CheckNumber string
	// This should potentially be a concatenation of all text about the transaction. Potentially.
	FullDescription string
	// These should likely stay unset.
	DateAdded, CategorizedDate string
}

func (row Row) AsCSVRow() []string {
	t := time.Time(row.Date)
	f := DateFmt(t)
	return []string{
		"",
		f,                        // "Date",
		row.Description,          // "Description",
		row.Category,             // "Category",
		USD(&row.Amount),         // "Amount",        // 5
		row.Tags,                 // "Tags",
		row.AccountDesc,          // "Account", // Account descriptive name
		row.AccountNum,           // "Account #",
		"",                       // "Institution",
		DateFmt(FirstOfMonth(t)), // "Month",          // 10
		DateFmt(FirstOfWeek(t)),  // "Week",
		row.TransactionID,        //"Transaction ID",
		row.AccountID,            // "Account ID",
		"",                       // "Check Number",
		row.FullDescription,      //  "Full Description", // 15
		f,                        // "Date Added",
		"",                       // Categorized Date",
	}
}

func (r Row) DateM() time.Time {
	return time.Time(r.Date)
}

type CSVDate time.Time

type Balance struct {
	Date        DateTime
	AccountName string
	AccountID   string
	Balance     big.Float
}

func (b Balance) DateM() time.Time {
	return time.Time(b.Date)
}

type Export struct {
	AccountID string
	Rows      []Row
	Balances  []Balance
}

// Caluclates the amount added or subtracted from an account type, given a debit
// and credit value.
func CalculateAmount(ty csv2.AccountType, d, c big.Float) big.Float {
	var ret big.Float

	// Take good note which accounts do what.
	//
	// For Assets: Debit means adding to the account.
	//             Credit means subtracting from the account.
	// For all other accounts, it's the reverse.
	if ty == csv2.AssetType {
		ret.Add(&ret, &d)
		ret.Sub(&ret, &c)
	} else {
		ret.Add(&ret, &c)
		ret.Sub(&ret, &d)
	}

	return ret
}

func New(i *index.Instance, c *cfg.Instance) *Export {
	var (
		ret Export
	)

	seen := map[string]struct{}{}

	for _, a := range i.StableKeys() {
		entries := i.GetEntries(a)

		asset, acName := index.FindFirstByType(csv2.AssetType, entries)
		if asset == nil {
			// Account like "Prepayment" doesn't have any pair.
			// Skip it, because I don't know how to express it in single
			// entry bookkeeping.
			continue
		}

		ac := i.GetAccountByName(acName)

		// Reuse or generate the account ID
		acId := c.GetAccID(acName)
		if acId == "" {
			acId = GenID(acName)
		}

		if _, ok := seen[acName]; !ok {
			var b Balance
			b.Date = DateTime(ac.MinDate)
			b.AccountName = acName
			b.AccountID = acId
			b.Balance = ac.BeginBalance
			seen[acName] = struct{}{}

			ret.Balances = append(ret.Balances, b)

			// TODO: filmil - Also needs the end balance I think.

		}
		if ac == nil {
			// This shouldn't happen, but I wonder how long before it does.
			panic("account is nil")
		}
		ret.AccountID = ac.Name
		es := index.FindOthers(csv2.AssetType, entries)

		txid := uuid.New()

		for _, e := range es {
			var r Row

			r.Date = DateTime(e.Tx.Date)
			r.Description = e.Tx.Description
			r.Category = c.GetCat(e.AccName)
			r.Amount = CalculateAmount(e.Ty, e.Tx.Debit, e.Tx.Credit)
			r.Tags = "Tax"
			r.AccountDesc = e.AccName
			//r.AccountNum = ""  // Dunno.
			//r.Institution = "" // Dunno.
			r.Month = FirstOfMonth(e.Tx.Date)
			r.Week = FirstOfWeek(e.Tx.Date)
			r.TransactionID = txid.String()
			r.AccountID = acId
			r.FullDescription = e.Tx.FullDescription()

			ret.Rows = append(ret.Rows, r)
		}
		// TODO: filmil - Find the account end balance.
	}

	// Sort rows.
	SortDater(ret.Rows)
	SortDater(ret.Balances)
	return &ret
}

// Writes the transaction rows.
func (e *Export) WriteRows(w io.Writer) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()
	h := ColumnNames
	if err := cw.Write(h); err != nil {
		return fmt.Errorf("could not write header: %w", err)
	}
	for _, row := range e.Rows {
		r := row.AsCSVRow()
		if err := cw.Write(r); err != nil {
			return fmt.Errorf("could not write row: %+v: %w", row, err)
		}
	}
	return nil
}

// DateFmt formats dates in the way used in Tiller.
func DateFmt(t time.Time) string {
	return t.Format(csv2.DateLayout)
}

// USD formats a big float as a dollar amount with sufficient financial precision.
func USD(b *big.Float) string {
	return fmt.Sprintf("$%v", b.Text('f', 6))
}

func (e *Export) WriteBalances(w io.Writer) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()
	// Write the headers of the account balances
	h := []string{
		"T",
		"Date",
		"Time",
		"Account",
		"Account #", // 5
		"Account ID",
		"Balance ID",
		"Institution",
		"Balance",
		"Month", // 10
		"Week",
		"Type",
		"Class",
		"Account Status",
		"Date Added", // 15
	}
	if err := cw.Write(h); err != nil {
		return fmt.Errorf("could not write header: %w", err)
	}
	for _, bal := range e.Balances {
		t := time.Time(bal.Date)
		f := t.Format(csv2.DateLayout)
		r := []string{
			"",
			f, // Date
			f, // Time
			bal.AccountName,
			"", // 5: Account #
			bal.AccountID,
			uuid.New().String(), // Balance ID (unique ID)
			"",                  // "Institution",
			USD(&bal.Balance),   // Balance,
			f,                   // Month
			f,                   // Week
			"Other",             // Type
			"Asset",             // "Class",
			"ACTIVE",            // Account status
			f,                   // Date Added
		}
		if err := cw.Write(r); err != nil {
			return fmt.Errorf("could not write balance: %+v: %w", bal, err)
		}
	}
	return nil
}

// First of month returns the time corresponding to the 1st of the month of t.
func FirstOfMonth(t time.Time) time.Time {
	var ret time.Time
	return ret.AddDate(t.Year()-1, int(t.Month())-1, 0)
}

// FirstOfWeek returns the time corresponding to the Sunday of the week of t.
func FirstOfWeek(t time.Time) time.Time {
	var ret time.Time
	return ret.AddDate(t.Year()-1, int(t.Month())-1, t.Day()-int(t.Weekday()+1))
}
