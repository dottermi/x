package render

// Cell represents a single character position in the terminal with styling attributes.
// Each cell contains a character, foreground and background colors, and text decorations.
// The zero value displays as a space with terminal default colors.
//
// Example:
//
//	cell := render.Cell{
//		Char:   'A',
//		FG:     render.RGB(255, 255, 255),
//		BG:     render.RGB(0, 0, 128),
//		Bold:   true,
//	}
type Cell struct {
	// Char is the Unicode character to display. Zero value renders as space.
	Char rune

	// FG is the foreground (text) color.
	FG Color

	// BG is the background color.
	BG Color

	// Text decoration attributes
	Bold      bool
	Dim       bool
	Italic    bool
	Underline bool
	Strike    bool
	Reverse   bool
}

// EmptyCell returns a Cell containing a space character with default colors.
// Use this as the baseline for cleared or empty screen regions.
func EmptyCell() Cell {
	return Cell{Char: ' '}
}

// Equal reports whether c and other have identical character, colors, and attributes.
// Used internally by [Buffer.Diff] to detect changes between frames.
func (c Cell) Equal(other Cell) bool {
	return c == other
}
