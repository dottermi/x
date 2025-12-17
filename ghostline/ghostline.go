// Package ghostline provides interactive readline functionality with "ghost text" suggestions.
// Displays completion suggestions inline as dimmed text that users can accept with Tab.
// Supports standard line editing with backspace and handles terminal raw mode automatically.
//
// Ghost text suggestions appear after the cursor as the user types, showing potential
// completions based on prefix matching against a configurable list of suggestions.
// The terminal is placed in raw mode during input to capture individual keystrokes.
//
// Example:
//
//	suggestions := []string{"help", "hello", "history", "exit"}
//	input := ghostline.NewInput(suggestions, nil, nil)
//	line, err := input.Readline("$ ")
//	if err != nil {
//		// handle ErrInterrupted, ErrEOF, or other errors
//	}
//	fmt.Println("You entered:", line)
package ghostline

import (
	"bufio"
	"io"
	"os"
	"unicode"

	"golang.org/x/term"
)

// Input manages interactive line input with ghost text suggestions.
// Handles raw terminal mode, keyboard input, and inline suggestion rendering.
// Create instances using NewInput rather than constructing directly.
//
// The zero value is not usable; always use NewInput to create Input instances.
//
// Input is not safe for concurrent use. Each Input instance must be used by
// a single goroutine only. If you need to read from multiple inputs
// concurrently, create a separate Input instance per goroutine.
type Input struct {
	buffer      []rune
	cursorPos   int
	suggestions []string
	prompt      string
	contPrompt  string // continuation prompt for multiline
	prevLines   int    // previous line count for clearing
	handlers    map[rune]keyHandler
	history     *History

	in  io.Reader
	out io.Writer

	fd       int
	oldState *term.State
}

// NewInput creates an Input configured with the given suggestions and I/O streams.
// Returns a ready-to-use Input for interactive line reading with ghost text.
//
// Parameters:
//   - suggestions: list of strings to use for prefix-based completion
//   - in: input source for reading keystrokes (nil defaults to os.Stdin)
//   - out: output destination for rendering (nil defaults to os.Stdout)
//
// Example:
//
//	// Using defaults for standard terminal interaction
//	input := ghostline.NewInput([]string{"commit", "checkout", "cherry-pick"}, nil, nil)
//
//	// Using custom streams for testing
//	input := ghostline.NewInput(suggestions, mockReader, mockWriter)
func NewInput(suggestions []string, in io.Reader, out io.Writer) *Input {
	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}

	fd := int(os.Stdin.Fd())
	if f, ok := in.(*os.File); ok {
		fd = int(f.Fd())
	}

	return &Input{
		suggestions: suggestions,
		handlers:    defaultHandlers(),
		history:     NewHistory(),
		in:          in,
		out:         out,
		fd:          fd,
	}
}

// AddHistory adds a line to the command history.
// Empty lines and duplicates of the last entry are ignored.
func (i *Input) AddHistory(line string) {
	i.history.Add(line)
}

// Readline reads a line of input with interactive ghost text suggestions.
// Returns the entered text and nil on success, or an error if aborted.
// Places the terminal in raw mode for the duration of input.
//
// Parameters:
//   - prompt: string displayed before the input area (e.g., "$ " or "> ")
//
// Keyboard controls:
//   - Tab: accept the current ghost text suggestion
//   - Enter: submit the current input
//   - Backspace/Delete: remove the last character
//   - Ctrl+C: abort input (returns ErrInterrupted)
//   - Ctrl+D: abort input when buffer is empty (returns ErrEOF)
//
// Example:
//
//	input := ghostline.NewInput([]string{"help", "history"}, nil, nil)
//	for {
//		line, err := input.Readline(">>> ")
//
//		if err == ghostline.ErrInterrupted {
//			fmt.Println("^C")
//			continue
//		}
//
//		if err == ghostline.ErrEOF {
//			fmt.Println("Goodbye!")
//			break
//		}
//
//		if err != nil {
//			fmt.Println("Error:", err)
//			break
//		}
//
//		fmt.Println("Command:", line)
//	}
func (i *Input) Readline(prompt string) (string, error) {
	i.prompt = prompt
	i.contPrompt = "... "
	i.buffer = []rune{}
	i.cursorPos = 0
	i.prevLines = 1
	i.history.Reset("")

	if err := i.enableRawMode(); err != nil {
		return "", err
	}
	defer i.disableRawMode()

	reader := bufio.NewReader(i.in)
	i.render()

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return "", err
		}

		if handler, exists := i.handlers[r]; exists {
			result, act := handler(i, reader)
			if act == actionContinue {
				continue
			}
			return result, actionErrors[act]
		}

		if unicode.IsPrint(r) {
			// Insert at cursor position
			i.buffer = append(i.buffer[:i.cursorPos], append([]rune{r}, i.buffer[i.cursorPos:]...)...)
			i.cursorPos++
			i.render()
		}
	}
}
