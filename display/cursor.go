package display

import (
	"image"
	"image/color"
	"strconv"

	"github.com/kriskowal/cops/vtcolor"
)

type Cursor struct {
	Position   image.Point
	Foreground color.RGBA
	Background color.RGBA
}

var (
	Origin  = image.ZP
	Unknown = image.Point{-1, -1}

	DefaultCursor = Cursor{
		Position:   Unknown,
		Foreground: vtcolor.Colors[7],
		Background: vtcolor.Colors[0],
	}
)

func (c Cursor) Hide(buf []byte) []byte {
	return append(buf, "\033[?25l"...)
}

func (c Cursor) Show(buf []byte) []byte {
	return append(buf, "\033[?25h"...)
}

func (c Cursor) Clear(buf []byte) ([]byte, Cursor) {
	// Clear implicitly invalidates the cursor position since its behavior is
	// consistent across terminal implementations.
	return append(buf, "\033[2J"...), Cursor{
		Position:   Unknown,
		Foreground: c.Foreground,
		Background: c.Background,
	}
}

func (c Cursor) Reset(buf []byte) ([]byte, Cursor) {
	if c == DefaultCursor {
		return buf, c
	}
	return append(buf, "\033[m"...), Cursor{
		Position:   c.Position,
		Foreground: vtcolor.Colors[7],
		Background: vtcolor.Colors[0],
	}
}

func (c Cursor) Home(buf []byte) ([]byte, Cursor) {
	c.Position = Origin
	return append(buf, "\033[H"...), c
}

func (c Cursor) Go(buf []byte, to image.Point) ([]byte, Cursor) {
	if c.Position == Unknown {
		// If the cursor position is completely unknown, move relative to
		// screen origin. This mode must be avoided to render relative to
		// cursor position inline with a scrolling log, by setting the cursor
		// position relative to an arbitrary origin before rendering.
		buf = append(buf, "\033["...)
		buf = append(buf, strconv.Itoa(to.Y)...)
		buf = append(buf, ";"...)
		buf = append(buf, strconv.Itoa(to.X)...)
		buf = append(buf, "H"...)
		c.Position = to
		return buf, c
	}

	if c.Position.X == -1 {
		// If only horizontal position is unknown, return to first column and
		// march forward.
		// Rendering a non-ASCII cell of unknown or indeterminite width may
		// invalidate the column number.
		// For example, a skin tone emoji may or may not render as a single
		// column glyph.
		buf = append(buf, "\r"...)
		c.Position.X = 0
		// Continue...
	}

	if to.X == 0 && to.Y == c.Position.Y+1 {
		buf, c = c.Reset(buf)
		buf = append(buf, "\r\n"...)
		c.Position.X = 0
		c.Position.Y++
	} else if to.X == 0 && c.Position.X != 0 {
		buf, c = c.Reset(buf)
		buf = append(buf, "\r"...)
		c.Position.X = 0

		// In addition to scrolling back to the first column generally, this
		// has the effect of resetting the column if writing a multi-byte
		// string invalidates the cursor's horizontal position.
		// For example, a skin tone emoji may or may not render as a single
		// column glyph.
	}

	if to.Y < c.Position.Y {
		buf = append(buf, "\033["...)
		buf = append(buf, strconv.Itoa(c.Position.Y-to.Y)...)
		buf = append(buf, "A"...)
	} else if to.Y > c.Position.Y {
		buf = append(buf, "\033["...)
		buf = append(buf, strconv.Itoa(to.Y-c.Position.Y)...)
		buf = append(buf, "B"...)
	}
	if to.X < c.Position.X {
		buf = append(buf, "\033["...)
		buf = append(buf, strconv.Itoa(c.Position.X-to.X)...)
		buf = append(buf, "D"...)
	} else if to.X > c.Position.X {
		buf = append(buf, "\033["...)
		buf = append(buf, strconv.Itoa(to.X-c.Position.X)...)
		buf = append(buf, "C"...)
	}

	c.Position = to
	return buf, c
}
