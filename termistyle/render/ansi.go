package render

import (
	"fmt"

	"github.com/dottermi/x/termistyle/style"
)

// ANSI escape codes.
const (
	Reset     = "\x1b[0m"
	Clear     = "\x1b[2J\x1b[H"
	fgFormat  = "\x1b[38;2;%d;%d;%dm"
	bgFormat  = "\x1b[48;2;%d;%d;%dm"
	fgDefault = "\x1b[39m"
	bgDefault = "\x1b[49m"

	// Bold enables bold/bright text.
	Bold = "\x1b[1m"
	// Dim enables dim/faint text.
	Dim = "\x1b[2m"
	// Italic enables italic text.
	Italic = "\x1b[3m"
	// Underline enables underlined text.
	Underline = "\x1b[4m"
	// Reverse enables reverse video (swap fg/bg).
	Reverse = "\x1b[7m"
	// Strikethrough enables strikethrough text.
	Strikethrough = "\x1b[9m"

	// NoBold disables bold and dim.
	NoBold = "\x1b[22m"
	// NoDim disables dim (same as NoBold).
	NoDim = "\x1b[22m"
	// NoItalic disables italic.
	NoItalic = "\x1b[23m"
	// NoUnderline disables underline.
	NoUnderline = "\x1b[24m"
	// NoReverse disables reverse video.
	NoReverse = "\x1b[27m"
	// NoStrike disables strikethrough.
	NoStrike = "\x1b[29m"
)

// colorCode generates ANSI escape sequence for foreground and background colors.
// Empty color means use terminal default (resets that color).
func colorCode(fg, bg style.Color) string {
	return fgCode(fg) + bgCode(bg)
}

// fgCode generates ANSI escape sequence for foreground color.
func fgCode(c style.Color) string {
	if !c.IsSet() {
		return fgDefault
	}
	return fmt.Sprintf(fgFormat, c.R(), c.G(), c.B())
}

// bgCode generates ANSI escape sequence for background color.
func bgCode(c style.Color) string {
	if !c.IsSet() {
		return bgDefault
	}
	return fmt.Sprintf(bgFormat, c.R(), c.G(), c.B())
}
