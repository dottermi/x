package render

import (
	"io"
	"strings"
)

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
func (t *Terminal) Render(new *Buffer) string {
	var sb strings.Builder
	t.RenderTo(new, &sb)
	return sb.String()
}

// RenderTo writes the diff-optimized ANSI output directly to w.
// Prefer this over [Terminal.Render] when writing to files or network connections
// to avoid intermediate string allocation.
//
// Parameters:
//   - new: the buffer representing the desired screen state
//   - w: destination for ANSI output
func (t *Terminal) RenderTo(new *Buffer, w io.Writer) {
	// Handle dimension mismatch
	if new.Width != t.width || new.Height != t.height {
		t.Resize(new.Width, new.Height)
	}

	changes := t.current.Diff(new)
	if len(changes) == 0 {
		return
	}

	var lastFG, lastBG Color
	var lastBold, lastDim, lastItalic, lastUnderline, lastStrike, lastReverse bool
	firstCell := true

	for _, change := range changes {
		// Move cursor if not at expected position
		if change.X != t.cursorX || change.Y != t.cursorY {
			io.WriteString(w, MoveCursor(change.X, change.Y))
			t.cursorX = change.X
			t.cursorY = change.Y
		}

		cell := change.Cell

		// Emit color codes only when changed
		if firstCell || cell.FG != lastFG {
			io.WriteString(w, cell.FG.FGCode())
			lastFG = cell.FG
		}
		if firstCell || cell.BG != lastBG {
			io.WriteString(w, cell.BG.BGCode())
			lastBG = cell.BG
		}

		// Emit attribute codes only when changed
		if firstCell || cell.Bold != lastBold {
			if cell.Bold {
				io.WriteString(w, BoldOn)
			} else {
				io.WriteString(w, BoldOff)
			}
			lastBold = cell.Bold
		}
		if firstCell || cell.Dim != lastDim {
			if cell.Dim {
				io.WriteString(w, DimOn)
			} else {
				io.WriteString(w, DimOff)
			}
			lastDim = cell.Dim
		}
		if firstCell || cell.Italic != lastItalic {
			if cell.Italic {
				io.WriteString(w, ItalicOn)
			} else {
				io.WriteString(w, ItalicOff)
			}
			lastItalic = cell.Italic
		}
		if firstCell || cell.Underline != lastUnderline {
			if cell.Underline {
				io.WriteString(w, UnderlineOn)
			} else {
				io.WriteString(w, UnderlineOff)
			}
			lastUnderline = cell.Underline
		}
		if firstCell || cell.Strike != lastStrike {
			if cell.Strike {
				io.WriteString(w, StrikeOn)
			} else {
				io.WriteString(w, StrikeOff)
			}
			lastStrike = cell.Strike
		}
		if firstCell || cell.Reverse != lastReverse {
			if cell.Reverse {
				io.WriteString(w, ReverseOn)
			} else {
				io.WriteString(w, ReverseOff)
			}
			lastReverse = cell.Reverse
		}

		firstCell = false

		// Write character
		if cell.Char == 0 {
			io.WriteString(w, " ")
		} else {
			io.WriteString(w, string(cell.Char))
		}

		// Update cursor position
		t.cursorX++
		if t.cursorX >= t.width {
			t.cursorX = 0
			t.cursorY++
		}
	}

	// Reset at end
	io.WriteString(w, Reset)

	// Update current buffer
	t.current = new.Clone()
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
	io.WriteString(w, MoveCursor(0, 0))

	var lastFG, lastBG Color
	var lastBold, lastDim, lastItalic, lastUnderline, lastStrike, lastReverse bool
	firstCell := true

	for y := 0; y < buf.Height; y++ {
		for x := 0; x < buf.Width; x++ {
			cell := buf.Get(x, y)

			// Emit color codes only when changed
			if firstCell || cell.FG != lastFG {
				io.WriteString(w, cell.FG.FGCode())
				lastFG = cell.FG
			}
			if firstCell || cell.BG != lastBG {
				io.WriteString(w, cell.BG.BGCode())
				lastBG = cell.BG
			}

			// Emit attribute codes only when changed
			if firstCell || cell.Bold != lastBold {
				if cell.Bold {
					io.WriteString(w, BoldOn)
				} else if !firstCell {
					io.WriteString(w, BoldOff)
				}
				lastBold = cell.Bold
			}
			if firstCell || cell.Dim != lastDim {
				if cell.Dim {
					io.WriteString(w, DimOn)
				} else if !firstCell {
					io.WriteString(w, DimOff)
				}
				lastDim = cell.Dim
			}
			if firstCell || cell.Italic != lastItalic {
				if cell.Italic {
					io.WriteString(w, ItalicOn)
				} else if !firstCell {
					io.WriteString(w, ItalicOff)
				}
				lastItalic = cell.Italic
			}
			if firstCell || cell.Underline != lastUnderline {
				if cell.Underline {
					io.WriteString(w, UnderlineOn)
				} else if !firstCell {
					io.WriteString(w, UnderlineOff)
				}
				lastUnderline = cell.Underline
			}
			if firstCell || cell.Strike != lastStrike {
				if cell.Strike {
					io.WriteString(w, StrikeOn)
				} else if !firstCell {
					io.WriteString(w, StrikeOff)
				}
				lastStrike = cell.Strike
			}
			if firstCell || cell.Reverse != lastReverse {
				if cell.Reverse {
					io.WriteString(w, ReverseOn)
				} else if !firstCell {
					io.WriteString(w, ReverseOff)
				}
				lastReverse = cell.Reverse
			}

			firstCell = false

			// Write character
			if cell.Char == 0 {
				io.WriteString(w, " ")
			} else {
				io.WriteString(w, string(cell.Char))
			}
		}

		// Newline between rows (except last)
		if y < buf.Height-1 {
			io.WriteString(w, "\n")
		}
	}

	io.WriteString(w, Reset)

	// Update state
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
