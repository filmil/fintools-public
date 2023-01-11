package xml

import (
	"fmt"
	"math"
	"testing"
)

func TestMonotonic(t *testing.T) {
	tests := []struct {
		x, y, z  float64
		expected bool
	}{
		{1, 1, 1, true},
		{1, 1, 2, true},
		{1, 2, 2, true},
		{1, 2, 3, true},
		{1, 2, math.Inf(-1), false},
		{1, 2, math.Inf(1), true},
		{1, 4, 3, false},
		{3, 2, 4, false},
		{math.Inf(-1), 1, math.Inf(0), true},
	}
	for _, test := range tests {
		test := test
		actual := Monotonic(test.x, test.y, test.z)
		if actual != test.expected {
			t.Errorf("Monotonic(%v,%v,%v)=%v, want=%v", test.x, test.y, test.z, actual, test.expected)
		}
	}
}

func TestIntersectingBBox(t *testing.T) {
	tests := []struct {
		b1, b2   BBox
		expected bool
	}{
		{
			b1:       BBox{Left: 0, Right: 2, Bottom: 0, Top: 2},
			b2:       BBox{Left: 1, Right: 3, Bottom: 1, Top: 3},
			expected: true,
		},
		{
			b1:       BBox{Left: 0, Right: 1, Bottom: 0, Top: 1},
			b2:       BBox{Left: 2, Right: 3, Bottom: 2, Top: 3},
			expected: false,
		},
		{
			b1:       BBox{Left: 1, Right: 2, Bottom: 1, Top: 4},
			b2:       BBox{Left: 0, Right: 3, Bottom: 2, Top: 3},
			expected: true,
		},
		{
			b1: BBox{
				Left:   445,
				Right:  474,
				Top:    2,
				Bottom: 1,
			},
			b2: BBox{
				Left:   475,
				Right:  math.Inf(1),
				Top:    2,
				Bottom: 1,
			},
			expected: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%v;%v", test.b1, test.b2),
			func(t *testing.T) {
				t.Parallel()
				actual := IntersectingBBox(test.b1, test.b2)
				if actual != test.expected {
					t.Errorf("IntersectingBBox(\n\t%+v,\n\t%+v)=%v, expected: %v", test.b1, test.b2, actual, test.expected)
				}
			})
	}
}
