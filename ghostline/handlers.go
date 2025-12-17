package ghostline

import (
	"bufio"
	"fmt"
)

// Key constants for terminal control characters.
// These values correspond to ASCII control codes sent by the terminal.
const (
	keyCtrlA     = 1   // Move cursor to beginning of line
	keyCtrlC     = 3   // Interrupt signal (abort input)
	keyCtrlD     = 4   // End of transmission (abort on empty buffer)
	keyCtrlE     = 5   // Move cursor to end of line
	keyTab       = 9   // Horizontal tab (accept suggestion)
	keyCtrlK     = 11  // Kill text from cursor to end of line
	keyEnter     = 13  // Carriage return (submit input)
	keyCtrlU     = 21  // Kill text from beginning to cursor
	keyCtrlW     = 23  // Delete word backward
	keyEscape    = 27  // Escape (start of escape sequences)
	keyBackspace = 127 // Delete previous character (DEL)
	keyDelete    = 8   // Backspace control character (BS)
)

type action int

const (
	actionContinue action = iota
	actionSubmit
	actionAbort
)

type keyHandler func(i *Input, reader *bufio.Reader) (string, action)

func handleCtrlC(i *Input, reader *bufio.Reader) (string, action) {
	i.buffer = nil
	i.render()
	return "", actionAbort
}

func handleCtrlD(i *Input, reader *bufio.Reader) (string, action) {
	if len(i.buffer) == 0 {
		return "", actionAbort
	}
	return "", actionContinue
}

func handleTab(i *Input, reader *bufio.Reader) (string, action) {
	if i.cursorPos == len(i.buffer) {
		if ghost := i.findGhost(); ghost != "" {
			i.buffer = append(i.buffer, []rune(ghost)...)
			i.cursorPos = len(i.buffer)
			i.render()
		}
	}
	return "", actionContinue
}

func handleEnter(i *Input, reader *bufio.Reader) (string, action) {
	_, _ = fmt.Fprintln(i.out)
	return string(i.buffer), actionSubmit
}

func handleBackspace(i *Input, reader *bufio.Reader) (string, action) {
	if i.cursorPos > 0 {
		i.buffer = append(i.buffer[:i.cursorPos-1], i.buffer[i.cursorPos:]...)
		i.cursorPos--
		i.render()
	}
	return "", actionContinue
}

func handleEscape(i *Input, reader *bufio.Reader) (string, action) {
	// Read escape sequence: ESC [ <code>
	b1, err := reader.ReadByte()
	if err != nil || b1 != '[' {
		return "", actionContinue
	}

	b2, err := reader.ReadByte()
	if err != nil {
		return "", actionContinue
	}

	i.handleCSI(b2, reader)
	return "", actionContinue
}

type csiHandler func(i *Input, reader *bufio.Reader)

var csiHandlers = map[byte]csiHandler{
	'C': handleRightArrow,
	'D': handleLeftArrow,
	'H': handleHome,
	'F': handleEnd,
	'3': handleDelete,
}

func (i *Input) handleCSI(code byte, reader *bufio.Reader) {
	if handler, ok := csiHandlers[code]; ok {
		handler(i, reader)
	}
}

func handleRightArrow(i *Input, _ *bufio.Reader) {
	if i.cursorPos < len(i.buffer) {
		i.cursorPos++
		i.render()
	}
}

func handleLeftArrow(i *Input, _ *bufio.Reader) {
	if i.cursorPos > 0 {
		i.cursorPos--
		i.render()
	}
}

func handleHome(i *Input, _ *bufio.Reader) {
	i.cursorPos = 0
	i.render()
}

func handleEnd(i *Input, _ *bufio.Reader) {
	i.cursorPos = len(i.buffer)
	i.render()
}

func handleDelete(i *Input, reader *bufio.Reader) {
	i.handleDeleteKey(reader)
}

func (i *Input) handleDeleteKey(reader *bufio.Reader) {
	b3, err := reader.ReadByte()
	if err != nil || b3 != '~' {
		return
	}
	if i.cursorPos < len(i.buffer) {
		i.buffer = append(i.buffer[:i.cursorPos], i.buffer[i.cursorPos+1:]...)
		i.render()
	}
}

func handleCtrlA(i *Input, reader *bufio.Reader) (string, action) {
	i.cursorPos = 0
	i.render()
	return "", actionContinue
}

func handleCtrlE(i *Input, reader *bufio.Reader) (string, action) {
	i.cursorPos = len(i.buffer)
	i.render()
	return "", actionContinue
}

func handleCtrlK(i *Input, reader *bufio.Reader) (string, action) {
	// Kill from cursor to end of line
	i.buffer = i.buffer[:i.cursorPos]
	i.render()
	return "", actionContinue
}

func handleCtrlU(i *Input, reader *bufio.Reader) (string, action) {
	// Kill from beginning to cursor
	i.buffer = i.buffer[i.cursorPos:]
	i.cursorPos = 0
	i.render()
	return "", actionContinue
}

func handleCtrlW(i *Input, reader *bufio.Reader) (string, action) {
	// Delete word backward
	if i.cursorPos == 0 {
		return "", actionContinue
	}

	// Skip trailing spaces
	pos := i.cursorPos
	for pos > 0 && i.buffer[pos-1] == ' ' {
		pos--
	}

	// Delete until next space or beginning
	for pos > 0 && i.buffer[pos-1] != ' ' {
		pos--
	}

	i.buffer = append(i.buffer[:pos], i.buffer[i.cursorPos:]...)
	i.cursorPos = pos
	i.render()
	return "", actionContinue
}

func defaultHandlers() map[rune]keyHandler {
	return map[rune]keyHandler{
		keyCtrlA:     handleCtrlA,
		keyCtrlC:     handleCtrlC,
		keyCtrlD:     handleCtrlD,
		keyCtrlE:     handleCtrlE,
		keyCtrlK:     handleCtrlK,
		keyCtrlU:     handleCtrlU,
		keyCtrlW:     handleCtrlW,
		keyTab:       handleTab,
		keyEnter:     handleEnter,
		keyBackspace: handleBackspace,
		keyDelete:    handleBackspace,
		keyEscape:    handleEscape,
	}
}
