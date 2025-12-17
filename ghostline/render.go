package ghostline

import "fmt"

// findGhost returns the completion suffix for the last word in the buffer.
// Returns empty string if no matching suggestion exists or buffer ends with space.
func (i *Input) render() {
	_, _ = fmt.Fprintf(i.out, "\r\033[K%s%s", i.prompt, string(i.buffer))

	ghost := i.findGhost()
	if ghost == "" {
		return
	}

	_, _ = fmt.Fprintf(i.out, "\033[2m%s\033[0m", ghost)
	_, _ = fmt.Fprintf(i.out, "\033[%dD", len(ghost))
}
