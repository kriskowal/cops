package display

import (
	"image"
	"image/color"
	"strconv"

	"github.com/kriskowal/cops"
	"github.com/kriskowal/cops/sheet"
)

var colorIndex map[color.RGBA]int

func init() {
	colorIndex = make(map[color.RGBA]int, 256)
	for i := 0; i < 256; i++ {
		colorIndex[cops.Colors[i]] = i
	}
}

func New(r image.Rectangle) *Display {
	return &Display{
		Background: image.NewRGBA(r),
		Foreground: image.NewRGBA(r),
		Sheet:      sheet.NewStrings(r),
		Rect:       r,
	}
}

type Display struct {
	Background *image.RGBA
	Foreground *image.RGBA
	Sheet      *sheet.Strings
	Rect       image.Rectangle
	// TODO underline and intensity
}

type Cell struct {
	Foreground color.RGBA
	Background color.RGBA
	Text       string
}

func (d Display) Write(x, y int, s string, f, b color.Color) (int, int) {
	w := d.Sheet.Rect.Dx()
	for _, c := range s {
		d.Sheet.Set(x, y, string(c))
		d.Foreground.Set(x, y, rgba(f))
		d.Background.Set(x, y, rgba(b))
		x++
		if x >= w {
			x = 0
			y++
		}
	}
	return x, y
}

func (d Display) Fill(s string, f, b color.Color) {
	for y := d.Rect.Min.Y; y < d.Rect.Max.Y; y++ {
		for x := d.Rect.Min.X; x < d.Rect.Max.X; x++ {
			d.Sheet.Set(x, y, s)
			d.Foreground.Set(x, y, rgba(f))
			d.Background.Set(x, y, rgba(b))
		}
	}
}

func (d Display) At(x, y int) Cell {
	return Cell{
		Foreground: rgba(d.Foreground.At(x, y)),
		Background: rgba(d.Background.At(x, y)),
		Text:       d.Sheet.At(x, y),
	}
}

type RenderFunc func(buf []byte, cursor Cursor, over, under *Display) ([]byte, Cursor)

type colorFunc func(buf []byte, c color.Color) []byte

func renderer(renderForegroundColor, renderBackgroundColor colorFunc) RenderFunc {
	return RenderFunc(func(buf []byte, cursor Cursor, over, under *Display) ([]byte, Cursor) {
		for y := over.Rect.Min.Y; y < over.Rect.Max.Y; y++ {
			for x := over.Rect.Min.X; x < over.Rect.Max.X; x++ {
				o := over.At(x, y)
				u := under.At(x, y)

				// nop if empty text or previous generation of the display
				// was already correct.
				if len(o.Text) == 0 || o == u {
					continue
				}

				buf, cursor = cursor.Go(buf, image.Pt(x, y))
				if o.Foreground != cursor.Foreground {
					buf = renderForegroundColor(buf, o.Foreground)
					cursor.Foreground = o.Foreground
				}
				if o.Background != cursor.Background {
					buf = renderBackgroundColor(buf, o.Background)
					cursor.Background = o.Background
				}
				buf = append(buf, o.Text...)

				if len(o.Text) == 1 {
					cursor.Position.X++
				} else if len(o.Text) > 1 {
					// Invalidate cursor position to force position reset
					// before next draw, if the string drawn might be longer
					// than one cell wide.
					cursor.Position = Unknown
				}

			}
		}
		buf, cursor = cursor.Reset(buf)
		return buf, cursor
	})
}

func rgba(c color.Color) color.RGBA {
	return color.RGBAModel.Convert(c).(color.RGBA)
}

var Render0 = renderer(noColor, noColor)
var Render3 = renderer(renderForegroundColor3, renderBackgroundColor3)
var Render4 = renderer(renderForegroundColor4, renderBackgroundColor4)
var Render8 = renderer(renderForegroundColor8, renderBackgroundColor8)
var Render24 = renderer(renderForegroundColor24, renderBackgroundColor24)

func noColor(buf []byte, c color.Color) []byte {
	return buf
}

func renderBackgroundColor3(buf []byte, c color.Color) []byte {
	return renderBackgroundColor(buf, cops.Palette3, c)
}

func renderForegroundColor3(buf []byte, c color.Color) []byte {
	return renderForegroundColor(buf, cops.Palette3, c)
}

func renderBackgroundColor4(buf []byte, c color.Color) []byte {
	return renderBackgroundColor(buf, cops.Palette4, c)
}

func renderForegroundColor4(buf []byte, c color.Color) []byte {
	return renderForegroundColor(buf, cops.Palette4, c)
}

func renderBackgroundColor8(buf []byte, c color.Color) []byte {
	return renderBackgroundColor(buf, cops.Palette8, c)
}

func renderForegroundColor8(buf []byte, c color.Color) []byte {
	return renderForegroundColor(buf, cops.Palette8, c)
}

func renderForegroundColor(buf []byte, p color.Palette, c color.Color) []byte {
	i := p.Index(c)
	return renderForegroundColorIndex(buf, i)
}

func renderBackgroundColor(buf []byte, p color.Palette, c color.Color) []byte {
	i := p.Index(c)
	return renderBackgroundColorIndex(buf, i)
}

func renderForegroundColorIndex(buf []byte, i int) []byte {
	if i < 8 {
		buf = append(buf, "\033["...)
		buf = append(buf, strconv.Itoa(int(30+i))...)
		buf = append(buf, "m"...)
	} else if i < 16 {
		buf = append(buf, "\033["...)
		buf = append(buf, strconv.Itoa(int(90-8+i))...)
		buf = append(buf, "m"...)
	} else {
		buf = append(buf, "\033[38;5;"...)
		buf = append(buf, strconv.Itoa(int(i))...)
		buf = append(buf, "m"...)
	}
	return buf
}

func renderBackgroundColorIndex(buf []byte, i int) []byte {
	if i < 8 {
		buf = append(buf, "\033["...)
		buf = append(buf, strconv.Itoa(int(40+i))...)
		buf = append(buf, "m"...)
	} else if i < 16 {
		buf = append(buf, "\033["...)
		buf = append(buf, strconv.Itoa(int(100-8+i))...)
		buf = append(buf, "m"...)
	} else {
		buf = append(buf, "\033[48;5;"...)
		buf = append(buf, strconv.Itoa(int(i))...)
		buf = append(buf, "m"...)
	}
	return buf
}

func renderForegroundColor24(buf []byte, c color.Color) []byte {
	if i, ok := colorIndex[rgba(c)]; ok {
		return renderForegroundColorIndex(buf, i)
	}
	return renderColor24(buf, "38", c)
}

func renderBackgroundColor24(buf []byte, c color.Color) []byte {
	if i, ok := colorIndex[rgba(c)]; ok {
		return renderBackgroundColorIndex(buf, i)
	}
	return renderColor24(buf, "48", c)
}

func renderColor24(buf []byte, code string, c color.Color) []byte {
	r, g, b, _ := c.RGBA()
	buf = append(buf, "\033["...)
	buf = append(buf, code...)
	buf = append(buf, ";2;"...)
	buf = append(buf, strconv.Itoa(int(r/255))...)
	buf = append(buf, ";"...)
	buf = append(buf, strconv.Itoa(int(g/255))...)
	buf = append(buf, ";"...)
	buf = append(buf, strconv.Itoa(int(b/255))...)
	buf = append(buf, "m"...)
	return buf
}
