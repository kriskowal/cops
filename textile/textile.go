// Package textile weaves strings into a text image.
package textile

import (
	"image"
)

// Textile represents every cell in a display as a string that ideally renders
// as a single glyph. Like images and slices, the textile is a thin header
// that can share allocated memory with other textiles.
type Textile struct {
	Strings []string
	Stride  int
	Rect    image.Rectangle
}

// New returns a Textile with the given rectangle.
// As with images, the rectangle need not rest at the origin.
func New(r image.Rectangle) *Textile {
	w, h := r.Dx(), r.Dy()
	count := w * h
	buf := make([]string, count)
	return &Textile{
		Strings: buf,
		Stride:  w,
		Rect:    r,
	}
}

// Bounds returns the bounding box of the textile.
func (t *Textile) Bounds() image.Rectangle {
	return t.Rect
}

// Draw writes the text from a source textile onto a destination textile,
// within a rectangle of the destination textile, offset by a position within
// the source textile.
func Draw(dst *Textile, r image.Rectangle, src *Textile, sp image.Point) {
	// internal.Clip(dst.Bounds(), &r, src.Bounds(), &sp, nil, nil)
	r = r.Intersect(dst.Bounds())
	if r.Empty() {
		return
	}
	w, h := r.Dx(), r.Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			t := src.At(sp.X+x, sp.Y+y)
			if t != "" {
				dst.Set(r.Min.X+x, r.Min.Y+y, t)
			}
		}
	}
}

// Fill overwrites every cell in the textile with the given string.
func (t *Textile) Fill(str string) {
	area := t.Rect
	for y := area.Min.Y; y < area.Max.Y; y++ {
		for x := area.Min.X; x < area.Max.X; x++ {
			t.Set(x, y, str)
		}
	}
}

// At returns the string at a given point.
func (t *Textile) At(x, y int) string {
	if !(image.Point{x, y}.In(t.Rect)) {
		return ""
	}
	i := t.StringsOffset(x, y)
	return t.Strings[i]
}

// Set overwrites the string at a point.
func (t *Textile) Set(x, y int, str string) {
	if !(image.Point{x, y}.In(t.Rect)) {
		return
	}
	i := t.StringsOffset(x, y)
	t.Strings[i] = str
}

// SubText returns a region of text within the textile.
func (t *Textile) SubText(r image.Rectangle) *Textile {
	r = r.Intersect(t.Rect)
	if r.Empty() {
		return &Textile{}
	}
	i := t.StringsOffset(r.Min.X, r.Min.Y)
	return &Textile{
		Strings: t.Strings[i:],
		Stride:  t.Stride,
		Rect:    r,
	}
}

// StringOffset is a utility for seeking a slice of the underlying strings
// starting at the given position within the allocation.
func (t *Textile) StringsOffset(x, y int) int {
	return (y-t.Rect.Min.Y)*t.Stride + (x - t.Rect.Min.X)
}
