package xml

// TODO(filmil): Bounding box tests.

import (
	goxml "encoding/xml"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// BBox is a bounding box.
type BBox struct {
	Left, Right, Top, Bottom float64
}

func (b BBox) String() string {
	return fmt.Sprintf("{bbox left:%v, right:%v, top:%v, bottom:%v}",
		b.Left, b.Right, b.Top, b.Bottom)
}

// NullBBox is an empty bounding box.
func NullBBox() BBox {
	return BBox{}
}

func (b BBox) BottomLeft() Point {
	return Point{X: b.Left, Y: b.Bottom}
}

func (b BBox) BottomRight() Point {
	return Point{X: b.Right, Y: b.Bottom}
}

func (b BBox) TopRight() Point {
	return Point{X: b.Right, Y: b.Top}
}

func (b BBox) TopLeft() Point {
	return Point{X: b.Left, Y: b.Top}
}

// Contains returns true if point is contained in the bounding box.
func (b BBox) Contains(p Point) bool {
	return (p.X-b.Left)*(p.X-b.Right) <= 0 &&
		(p.Y-b.Top)*(p.Y-b.Bottom) <= 0
}

// HasCornerOf returns true if box a has a corner in the box b.
func (b BBox) HasCornerOf(a BBox) bool {
	return b.Contains(a.BottomLeft()) ||
		b.Contains(a.BottomRight()) ||
		b.Contains(a.TopLeft()) ||
		b.Contains(a.TopRight())
}

// SubsumesHeight returns true f the height of box B subsumes the height of box
// a.
func (b BBox) SubsumesHeight(a BBox) bool {
	return Monotonic(b.Bottom, a.Bottom, b.Top) &&
		Monotonic(b.Bottom, a.Top, b.Top)
}

// SubsumesWidth returns true if the width of box b subsumes the width of box a.
func (b BBox) SubsumesWidth(a BBox) bool {
	return Monotonic(b.Left, a.Left, b.Right) &&
		Monotonic(b.Left, a.Left, b.Right)
}

func (b BBox) ExtendLeft() BBox {
	b.Left = math.Inf(-1)
	return b
}

func (b BBox) ExtendRight() BBox {
	b.Right = math.Inf(1)
	return b
}

func (b BBox) ExtendTop() BBox {
	b.Top = math.Inf(1)
	return b
}

func (b BBox) ExtendBottom() BBox {
	b.Bottom = math.Inf(-1)
	return b
}

func (b BBox) Widen(delta float64) BBox {
	b.Left -= delta
	b.Right += delta
	return b
}

func (b BBox) Heighten(delta float64) BBox {
	b.Top += delta
	b.Bottom -= delta
	return b
}

func (b BBox) ExtendDownTo(d float64) BBox {
	b.Bottom = d
	return b
}

// RightOf returns a box that is strictly right of b.
func (b BBox) RightOf() BBox {
	n := b.ExtendRight()
	n.Left = b.Right + eps
	return n
}

func (b BBox) Below(a BBox) BBox {
	b.Top = a.Bottom - eps
	return b
}

// IntersectingBBox returns true if boxes a and b are intersecting.
func IntersectingBBox(a, b BBox) bool {
	r := a.HasCornerOf(b) ||
		b.HasCornerOf(a) ||
		((b.SubsumesWidth(a) || b.SubsumesWidth(a)) &&
			(a.SubsumesHeight(b) || a.SubsumesHeight(b)))
	return r
}

// UnmarshallerAttr implements goxml.UnmarshallerAttr.
func (b *BBox) UnmarshalXMLAttr(attr goxml.Attr) error {
	s := strings.Split(attr.Value, ",")
	if len(s) != 4 {
		return fmt.Errorf("value not a BBox: %v", attr.Value)
	}
	var err error
	b.Left, err = strconv.ParseFloat(s[0], 64)
	if err != nil {
		return fmt.Errorf("left %v on BBox: %v", s[0], attr.Value)
	}
	b.Bottom, err = strconv.ParseFloat(s[1], 64)
	if err != nil {
		return fmt.Errorf("bottom %v on BBox: %v", s[0], attr.Value)
	}
	b.Right, err = strconv.ParseFloat(s[2], 64)
	if err != nil {
		return fmt.Errorf("right %v on BBox: %v", s[0], attr.Value)
	}
	b.Top, err = strconv.ParseFloat(s[3], 64)
	if err != nil {
		return fmt.Errorf("top %v on BBox: %v", s[0], attr.Value)
	}
	return nil
}
