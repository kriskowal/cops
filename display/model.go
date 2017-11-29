package display

import (
	"image/color"
)

// Model is the interface for a terminal color rendering model.
type Model interface {
	// Render appends the ANSI sequence for changing the foreground and
	// background color to the nearest colors supported by the terminal color
	// model.
	Render(buf []byte, cur Cursor, fg, bg color.Color) ([]byte, Cursor)
}

type model struct {
	foreground func([]byte, color.Color) []byte
	background func([]byte, color.Color) []byte
}

func (m model) Render(buf []byte, cur Cursor, fg, bg color.Color) ([]byte, Cursor) {
	if fg != cur.Foreground {
		buf = m.foreground(buf, fg)
		cur.Foreground = rgba(fg)
	}
	if bg != cur.Background {
		buf = m.background(buf, bg)
		cur.Background = rgba(bg)
	}
	return buf, cur
}

var (
	// Model0 is the monochrome color model, which does not print escape
	// sequences for any colors.
	Model0 = model{renderNoColor, renderNoColor}
	// Model3 supports the first 8 color terminal palette.
	Model3 = model{renderForegroundColor3, renderBackgroundColor3}
	// Model4 supports the first 16 color terminal palette, the same as Model3
	// but doubled for high intensity variants.
	Model4 = model{renderForegroundColor4, renderBackgroundColor4}
	// Model8 supports a 256 color terminal palette, comprised of the 16
	// previous colors, a 6x6x6 color cube, and a 24 gray scale.
	Model8 = model{renderForegroundColor8, renderBackgroundColor8}
	// Model24 supports all 24 bit colors, using palette colors only for exact
	// matches.
	Model24 = model{renderForegroundColor24, renderBackgroundColor24}
)

func rgba(c color.Color) color.RGBA {
	return color.RGBAModel.Convert(c).(color.RGBA)
}
