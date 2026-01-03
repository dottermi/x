// Package render converts draw buffers to ANSI escape sequences for terminal output.
package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/dottermi/x/termistyle/draw"
	"github.com/dottermi/x/termistyle/style"
)

// Render converts a Buffer to a string with ANSI escape codes.
func Render(buf *draw.Buffer) string {
	var sb strings.Builder
	To(buf, &sb)
	return sb.String()
}

// textAttrs tracks current text styling attributes.
type textAttrs struct {
	bold      bool
	italic    bool
	underline bool
	strike    bool
	dim       bool
	reverse   bool
}

// attrMapping maps text attribute changes to ANSI codes.
type attrMapping struct {
	fromField bool
	toField   bool
	on        string
	off       string
}

// To writes the buffer content with ANSI codes to the given writer.
func To(buf *draw.Buffer, w io.Writer) {
	var lastFg, lastBg style.Color
	var lastAttrs textAttrs
	hasStyle := false

	for y := range buf.Height {
		for x := range buf.Width {
			cell := buf.Get(x, y)
			attrs := textAttrs{
				bold:      cell.Bold,
				italic:    cell.Italic,
				underline: cell.Underline,
				strike:    cell.Strike,
				dim:       cell.Dim,
				reverse:   cell.Reverse,
			}

			colorChanged := !hasStyle || cell.Foreground != lastFg || cell.Background != lastBg
			attrsChanged := !hasStyle || attrs != lastAttrs

			if colorChanged {
				_, _ = fmt.Fprint(w, colorCode(cell.Foreground, cell.Background))
				lastFg = cell.Foreground
				lastBg = cell.Background
			}

			if attrsChanged {
				_, _ = fmt.Fprint(w, attrCode(lastAttrs, attrs))
				lastAttrs = attrs
			}

			hasStyle = true

			if cell.Char == 0 {
				_, _ = fmt.Fprint(w, " ")
			} else {
				_, _ = fmt.Fprintf(w, "%c", cell.Char)
			}
		}
		_, _ = fmt.Fprint(w, Reset)
		hasStyle = false
		lastAttrs = textAttrs{}
		if y < buf.Height-1 {
			_, _ = fmt.Fprintln(w)
		}
	}
}

// attrCode generates ANSI codes for text attribute changes.
func attrCode(from, to textAttrs) string {
	var sb strings.Builder

	mappings := []attrMapping{
		{from.bold, to.bold, Bold, NoBold},
		{from.dim, to.dim, Dim, NoDim},
		{from.italic, to.italic, Italic, NoItalic},
		{from.underline, to.underline, Underline, NoUnderline},
		{from.reverse, to.reverse, Reverse, NoReverse},
		{from.strike, to.strike, Strikethrough, NoStrike},
	}

	for _, m := range mappings {
		if m.toField && !m.fromField {
			sb.WriteString(m.on)
		} else if !m.toField && m.fromField {
			sb.WriteString(m.off)
		}
	}

	return sb.String()
}

// ClearScreen returns the ANSI code to clear the screen.
func ClearScreen(w io.Writer) {
	_, _ = fmt.Fprint(w, Clear)
}

// MoveCursor moves the cursor to position (x, y).
func MoveCursor(w io.Writer, x, y int) {
	_, _ = fmt.Fprintf(w, "\x1b[%d;%dH", y+1, x+1)
}
