package csv2

import (
	"encoding/csv"
	"fmt"
	"io"
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
		ret = append(ret, fmt.Sprintf("%010d: %v", n, strings.Join(l, " , ")))
	}
	return strings.Join(ret, "\n")
}

func (d CSVData) Size() int {
	return len(d.data)
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
