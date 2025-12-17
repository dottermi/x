package ghostline

import "errors"

// Readline errors for distinguishing abort types.
var (
	// ErrInterrupted is returned when the user presses Ctrl+C.
	ErrInterrupted = errors.New("interrupted")

	// ErrEOF is returned when the user presses Ctrl+D on an empty line.
	ErrEOF = errors.New("eof")
)
