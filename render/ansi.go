package render

import "fmt"

// ANSI escape sequences for terminal control and text styling.
// These constants can be written directly to the terminal output.
//
// Example:
//
//	fmt.Print(render.ClearScreen + render.CursorHome)
//	fmt.Print(render.BoldOn + "Bold text" + render.Reset)
const (
	// Reset clears all text attributes and colors.
	Reset = "\x1b[0m"

	// ClearScreen clears the entire screen.
	ClearScreen = "\x1b[2J"
	ClearLine   = "\x1b[2K" // Clear current line
	CursorHome  = "\x1b[H"  // Move cursor to top-left (0, 0)
	CursorHide  = "\x1b[?25l"
	CursorShow  = "\x1b[?25h"

	// AltScreenEnter switches to alternate screen buffer for fullscreen apps (like vim, htop).
	AltScreenEnter = "\x1b[?1049h"
	AltScreenLeave = "\x1b[?1049l"

	// BoldOn enables bold text attribute.
	BoldOn      = "\x1b[1m"
	DimOn       = "\x1b[2m"
	ItalicOn    = "\x1b[3m"
	UnderlineOn = "\x1b[4m"
	ReverseOn   = "\x1b[7m" // Swap foreground and background
	StrikeOn    = "\x1b[9m"

	// BoldOff disables bold text attribute (also disables Dim).
	BoldOff      = "\x1b[22m"
	DimOff       = "\x1b[22m" // Also disables Bold
	ItalicOff    = "\x1b[23m"
	UnderlineOff = "\x1b[24m"
	ReverseOff   = "\x1b[27m"
	StrikeOff    = "\x1b[29m"
)

// MoveCursor returns the ANSI escape sequence to position the cursor.
// Coordinates are 0-indexed; the function handles conversion to 1-indexed ANSI positions.
//
// Parameters:
//   - x: column (0-indexed)
//   - y: row (0-indexed)
//
// Example:
//
//	fmt.Print(render.MoveCursor(10, 5)) // Move to column 10, row 5
func MoveCursor(x, y int) string {
	return fmt.Sprintf("\x1b[%d;%dH", y+1, x+1)
}

// HideCursor returns the ANSI escape sequence to hide the cursor.
// Pair with [ShowCursor] when rendering is complete.
func HideCursor() string {
	return CursorHide
}

// ShowCursor returns the ANSI escape sequence to show the cursor.
// Call this to restore cursor visibility after [HideCursor].
func ShowCursor() string {
	return CursorShow
}
