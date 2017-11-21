package terminal

import (
	"fmt"
	"image"
	"syscall"
	"unsafe"

	"github.com/pkg/term/termios"
)

type Terminal struct {
	fd       uintptr
	old, now syscall.Termios
}

func New(fd uintptr) Terminal {
	t := Terminal{fd: fd}
	termios.Tcgetattr(fd, &t.old)
	t.now = t.old
	return t
}

func (t Terminal) set() {
	termios.Tcsetattr(t.fd, termios.TCSANOW, &t.now)
}

func (t *Terminal) Restore() {
	termios.Tcsetattr(t.fd, termios.TCSANOW, &t.old)
}

func (t Terminal) SetNoEcho() {
	t.now.Lflag &^= syscall.ECHO
	t.set()
}

func (t Terminal) SetRaw() {
	termios.Cfmakeraw(&t.now)
	t.set()
}

func (t Terminal) Bounds() (image.Rectangle, error) {
	return bounds(t.fd)
}

func (t Terminal) Size() (image.Point, error) {
	return size(t.fd)
}

func (t Terminal) SetSize(size image.Point) error {
	return setSize(t.fd, size)
}

func bounds(fd uintptr) (image.Rectangle, error) {
	size, err := size(fd)
	if err != nil {
		return image.Rectangle{}, err
	}
	return image.Rect(0, 0, size.X, size.Y), nil
}

func size(fd uintptr) (size image.Point, err error) {
	var dimensions [4]uint16
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		fd, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&dimensions)))
	if errno != 0 {
		return image.Point{}, fmt.Errorf("ioctl errno %d", errno)
	}
	return image.Pt(int(dimensions[1]), int(dimensions[0])), nil
}

func setSize(fd uintptr, size image.Point) (err error) {
	dimensions := [4]uint16{uint16(size.Y), uint16(size.X), 0, 0}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		fd, syscall.TIOCSWINSZ, uintptr(unsafe.Pointer(&dimensions)))
	if errno != 0 {
		return fmt.Errorf("ioctl errno %d", errno)
	}
	return nil
}
