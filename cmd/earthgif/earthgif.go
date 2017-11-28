package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"os"
	"time"

	"github.com/disintegration/imaging"
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
	term := terminal.New(os.Stdin.Fd())
	defer term.Restore()
	term.SetRaw()

	data, err := Asset("earth.gif")
	if err != nil {
		return err
	}

	bounds, err := term.Bounds()
	if err != nil {
		return err
	}

	r := bytes.NewReader(data)
	imgs, err := gif.DecodeAll(r)
	if len(imgs.Image) == 0 {
		return fmt.Errorf("no frames")
	}

	front, back := display.New2(bounds)

	front.Text.Fill(" ")
	draw.Draw(front.Background, bounds, &image.Uniform{vtcolor.Colors[0]}, image.ZP, draw.Src)
	draw.Draw(front.Foreground, bounds, &image.Uniform{vtcolor.Colors[0]}, image.ZP, draw.Src)

	// Clear Home Hide
	var buf []byte
	cur := cursor.Start
	buf, cur = cur.Hide(buf)
	buf, cur = cur.Clear(buf)
	buf, cur = cur.Home(buf)

	base := imgs.Image[0]
	projection := projectCenterPreserveAspect(base.Bounds().Size(), bounds.Size())

	// Await async keypress
	keypress := make(chan byte, 1)
	go func() {
		var rbuf [1]byte
		for {
			os.Stdin.Read(rbuf[0:1])
			keypress <- rbuf[0]
		}
	}()

Loop:
	for i := 0; ; i = (i + 1) % len(imgs.Image) {
		img := imgs.Image[i]
		// Resize image and draw onto display background
		img2 := imaging.Resize(img, projection.Dx(), projection.Dy(), imaging.Lanczos)
		draw.Draw(front.Background, projection, img2, img2.Bounds().Min, draw.Over)

		// Draw frame
		buf, cur = display.Render(buf, cur, front, back, vtcolor.Model24)
		front, back = back, front
		buf, cur = cur.Home(buf)
		os.Stdout.Write(buf)
		buf = buf[0:0]

		delay := time.Duration(imgs.Delay[i]) * time.Millisecond * 10
		timer := time.NewTimer(delay)

		select {
		case <-keypress:
			if !timer.Stop() {
				<-timer.C
			}
			break Loop
		case <-timer.C:
		}
	}

	// Restore
	buf, cur = cur.Home(buf)
	buf, cur = cur.Clear(buf)
	buf, cur = cur.Show(buf)
	os.Stdout.Write(buf)
	buf = buf[0:0]

	return nil
}

func projectCenterPreserveAspect(inner, outer image.Point) image.Rectangle {
	// Account for aspect of terminal cell
	inner.X *= 2

	// Scale down, into display
	if inner.X > outer.X {
		inner.Y = inner.Y * outer.X / inner.X
		inner.X = outer.X
	}
	if inner.Y > outer.Y {
		inner.X = inner.X * outer.Y / inner.Y
		inner.Y = outer.Y
	}

	// Offset center
	offset := image.Pt(
		outer.X/2-inner.X/2,
		outer.Y/2-inner.Y/2,
	)

	return image.Rect(offset.X, offset.Y, inner.X+offset.X, inner.Y+offset.Y)
}
