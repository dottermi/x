// Package termistyle provides CSS-like styling for terminal UIs in Go.
// Enables flexbox layouts, 24-bit colors, borders, and spacing.
// Build box trees, then render them to any io.Writer.
//
// Example:
//
//	container := ts.NewBox(ts.Style{
//		Width:      50,
//		Height:     12,
//		Display:    ts.Flex,
//		Direction:  ts.Row,
//		Background: ts.Color("#1a1a1a"),
//	})
//	tss.Print(os.Stdout, container)
package termistyle

import (
	"io"

	"github.com/dottermi/x/render"
	"github.com/dottermi/x/termistyle/draw"
	"github.com/dottermi/x/termistyle/layout"
	"github.com/dottermi/x/termistyle/style"
)

// Core types re-exported for convenience.
// See sub-packages for full documentation.
type (
	// Style holds visual and layout properties for a box element.
	Style = style.Style
	// Color represents a 24-bit RGB color as hex string ("#RRGGBB").
	Color = style.Color
	// Spacing defines space around an element's edges (padding/margin).
	Spacing = style.Spacing
	// Box is a rectangular layout element that can contain children.
	Box = layout.Box
	// Buffer is a 2D grid of cells for composing terminal output.
	Buffer = render.Buffer
	// Cell represents a single character position with colors.
	Cell = render.Cell
	// ClipRect defines a rectangular clipping region.
	ClipRect = render.ClipRect
	// Terminal manages double-buffered rendering with diff-based optimization.
	Terminal = render.Terminal
	// FontWeight represents text boldness (100-900, >= 600 is bold).
	FontWeight = style.FontWeight
	// FontStyle represents text style (normal, italic, oblique).
	FontStyle = style.FontStyle
	// TextDecoration represents text decorations (underline, line-through).
	TextDecoration = style.TextDecoration
	// TextTransform represents text case transformation.
	TextTransform = style.TextTransform
	// TextAlign represents horizontal text alignment.
	TextAlign = style.TextAlign
	// TextWrap represents text wrapping behavior.
	TextWrap = style.TextWrap
	// TextOverflow represents how overflowed text is handled.
	TextOverflow = style.TextOverflow
	// Overflow represents content overflow handling.
	Overflow = style.Overflow
	// Justify controls distribution of children along the main axis.
	Justify = style.Justify
	// Align controls positioning of children along the cross axis.
	Align = style.Align
	// FlexWrap controls whether flex items wrap onto multiple lines.
	FlexWrap = style.FlexWrap
	// Border specifies border style and color for each side.
	Border = style.Border
	// BorderEdge defines a border edge with style and color.
	BorderEdge = style.BorderEdge
	// BorderText represents a text element on a horizontal border edge.
	BorderText = style.BorderText
)

// Display mode constants control child layout behavior.
const (
	// Flex enables flexbox layout with configurable direction.
	Flex = style.Flex
	// Block stacks children vertically.
	Block = style.Block
)

// Position constants control element positioning.
const (
	// Relative positions element in normal flow (default).
	Relative = style.Relative
	// Absolute positions element relative to its container, removed from normal flow.
	Absolute = style.Absolute
)

// Flex direction constants set the main axis for flex layout.
const (
	// Row arranges children horizontally.
	Row = style.Row
	// Column arranges children vertically.
	Column = style.Column
)

// Constructor functions re-exported for convenience.
var (
	// SpacingAll creates uniform spacing on all sides.
	SpacingAll = style.SpacingAll
	// NewBox creates a box with the given style.
	NewBox = layout.NewBox
	// Text creates a text element with content and style.
	Text = layout.NewText
	// NewBuffer creates a buffer filled with space characters.
	NewBuffer = render.NewBuffer
	// NewTerminal creates a Terminal renderer for the specified dimensions.
	NewTerminal = render.NewTerminal
	// RGB creates a render.Color with the specified RGB values.
	RGB = render.RGB
)

// Border style constants select box-drawing character sets.
var (
	// BorderNone disables border rendering.
	BorderNone = style.BorderNone
	// BorderHidden renders invisible border (occupies space but uses spaces).
	BorderHidden = style.BorderHidden
	// BorderSingle uses thin box-drawing characters.
	BorderSingle = style.BorderSingle
	// BorderDouble uses double-line box-drawing characters.
	BorderDouble = style.BorderDouble
	// BorderRound uses single lines with rounded corners.
	BorderRound = style.BorderRound
	// BorderBold uses thick box-drawing characters.
	BorderBold = style.BorderBold
	// BorderAll creates a Border with the same style on all sides.
	BorderAll = style.BorderAll
	// BorderAllWithColor creates a Border with the same style and color on all sides.
	BorderAllWithColor = style.BorderAllWithColor
	// BorderXY creates a BorderSides with horizontal and vertical styles.
	BorderXY = style.BorderXY
	// BorderTRBL creates a BorderSides with top, right, bottom, left styles.
	BorderTRBL = style.BorderTRBL
	// BorderAllWithTitle creates a border with a title on the top edge.
	BorderAllWithTitle = style.BorderAllWithTitle
)

// JustifyContent constants control distribution along the main axis.
const (
	// JustifyStart aligns children to the start (flex-start).
	JustifyStart = style.JustifyStart
	// JustifyCenter centers children along the main axis.
	JustifyCenter = style.JustifyCenter
	// JustifyEnd aligns children to the end (flex-end).
	JustifyEnd = style.JustifyEnd
	// JustifyBetween distributes with equal space between children.
	JustifyBetween = style.JustifyBetween
	// JustifyAround distributes with equal space around children.
	JustifyAround = style.JustifyAround
	// SpaceBetween is a CSS-like alias for JustifyBetween.
	SpaceBetween = style.SpaceBetween
	// SpaceAround is a CSS-like alias for JustifyAround.
	SpaceAround = style.SpaceAround
)

// AlignItems constants control positioning along the cross axis.
const (
	// AlignStart positions children at cross-axis start (flex-start).
	AlignStart = style.AlignStart
	// AlignCenter centers children on the cross axis.
	AlignCenter = style.AlignCenter
	// AlignEnd positions children at cross-axis end (flex-end).
	AlignEnd = style.AlignEnd
	// AlignStretch expands children to fill the cross axis.
	AlignStretch = style.AlignStretch
	// Stretch is a CSS-like alias for AlignStretch.
	Stretch = style.Stretch
)

// FontWeight constants for text boldness.
const (
	WeightNormal = style.WeightNormal
	WeightBold   = style.WeightBold
)

// FontStyle constants for text style.
const (
	StyleNormal  = style.StyleNormal
	StyleItalic  = style.StyleItalic
	StyleOblique = style.StyleOblique
)

// TextDecoration constants for text decorations.
const (
	DecorationNone        = style.DecorationNone
	DecorationUnderline   = style.DecorationUnderline
	DecorationLineThrough = style.DecorationLineThrough
)

// TextTransform constants for text case transformation.
const (
	TransformNone       = style.TransformNone
	TransformUppercase  = style.TransformUppercase
	TransformLowercase  = style.TransformLowercase
	TransformCapitalize = style.TransformCapitalize
)

// TextAlign constants for horizontal text alignment.
const (
	TextAlignLeft   = style.TextAlignLeft
	TextAlignCenter = style.TextAlignCenter
	TextAlignRight  = style.TextAlignRight
)

// TextWrap constants for text wrapping behavior.
const (
	WrapNone = style.WrapNone
	WrapWord = style.WrapWord
	WrapChar = style.WrapChar
)

// TextOverflow constants for text overflow handling.
const (
	// TextOverflowClip clips text at the boundary (default).
	TextOverflowClip = style.TextOverflowClip
	// TextOverflowEllipsis truncates with "..." when text overflows.
	TextOverflowEllipsis = style.TextOverflowEllipsis
)

// Overflow constants for content overflow handling.
const (
	OverflowVisible = style.OverflowVisible
	OverflowHidden  = style.OverflowHidden
)

// FlexWrap constants for item wrapping behavior.
const (
	// NoWrap keeps all items on a single line (default).
	NoWrap = style.NoWrap
	// Wrap allows items to wrap onto multiple lines.
	Wrap = style.Wrap
)

// Draw computes layout and renders a box tree to a buffer.
// Calculates positions for all boxes, then draws backgrounds and borders.
// Returns a buffer ready for rendering to a terminal.
//
// Example:
//
//	box := termistyle.NewBox(termistyle.Style{Width: 40, Height: 10})
//	buf := termistyle.Draw(box)
//	term.Render(buf)
func Draw(box *layout.Box) *render.Buffer {
	box.Calculate()

	buf := render.NewBuffer(box.W, box.H)
	// Start with full buffer as clip bounds
	clip := render.ClipRect{X: 0, Y: 0, W: box.W, H: box.H}
	drawBoxClipped(buf, box, clip)

	return buf
}

// borderOffsets returns the offsets (0 or 1) for each border side.
func borderOffsets(border style.Border) (top, right, bottom, left int) {
	if border.Top.IsSet() {
		top = 1
	}
	if border.Right.IsSet() {
		right = 1
	}
	if border.Bottom.IsSet() {
		bottom = 1
	}
	if border.Left.IsSet() {
		left = 1
	}
	return top, right, bottom, left
}

// drawBackground fills the box interior with the background color.
func drawBackground(buf *render.Buffer, box *layout.Box, clip render.ClipRect) {
	if !box.Style.Background.IsSet() {
		return
	}
	top, right, bottom, left := borderOffsets(box.Style.Border)
	for y := top; y < box.H-bottom; y++ {
		for x := left; x < box.W-right; x++ {
			buf.SetClipped(box.X+x, box.Y+y, render.Cell{
				Char: ' ',
				BG:   box.Style.Background.ToRender(),
				FG:   box.Style.Foreground.ToRender(),
			}, clip)
		}
	}
}

// drawTextContent renders text content within the box.
func drawTextContent(buf *render.Buffer, box *layout.Box, clip render.ClipRect) {
	if box.Content == "" {
		return
	}
	top, right, bottom, left := borderOffsets(box.Style.Border)
	innerX := box.X + left
	innerY := box.Y + top
	innerWidth := box.W - left - right
	innerHeight := box.H - top - bottom
	draw.DrawStyledTextInBoxClipped(buf, innerX, innerY, box.Content, innerWidth, innerHeight, box.Style, clip)
}

// calculateChildClip computes the clip bounds for children.
func calculateChildClip(box *layout.Box, parentClip render.ClipRect) render.ClipRect {
	if box.Style.Overflow != style.OverflowHidden {
		return parentClip
	}
	top, right, bottom, left := borderOffsets(box.Style.Border)
	innerX := box.X + left
	innerY := box.Y + top
	innerW := box.W - left - right
	innerH := box.H - top - bottom
	return intersectClipRect(parentClip, render.ClipRect{X: innerX, Y: innerY, W: innerW, H: innerH})
}

func drawBoxClipped(buf *render.Buffer, box *layout.Box, clip render.ClipRect) {
	drawBackground(buf, box, clip)
	draw.DrawBorderClipped(buf, box.X, box.Y, box.W, box.H, box.Style.Border, clip)
	drawTextContent(buf, box, clip)

	childClip := calculateChildClip(box, clip)

	// Sort children by z-index (lower z-index rendered first, higher on top)
	children := sortByZIndex(box.Children)
	for _, child := range children {
		drawBoxClipped(buf, child, childClip)
	}
}

// sortByZIndex returns children sorted by z-index (ascending).
// Lower z-index elements are drawn first, higher ones on top.
func sortByZIndex(children []*layout.Box) []*layout.Box {
	if len(children) <= 1 {
		return children
	}

	// Check if any child has non-zero z-index
	hasZIndex := false
	for _, child := range children {
		if child.Style.ZIndex != 0 {
			hasZIndex = true
			break
		}
	}
	if !hasZIndex {
		return children
	}

	// Create a copy to avoid modifying original slice
	sorted := make([]*layout.Box, len(children))
	copy(sorted, children)

	// Simple insertion sort (stable, good for small n)
	for i := 1; i < len(sorted); i++ {
		j := i
		for j > 0 && sorted[j-1].Style.ZIndex > sorted[j].Style.ZIndex {
			sorted[j-1], sorted[j] = sorted[j], sorted[j-1]
			j--
		}
	}
	return sorted
}

// intersectClipRect returns the intersection of two clip rects.
func intersectClipRect(a, b render.ClipRect) render.ClipRect {
	x := max(a.X, b.X)
	y := max(a.Y, b.Y)
	x2 := min(a.X+a.W, b.X+b.W)
	y2 := min(a.Y+a.H, b.Y+b.H)

	w := x2 - x
	h := y2 - y
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}

	return render.ClipRect{X: x, Y: y, W: w, H: h}
}

// Print renders a box tree directly to a writer in one step.
// Combines Draw and RenderFullTo for common use cases.
// Outputs ANSI escape sequences for colors and positioning.
//
// Example:
//
//	container := ts.NewBox(ts.Style{
//		Width:      40,
//		Height:     10,
//		Background: ts.Color("#1a1a2e"),
//	})
//	ts.Print(os.Stdout, container)
func Print(w io.Writer, box *layout.Box) {
	buf := Draw(box)
	term := render.NewTerminal(buf.Width, buf.Height)
	term.RenderFullTo(buf, w)
}

// Println renders a box tree to a writer and appends a newline.
// Combines Draw and RenderFullTo, adding a newline after rendering.
//
// Example:
//
//	container := ts.NewBox(ts.Style{
//		Width:      40,
//		Height:     10,
//		Background: ts.Color("#1a1a2e"),
//	})
//	ts.Println(os.Stdout, container)
func Println(w io.Writer, box *layout.Box) {
	buf := Draw(box)
	term := render.NewTerminal(buf.Width, buf.Height)
	term.RenderFullTo(buf, w)
	_, _ = w.Write([]byte("\n"))
}

// Render converts a buffer to an ANSI-escaped string using full rendering.
// This is a convenience function for simple cases.
func Render(buf *render.Buffer) string {
	term := render.NewTerminal(buf.Width, buf.Height)
	return term.RenderFull(buf)
}

// RenderTo writes a buffer to an io.Writer with ANSI escape codes.
// This is a convenience function for simple cases.
func RenderTo(buf *render.Buffer, w io.Writer) {
	term := render.NewTerminal(buf.Width, buf.Height)
	term.RenderFullTo(buf, w)
}
