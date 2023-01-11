package xml

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/filmil/fintools/pkg/tx"
	"github.com/google/go-cmp/cmp"
)

type TestSpec struct {
	// Input is the filename with paystub data to read
	Input    string         `json:",omitempty"`
	Expected tx.Transaction `json:",omitempty"`
}

type TestSpecSeq []TestSpec

func ReadJSON(r io.Reader) (TestSpecSeq, error) {
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	var ret TestSpecSeq
	if err := d.Decode(&ret); err != nil {
		return ret, fmt.Errorf("could not parse JSON: %w", err)
	}
	return ret, nil
}

func TestParsingJSON(t *testing.T) {
	const filename = "testdata_private/test_spec.json"
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("could not open, see README.md in this dir: %v: %v", filename, err)
	}
	tests, err := ReadJSON(f)
	if err != nil {
		t.Fatalf("could not parse JSON: %v: %v", filename, err)
	}

	fmt.Printf("specs: %+v", tests)

	opts := cmp.Options{}
	for _, test := range tests {
		test := test
		t.Run(test.Input, func(t *testing.T) {
			t.Parallel()
			defer func() {
				if r := recover(); r != nil {
					panic(fmt.Sprintf("Panic in test: %v", test.Input))
				}
			}()
			f, err := os.Open(test.Input)
			if err != nil {
				t.Fatalf("Open: unexpected error: %v", err)
			}
			defer f.Close()
			tr, err := Parse(f)
			if err != nil {
				t.Fatalf("Run: unexpected error: %v", err)
			}
			if !cmp.Equal(test.Expected, tr, opts) {
				t.Errorf("Run(_)=%+v\nwant:\n%+v\ndiff:\n%+v",
					tr, test.Expected, cmp.Diff(test.Expected, tr, opts))
			}
		})
	}
}
