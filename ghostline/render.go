package ghostline

import "fmt"

// render displays the prompt, buffer, and ghost text with proper cursor positioning.
// Handles cursor placement at any position within the buffer.
func (i *Input) render() {
	// Clear line and write prompt + buffer
	_, _ = fmt.Fprintf(i.out, "\r\033[K%s%s", i.prompt, string(i.buffer))

	// Get ghost text (only shown when cursor is at end)
	ghost := ""
	if i.cursorPos == len(i.buffer) {
		ghost = i.findGhost()
	}

	// Show ghost text if available
	if ghost != "" {
		_, _ = fmt.Fprintf(i.out, "\033[2m%s\033[0m", ghost)
	}

	// Position cursor correctly
	moveBack := len(i.buffer) - i.cursorPos + len(ghost)
	if moveBack > 0 {
		_, _ = fmt.Fprintf(i.out, "\033[%dD", moveBack)
	}
}
