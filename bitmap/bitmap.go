package bitmap

import (
	"image"
	"image/color"
)

// Bitmap is a compact bitmap image with a two-color palette.
type Bitmap struct {
	Bytes   []byte
	Stride  int
	Rect    image.Rectangle
	Palette color.Palette
}

// New returns a bitmap with the given rectangle and two-color palette.
func New(r image.Rectangle, off, on color.Color) *Bitmap {
	w, h := r.Dx(), r.Dy()
	stride := (w + 7) / 8
	count := stride * h
	return &Bitmap{
		Bytes:   make([]byte, count),
		Stride:  stride,
		Rect:    r,
		Palette: color.Palette{off, on},
	}
}

// At returns the color at a point.
func (b *Bitmap) At(x, y int) color.Color {
	if b.BitAt(x, y) {
		return b.Palette[1]
	}
	return b.Palette[0]
}

// Set sets the color at a point.
func (b *Bitmap) Set(x, y int, c color.Color) {
	b.BitSet(x, y, color.Model(b.Palette).Convert(c) != b.Palette[0])
}

// ColorModel returns the bitmap's palette.
func (b *Bitmap) ColorModel() color.Model {
	return b.Palette
}

// BitAt returns whether the bit is set at a point.
func (b *Bitmap) BitAt(x, y int) bool {
	if !image.Pt(x, y).In(b.Rect) {
		return false
	}
	i := y*b.Stride + x/8
	by := b.Bytes[i]
	return by&(1<<uint(x)&0x7) != 0
}

// BitSet sets or resets the bit at a point.
func (b *Bitmap) BitSet(x, y int, bit bool) {
	if !image.Pt(x, y).In(b.Rect) {
		return
	}
	i := y*b.Stride + x/8
	if bit {
		b.Bytes[i] |= 1 << uint(x) & 07
	} else {
		b.Bytes[i] &^= 1 << uint(x) & 07
	}
}
