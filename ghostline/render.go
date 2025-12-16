package ghostline

import "fmt"

// render redraws the entire input line with prompt, buffer, and ghost text.
// Uses ANSI escape codes to clear the line and position the cursor correctly.
func (i *Input) render() {
	fmt.Fprintf(i.out, "\r\033[K%s%s", i.prompt, string(i.buffer))

	ghost := i.findGhost()
	if ghost != "" {
		// Render ghost text in dim color (ANSI code \033[2m)
		fmt.Fprintf(i.out, "\033[2m%s\033[0m", ghost)

		// Move cursor back to the end of the buffer (before ghost)
		if len(ghost) > 0 {
			fmt.Fprintf(i.out, "\033[%dD", len(ghost))
		}
	}
}
