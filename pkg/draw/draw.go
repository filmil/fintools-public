// Package draw draws bounding boxes.
package draw

import (
	"image"

	"github.com/filmil/fintools-public/pkg/xml"
	"github.com/llgcode/draw2d"
)

// Box draws a box into the given graphic context.
func FillBox(gc draw2d.GraphicContext, b xml.BBox) {
	gc.BeginPath()
	gc.MoveTo(b.Left, b.Bottom)
	gc.LineTo(b.Right, b.Bottom)
	gc.LineTo(b.Right, b.Top)
	gc.LineTo(b.Left, b.Top)
	gc.LineTo(b.Left, b.Bottom)
	gc.Close()
	gc.FillStroke()
}

// Box draws a box into the given graphic context.
func Box(gc draw2d.GraphicContext, b xml.BBox) {
	gc.BeginPath()
	gc.MoveTo(b.Left, b.Bottom)
	gc.LineTo(b.Right, b.Bottom)
	gc.LineTo(b.Right, b.Top)
	gc.LineTo(b.Left, b.Top)
	gc.LineTo(b.Left, b.Bottom)
	gc.Close()
	gc.Stroke()
}

func WithCtx(gc draw2d.GraphicContext, f func(gc draw2d.GraphicContext, b xml.BBox)) func(xml.BBox) {
	return func(b xml.BBox) {
		f(gc, b)
	}
}

func ImageForPage(p xml.Page) *image.RGBA {
	b := p.BBox
	const ppi = 10
	r := image.Rect(
		int(ppi*b.Left),
		int(ppi*b.Bottom),
		int(ppi*b.Right),
		int(ppi*b.Top))
	return image.NewRGBA(r)
}
