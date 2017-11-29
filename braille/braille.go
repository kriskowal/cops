// Package braille composites bitmaps as braille.
package braille

import (
	"image"
	"image/color"

	"github.com/kriskowal/cops"
	"github.com/kriskowal/cops/bitmap"
	"github.com/kriskowal/cops/display"
)

// BrailleAt returns the braille bitmap that coresponds to the 2x4 grid at the
// given point in a bitmap.
func BrailleAt(src cops.BitmapReader, sp image.Point) string {
	var r rune
	if src.BitAt(sp.X, sp.Y) {
		r |= 0x1
	}
	if src.BitAt(sp.X, sp.Y+1) {
		r |= 0x2
	}
	if src.BitAt(sp.X, sp.Y+2) {
		r |= 0x4
	}
	if src.BitAt(sp.X, sp.Y+3) {
		r |= 0x40
	}
	if src.BitAt(sp.X+1, sp.Y) {
		r |= 0x8
	}
	if src.BitAt(sp.X+1, sp.Y+1) {
		r |= 0x10
	}
	if src.BitAt(sp.X+1, sp.Y+2) {
		r |= 0x20
	}
	if src.BitAt(sp.X+1, sp.Y+3) {
		r |= 0x80
	}
	return string(0x2800 + r)
}

// Draw composites a bitmap into the text and foreground layer of a display
// based on whether the colors of the source image more closely resemble the on
// or off colors of a bitmap palette.
func Draw(dst *display.Display, r image.Rectangle, src image.Image, sp image.Point, off, on color.Color) {
	// internal.Clip(dst.Bounds(), &r, src.Bounds(), &sp, nil, nil)
	r = r.Intersect(dst.Bounds())
	if r.Empty() {
		return
	}

	bits := bitmap.NewPaletted(src, off, on)
	w, h := r.Dx(), r.Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := image.Pt(x*2, y*4).Add(sp)
			dst.Text.Set(r.Min.X+x, r.Min.Y+y, BrailleAt(bits, p))
			dst.Foreground.Set(r.Min.X+x, r.Min.Y+y, on)
		}
	}
}

// Bounds takes a rectangle describing cells on a display to the cells of a
// braille bitmap covering the cells of the display.
func Bounds(r image.Rectangle) image.Rectangle {
	w, h := r.Dx(), r.Dy()
	return image.Rectangle{
		r.Min,
		r.Min.Add(image.Pt(w*2, h*4)),
	}
}
