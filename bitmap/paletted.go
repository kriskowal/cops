package bitmap

import (
	"image"
	"image/color"
)

// NewPaletted adapts an arbitrary image to a readable bitmap based on a
// two-color palette.
func NewPaletted(img image.Image, off, on color.Color) *PalettedBitmap {
	return &PalettedBitmap{
		Image:   img,
		Palette: color.Palette{off, on},
	}
}

// PalettedBitmap is a view of another image, rendered down to the nearest of
// two colors.
type PalettedBitmap struct {
	Image   image.Image
	Palette color.Palette
}

// BitAt returns whether the bit at a point more closely resembles the "on"
// color in the palette.
func (p *PalettedBitmap) BitAt(x, y int) bool {
	c := p.Image.At(x, y)
	return color.Model(p.Palette).Convert(c) != p.Palette[0]
}

// Bounds returns the bounds of the underlying image.
func (p *PalettedBitmap) Bounds() image.Rectangle {
	return p.Image.Bounds()
}
