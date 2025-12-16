package ghostline

import (
	"bufio"
	"fmt"
)

// Key constants for terminal control characters.
// These values correspond to ASCII control codes sent by the terminal.
const (
	keyCtrlC     = 3   // Interrupt signal (abort input)
	keyCtrlD     = 4   // End of transmission (abort on empty buffer)
	keyTab       = 9   // Horizontal tab (accept suggestion)
	keyEnter     = 13  // Carriage return (submit input)
	keyEscape    = 27  // Escape (start of escape sequences)
	keyBackspace = 127 // Delete previous character (DEL)
	keyDelete    = 8   // Backspace control character (BS)
)

// keyHandler processes a key press and returns the result.
// Returns: (result string, done bool, ok bool)
//   - result: the text to return if done
//   - done: whether to exit the input loop
//   - ok: whether input was successful (true) or aborted (false)
type keyHandler func(i *Input, reader *bufio.Reader) (string, bool, bool)

func handleCtrlC(i *Input, reader *bufio.Reader) (string, bool, bool) {
	i.buffer = nil
	i.render()
	return "", true, false
}

func handleCtrlD(i *Input, reader *bufio.Reader) (string, bool, bool) {
	if len(i.buffer) == 0 {
		return "", true, false
	}
	return "", false, false
}

func handleTab(i *Input, reader *bufio.Reader) (string, bool, bool) {
	ghost := i.findGhost()
	if ghost != "" {
		i.buffer = append(i.buffer, []rune(ghost)...)
		i.render()
	}
	return "", false, false
}

func handleEnter(i *Input, reader *bufio.Reader) (string, bool, bool) {
	fmt.Fprintln(i.out)
	return string(i.buffer), true, true
}

func handleBackspace(i *Input, reader *bufio.Reader) (string, bool, bool) {
	if len(i.buffer) > 0 {
		i.buffer = i.buffer[:len(i.buffer)-1]
		i.render()
	}
	return "", false, false
}

func handleEscape(i *Input, reader *bufio.Reader) (string, bool, bool) {
	if reader.Buffered() > 0 {
		reader.ReadByte()
	}
	if reader.Buffered() > 0 {
		reader.ReadByte()
	}
	return "", false, false
}

func defaultHandlers() map[rune]keyHandler {
	return map[rune]keyHandler{
		keyCtrlC:     handleCtrlC,
		keyCtrlD:     handleCtrlD,
		keyTab:       handleTab,
		keyEnter:     handleEnter,
		keyBackspace: handleBackspace,
		keyDelete:    handleBackspace,
		keyEscape:    handleEscape,
	}
}
