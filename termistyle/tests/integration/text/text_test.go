package text_test

import (
	"testing"

	"github.com/dottermi/termitest/assert"
	"github.com/dottermi/x/render"
	"github.com/dottermi/x/termistyle/draw"
	"github.com/dottermi/x/termistyle/style"
	. "github.com/dottermi/x/termistyle/tests/integration/helper"
)

func TestGolden_Text(t *testing.T) {
	buf := render.NewBuffer(30, 20)

	// Basic text with colors
	draw.Text(buf, 0, 0, "Basic Text", ColorText, ColorBg)
	draw.Text(buf, 0, 1, "Cyan", ColorCyan, ColorBg)
	draw.Text(buf, 0, 2, "Cyan2", ColorCyan2, ColorBg)
	draw.Text(buf, 0, 3, "Blue", ColorBlue, ColorBg)

	// Bold
	draw.StyledText(buf, 0, 5, "Bold", style.Style{
		Foreground: ColorText,
		FontWeight: style.WeightBold,
	})

	// Italic
	draw.StyledText(buf, 10, 5, "Italic", style.Style{
		Foreground: ColorText,
		FontStyle:  style.StyleItalic,
	})

	// Underline
	draw.StyledText(buf, 0, 6, "Underline", style.Style{
		Foreground:     ColorText,
		TextDecoration: style.DecorationUnderline,
	})

	// Strikethrough
	draw.StyledText(buf, 12, 6, "Strike", style.Style{
		Foreground:     ColorText,
		TextDecoration: style.DecorationLineThrough,
	})

	// Dim and Reverse
	draw.StyledText(buf, 0, 7, "Dim", style.Style{
		Foreground: ColorText,
		Dim:        true,
	})
	draw.StyledText(buf, 10, 7, "Reverse", style.Style{
		Foreground: ColorText,
		Reverse:    true,
	})

	// Text transform
	draw.StyledText(buf, 0, 9, "hello", style.Style{
		Foreground:    ColorText,
		TextTransform: style.TransformUppercase,
	})
	draw.StyledText(buf, 10, 9, "WORLD", style.Style{
		Foreground:    ColorText,
		TextTransform: style.TransformLowercase,
	})

	// Text align
	draw.StyledTextInWidth(buf, 0, 11, "Left", 15, style.Style{
		Foreground: ColorMuted,
		TextAlign:  style.TextAlignLeft,
	})
	draw.StyledTextInWidth(buf, 0, 12, "Center", 15, style.Style{
		Foreground: ColorMuted,
		TextAlign:  style.TextAlignCenter,
	})
	draw.StyledTextInWidth(buf, 0, 13, "Right", 15, style.Style{
		Foreground: ColorMuted,
		TextAlign:  style.TextAlignRight,
	})

	// Word wrap
	draw.StyledTextInBox(buf, 0, 15, "Word wrap test here", 10, 2, style.Style{
		Foreground: ColorCyan,
		TextWrap:   style.WrapWord,
	})

	// Ellipsis
	draw.StyledTextInWidth(buf, 0, 18, "Text with ellipsis overflow", 15, style.Style{
		Foreground:   ColorAccent,
		TextOverflow: style.TextOverflowEllipsis,
	})

	term := render.NewTerminal(buf.Width, buf.Height)
	assert.Golden(t, term.RenderFull(buf))
}
