// Package render provides efficient terminal rendering with double-buffering and diff optimization.
// Inspired by the Ratatui immediate-mode rendering pattern, it minimizes ANSI escape code output
// by only updating cells that have changed between frames.
//
// The package centers around four types:
//   - [Color]: RGB colors with optional terminal default fallback
//   - [Cell]: A single terminal position with character and styling
//   - [Buffer]: A 2D grid of cells representing the screen
//   - [Terminal]: Double-buffered renderer that computes minimal diffs
//
// Basic usage:
//
//	term := render.NewTerminal(80, 24)
//	buf := render.NewBuffer(80, 24)
//
//	// Draw to buffer
//	cell := render.Cell{Char: 'X', FG: render.RGB(255, 0, 0)}
//	buf.Set(10, 5, cell)
//
//	// Render only changed cells
//	output := term.Render(buf)
//	fmt.Print(output)
package render

import (
	"io"
	"strconv"
)

// Color represents an RGB color for terminal rendering with optional default fallback.
// Use [RGB] to create explicit colors or [Default] for terminal defaults.
// The zero value is equivalent to Default().
//
// Example:
//
//	red := render.RGB(255, 0, 0)
//	transparent := render.Default()
type Color struct {
	R, G, B uint8
	Set     bool
}

// RGB creates a Color with the specified red, green, and blue components.
// Values range from 0 to 255.
//
// Example:
//
//	white := render.RGB(255, 255, 255)
//	coral := render.RGB(255, 127, 80)
func RGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b, Set: true}
}

// Default returns a Color that uses the terminal's default foreground or background.
// Equivalent to the zero value.
func Default() Color {
	return Color{}
}

// IsSet reports whether the color was explicitly set via [RGB].
// Returns false for colors created with [Default] or zero values.
func (c Color) IsSet() bool {
	return c.Set
}

// FGCode returns the ANSI escape sequence for setting this color as foreground.
// Returns the default foreground reset code if the color is not set.
func (c Color) FGCode() string {
	if !c.Set {
		return "\x1b[39m"
	}
	return c.buildCode("\x1b[38;2;")
}

// BGCode returns the ANSI escape sequence for setting this color as background.
// Returns the default background reset code if the color is not set.
func (c Color) BGCode() string {
	if !c.Set {
		return "\x1b[49m"
	}
	return c.buildCode("\x1b[48;2;")
}

// buildCode constructs an ANSI color code with the given prefix.
func (c Color) buildCode(prefix string) string {
	// "\x1b[38;2;255;255;255m" = 19 bytes max
	buf := make([]byte, 0, 19)
	buf = append(buf, prefix...)
	buf = strconv.AppendUint(buf, uint64(c.R), 10)
	buf = append(buf, ';')
	buf = strconv.AppendUint(buf, uint64(c.G), 10)
	buf = append(buf, ';')
	buf = strconv.AppendUint(buf, uint64(c.B), 10)
	buf = append(buf, 'm')
	return string(buf)
}

// WriteFGTo writes the foreground ANSI escape sequence directly to w.
// This avoids string allocation compared to FGCode().
func (c Color) WriteFGTo(w io.Writer) {
	if !c.Set {
		_, _ = io.WriteString(w, "\x1b[39m")
		return
	}
	c.writeCodeTo(w, "\x1b[38;2;")
}

// WriteBGTo writes the background ANSI escape sequence directly to w.
// This avoids string allocation compared to BGCode().
func (c Color) WriteBGTo(w io.Writer) {
	if !c.Set {
		_, _ = io.WriteString(w, "\x1b[49m")
		return
	}
	c.writeCodeTo(w, "\x1b[48;2;")
}

// writeCodeTo writes an ANSI color code with the given prefix to w.
func (c Color) writeCodeTo(w io.Writer, prefix string) {
	var buf [19]byte
	b := buf[:0]
	b = append(b, prefix...)
	b = strconv.AppendUint(b, uint64(c.R), 10)
	b = append(b, ';')
	b = strconv.AppendUint(b, uint64(c.G), 10)
	b = append(b, ';')
	b = strconv.AppendUint(b, uint64(c.B), 10)
	b = append(b, 'm')
	_, _ = w.Write(b)
}
