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
//	line, ok := input.Readline("$ ")
//	if ok {
//		fmt.Println("You entered:", line)
//	}
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
type Input struct {
	buffer      []rune
	suggestions []string
	prompt      string
	handlers    map[rune]keyHandler

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
	return &Input{
		suggestions: suggestions,
		handlers:    defaultHandlers(),
		in:          in,
		out:         out,
		fd:          int(os.Stdin.Fd()),
	}
}

// Readline reads a line of input with interactive ghost text suggestions.
// Returns the entered text and true on success, or empty string and false if aborted.
// Places the terminal in raw mode for the duration of input.
//
// Parameters:
//   - prompt: string displayed before the input area (e.g., "$ " or "> ")
//
// Keyboard controls:
//   - Tab: accept the current ghost text suggestion
//   - Enter: submit the current input
//   - Backspace/Delete: remove the last character
//   - Ctrl+C: abort input immediately
//   - Ctrl+D: abort input when buffer is empty
//
// Returns false when the user aborts with Ctrl+C, Ctrl+D on empty input,
// or when an I/O error occurs.
//
// Example:
//
//	input := ghostline.NewInput([]string{"help", "history"}, nil, nil)
//	for {
//		line, ok := input.Readline(">>> ")
//		if !ok {
//			break
//		}
//		fmt.Println("Command:", line)
//	}
func (i *Input) Readline(prompt string) (string, bool) {
	i.prompt = prompt
	i.buffer = []rune{}

	if err := i.enableRawMode(); err != nil {
		return "", false
	}
	defer i.disableRawMode()

	reader := bufio.NewReader(i.in)
	i.render()

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return "", false
		}

		if handler, exists := i.handlers[r]; exists {
			result, act := handler(i, reader)
			if act != actionContinue {
				return result, act == actionSubmit
			}
			continue
		}

		if unicode.IsPrint(r) {
			i.buffer = append(i.buffer, r)
			i.render()
		}
	}
}
