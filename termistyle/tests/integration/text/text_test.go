package text_test

import (
	"testing"

	"github.com/dottermi/termitest/assert"
	"github.com/dottermi/x/render"

	"github.com/dottermi/x/termistyle/draw"
	"github.com/dottermi/x/termistyle/style"
)

func TestGolden_Text(t *testing.T) {
	buf := render.NewBuffer(30, 20)

	// Basic text with colors
	draw.DrawText(buf, 0, 0, "Basic Text", style.Color("#FFFFFF"), style.Color(""))
	draw.DrawText(buf, 0, 1, "Red", style.Color("#FF0000"), style.Color(""))
	draw.DrawText(buf, 0, 2, "Green", style.Color("#00FF00"), style.Color(""))
	draw.DrawText(buf, 0, 3, "Blue", style.Color("#0000FF"), style.Color(""))

	// Bold
	draw.DrawStyledText(buf, 0, 5, "Bold", style.Style{
		Foreground: style.Color("#FFFFFF"),
		FontWeight: style.WeightBold,
	})

	// Italic
	draw.DrawStyledText(buf, 10, 5, "Italic", style.Style{
		Foreground: style.Color("#FFFFFF"),
		FontStyle:  style.StyleItalic,
	})

	// Underline
	draw.DrawStyledText(buf, 0, 6, "Underline", style.Style{
		Foreground:     style.Color("#FFFFFF"),
		TextDecoration: style.DecorationUnderline,
	})

	// Strikethrough
	draw.DrawStyledText(buf, 12, 6, "Strike", style.Style{
		Foreground:     style.Color("#FFFFFF"),
		TextDecoration: style.DecorationLineThrough,
	})

	// Dim and Reverse
	draw.DrawStyledText(buf, 0, 7, "Dim", style.Style{
		Foreground: style.Color("#FFFFFF"),
		Dim:        true,
	})
	draw.DrawStyledText(buf, 10, 7, "Reverse", style.Style{
		Foreground: style.Color("#FFFFFF"),
		Reverse:    true,
	})

	// Text transform
	draw.DrawStyledText(buf, 0, 9, "hello", style.Style{
		Foreground:    style.Color("#FFFFFF"),
		TextTransform: style.TransformUppercase,
	})
	draw.DrawStyledText(buf, 10, 9, "WORLD", style.Style{
		Foreground:    style.Color("#FFFFFF"),
		TextTransform: style.TransformLowercase,
	})

	// Text align
	draw.DrawStyledTextInWidth(buf, 0, 11, "Left", 15, style.Style{
		Foreground: style.Color("#FFFFFF"),
		TextAlign:  style.TextAlignLeft,
	})
	draw.DrawStyledTextInWidth(buf, 0, 12, "Center", 15, style.Style{
		Foreground: style.Color("#FFFFFF"),
		TextAlign:  style.TextAlignCenter,
	})
	draw.DrawStyledTextInWidth(buf, 0, 13, "Right", 15, style.Style{
		Foreground: style.Color("#FFFFFF"),
		TextAlign:  style.TextAlignRight,
	})

	// Word wrap
	draw.DrawStyledTextInBox(buf, 0, 15, "Word wrap test here", 10, 2, style.Style{
		Foreground: style.Color("#FFFFFF"),
		TextWrap:   style.WrapWord,
	})

	// Ellipsis
	draw.DrawStyledTextInWidth(buf, 0, 18, "Text with ellipsis overflow", 15, style.Style{
		Foreground:   style.Color("#FFFFFF"),
		TextOverflow: style.TextOverflowEllipsis,
	})

	term := render.NewTerminal(buf.Width, buf.Height)
	assert.Golden(t, term.RenderFull(buf))
}
