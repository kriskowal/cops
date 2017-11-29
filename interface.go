package cops

import (
	"image"
	"image/color"
)

type DisplayWriter interface {
	Set(x, y int, t string, f, b color.Color)
	Bounds() image.Rectangle
}

type DisplayReader interface {
	At(x, y int) (t string, f, b color.Color)
	Bounds() image.Rectangle
}

type TextileReader interface {
	At(x, y int) string
	Bounds() image.Rectangle
}

type TextileWriter interface {
	Set(x, y int, t string)
	Bounds() image.Rectangle
}

type BitmapReader interface {
	BitAt(x, y int) bool
	Bounds() image.Rectangle
}

type BitmapWriter interface {
	SetBit(x, y int, bit bool)
	Bounds() image.Rectangle
}
