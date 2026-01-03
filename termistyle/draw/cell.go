// Package draw provides a 2D cell buffer for terminal rendering.
// Supports text, borders, and colored regions.
package draw

import "github.com/dottermi/x/termistyle/style"

// Cell represents a single character position in the terminal.
// Contains the character to display, colors, and text styling attributes.
type Cell struct {
	Char       rune
	Foreground style.Color
	Background style.Color
	Bold       bool
	Italic     bool
	Underline  bool
	Strike     bool
	Dim        bool
	Reverse    bool
}
