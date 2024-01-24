package csv2

import (
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
