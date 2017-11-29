package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"github.com/disintegration/imaging"
	"github.com/kriskowal/cops/braille"
	"github.com/kriskowal/cops/cursor"
	"github.com/kriskowal/cops/display"
	"github.com/kriskowal/cops/terminal"
	"github.com/kriskowal/cops/vtcolor"
)

func main() {
	if err := Main(); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func Main() error {
	term := terminal.New(os.Stdout.Fd())
	defer term.Restore()
	term.SetRaw()

	bounds, err := term.Bounds()
	if err != nil {
		return err
	}

	front, back := display.New2(bounds)

	// Draw a circle as a stand-alone image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if float64(y) > 50+50*math.Sin(float64(x)*math.Pi*2/100) {
				img.Set(x, y, color.White)
			}
		}
	}

	// Size that image down and write it in braille to the display.
	bsz := braille.Bounds(bounds)
	img2 := imaging.Resize(img, bsz.Dx(), bsz.Dy(), imaging.Lanczos)
	braille.Draw(front, front.Bounds(), img2, image.ZP, color.Black, color.White)

	var buf []byte
	cur := cursor.Start
	buf, cur = cur.Hide(buf)
	buf, cur = cur.Home(buf)
	buf, cur = cur.Clear(buf)
	buf, cur = display.Render(buf, cur, front, back, vtcolor.Model24)
	os.Stdout.Write(buf)
	front, back = back, front
	buf = buf[0:0]

	var rbuf [1]byte
	os.Stdin.Read(rbuf[0:1])

	buf, cur = cur.Home(buf)
	buf, cur = cur.Clear(buf)
	buf, cur = cur.Show(buf)
	os.Stdout.Write(buf)
	buf = buf[0:0]

	return nil
}
