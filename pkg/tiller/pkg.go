// Package tiller provides tiller output formatting.
package tiller

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"sort"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/filmil/fintools/pkg/cfg"
	"github.com/filmil/fintools/pkg/csv2"
	"github.com/filmil/fintools/pkg/index"
	"github.com/google/uuid"
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

// GenId generates a stable ID from a given string (account name really)
func GenId(n string) string {
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
	time.Time
	// A transaction UUID. Each transaction leg has the same UUID.
	TransactionID string
	AccountID,
	CheckNumber string
	// This should potentially be a concatenation of all text about the transaction. Potentially.
	FullDescription string
	// These should likely stay unset.
	DateAdded, CategorizedDate string
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

		acId := GenId(acName)
		if _, ok := seen[acName]; !ok {
			var b Balance
			b.Date = DateTime(ac.MinDate)
			b.AccountName = acName
			b.AccountID = acId
			b.Balance = ac.BeginBalance
			ret.Balances = append(ret.Balances, b)
			seen[acName] = struct{}{}

			// TODO: filmil - Also needs the end balance I think.

		}
		if ac == nil {
			// This shouldn't happen, but I wonder how long before it does.
			panic("account is nil")
		}
		ret.AccountID = ac.Name
		es := index.FindOthers(csv2.AssetType, entries)
		log.Printf("Line item -->\nasset:\n\t%+v\ntxs:\n\t%+v", asset, spew.Sdump(es))

		txid := uuid.New()

		for _, e := range es {
			var r Row

			r.Date = DateTime(e.Tx.Date)
			r.Description = e.Tx.Description
			r.Category = c.GetCat(e.AccName)
			r.Amount = CalculateAmount(e.Ty, e.Tx.Debit, e.Tx.Credit)
			// TODO: filmil - fill out amount
			r.Tags = "Tax"
			r.AccountDesc = e.AccName
			r.AccountNum = ""   // Dunno.
			r.Institution = ""  // Dunno.
			r.Month = e.Tx.Date // The first of the month of this date
			r.Week = e.Tx.Date
			r.TransactionID = txid.String()
			r.AccountID = GenId(acName)
			r.FullDescription = e.Tx.Description // Dunno.
			// TODO: filmil - Dunno the rest, see what can be done.

			ret.Rows = append(ret.Rows, r)
		}

		// TODO: filmil - Process transactions.
		// TODO: filmil - Find the account end balance.

	}

	// Sort rows.
	SortDater(ret.Rows)
	SortDater(ret.Balances)
	return &ret
}
