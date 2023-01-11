package xml

import (
	goxml "encoding/xml"
)

// Paystub is the representation of a paystub page.
type Paystub struct {
	XMLName goxml.Name `xml:"pages"`
	Pages   []Page     `xml:"page"`
}

// ForAllBBox runs the function do on every bounding box found in a paystub.
func (p Paystub) ForAllBBox(do func(b BBox)) {
	for _, pg := range p.Pages {
		pg.ForAllBBox(do)
	}
}

// Page contains the XML bounding box encoding of a paystub page.
type Page struct {
	ID        string    `xml:"id,attr,omitempty"`
	BBox      BBox      `xml:"bbox,attr,omitempty"`
	Rotate    float64   `xml:"rotate,attr,omitempty"`
	Textboxes []Textbox `xml:"textbox,omitempty"`
	Rects     []Rect    `xml:"rect,omitempty"`
	Figures   []Figure  `xml:"figure"`
	Layout    Layout    `xml:"layout"`
}

// ForAllBBox runs the function do on all bounding boxes on a page.
func (p Page) ForAllBBox(do func(b BBox)) {
	if p.BBox != NullBBox() {
		do(p.BBox)
	}
	for _, tb := range p.Textboxes {
		tb.ForAllBBox(do)
	}
	for _, r := range p.Rects {
		r.ForAllBBox(do)
	}
	for _, f := range p.Figures {
		f.ForAllBBox(do)
	}
	p.Layout.ForAllBBox(do)
}

// ForAllTextlines runs the function do on all textlines on a page.
func (p Page) ForAllTextlines(do func(t Textline)) {
	for _, b := range p.Textboxes {
		b.ForAllTextlines(do)
	}
}

// Textbox is a text with its bounding box.
type Textbox struct {
	// ID is the unique ID of the textbox
	ID int `xml:"id,attr"`
	// BBox is the bounding box of this text.
	BBox BBox `xml:"bbox,attr,omitempty"`
	// Textlines are the lines comprising the text in this bounding box.
	Textlines []Textline `xml:"textline,omitempty"`
}

// ForAllBBox runs the function do on every bounding box inside a textbox.
func (t Textbox) ForAllBBox(do func(b BBox)) {
	if t.BBox != NullBBox() {
		do(t.BBox)
	}
	for _, l := range t.Textlines {
		l.ForAllBBox(do)
	}
}

// ForAllTextlines runs the function do on all textlines on a page.
func (t Textbox) ForAllTextlines(do func(l Textline)) {
	for _, l := range t.Textlines {
		do(l)
	}
}

// Text is the text that is rendered on a page.  Each letter on the page will
// typically have its own bounding box.
type Text struct {
	Font string  `xml:"font,attr,omitempty"`
	BBox BBox    `xml:"bbox,attr,omitempty"`
	Size float64 `xml:"size,attr,omitempty"`
	T    string  `xml:",chardata"`
}

func (t Text) ForAllBBox(do func(b BBox)) {
	do(t.BBox)
}

type Rect struct {
	Linewidth float64 `xml:"linewidth,attr,omitempty"`
	BBox      BBox    `xml:"bbox,attr,omitempty"`
}

func (r Rect) ForAllBBox(do func(BBox)) {
	if r.BBox == NullBBox() {
		return
	}
	do(r.BBox)
}

type Figure struct {
	Name  string  `xml:"name,attr,omitempty"'`
	BBox  BBox    `xml:"bbox,attr,omitempty"`
	Image []Image `xml:"image,omitempty"`
}

func (f Figure) ForAllBBox(do func(BBox)) {
	if f.BBox == NullBBox() {
		return
	}
	do(f.BBox)
}

type Layout struct {
	Textgroups []Textgroup `xml:"textgroup,omitempty"`
}

func (l Layout) ForAllBBox(do func(BBox)) {
	for _, t := range l.Textgroups {
		t.ForAllBBox(do)
	}
}

type Textgroup struct {
	BBox       BBox        `xml:"bbox,attr,omitempty"`
	Textgroups []Textgroup `xml:"textgroup,omitempty"`
	Textboxes  []Textbox   `xml:"textbox,omitempty"`
}

func (t Textgroup) ForAllBBox(do func(BBox)) {
	if t.BBox != NullBBox() {
		do(t.BBox)
	}
	for _, g := range t.Textgroups {
		g.ForAllBBox(do)
	}
	for _, b := range t.Textboxes {
		b.ForAllBBox(do)
	}
}

// Image is an image on the page. We remember only its dimensions, not the
// actual content.
type Image struct {
	Width  float64 `xml:"width,omitempty"`
	Height float64 `xml:"height,omitempty"`
}

var _ goxml.UnmarshalerAttr = &BBox{}

// Point represents an (X,Y) coordinate pair.
type Point struct {
	X, Y float64
}

// Monotonic returns true if the three numbers x y and z are monotonic.
func Monotonic(x, y, z float64) bool {
	return (y-z)*(y-x) <= 0
}
