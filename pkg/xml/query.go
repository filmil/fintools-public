package xml

import (
	"fmt"
	"strings"
)

// Querying things on a Page.

// Match returns all textlines matching s.
func Match(p Page, s string) []Textline {
	var t []Textline
	p.ForAllTextlines(func(l Textline) {
		if l.Text() == s {
			t = append(t, l)
		}
	})
	return t
}

func Textlines(p Page) []Textline {
	var t []Textline
	p.ForAllTextlines(func(l Textline) {
		t = append(t, l)
	})
	return t
}

type Predicate func(l Textline) bool

// MatchPredicate retains textlines that match every Predicate that is passed
// in.
func MatchPredicate(
	tl []Textline, predicates ...Predicate) []Textline {
	var t []Textline
	f := func(l Textline) bool {
		for _, p := range predicates {
			if p(l) == false {
				return false
			}
		}
		return true
	}
	for _, l := range tl {
		if f(l) {
			t = append(t, l)
		}

	}
	return t
}

// TextOf returns the texts corresponding to the given textlines.
func TextOf(ts []Textline) []string {
	var s []string
	for _, t := range ts {
		s = append(s, t.Text())
	}
	return s
}

func OneTextline(ts []Textline) (Textline, error) {
	if len(ts) != 1 {
		return Textline{}, fmt.Errorf("Not a singleton textline: %+v", ts)
	}
	return ts[0], nil
}

// BindBBox returns a Predicate which requires that a textline intersects with
// the bounding box.
func BindBBox(b BBox, f func(BBox, Textline) bool) Predicate {
	return func(t Textline) bool { // Predicate
		return f(b, t)
	}
}

// IntersectingBBoxTextline returns true if t intersects box b.
func IntersectingBBoxTextline(b BBox, t Textline) bool {
	return IntersectingBBox(b, t.BBox)
}

// BindText returns a Predicate which matches a textline that contains the
// "text".
func BindText(text string, f func(string, Textline) bool) Predicate {
	return func(t Textline) bool {
		return f(text, t)
	}
}

// MatchingText returns true if textline matches the text exactly.
func MatchingText(text string, t Textline) bool {
	return text == t.Text()
}

// MatchingPrefix returns true if textline prefix matches the given text.
func MatchingPrefix(prefix string, t Textline) bool {
	return strings.HasPrefix(t.Text(), prefix)
}

// MatchingSuffix returns true if textline suffix matches the given text.
func MatchingSuffix(prefix string, t Textline) bool {
	return strings.HasSuffix(t.Text(), prefix)
}

/// Below functions build on the primitives to get the more commonly useful
/// functionality.

// FindOneTL finds the (known) single textline containing exactly text.
func FindOneTL(tls []Textline, text string) (Textline, error) {
	return OneTextline(MatchPredicate(tls, BindText(text, MatchingText)))
}

// FindOneTLPrefix finds the (known) single textline containing the prefix.
func FindOneTLPrefix(tls []Textline, prefix string) (Textline, error) {
	return OneTextline(MatchPredicate(tls, BindText(prefix, MatchingPrefix)))
}

// FindOneTLInBBox finds the (known) single textline containing exactly text,
// within the extents of the given bounding box.
func FindOneTLInBBox(tls []Textline, text string, bbox BBox) (Textline, error) {
	return OneTextline(
		MatchPredicate(tls,
			BindText(text, MatchingText),
			BindBBox(bbox, IntersectingBBoxTextline)))
}

// FindOneTLInBBoxWithSuffix finds the known single textline with text that
// has suffix matching the given suffix.
func FindOneTLInBBoxWithSuffix(tls []Textline, suffix string, bbox BBox) (Textline, error) {
	return OneTextline(
		MatchPredicate(tls,
			BindText(suffix, MatchingSuffix),
			BindBBox(bbox, IntersectingBBoxTextline)))
}

// FindInBBox filters textlines to only those that intersect with bbox.
func FindInBBox(tls []Textline, bbox BBox) []Textline {
	return MatchPredicate(tls, BindBBox(bbox, IntersectingBBoxTextline))
}
