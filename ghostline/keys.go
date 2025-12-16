package ghostline

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
