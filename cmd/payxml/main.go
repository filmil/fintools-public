// Package main contains an experimental XML parsing program that produces
// an image of the decoded bounding boxes as a PNG file.
//
// Usage:
//   payxml -input=<xml_file> -output=<png_file>
package main

import (
	"flag"
	"image/color"
	"os"

	"github.com/filmil/fintools/pkg/draw"
	"github.com/filmil/fintools/pkg/xml"
	"github.com/golang/glog"
	"github.com/llgcode/draw2d/draw2dimg"
)

var (
	input  = flag.String("input", "", "Input filename")
	output = flag.String("output", "output.png", "output filename")
)

func main() {
	flag.Parse()

	if *input == "" {
		glog.Fatalf("--input=... is mandatory")
	}

	file, err := os.Open(*input)
	if err != nil {
		glog.Fatalf("can not open %q: %v", *input, err)
	}

	paystub, err := xml.Decode(file)
	if err != nil {
		glog.Fatalf("xml.Parse(%q)=%v", *input, err)
	}
	page := paystub.Pages[0]
	dest := draw.ImageForPage(page)
	gc := draw2dimg.NewGraphicContext(dest)
	gc.SetDPI(72)
	// Set some properties
	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.Scale(10, -10)
	gc.Translate(0, -page.BBox.Top)
	gc.SetLineWidth(1)
	draw.FillBox(gc, page.BBox)
	page.ForAllBBox(draw.WithCtx(gc, draw.Box))

	// Save to file
	draw2dimg.SaveToPngFile(*output, dest)
}
