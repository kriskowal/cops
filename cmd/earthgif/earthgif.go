package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"os"
	"time"

	"github.com/kriskowal/cops/display"
	"github.com/kriskowal/cops/terminal"
	"github.com/kriskowal/cops/vtcolor"
	"github.com/nfnt/resize"
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

	front := display.New(bounds)
	back := display.New(bounds)

	front.Sheet.Fill(" ")
	draw.Draw(front.Background, bounds, &image.Uniform{vtcolor.Colors[0]}, image.ZP, draw.Src)
	draw.Draw(front.Foreground, bounds, &image.Uniform{vtcolor.Colors[0]}, image.ZP, draw.Src)

	// Clear Home Hide
	var buf []byte
	cursor := display.DefaultCursor
	buf = cursor.Hide(buf)
	buf, cursor = cursor.Clear(buf)
	buf, cursor = cursor.Home(buf)

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
		img2 := resize.Resize(uint(projection.Dx()), uint(projection.Dy()), img, resize.Lanczos3)
		draw.Draw(front.Background, projection, img2, image.ZP, draw.Over)

		// Draw frame
		buf, cursor = display.Render24(buf, cursor, front, back)
		front, back = back, front
		buf, cursor = cursor.Home(buf)
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
	buf, cursor = cursor.Home(buf)
	buf, cursor = cursor.Clear(buf)
	buf = cursor.Show(buf)
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
