package render

import (
	"io"
	"strings"
)

// cellStyle tracks the current styling state for incremental ANSI output.
type cellStyle struct {
	FG, BG                                        Color
	Bold, Dim, Italic, Underline, Strike, Reverse bool
}

// writeAttr writes the on or off code based on current value.
func writeAttr(w io.Writer, val bool, on, off string) {
	if val {
		_, _ = io.WriteString(w, on)
	} else {
		_, _ = io.WriteString(w, off)
	}
}

// writeDiffColors writes ANSI codes for foreground and background color differences.
func writeDiffColors(w io.Writer, last *cellStyle, cell Cell, firstCell bool) {
	if firstCell || cell.FG != last.FG {
		_, _ = io.WriteString(w, cell.FG.FGCode())
		last.FG = cell.FG
	}
	if firstCell || cell.BG != last.BG {
		_, _ = io.WriteString(w, cell.BG.BGCode())
		last.BG = cell.BG
	}
}

// writeDiffAttrs writes ANSI codes for text attribute differences.
func writeDiffAttrs(w io.Writer, last *cellStyle, cell Cell, firstCell bool) {
	if firstCell || cell.Bold != last.Bold {
		writeAttr(w, cell.Bold, BoldOn, BoldOff)
		last.Bold = cell.Bold
	}
	if firstCell || cell.Dim != last.Dim {
		writeAttr(w, cell.Dim, DimOn, DimOff)
		last.Dim = cell.Dim
	}
	if firstCell || cell.Italic != last.Italic {
		writeAttr(w, cell.Italic, ItalicOn, ItalicOff)
		last.Italic = cell.Italic
	}
	if firstCell || cell.Underline != last.Underline {
		writeAttr(w, cell.Underline, UnderlineOn, UnderlineOff)
		last.Underline = cell.Underline
	}
	if firstCell || cell.Strike != last.Strike {
		writeAttr(w, cell.Strike, StrikeOn, StrikeOff)
		last.Strike = cell.Strike
	}
	if firstCell || cell.Reverse != last.Reverse {
		writeAttr(w, cell.Reverse, ReverseOn, ReverseOff)
		last.Reverse = cell.Reverse
	}
}

// writeDiff emits ANSI codes for the differences between last and cell styles.
// Updates last to match cell and returns it.
func writeDiff(w io.Writer, last cellStyle, cell Cell, firstCell bool) cellStyle {
	writeDiffColors(w, &last, cell, firstCell)
	writeDiffAttrs(w, &last, cell, firstCell)
	return last
}

// writeChar writes the cell's character, defaulting to space for zero value.
func writeChar(w io.Writer, char rune) {
	if char == 0 {
		_, _ = io.WriteString(w, " ")
	} else {
		_, _ = io.WriteString(w, string(char))
	}
}

// Terminal manages double-buffered rendering with diff-based optimization.
// It tracks the previously rendered state and outputs only the ANSI sequences
// needed to update changed cells, minimizing terminal I/O.
//
// Typical usage follows an immediate-mode pattern:
//
//	term := render.NewTerminal(80, 24)
//	for {
//		buf := render.NewBuffer(80, 24)
//		// Draw your UI into buf...
//		output := term.Render(buf)
//		os.Stdout.WriteString(output)
//	}
type Terminal struct {
	width   int
	height  int
	current *Buffer

	// Cursor tracking for optimization
	cursorX int
	cursorY int
}

// NewTerminal creates a Terminal renderer for the specified dimensions.
// The internal buffer is initialized to empty cells.
//
// Parameters:
//   - width: terminal width in columns
//   - height: terminal height in rows
//
// Example:
//
//	term := render.NewTerminal(80, 24)
func NewTerminal(width, height int) *Terminal {
	return &Terminal{
		width:   width,
		height:  height,
		current: NewBuffer(width, height),
	}
}

// Resize updates the terminal dimensions and clears the internal buffer.
// Call this when the terminal window size changes.
//
// Parameters:
//   - width: new terminal width in columns
//   - height: new terminal height in rows
func (t *Terminal) Resize(width, height int) {
	t.width = width
	t.height = height
	t.current = NewBuffer(width, height)
	t.cursorX = 0
	t.cursorY = 0
}

// Render computes the diff between the current state and new buffer,
// returning the minimal ANSI escape sequences needed to update the terminal.
// Updates the internal state to match new after rendering.
//
// Returns an empty string if no cells have changed.
//
// Example:
//
//	buf := render.NewBuffer(80, 24)
//	buf.Set(0, 0, render.Cell{Char: 'X'})
//	output := term.Render(buf)
//	os.Stdout.WriteString(output)
func (t *Terminal) Render(next *Buffer) string {
	var sb strings.Builder
	t.RenderTo(next, &sb)
	return sb.String()
}

// RenderTo writes the diff-optimized ANSI output directly to w.
// Prefer this over [Terminal.Render] when writing to files or network connections
// to avoid intermediate string allocation.
//
// Parameters:
//   - next: the buffer representing the desired screen state
//   - w: destination for ANSI output
func (t *Terminal) RenderTo(next *Buffer, w io.Writer) {
	// Handle dimension mismatch
	if next.Width != t.width || next.Height != t.height {
		t.Resize(next.Width, next.Height)
	}

	changes := t.current.Diff(next)
	if len(changes) == 0 {
		return
	}

	var style cellStyle
	firstCell := true

	for _, change := range changes {
		// Move cursor if not at expected position
		if change.X != t.cursorX || change.Y != t.cursorY {
			_, _ = io.WriteString(w, MoveCursor(change.X, change.Y))
			t.cursorX = change.X
			t.cursorY = change.Y
		}

		style = writeDiff(w, style, change.Cell, firstCell)
		firstCell = false
		writeChar(w, change.Cell.Char)

		// Update cursor position
		t.cursorX++
		if t.cursorX >= t.width {
			t.cursorX = 0
			t.cursorY++
		}
	}

	_, _ = io.WriteString(w, Reset)
	t.current = next.Clone()
}

// RenderFull renders the entire buffer without diff computation.
// Use this for the initial frame or after clearing the screen.
// Subsequent frames should use [Terminal.Render] for efficiency.
//
// Example:
//
//	fmt.Print(render.ClearScreen)
//	fmt.Print(term.RenderFull(buf))
func (t *Terminal) RenderFull(buf *Buffer) string {
	var sb strings.Builder
	t.RenderFullTo(buf, &sb)
	return sb.String()
}

// RenderFullTo writes the entire buffer to w without diff computation.
// Prefer this over [Terminal.RenderFull] to avoid intermediate string allocation.
//
// Parameters:
//   - buf: the buffer to render
//   - w: destination for ANSI output
func (t *Terminal) RenderFullTo(buf *Buffer, w io.Writer) {
	_, _ = io.WriteString(w, MoveCursor(0, 0))

	var style cellStyle
	firstCell := true

	for y := range buf.Height {
		for x := range buf.Width {
			cell := buf.Get(x, y)
			style = writeDiff(w, style, cell, firstCell)
			firstCell = false
			writeChar(w, cell.Char)
		}

		// Newline between rows (except last)
		if y < buf.Height-1 {
			_, _ = io.WriteString(w, "\n")
		}
	}

	_, _ = io.WriteString(w, Reset)
	t.current = buf.Clone()
	t.cursorX = buf.Width - 1
	t.cursorY = buf.Height - 1
}

// Clear resets the internal buffer and returns ANSI codes to clear the screen.
// The cursor is moved to the home position (0, 0).
//
// Example:
//
//	os.Stdout.WriteString(term.Clear())
func (t *Terminal) Clear() string {
	t.current = NewBuffer(t.width, t.height)
	t.cursorX = 0
	t.cursorY = 0
	return ClearScreen + CursorHome
}
