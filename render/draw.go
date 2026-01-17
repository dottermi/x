package render

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
// Starts at (x, y) and extends rightward for length cells.
func (b *Buffer) FillHorizontal(x, y, length int, c Cell) {
	for i := range length {
		b.Set(x+i, y, c)
	}
}

// FillVertical draws a vertical line of cells.
// Starts at (x, y) and extends downward for length cells.
func (b *Buffer) FillVertical(x, y, length int, c Cell) {
	for i := range length {
		b.Set(x, y+i, c)
	}
}
