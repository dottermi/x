package layout

import "github.com/dottermi/x/termistyle/style"

// flowChildren returns children that are in normal flow (not absolute positioned).
func (b *Box) flowChildren() []*Box {
	var result []*Box
	for _, child := range b.Children {
		if child.Style.Position != style.Absolute {
			result = append(result, child)
		}
	}
	return result
}

// calculateFlexRow positions children horizontally with justify and align.
// Respects each child's margin for positioning.
func (b *Box) calculateFlexRow(startX, startY, innerW, innerH int) {
	children := b.flowChildren()
	if len(children) == 0 {
		return
	}

	// Apply FlexBasis as initial width for children
	b.applyFlexBasisRow(children)

	// Stretch height FIRST (before calculating children without FlexGrow)
	// Account for vertical margins when stretching
	for _, child := range children {
		if child.Style.Height == 0 {
			verticalMargin := child.Style.Margin.Top + child.Style.Margin.Bottom
			child.Style.Height = innerH - verticalMargin
			if child.Style.Height < 0 {
				child.Style.Height = 0
			}
		}
	}

	// First pass: distribute grow/shrink space to set child dimensions BEFORE they calculate
	b.distributeFlexSpaceRow(children, innerW)

	totalWidth := totalChildrenWidth(children)
	totalGap := b.Style.Gap * (len(children) - 1)

	x, spacing := b.justifyRowX(startX, innerW, totalWidth, totalGap, len(children))

	for _, child := range children {
		// Apply left margin to X position
		child.X = b.X + x + child.marginLeft()
		// Apply top margin to Y position
		childY, childH := b.alignRowY(startY, innerH, child.Style.Height+child.Style.Margin.Top+child.Style.Margin.Bottom)
		child.Y = childY + child.marginTop()
		child.Style.Height = childH - child.Style.Margin.Top - child.Style.Margin.Bottom
		if child.Style.Height < 0 {
			child.Style.Height = 0
		}

		child.Calculate()

		// Advance by outer width (including margins) plus spacing
		x += child.outerWidth() + b.childRowSpacing(spacing)
	}
}

// applyFlexBasisRow sets initial widths from FlexBasis when specified.
func (b *Box) applyFlexBasisRow(children []*Box) {
	for _, child := range children {
		if child.Style.FlexBasis > 0 {
			child.Style.Width = child.Style.FlexBasis
		}
	}
}

// distributeFlexSpaceRow distributes extra space (grow) or reduces space (shrink).
// Accounts for children's horizontal margins when calculating available space.
func (b *Box) distributeFlexSpaceRow(children []*Box, innerW int) {
	totalGap := b.Style.Gap * (len(children) - 1)

	// Calculate total horizontal margins
	totalMargins := 0
	for _, child := range children {
		totalMargins += child.Style.Margin.Left + child.Style.Margin.Right
	}

	// First, calculate widths for children WITHOUT FlexGrow
	// so we know how much space they actually need
	baseWidth := 0
	for _, child := range children {
		if child.Style.FlexGrow == 0 {
			// Temporarily calculate to get actual width
			child.Calculate()
			baseWidth += child.W
		} else {
			baseWidth += child.Style.Width
		}
	}

	// Available space excludes gaps and margins
	availableSpace := innerW - baseWidth - totalGap - totalMargins

	if availableSpace > 0 {
		b.distributeGrowRow(children, innerW-totalMargins, baseWidth)
	} else if availableSpace < 0 {
		b.distributeShrinkRow(children, innerW-totalMargins)
	}
}

// distributeGrowRow distributes extra horizontal space among children with FlexGrow > 0.
func (b *Box) distributeGrowRow(children []*Box, innerW, baseWidth int) {
	totalGap := b.Style.Gap * (len(children) - 1)

	// Calculate total grow
	var totalGrow float64
	for _, child := range children {
		totalGrow += child.Style.FlexGrow
	}

	if totalGrow == 0 {
		return // No grow, nothing to distribute
	}

	extraSpace := innerW - baseWidth - totalGap
	if extraSpace <= 0 {
		return // No extra space to distribute
	}

	// Distribute extra space proportionally, tracking remainder for rounding
	distributed := 0
	var lastGrowChild *Box
	for _, child := range children {
		if child.Style.FlexGrow > 0 {
			extra := int(float64(extraSpace) * (child.Style.FlexGrow / totalGrow))
			child.Style.Width += extra
			distributed += extra
			lastGrowChild = child
		}
	}

	// Give any remaining pixels to the last growing child to avoid gaps
	if lastGrowChild != nil && distributed < extraSpace {
		lastGrowChild.Style.Width += extraSpace - distributed
	}
}

// distributeShrinkRow reduces widths when container is smaller than children need.
// FlexShrink defaults to 1 when not set (CSS behavior).
func (b *Box) distributeShrinkRow(children []*Box, innerW int) {
	totalGap := b.Style.Gap * (len(children) - 1)

	// Calculate total shrink factor and base width
	var totalShrink float64
	var baseWidth int
	for _, child := range children {
		shrink := child.Style.FlexShrink
		if shrink == 0 {
			shrink = 1 // Default shrink is 1
		}
		totalShrink += shrink * float64(child.Style.Width)
		baseWidth += child.Style.Width
	}

	if totalShrink == 0 {
		return // No shrink, nothing to reduce
	}

	overflow := baseWidth + totalGap - innerW
	if overflow <= 0 {
		return // No overflow to handle
	}

	// Distribute shrink proportionally (weighted by width * shrink factor)
	for _, child := range children {
		shrink := child.Style.FlexShrink
		if shrink == 0 {
			shrink = 1
		}
		shrinkFactor := shrink * float64(child.Style.Width)
		reduction := int(float64(overflow) * (shrinkFactor / totalShrink))
		child.Style.Width -= reduction
		if child.Style.Width < 0 {
			child.Style.Width = 0
		}
	}
}

// totalChildrenWidth sums the widths of all children including horizontal margins.
func totalChildrenWidth(children []*Box) int {
	total := 0
	for _, child := range children {
		total += child.Style.Width + child.Style.Margin.Left + child.Style.Margin.Right
	}
	return total
}

// justifyRowX computes the starting X and spacing for horizontal justify.
func (b *Box) justifyRowX(startX, innerW, totalWidth, totalGap, numChildren int) (x, spacing int) {
	if b.Style.JustifyContent == style.JustifyCenter {
		return startX + (innerW-totalWidth-totalGap)/2, 0
	}
	if b.Style.JustifyContent == style.JustifyEnd {
		return startX + innerW - totalWidth - totalGap, 0
	}
	if b.Style.JustifyContent == style.JustifyBetween && numChildren > 1 {
		return startX, (innerW - totalWidth) / (numChildren - 1)
	}
	if b.Style.JustifyContent == style.JustifyAround && numChildren > 0 {
		spacing := (innerW - totalWidth) / numChildren
		return startX + spacing/2, spacing
	}

	return startX, 0
}

// alignRowY computes Y position and height for vertical alignment.
func (b *Box) alignRowY(startY, innerH, childH int) (y, height int) {
	if b.Style.AlignItems == style.AlignCenter {
		return b.Y + startY + (innerH-childH)/2, childH
	}
	if b.Style.AlignItems == style.AlignEnd {
		return b.Y + startY + innerH - childH, childH
	}
	if b.Style.AlignItems == style.AlignStretch {
		return b.Y + startY, innerH
	}

	return b.Y + startY, childH
}

// childRowSpacing returns the horizontal gap between children.
func (b *Box) childRowSpacing(spacing int) int {
	if b.Style.JustifyContent == style.JustifyBetween || b.Style.JustifyContent == style.JustifyAround {
		return spacing
	}
	return b.Style.Gap
}

// calculateFlexRowWrap positions children horizontally with wrapping to multiple lines.
// Respects each child's margin for positioning and line splitting.
func (b *Box) calculateFlexRowWrap(startX, startY, innerW, innerH int) {
	children := b.flowChildren()
	if len(children) == 0 {
		return
	}

	// Apply FlexBasis as initial width for children
	b.applyFlexBasisRow(children)

	// Split children into lines based on available width (considering margins)
	lines := splitIntoRowLines(children, innerW, b.Style.Gap)
	if len(lines) == 0 {
		return
	}

	// Calculate line heights (including vertical margins)
	lineHeights := make([]int, len(lines))
	for i, line := range lines {
		maxH := 0
		for _, child := range line {
			outerH := child.Style.Height + child.Style.Margin.Top + child.Style.Margin.Bottom
			if outerH > maxH {
				maxH = outerH
			}
		}
		if maxH == 0 {
			maxH = 1 // Minimum height
		}
		lineHeights[i] = maxH
	}

	// Position each line
	y := startY
	for i, line := range lines {
		lineH := lineHeights[i]

		// Calculate total width including margins and gap for this line
		totalWidth := 0
		for _, child := range line {
			totalWidth += child.Style.Width + child.Style.Margin.Left + child.Style.Margin.Right
		}
		totalGap := b.Style.Gap * (len(line) - 1)

		x, spacing := b.justifyRowXForLine(startX, innerW, totalWidth, totalGap, len(line))

		for _, child := range line {
			// Apply margins to position
			child.X = b.X + x + child.marginLeft()
			childY, childH := b.alignRowYForLine(y, lineH, child.Style.Height+child.Style.Margin.Top+child.Style.Margin.Bottom)
			child.Y = childY + child.marginTop()
			child.Style.Height = childH - child.Style.Margin.Top - child.Style.Margin.Bottom
			if child.Style.Height < 0 {
				child.Style.Height = 0
			}

			child.Calculate()

			// Advance by outer width plus spacing
			x += child.outerWidth() + b.childRowSpacing(spacing)
		}

		y += lineH + b.Style.Gap
	}
}

// splitIntoRowLines divides children into lines that fit within innerW.
// Considers children's horizontal margins when calculating line widths.
func splitIntoRowLines(children []*Box, innerW, gap int) [][]*Box {
	var lines [][]*Box
	var currentLine []*Box
	currentWidth := 0

	for _, child := range children {
		// Use outer width (including margins)
		childW := child.Style.Width + child.Style.Margin.Left + child.Style.Margin.Right
		if childW == 0 {
			childW = 1 // Minimum width
		}

		gapNeeded := 0
		if len(currentLine) > 0 {
			gapNeeded = gap
		}

		// Check if child fits in current line
		if len(currentLine) > 0 && currentWidth+gapNeeded+childW > innerW {
			// Start new line
			lines = append(lines, currentLine)
			currentLine = []*Box{child}
			currentWidth = childW
		} else {
			// Add to current line
			currentLine = append(currentLine, child)
			currentWidth += gapNeeded + childW
		}
	}

	// Don't forget the last line
	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

// justifyRowXForLine computes starting X for a specific line.
func (b *Box) justifyRowXForLine(startX, innerW, totalWidth, totalGap, numChildren int) (x, spacing int) {
	if b.Style.JustifyContent == style.JustifyCenter {
		return startX + (innerW-totalWidth-totalGap)/2, 0
	}
	if b.Style.JustifyContent == style.JustifyEnd {
		return startX + innerW - totalWidth - totalGap, 0
	}
	if b.Style.JustifyContent == style.JustifyBetween && numChildren > 1 {
		return startX, (innerW - totalWidth) / (numChildren - 1)
	}
	if b.Style.JustifyContent == style.JustifyAround && numChildren > 0 {
		spacing := (innerW - totalWidth) / numChildren
		return startX + spacing/2, spacing
	}

	return startX, 0
}

// alignRowYForLine computes Y position within a line.
func (b *Box) alignRowYForLine(lineY, lineH, childH int) (y, height int) {
	if b.Style.AlignItems == style.AlignCenter {
		return b.Y + lineY + (lineH-childH)/2, childH
	}
	if b.Style.AlignItems == style.AlignEnd {
		return b.Y + lineY + lineH - childH, childH
	}
	if b.Style.AlignItems == style.AlignStretch {
		return b.Y + lineY, lineH
	}

	return b.Y + lineY, childH
}
