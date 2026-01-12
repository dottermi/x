package render

// Buffer represents a 2D grid of [Cell] values for terminal rendering.
// Cells are stored in row-major order and accessed via (x, y) coordinates.
// Use [NewBuffer] to create a buffer pre-filled with empty cells.
//
// Example:
//
//	buf := render.NewBuffer(80, 24)
//	buf.Set(0, 0, render.Cell{Char: 'H', FG: render.RGB(0, 255, 0)})
//	buf.Set(1, 0, render.Cell{Char: 'i', FG: render.RGB(0, 255, 0)})
type Buffer struct {
	Width  int
	Height int
	Cells  []Cell // row-major: index = y*Width + x
}

// NewBuffer creates a Buffer of the specified dimensions filled with [EmptyCell] values.
//
// Parameters:
//   - width: number of columns
//   - height: number of rows
//
// Example:
//
//	buf := render.NewBuffer(80, 24)
func NewBuffer(width, height int) *Buffer {
	cells := make([]Cell, width*height)
	for i := range cells {
		cells[i] = EmptyCell()
	}
	return &Buffer{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// index converts (x, y) to flat array index.
func (b *Buffer) index(x, y int) int {
	return y*b.Width + x
}

// inBounds checks if (x, y) is within buffer bounds.
func (b *Buffer) inBounds(x, y int) bool {
	return x >= 0 && x < b.Width && y >= 0 && y < b.Height
}

// Set writes a cell at the specified position.
// Out-of-bounds coordinates are silently ignored, enabling safe clipping.
//
// Parameters:
//   - x: column index (0-based)
//   - y: row index (0-based)
//   - c: cell to write
func (b *Buffer) Set(x, y int, c Cell) {
	if b.inBounds(x, y) {
		b.Cells[b.index(x, y)] = c
	}
}

// Get returns the cell at the specified position.
// Returns [EmptyCell] for out-of-bounds coordinates.
//
// Parameters:
//   - x: column index (0-based)
//   - y: row index (0-based)
func (b *Buffer) Get(x, y int) Cell {
	if b.inBounds(x, y) {
		return b.Cells[b.index(x, y)]
	}
	return EmptyCell()
}

// Fill sets every cell in the buffer to the specified value.
// Useful for clearing the buffer or setting a uniform background.
//
// Example:
//
//	bg := render.Cell{Char: ' ', BG: render.RGB(0, 0, 64)}
//	buf.Fill(bg)
func (b *Buffer) Fill(c Cell) {
	for i := range b.Cells {
		b.Cells[i] = c
	}
}

// Clone creates a deep copy of the buffer.
// The returned buffer shares no memory with the original.
func (b *Buffer) Clone() *Buffer {
	cells := make([]Cell, len(b.Cells))
	copy(cells, b.Cells)
	return &Buffer{
		Width:  b.Width,
		Height: b.Height,
		Cells:  cells,
	}
}

// CellChange represents a modified cell at a specific screen position.
// Used by [Buffer.Diff] to report which cells need updating.
type CellChange struct {
	X, Y int
	Cell Cell
}

// Diff compares this buffer with new and returns cells that differ.
// Only the intersection of both buffer dimensions is compared.
// Returns nil if buffers are identical.
//
// This is the core optimization: instead of redrawing the entire screen,
// only changed cells need ANSI output.
//
// Example:
//
//	changes := oldBuf.Diff(newBuf)
//	for _, c := range changes {
//		fmt.Printf("Cell at (%d,%d) changed to %q\n", c.X, c.Y, c.Cell.Char)
//	}
func (b *Buffer) Diff(new *Buffer) []CellChange {
	var changes []CellChange

	// Handle mismatched dimensions
	maxY := b.Height
	if new.Height < maxY {
		maxY = new.Height
	}
	maxX := b.Width
	if new.Width < maxX {
		maxX = new.Width
	}

	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			oldCell := b.Get(x, y)
			newCell := new.Get(x, y)
			if !oldCell.Equal(newCell) {
				changes = append(changes, CellChange{X: x, Y: y, Cell: newCell})
			}
		}
	}

	return changes
}
