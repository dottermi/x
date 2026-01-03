// Package layout computes positions for nested box elements.
// Supports block and flexbox layout modes.
package layout

import "github.com/dottermi/x/termistyle/style"

// Box represents a rectangular layout element with children.
// After Calculate is called, X, Y, W, H contain computed positions.
type Box struct {
	Style    style.Style
	Children []*Box
	Content  string // Text content for text elements

	// Computed values set by Calculate
	X, Y int
	W, H int
}

// NewBox creates a box with the given style.
//
// Example:
//
//	box := layout.NewBox(style.Style{Width: 40, Height: 10})
func NewBox(s style.Style) *Box {
	return &Box{Style: s}
}

// NewText creates a text element with the given content and style.
// Width and height are auto-calculated from content if not specified.
//
// Example:
//
//	text := layout.NewText("Hello", style.Style{Foreground: style.Color("#FFF")})
func NewText(content string, s style.Style) *Box {
	// Auto-calculate dimensions if not specified
	if s.Width == 0 {
		s.Width = len([]rune(content))
	}
	if s.Height == 0 {
		s.Height = 1
	}
	return &Box{Style: s, Content: content}
}

// AddChild appends a child box to this box's children.
func (b *Box) AddChild(child *Box) {
	b.Children = append(b.Children, child)
}

// Calculate computes positions for this box and all descendants.
// Uses flex or block layout based on the Display property.
// If Width or Height is 0, auto-calculates based on children.
// Must be called after building the box tree.
//
// Example:
//
//	root := layout.NewBox(rootStyle)
//	root.AddChild(layout.NewBox(childStyle))
//	root.Calculate()
func (b *Box) Calculate() {
	// Calculate padding and border sizes
	padding := b.Style.Padding
	borderTop, borderRight, borderBottom, borderLeft := 0, 0, 0, 0
	if b.Style.Border.Top.IsSet() {
		borderTop = 1
	}
	if b.Style.Border.Right.IsSet() {
		borderRight = 1
	}
	if b.Style.Border.Bottom.IsSet() {
		borderBottom = 1
	}
	if b.Style.Border.Left.IsSet() {
		borderLeft = 1
	}

	extraW := padding.Left + padding.Right + borderLeft + borderRight
	extraH := padding.Top + padding.Bottom + borderTop + borderBottom

	// Handle auto-sizing when Width or Height is 0
	b.W = b.Style.Width
	b.H = b.Style.Height
	if b.Style.Width == 0 || b.Style.Height == 0 {
		autoW, autoH := b.calculateAutoSize()
		if b.Style.Width == 0 {
			b.W = autoW + extraW
		}
		if b.Style.Height == 0 {
			b.H = autoH + extraH
		}
	}

	if len(b.Children) == 0 {
		return
	}

	startX := padding.Left + borderLeft
	startY := padding.Top + borderTop
	innerW := b.W - extraW
	innerH := b.H - extraH

	if b.Style.Display == style.Flex {
		b.calculateFlex(startX, startY, innerW, innerH)
	} else {
		b.calculateBlock(startX, startY)
	}
}

// calculateAutoSize calculates the minimum size needed to fit all children.
// Returns the inner width and height (excluding padding and border).
// Considers children's margins when calculating total size.
func (b *Box) calculateAutoSize() (int, int) {
	if len(b.Children) == 0 {
		// For text content, use content size
		if b.Content != "" {
			return len([]rune(b.Content)), 1
		}
		return 0, 0
	}

	// First, recursively calculate sizes for all children
	for _, child := range b.Children {
		child.Calculate()
	}

	gap := b.Style.Gap
	isRow := b.Style.Display == style.Flex && b.Style.FlexDirection == style.Row

	var totalW, totalH, maxW, maxH int
	for i, child := range b.Children {
		// Use outer dimensions (including margins) for layout calculations
		outerW := child.outerWidth()
		outerH := child.outerHeight()

		if outerW > maxW {
			maxW = outerW
		}
		if outerH > maxH {
			maxH = outerH
		}
		totalW += outerW
		totalH += outerH
		if i > 0 {
			totalW += gap
			totalH += gap
		}
	}

	if isRow {
		// Row: sum widths, max height
		return totalW, maxH
	}
	// Column or Block: max width, sum heights
	return maxW, totalH
}

// calculateBlock positions children in vertical stack order.
// Respects each child's margin for positioning.
func (b *Box) calculateBlock(startX, startY int) {
	y := startY
	for _, child := range b.Children {
		// Skip absolute positioned elements in normal flow
		if child.Style.Position == style.Absolute {
			continue
		}
		// Apply margins to position
		child.X = b.X + startX + child.marginLeft()
		child.Y = b.Y + y + child.marginTop()
		child.Calculate()
		// Advance by outer height (including margins) plus gap
		y += child.outerHeight() + b.Style.Gap
	}
	// Position absolute children relative to this container
	b.positionAbsoluteChildren(startX, startY)
}

// calculateFlex dispatches to row or column flex calculation.
func (b *Box) calculateFlex(startX, startY, innerW, innerH int) {
	if b.Style.FlexWrap == style.Wrap {
		if b.Style.FlexDirection == style.Column {
			b.calculateFlexColumnWrap(startX, startY, innerW, innerH)
		} else {
			b.calculateFlexRowWrap(startX, startY, innerW, innerH)
		}
		b.positionAbsoluteChildren(startX, startY)
		return
	}

	if b.Style.FlexDirection == style.Column {
		b.calculateFlexColumn(startX, startY, innerW, innerH)
	} else {
		b.calculateFlexRow(startX, startY, innerW, innerH)
	}
	b.positionAbsoluteChildren(startX, startY)
}

// outerWidth returns the total width including margins.
func (b *Box) outerWidth() int {
	return b.W + b.Style.Margin.Left + b.Style.Margin.Right
}

// outerHeight returns the total height including margins.
func (b *Box) outerHeight() int {
	return b.H + b.Style.Margin.Top + b.Style.Margin.Bottom
}

// marginLeft returns the left margin value.
func (b *Box) marginLeft() int {
	return b.Style.Margin.Left
}

// marginTop returns the top margin value.
func (b *Box) marginTop() int {
	return b.Style.Margin.Top
}

// positionAbsoluteChildren positions children with Position: Absolute.
// Absolute children are positioned relative to the parent container,
// using their X, Y fields as offsets from the container's content area.
func (b *Box) positionAbsoluteChildren(startX, startY int) {
	for _, child := range b.Children {
		if child.Style.Position != style.Absolute {
			continue
		}
		// Use child's X, Y as offset from parent's content area
		child.X = b.X + startX + child.X
		child.Y = b.Y + startY + child.Y
		// Set dimensions
		child.W = child.Style.Width
		child.H = child.Style.Height
		// Recursively calculate children of absolute element
		child.Calculate()
	}
}
