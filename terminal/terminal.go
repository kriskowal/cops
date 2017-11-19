package terminal

import (
	"fmt"
	"image"
	"os"
	"syscall"
	"unsafe"

	"github.com/pkg/term/termios"
)

func NoEcho(fd uintptr) func() {
	var oldcaps syscall.Termios
	var newcaps syscall.Termios
	termios.Tcgetattr(fd, &oldcaps)
	newcaps = oldcaps
	newcaps.Lflag &^= syscall.ECHO
	termios.Tcsetattr(fd, termios.TCSANOW, &newcaps)
	return func() {
		termios.Tcsetattr(os.Stdin.Fd(), termios.TCSANOW, &oldcaps)
	}
}

func Raw(fd uintptr) func() {
	var oldcaps syscall.Termios
	var newcaps syscall.Termios
	termios.Tcgetattr(fd, &oldcaps)
	newcaps = oldcaps
	termios.Cfmakeraw(&newcaps)
	termios.Tcsetattr(fd, termios.TCSANOW, &newcaps)
	return func() {
		termios.Tcsetattr(fd, termios.TCSANOW, &oldcaps)
	}
}

func GetBounds(fd uintptr) (image.Rectangle, error) {
	size, err := GetSize(fd)
	if err != nil {
		return image.Rectangle{}, err
	}
	return image.Rect(0, 0, size.X, size.Y), nil
}

func GetSize(fd uintptr) (size image.Point, err error) {
	var dimensions [4]uint16
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		fd, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&dimensions)))
	if errno != 0 {
		return image.Point{}, fmt.Errorf("ioctl errno %d", errno)
	}
	return image.Pt(int(dimensions[1]), int(dimensions[0])), nil
}

func SetSize(fd uintptr, size image.Point) (err error) {
	dimensions := [4]uint16{uint16(size.Y), uint16(size.X), 0, 0}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		fd, syscall.TIOCSWINSZ, uintptr(unsafe.Pointer(&dimensions)))
	if errno != 0 {
		return fmt.Errorf("ioctl errno %d", errno)
	}
	return nil
}
