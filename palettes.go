package cops

import "image/color"

var Palette3 color.Palette
var Palette4 color.Palette
var Palette8 color.Palette

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
}
