// Package tiller provides tiller output formatting.
package tiller

import (
	"fmt"
	"log"
	"math/big"
	"sort"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/filmil/fintools/pkg/cfg"
	"github.com/filmil/fintools/pkg/csv2"
	"github.com/filmil/fintools/pkg/index"
)

type Dater interface {
	DateM() time.Time
}

func SortDater[D Dater](d []D) {
	sort.Slice(d, func(i, j int) bool {
		return d[i].DateM().Before(d[j].DateM())
	})
}

type DateTime time.Time

func (d DateTime) String() string {
	return fmt.Sprintf("%v", time.Time(d).Format("1/2/2006"))
}

// Row represents one row of a tiller structure.
type Row struct {
	Date DateTime
	Empty,
	Description,
	Category,
	Amount,
	Tags,
	Account,
	AccountNum,
	Institution,
	Month,
	Week,
	TransactionID,
	AccountID,
	CheckNumber,
	FullDescription,
	DateAdded,
	CategorizedDate string
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

		account := i.GetAccountByName(acName)

		if _, ok := seen[acName]; !ok {
			var b Balance
			b.Date = DateTime(account.MinDate)
			b.AccountName = acName
			b.AccountID = c.GetID()
			b.Balance = account.BeginBalance
			ret.Balances = append(ret.Balances, b)
			seen[acName] = struct{}{}

		}
		if account == nil {
			// This shouldn't happen, but I wonder how long before it does.
			panic("account is nil")
		}
		ret.AccountID = account.Name
		txs := index.FindOthers(csv2.AssetType, entries)
		log.Printf("Line item -->\nasset:\n\t%+v\ntxs:\n\t%+v", asset, spew.Sdump(txs))
		// TODO: fmil - Find the account end balance.

	}
	// Sort rows.
	SortDater(ret.Rows)
	SortDater(ret.Balances)
	return &ret
}
