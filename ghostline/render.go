package ghostline

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

// ansiRegex matches ANSI escape sequences.
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// visibleWidth returns the display width of a string, ignoring ANSI codes.
func visibleWidth(s string) int {
	return runewidth.StringWidth(ansiRegex.ReplaceAllString(s, ""))
}

// render redraws the entire input area including prompt, text, and ghost suggestions.
// Clears previous output, displays each line with appropriate prompts, and positions
// the cursor correctly for multiline editing.
//
// Uses ANSI escape sequences for terminal control:
//   - ESC[nA: move cursor up n lines
//   - ESC[J:  clear from cursor to end of screen
//   - ESC[nG: move cursor to column n
//   - ESC[38;2;r;g;bm: set 24-bit foreground color
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

	// Show ghost text and dropdown (only when cursor is at end of last line)
	if i.cursorPos == len(i.buffer) {
		// Ghost text (color: #6b7280)
		if ghost := i.findGhost(); ghost != "" {
			_, _ = fmt.Fprintf(i.out, "\033[38;2;107;114;128m%s\033[0m", ghost)
		}
		// Dropdown hints
		i.renderDropdown()
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
	col := visibleWidth(prompt) + cursorCol
	_, _ = fmt.Fprintf(i.out, "\r\033[%dG", col+1) // \r to start, \033[nG is 1-indexed
}

// getCursorPosition returns the cursor's line and column within the buffer.
// Line is 0-indexed. Column is measured in display width (handles wide characters).
func (i *Input) getCursorPosition() (line, col int) {
	for _, r := range i.buffer[:i.cursorPos] {
		if r == '\n' {
			line++
			col = 0
		} else {
			col += runewidth.RuneWidth(r)
		}
	}
	return line, col
}
