package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/kriskowal/cops/display"
	"github.com/kriskowal/cops/terminal"
	"github.com/kriskowal/cops/vtio"
	"github.com/pkg/term/termios"
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

	leader, follower, err := termios.Pty()
	if err != nil {
		return err
	}

	bounds, err := term.Bounds()
	if err != nil {
		return err
	}

	if err := terminal.SetSize(follower.Fd(), bounds.Max); err != nil {
		return err
	}

	cmd := exec.Command("htop")
	cmd.Stdin = follower
	cmd.Stdout = follower
	cmd.Stderr = follower
	if err := cmd.Start(); err != nil {
		return err
	}

	vtw := vtio.NewDisplayWriter(bounds)
	go io.Copy(vtw, leader)

	front, back := display.New2(bounds)

	var buf []byte
	cur := display.Start
	buf, cur = cur.Reset(buf)
	buf, cur = cur.Home(buf)
	buf, cur = cur.Clear(buf)
	buf, cur = cur.Hide(buf)

	// Wait for keypress
	r := make(chan struct{}, 0)
	go func() {
		var rbuf [1]byte
		os.Stdin.Read(rbuf[0:1])
		close(r)
	}()

DrawLoop:
	for {
		select {
		case <-vtw.C():
			vtw.Draw(front, bounds)
			buf, cur = display.RenderOver(buf, cur, front, back, display.Model24)
			front, back = back, front
			// fmt.Printf("%q\r\n", buf)
			os.Stdout.Write(buf)
			buf = buf[0:0]
		case <-r:
			break DrawLoop
		}
	}

	buf, cur = cur.Reset(buf)
	buf, cur = cur.Home(buf)
	buf, cur = cur.Clear(buf)
	buf, cur = cur.Show(buf)
	os.Stdout.Write(buf)
	buf = buf[0:0]

	return nil
}
