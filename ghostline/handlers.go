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
	keyCtrlJ     = 10  // Insert newline (multiline editing)
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
	actionInterrupted
	actionEOF
)

var actionErrors = map[action]error{
	actionSubmit:      nil,
	actionInterrupted: ErrInterrupted,
	actionEOF:         ErrEOF,
}

type keyHandler func(i *Input, reader *bufio.Reader) (string, action)

func handleCtrlC(i *Input, reader *bufio.Reader) (string, action) {
	i.buffer = nil
	i.render()
	return "", actionInterrupted
}

func handleCtrlD(i *Input, reader *bufio.Reader) (string, action) {
	if len(i.buffer) == 0 {
		return "", actionEOF
	}
	return "", actionContinue
}

func handleTab(i *Input, reader *bufio.Reader) (string, action) {
	if i.cursorPos == len(i.buffer) {
		if match := i.findMatch(); match != "" {
			// Replace last word with the full match (preserves suggestion's case)
			start := i.lastWordStart()
			i.buffer = append(i.buffer[:start], []rune(match)...)
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
	'A': handleUpArrow,
	'B': handleDownArrow,
	'C': handleRightArrow,
	'D': handleLeftArrow,
	'3': handleDelete,
}

func (i *Input) handleCSI(code byte, reader *bufio.Reader) {
	if handler, ok := csiHandlers[code]; ok {
		handler(i, reader)
	}
}

func handleUpArrow(i *Input, _ *bufio.Reader) {
	// Find current line start and check if we're on first line
	lineStart := i.findLineStart()
	if lineStart == 0 {
		// On first line, navigate history
		entry, ok := i.history.Previous(string(i.buffer))
		if ok {
			i.buffer = []rune(entry)
			i.cursorPos = len(i.buffer)
			i.render()
		}
		return
	}

	// Move to previous line at same column
	col := i.cursorPos - lineStart
	prevLineEnd := lineStart - 1
	prevLineStart := i.findLineStartFrom(prevLineEnd)
	prevLineLen := prevLineEnd - prevLineStart

	if col > prevLineLen {
		i.cursorPos = prevLineEnd
	} else {
		i.cursorPos = prevLineStart + col
	}
	i.render()
}

func handleDownArrow(i *Input, _ *bufio.Reader) {
	// Find current line end and check if we're on last line
	lineEnd := i.findLineEnd()
	if lineEnd == len(i.buffer) {
		// On last line, navigate history
		entry, ok := i.history.Next()
		if ok {
			i.buffer = []rune(entry)
			i.cursorPos = len(i.buffer)
			i.render()
		}
		return
	}

	// Move to next line at same column
	lineStart := i.findLineStart()
	col := i.cursorPos - lineStart
	nextLineStart := lineEnd + 1
	nextLineEnd := i.findLineEndFrom(nextLineStart)
	nextLineLen := nextLineEnd - nextLineStart

	if col > nextLineLen {
		i.cursorPos = nextLineEnd
	} else {
		i.cursorPos = nextLineStart + col
	}
	i.render()
}

func (i *Input) findLineStart() int {
	return i.findLineStartFrom(i.cursorPos)
}

func (i *Input) findLineStartFrom(pos int) int {
	for pos > 0 && i.buffer[pos-1] != '\n' {
		pos--
	}
	return pos
}

func (i *Input) findLineEnd() int {
	return i.findLineEndFrom(i.cursorPos)
}

func (i *Input) findLineEndFrom(pos int) int {
	for pos < len(i.buffer) && i.buffer[pos] != '\n' {
		pos++
	}
	return pos
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
	// Move to beginning of current line
	i.cursorPos = i.findLineStart()
	i.render()
	return "", actionContinue
}

func handleCtrlE(i *Input, reader *bufio.Reader) (string, action) {
	// Move to end of current line
	i.cursorPos = i.findLineEnd()
	i.render()
	return "", actionContinue
}

func handleCtrlK(i *Input, reader *bufio.Reader) (string, action) {
	// Kill from cursor to end of current line
	lineEnd := i.findLineEnd()
	i.buffer = append(i.buffer[:i.cursorPos], i.buffer[lineEnd:]...)
	i.render()
	return "", actionContinue
}

func handleCtrlU(i *Input, reader *bufio.Reader) (string, action) {
	// Kill from beginning of current line to cursor
	lineStart := i.findLineStart()
	i.buffer = append(i.buffer[:lineStart], i.buffer[i.cursorPos:]...)
	i.cursorPos = lineStart
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

func handleCtrlJ(i *Input, reader *bufio.Reader) (string, action) {
	// Insert newline at cursor position
	i.buffer = append(i.buffer[:i.cursorPos], append([]rune{'\n'}, i.buffer[i.cursorPos:]...)...)
	i.cursorPos++
	i.render()
	return "", actionContinue
}

func defaultHandlers() map[rune]keyHandler {
	return map[rune]keyHandler{
		keyCtrlA:     handleCtrlA,
		keyCtrlC:     handleCtrlC,
		keyCtrlD:     handleCtrlD,
		keyCtrlE:     handleCtrlE,
		keyCtrlJ:     handleCtrlJ,
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
