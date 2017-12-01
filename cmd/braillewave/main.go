package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/kriskowal/cops/braille"
	"github.com/kriskowal/cops/display"
	"github.com/kriskowal/cops/terminal"
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

	ticker := time.NewTicker(16 * time.Millisecond)

	stopper := make(chan struct{}, 0)
	go func() {
		var rbuf [1]byte
		os.Stdin.Read(rbuf[0:1])
		close(stopper)
	}()

	img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	var buf []byte
	cur := display.Start
	buf, cur = cur.Hide(buf)
	buf, cur = cur.Home(buf)
	buf, cur = cur.Clear(buf)

Loop:
	for {
		t := int(time.Now().UnixNano() / 5000000)

		for x := 0; x < 1000; x++ {
			z := int(500 + 200*math.Sin(float64(t+x)*math.Pi*2/200))
			for y := 0; y < 1000; y++ {
				if y < z+10 && y > z-10 {
					img.Set(x, y, color.White)
				} else {
					img.Set(x, y, color.Black)
				}
			}
		}

		// Size that image down and write it in braille to the display.
		dis := display.New(bounds)
		bb := braille.Bounds(bounds)
		bimg := imaging.Resize(img, bb.Dx(), bb.Dy(), imaging.Lanczos)
		braille.Draw(dis, bounds, bimg, image.ZP, color.RGBA{191, 191, 127, 255}, color.Black)

		buf, cur = display.Render(buf, cur, dis, display.Model24)
		os.Stdout.Write(buf)
		buf = buf[0:0]

		select {
		case <-ticker.C:
		case <-stopper:
			break Loop
		}
	}

	ticker.Stop()

	buf, cur = cur.Home(buf)
	buf, cur = cur.Clear(buf)
	buf, cur = cur.Show(buf)
	os.Stdout.Write(buf)
	buf = buf[0:0]

	return nil
}
