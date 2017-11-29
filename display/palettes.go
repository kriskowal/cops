package display

import (
	"image/color"
)

var (
	// Palette3 contains the first 8 Colors.
	Palette3 color.Palette
	// Palette4 contains the first 16 Colors.
	Palette4 color.Palette
	// Palette8 contains all 256 paletted virtual terminal colors.
	Palette8 color.Palette

	// colorIndex maps colors back to their palette index,
	// suitable for mapping arbitrary colors back to palette indexes in the 24
	// bit color model.
	colorIndex map[color.RGBA]int
)

func init() {
	for i := 0; i < 8; i++ {
		Palette3 = append(Palette3, color.Color(Colors[i]))
	}
	for i := 0; i < 16; i++ {
		Palette4 = append(Palette4, color.Color(Colors[i]))
	}
	for i := 0; i < 256; i++ {
		Palette8 = append(Palette8, color.Color(Colors[i]))
	}

	colorIndex = make(map[color.RGBA]int, 256)
	for i := 0; i < 256; i++ {
		colorIndex[Colors[i]] = i
	}
}
