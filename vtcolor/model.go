package vtcolor

import "image/color"

// Model is the interface for a terminal color rendering model.
type Model interface {
	// RenderForegroundColor appends the byte sequence for changing the
	// foreground color to the nearest color supported by the model.
	RenderForegroundColor([]byte, color.Color) []byte
	// RenderBackgroundColor appends the byte sequence for changing the
	// background color to the nearest color supported by the model.
	RenderBackgroundColor([]byte, color.Color) []byte
}

type model struct {
	foreground func([]byte, color.Color) []byte
	background func([]byte, color.Color) []byte
}

func (m model) RenderForegroundColor(buf []byte, col color.Color) []byte {
	return m.foreground(buf, col)
}

func (m model) RenderBackgroundColor(buf []byte, col color.Color) []byte {
	return m.background(buf, col)
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
