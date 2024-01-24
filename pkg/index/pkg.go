package index

import (
	"math/big"
	"time"

	"github.com/filmil/fintools/pkg/csv2"
)

type Entry struct {
	AccountName string
	Ty          csv2.AccountType
	Tr          *csv2.Transaction
	//Acc *csv2.Account
}

func FindFirstByType(t csv2.AccountType, es []Entry) (*csv2.Transaction, string) {
	for _, e := range es {
		if e.Ty == t {
			return e.Tr, e.AccountName
		}
	}
	return nil, ""
}

func FindOthers(t csv2.AccountType, es []Entry) []*csv2.Transaction {
	var ret []*csv2.Transaction
	for _, e := range es {
		if e.Ty != t {
			ret = append(ret, e.Tr)
		}
	}
	return ret
}

// ID returns a (non-unique) identifier.
func (e Entry) ID() string {
	if e.Tr == nil {
		panic("nil transaction in ID()")
	}
	return e.Tr.ID()
}

// Instance is an index of all the transactions
type Instance struct {
	Name string
	// The earliest date seen in the reports.
	MinDate  time.Time
	Balance  big.Float
	entries  map[string][]Entry
	accounts map[string]*csv2.Account
}

func (i Instance) StableKeys() []string {
	return csv2.StableKeys(i.entries)
}

func (i Instance) GetEntries(k string) []Entry {
	return i.entries[k]
}

func (i Instance) GetAccountByName(n string) *csv2.Account {
	return i.accounts[n]
}

// Creates a new index instance based on the Report.
func New(r *csv2.Report) *Instance {
	var in Instance
	in.entries = map[string][]Entry{}
	in.accounts = map[string]*csv2.Account{}
	in.Name = r.Name
	in.MinDate = r.MinDate
	in.Balance = r.Balance

	for i := range r.Accounts {
		rai := &r.Accounts[i]
		in.accounts[rai.Name] = rai
		for j, _ := range rai.Transactions {
			t := &rai.Transactions[j]
			var e Entry
			e.Tr = t
			e.AccountName = rai.Name
			e.Ty = rai.Type

			eid := e.ID()
			ee := in.entries[eid]
			ee = append(ee, e)
			in.entries[eid] = ee
		}
	}

	return &in
}
