package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dottermi/x/termistyle/draw"
	"github.com/dottermi/x/termistyle/style"
	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	t.Parallel()

	t.Run("Reset constant should be correct ANSI reset code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[0m", Reset)
	})

	t.Run("Clear constant should contain clear screen and cursor home codes", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[2J\x1b[H", Clear)
	})

	t.Run("Bold constant should be correct ANSI code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[1m", Bold)
	})

	t.Run("Dim constant should be correct ANSI code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[2m", Dim)
	})

	t.Run("Italic constant should be correct ANSI code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[3m", Italic)
	})

	t.Run("Underline constant should be correct ANSI code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[4m", Underline)
	})

	t.Run("Reverse constant should be correct ANSI code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[7m", Reverse)
	})

	t.Run("Strikethrough constant should be correct ANSI code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[9m", Strikethrough)
	})

	t.Run("NoBold constant should be correct ANSI code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[22m", NoBold)
	})

	t.Run("NoDim constant should be correct ANSI code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[22m", NoDim)
	})

	t.Run("NoItalic constant should be correct ANSI code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[23m", NoItalic)
	})

	t.Run("NoUnderline constant should be correct ANSI code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[24m", NoUnderline)
	})

	t.Run("NoReverse constant should be correct ANSI code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[27m", NoReverse)
	})

	t.Run("NoStrike constant should be correct ANSI code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "\x1b[29m", NoStrike)
	})

	t.Run("NoBold and NoDim should be the same code", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, NoBold, NoDim)
	})
}

func Test_fgCode(t *testing.T) {
	t.Parallel()

	t.Run("should return default foreground code when color is empty", func(t *testing.T) {
		t.Parallel()

		color := style.Color("")

		result := fgCode(color)

		assert.Equal(t, "\x1b[39m", result)
	})

	t.Run("should generate RGB foreground code for red color", func(t *testing.T) {
		t.Parallel()

		color := style.Color("#FF0000")

		result := fgCode(color)

		assert.Equal(t, "\x1b[38;2;255;0;0m", result)
	})

	t.Run("should generate RGB foreground code for green color", func(t *testing.T) {
		t.Parallel()

		color := style.Color("#00FF00")

		result := fgCode(color)

		assert.Equal(t, "\x1b[38;2;0;255;0m", result)
	})

	t.Run("should generate RGB foreground code for blue color", func(t *testing.T) {
		t.Parallel()

		color := style.Color("#0000FF")

		result := fgCode(color)

		assert.Equal(t, "\x1b[38;2;0;0;255m", result)
	})

	t.Run("should generate RGB foreground code for white color", func(t *testing.T) {
		t.Parallel()

		color := style.Color("#FFFFFF")

		result := fgCode(color)

		assert.Equal(t, "\x1b[38;2;255;255;255m", result)
	})

	t.Run("should generate RGB foreground code for black color", func(t *testing.T) {
		t.Parallel()

		color := style.Color("#000000")

		result := fgCode(color)

		assert.Equal(t, "\x1b[38;2;0;0;0m", result)
	})

	t.Run("should handle hex color without hash prefix", func(t *testing.T) {
		t.Parallel()

		color := style.Color("AABBCC")

		result := fgCode(color)

		assert.Equal(t, "\x1b[38;2;170;187;204m", result)
	})

	t.Run("should generate RGB foreground code for arbitrary color", func(t *testing.T) {
		t.Parallel()

		color := style.Color("#123456")

		result := fgCode(color)

		assert.Equal(t, "\x1b[38;2;18;52;86m", result)
	})
}

func Test_bgCode(t *testing.T) {
	t.Parallel()

	t.Run("should return default background code when color is empty", func(t *testing.T) {
		t.Parallel()

		color := style.Color("")

		result := bgCode(color)

		assert.Equal(t, "\x1b[49m", result)
	})

	t.Run("should generate RGB background code for red color", func(t *testing.T) {
		t.Parallel()

		color := style.Color("#FF0000")

		result := bgCode(color)

		assert.Equal(t, "\x1b[48;2;255;0;0m", result)
	})

	t.Run("should generate RGB background code for green color", func(t *testing.T) {
		t.Parallel()

		color := style.Color("#00FF00")

		result := bgCode(color)

		assert.Equal(t, "\x1b[48;2;0;255;0m", result)
	})

	t.Run("should generate RGB background code for blue color", func(t *testing.T) {
		t.Parallel()

		color := style.Color("#0000FF")

		result := bgCode(color)

		assert.Equal(t, "\x1b[48;2;0;0;255m", result)
	})

	t.Run("should generate RGB background code for white color", func(t *testing.T) {
		t.Parallel()

		color := style.Color("#FFFFFF")

		result := bgCode(color)

		assert.Equal(t, "\x1b[48;2;255;255;255m", result)
	})

	t.Run("should generate RGB background code for black color", func(t *testing.T) {
		t.Parallel()

		color := style.Color("#000000")

		result := bgCode(color)

		assert.Equal(t, "\x1b[48;2;0;0;0m", result)
	})

	t.Run("should handle hex color without hash prefix", func(t *testing.T) {
		t.Parallel()

		color := style.Color("AABBCC")

		result := bgCode(color)

		assert.Equal(t, "\x1b[48;2;170;187;204m", result)
	})
}

func Test_colorCode(t *testing.T) {
	t.Parallel()

	t.Run("should return both default codes when both colors are empty", func(t *testing.T) {
		t.Parallel()

		fg := style.Color("")
		bg := style.Color("")

		result := colorCode(fg, bg)

		assert.Equal(t, "\x1b[39m\x1b[49m", result)
	})

	t.Run("should combine foreground and default background codes", func(t *testing.T) {
		t.Parallel()

		fg := style.Color("#FF0000")
		bg := style.Color("")

		result := colorCode(fg, bg)

		assert.Equal(t, "\x1b[38;2;255;0;0m\x1b[49m", result)
	})

	t.Run("should combine default foreground and background codes", func(t *testing.T) {
		t.Parallel()

		fg := style.Color("")
		bg := style.Color("#0000FF")

		result := colorCode(fg, bg)

		assert.Equal(t, "\x1b[39m\x1b[48;2;0;0;255m", result)
	})

	t.Run("should combine both RGB foreground and background codes", func(t *testing.T) {
		t.Parallel()

		fg := style.Color("#FF0000")
		bg := style.Color("#00FF00")

		result := colorCode(fg, bg)

		assert.Equal(t, "\x1b[38;2;255;0;0m\x1b[48;2;0;255;0m", result)
	})
}

func Test_attrCode(t *testing.T) {
	t.Parallel()

	t.Run("should return empty string when no attributes change", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{}
		to := textAttrs{}

		result := attrCode(from, to)

		assert.Empty(t, result)
	})

	t.Run("should enable bold when transitioning from off to on", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{bold: false}
		to := textAttrs{bold: true}

		result := attrCode(from, to)

		assert.Equal(t, Bold, result)
	})

	t.Run("should disable bold when transitioning from on to off", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{bold: true}
		to := textAttrs{bold: false}

		result := attrCode(from, to)

		assert.Equal(t, NoBold, result)
	})

	t.Run("should enable dim when transitioning from off to on", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{dim: false}
		to := textAttrs{dim: true}

		result := attrCode(from, to)

		assert.Equal(t, Dim, result)
	})

	t.Run("should disable dim when transitioning from on to off", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{dim: true}
		to := textAttrs{dim: false}

		result := attrCode(from, to)

		assert.Equal(t, NoDim, result)
	})

	t.Run("should enable italic when transitioning from off to on", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{italic: false}
		to := textAttrs{italic: true}

		result := attrCode(from, to)

		assert.Equal(t, Italic, result)
	})

	t.Run("should disable italic when transitioning from on to off", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{italic: true}
		to := textAttrs{italic: false}

		result := attrCode(from, to)

		assert.Equal(t, NoItalic, result)
	})

	t.Run("should enable underline when transitioning from off to on", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{underline: false}
		to := textAttrs{underline: true}

		result := attrCode(from, to)

		assert.Equal(t, Underline, result)
	})

	t.Run("should disable underline when transitioning from on to off", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{underline: true}
		to := textAttrs{underline: false}

		result := attrCode(from, to)

		assert.Equal(t, NoUnderline, result)
	})

	t.Run("should enable reverse when transitioning from off to on", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{reverse: false}
		to := textAttrs{reverse: true}

		result := attrCode(from, to)

		assert.Equal(t, Reverse, result)
	})

	t.Run("should disable reverse when transitioning from on to off", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{reverse: true}
		to := textAttrs{reverse: false}

		result := attrCode(from, to)

		assert.Equal(t, NoReverse, result)
	})

	t.Run("should enable strikethrough when transitioning from off to on", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{strike: false}
		to := textAttrs{strike: true}

		result := attrCode(from, to)

		assert.Equal(t, Strikethrough, result)
	})

	t.Run("should disable strikethrough when transitioning from on to off", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{strike: true}
		to := textAttrs{strike: false}

		result := attrCode(from, to)

		assert.Equal(t, NoStrike, result)
	})

	t.Run("should enable multiple attributes simultaneously", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{}
		to := textAttrs{bold: true, italic: true, underline: true}

		result := attrCode(from, to)

		assert.Contains(t, result, Bold)
		assert.Contains(t, result, Italic)
		assert.Contains(t, result, Underline)
	})

	t.Run("should disable multiple attributes simultaneously", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{bold: true, italic: true, underline: true}
		to := textAttrs{}

		result := attrCode(from, to)

		assert.Contains(t, result, NoBold)
		assert.Contains(t, result, NoItalic)
		assert.Contains(t, result, NoUnderline)
	})

	t.Run("should enable some and disable other attributes", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{bold: true, italic: false}
		to := textAttrs{bold: false, italic: true}

		result := attrCode(from, to)

		assert.Contains(t, result, NoBold)
		assert.Contains(t, result, Italic)
	})

	t.Run("should not output codes for unchanged attributes", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{bold: true, italic: true}
		to := textAttrs{bold: true, italic: true}

		result := attrCode(from, to)

		assert.Empty(t, result)
	})

	t.Run("should enable all attributes at once", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{}
		to := textAttrs{
			bold:      true,
			dim:       true,
			italic:    true,
			underline: true,
			reverse:   true,
			strike:    true,
		}

		result := attrCode(from, to)

		assert.Contains(t, result, Bold)
		assert.Contains(t, result, Dim)
		assert.Contains(t, result, Italic)
		assert.Contains(t, result, Underline)
		assert.Contains(t, result, Reverse)
		assert.Contains(t, result, Strikethrough)
	})

	t.Run("should disable all attributes at once", func(t *testing.T) {
		t.Parallel()

		from := textAttrs{
			bold:      true,
			dim:       true,
			italic:    true,
			underline: true,
			reverse:   true,
			strike:    true,
		}
		to := textAttrs{}

		result := attrCode(from, to)

		assert.Contains(t, result, NoBold)
		assert.Contains(t, result, NoDim)
		assert.Contains(t, result, NoItalic)
		assert.Contains(t, result, NoUnderline)
		assert.Contains(t, result, NoReverse)
		assert.Contains(t, result, NoStrike)
	})
}

func TestTo(t *testing.T) {
	t.Parallel()

	t.Run("should render empty buffer with default colors", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, "\x1b[39m")
		assert.Contains(t, result, "\x1b[49m")
		assert.Contains(t, result, Reset)
	})

	t.Run("should render null char as space", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{Char: 0})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, " ")
	})

	t.Run("should render single character", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{Char: 'A'})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, "A")
	})

	t.Run("should render character with foreground color", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{
			Char:       'X',
			Foreground: style.Color("#FF0000"),
		})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, "\x1b[38;2;255;0;0m")
		assert.Contains(t, result, "X")
	})

	t.Run("should render character with background color", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{
			Char:       'Y',
			Background: style.Color("#00FF00"),
		})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, "\x1b[48;2;0;255;0m")
		assert.Contains(t, result, "Y")
	})

	t.Run("should render character with both foreground and background colors", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{
			Char:       'Z',
			Foreground: style.Color("#FF0000"),
			Background: style.Color("#0000FF"),
		})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, "\x1b[38;2;255;0;0m")
		assert.Contains(t, result, "\x1b[48;2;0;0;255m")
		assert.Contains(t, result, "Z")
	})

	t.Run("should render bold character", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{Char: 'B', Bold: true})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, Bold)
		assert.Contains(t, result, "B")
	})

	t.Run("should render italic character", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{Char: 'I', Italic: true})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, Italic)
		assert.Contains(t, result, "I")
	})

	t.Run("should render underlined character", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{Char: 'U', Underline: true})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, Underline)
		assert.Contains(t, result, "U")
	})

	t.Run("should render strikethrough character", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{Char: 'S', Strike: true})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, Strikethrough)
		assert.Contains(t, result, "S")
	})

	t.Run("should render dim character", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{Char: 'D', Dim: true})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, Dim)
		assert.Contains(t, result, "D")
	})

	t.Run("should render reverse video character", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{Char: 'R', Reverse: true})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, Reverse)
		assert.Contains(t, result, "R")
	})

	t.Run("should render multiple characters on single line", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(3, 1)
		buf.Set(0, 0, draw.Cell{Char: 'A'})
		buf.Set(1, 0, draw.Cell{Char: 'B'})
		buf.Set(2, 0, draw.Cell{Char: 'C'})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, "A")
		assert.Contains(t, result, "B")
		assert.Contains(t, result, "C")
	})

	t.Run("should render multiple lines separated by newline", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 2)
		buf.Set(0, 0, draw.Cell{Char: 'X'})
		buf.Set(0, 1, draw.Cell{Char: 'Y'})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, "X")
		assert.Contains(t, result, "Y")
		assert.Contains(t, result, "\n")
	})

	t.Run("should not add trailing newline after last line", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 2)
		buf.Set(0, 0, draw.Cell{Char: 'A'})
		buf.Set(0, 1, draw.Cell{Char: 'B'})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.False(t, strings.HasSuffix(result, "\n"))
	})

	t.Run("should reset style at end of each line", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 2)
		buf.Set(0, 0, draw.Cell{Char: 'A', Bold: true})
		buf.Set(0, 1, draw.Cell{Char: 'B'})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		lines := strings.Split(result, "\n")
		assert.True(t, strings.HasSuffix(lines[0], Reset))
	})

	t.Run("should optimize color codes by not repeating unchanged colors", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(2, 1)
		red := style.Color("#FF0000")
		buf.Set(0, 0, draw.Cell{Char: 'A', Foreground: red})
		buf.Set(1, 0, draw.Cell{Char: 'B', Foreground: red})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		colorCode := "\x1b[38;2;255;0;0m"
		count := strings.Count(result, colorCode)
		assert.Equal(t, 1, count, "should only emit color code once for consecutive same-colored cells")
	})

	t.Run("should emit new color code when color changes", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(2, 1)
		buf.Set(0, 0, draw.Cell{Char: 'A', Foreground: style.Color("#FF0000")})
		buf.Set(1, 0, draw.Cell{Char: 'B', Foreground: style.Color("#00FF00")})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, "\x1b[38;2;255;0;0m")
		assert.Contains(t, result, "\x1b[38;2;0;255;0m")
	})

	t.Run("should render character with all styling attributes", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{
			Char:       'X',
			Foreground: style.Color("#FFFFFF"),
			Background: style.Color("#000000"),
			Bold:       true,
			Italic:     true,
			Underline:  true,
			Strike:     true,
			Dim:        true,
			Reverse:    true,
		})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, "\x1b[38;2;255;255;255m")
		assert.Contains(t, result, "\x1b[48;2;0;0;0m")
		assert.Contains(t, result, Bold)
		assert.Contains(t, result, Italic)
		assert.Contains(t, result, Underline)
		assert.Contains(t, result, Strikethrough)
		assert.Contains(t, result, Dim)
		assert.Contains(t, result, Reverse)
		assert.Contains(t, result, "X")
	})

	t.Run("should handle attribute transition between cells", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(2, 1)
		buf.Set(0, 0, draw.Cell{Char: 'A', Bold: true})
		buf.Set(1, 0, draw.Cell{Char: 'B', Italic: true})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, Bold)
		assert.Contains(t, result, NoBold)
		assert.Contains(t, result, Italic)
	})

	t.Run("should render unicode characters", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(3, 1)
		buf.Set(0, 0, draw.Cell{Char: '日'})
		buf.Set(1, 0, draw.Cell{Char: '本'})
		buf.Set(2, 0, draw.Cell{Char: '語'})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, "日")
		assert.Contains(t, result, "本")
		assert.Contains(t, result, "語")
	})

	t.Run("should render emoji characters", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{Char: '🎉'})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, "🎉")
	})

	t.Run("should render box drawing characters", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(3, 3)
		buf.Set(0, 0, draw.Cell{Char: '┌'})
		buf.Set(1, 0, draw.Cell{Char: '─'})
		buf.Set(2, 0, draw.Cell{Char: '┐'})
		buf.Set(0, 1, draw.Cell{Char: '│'})
		buf.Set(2, 1, draw.Cell{Char: '│'})
		buf.Set(0, 2, draw.Cell{Char: '└'})
		buf.Set(1, 2, draw.Cell{Char: '─'})
		buf.Set(2, 2, draw.Cell{Char: '┘'})
		var writer bytes.Buffer

		To(buf, &writer)
		result := writer.String()

		assert.Contains(t, result, "┌")
		assert.Contains(t, result, "─")
		assert.Contains(t, result, "┐")
		assert.Contains(t, result, "│")
		assert.Contains(t, result, "└")
		assert.Contains(t, result, "┘")
	})
}

func TestRender(t *testing.T) {
	t.Parallel()

	t.Run("should return string representation of buffer", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(3, 1)
		buf.Set(0, 0, draw.Cell{Char: 'A'})
		buf.Set(1, 0, draw.Cell{Char: 'B'})
		buf.Set(2, 0, draw.Cell{Char: 'C'})

		result := Render(buf)

		assert.Contains(t, result, "A")
		assert.Contains(t, result, "B")
		assert.Contains(t, result, "C")
	})

	t.Run("should return empty string representation for empty buffer", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(0, 0)

		result := Render(buf)

		assert.NotNil(t, result)
	})

	t.Run("should include color codes in string output", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{
			Char:       'X',
			Foreground: style.Color("#FF0000"),
		})

		result := Render(buf)

		assert.Contains(t, result, "\x1b[38;2;255;0;0m")
		assert.Contains(t, result, "X")
	})

	t.Run("should include attribute codes in string output", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{Char: 'B', Bold: true})

		result := Render(buf)

		assert.Contains(t, result, Bold)
		assert.Contains(t, result, "B")
	})

	t.Run("should include reset code at end of lines", func(t *testing.T) {
		t.Parallel()

		buf := draw.NewBuffer(1, 1)
		buf.Set(0, 0, draw.Cell{Char: 'A'})

		result := Render(buf)

		assert.Contains(t, result, Reset)
	})
}

func TestClearScreen(t *testing.T) {
	t.Parallel()

	t.Run("should write clear screen ANSI code to writer", func(t *testing.T) {
		t.Parallel()

		var writer bytes.Buffer

		ClearScreen(&writer)

		assert.Equal(t, Clear, writer.String())
	})

	t.Run("should write clear and home cursor codes", func(t *testing.T) {
		t.Parallel()

		var writer bytes.Buffer

		ClearScreen(&writer)
		result := writer.String()

		assert.Contains(t, result, "\x1b[2J")
		assert.Contains(t, result, "\x1b[H")
	})
}

func TestMoveCursor(t *testing.T) {
	t.Parallel()

	t.Run("should move cursor to position 0,0", func(t *testing.T) {
		t.Parallel()

		var writer bytes.Buffer

		MoveCursor(&writer, 0, 0)

		assert.Equal(t, "\x1b[1;1H", writer.String())
	})

	t.Run("should move cursor to position 10,5", func(t *testing.T) {
		t.Parallel()

		var writer bytes.Buffer

		MoveCursor(&writer, 10, 5)

		assert.Equal(t, "\x1b[6;11H", writer.String())
	})

	t.Run("should move cursor to position 79,23", func(t *testing.T) {
		t.Parallel()

		var writer bytes.Buffer

		MoveCursor(&writer, 79, 23)

		assert.Equal(t, "\x1b[24;80H", writer.String())
	})

	t.Run("should use 1-based coordinates in ANSI code", func(t *testing.T) {
		t.Parallel()

		var writer bytes.Buffer

		MoveCursor(&writer, 0, 0)
		result := writer.String()

		assert.Equal(t, "\x1b[1;1H", result)
		assert.NotContains(t, result, ";0")
		assert.NotContains(t, result, "[0")
	})

	t.Run("should handle large coordinates", func(t *testing.T) {
		t.Parallel()

		var writer bytes.Buffer

		MoveCursor(&writer, 999, 999)

		assert.Equal(t, "\x1b[1000;1000H", writer.String())
	})
}
