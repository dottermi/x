package ghostline

import (
	"bufio"
	"fmt"
)

// Terminal control character constants (ASCII codes).
// These correspond to the raw bytes sent when keys are pressed in raw mode.
const (
	keyCtrlA     = 1   // Ctrl+A: move to line start
	keyCtrlC     = 3   // Ctrl+C: interrupt/abort
	keyCtrlD     = 4   // Ctrl+D: EOF on empty buffer
	keyCtrlE     = 5   // Ctrl+E: move to line end
	keyTab       = 9   // Tab: accept suggestion
	keyCtrlJ     = 10  // Ctrl+J: insert newline
	keyCtrlK     = 11  // Ctrl+K: kill to end of line
	keyEnter     = 13  // Enter: submit input
	keyCtrlU     = 21  // Ctrl+U: kill to start of line
	keyCtrlW     = 23  // Ctrl+W: delete word backward
	keyEscape    = 27  // Escape: start of CSI sequences
	keyBackspace = 127 // Backspace: delete previous char
	keyDelete    = 8   // Delete: alternate backspace
)

// action represents the result of processing a keystroke.
type action int

const (
	actionContinue    action = iota // continue reading input
	actionSubmit                    // submit buffer contents
	actionInterrupted               // user pressed Ctrl+C
	actionEOF                       // user pressed Ctrl+D on empty line
)

// actionErrors maps actions to their corresponding errors.
var actionErrors = map[action]error{
	actionSubmit:      nil,
	actionInterrupted: ErrInterrupted,
	actionEOF:         ErrEOF,
}

// keyHandler processes a keystroke and returns the result text and action.
// Handlers may consume additional bytes from reader for escape sequences.
type keyHandler func(i *Input, reader *bufio.Reader) (string, action)

// handleCtrlC clears the buffer and signals an interrupt.
func handleCtrlC(i *Input, reader *bufio.Reader) (string, action) {
	i.buffer = nil
	i.render()
	return "", actionInterrupted
}

// handleCtrlD signals EOF only when the buffer is empty.
// Does nothing if there is text in the buffer.
func handleCtrlD(i *Input, reader *bufio.Reader) (string, action) {
	if len(i.buffer) == 0 {
		return "", actionEOF
	}
	return "", actionContinue
}

// handleTab accepts the current ghost suggestion when cursor is at end.
// Replaces the last word with the full suggestion text.
func handleTab(i *Input, reader *bufio.Reader) (string, action) {
	if i.cursorPos == len(i.buffer) {
		if match := i.findMatch(); match != "" {
			// Accept current suggestion
			start := i.lastWordStart()
			i.buffer = append(i.buffer[:start], []rune(match)...)
			i.cursorPos = len(i.buffer)
			i.matchIndex = 0
			i.render()
		}
	}
	return "", actionContinue
}

// handleEnter submits the current buffer contents.
func handleEnter(i *Input, reader *bufio.Reader) (string, action) {
	// Move cursor to column 1, then to next line
	_, _ = fmt.Fprint(i.out, "\r\n")
	return string(i.buffer), actionSubmit
}

// handleBackspace deletes the character before the cursor.
func handleBackspace(i *Input, reader *bufio.Reader) (string, action) {
	if i.cursorPos > 0 {
		i.buffer = append(i.buffer[:i.cursorPos-1], i.buffer[i.cursorPos:]...)
		i.cursorPos--
		i.matchIndex = 0 // Reset match cycling on buffer change
		i.render()
	}
	return "", actionContinue
}

// handleEscape processes CSI escape sequences (arrow keys, delete, etc.).
// Reads the sequence bytes and delegates to the appropriate CSI handler.
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

// csiHandler processes a CSI (Control Sequence Introducer) escape sequence.
type csiHandler func(i *Input, reader *bufio.Reader)

// csiHandlers maps CSI final bytes to their handlers.
var csiHandlers = map[byte]csiHandler{
	'A': handleUpArrow,
	'B': handleDownArrow,
	'C': handleRightArrow,
	'D': handleLeftArrow,
	'3': handleDelete,
}

// handleCSI dispatches CSI escape sequences to specific handlers.
// Handles both simple sequences (ESC [ A) and extended ones (ESC [ 1 ; 5 C).
func (i *Input) handleCSI(code byte, reader *bufio.Reader) {
	// Handle extended sequences like ESC [ 1 ; 5 C (Ctrl+Right)
	if code == '1' {
		b2, err := reader.ReadByte()
		if err != nil {
			return
		}
		if b2 == ';' {
			mod, _ := reader.ReadByte() // modifier (5 = Ctrl)
			dir, _ := reader.ReadByte() // direction
			if mod == '5' && dir == 'C' {
				handleCtrlRight(i, reader)
				return
			}
			if mod == '5' && dir == 'D' {
				handleCtrlLeft(i, reader)
				return
			}
		}
		return
	}

	if handler, ok := csiHandlers[code]; ok {
		handler(i, reader)
	}
}

// handleUpArrow cycles suggestions, navigates history, or moves up a line.
// Behavior depends on cursor position and available matches.
func handleUpArrow(i *Input, _ *bufio.Reader) {
	// At end of buffer with multiple matches: cycle suggestions
	if i.cursorPos == len(i.buffer) {
		matches := i.getMatches()
		if len(matches) > 1 {
			i.matchIndex = (i.matchIndex - 1 + len(matches)) % len(matches)
			i.render()
			return
		}
	}

	// Find current line start and check if we're on first line
	lineStart := i.findLineStart()
	if lineStart == 0 {
		// On first line, navigate history
		entry, ok := i.history.Previous(string(i.buffer))
		if ok {
			i.buffer = []rune(entry)
			i.cursorPos = len(i.buffer)
			i.matchIndex = 0
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

// handleDownArrow cycles suggestions, navigates history, or moves down a line.
// Behavior depends on cursor position and available matches.
func handleDownArrow(i *Input, _ *bufio.Reader) {
	// At end of buffer with multiple matches: cycle suggestions
	if i.cursorPos == len(i.buffer) {
		matches := i.getMatches()
		if len(matches) > 1 {
			i.matchIndex = (i.matchIndex + 1) % len(matches)
			i.render()
			return
		}
	}

	// Find current line end and check if we're on last line
	lineEnd := i.findLineEnd()
	if lineEnd == len(i.buffer) {
		// On last line, navigate history
		entry, ok := i.history.Next()
		if ok {
			i.buffer = []rune(entry)
			i.cursorPos = len(i.buffer)
			i.matchIndex = 0
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

// findLineStart returns the buffer index of the current line's start.
func (i *Input) findLineStart() int {
	return i.findLineStartFrom(i.cursorPos)
}

// findLineStartFrom returns the start of the line containing the given position.
func (i *Input) findLineStartFrom(pos int) int {
	for pos > 0 && i.buffer[pos-1] != '\n' {
		pos--
	}
	return pos
}

// findLineEnd returns the buffer index of the current line's end.
func (i *Input) findLineEnd() int {
	return i.findLineEndFrom(i.cursorPos)
}

// findLineEndFrom returns the end of the line containing the given position.
func (i *Input) findLineEndFrom(pos int) int {
	for pos < len(i.buffer) && i.buffer[pos] != '\n' {
		pos++
	}
	return pos
}

// handleCtrlRight accepts the next word from ghost text or jumps forward a word.
// At end of buffer: accepts partial ghost suggestion up to next word boundary.
// Mid-buffer: moves cursor to start of next word.
func handleCtrlRight(i *Input, _ *bufio.Reader) {
	if i.cursorPos == len(i.buffer) {
		// At end: accept next word from ghost
		ghost := i.findGhost()
		if ghost == "" {
			return
		}
		// Find end of next word in ghost
		wordEnd := 0
		inWord := false
		for idx, r := range ghost {
			if r == ' ' || r == '\t' {
				if inWord {
					wordEnd = idx
					break
				}
			} else {
				inWord = true
				wordEnd = idx + 1
			}
		}
		if wordEnd > 0 {
			match := i.findMatch()
			text := string(i.buffer)
			lastWord := extractLastWord(text)
			// Accept lastWord + partial ghost
			accepted := match[:len(lastWord)+wordEnd]
			start := i.lastWordStart()
			i.buffer = append(i.buffer[:start], []rune(accepted)...)
			i.cursorPos = len(i.buffer)
			i.matchIndex = 0
			i.render()
		}
		return
	}

	// Move cursor to end of next word
	pos := i.cursorPos
	// Skip current word
	for pos < len(i.buffer) && i.buffer[pos] != ' ' {
		pos++
	}
	// Skip spaces
	for pos < len(i.buffer) && i.buffer[pos] == ' ' {
		pos++
	}
	i.cursorPos = pos
	i.render()
}

// handleCtrlLeft moves cursor backward to the start of the previous word.
func handleCtrlLeft(i *Input, _ *bufio.Reader) {
	// Move cursor to start of previous word
	pos := i.cursorPos
	// Skip spaces before cursor
	for pos > 0 && i.buffer[pos-1] == ' ' {
		pos--
	}
	// Skip word
	for pos > 0 && i.buffer[pos-1] != ' ' {
		pos--
	}
	i.cursorPos = pos
	i.render()
}

// handleRightArrow moves cursor right or accepts suggestion at end of buffer.
func handleRightArrow(i *Input, _ *bufio.Reader) {
	if i.cursorPos < len(i.buffer) {
		i.cursorPos++
		i.render()
	} else if match := i.findMatch(); match != "" {
		// At end of buffer: accept current suggestion
		start := i.lastWordStart()
		i.buffer = append(i.buffer[:start], []rune(match)...)
		i.cursorPos = len(i.buffer)
		i.matchIndex = 0
		i.render()
	}
}

// handleLeftArrow moves cursor one position to the left.
func handleLeftArrow(i *Input, _ *bufio.Reader) {
	if i.cursorPos > 0 {
		i.cursorPos--
		i.render()
	}
}

// handleDelete processes the Delete key CSI sequence.
func handleDelete(i *Input, reader *bufio.Reader) {
	i.handleDeleteKey(reader)
}

// handleDeleteKey deletes the character under the cursor.
func (i *Input) handleDeleteKey(reader *bufio.Reader) {
	b3, err := reader.ReadByte()
	if err != nil || b3 != '~' {
		return
	}
	if i.cursorPos < len(i.buffer) {
		i.buffer = append(i.buffer[:i.cursorPos], i.buffer[i.cursorPos+1:]...)
		i.matchIndex = 0
		i.render()
	}
}

// handleCtrlA moves cursor to the beginning of the current line.
func handleCtrlA(i *Input, reader *bufio.Reader) (string, action) {
	i.cursorPos = i.findLineStart()
	i.render()
	return "", actionContinue
}

// handleCtrlE moves cursor to the end of the current line.
func handleCtrlE(i *Input, reader *bufio.Reader) (string, action) {
	i.cursorPos = i.findLineEnd()
	i.render()
	return "", actionContinue
}

// handleCtrlK deletes text from cursor to end of line (kill forward).
func handleCtrlK(i *Input, reader *bufio.Reader) (string, action) {
	lineEnd := i.findLineEnd()
	i.buffer = append(i.buffer[:i.cursorPos], i.buffer[lineEnd:]...)
	i.matchIndex = 0
	i.render()
	return "", actionContinue
}

// handleCtrlU deletes text from start of line to cursor (kill backward).
func handleCtrlU(i *Input, reader *bufio.Reader) (string, action) {
	lineStart := i.findLineStart()
	i.buffer = append(i.buffer[:lineStart], i.buffer[i.cursorPos:]...)
	i.cursorPos = lineStart
	i.matchIndex = 0
	i.render()
	return "", actionContinue
}

// handleCtrlW deletes the word before the cursor (unix-word-rubout).
func handleCtrlW(i *Input, reader *bufio.Reader) (string, action) {
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
	i.matchIndex = 0
	i.render()
	return "", actionContinue
}

// handleCtrlJ inserts a newline at the cursor for multiline input.
func handleCtrlJ(i *Input, reader *bufio.Reader) (string, action) {
	i.buffer = append(i.buffer[:i.cursorPos], append([]rune{'\n'}, i.buffer[i.cursorPos:]...)...)
	i.cursorPos++
	i.matchIndex = 0
	i.render()
	return "", actionContinue
}

// defaultHandlers returns the standard key bindings for readline input.
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
