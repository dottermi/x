package draw

// Buffer is a 2D grid of cells for composing terminal output.
// Coordinates use (0,0) as top-left corner.
type Buffer struct {
	Width  int
	Height int
	Cells  [][]Cell
}

// NewBuffer creates a buffer filled with space characters.
// All cells start with default (zero) colors.
//
// Example:
//
//	buf := draw.NewBuffer(80, 24)
func NewBuffer(width, height int) *Buffer {
	cells := make([][]Cell, height)
	for y := range cells {
		cells[y] = make([]Cell, width)
		for x := range cells[y] {
			cells[y][x] = Cell{Char: ' '}
		}
	}

	return &Buffer{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// Set writes a cell at the given position.
// Out-of-bounds coordinates are silently ignored.
func (b *Buffer) Set(x, y int, c Cell) {
	if x >= 0 && x < b.Width && y >= 0 && y < b.Height {
		b.Cells[y][x] = c
	}
}

// Get returns the cell at the given position.
// Returns an empty Cell for out-of-bounds coordinates.
func (b *Buffer) Get(x, y int) Cell {
	if x >= 0 && x < b.Width && y >= 0 && y < b.Height {
		return b.Cells[y][x]
	}
	return Cell{}
}

// ClipRect defines a rectangular clipping region.
type ClipRect struct {
	X, Y, W, H int
}

// Contains returns true if the point is within the clip rect.
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

// Fill sets every cell in the buffer to the given cell.
func (b *Buffer) Fill(c Cell) {
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			b.Cells[y][x] = c
		}
	}
}

// FillRect fills a rectangular region with the given cell.
// Cells outside buffer bounds are silently skipped.
func (b *Buffer) FillRect(x, y, width, height int, c Cell) {
	for dy := range height {
		for dx := range width {
			b.Set(x+dx, y+dy, c)
		}
	}
}

// FillHorizontal draws a horizontal line of cells.
// Starts at (x, y) and extends rightward.
func (b *Buffer) FillHorizontal(x, y, length int, c Cell) {
	for i := range length {
		b.Set(x+i, y, c)
	}
}

// FillVertical draws a vertical line of cells.
// Starts at (x, y) and extends downward.
func (b *Buffer) FillVertical(x, y, length int, c Cell) {
	for i := range length {
		b.Set(x, y+i, c)
	}
}
