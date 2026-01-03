package style

// Spacing defines space around an element's edges.
// Used for padding and margin configuration.
type Spacing struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// SpacingAll creates uniform spacing on all sides.
//
// Example:
//
//	padding := style.SpacingAll(2)  // 2 cells on every side
func SpacingAll(v int) Spacing {
	return Spacing{Top: v, Right: v, Bottom: v, Left: v}
}

// SpacingXY creates spacing with separate horizontal and vertical values.
//
// Parameters:
//   - x: horizontal spacing (left and right)
//   - y: vertical spacing (top and bottom)
//
// Example:
//
//	padding := style.SpacingXY(4, 1)  // 4 horizontal, 1 vertical
func SpacingXY(x, y int) Spacing {
	return Spacing{Top: y, Right: x, Bottom: y, Left: x}
}
