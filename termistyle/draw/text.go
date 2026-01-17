package draw

import (
	"strings"

	"github.com/dottermi/x/render"
	"github.com/dottermi/x/termistyle/style"
)

// Text renders a string horizontally starting at the given position.
// Each rune occupies one cell; text extending beyond buffer bounds is clipped.
//
// Example:
//
//	draw.Text(buf, 5, 2, "Hello", fg, bg)
func Text(b *render.Buffer, x, y int, text string, fg, bg style.Color) {
	for i, ch := range []rune(text) {
		b.Set(x+i, y, render.Cell{Char: ch, FG: fg.ToRender(), BG: bg.ToRender()})
	}
}

// StyledText renders a string with full text styling support.
// Applies text-transform, font-weight, font-style, text-decoration, dim, and reverse.
//
// Example:
//
//	draw.StyledText(buf, 5, 2, "Hello", s)
func StyledText(b *render.Buffer, x, y int, text string, s style.Style) {
	StyledTextInBox(b, x, y, text, 0, 1, s)
}

// StyledTextInWidth renders a string with text alignment within a container width.
// If containerWidth is 0, text is drawn at x position without alignment.
// Applies text-transform, font-weight, font-style, text-decoration, dim, and reverse.
//
// Example:
//
//	draw.StyledTextInWidth(buf, 5, 2, "Hello", 20, s) // aligns within 20 chars
func StyledTextInWidth(b *render.Buffer, x, y int, text string, containerWidth int, s style.Style) {
	StyledTextInBox(b, x, y, text, containerWidth, 1, s)
}

// StyledTextInBox renders text with wrapping, alignment within a box.
// containerWidth of 0 disables wrapping and alignment.
// containerHeight limits the number of lines rendered.
//
// Example:
//
//	draw.StyledTextInBox(buf, 5, 2, "Long text here", 20, 5, s)
func StyledTextInBox(b *render.Buffer, x, y int, text string, containerWidth, containerHeight int, s style.Style) {
	StyledTextInBoxClipped(b, x, y, text, containerWidth, containerHeight, s, render.ClipRect{X: 0, Y: 0, W: b.Width, H: b.Height})
}

// StyledTextInBoxClipped renders text with wrapping, alignment, and clipping.
func StyledTextInBoxClipped(b *render.Buffer, x, y int, text string, containerWidth, containerHeight int, s style.Style, clip render.ClipRect) {
	text = s.TextTransform.Apply(text)

	lines, needsEllipsis := prepareTextLines(text, containerWidth, containerHeight, s.TextWrap)

	cell := render.Cell{
		FG:        s.Foreground.ToRender(),
		BG:        s.Background.ToRender(),
		Bold:      s.FontWeight.IsBold(),
		Italic:    s.FontStyle.IsItalic(),
		Underline: s.TextDecoration.HasUnderline(),
		Strike:    s.TextDecoration.HasLineThrough(),
		Dim:       s.Dim,
		Reverse:   s.Reverse,
	}

	for lineIdx, line := range lines {
		runes := []rune(line)
		isLastLine := lineIdx == len(lines)-1

		runes = applyTextOverflow(runes, containerWidth, isLastLine, needsEllipsis, s.TextOverflow)
		offsetX := calculateTextOffset(len(runes), containerWidth, s.TextAlign)

		for i, ch := range runes {
			cell.Char = ch
			b.SetClipped(x+offsetX+i, y+lineIdx, cell, clip)
		}
	}
}

// prepareTextLines wraps and truncates text into lines.
func prepareTextLines(text string, containerWidth, containerHeight int, wrap style.TextWrap) ([]string, bool) {
	var lines []string
	if containerWidth > 0 && wrap != style.WrapNone {
		lines = wrapText(text, containerWidth, wrap)
	} else {
		lines = []string{text}
	}

	needsEllipsis := false
	if containerHeight > 0 && len(lines) > containerHeight {
		needsEllipsis = true
		lines = lines[:containerHeight]
	}
	return lines, needsEllipsis
}

// applyTextOverflow applies ellipsis truncation if needed.
func applyTextOverflow(runes []rune, containerWidth int, isLastLine, needsEllipsis bool, overflow style.TextOverflow) []rune {
	if overflow != style.TextOverflowEllipsis || containerWidth <= 0 {
		return runes
	}
	lineLen := len(runes)
	if lineLen > containerWidth {
		return applyEllipsis(runes, containerWidth)
	}
	if isLastLine && needsEllipsis && lineLen > 0 {
		return applyEllipsis(runes, containerWidth)
	}
	return runes
}

// calculateTextOffset calculates x offset for text alignment.
func calculateTextOffset(lineLen, containerWidth int, align style.TextAlign) int {
	if containerWidth <= 0 || lineLen >= containerWidth {
		return 0
	}
	switch align {
	case style.TextAlignLeft:
		return 0
	case style.TextAlignCenter:
		return (containerWidth - lineLen) / 2
	case style.TextAlignRight:
		return containerWidth - lineLen
	}
	return 0
}

// applyEllipsis truncates runes to fit maxWidth with "..." at the end.
func applyEllipsis(runes []rune, maxWidth int) []rune {
	if len(runes) <= maxWidth {
		return runes
	}
	if maxWidth <= 3 {
		return []rune("...")[:maxWidth]
	}
	return append(runes[:maxWidth-3], '.', '.', '.')
}

// wrapText breaks text into lines that fit within maxWidth.
func wrapText(text string, maxWidth int, wrapMode style.TextWrap) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	if wrapMode == style.WrapChar {
		return wrapByChar(text, maxWidth)
	}

	return wrapByWord(text, maxWidth)
}

// wrapByWord wraps text at word boundaries.
func wrapByWord(text string, maxWidth int) []string {
	var lines []string
	words := strings.Fields(text)

	if len(words) == 0 {
		return []string{""}
	}

	currentLine := ""
	for _, word := range words {
		wordLen := len([]rune(word))

		// If word is longer than maxWidth, break it
		if wordLen > maxWidth {
			// Flush current line if not empty
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = ""
			}
			// Break the long word
			lines = append(lines, wrapByChar(word, maxWidth)...)
			continue
		}

		switch {
		case currentLine == "":
			currentLine = word
		case len([]rune(currentLine))+1+wordLen <= maxWidth:
			currentLine += " " + word
		default:
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// wrapByChar wraps text at character boundaries.
func wrapByChar(text string, maxWidth int) []string {
	var lines []string
	runes := []rune(text)

	for len(runes) > 0 {
		end := maxWidth
		if end > len(runes) {
			end = len(runes)
		}
		lines = append(lines, string(runes[:end]))
		runes = runes[end:]
	}

	return lines
}
