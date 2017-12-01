// Package display models, composes, and renders virtual terminal displays
// using ANSI escape sequences.
// Models displays as three layers: a text layer and foreground and background
// color layers as images in any logical color space.
// The package includes colors, palettes, and rendering models for terminal
// displays supporting 0, 3, 4, 8, and 24 bit color.
// The package also includes a cursor that tracks the known position and colors
// of the cursor, appending ANSI escape sequences to incrementally modify that
// state.
package display

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/kriskowal/cops/internal"
	"github.com/kriskowal/cops/textile"
)

// New returns a new display with the given bounding rectangle.
// The rectangle need not rest at the origin.
func New(r image.Rectangle) *Display {
	return &Display{
		Background: image.NewRGBA(r),
		Foreground: image.NewRGBA(r),
		Text:       textile.New(r),
		Rect:       r,
	}
}

// New2 returns a pair of displays with the same rectangle,
// suitable for creating front and back buffers.
//
//  bounds := term.Bounds()
//	front, back := display.New2(bounds)
func New2(r image.Rectangle) (*Display, *Display) {
	return New(r), New(r)
}

// Display models a terminal display's state as three images.
type Display struct {
	Background *image.RGBA
	Foreground *image.RGBA
	Text       *textile.Textile
	Rect       image.Rectangle
	// TODO underline and intensity
}

// SubDisplay returns a mutable sub-region within the display, sharing the same
// memory.
func (d *Display) SubDisplay(r image.Rectangle) *Display {
	r = r.Intersect(d.Rect)
	return &Display{
		Background: d.Background.SubImage(r).(*image.RGBA),
		Foreground: d.Foreground.SubImage(r).(*image.RGBA),
		Text:       d.Text.SubText(r),
		Rect:       r,
	}
}

// Fill overwrites every cell with the given text and foreground and background
// colors.
func (d *Display) Fill(r image.Rectangle, t string, f, b color.Color) {
	r = r.Intersect(d.Rect)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			d.Set(x, y, t, f, b)
		}
	}
}

// Clear fills the display with transparent cells.
func (d *Display) Clear(r image.Rectangle) {
	d.Fill(r, "", color.Transparent, color.Transparent)
}

// Set overwrites the text and foreground and background colors of the cell at
// the given position.
func (d *Display) Set(x, y int, t string, f, b color.Color) {
	d.Text.Set(x, y, t)
	d.Foreground.Set(x, y, rgba(f))
	d.Background.Set(x, y, rgba(b))
}

// Draw composes one display over another. The bounds dictate the region of the
// destination.  The offset dictates the position within the source. Draw will:
//
// Overwrite the text layer for all non-empty text cells inside the rectangle.
// Fill the text with space " " to overdraw all cells.
//
// Draw the foreground of the source over the foreground of the destination
// image.  Typically, the foreground is transparent for all cells empty of
// text.  Otherwise, this operation can have interesting results.
//
// Draw the background of the source over the *background* of the destination
// image.  This allows for translucent background colors on the source image
// partially obscuring the text of the destination image.
//
// Draw the background of the source over the background of the destination
// image.
func Draw(dst *Display, r image.Rectangle, src *Display, sp image.Point, op draw.Op) {
	internal.Clip(dst.Bounds(), &r, src.Bounds(), &sp, nil, nil)
	if r.Empty() {
		return
	}
	draw.Draw(dst.Background, r, src.Background, sp, op)
	draw.Draw(dst.Foreground, r, src.Background, sp, op)
	draw.Draw(dst.Foreground, r, src.Foreground, sp, op)
	textile.Draw(dst.Text, r, src.Text, sp)
}

// At returns the text and foreground and background colors at the given
// coordinates.
func (d *Display) At(x, y int) (t string, f, b color.Color) {
	if d == nil {
		return "", Colors[7], color.Transparent
	}
	return d.Text.At(x, y), rgba(d.Foreground.At(x, y)), rgba(d.Background.At(x, y))
}

// Bounds returns the bounding rectangle of the display.
func (d *Display) Bounds() image.Rectangle {
	return d.Rect
}

// Render appends ANSI escape sequences to a byte slice to overwrite an entire
// terminal window, using the best matching colors in the terminal color model.
func Render(buf []byte, cur Cursor, over *Display, model Model) ([]byte, Cursor) {
	return RenderOver(buf, cur, over, nil, model)
}

// RenderOver appends ANSI escape sequences to a byte slice to update a
// terminal display to look like the front model, skipping cells that are the
// same in the back model, using escape sequences and the nearest matching
// colors in the given color model.
func RenderOver(buf []byte, cur Cursor, over, under *Display, model Model) ([]byte, Cursor) {
	for y := over.Rect.Min.Y; y < over.Rect.Max.Y; y++ {
		for x := over.Rect.Min.X; x < over.Rect.Max.X; x++ {
			ot, of, ob := over.At(x, y)
			ut, uf, ub := under.At(x, y)
			if len(ot) == 0 {
				ot = " "
			}
			if len(ut) == 0 {
				ut = " "
			}
			if ot == ut && of == uf && ob == ub {
				continue
			}
			buf, cur = cur.Go(buf, image.Pt(x, y))
			buf, cur = model.Render(buf, cur, of, ob)
			buf, cur = cur.WriteGlyph(buf, ot)
		}
	}
	return buf, cur
}
