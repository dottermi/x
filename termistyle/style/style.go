// Package style defines visual properties for terminal UI elements.
// Provides CSS-like styling with flexbox layout support.
package style

// Style holds all visual and layout properties for a UI element.
// Combines dimensions, spacing, colors, and flex layout configuration.
//
// Example:
//
//	s := style.Style{
//		Width:          40,
//		Height:         10,
//		Display:        style.Flex,
//		FlexDirection:  style.Row,
//		Background:     style.Color("#1a1a2e"),
//	}
type Style struct {
	Display  Display
	Position PositionType
	ZIndex   int

	Width     int
	Height    int
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int

	Padding Spacing
	Margin  Spacing

	Border Border

	Background Color
	Foreground Color

	FontWeight     FontWeight
	FontStyle      FontStyle
	TextDecoration TextDecoration
	TextTransform  TextTransform
	TextAlign      TextAlign
	TextWrap       TextWrap
	TextOverflow   TextOverflow
	Dim            bool
	Reverse        bool

	FlexDirection  FlexDirection
	FlexWrap       FlexWrap
	JustifyContent Justify
	AlignItems     Align
	AlignSelf      Align // Per-item override of AlignItems
	Gap            int
	FlexGrow       float64
	FlexShrink     float64
	FlexBasis      int
	AspectRatio    float64 // Width/Height ratio constraint
	Overflow       Overflow
}
