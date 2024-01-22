package csv2

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFindFirst(t *testing.T) {

	tests := []struct {
		name string
		data [][]string

		find string

		expected int
	}{
		{
			name: "notfound",
			data: [][]string{
				{},
				{"foo", "", "baz"},
			},

			find:     "baz",
			expected: NotFound,
		},
		{
			name: "first",
			data: [][]string{
				{},
				{"foo", "", "baz"},
			},

			find:     "foo",
			expected: 1,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			d := CSVData{data: test.data}

			a := d.FindFirst(test.find)
			if a != test.expected {
				t.Errorf("want: %v; got: %v", test.expected, a)
			}
		})
	}
}

func TestFindFrom(t *testing.T) {

	tests := []struct {
		name string
		data [][]string

		find string
		from int

		expected int
	}{
		{
			name: "notfound",
			data: [][]string{
				{},
				{"foo", "", "baz"},
				{},
				{},
				{"foo", "", "baz"},
			},

			find:     "foo",
			from:     1,
			expected: 1,
		},
		{
			name: "notfound",
			data: [][]string{
				{},
				{"foo", "", "baz"},
				{},
				{},
				{"foo", "", "baz"},
			},

			find:     "foo",
			from:     2,
			expected: 4,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			d := CSVData{data: test.data}

			a := d.FindFrom(test.from, test.find)
			if a != test.expected {
				t.Errorf("want: %v; got: %v", test.expected, a)
			}
		})
	}
}

func TestFindFromTo(t *testing.T) {

	tests := []struct {
		name string
		data [][]string

		find     string
		from, to int

		expected int
	}{
		{
			name: "found1",
			data: [][]string{
				{},
				{"foo", "", "baz"},
				{},
				{},
				{"foo", "", "baz"},
			},

			find:     "foo",
			from:     0,
			to:       2,
			expected: 1,
		},
		{
			name: "found4",
			data: [][]string{
				{},
				{"foo", "", "baz"},
				{},
				{},
				{"foo", "", "baz"},
			},

			find:     "foo",
			from:     2,
			to:       5,
			expected: 4,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			d := CSVData{data: test.data}

			a := d.FindFromTo(test.from, test.to, test.find)
			if a != test.expected {
				t.Errorf("want: %v; got: %v", test.expected, a)
			}
		})
	}
}

func TestCell(t *testing.T) {

	tests := []struct {
		name     string
		data     [][]string
		r, c     int
		expected string
	}{
		{
			name: "cell1,0",
			data: [][]string{
				{},
				{"foo", "", "baz"},
				{},
				{},
				{"foo", "", "baz"},
			},
			r: 1, c: 0,
			expected: "foo",
		},
		{
			name: "cell1,1",
			data: [][]string{
				{},
				{"foo", "", "baz"},
				{},
				{},
				{"foo", "", "baz"},
			},
			r: 1, c: 1,
			expected: "",
		},
		{
			name: "way out",
			data: [][]string{
				{},
				{"foo", "", "baz"},
				{},
				{},
				{"foo", "", "baz"},
			},
			r: 1, c: 1000,
			expected: "",
		},
		{
			name: "way out r",
			data: [][]string{
				{},
				{"foo", "", "baz"},
				{},
				{},
				{"foo", "", "baz"},
			},
			r: 1000, c: 0,
			expected: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			d := CSVData{data: test.data}

			a := d.Cell(test.r, test.c)
			if a != test.expected {
				t.Errorf("want: %v; got: %v", test.expected, a)
			}
		})
	}
}

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
				Name: "ReportName",
				Accounts: []Account{
					{
						Name:         "AccountName",
						Type:         AssetType,
						BeginBalance: *Must(USDToBFloat(("$1.00"))),
						EndBalance:   *Must(USDToBFloat(("$2.00"))),
						Transactions: []Transaction{
							{
								Date:        "2/8/2023",
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
				{"Total for AccountName", "", "", "", "", "", "", "$2.00"},
				{"Total for Assets"},
				{},
			},
			expected: Report{
				Name: "ReportName",
				Accounts: []Account{
					{
						Name:         "AccountName",
						Type:         AssetType,
						BeginBalance: *Must(USDToBFloat(("$1.00"))),
						EndBalance:   *Must(USDToBFloat(("$2.00"))),
						Transactions: []Transaction{
							{
								Date:        "2/8/2023",
								Type:        "Check 12198",
								Unit:        "U",
								Name:        "SD1",
								Description: "9009007300-006_12-01_12-31-22",
								Debit:       *Must(USDToBFloat("$57.40")),
								Credit:      *Must(USDToBFloat("$0.00")),
								Balance:     *Must(USDToBFloat("$57.40")),
							},
							{
								Date:        "2/8/2024",
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

			if !reflect.DeepEqual(fmt.Sprintf("%+v", *r), fmt.Sprintf("%+v", test.expected)) {
				t.Errorf("want:\n\t%+v\ngot:\n\t%+v", test.expected, *r)
			}

		})
	}
}
