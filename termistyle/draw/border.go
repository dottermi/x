package draw

import "github.com/dottermi/x/termistyle/style"

// DrawBorder renders a rectangular border using box-drawing characters.
// Each side can have a different style and color.
//
// Parameters:
//   - x, y: top-left corner position
//   - width, height: outer dimensions including the border
//   - border: border configuration for each side
//
// Example:
//
//	buf.DrawBorder(0, 0, 20, 10, style.BorderAll(style.BorderRound, style.Color("#FFFFFF")))
func (b *Buffer) DrawBorder(x, y, width, height int, border style.Border) {
	b.DrawBorderClipped(x, y, width, height, border, ClipRect{X: 0, Y: 0, W: b.Width, H: b.Height})
}

// DrawBorderClipped renders a rectangular border with clipping.
func (b *Buffer) DrawBorderClipped(x, y, width, height int, border style.Border, clip ClipRect) {
	if !border.HasAny() || width < 2 || height < 2 {
		return
	}

	d := &borderDrawerClipped{
		buf:       b,
		clip:      clip,
		x:         x,
		y:         y,
		width:     width,
		height:    height,
		hasTop:    border.Top.IsSet(),
		hasBottom: border.Bottom.IsSet(),
		hasLeft:   border.Left.IsSet(),
		hasRight:  border.Right.IsSet(),
		border:    border,
	}

	d.drawCorners()

	if d.hasTop {
		d.drawHorizontal(border.Top, y)
	}
	if d.hasBottom {
		d.drawHorizontal(border.Bottom, y+height-1)
	}
	if d.hasLeft {
		d.drawVertical(border.Left, x)
	}
	if d.hasRight {
		d.drawVertical(border.Right, x+width-1)
	}
}

// borderDrawerClipped holds state for drawing a border with clipping.
type borderDrawerClipped struct {
	buf                                  *Buffer
	clip                                 ClipRect
	x, y, width, height                  int
	hasTop, hasBottom, hasLeft, hasRight bool
	border                               style.Border
}

func (d *borderDrawerClipped) cell(ch rune, color style.Color) Cell {
	return Cell{Char: ch, Foreground: color}
}

func (d *borderDrawerClipped) drawCorners() {
	if d.hasTop && d.hasLeft {
		if chars, ok := style.Borders[d.border.Top.Style]; ok {
			d.buf.SetClipped(d.x, d.y, d.cell(chars.TopLeft, d.border.Top.Color), d.clip)
		}
	}
	if d.hasTop && d.hasRight {
		if chars, ok := style.Borders[d.border.Top.Style]; ok {
			d.buf.SetClipped(d.x+d.width-1, d.y, d.cell(chars.TopRight, d.border.Top.Color), d.clip)
		}
	}
	if d.hasBottom && d.hasLeft {
		if chars, ok := style.Borders[d.border.Bottom.Style]; ok {
			d.buf.SetClipped(d.x, d.y+d.height-1, d.cell(chars.BottomLeft, d.border.Bottom.Color), d.clip)
		}
	}
	if d.hasBottom && d.hasRight {
		if chars, ok := style.Borders[d.border.Bottom.Style]; ok {
			d.buf.SetClipped(d.x+d.width-1, d.y+d.height-1, d.cell(chars.BottomRight, d.border.Bottom.Color), d.clip)
		}
	}
}

func (d *borderDrawerClipped) drawHorizontal(edge style.BorderEdge, posY int) {
	chars, ok := style.Borders[edge.Style]
	if !ok {
		return
	}
	startX, endX := d.x, d.x+d.width
	if d.hasLeft {
		startX++
	}
	if d.hasRight {
		endX--
	}

	// Get all texts (from Texts slice or legacy Title field)
	texts := edge.GetTexts()
	if len(texts) == 0 {
		// No texts, just draw horizontal line
		for x := startX; x < endX; x++ {
			d.buf.SetClipped(x, posY, d.cell(chars.Horizontal, edge.Color), d.clip)
		}
		return
	}

	// Draw with texts
	d.drawHorizontalWithTexts(edge, posY, startX, endX, chars, texts)
}

// renderedText holds a text segment positioned on the border line.
type renderedText struct {
	runes    []rune
	startPos int // relative to startX
	style    style.Style
}

func (d *borderDrawerClipped) drawHorizontalWithTexts(
	edge style.BorderEdge, posY, startX, endX int,
	chars style.BorderChars, texts []style.BorderText,
) {
	availableWidth := endX - startX
	if availableWidth < 3 {
		// Too narrow for any text
		for x := startX; x < endX; x++ {
			d.buf.SetClipped(x, posY, d.cell(chars.Horizontal, edge.Color), d.clip)
		}
		return
	}

	// Group texts by alignment
	var leftTexts, centerTexts, rightTexts []style.BorderText
	for _, t := range texts {
		switch t.Align {
		case style.TextAlignLeft:
			leftTexts = append(leftTexts, t)
		case style.TextAlignCenter:
			centerTexts = append(centerTexts, t)
		case style.TextAlignRight:
			rightTexts = append(rightTexts, t)
		default:
			leftTexts = append(leftTexts, t)
		}
	}

	// Calculate positions for all texts
	rendered := calculateTextPositions(leftTexts, centerTexts, rightTexts, availableWidth)

	// Build a map of positions to cells
	charMap := make(map[int]Cell)
	for _, rt := range rendered {
		fg := rt.style.Foreground
		if !fg.IsSet() {
			fg = edge.Color
		}
		for i, r := range rt.runes {
			pos := rt.startPos + i
			if pos >= 0 && pos < availableWidth {
				charMap[pos] = Cell{
					Char:       r,
					Foreground: fg,
					Background: rt.style.Background,
					Bold:       rt.style.FontWeight.IsBold(),
					Italic:     rt.style.FontStyle.IsItalic(),
					Underline:  rt.style.TextDecoration.HasUnderline(),
					Strike:     rt.style.TextDecoration.HasLineThrough(),
					Dim:        rt.style.Dim,
					Reverse:    rt.style.Reverse,
				}
			}
		}
	}

	// Draw the line with texts
	for x := startX; x < endX; x++ {
		relPos := x - startX
		if cell, ok := charMap[relPos]; ok {
			d.buf.SetClipped(x, posY, cell, d.clip)
		} else {
			d.buf.SetClipped(x, posY, d.cell(chars.Horizontal, edge.Color), d.clip)
		}
	}
}

func calculateTextPositions(
	left, center, right []style.BorderText,
	width int,
) []renderedText {
	leftWidth, rightWidth := calculateZoneWidths(left, right, width)

	var result []renderedText
	result = placeLeftTexts(result, left, leftWidth)
	result = placeRightTexts(result, right, width, leftWidth)
	result = placeCenterTexts(result, center, leftWidth, width-rightWidth)

	return result
}

func calculateZoneWidths(left, right []style.BorderText, width int) (leftWidth, rightWidth int) {
	for _, t := range left {
		leftWidth += len([]rune(" " + t.Text + " "))
	}
	for _, t := range right {
		rightWidth += len([]rune(" " + t.Text + " "))
	}
	// Handle collision - proportionally reduce both sides
	if leftWidth+rightWidth > width {
		total := leftWidth + rightWidth
		leftWidth = (leftWidth * width) / total
		rightWidth = width - leftWidth
	}
	return leftWidth, rightWidth
}

func placeLeftTexts(result []renderedText, texts []style.BorderText, maxWidth int) []renderedText {
	pos := 1
	for _, t := range texts {
		paddedText := " " + t.Text + " "
		runes := []rune(paddedText)
		maxLen := maxWidth - pos + 1
		if maxLen > 0 && len(runes) > maxLen {
			runes = truncateRunes(runes, maxLen)
		}
		if len(runes) > 0 {
			result = append(result, renderedText{runes: runes, startPos: pos, style: t.Style})
			pos += len(runes)
		}
	}
	return result
}

func placeRightTexts(result []renderedText, texts []style.BorderText, width, leftWidth int) []renderedText {
	rightPos := width - 1
	for i := len(texts) - 1; i >= 0; i-- {
		t := texts[i]
		paddedText := " " + t.Text + " "
		runes := []rune(paddedText)
		textStart := rightPos - len(runes)
		if textStart < leftWidth {
			maxLen := rightPos - leftWidth
			if maxLen <= 0 {
				continue
			}
			runes = truncateRunesFromStart(runes, maxLen)
			textStart = rightPos - len(runes)
		}
		if len(runes) > 0 {
			result = append(result, renderedText{runes: runes, startPos: textStart, style: t.Style})
			rightPos = textStart
		}
	}
	return result
}

func placeCenterTexts(result []renderedText, texts []style.BorderText, centerStart, centerEnd int) []renderedText {
	centerWidth := centerEnd - centerStart
	for _, t := range texts {
		paddedText := " " + t.Text + " "
		runes := []rune(paddedText)
		if len(runes) > centerWidth {
			runes = truncateRunes(runes, centerWidth)
		}
		if len(runes) > 0 && centerWidth > 0 {
			textStart := centerStart + (centerWidth-len(runes))/2
			result = append(result, renderedText{runes: runes, startPos: textStart, style: t.Style})
		}
	}
	return result
}

func truncateRunes(runes []rune, maxLen int) []rune {
	if len(runes) <= maxLen {
		return runes
	}
	if maxLen <= 3 {
		return runes[:maxLen]
	}
	return append(runes[:maxLen-3], '.', '.', '.')
}

func truncateRunesFromStart(runes []rune, maxLen int) []rune {
	if len(runes) <= maxLen {
		return runes
	}
	if maxLen <= 3 {
		return runes[len(runes)-maxLen:]
	}
	return append([]rune{'.', '.', '.'}, runes[len(runes)-maxLen+3:]...)
}

func (d *borderDrawerClipped) drawVertical(edge style.BorderEdge, posX int) {
	chars, ok := style.Borders[edge.Style]
	if !ok {
		return
	}
	startY, endY := d.y, d.y+d.height
	if d.hasTop {
		startY++
	}
	if d.hasBottom {
		endY--
	}
	for y := startY; y < endY; y++ {
		d.buf.SetClipped(posX, y, d.cell(chars.Vertical, edge.Color), d.clip)
	}
}
