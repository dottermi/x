package draw

import (
	"testing"

	"github.com/dottermi/x/render"
	"github.com/dottermi/x/termistyle/style"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// text.go tests

func TestText(t *testing.T) {
	t.Parallel()

	t.Run("should render text at position", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(20, 5)
		fg := style.Color("#FFFFFF")
		bg := style.Color("#000000")

		Text(buf, 5, 2, "Hello", fg, bg)

		assert.Equal(t, 'H', buf.Get(5, 2).Char)
		assert.Equal(t, 'e', buf.Get(6, 2).Char)
		assert.Equal(t, 'l', buf.Get(7, 2).Char)
		assert.Equal(t, 'l', buf.Get(8, 2).Char)
		assert.Equal(t, 'o', buf.Get(9, 2).Char)
		assert.True(t, buf.Get(5, 2).FG.IsSet())
		assert.True(t, buf.Get(5, 2).BG.IsSet())
	})

	t.Run("should clip text beyond buffer width", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)

		Text(buf, 7, 2, "Hello", style.Color(""), style.Color(""))

		assert.Equal(t, 'H', buf.Get(7, 2).Char)
		assert.Equal(t, 'e', buf.Get(8, 2).Char)
		assert.Equal(t, 'l', buf.Get(9, 2).Char)
	})

	t.Run("should handle empty string", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)

		Text(buf, 5, 2, "", style.Color(""), style.Color(""))

		assert.Equal(t, ' ', buf.Get(5, 2).Char)
	})

	t.Run("should handle unicode characters", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(20, 5)

		Text(buf, 0, 0, "日本語", style.Color(""), style.Color(""))

		assert.Equal(t, '日', buf.Get(0, 0).Char)
		assert.Equal(t, '本', buf.Get(1, 0).Char)
		assert.Equal(t, '語', buf.Get(2, 0).Char)
	})

	t.Run("should render at origin", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)

		Text(buf, 0, 0, "ABC", style.Color(""), style.Color(""))

		assert.Equal(t, 'A', buf.Get(0, 0).Char)
		assert.Equal(t, 'B', buf.Get(1, 0).Char)
		assert.Equal(t, 'C', buf.Get(2, 0).Char)
	})
}

func TestStyledText(t *testing.T) {
	t.Parallel()

	t.Run("should render text with style attributes", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(20, 5)
		s := style.Style{
			Foreground:     style.Color("#FF0000"),
			Background:     style.Color("#00FF00"),
			FontWeight:     style.WeightBold,
			FontStyle:      style.StyleItalic,
			TextDecoration: style.DecorationUnderline,
		}

		StyledText(buf, 0, 0, "Test", s)

		cell := buf.Get(0, 0)
		assert.Equal(t, 'T', cell.Char)
		assert.True(t, cell.FG.IsSet())
		assert.True(t, cell.BG.IsSet())
		assert.True(t, cell.Bold)
		assert.True(t, cell.Italic)
		assert.True(t, cell.Underline)
	})

	t.Run("should apply text transform uppercase", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(20, 5)
		s := style.Style{
			TextTransform: style.TransformUppercase,
		}

		StyledText(buf, 0, 0, "hello", s)

		assert.Equal(t, 'H', buf.Get(0, 0).Char)
		assert.Equal(t, 'E', buf.Get(1, 0).Char)
		assert.Equal(t, 'L', buf.Get(2, 0).Char)
		assert.Equal(t, 'L', buf.Get(3, 0).Char)
		assert.Equal(t, 'O', buf.Get(4, 0).Char)
	})

	t.Run("should apply text transform lowercase", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(20, 5)
		s := style.Style{
			TextTransform: style.TransformLowercase,
		}

		StyledText(buf, 0, 0, "HELLO", s)

		assert.Equal(t, 'h', buf.Get(0, 0).Char)
		assert.Equal(t, 'e', buf.Get(1, 0).Char)
	})

	t.Run("should render with dim and reverse attributes", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(20, 5)
		s := style.Style{
			Dim:     true,
			Reverse: true,
		}

		StyledText(buf, 0, 0, "X", s)

		cell := buf.Get(0, 0)
		assert.True(t, cell.Dim)
		assert.True(t, cell.Reverse)
	})

	t.Run("should render with strikethrough decoration", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(20, 5)
		s := style.Style{
			TextDecoration: style.DecorationLineThrough,
		}

		StyledText(buf, 0, 0, "Strike", s)

		assert.True(t, buf.Get(0, 0).Strike)
	})
}

func Test_calculateTextOffset(t *testing.T) {
	t.Parallel()

	t.Run("should return 0 for left alignment", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(10, 20, style.TextAlignLeft)

		assert.Equal(t, 0, offset)
	})

	t.Run("should center text correctly", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(10, 20, style.TextAlignCenter)

		assert.Equal(t, 5, offset)
	})

	t.Run("should right align text", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(10, 20, style.TextAlignRight)

		assert.Equal(t, 10, offset)
	})

	t.Run("should return 0 when line fills container", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(20, 20, style.TextAlignCenter)

		assert.Equal(t, 0, offset)
	})

	t.Run("should return 0 when line exceeds container", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(25, 20, style.TextAlignCenter)

		assert.Equal(t, 0, offset)
	})

	t.Run("should return 0 for zero container width", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(10, 0, style.TextAlignCenter)

		assert.Equal(t, 0, offset)
	})

	t.Run("should return 0 for negative container width", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(10, -5, style.TextAlignCenter)

		assert.Equal(t, 0, offset)
	})

	t.Run("should handle odd centering", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(5, 10, style.TextAlignCenter)

		assert.Equal(t, 2, offset)
	})

	t.Run("should return 0 for unknown alignment", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(10, 20, style.TextAlign(99))

		assert.Equal(t, 0, offset)
	})
}

func Test_applyEllipsis(t *testing.T) {
	t.Parallel()

	t.Run("should return unchanged when fits", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := applyEllipsis(runes, 10)

		assert.Equal(t, []rune("hello"), result)
	})

	t.Run("should return unchanged when exact fit", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := applyEllipsis(runes, 5)

		assert.Equal(t, []rune("hello"), result)
	})

	t.Run("should truncate with ellipsis", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello world")

		result := applyEllipsis(runes, 8)

		assert.Equal(t, []rune("hello..."), result)
	})

	t.Run("should handle maxWidth of 3", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := applyEllipsis(runes, 3)

		assert.Equal(t, []rune("..."), result)
	})

	t.Run("should handle maxWidth of 2", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := applyEllipsis(runes, 2)

		assert.Equal(t, []rune(".."), result)
	})

	t.Run("should handle maxWidth of 1", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := applyEllipsis(runes, 1)

		assert.Equal(t, []rune("."), result)
	})

	t.Run("should handle unicode text", func(t *testing.T) {
		t.Parallel()

		runes := []rune("日本語テスト")

		result := applyEllipsis(runes, 5)

		assert.Equal(t, []rune("日本..."), result)
	})
}

func Test_wrapText(t *testing.T) {
	t.Parallel()

	t.Run("should return single line when no wrapping needed", func(t *testing.T) {
		t.Parallel()

		lines := wrapText("hello", 10, style.WrapWord)

		assert.Equal(t, []string{"hello"}, lines)
	})

	t.Run("should return unchanged for zero width", func(t *testing.T) {
		t.Parallel()

		lines := wrapText("hello", 0, style.WrapWord)

		assert.Equal(t, []string{"hello"}, lines)
	})

	t.Run("should return unchanged for negative width", func(t *testing.T) {
		t.Parallel()

		lines := wrapText("hello", -5, style.WrapWord)

		assert.Equal(t, []string{"hello"}, lines)
	})

	t.Run("should use character wrapping when WrapChar", func(t *testing.T) {
		t.Parallel()

		lines := wrapText("abcdefghij", 4, style.WrapChar)

		assert.Equal(t, []string{"abcd", "efgh", "ij"}, lines)
	})

	t.Run("should use word wrapping for WrapWord", func(t *testing.T) {
		t.Parallel()

		lines := wrapText("hello world", 6, style.WrapWord)

		assert.Equal(t, []string{"hello", "world"}, lines)
	})
}

func Test_wrapByWord(t *testing.T) {
	t.Parallel()

	t.Run("should wrap at word boundaries", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("hello world foo", 10)

		assert.Equal(t, []string{"hello", "world foo"}, lines)
	})

	t.Run("should handle single word", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("hello", 10)

		assert.Equal(t, []string{"hello"}, lines)
	})

	t.Run("should handle empty string", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("", 10)

		assert.Equal(t, []string{""}, lines)
	})

	t.Run("should break long word with character wrapping", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("abcdefghij", 4)

		assert.Equal(t, []string{"abcd", "efgh", "ij"}, lines)
	})

	t.Run("should handle word longer than maxWidth in middle", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("hi superlongword bye", 5)

		require.Len(t, lines, 5)
		assert.Equal(t, "hi", lines[0])
		assert.Equal(t, "super", lines[1])
		assert.Equal(t, "longw", lines[2])
		assert.Equal(t, "ord", lines[3])
		assert.Equal(t, "bye", lines[4])
	})

	t.Run("should join words on same line when space permits", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("a b c d", 5)

		assert.Equal(t, []string{"a b c", "d"}, lines)
	})

	t.Run("should handle multiple spaces between words", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("hello   world", 10)

		assert.Equal(t, []string{"hello", "world"}, lines)
	})

	t.Run("should handle only whitespace", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("   ", 10)

		assert.Equal(t, []string{""}, lines)
	})
}

func Test_wrapByChar(t *testing.T) {
	t.Parallel()

	t.Run("should wrap at character boundaries", func(t *testing.T) {
		t.Parallel()

		lines := wrapByChar("abcdefghij", 3)

		assert.Equal(t, []string{"abc", "def", "ghi", "j"}, lines)
	})

	t.Run("should handle text shorter than width", func(t *testing.T) {
		t.Parallel()

		lines := wrapByChar("hi", 10)

		assert.Equal(t, []string{"hi"}, lines)
	})

	t.Run("should handle exact width", func(t *testing.T) {
		t.Parallel()

		lines := wrapByChar("abcd", 4)

		assert.Equal(t, []string{"abcd"}, lines)
	})

	t.Run("should handle empty string", func(t *testing.T) {
		t.Parallel()

		lines := wrapByChar("", 5)

		assert.Empty(t, lines)
	})

	t.Run("should handle unicode characters", func(t *testing.T) {
		t.Parallel()

		lines := wrapByChar("日本語テスト", 2)

		assert.Equal(t, []string{"日本", "語テ", "スト"}, lines)
	})

	t.Run("should handle width of 1", func(t *testing.T) {
		t.Parallel()

		lines := wrapByChar("abc", 1)

		assert.Equal(t, []string{"a", "b", "c"}, lines)
	})
}

// border.go tests

func TestBorder(t *testing.T) {
	t.Parallel()

	t.Run("should draw single border", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)
		border := style.BorderAll(style.BorderSingle)

		Border(buf, 0, 0, 10, 5, border)

		assert.Equal(t, '┌', buf.Get(0, 0).Char)
		assert.Equal(t, '┐', buf.Get(9, 0).Char)
		assert.Equal(t, '└', buf.Get(0, 4).Char)
		assert.Equal(t, '┘', buf.Get(9, 4).Char)
		assert.Equal(t, '─', buf.Get(5, 0).Char)
		assert.Equal(t, '│', buf.Get(0, 2).Char)
	})

	t.Run("should draw round border", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)
		border := style.BorderAll(style.BorderRound)

		Border(buf, 0, 0, 10, 5, border)

		assert.Equal(t, '╭', buf.Get(0, 0).Char)
		assert.Equal(t, '╮', buf.Get(9, 0).Char)
		assert.Equal(t, '╰', buf.Get(0, 4).Char)
		assert.Equal(t, '╯', buf.Get(9, 4).Char)
	})

	t.Run("should draw double border", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)
		border := style.BorderAll(style.BorderDouble)

		Border(buf, 0, 0, 10, 5, border)

		assert.Equal(t, '╔', buf.Get(0, 0).Char)
		assert.Equal(t, '╗', buf.Get(9, 0).Char)
		assert.Equal(t, '╚', buf.Get(0, 4).Char)
		assert.Equal(t, '╝', buf.Get(9, 4).Char)
		assert.Equal(t, '═', buf.Get(5, 0).Char)
		assert.Equal(t, '║', buf.Get(0, 2).Char)
	})

	t.Run("should draw bold border", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)
		border := style.BorderAll(style.BorderBold)

		Border(buf, 0, 0, 10, 5, border)

		assert.Equal(t, '┏', buf.Get(0, 0).Char)
		assert.Equal(t, '┓', buf.Get(9, 0).Char)
		assert.Equal(t, '┗', buf.Get(0, 4).Char)
		assert.Equal(t, '┛', buf.Get(9, 4).Char)
		assert.Equal(t, '━', buf.Get(5, 0).Char)
		assert.Equal(t, '┃', buf.Get(0, 2).Char)
	})

	t.Run("should apply border color", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)
		border := style.BorderAllWithColor(style.BorderSingle, style.Color("#FF0000"))

		Border(buf, 0, 0, 10, 5, border)

		assert.True(t, buf.Get(0, 0).FG.IsSet())
		assert.True(t, buf.Get(5, 0).FG.IsSet())
		assert.True(t, buf.Get(0, 2).FG.IsSet())
	})

	t.Run("should not draw border with no style", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)
		border := style.Border{}

		Border(buf, 0, 0, 10, 5, border)

		assert.Equal(t, ' ', buf.Get(0, 0).Char)
	})

	t.Run("should not draw border when width less than 2", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)
		border := style.BorderAll(style.BorderSingle)

		Border(buf, 0, 0, 1, 5, border)

		assert.Equal(t, ' ', buf.Get(0, 0).Char)
	})

	t.Run("should not draw border when height less than 2", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)
		border := style.BorderAll(style.BorderSingle)

		Border(buf, 0, 0, 10, 1, border)

		assert.Equal(t, ' ', buf.Get(0, 0).Char)
	})

	t.Run("should draw minimum 2x2 border", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)
		border := style.BorderAll(style.BorderSingle)

		Border(buf, 0, 0, 2, 2, border)

		assert.Equal(t, '┌', buf.Get(0, 0).Char)
		assert.Equal(t, '┐', buf.Get(1, 0).Char)
		assert.Equal(t, '└', buf.Get(0, 1).Char)
		assert.Equal(t, '┘', buf.Get(1, 1).Char)
	})

	t.Run("should draw border at offset position", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(15, 10)
		border := style.BorderAll(style.BorderSingle)

		Border(buf, 3, 2, 8, 5, border)

		assert.Equal(t, '┌', buf.Get(3, 2).Char)
		assert.Equal(t, '┐', buf.Get(10, 2).Char)
		assert.Equal(t, '└', buf.Get(3, 6).Char)
		assert.Equal(t, '┘', buf.Get(10, 6).Char)
		assert.Equal(t, ' ', buf.Get(0, 0).Char)
	})

	t.Run("should draw border with only top and bottom", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)
		edge := style.BorderEdge{Style: style.BorderSingle}
		border := style.Border{Top: edge, Bottom: edge}

		Border(buf, 0, 0, 10, 5, border)

		assert.Equal(t, '─', buf.Get(0, 0).Char)
		assert.Equal(t, '─', buf.Get(0, 4).Char)
		assert.Equal(t, ' ', buf.Get(0, 2).Char)
	})

	t.Run("should draw border with only left and right", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(10, 5)
		edge := style.BorderEdge{Style: style.BorderSingle}
		border := style.Border{Left: edge, Right: edge}

		Border(buf, 0, 0, 10, 5, border)

		assert.Equal(t, '│', buf.Get(0, 0).Char)
		assert.Equal(t, '│', buf.Get(9, 0).Char)
		assert.Equal(t, ' ', buf.Get(5, 0).Char)
	})
}

func TestBorder_WithTitle(t *testing.T) {
	t.Parallel()

	t.Run("should render left-aligned title on top border", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(20, 5)
		border := style.BorderAllWithTitle(style.BorderSingle, style.Color("#FFFFFF"), "Title")

		Border(buf, 0, 0, 20, 5, border)

		assert.Equal(t, '┌', buf.Get(0, 0).Char)
		assert.Equal(t, '─', buf.Get(1, 0).Char)
		assert.Equal(t, ' ', buf.Get(2, 0).Char)
		assert.Equal(t, 'T', buf.Get(3, 0).Char)
		assert.Equal(t, 'i', buf.Get(4, 0).Char)
		assert.Equal(t, 't', buf.Get(5, 0).Char)
		assert.Equal(t, 'l', buf.Get(6, 0).Char)
		assert.Equal(t, 'e', buf.Get(7, 0).Char)
		assert.Equal(t, ' ', buf.Get(8, 0).Char)
		assert.Equal(t, '─', buf.Get(9, 0).Char)
	})

	t.Run("should render right-aligned text on border", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(20, 5)
		topEdge := style.BorderEdge{Style: style.BorderSingle}
		topEdge = topEdge.AddText("Right", style.TextAlignRight)
		border := style.BorderAll(style.BorderSingle)
		border.Top = topEdge

		Border(buf, 0, 0, 20, 5, border)

		assert.Equal(t, 't', buf.Get(16, 0).Char)
		assert.Equal(t, ' ', buf.Get(17, 0).Char)
		assert.Equal(t, '┐', buf.Get(19, 0).Char)
	})

	t.Run("should render center-aligned text on border", func(t *testing.T) {
		t.Parallel()

		buf := render.NewBuffer(20, 5)
		topEdge := style.BorderEdge{Style: style.BorderSingle}
		topEdge = topEdge.AddText("Hi", style.TextAlignCenter)
		border := style.BorderAll(style.BorderSingle)
		border.Top = topEdge

		Border(buf, 0, 0, 20, 5, border)

		foundH := false
		for x := 1; x < 19; x++ {
			if buf.Get(x, 0).Char == 'H' {
				foundH = true
				assert.Equal(t, 'i', buf.Get(x+1, 0).Char)
				break
			}
		}
		assert.True(t, foundH, "should find 'Hi' text in border")
	})
}

func Test_truncateRunes(t *testing.T) {
	t.Parallel()

	t.Run("should return unchanged when fits", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunes(runes, 10)

		assert.Equal(t, []rune("hello"), result)
	})

	t.Run("should return unchanged when exact fit", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunes(runes, 5)

		assert.Equal(t, []rune("hello"), result)
	})

	t.Run("should truncate with ellipsis when too long", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello world")

		result := truncateRunes(runes, 8)

		assert.Equal(t, []rune("hello..."), result)
	})

	t.Run("should truncate without ellipsis when maxLen <= 3", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunes(runes, 3)

		assert.Equal(t, []rune("hel"), result)
	})

	t.Run("should truncate to 2 chars without ellipsis", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunes(runes, 2)

		assert.Equal(t, []rune("he"), result)
	})

	t.Run("should truncate to 1 char", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunes(runes, 1)

		assert.Equal(t, []rune("h"), result)
	})

	t.Run("should handle unicode", func(t *testing.T) {
		t.Parallel()

		runes := []rune("日本語テスト")

		result := truncateRunes(runes, 5)

		assert.Equal(t, []rune("日本..."), result)
	})
}

func Test_truncateRunesFromStart(t *testing.T) {
	t.Parallel()

	t.Run("should return unchanged when fits", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunesFromStart(runes, 10)

		assert.Equal(t, []rune("hello"), result)
	})

	t.Run("should return unchanged when exact fit", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunesFromStart(runes, 5)

		assert.Equal(t, []rune("hello"), result)
	})

	t.Run("should truncate from start with ellipsis", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello world")

		result := truncateRunesFromStart(runes, 8)

		assert.Equal(t, []rune("...world"), result)
	})

	t.Run("should truncate without ellipsis when maxLen <= 3", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunesFromStart(runes, 3)

		assert.Equal(t, []rune("llo"), result)
	})

	t.Run("should truncate to last 2 chars", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunesFromStart(runes, 2)

		assert.Equal(t, []rune("lo"), result)
	})

	t.Run("should truncate to last 1 char", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunesFromStart(runes, 1)

		assert.Equal(t, []rune("o"), result)
	})

	t.Run("should handle unicode", func(t *testing.T) {
		t.Parallel()

		runes := []rune("日本語テスト")

		result := truncateRunesFromStart(runes, 5)

		assert.Equal(t, []rune("...スト"), result)
	})
}
