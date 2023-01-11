package xml

import (
	"fmt"
	"sort"
	"strings"
)

// Textline represents a text that is laid out in a single physical line on a
// page.  A textline has a bounding box that contains all characters inside the
// line.
type Textline struct {
	BBox  BBox   `xml:"bbox,attr"`
	Texts []Text `xml:"text,omitempty"`
}

// SortLeft sorts textline in a nonincreasing order of left boundaries.
func SortLeft(tl []Textline) []Textline {
	sort.Slice(tl, func(i, j int) bool {
		return tl[i].BBox.Left <= tl[j].BBox.Left
	})
	return tl
}

// SortTop sorts textline in a decreasing order of top boundaries (top down).
func SortTop(tl []Textline) []Textline {
	sort.Slice(tl, func(i, j int) bool {
		return tl[i].BBox.Top >= tl[j].BBox.Top
	})
	return tl
}

// ForAllBBox runs the function do on every bounding box inside a textline.
func (t Textline) ForAllBBox(do func(b BBox)) {
	if t.BBox != NullBBox() {
		do(t.BBox)
	}
	for _, tx := range t.Texts {
		tx.ForAllBBox(do)
	}
}

// String implements Stringer.
func (t Textline) String() string {
	return fmt.Sprintf("<%q BBox:%+v>", t.Text(), t.BBox)
}

// Text returns the text contained in this text line.
func (t Textline) Text() string {
	var b strings.Builder
	for _, tx := range t.Texts {
		fmt.Fprintf(&b, "%s", tx.T)
	}
	return strings.Trim(b.String(), "\n")
}
