package internal

import "image"

// Clip clips r against each image's bounds (after translating into the
// destination image's coordinate space) and shifts the points sp and mp by
// the same amount as the change in r.Min.
// Borrowed from "image".
func Clip(dst image.Rectangle, r *image.Rectangle, src image.Rectangle, sp *image.Point, mask *image.Rectangle, mp *image.Point) {
	orig := r.Min
	*r = r.Intersect(dst)
	*r = r.Intersect(src.Add(orig.Sub(*sp)))
	if mask != nil {
		*r = r.Intersect(mask.Add(orig.Sub(*mp)))
	}
	dx := r.Min.X - orig.X
	dy := r.Min.Y - orig.Y
	if dx == 0 && dy == 0 {
		return
	}
	sp.X += dx
	sp.Y += dy
	if mp != nil {
		mp.X += dx
		mp.Y += dy
	}
}
