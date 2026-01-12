// Package layout computes positions for nested box elements.
// Supports block and flexbox layout modes.
package layout

import (
	"github.com/dottermi/x/termistyle/style"
	"github.com/kjk/flex"
)

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
// Uses the Yoga flex layout engine for accurate CSS flexbox calculations.
// If Width or Height is 0, auto-calculates based on children.
// Must be called after building the box tree.
//
// Example:
//
//	root := layout.NewBox(rootStyle)
//	root.AddChild(layout.NewBox(childStyle))
//	root.Calculate()
func (b *Box) Calculate() {
	// Handle text content auto-sizing before flex calculation
	b.prepareTextContent()

	// Build flex tree from box tree
	root := buildFlexTree(b, 0, true)

	// Determine parent constraints
	parentWidth := float32(b.Style.Width)
	parentHeight := float32(b.Style.Height)
	if parentWidth == 0 {
		parentWidth = flex.Undefined
	}
	if parentHeight == 0 {
		parentHeight = flex.Undefined
	}

	// Calculate layout using Yoga
	flex.CalculateLayout(root, parentWidth, parentHeight, flex.DirectionLTR)

	// Extract computed layout back to box tree
	extractLayout(root, b, 0, 0)

	// Position absolute children (handled separately from flex flow)
	// Calculate content area offsets from padding and border
	padding := b.Style.Padding
	borderLeft, borderTop := 0, 0
	if b.Style.Border.Left.IsSet() {
		borderLeft = 1
	}
	if b.Style.Border.Top.IsSet() {
		borderTop = 1
	}
	startX := padding.Left + borderLeft
	startY := padding.Top + borderTop
	b.positionAbsoluteChildren(startX, startY)
}

// prepareTextContent recursively prepares text content sizing.
// Text nodes need their content size set before flex calculation.
func (b *Box) prepareTextContent() {
	// Set content-based size for text nodes
	if b.Content != "" && b.Style.Width == 0 {
		b.Style.Width = len([]rune(b.Content))
	}
	if b.Content != "" && b.Style.Height == 0 {
		b.Style.Height = 1
	}

	// Recurse to children
	for _, child := range b.Children {
		child.prepareTextContent()
	}
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
