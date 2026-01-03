package layout

import "github.com/dottermi/x/termistyle/style"

// calculateFlexColumn positions children vertically with justify and align.
// Respects each child's margin for positioning.
func (b *Box) calculateFlexColumn(startX, startY, innerW, innerH int) {
	children := b.flowChildren()
	if len(children) == 0 {
		return
	}

	// Apply FlexBasis as initial height for children
	applyFlexBasisColumn(children)

	// Stretch width FIRST (before calculating children without FlexGrow)
	// Account for horizontal margins when stretching
	for _, child := range children {
		if child.Style.Width == 0 {
			horizontalMargin := child.Style.Margin.Left + child.Style.Margin.Right
			child.Style.Width = innerW - horizontalMargin
			if child.Style.Width < 0 {
				child.Style.Width = 0
			}
		}
	}

	// First pass: distribute grow/shrink space to set child dimensions BEFORE they calculate
	b.distributeFlexSpaceColumn(children, innerH)

	totalHeight := totalChildrenHeight(children)
	totalGap := b.Style.Gap * (len(children) - 1)

	y, spacing := b.justifyColY(startY, innerH, totalHeight, totalGap, len(children))

	for _, child := range children {
		// Apply top margin to Y position
		child.Y = b.Y + y + child.marginTop()
		// Apply left margin to X position
		childX, childW := b.alignColX(startX, innerW, child.Style.Width+child.Style.Margin.Left+child.Style.Margin.Right)
		child.X = childX + child.marginLeft()
		child.Style.Width = childW - child.Style.Margin.Left - child.Style.Margin.Right
		if child.Style.Width < 0 {
			child.Style.Width = 0
		}

		child.Calculate()

		// Advance by outer height (including margins) plus spacing
		y += child.outerHeight() + b.childColSpacing(spacing)
	}
}

// applyFlexBasisColumn sets initial heights from FlexBasis when specified.
func applyFlexBasisColumn(children []*Box) {
	for _, child := range children {
		if child.Style.FlexBasis > 0 {
			child.Style.Height = child.Style.FlexBasis
		}
	}
}

// distributeFlexSpaceColumn distributes extra space (grow) or reduces space (shrink).
// Accounts for children's vertical margins when calculating available space.
func (b *Box) distributeFlexSpaceColumn(children []*Box, innerH int) {
	totalGap := b.Style.Gap * (len(children) - 1)

	// Calculate total vertical margins
	totalMargins := 0
	for _, child := range children {
		totalMargins += child.Style.Margin.Top + child.Style.Margin.Bottom
	}

	// First, calculate heights for children WITHOUT FlexGrow
	// so we know how much space they actually need
	baseHeight := 0
	for _, child := range children {
		if child.Style.FlexGrow == 0 {
			// Temporarily calculate to get actual height
			child.Calculate()
			baseHeight += child.H
		} else {
			baseHeight += child.Style.Height
		}
	}

	// Available space excludes gaps and margins
	availableSpace := innerH - baseHeight - totalGap - totalMargins

	if availableSpace > 0 {
		b.distributeGrowColumn(children, innerH-totalMargins, baseHeight)
	} else if availableSpace < 0 {
		b.distributeShrinkColumn(children, innerH-totalMargins)
	}
}

// distributeGrowColumn distributes extra vertical space among children with FlexGrow > 0.
func (b *Box) distributeGrowColumn(children []*Box, innerH, baseHeight int) {
	totalGap := b.Style.Gap * (len(children) - 1)

	// Calculate total grow
	var totalGrow float64
	for _, child := range children {
		totalGrow += child.Style.FlexGrow
	}

	if totalGrow == 0 {
		return // No grow, nothing to distribute
	}

	extraSpace := innerH - baseHeight - totalGap
	if extraSpace <= 0 {
		return // No extra space to distribute
	}

	// Distribute extra space proportionally, tracking remainder for rounding
	distributed := 0
	var lastGrowChild *Box
	for _, child := range children {
		if child.Style.FlexGrow > 0 {
			extra := int(float64(extraSpace) * (child.Style.FlexGrow / totalGrow))
			child.Style.Height += extra
			distributed += extra
			lastGrowChild = child
		}
	}

	// Give any remaining pixels to the last growing child to avoid gaps
	if lastGrowChild != nil && distributed < extraSpace {
		lastGrowChild.Style.Height += extraSpace - distributed
	}
}

// distributeShrinkColumn reduces heights when container is smaller than children need.
// FlexShrink defaults to 1 when not set (CSS behavior).
func (b *Box) distributeShrinkColumn(children []*Box, innerH int) {
	totalGap := b.Style.Gap * (len(children) - 1)

	// Calculate total shrink factor and base height
	var totalShrink float64
	var baseHeight int
	for _, child := range children {
		shrink := child.Style.FlexShrink
		if shrink == 0 {
			shrink = 1 // Default shrink is 1
		}
		totalShrink += shrink * float64(child.Style.Height)
		baseHeight += child.Style.Height
	}

	if totalShrink == 0 {
		return // No shrink, nothing to reduce
	}

	overflow := baseHeight + totalGap - innerH
	if overflow <= 0 {
		return // No overflow to handle
	}

	// Distribute shrink proportionally (weighted by height * shrink factor)
	for _, child := range children {
		shrink := child.Style.FlexShrink
		if shrink == 0 {
			shrink = 1
		}
		shrinkFactor := shrink * float64(child.Style.Height)
		reduction := int(float64(overflow) * (shrinkFactor / totalShrink))
		child.Style.Height -= reduction
		if child.Style.Height < 0 {
			child.Style.Height = 0
		}
	}
}

// totalChildrenHeight sums the heights of all children including vertical margins.
func totalChildrenHeight(children []*Box) int {
	total := 0
	for _, child := range children {
		total += child.Style.Height + child.Style.Margin.Top + child.Style.Margin.Bottom
	}
	return total
}

// justifyColY computes the starting Y and spacing for vertical justify.
func (b *Box) justifyColY(startY, innerH, totalHeight, totalGap, numChildren int) (y, spacing int) {
	if b.Style.JustifyContent == style.JustifyCenter {
		return startY + (innerH-totalHeight-totalGap)/2, 0
	}
	if b.Style.JustifyContent == style.JustifyEnd {
		return startY + innerH - totalHeight - totalGap, 0
	}
	if b.Style.JustifyContent == style.JustifyBetween && numChildren > 1 {
		return startY, (innerH - totalHeight) / (numChildren - 1)
	}
	if b.Style.JustifyContent == style.JustifyAround && numChildren > 0 {
		spacing := (innerH - totalHeight) / numChildren
		return startY + spacing/2, spacing
	}

	return startY, 0
}

// alignColX computes X position and width for horizontal alignment.
func (b *Box) alignColX(startX, innerW, childW int) (x, width int) {
	if b.Style.AlignItems == style.AlignCenter {
		return b.X + startX + (innerW-childW)/2, childW
	}
	if b.Style.AlignItems == style.AlignEnd {
		return b.X + startX + innerW - childW, childW
	}
	if b.Style.AlignItems == style.AlignStretch {
		return b.X + startX, innerW
	}

	return b.X + startX, childW
}

// childColSpacing returns the vertical gap between children.
func (b *Box) childColSpacing(spacing int) int {
	if b.Style.JustifyContent == style.JustifyBetween || b.Style.JustifyContent == style.JustifyAround {
		return spacing
	}
	return b.Style.Gap
}

// calculateFlexColumnWrap positions children vertically with wrapping to multiple columns.
// Respects each child's margin for positioning and column splitting.
func (b *Box) calculateFlexColumnWrap(startX, startY, innerW, innerH int) {
	children := b.flowChildren()
	if len(children) == 0 {
		return
	}

	// Apply FlexBasis as initial height for children
	applyFlexBasisColumn(children)

	// Split children into columns based on available height (considering margins)
	columns := splitIntoColumnLines(children, innerH, b.Style.Gap)
	if len(columns) == 0 {
		return
	}

	// Calculate column widths (including horizontal margins)
	colWidths := make([]int, len(columns))
	for i, col := range columns {
		maxW := 0
		for _, child := range col {
			outerW := child.Style.Width + child.Style.Margin.Left + child.Style.Margin.Right
			if outerW > maxW {
				maxW = outerW
			}
		}
		if maxW == 0 {
			maxW = 1 // Minimum width
		}
		colWidths[i] = maxW
	}

	// Position each column
	x := startX
	for i, col := range columns {
		colW := colWidths[i]

		// Calculate total height including margins and gap for this column
		totalHeight := 0
		for _, child := range col {
			totalHeight += child.Style.Height + child.Style.Margin.Top + child.Style.Margin.Bottom
		}
		totalGap := b.Style.Gap * (len(col) - 1)

		y, spacing := b.justifyColYForColumn(startY, innerH, totalHeight, totalGap, len(col))

		for _, child := range col {
			// Apply margins to position
			child.Y = b.Y + y + child.marginTop()
			childX, childW := b.alignColXForColumn(x, colW, child.Style.Width+child.Style.Margin.Left+child.Style.Margin.Right)
			child.X = childX + child.marginLeft()
			child.Style.Width = childW - child.Style.Margin.Left - child.Style.Margin.Right
			if child.Style.Width < 0 {
				child.Style.Width = 0
			}

			child.Calculate()

			// Advance by outer height plus spacing
			y += child.outerHeight() + b.childColSpacing(spacing)
		}

		x += colW + b.Style.Gap
	}
}

// splitIntoColumnLines divides children into columns that fit within innerH.
// Considers children's vertical margins when calculating column heights.
func splitIntoColumnLines(children []*Box, innerH, gap int) [][]*Box {
	var columns [][]*Box
	var currentCol []*Box
	currentHeight := 0

	for _, child := range children {
		// Use outer height (including margins)
		childH := child.Style.Height + child.Style.Margin.Top + child.Style.Margin.Bottom
		if childH == 0 {
			childH = 1 // Minimum height
		}

		gapNeeded := 0
		if len(currentCol) > 0 {
			gapNeeded = gap
		}

		// Check if child fits in current column
		if len(currentCol) > 0 && currentHeight+gapNeeded+childH > innerH {
			// Start new column
			columns = append(columns, currentCol)
			currentCol = []*Box{child}
			currentHeight = childH
		} else {
			// Add to current column
			currentCol = append(currentCol, child)
			currentHeight += gapNeeded + childH
		}
	}

	// Don't forget the last column
	if len(currentCol) > 0 {
		columns = append(columns, currentCol)
	}

	return columns
}

// justifyColYForColumn computes starting Y for a specific column.
func (b *Box) justifyColYForColumn(startY, innerH, totalHeight, totalGap, numChildren int) (y, spacing int) {
	if b.Style.JustifyContent == style.JustifyCenter {
		return startY + (innerH-totalHeight-totalGap)/2, 0
	}
	if b.Style.JustifyContent == style.JustifyEnd {
		return startY + innerH - totalHeight - totalGap, 0
	}
	if b.Style.JustifyContent == style.JustifyBetween && numChildren > 1 {
		return startY, (innerH - totalHeight) / (numChildren - 1)
	}
	if b.Style.JustifyContent == style.JustifyAround && numChildren > 0 {
		spacing := (innerH - totalHeight) / numChildren
		return startY + spacing/2, spacing
	}

	return startY, 0
}

// alignColXForColumn computes X position within a column.
func (b *Box) alignColXForColumn(colX, colW, childW int) (x, width int) {
	if b.Style.AlignItems == style.AlignCenter {
		return b.X + colX + (colW-childW)/2, childW
	}
	if b.Style.AlignItems == style.AlignEnd {
		return b.X + colX + colW - childW, childW
	}
	if b.Style.AlignItems == style.AlignStretch {
		return b.X + colX, colW
	}

	return b.X + colX, childW
}
