// Package geometry provides functions for manipulating image points and
// rectangles, particularly applicable to terminal user interfaces.
package geometry

import "image"

func MiddleCenter(inner, outer image.Rectangle) image.Rectangle {
	os := outer.Size()
	is := inner.Size()
	corner := scalePt(os, 1, 2).Sub(scalePt(is, 1, 2)).Add(outer.Min)
	return image.Rectangle{corner, corner.Add(is)}
}

func Outset(r image.Rectangle, x, y int) image.Rectangle {
	return image.Rectangle{
		r.Min.Sub(image.Pt(x, y)),
		r.Max.Add(image.Pt(x, y)),
	}
}

func scalePt(p image.Point, up int, down int) image.Point {
	return image.Pt(p.X*up/down, p.Y*up/down)
}
