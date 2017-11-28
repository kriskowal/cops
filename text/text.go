// Package text measures and cuts raw text for terminal displays.
// Rather that implement the gaummut of virtual terminal commands,
// the text package recognizes only newline "\n", tab "\t", and space " ",
// assuming all other characters are printable in a single cell.
// The text package treats white space as transparent, only writing the text
// and foreground color layer for each cell that contains opaque text.
package text

import (
	"image"
	"image/color"

	"github.com/kriskowal/cops/display"
)

const tabStopWidth = 8

// Bounds measures a bounding box that would hold the given message.
func Bounds(str string) image.Rectangle {
	width, height := 0, 0
	x, y := 0, 1
	for _, r := range str {
		if r == '\n' {
			y++
			x = 0
		} else if r == '\t' {
			x = ((x + tabStopWidth) / tabStopWidth) * tabStopWidth
		} else {
			x++
			if x > width {
				width = x
			}
			if y > height {
				height = y
			}
		}
	}
	return image.Rect(0, 0, width, height)
}

// Write draws a message onto a display in the given bounds and with the given color.
func Write(dst *display.Display, bounds image.Rectangle, str string, f color.Color) {
	x, y := 0, 0
	for _, r := range str {
		if r == '\n' {
			y++
			x = 0
		} else if r == '\r' {
		} else if r == '\t' {
			x = ((x + tabStopWidth) / tabStopWidth) * tabStopWidth
		} else if r == ' ' {
			x++
		} else {
			pt := image.Pt(x, y).Add(bounds.Min)
			if pt.In(bounds) {
				dst.Text.Set(pt.X, pt.Y, string(r))
				dst.Foreground.Set(pt.X, pt.Y, f)
			}
			x++
		}
	}
}
