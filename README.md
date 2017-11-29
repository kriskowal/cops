
# cops

Cops is a Go library for rendering terminal user interfaces.
Cops supports 24 bit color and pairs well with Go's `color`, `image`, and
`image/draw` packages.

Cops models a terminal Display as an image with three layers:

- `Text *"github.com/kriskowal/cops/textile".Textile`
- `Foreground *"image".RGBA`
- `Background *"image".RGBA`

The display package provides a display type that models these three layers.
Since the foreground and background layers are standard Go images,
we can use Go's `draw` package and third-party image processing packages to
composite color layers.

```go
bounds := image.Rect(0, 0, 80, 23)
disp := display.New(bounds)
```

## draw

Cops can draw displays with alpha transparency channels over lower display
layers.
Use the `Draw` method to compose displays in layers.

```go
display.Draw(
    dst *display.Display,
    r image.Rectangle,
    src *display.Display,
    sp image.Point,
    op draw.Op, // draw.Over or draw.Src
)
```

Drawing will:

- Overwrite the text layer for all non-empty text cells inside the rectangle.
  Fill the text with space " " to overdraw all cells.
- Draw the foreground of the source over the foreground of the destination image.
  Typically, the foreground is transparent for all cells empty of text.
  Otherwise, this operation can have interesting results.
- Draw the background of the source over the *foreground* of the destination image.
  This allows for translucent background colors on the source image partially
  obscuring the text of the destination image.
- Draw the background of the source over the background of the destination image.

Cops defers the decision to render to 3, 4, 8, or 24 bit terminal color model
to the very last phase of rendering, so application authors are free to use the
gammut of any color model supported by Go, including third-party color models
like HUSLuv, or pluck from the terminal 256 color palette in `display.Colors`.

## render

The display package provides `Render` and `RenderOver` methods.
The render method produces a sequence of bytes to write that will update a
terminal, skipping over cells that have not changed.

```go
Render(
    buf []byte,
    cur display.Cursor,
    dis *display.Display,
    model display.Model,
) (
    buf []byte,
    cur display.Cursor,
)

RenderOver(
    buf []byte,
    cur display.Cursor,
    over, under *display.Display,
    model display.Model,
) (
    buf []byte,
    cur display.Cursor,
)
```

The default display is blank and rendering it will cause no changes.
Rendering a non-blank display over a blank display will effect a full display
rewrite.

Render, like `append`, accepts and returns a slice of bytes, prefering to reuse
the prior allocation, growing the allocation only when necessary.

Typical terminal applications swap a front and back display, drawing
each frame over the previous.

```go
var buf []byte
cur := display.Start
front, back := display.New2(bounds)
buf, cur = display.Render(buf, cur, front, back, display.Model24)
front, back = back, front
buf = buf[0:0]
```

Although the display models all colors in 32 bit RGBA, the color model
samples these colors down to the terminal's supported color model.

## cursor

Render accepts the current cursor state and returns the cursor state after
applying the rendered bytes to the terminal.
The `display` package provides the cursor type, which has methods for updating
the cursor's position, foreground color, and background color.  Each of these
methods append to a buffer and return the resulting cursor state.

The initial cursor state is unknown, assuming nothing about the cursor
position or coloring.

```go
cur := display.Start
```

The differential update will attempt to use relative cursor position changes
whenever possible, resorting to changes relative to the beginning of the same
line if it loses track of its horizontal position, or relative to the home or
display origin, only when the cursor position is wholely unknown.

Partial-display or log-leading renders are possible by postulating that the
current cursor position at the origin and drawing around it.

```go
cur := display.Reset
```

The cursor has methods to produce the commands that will show and hide the
cursor, clear the display, reset its state, seek to the origin, or move to
another cell's coordinates.

```go
var buf []byte
cur := display.Start
buf, cur = cur.Hide(buf)
buf, cur = cur.Clear(buf)
buf, cur = cur.Home(buf)
os.Stdout.Write(buf)
buf = buf[0:0]

cur, buf = cur.Go(buf, image.Pt(10, 20))

buf, cur = cur.Home(buf)
buf, cur = cur.Clear(buf)
buf, cur = cur.Show(buf)
os.Stdout.Write(buf)
buf = buf[0:0]
```

## colors

The `display` package provides the terminal colors, palettes, terminal
rendering color models.

- `Colors` is an array of 256 palette colors.
  - The first 8 are the 3-bit color palette.
  - The second 8 (8-15) are their bright versions from the 4-bit color palette.
  - The remaining colors form the 6x6x6 color cube and 24 grays scale.
- `Palette3`, `Palette4`, and `Palette8` are Go `"image".Palette` instances for
  colors expressible in those ranges.
- `Model0`, `Model3`, `Model4`, `Model8`, and `Model24` are virtual terminal
  color depth models with methods for rendering background and foreground
  colors to ANSI escape sequences, as used by `"display".Render`.
  `Model0` is monochrome and does not render color. `Model24` uses
  paletted colors only for exact matches.

## textile

Displays have a text image or "textile" of strings.
The `textile` package implements a text image, modeled after Go's `"image"`
package.

```go
text := textile.New(image.Rect(0, 0, 80, 23))
```

Just as with Go's images, the bounding box for the textile does not need
to start at the origin, and subtexts share the same memory.

```go
text.Subtext(image.Rect(1, 1, 79, 22)).Fill(".")
```

The default subtext is a matrix of nil strings.
Nil strings are "transparent". 
Drawing a blank text over another textile effects no change.

Each cell of the textile may be a string of arbitrary length, so it is possible
to model cells as sequences of UTF-8 including multiple code points with
joiners.

The display render function is sensitive to the uncertainty whether these
characters will necessarily be merged on the terminal display, invalidating the
cursor's horizontal position after rendering each cell that contains more than
one byte of text, then seeking to the next cell before rendering another.

## text

The `text` package provides a convenience for rendering plain text
onto a display.
The `Bounds(string)` method returns a bounding rectangle by measuring
how much space the string would need.
The `Write(*Display, Rectangle, string, Color)` method can then
write the string into a display.

```go
front, back := display.New2(bounds)

msg := "Hello, Cops!"
msgbox := text.Bounds(msg)
center := rectangle.MiddleCenter(msgbox, bounds)
text.Write(front, center, msg, display.Colors[7])
```

The text package executes only the smallest subset of the terminal language,
respecting "\n", "\t", and " ". Newline advances to the first column of the
next line. Tab and space advance the cursor without drawing, leaving
a transparent gap in the display.

Fill the bounds with " " or add a translucent or opaque background color to the
display to occlude holes left by transparent space.

## terminal

The `terminal` package is a thin wrapper around terminal control, to make
terminal capability changes and restoration easy and idiomatic to Go.

```go
term := terminal.New(os.Stdin.Fd())
defer term.Restore()
term.SetRaw()
term.SetNoEcho()
```

The `Bounds()` method returns an `"image".Rectangle` from the terminal size,
suitable for constructing a virtual display of the same size.

```go
bounds := term.Bounds()
front, back := display.New2(bounds)
```

## bitmap

The `bitmap` package provides a memory compact image type for images with only
two colors, as well as an image transformation layer that interprets another
image as a bitmap of the closer match of two colors.

## braille

The `braille` package draws bitmaps as matrices of braille dots.
See `cmd/braille/` for a demonstration.

# Tips / Tricks

## How to fill the background color for a text panel

```go
panel := text.Display("Press any key to continue...", display.Colors[7])
draw.Draw(panel.Background, panel.Bounds(), &image.Uniform{color.NRGBA{63, 63, 63, 128}}, image.ZP, draw.Over)
```

The background color can be translucent.
Drawing a display over another display:

- overrides the foreground color for all cells with text, except empty and
  space cells.
- draws the panel's background color over the foreground *and* background color
  of the underlying cell. This makes it possible to render translucent panels.
  The underlying text shows through, but faded by the overlying background
  color. 

```go
display.Draw(front, panel.Bounds(), panel, image.ZP, draw.Over)
```

---

Copyright 2017 Kristopher Michael Kowal and contributors.
Apache 2.0 license.
