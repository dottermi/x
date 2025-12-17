package ghostline

import (
	"fmt"
	"strings"
)

// render displays the prompt, buffer, and ghost text with proper cursor positioning.
// Handles multiline input with continuation prompts.
func (i *Input) render() {
	// Move cursor up to first line if we had multiple lines before
	if i.prevLines > 1 {
		_, _ = fmt.Fprintf(i.out, "\033[%dA", i.prevLines-1)
	}

	// Clear from cursor to end of screen
	_, _ = fmt.Fprintf(i.out, "\r\033[J")

	// Split buffer into lines
	lines := strings.Split(string(i.buffer), "\n")
	i.prevLines = len(lines)

	// Find cursor line and column
	cursorLine, cursorCol := i.getCursorPosition()

	// Render each line
	for idx, line := range lines {
		if idx > 0 {
			_, _ = fmt.Fprint(i.out, "\n\r") // \n moves down, \r goes to column 0
		}
		if idx == 0 {
			_, _ = fmt.Fprintf(i.out, "%s%s", i.prompt, line)
		} else {
			_, _ = fmt.Fprintf(i.out, "%s%s", i.contPrompt, line)
		}
	}

	// Get ghost text (only shown when cursor is at end of last line)
	ghost := ""
	if i.cursorPos == len(i.buffer) {
		ghost = i.findGhost()
	}

	// Show ghost text if available
	if ghost != "" {
		_, _ = fmt.Fprintf(i.out, "\033[2m%s\033[0m", ghost)
	}

	// Position cursor correctly
	// Move to the cursor line (from the last line)
	linesFromEnd := len(lines) - 1 - cursorLine
	if linesFromEnd > 0 {
		_, _ = fmt.Fprintf(i.out, "\033[%dA", linesFromEnd)
	}

	// Move to correct column using absolute positioning
	prompt := i.prompt
	if cursorLine > 0 {
		prompt = i.contPrompt
	}
	col := len(prompt) + cursorCol
	_, _ = fmt.Fprintf(i.out, "\033[%dG", col+1) // \033[nG is 1-indexed
}

// getCursorPosition returns the line number and column of the cursor.
func (i *Input) getCursorPosition() (line, col int) {
	for _, r := range i.buffer[:i.cursorPos] {
		if r == '\n' {
			line++
			col = 0
		} else {
			col++
		}
	}
	return line, col
}
