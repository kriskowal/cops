package display

import (
	"image"
	"image/color"
	"strconv"

	"github.com/kriskowal/cops"
)

type Cursor struct {
	Position   image.Point
	Foreground color.RGBA
	Background color.RGBA
}

var (
	Origin  = image.Point{}
	Unknown = image.Point{-1, -1}

	DefaultCursor = Cursor{
		Position:   Unknown,
		Foreground: cops.Colors[7],
		Background: cops.Colors[0],
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
		Foreground: cops.Colors[7],
		Background: cops.Colors[0],
	}
}

func (c Cursor) Home(buf []byte) ([]byte, Cursor) {
	return c.Go(buf, Origin)
}

func (c Cursor) Go(buf []byte, to image.Point) ([]byte, Cursor) {
	switch {
	case to.X == 0 && to.Y == c.Position.Y+1 && c.Position != Unknown:
		buf, c = c.Reset(buf)
		buf = append(buf, "\r\n"...)
	case to.Sub(image.Point{1, 0}) == c.Position:
		buf = append(buf, "\b"...)
	case to == c.Position:
	case to == Origin:
		buf = append(buf, "\033[H"...)
	default:
		buf = append(buf, "\033["...)
		buf = append(buf, strconv.Itoa(to.Y)...)
		buf = append(buf, ";"...)
		buf = append(buf, strconv.Itoa(to.X)...)
		buf = append(buf, "H"...)
	}
	c.Position = to
	return buf, c
}
