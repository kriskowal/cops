package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/kriskowal/cops"
	"github.com/kriskowal/cops/display"
	"github.com/kriskowal/cops/terminal"
)

func main() {
	if err := Main(); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func Main() error {
	bounds, err := terminal.GetBounds(0)
	if err != nil {
		return err
	}

	front := display.New(bounds)
	back := display.New(bounds)

	front.Fill(".", color.RGBA{31, 31, 31, 255}, cops.Colors[0])
	front.Write(0, bounds.Max.Y/2, "Press any key to continue...", color.RGBA{0, 127, 127, 255}, cops.Colors[0])

	var buf []byte
	cursor := display.DefaultCursor
	buf = cursor.Hide(buf)
	buf, cursor = cursor.Clear(buf)
	buf, cursor = cursor.Home(buf)
	buf, cursor = display.Render24(buf, cursor, front, back)
	back, front = front, back
	buf, cursor = cursor.Home(buf)
	os.Stdout.Write(buf)
	buf = buf[0:0]

	restore := terminal.Raw(os.Stdin.Fd())
	defer restore()
	var input [1]byte
	os.Stdin.Read(input[0:1])

	buf, cursor = cursor.Home(buf)
	buf, cursor = cursor.Clear(buf)
	buf = cursor.Show(buf)
	os.Stdout.Write(buf)
	buf = buf[0:0]

	return nil
}
