package vtio

import (
	// "fmt"

	"image"
	"image/color"
	"image/draw"
	"sync"

	ansiterm "github.com/Azure/go-ansiterm"
	"github.com/kriskowal/cops/display"
)

func NewDisplayWriter(rect image.Rectangle) *DisplayWriter {
	dis := display.New(rect)
	handler := &displayWriterHandler{
		dis: dis,
		fg:  display.Colors[7],
		bg:  display.Colors[0],
		c:   make(chan struct{}, 1),
	}
	par := ansiterm.CreateParser("Ground", handler)
	return &DisplayWriter{
		parser:  par,
		handler: handler,
	}
}

type DisplayWriter struct {
	parser  *ansiterm.AnsiParser
	handler *displayWriterHandler
}

func (d *DisplayWriter) C() <-chan struct{} {
	return d.handler.c
}

func (d *DisplayWriter) Write(buf []byte) (int, error) {
	d.handler.lock.Lock()
	defer d.handler.lock.Unlock()

	count, err := d.parser.Parse(buf)

	select {
	case d.handler.c <- struct{}{}:
	default:
	}

	return count, err
}

func (d *DisplayWriter) Draw(e *display.Display, r image.Rectangle) {
	d.handler.lock.RLock()
	defer d.handler.lock.RUnlock()
	display.Draw(e, r, d.handler.dis, image.ZP, draw.Src)
}

func (d *DisplayWriter) Resize(rect image.Rectangle) {
	d.handler.lock.Lock()
	defer d.handler.lock.Unlock()
	dis := display.New(rect)
	display.Draw(dis, rect, d.handler.dis, image.ZP, draw.Src)
	d.handler.dis = dis
	d.handler.rect = rect
}

type displayWriterHandler struct {
	lock sync.RWMutex
	c    chan struct{}
	dis  *display.Display
	pos  image.Point
	rect image.Rectangle
	fg   color.Color
	bg   color.Color
	buf  []byte
}

func (h *displayWriterHandler) scroll(dy int) {
	display.Draw(h.dis, h.rect, h.dis, image.Pt(0, dy), draw.Src)
	h.dis.Clear(image.Rect(0, h.rect.Max.Y-dy, h.rect.Max.X, h.rect.Max.Y))
}

func (h *displayWriterHandler) Flush() error {
	// fmt.Printf("F %q\r\n", string(h.buf))
	for _, r := range string(h.buf) {
		// TODO deal with zero-width-joiner and combined characters
		if h.pos.X > h.dis.Rect.Max.X {
			h.pos.X = 0
			h.pos.Y++
			if h.pos.Y > h.dis.Rect.Max.Y {
				h.scroll(1)
				h.pos.Y--
			}
		}
		h.dis.Set(h.pos.X, h.pos.Y, string(r), h.fg, h.bg)
		h.pos.X++
	}
	h.buf = h.buf[0:0]

	return nil
}

func (h *displayWriterHandler) Print(b byte) error {
	h.buf = append(h.buf, b)
	return nil
}

func (h *displayWriterHandler) Execute(b byte) error {
	// fmt.Printf("E %q\n", string(b))
	h.Flush()
	switch b {
	case '\n':
		h.pos.Y++
	case '\r':
		h.pos.X = 0
	case '\t':
		h.pos.X = (h.pos.X + 8) / 8 * 8
	case '\b':
		if h.pos.X > 0 {
			h.pos.X--
		}
	}
	return nil
}

// Cursor up
func (h *displayWriterHandler) CUU(i int) error {
	h.Flush()
	// fmt.Printf("CUU\n")
	h.pos.Y -= i
	// if h.pos.Y < h.dis.Rect.Min.Y {
	// 	h.pos.Y = h.dis.Rect.Min.Y
	// }
	return nil
}

// Cursor down
func (h *displayWriterHandler) CUD(i int) error {
	h.Flush()
	// fmt.Printf("CUD\n")
	h.pos.Y += i
	// if h.pos.Y >= h.dis.Rect.Max.Y {
	// 	h.pos.Y = h.dis.Rect.Max.Y - 1
	// }
	return nil
}

// Cursor forward
func (h *displayWriterHandler) CUF(i int) error {
	h.Flush()
	// fmt.Printf("CUF\n")
	h.pos.X += i
	// if h.pos.X >= h.dis.Rect.Max.X {
	// 	h.pos.X = h.dis.Rect.Max.X - 1
	// }
	return nil
}

// Cursor backward
func (h *displayWriterHandler) CUB(i int) error {
	h.Flush()
	// fmt.Printf("CUB\n")
	h.pos.X -= i
	// if h.pos.X < h.dis.Rect.Min.X {
	// 	h.pos.X = h.dis.Rect.Min.X
	// }
	return nil
}

// Cursor new line?
func (h *displayWriterHandler) CNL(i int) error {
	h.Flush()
	// fmt.Printf("CNL\n")
	return nil
}

func (h *displayWriterHandler) CPL(int) error {
	h.Flush()
	// fmt.Printf("CPL\n")
	return nil
}

func (h *displayWriterHandler) CHA(i int) error {
	h.Flush()
	h.pos.X = i - 1
	return nil
}

// Vertical line position absolute
func (h *displayWriterHandler) VPA(i int) error {
	// fmt.Printf("VPA\r\n")
	h.pos.Y = i - 1
	return nil
}

// Cursor update
func (h *displayWriterHandler) CUP(y, x int) error {
	// fmt.Printf("CUP\r\n")
	h.Flush()
	h.pos.X = x - 1
	h.pos.Y = y - 1
	return nil
}

// Horizontal vertical position
func (h *displayWriterHandler) HVP(y, x int) error {
	h.CUP(y, x)
	return nil
}

func (h *displayWriterHandler) DECTCEM(bool) error {
	return nil
}

func (h *displayWriterHandler) DECOM(bool) error {
	return nil
}

func (h *displayWriterHandler) DECCOLM(bool) error {
	return nil
}

// Erase display
func (h *displayWriterHandler) ED(i int) error {
	// fmt.Printf("ED\r\n")
	h.Flush()
	return nil
}

// Erase line
func (h *displayWriterHandler) EL(i int) error {
	// fmt.Printf("EL %d\r\n", i)
	switch i {
	case 0:
		h.dis.Fill(image.Rect(h.pos.X, h.pos.Y, h.rect.Max.X, h.pos.Y+1), " ", h.fg, h.bg)
	case 1:
		h.dis.Fill(image.Rect(0, h.pos.Y, h.pos.X, h.pos.Y+1), " ", h.fg, h.bg)
	case 2:
		h.dis.Fill(image.Rect(0, h.pos.Y, h.rect.Max.X, h.pos.Y+1), " ", h.fg, h.bg)
	}
	return nil
}

// Insert line
func (h *displayWriterHandler) IL(int) error {
	// fmt.Printf("IL\r\n")
	return nil
}

// Delete line
func (h *displayWriterHandler) DL(int) error {
	// fmt.Printf("DL\r\n")
	return nil
}

// Insert column
func (h *displayWriterHandler) ICH(int) error {
	// fmt.Printf("ICH\r\n")
	return nil
}

// Delete column
func (h *displayWriterHandler) DCH(int) error {
	// fmt.Printf("DCH\r\n")
	return nil
}

// Set graphics rendition
func (h *displayWriterHandler) SGR(codes []int) error {
	// fmt.Printf("SGR %#v\r\n", codes)
	h.Flush()

	if len(codes) == 0 {
		h.fg = display.Colors[7]
		h.bg = display.Colors[0]
	}

	for len(codes) > 0 {
		code := codes[0]
		codes = codes[1:]
		switch {

		case code == 0: // reset
			h.fg = display.Colors[7]
			h.bg = display.Colors[0]
		case code == 1: // TODO high intensity, right?

		case code >= 30 && code < 38: // set foreground color
			h.fg = display.Colors[code-30]
		case code >= 90 && code < 98: // set high intensity foreground color
			h.fg = display.Colors[code-90+8]
		case code == 39:
			h.fg = display.Colors[7]
		case code == 38: // set foreground color
			h.fg, codes = colorForCodes(codes)

		case code >= 40 && code < 48: // set background color
			h.bg = display.Colors[code-40]
		case code >= 100 && code < 108: // set high intensity background color
			h.bg = display.Colors[code-100+8]
		case code == 48: // set background color
			h.bg, codes = colorForCodes(codes)
		case code == 49:
			h.bg = display.Colors[0]
		}
	}
	return nil
}

func colorForCodes(codes []int) (color.RGBA, []int) {
	if len(codes) == 0 {
		return display.Colors[0], codes
	}
	code := codes[0]
	codes = codes[1:]
	switch {
	case code == 5:
		if len(codes) < 1 {
			return display.Colors[0], codes
		}
		return display.Colors[codes[0]], codes[1:]
	case code == 2:
		if len(codes) < 3 {
			return display.Colors[0], codes
		}
		return color.RGBA{
			byte(codes[0]),
			byte(codes[1]),
			byte(codes[2]),
			255,
		}, codes[3:]
	}

	return display.Colors[0], codes
}

// Scroll up
func (h *displayWriterHandler) SU(int) error {
	// fmt.Printf("SU\r\n")
	return nil
}

// Scroll down
func (h *displayWriterHandler) SD(int) error {
	// fmt.Printf("SD\r\n")
	return nil
}

func (h *displayWriterHandler) DA([]string) error {
	// fmt.Printf("DA\r\n")
	return nil
}

func (h *displayWriterHandler) DECSTBM(t, b int) error {
	// fmt.Printf("DECSTBM %d %d\r\n", t, b)
	return nil
}

func (h *displayWriterHandler) IND() error {
	// fmt.Printf("IND\r\n")
	return nil
}

func (h *displayWriterHandler) RI() error {
	// fmt.Printf("RI\r\n")
	return nil
}
