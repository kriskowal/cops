package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/kriskowal/cops/bitmap"
	"github.com/kriskowal/cops/braille"
	"github.com/kriskowal/cops/display"
)

func main() {
	if err := Main(); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func Main() error {

	w, h := 16, 16
	pb := image.Rect(0, 0, w, h)
	ib := braille.Bounds(pb)
	rb := ib
	pb.Max.X += 2
	rb = rb.Add(image.Pt(2, 0))
	page := display.New(pb)
	img := bitmap.New(ib, color.Black, color.White)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			b := w*y + x
			if b&0x01 != 0 {
				img.Set(x*2, y*4, color.White)
			}
			if b&0x02 != 0 {
				img.Set(x*2+1, y*4, color.White)
			}
			if b&0x04 != 0 {
				img.Set(x*2, y*4+1, color.White)
			}
			if b&0x08 != 0 {
				img.Set(x*2+1, y*4+1, color.White)
			}

			if b&0x10 != 0 {
				img.Set(x*2, y*4+2, color.White)
			}
			if b&0x20 != 0 {
				img.Set(x*2+1, y*4+2, color.White)
			}
			if b&0x40 != 0 {
				img.Set(x*2, y*4+3, color.White)
			}
			if b&0x80 != 0 {
				img.Set(x*2+1, y*4+3, color.White)
			}
		}
	}

	page.SubDisplay(image.Rect(0, 0, 1, h)).Fill(string(0x28ff), display.Colors[8], color.Transparent)
	braille.Draw(page, rb, img, image.ZP, color.White, color.Transparent)

	var buf []byte
	cur := display.Reset
	buf, cur = display.Render(buf, cur, page, display.Model8)
	buf = append(buf, "\r\n"...)
	os.Stdout.Write(buf)

	return nil
}
