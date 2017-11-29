// Package cops contains packages for rendering terminal user interfaces.
// Cops supports 24 bit color and pairs well with the "color", "image", and
// "image/draw" packages.
//
// The "display" package models a display as foreground, background, and text
// layers. The display package provides Draw for composing display layers and
// Render for producing ANSI escape sequences to differentially update a
// terminal.
//
// The "cursor" package models cursor state and has methods to modify color,
// position, and visibility, appending the corresponding ANSI escape sequences
// to a write buffer.
//
// The "textile" package implements a text layer, like Go's own "image"
// package.
//
// The "text" package measures and cuts text onto a display.
//
// The "terminal" package provides an idiomatic Go interface for terminal
// capabilities ("raw mode", "no echo", getting and setting size).
//
// The "vtcolor" package supplements "color" with virtual terminal color
// palettes and rendering models for 0, 3, 4, 8, and 24 bit color.
//
// The "rectangle" package provides conveniences for manipulating image
// rectangles for display composition.
//
// The "bitmap" package provides a compact representation of bitmap images,
// suitable for use as masks or sources for braille bitmap displays.
//
// The "braille" package draws bitmap images onto displays as a matrix of
// braille text.
package cops
