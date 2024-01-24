package csv2

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"
)

// USDToBFloat parses a big float out of a string such as "$42.00"
func USDToBFloat(s string) (*big.Float, error) {
	if s == "" {
		s = "$0.00"
	}
	// Skip `$` and `,` if present.
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")

	bf, _, err := big.ParseFloat(s, 10, 100, big.ToNearestEven)
	if err != nil {
		return nil, fmt.Errorf("could not parse: %q: %w", s, err)
	}
	return bf, nil
}

type Transaction struct {
	Date time.Time

	Type, Unit, Name, Description string
	Debit, Credit, Balance        big.Float
}

func (t Transaction) String() string {
	return strings.Join(
		[]string{
			t.Date.Format(DateLayout),
			t.Type, t.Unit, t.Name, t.Description,
			fmt.Sprintf("$%v", t.Debit.Text('f', 6)),
			fmt.Sprintf("$%v", t.Credit.Text('f', 6)),
			fmt.Sprintf("$%v", t.Balance.Text('f', 6)),
		},
		" / ")
}

// ID returns the transaction ID. This is shared by all transactions that have
// the same characteristics.
func (t Transaction) ID() string {
	h := sha256.New()

	h.Write([]byte(t.Type))
	h.Write([]byte(t.Date.Format(DateLayout)))
	h.Write([]byte(t.Unit))
	h.Write([]byte(t.Name))
	h.Write([]byte(t.Description))

	return fmt.Sprintf("%x", h.Sum(nil))
}

const DateLayout = `1/2/2006`

// ParseTransaction extracts one transaction or errors out.
func ParseTransaction(l []string) (Transaction, error) {
	var t Transaction
	td, err := time.Parse(DateLayout, l[0])
	if err != nil {
		return t, fmt.Errorf("could not parse date: %v: %w", l[0], err)
	}
	t.Date = td
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

// COrd is a comparable interface.
// Not sure why this does not work out of the box.
type COrd interface {
	~int | ~string
}

// StableKeys returns the map keys of `m` in a stable ordering.
func StableKeys[K COrd, V any](m map[K]V) []K {
	var ret []K
	for k := range m {
		ret = append(ret, k)
	}
	// Sort the keys, since by default iteration over a map has a random order.
	sort.SliceStable(ret, func(i, j int) bool {
		return ret[i] < ret[j]

	})
	return ret
}

type AccountType int

var AccountTypeToStr = map[AccountType]string{
	AssetType:     "Assets",
	LiabilityType: "Liabilities",
	EquityType:    "Equity",
	IncomeType:    "Income",
	ExpensesType:  "Expenses",
}

func (a AccountType) String() string {
	return AccountTypeToStr[a]
}

const (
	// AssetType is an account type of an asset.
	AssetType AccountType = iota
	// LiabilityType is an account type of a liability.
	LiabilityType
	EquityType
	IncomeType
	ExpensesType
)

// A description of a single account.
type Account struct {
	Name                     string
	Type                     AccountType
	BeginBalance, EndBalance big.Float
	TotalCredit, TotalDebit  big.Float
	Transactions             []Transaction
	MinDate                  time.Time
}

// Report is a structured report from an account.
type Report struct {
	Name     string
	MinDate  time.Time
	Balance  big.Float
	Accounts []Account
}

func NewReport(from *CSVData) (*Report, error) {
	const (
		AccountNameCol = 0
		TotalCreditCol = 5
		TotalDebitCol  = 6
		BalanceCol     = 7
		AccountNameRow = 2
	)

	var ret Report
	ret.Name = from.Cell(AccountNameRow, AccountNameCol)
	ret.MinDate, _ = time.Parse("01/02/2006", "12/31/4001")

	// Run for all account types.
	for _, act := range StableKeys(AccountTypeToStr) {
		acts := AccountTypeToStr[act]
		f := from.FindFirst(acts)
		if f == NotFound {
			continue
		}
		l := from.FindFirst(fmt.Sprintf("Total for %s", acts))
		if l == NotFound {
			continue
		}

		i := f
		for {
			// Find the first account starting from 'i'
			p := from.FindFromTo(i, l, "Previous balance")
			if p == NotFound {
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

			// Fill out a new account.
			var ac Account
			ac.Name = aName
			ac.Type = act

			eBalStr := from.Cell(r, BalanceCol)
			acEB, err := USDToBFloat(eBalStr)
			if err != nil {
				return nil, fmt.Errorf("could not parse: %q: %w", eBalStr, err)
			}
			ac.EndBalance = *acEB

			totalCreditStr := from.Cell(r, TotalCreditCol)
			totalCredit, err := USDToBFloat(totalCreditStr)
			if err != nil {
				return nil, fmt.Errorf("could not parse total credit: %q: %w", totalCreditStr, err)
			}
			ac.TotalCredit = *totalCredit

			totalDebitStr := from.Cell(r, TotalDebitCol)
			totalDebit, err := USDToBFloat(totalDebitStr)
			if err != nil {
				return nil, fmt.Errorf("could not parse total debit: %q: %w", totalCreditStr, err)
			}
			ac.TotalDebit = *totalDebit

			acBB, err := USDToBFloat(bBalStr)
			ac.BeginBalance = *acBB
			if act == AssetType {
				ret.Balance = ac.BeginBalance
			}
			if err != nil {
				return nil, fmt.Errorf("could not parse: %q", bBalStr)
			}
			err = from.ForEach(p+1, r, func(l []string) error {
				t, err := ParseTransaction(l)
				if err != nil {
					return fmt.Errorf("could not parse: %v", l)
				}
				// Update the minimum date observed. But only for the asset type.
				if act == AssetType && t.Date.Before(ret.MinDate) {
					ret.MinDate = t.Date
				}
				ac.Transactions = append(ac.Transactions, t)
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("while parsing: %v: %w", aName, err)
			}

			ret.Accounts = append(ret.Accounts, ac)

			// Advance the iteration here.
			i = r
		}
	}

	// Process assets accounts

	return &ret, nil
}
