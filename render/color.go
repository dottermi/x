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

import "fmt"

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
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
}

// BGCode returns the ANSI escape sequence for setting this color as background.
// Returns the default background reset code if the color is not set.
func (c Color) BGCode() string {
	if !c.Set {
		return "\x1b[49m"
	}
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
}
