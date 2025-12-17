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
	if ghost := i.findGhost(); ghost != "" {
		i.buffer = append(i.buffer, []rune(ghost)...)
		i.render()
	}
	return "", actionContinue
}

func handleEnter(i *Input, reader *bufio.Reader) (string, action) {
	fmt.Fprintln(i.out)
	return string(i.buffer), actionSubmit
}

func handleBackspace(i *Input, reader *bufio.Reader) (string, action) {
	if len(i.buffer) > 0 {
		i.buffer = i.buffer[:len(i.buffer)-1]
		i.render()
	}
	return "", actionContinue
}

func handleEscape(i *Input, reader *bufio.Reader) (string, action) {
	if reader.Buffered() > 0 {
		reader.ReadByte()
	}
	if reader.Buffered() > 0 {
		reader.ReadByte()
	}
	return "", actionContinue
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
