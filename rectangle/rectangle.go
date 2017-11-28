// Package rectangle provides functions for manipulating image points and
// rectangles, particularly applicable to terminal user interfaces.
package rectangle

import "image"

// MiddleCenter finds a rectangle of the size of the inner rectangle but
// transposed into the vertical and horizontal center of the outer rectangle.
func MiddleCenter(inner, outer image.Rectangle) image.Rectangle {
	os := outer.Size()
	is := inner.Size()
	corner := os.Div(2).Sub(is.Div(2)).Add(outer.Min)
	return image.Rectangle{corner, corner.Add(is)}
}

// Outset grows
func Outset(r image.Rectangle, x, y int) image.Rectangle {
	return image.Rectangle{
		r.Min.Sub(image.Pt(x, y)),
		r.Max.Add(image.Pt(x, y)),
	}
}

// Capture returns the rectangle expanded to enclose a point.
func Capture(r image.Rectangle, pt image.Point) image.Rectangle {
	if pt.X < r.Min.X {
		r.Min.X = pt.X
	}
	if pt.X > r.Max.X {
		r.Max.X = pt.X
	}
	if pt.Y < r.Min.Y {
		r.Min.Y = pt.Y
	}
	if pt.Y > r.Max.Y {
		r.Max.Y = pt.Y
	}
	return r
}
