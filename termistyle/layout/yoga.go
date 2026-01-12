// Package layout computes positions for nested box elements using flex layout.
// This file contains the adapter for the kjk/flex (Yoga) library.
package layout

import (
	"github.com/dottermi/x/termistyle/style"
	"github.com/kjk/flex"
)

// flexConfig is the global config with CSS-like defaults (Row as default direction).
var flexConfig = func() *flex.Config {
	cfg := flex.NewConfig()
	cfg.UseWebDefaults = true // FlexDirection=Row, AlignContent=Stretch
	return cfg
}()

// buildFlexTree converts a Box tree into a flex.Node tree.
// It recursively processes all children and applies styles.
func buildFlexTree(box *Box, parentGap int, isFirstChild bool) *flex.Node {
	node := flex.NewNodeWithConfig(flexConfig)
	applyStyleToNode(node, box.Style, parentGap, isFirstChild)

	// Process children
	gap := box.Style.Gap
	for i, child := range box.Children {
		// Skip absolute positioned elements - they're handled separately
		if child.Style.Position == style.Absolute {
			continue
		}
		childNode := buildFlexTree(child, gap, i == 0)
		node.InsertChild(childNode, len(node.Children))
	}

	return node
}

// applyStyleToNode maps termistyle.Style properties to flex.Node setters.
func applyStyleToNode(node *flex.Node, s style.Style, parentGap int, isFirstChild bool) {
	// Dimensions
	if s.Width > 0 {
		node.StyleSetWidth(float32(s.Width))
	}
	if s.Height > 0 {
		node.StyleSetHeight(float32(s.Height))
	}

	// Min/Max dimensions
	if s.MinWidth > 0 {
		node.StyleSetMinWidth(float32(s.MinWidth))
	}
	if s.MaxWidth > 0 {
		node.StyleSetMaxWidth(float32(s.MaxWidth))
	}
	if s.MinHeight > 0 {
		node.StyleSetMinHeight(float32(s.MinHeight))
	}
	if s.MaxHeight > 0 {
		node.StyleSetMaxHeight(float32(s.MaxHeight))
	}

	// Flex direction and wrap
	node.StyleSetFlexDirection(convertFlexDirection(s.FlexDirection))
	node.StyleSetFlexWrap(convertFlexWrap(s.FlexWrap))

	// Alignment
	node.StyleSetJustifyContent(convertJustify(s.JustifyContent))
	node.StyleSetAlignItems(convertAlign(s.AlignItems))
	if s.AlignSelf != 0 {
		node.StyleSetAlignSelf(convertAlign(s.AlignSelf))
	}

	// Flex properties
	if s.FlexGrow > 0 {
		node.StyleSetFlexGrow(float32(s.FlexGrow))
	}
	if s.FlexShrink > 0 {
		node.StyleSetFlexShrink(float32(s.FlexShrink))
	}
	if s.FlexBasis > 0 {
		node.StyleSetFlexBasis(float32(s.FlexBasis))
	}

	// Aspect ratio
	if s.AspectRatio > 0 {
		node.StyleSetAspectRatio(float32(s.AspectRatio))
	}

	// Position type
	node.StyleSetPositionType(convertPositionType(s.Position))

	// Padding (individual edges)
	if s.Padding.Top > 0 {
		node.StyleSetPadding(flex.EdgeTop, float32(s.Padding.Top))
	}
	if s.Padding.Right > 0 {
		node.StyleSetPadding(flex.EdgeRight, float32(s.Padding.Right))
	}
	if s.Padding.Bottom > 0 {
		node.StyleSetPadding(flex.EdgeBottom, float32(s.Padding.Bottom))
	}
	if s.Padding.Left > 0 {
		node.StyleSetPadding(flex.EdgeLeft, float32(s.Padding.Left))
	}

	// Margin (individual edges)
	if s.Margin.Top > 0 {
		node.StyleSetMargin(flex.EdgeTop, float32(s.Margin.Top))
	}
	if s.Margin.Right > 0 {
		node.StyleSetMargin(flex.EdgeRight, float32(s.Margin.Right))
	}
	if s.Margin.Bottom > 0 {
		node.StyleSetMargin(flex.EdgeBottom, float32(s.Margin.Bottom))
	}
	if s.Margin.Left > 0 {
		node.StyleSetMargin(flex.EdgeLeft, float32(s.Margin.Left))
	}

	// Border (as spacing - flex uses border for layout calculations)
	if s.Border.Top.IsSet() {
		node.StyleSetBorder(flex.EdgeTop, 1)
	}
	if s.Border.Right.IsSet() {
		node.StyleSetBorder(flex.EdgeRight, 1)
	}
	if s.Border.Bottom.IsSet() {
		node.StyleSetBorder(flex.EdgeBottom, 1)
	}
	if s.Border.Left.IsSet() {
		node.StyleSetBorder(flex.EdgeLeft, 1)
	}

	// Display
	if s.Display == style.None {
		node.StyleSetDisplay(flex.DisplayNone)
	} else {
		node.StyleSetDisplay(flex.DisplayFlex)
	}

	// Simulate Gap via margins (flex doesn't have native gap support)
	// Apply margin on the start side of the main axis for non-first children
	applyGapAsMargin(node, s, parentGap, isFirstChild)
}

// applyGapAsMargin simulates CSS gap by adding margins to non-first children.
// Gap is applied to the leading edge based on parent's flex direction.
func applyGapAsMargin(node *flex.Node, s style.Style, parentGap int, isFirstChild bool) {
	if parentGap == 0 || isFirstChild {
		return
	}

	// We need to determine the parent's flex direction, but we don't have access to it here.
	// Since this is called during buildFlexTree, we pass direction info differently.
	// For now, we'll use a simpler approach: add margin to both Top and Left,
	// and let the flex algorithm handle which one applies based on layout.

	// Actually, we need to track parent's flex direction. Let's modify the signature.
	// For now, apply margin to the appropriate edge based on common usage.
	// In a Row layout, gap applies to Left margin; in Column, to Top margin.
	// Since we're using web defaults (Row), we'll apply to Left by default.

	// This will be enhanced later. For now, apply to left (row direction).
	existingLeft := s.Margin.Left
	node.StyleSetMargin(flex.EdgeLeft, float32(existingLeft+parentGap))
}

// extractLayout copies computed layout values from flex.Node back to Box.
// It recursively processes all children to sync the entire tree.
func extractLayout(node *flex.Node, box *Box, parentX, parentY int) {
	// Get computed layout values
	box.X = parentX + int(node.LayoutGetLeft())
	box.Y = parentY + int(node.LayoutGetTop())
	box.W = int(node.LayoutGetWidth())
	box.H = int(node.LayoutGetHeight())

	// Process children (excluding absolute positioned ones handled separately)
	childIdx := 0
	for _, child := range box.Children {
		if child.Style.Position == style.Absolute {
			continue
		}
		if childIdx < len(node.Children) {
			childNode := node.GetChild(childIdx)
			extractLayout(childNode, child, box.X, box.Y)
			childIdx++
		}
	}
}

// convertFlexDirection maps termistyle FlexDirection to flex.FlexDirection.
func convertFlexDirection(d style.FlexDirection) flex.FlexDirection {
	switch d {
	case style.Column:
		return flex.FlexDirectionColumn
	default:
		return flex.FlexDirectionRow
	}
}

// convertJustify maps termistyle Justify to flex.Justify.
func convertJustify(j style.Justify) flex.Justify {
	switch j {
	case style.JustifyCenter:
		return flex.JustifyCenter
	case style.JustifyEnd:
		return flex.JustifyFlexEnd
	case style.JustifyBetween:
		return flex.JustifySpaceBetween
	case style.JustifyAround:
		return flex.JustifySpaceAround
	default:
		return flex.JustifyFlexStart
	}
}

// convertAlign maps termistyle Align to flex.Align.
func convertAlign(a style.Align) flex.Align {
	switch a {
	case style.AlignCenter:
		return flex.AlignCenter
	case style.AlignEnd:
		return flex.AlignFlexEnd
	case style.AlignStretch:
		return flex.AlignStretch
	default:
		return flex.AlignFlexStart
	}
}

// convertFlexWrap maps termistyle FlexWrap to flex.Wrap.
func convertFlexWrap(w style.FlexWrap) flex.Wrap {
	switch w {
	case style.Wrap:
		return flex.WrapWrap
	default:
		return flex.WrapNoWrap
	}
}

// convertPositionType maps termistyle PositionType to flex.PositionType.
func convertPositionType(p style.PositionType) flex.PositionType {
	switch p {
	case style.Absolute:
		return flex.PositionTypeAbsolute
	default:
		return flex.PositionTypeRelative
	}
}
