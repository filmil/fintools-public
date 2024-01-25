package tiller

import (
	"fmt"
	"testing"
	"time"

	"github.com/filmil/fintools-public/pkg/csv2"
)

func TestFirstOf(t *testing.T) {
	t.Parallel()
	tests := []struct {
		date     string
		expected string
		fn       func(time.Time) time.Time
	}{
		{"1/1/2002", "1/1/2002", FirstOfMonth},
		{"10/10/2002", "10/1/2002", FirstOfMonth},
		{"10/11/2002", "10/1/2002", FirstOfMonth},
		{"10/12/2002", "10/1/2002", FirstOfMonth},
		{"10/13/2002", "10/1/2002", FirstOfMonth},
		{"10/14/2002", "10/1/2002", FirstOfMonth},

		{"1/21/2024", "1/21/2024", FirstOfWeek}, // Sunday
		{"1/22/2024", "1/21/2024", FirstOfWeek},
		{"1/23/2024", "1/21/2024", FirstOfWeek},
		{"1/24/2024", "1/21/2024", FirstOfWeek},
		{"1/25/2024", "1/21/2024", FirstOfWeek},
		{"1/26/2024", "1/21/2024", FirstOfWeek},
		{"1/27/2024", "1/21/2024", FirstOfWeek},
		{"1/28/2024", "1/28/2024", FirstOfWeek},
	}

	for i, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			d, err := time.Parse(csv2.DateLayout, test.date)
			if err != nil {
				t.Fatalf("could not parse: %v: %v", test.date, err)
			}
			actual := test.fn(d).Format(csv2.DateLayout)
			if actual != test.expected {
				t.Errorf("want: %v, got: %v", test.expected, actual)
			}
		})
	}
}
