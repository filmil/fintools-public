package csv2

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Must[T any](t T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("Must failed with: %v", err))
	}
	return t
}

func TestReport(t *testing.T) {

	tests := []struct {
		name     string
		data     [][]string
		expected Report
	}{
		{
			name: "first parse",
			data: [][]string{
				{},
				{},
				{"ReportName"},
				{"Assets"},
				{"AccountName"},
				{"Previous balance", "", "", "", "", "", "", "$1.00"},
				{"2/8/2023", "Check 12198", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"Total for AccountName", "", "", "", "", "", "", "$2.00"},
				{"Total for Assets"},
				{},
			},
			expected: Report{
				Name:    "ReportName",
				MinDate: Must(time.Parse(DateLayout, "2/8/2023")),
				Balance: *Must(USDToBFloat("$1")),
				Accounts: []Account{
					{
						Name:         "AccountName",
						Type:         AssetType,
						BeginBalance: *Must(USDToBFloat(("$1.00"))),
						EndBalance:   *Must(USDToBFloat(("$2.00"))),
						Transactions: []Transaction{
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2023")),
								Type:        "Check 12198",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
						},
					},
				},
			},
		},
		{
			name: "second parse",
			data: [][]string{
				{},
				{},
				{"ReportName"},
				{"Assets"},
				{"AccountName"},
				{"Previous balance", "", "", "", "", "", "", "$1.00"},
				{"2/8/2023", "Check 12198", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"2/8/2024", "Check 12199", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"Total for AccountName", "", "", "", "", "$3.00", "$4.00", "$2.00"},
				{"Total for Assets"},
				{},
			},
			expected: Report{
				Name:    "ReportName",
				MinDate: Must(time.Parse(DateLayout, "2/8/2023")),
				Balance: *Must(USDToBFloat("$1")),
				Accounts: []Account{
					{
						Name:         "AccountName",
						Type:         AssetType,
						BeginBalance: *Must(USDToBFloat(("$1.00"))),
						EndBalance:   *Must(USDToBFloat(("$2.00"))),
						TotalCredit:  *Must(USDToBFloat(("$3.00"))),
						TotalDebit:   *Must(USDToBFloat(("$4.00"))),
						Transactions: []Transaction{
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2023")),
								Type:        "Check 12198",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2024")),
								Type:        "Check 12199",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
						},
					},
				},
			},
		},
		{
			name: "multiple accounts",
			data: [][]string{
				{},
				{},
				{"ReportName"},
				{"Assets"},
				{"AccountName"},
				{"Previous balance", "", "", "", "", "", "", "$1.00"},
				{"2/8/2023", "Check 12198", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"2/8/2024", "Check 12199", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"Total for AccountName", "", "", "", "", "", "", "$2.00"},
				{"AccountName2"},
				{"Previous balance", "", "", "", "", "", "", "$1.00"},
				{"2/8/2023", "Check 12198", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"2/8/2024", "Check 12199", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"Total for AccountName2", "", "", "", "", "", "", "$2.00"},
				{"Total for Assets"},
				{},
			},
			expected: Report{
				Name:    "ReportName",
				MinDate: Must(time.Parse(DateLayout, "2/8/2023")),
				Balance: *Must(USDToBFloat("$1")),
				Accounts: []Account{
					{
						Name:         "AccountName",
						Type:         AssetType,
						BeginBalance: *Must(USDToBFloat(("$1.00"))),
						EndBalance:   *Must(USDToBFloat(("$2.00"))),
						Transactions: []Transaction{
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2023")),
								Type:        "Check 12198",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2024")),
								Type:        "Check 12199",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
						},
					},
					{
						Name:         "AccountName2",
						Type:         AssetType,
						BeginBalance: *Must(USDToBFloat(("$1.00"))),
						EndBalance:   *Must(USDToBFloat(("$2.00"))),
						Transactions: []Transaction{
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2023")),
								Type:        "Check 12198",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2024")),
								Type:        "Check 12199",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
						},
					},
				},
			},
		},
		{
			name: "multiple accounts multiple asset classes",
			data: [][]string{
				{},
				{},
				{"ReportName"},
				{"Assets"},
				{"AccountName"},
				{"Previous balance", "", "", "", "", "", "", "$1.00"},
				{"2/8/2023", "Check 12198", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"2/8/2024", "Check 12199", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"Total for AccountName", "", "", "", "", "", "", "$2.00"},
				{"AccountName2"},
				{"Previous balance", "", "", "", "", "", "", "$1.00"},
				{"2/8/2023", "Check 12198", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"2/8/2024", "Check 12199", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"Total for AccountName2", "", "", "", "", "", "", "$2.00"},
				{"Total for Assets"},
				{"Liabilities"},
				{"AccountName3"},
				{"Previous balance", "", "", "", "", "", "", "$1.00"},
				{"2/8/2023", "Check 12198", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"2/8/2024", "Check 12199", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"Total for AccountName3", "", "", "", "", "", "", "$2.00"},
				{"AccountName4"},
				{"Previous balance", "", "", "", "", "", "", "$1.00"},
				{"2/8/2023", "Check 12198", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"2/8/2024", "Check 12199", "U", "SD1", "9009007300-006_12-01_12-31-22", "$57.40", "", "$57.40"},
				{"Total for AccountName4", "", "", "", "", "", "", "$2.00"},
				{"Total for Liabilities"},
				{},
			},
			expected: Report{
				Name:    "ReportName",
				MinDate: Must(time.Parse(DateLayout, "2/8/2023")),
				Balance: *Must(USDToBFloat("$1")),
				Accounts: []Account{
					{
						Name:         "AccountName",
						Type:         AssetType,
						BeginBalance: *Must(USDToBFloat(("$1.00"))),
						EndBalance:   *Must(USDToBFloat(("$2.00"))),
						Transactions: []Transaction{
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2023")),
								Type:        "Check 12198",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2024")),
								Type:        "Check 12199",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
						},
					},
					{
						Name:         "AccountName2",
						Type:         AssetType,
						BeginBalance: *Must(USDToBFloat(("$1.00"))),
						EndBalance:   *Must(USDToBFloat(("$2.00"))),
						Transactions: []Transaction{
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2023")),
								Type:        "Check 12198",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2024")),
								Type:        "Check 12199",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
						},
					},
					{
						Name:         "AccountName3",
						Type:         LiabilityType,
						BeginBalance: *Must(USDToBFloat(("$1.00"))),
						EndBalance:   *Must(USDToBFloat(("$2.00"))),
						Transactions: []Transaction{
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2023")),
								Type:        "Check 12198",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2024")),
								Type:        "Check 12199",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
						},
					},
					{
						Name:         "AccountName4",
						Type:         LiabilityType,
						BeginBalance: *Must(USDToBFloat(("$1.00"))),
						EndBalance:   *Must(USDToBFloat(("$2.00"))),
						Transactions: []Transaction{
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2023")),
								Type:        "Check 12198",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
							{
								Date:        Must(time.Parse(DateLayout, "2/8/2024")),
								Type:        "Check 12199",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			d := CSVData{data: test.data}

			r, err := NewReport(&d)
			if err != nil {
				t.Fatalf("didn't parse: %v", err)
			}

			if !cmp.Equal(*r, test.expected, cmp.Comparer(cmpFloat)) {
				t.Errorf("want:\n\t%+v\ngot:\n\t%+v\n\tdiff:\n\t%v",
					test.expected, *r, cmp.Diff(test.expected, *r, cmp.Comparer(cmpFloat)))
			}

		})
	}
}

var cmpFloat = func(a, b big.Float) bool {
	return a.Cmp(&b) == 0
}

func TestBigFloat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected big.Float
	}{
		{"$1.0", *Must(USDToBFloat("1"))},
		{"$1.000", *Must(USDToBFloat("1"))},
		{"$1000", *Must(USDToBFloat("1000"))},
		{"$1,000", *Must(USDToBFloat("1000"))},
		// Not intended, but fine, I suppose.
		{"$1,000$", *Must(USDToBFloat("1000"))},
	}
	for _, test := range tests {
		test := test
		t.Run("", func(t *testing.T) {
			actual := *Must(USDToBFloat(test.input))
			if actual.Cmp(&test.expected) != 0 {
				t.Errorf("want: %v, got: %v", test.expected, actual)
			}
		})
	}
}
