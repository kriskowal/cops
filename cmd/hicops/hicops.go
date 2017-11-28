package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/kriskowal/cops/cursor"
	"github.com/kriskowal/cops/display"
	"github.com/kriskowal/cops/geometry"
	"github.com/kriskowal/cops/terminal"
	"github.com/kriskowal/cops/text"
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

	bounds, err := term.Bounds()
	if err != nil {
		return err
	}

	front := display.New(bounds)
	back := display.New(bounds)

	front.Fill("/", color.RGBA{192, 0, 0, 255}, color.RGBA{30, 20, 40, 255})

	msg := "Press any key to continue..."
	msgbox := text.Bounds(msg)
	inset := geometry.MiddleCenter(msgbox, bounds)
	outset := geometry.Outset(inset, 4, 2)
	panel := display.New(outset)
	// Fill the panel with a translucent background color.
	draw.Draw(panel.Background, outset, &image.Uniform{color.NRGBA{63, 63, 127, 127}}, image.ZP, draw.Over)
	// Draw our text in the panel.
	text.Write(panel, inset, msg, vtcolor.Colors[7])
	display.Draw(front, outset, panel, outset.Min, draw.Over)

	var buf []byte
	cur := cursor.Start
	buf, cur = cur.Hide(buf)
	buf, cur = cur.Clear(buf)
	buf, cur = cur.Home(buf)
	buf, cur = display.Render(buf, cur, front, back, vtcolor.Model24)
	back, front = front, back
	buf, cur = cur.Home(buf)
	os.Stdout.Write(buf)
	buf = buf[0:0]

	var input [1]byte
	os.Stdin.Read(input[0:1])

	buf, cur = cur.Home(buf)
	buf, cur = cur.Clear(buf)
	buf, cur = cur.Show(buf)
	os.Stdout.Write(buf)
	buf = buf[0:0]

	return nil
}
