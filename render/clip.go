package render

// ClipRect defines a rectangular clipping region.
// Coordinates use (0,0) as top-left corner.
type ClipRect struct {
	X, Y, W, H int
}

// Contains returns true if the point (x, y) is within the clip rect bounds.
func (c ClipRect) Contains(x, y int) bool {
	return x >= c.X && x < c.X+c.W && y >= c.Y && y < c.Y+c.H
}

// SetClipped writes a cell at the given position only if within clip bounds.
// Both buffer bounds and clip bounds are checked.
func (b *Buffer) SetClipped(x, y int, c Cell, clip ClipRect) {
	if clip.Contains(x, y) {
		b.Set(x, y, c)
	}
}
