package display

import (
	"image/color"
	"strconv"
)

func renderNoColor(buf []byte, c color.Color) []byte {
	return buf
}

func renderBackgroundColor3(buf []byte, c color.Color) []byte {
	return renderBackgroundColor(buf, Palette3, c)
}

func renderForegroundColor3(buf []byte, c color.Color) []byte {
	return renderForegroundColor(buf, Palette3, c)
}

func renderBackgroundColor4(buf []byte, c color.Color) []byte {
	return renderBackgroundColor(buf, Palette4, c)
}

func renderForegroundColor4(buf []byte, c color.Color) []byte {
	return renderForegroundColor(buf, Palette4, c)
}

func renderBackgroundColor8(buf []byte, c color.Color) []byte {
	return renderBackgroundColor(buf, Palette8, c)
}

func renderForegroundColor8(buf []byte, c color.Color) []byte {
	return renderForegroundColor(buf, Palette8, c)
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
	buf = append(buf, strconv.Itoa(int(r/256))...)
	buf = append(buf, ";"...)
	buf = append(buf, strconv.Itoa(int(g/256))...)
	buf = append(buf, ";"...)
	buf = append(buf, strconv.Itoa(int(b/256))...)
	buf = append(buf, "m"...)
	return buf
}
