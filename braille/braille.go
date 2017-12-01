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
	if r == 0 {
		return ""
	}
	return string(0x2800 + r)
}

// Draw composites a bitmap into the text and foreground layer of a display
// based on whether the colors of the source image more closely resemble the on
// or off colors of a bitmap palette.
func Draw(dst *display.Display, r image.Rectangle, src image.Image, sp image.Point, on, off color.Color) {
	// internal.Clip(dst.Bounds(), &r, src.Bounds(), &sp, nil, nil)
	r = r.Intersect(dst.Bounds())
	if r.Empty() {
		return
	}

	bits := bitmap.NewPaletted(src, off, on)
	w, h := r.Dx(), r.Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			pt := image.Pt(x*3, y*6).Add(sp)
			dx := r.Min.X + x
			dy := r.Min.Y + y
			br := BrailleAt(bits, pt)
			if br != "" {
				dst.Text.Set(dx, dy, br)
				dst.Foreground.Set(dx, dy, on)
			}
		}
	}
}

// Bounds takes a rectangle describing cells on a display to the cells of a
// braille bitmap covering the cells of the display.
func Bounds(r image.Rectangle) image.Rectangle {
	w, h := r.Dx(), r.Dy()
	return image.Rectangle{
		r.Min,
		r.Min.Add(image.Pt(w*3, h*6)).Sub(image.Pt(1, 2)),
	}
}
