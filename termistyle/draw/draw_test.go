package draw

import (
	"testing"

	"github.com/dottermi/x/termistyle/style"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBuffer(t *testing.T) {
	t.Parallel()

	t.Run("should create buffer with specified dimensions", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(80, 24)

		assert.Equal(t, 80, buf.Width)
		assert.Equal(t, 24, buf.Height)
		assert.Len(t, buf.Cells, 24)
		assert.Len(t, buf.Cells[0], 80)
	})

	t.Run("should fill all cells with space character", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(5, 3)

		for y := 0; y < buf.Height; y++ {
			for x := 0; x < buf.Width; x++ {
				assert.Equal(t, ' ', buf.Cells[y][x].Char)
			}
		}
	})

	t.Run("should create buffer with zero dimensions", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(0, 0)

		assert.Equal(t, 0, buf.Width)
		assert.Equal(t, 0, buf.Height)
		assert.Empty(t, buf.Cells)
	})

	t.Run("should create buffer with only width", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 0)

		assert.Equal(t, 10, buf.Width)
		assert.Equal(t, 0, buf.Height)
		assert.Empty(t, buf.Cells)
	})

	t.Run("should create buffer with only height", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(0, 5)

		assert.Equal(t, 0, buf.Width)
		assert.Equal(t, 5, buf.Height)
		assert.Len(t, buf.Cells, 5)
		for y := 0; y < buf.Height; y++ {
			assert.Empty(t, buf.Cells[y])
		}
	})

	t.Run("should initialize cells with default colors", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(2, 2)

		cell := buf.Cells[0][0]
		assert.Equal(t, style.Color(""), cell.Foreground)
		assert.Equal(t, style.Color(""), cell.Background)
	})
}

func TestBuffer_Set(t *testing.T) {
	t.Parallel()

	t.Run("should write cell at valid position", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: 'X', Foreground: style.Color("#FF0000")}

		buf.Set(5, 5, cell)

		assert.Equal(t, 'X', buf.Cells[5][5].Char)
		assert.Equal(t, style.Color("#FF0000"), buf.Cells[5][5].Foreground)
	})

	t.Run("should write cell at top-left corner", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: 'A'}

		buf.Set(0, 0, cell)

		assert.Equal(t, 'A', buf.Cells[0][0].Char)
	})

	t.Run("should write cell at bottom-right corner", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: 'Z'}

		buf.Set(9, 9, cell)

		assert.Equal(t, 'Z', buf.Cells[9][9].Char)
	})

	t.Run("should ignore negative x coordinate", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: 'X'}

		buf.Set(-1, 5, cell)

		for y := 0; y < buf.Height; y++ {
			for x := 0; x < buf.Width; x++ {
				assert.Equal(t, ' ', buf.Cells[y][x].Char)
			}
		}
	})

	t.Run("should ignore negative y coordinate", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: 'X'}

		buf.Set(5, -1, cell)

		for y := 0; y < buf.Height; y++ {
			for x := 0; x < buf.Width; x++ {
				assert.Equal(t, ' ', buf.Cells[y][x].Char)
			}
		}
	})

	t.Run("should ignore x coordinate beyond width", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: 'X'}

		buf.Set(10, 5, cell)

		for y := 0; y < buf.Height; y++ {
			for x := 0; x < buf.Width; x++ {
				assert.Equal(t, ' ', buf.Cells[y][x].Char)
			}
		}
	})

	t.Run("should ignore y coordinate beyond height", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: 'X'}

		buf.Set(5, 10, cell)

		for y := 0; y < buf.Height; y++ {
			for x := 0; x < buf.Width; x++ {
				assert.Equal(t, ' ', buf.Cells[y][x].Char)
			}
		}
	})

	t.Run("should overwrite existing cell", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		buf.Set(5, 5, Cell{Char: 'A'})

		buf.Set(5, 5, Cell{Char: 'B'})

		assert.Equal(t, 'B', buf.Cells[5][5].Char)
	})
}

func TestBuffer_Get(t *testing.T) {
	t.Parallel()

	t.Run("should return cell at valid position", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		buf.Cells[3][7] = Cell{Char: 'Q', Foreground: style.Color("#00FF00")}

		cell := buf.Get(7, 3)

		assert.Equal(t, 'Q', cell.Char)
		assert.Equal(t, style.Color("#00FF00"), cell.Foreground)
	})

	t.Run("should return empty cell for negative x", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)

		cell := buf.Get(-1, 5)

		assert.Equal(t, Cell{}, cell)
	})

	t.Run("should return empty cell for negative y", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)

		cell := buf.Get(5, -1)

		assert.Equal(t, Cell{}, cell)
	})

	t.Run("should return empty cell for x beyond width", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)

		cell := buf.Get(10, 5)

		assert.Equal(t, Cell{}, cell)
	})

	t.Run("should return empty cell for y beyond height", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)

		cell := buf.Get(5, 10)

		assert.Equal(t, Cell{}, cell)
	})

	t.Run("should return cell at top-left corner", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		buf.Cells[0][0] = Cell{Char: 'T'}

		cell := buf.Get(0, 0)

		assert.Equal(t, 'T', cell.Char)
	})

	t.Run("should return cell at bottom-right corner", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		buf.Cells[9][9] = Cell{Char: 'B'}

		cell := buf.Get(9, 9)

		assert.Equal(t, 'B', cell.Char)
	})
}

func TestClipRect_Contains(t *testing.T) {
	t.Parallel()

	t.Run("should return true for point inside rect", func(t *testing.T) {
		t.Parallel()

		clip := ClipRect{X: 10, Y: 10, W: 20, H: 20}

		assert.True(t, clip.Contains(15, 15))
	})

	t.Run("should return true for point at top-left corner", func(t *testing.T) {
		t.Parallel()

		clip := ClipRect{X: 10, Y: 10, W: 20, H: 20}

		assert.True(t, clip.Contains(10, 10))
	})

	t.Run("should return false for point at bottom-right corner (exclusive)", func(t *testing.T) {
		t.Parallel()

		clip := ClipRect{X: 10, Y: 10, W: 20, H: 20}

		assert.False(t, clip.Contains(30, 30))
	})

	t.Run("should return true for point just inside bottom-right", func(t *testing.T) {
		t.Parallel()

		clip := ClipRect{X: 10, Y: 10, W: 20, H: 20}

		assert.True(t, clip.Contains(29, 29))
	})

	t.Run("should return false for point left of rect", func(t *testing.T) {
		t.Parallel()

		clip := ClipRect{X: 10, Y: 10, W: 20, H: 20}

		assert.False(t, clip.Contains(9, 15))
	})

	t.Run("should return false for point above rect", func(t *testing.T) {
		t.Parallel()

		clip := ClipRect{X: 10, Y: 10, W: 20, H: 20}

		assert.False(t, clip.Contains(15, 9))
	})

	t.Run("should return false for point right of rect", func(t *testing.T) {
		t.Parallel()

		clip := ClipRect{X: 10, Y: 10, W: 20, H: 20}

		assert.False(t, clip.Contains(30, 15))
	})

	t.Run("should return false for point below rect", func(t *testing.T) {
		t.Parallel()

		clip := ClipRect{X: 10, Y: 10, W: 20, H: 20}

		assert.False(t, clip.Contains(15, 30))
	})

	t.Run("should handle zero-size rect", func(t *testing.T) {
		t.Parallel()

		clip := ClipRect{X: 10, Y: 10, W: 0, H: 0}

		assert.False(t, clip.Contains(10, 10))
	})

	t.Run("should handle rect at origin", func(t *testing.T) {
		t.Parallel()

		clip := ClipRect{X: 0, Y: 0, W: 10, H: 10}

		assert.True(t, clip.Contains(0, 0))
		assert.True(t, clip.Contains(5, 5))
		assert.False(t, clip.Contains(10, 10))
	})
}

func TestBuffer_SetClipped(t *testing.T) {
	t.Parallel()

	t.Run("should write cell when within clip bounds", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 20)
		clip := ClipRect{X: 5, Y: 5, W: 10, H: 10}
		cell := Cell{Char: 'X'}

		buf.SetClipped(10, 10, cell, clip)

		assert.Equal(t, 'X', buf.Cells[10][10].Char)
	})

	t.Run("should not write cell when outside clip bounds", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 20)
		clip := ClipRect{X: 5, Y: 5, W: 10, H: 10}
		cell := Cell{Char: 'X'}

		buf.SetClipped(2, 2, cell, clip)

		assert.Equal(t, ' ', buf.Cells[2][2].Char)
	})

	t.Run("should respect both buffer and clip bounds", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		clip := ClipRect{X: 0, Y: 0, W: 20, H: 20}
		cell := Cell{Char: 'X'}

		buf.SetClipped(15, 15, cell, clip)

		for y := 0; y < buf.Height; y++ {
			for x := 0; x < buf.Width; x++ {
				assert.Equal(t, ' ', buf.Cells[y][x].Char)
			}
		}
	})

	t.Run("should write at clip edge", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 20)
		clip := ClipRect{X: 5, Y: 5, W: 10, H: 10}
		cell := Cell{Char: 'E'}

		buf.SetClipped(5, 5, cell, clip)

		assert.Equal(t, 'E', buf.Cells[5][5].Char)
	})

	t.Run("should not write at clip exclusive boundary", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 20)
		clip := ClipRect{X: 5, Y: 5, W: 10, H: 10}
		cell := Cell{Char: 'X'}

		buf.SetClipped(15, 15, cell, clip)

		assert.Equal(t, ' ', buf.Cells[15][15].Char)
	})
}

func TestBuffer_Fill(t *testing.T) {
	t.Parallel()

	t.Run("should fill entire buffer with cell", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(5, 5)
		cell := Cell{Char: '#', Background: style.Color("#0000FF")}

		buf.Fill(cell)

		for y := 0; y < buf.Height; y++ {
			for x := 0; x < buf.Width; x++ {
				assert.Equal(t, '#', buf.Cells[y][x].Char)
				assert.Equal(t, style.Color("#0000FF"), buf.Cells[y][x].Background)
			}
		}
	})

	t.Run("should overwrite existing cells", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(3, 3)
		buf.Cells[1][1] = Cell{Char: 'A'}

		buf.Fill(Cell{Char: 'Z'})

		assert.Equal(t, 'Z', buf.Cells[1][1].Char)
	})

	t.Run("should handle empty buffer", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(0, 0)

		buf.Fill(Cell{Char: 'X'})

		assert.Empty(t, buf.Cells)
	})
}

func TestBuffer_FillRect(t *testing.T) {
	t.Parallel()

	t.Run("should fill rectangular region", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: '*'}

		buf.FillRect(2, 3, 4, 3, cell)

		for y := 3; y < 6; y++ {
			for x := 2; x < 6; x++ {
				assert.Equal(t, '*', buf.Cells[y][x].Char)
			}
		}
		assert.Equal(t, ' ', buf.Cells[0][0].Char)
		assert.Equal(t, ' ', buf.Cells[9][9].Char)
	})

	t.Run("should clip to buffer bounds", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: 'X'}

		buf.FillRect(8, 8, 5, 5, cell)

		assert.Equal(t, 'X', buf.Cells[8][8].Char)
		assert.Equal(t, 'X', buf.Cells[9][9].Char)
	})

	t.Run("should handle rect starting outside buffer", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: 'X'}

		buf.FillRect(-2, -2, 5, 5, cell)

		assert.Equal(t, 'X', buf.Cells[0][0].Char)
		assert.Equal(t, 'X', buf.Cells[2][2].Char)
		assert.Equal(t, ' ', buf.Cells[3][3].Char)
	})

	t.Run("should handle zero dimensions", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: 'X'}

		buf.FillRect(5, 5, 0, 0, cell)

		for y := 0; y < buf.Height; y++ {
			for x := 0; x < buf.Width; x++ {
				assert.Equal(t, ' ', buf.Cells[y][x].Char)
			}
		}
	})

	t.Run("should fill single cell rect", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: 'S'}

		buf.FillRect(5, 5, 1, 1, cell)

		assert.Equal(t, 'S', buf.Cells[5][5].Char)
		assert.Equal(t, ' ', buf.Cells[5][4].Char)
		assert.Equal(t, ' ', buf.Cells[5][6].Char)
	})
}

func TestBuffer_FillHorizontal(t *testing.T) {
	t.Parallel()

	t.Run("should draw horizontal line", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: '-'}

		buf.FillHorizontal(2, 5, 5, cell)

		for x := 2; x < 7; x++ {
			assert.Equal(t, '-', buf.Cells[5][x].Char)
		}
		assert.Equal(t, ' ', buf.Cells[5][1].Char)
		assert.Equal(t, ' ', buf.Cells[5][7].Char)
	})

	t.Run("should clip to buffer width", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: '-'}

		buf.FillHorizontal(7, 5, 10, cell)

		assert.Equal(t, '-', buf.Cells[5][7].Char)
		assert.Equal(t, '-', buf.Cells[5][9].Char)
	})

	t.Run("should handle zero length", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: '-'}

		buf.FillHorizontal(5, 5, 0, cell)

		assert.Equal(t, ' ', buf.Cells[5][5].Char)
	})

	t.Run("should handle negative start position", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: '-'}

		buf.FillHorizontal(-3, 5, 5, cell)

		assert.Equal(t, '-', buf.Cells[5][0].Char)
		assert.Equal(t, '-', buf.Cells[5][1].Char)
		assert.Equal(t, ' ', buf.Cells[5][2].Char)
	})
}

func TestBuffer_FillVertical(t *testing.T) {
	t.Parallel()

	t.Run("should draw vertical line", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: '|'}

		buf.FillVertical(5, 2, 5, cell)

		for y := 2; y < 7; y++ {
			assert.Equal(t, '|', buf.Cells[y][5].Char)
		}
		assert.Equal(t, ' ', buf.Cells[1][5].Char)
		assert.Equal(t, ' ', buf.Cells[7][5].Char)
	})

	t.Run("should clip to buffer height", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: '|'}

		buf.FillVertical(5, 7, 10, cell)

		assert.Equal(t, '|', buf.Cells[7][5].Char)
		assert.Equal(t, '|', buf.Cells[9][5].Char)
	})

	t.Run("should handle zero length", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: '|'}

		buf.FillVertical(5, 5, 0, cell)

		assert.Equal(t, ' ', buf.Cells[5][5].Char)
	})

	t.Run("should handle negative start position", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 10)
		cell := Cell{Char: '|'}

		buf.FillVertical(5, -3, 5, cell)

		assert.Equal(t, '|', buf.Cells[0][5].Char)
		assert.Equal(t, '|', buf.Cells[1][5].Char)
		assert.Equal(t, ' ', buf.Cells[2][5].Char)
	})
}

// text.go tests

func TestBuffer_DrawText(t *testing.T) {
	t.Parallel()

	t.Run("should render text at position", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 5)
		fg := style.Color("#FFFFFF")
		bg := style.Color("#000000")

		buf.DrawText(5, 2, "Hello", fg, bg)

		assert.Equal(t, 'H', buf.Cells[2][5].Char)
		assert.Equal(t, 'e', buf.Cells[2][6].Char)
		assert.Equal(t, 'l', buf.Cells[2][7].Char)
		assert.Equal(t, 'l', buf.Cells[2][8].Char)
		assert.Equal(t, 'o', buf.Cells[2][9].Char)
		assert.Equal(t, fg, buf.Cells[2][5].Foreground)
		assert.Equal(t, bg, buf.Cells[2][5].Background)
	})

	t.Run("should clip text beyond buffer width", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)

		buf.DrawText(7, 2, "Hello", style.Color(""), style.Color(""))

		assert.Equal(t, 'H', buf.Cells[2][7].Char)
		assert.Equal(t, 'e', buf.Cells[2][8].Char)
		assert.Equal(t, 'l', buf.Cells[2][9].Char)
	})

	t.Run("should handle empty string", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)

		buf.DrawText(5, 2, "", style.Color(""), style.Color(""))

		assert.Equal(t, ' ', buf.Cells[2][5].Char)
	})

	t.Run("should handle unicode characters", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 5)

		buf.DrawText(0, 0, "日本語", style.Color(""), style.Color(""))

		assert.Equal(t, '日', buf.Cells[0][0].Char)
		assert.Equal(t, '本', buf.Cells[0][3].Char)
		assert.Equal(t, '語', buf.Cells[0][6].Char)
	})

	t.Run("should render at origin", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)

		buf.DrawText(0, 0, "ABC", style.Color(""), style.Color(""))

		assert.Equal(t, 'A', buf.Cells[0][0].Char)
		assert.Equal(t, 'B', buf.Cells[0][1].Char)
		assert.Equal(t, 'C', buf.Cells[0][2].Char)
	})
}

func TestBuffer_DrawStyledText(t *testing.T) {
	t.Parallel()

	t.Run("should render text with style attributes", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 5)
		s := style.Style{
			Foreground:     style.Color("#FF0000"),
			Background:     style.Color("#00FF00"),
			FontWeight:     style.WeightBold,
			FontStyle:      style.StyleItalic,
			TextDecoration: style.DecorationUnderline,
		}

		buf.DrawStyledText(0, 0, "Test", s)

		cell := buf.Cells[0][0]
		assert.Equal(t, 'T', cell.Char)
		assert.Equal(t, style.Color("#FF0000"), cell.Foreground)
		assert.Equal(t, style.Color("#00FF00"), cell.Background)
		assert.True(t, cell.Bold)
		assert.True(t, cell.Italic)
		assert.True(t, cell.Underline)
	})

	t.Run("should apply text transform uppercase", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 5)
		s := style.Style{
			TextTransform: style.TransformUppercase,
		}

		buf.DrawStyledText(0, 0, "hello", s)

		assert.Equal(t, 'H', buf.Cells[0][0].Char)
		assert.Equal(t, 'E', buf.Cells[0][1].Char)
		assert.Equal(t, 'L', buf.Cells[0][2].Char)
		assert.Equal(t, 'L', buf.Cells[0][3].Char)
		assert.Equal(t, 'O', buf.Cells[0][4].Char)
	})

	t.Run("should apply text transform lowercase", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 5)
		s := style.Style{
			TextTransform: style.TransformLowercase,
		}

		buf.DrawStyledText(0, 0, "HELLO", s)

		assert.Equal(t, 'h', buf.Cells[0][0].Char)
		assert.Equal(t, 'e', buf.Cells[0][1].Char)
	})

	t.Run("should render with dim and reverse attributes", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 5)
		s := style.Style{
			Dim:     true,
			Reverse: true,
		}

		buf.DrawStyledText(0, 0, "X", s)

		cell := buf.Cells[0][0]
		assert.True(t, cell.Dim)
		assert.True(t, cell.Reverse)
	})

	t.Run("should render with strikethrough decoration", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 5)
		s := style.Style{
			TextDecoration: style.DecorationLineThrough,
		}

		buf.DrawStyledText(0, 0, "Strike", s)

		assert.True(t, buf.Cells[0][0].Strike)
	})
}

func Test_calculateTextOffset(t *testing.T) {
	t.Parallel()

	t.Run("should return 0 for left alignment", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(10, 20, style.TextAlignLeft)

		assert.Equal(t, 0, offset)
	})

	t.Run("should center text correctly", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(10, 20, style.TextAlignCenter)

		assert.Equal(t, 5, offset)
	})

	t.Run("should right align text", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(10, 20, style.TextAlignRight)

		assert.Equal(t, 10, offset)
	})

	t.Run("should return 0 when line fills container", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(20, 20, style.TextAlignCenter)

		assert.Equal(t, 0, offset)
	})

	t.Run("should return 0 when line exceeds container", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(25, 20, style.TextAlignCenter)

		assert.Equal(t, 0, offset)
	})

	t.Run("should return 0 for zero container width", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(10, 0, style.TextAlignCenter)

		assert.Equal(t, 0, offset)
	})

	t.Run("should return 0 for negative container width", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(10, -5, style.TextAlignCenter)

		assert.Equal(t, 0, offset)
	})

	t.Run("should handle odd centering", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(5, 10, style.TextAlignCenter)

		assert.Equal(t, 2, offset)
	})

	t.Run("should return 0 for unknown alignment", func(t *testing.T) {
		t.Parallel()

		offset := calculateTextOffset(10, 20, style.TextAlign(99))

		assert.Equal(t, 0, offset)
	})
}

func Test_applyEllipsis(t *testing.T) {
	t.Parallel()

	t.Run("should return unchanged when fits", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := applyEllipsis(runes, 10)

		assert.Equal(t, []rune("hello"), result)
	})

	t.Run("should return unchanged when exact fit", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := applyEllipsis(runes, 5)

		assert.Equal(t, []rune("hello"), result)
	})

	t.Run("should truncate with ellipsis", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello world")

		result := applyEllipsis(runes, 8)

		assert.Equal(t, []rune("hello..."), result)
	})

	t.Run("should handle maxWidth of 3", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := applyEllipsis(runes, 3)

		assert.Equal(t, []rune("..."), result)
	})

	t.Run("should handle maxWidth of 2", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := applyEllipsis(runes, 2)

		assert.Equal(t, []rune(".."), result)
	})

	t.Run("should handle maxWidth of 1", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := applyEllipsis(runes, 1)

		assert.Equal(t, []rune("."), result)
	})

	t.Run("should handle unicode text", func(t *testing.T) {
		t.Parallel()

		runes := []rune("日本語テスト")

		result := applyEllipsis(runes, 5)

		assert.Equal(t, []rune("日本..."), result)
	})
}

func Test_wrapText(t *testing.T) {
	t.Parallel()

	t.Run("should return single line when no wrapping needed", func(t *testing.T) {
		t.Parallel()

		lines := wrapText("hello", 10, style.WrapWord)

		assert.Equal(t, []string{"hello"}, lines)
	})

	t.Run("should return unchanged for zero width", func(t *testing.T) {
		t.Parallel()

		lines := wrapText("hello", 0, style.WrapWord)

		assert.Equal(t, []string{"hello"}, lines)
	})

	t.Run("should return unchanged for negative width", func(t *testing.T) {
		t.Parallel()

		lines := wrapText("hello", -5, style.WrapWord)

		assert.Equal(t, []string{"hello"}, lines)
	})

	t.Run("should use character wrapping when WrapChar", func(t *testing.T) {
		t.Parallel()

		lines := wrapText("abcdefghij", 4, style.WrapChar)

		assert.Equal(t, []string{"abcd", "efgh", "ij"}, lines)
	})

	t.Run("should use word wrapping for WrapWord", func(t *testing.T) {
		t.Parallel()

		lines := wrapText("hello world", 6, style.WrapWord)

		assert.Equal(t, []string{"hello", "world"}, lines)
	})
}

func Test_wrapByWord(t *testing.T) {
	t.Parallel()

	t.Run("should wrap at word boundaries", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("hello world foo", 10)

		assert.Equal(t, []string{"hello", "world foo"}, lines)
	})

	t.Run("should handle single word", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("hello", 10)

		assert.Equal(t, []string{"hello"}, lines)
	})

	t.Run("should handle empty string", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("", 10)

		assert.Equal(t, []string{""}, lines)
	})

	t.Run("should break long word with character wrapping", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("abcdefghij", 4)

		assert.Equal(t, []string{"abcd", "efgh", "ij"}, lines)
	})

	t.Run("should handle word longer than maxWidth in middle", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("hi superlongword bye", 5)

		require.Len(t, lines, 5)
		assert.Equal(t, "hi", lines[0])
		assert.Equal(t, "super", lines[1])
		assert.Equal(t, "longw", lines[2])
		assert.Equal(t, "ord", lines[3])
		assert.Equal(t, "bye", lines[4])
	})

	t.Run("should join words on same line when space permits", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("a b c d", 5)

		assert.Equal(t, []string{"a b c", "d"}, lines)
	})

	t.Run("should handle multiple spaces between words", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("hello   world", 10)

		assert.Equal(t, []string{"hello", "world"}, lines)
	})

	t.Run("should handle only whitespace", func(t *testing.T) {
		t.Parallel()

		lines := wrapByWord("   ", 10)

		assert.Equal(t, []string{""}, lines)
	})
}

func Test_wrapByChar(t *testing.T) {
	t.Parallel()

	t.Run("should wrap at character boundaries", func(t *testing.T) {
		t.Parallel()

		lines := wrapByChar("abcdefghij", 3)

		assert.Equal(t, []string{"abc", "def", "ghi", "j"}, lines)
	})

	t.Run("should handle text shorter than width", func(t *testing.T) {
		t.Parallel()

		lines := wrapByChar("hi", 10)

		assert.Equal(t, []string{"hi"}, lines)
	})

	t.Run("should handle exact width", func(t *testing.T) {
		t.Parallel()

		lines := wrapByChar("abcd", 4)

		assert.Equal(t, []string{"abcd"}, lines)
	})

	t.Run("should handle empty string", func(t *testing.T) {
		t.Parallel()

		lines := wrapByChar("", 5)

		assert.Empty(t, lines)
	})

	t.Run("should handle unicode characters", func(t *testing.T) {
		t.Parallel()

		lines := wrapByChar("日本語テスト", 2)

		assert.Equal(t, []string{"日本", "語テ", "スト"}, lines)
	})

	t.Run("should handle width of 1", func(t *testing.T) {
		t.Parallel()

		lines := wrapByChar("abc", 1)

		assert.Equal(t, []string{"a", "b", "c"}, lines)
	})
}

// border.go tests

func TestBuffer_DrawBorder(t *testing.T) {
	t.Parallel()

	t.Run("should draw single border", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)
		border := style.BorderAll(style.BorderSingle)

		buf.DrawBorder(0, 0, 10, 5, border)

		assert.Equal(t, '┌', buf.Cells[0][0].Char)
		assert.Equal(t, '┐', buf.Cells[0][9].Char)
		assert.Equal(t, '└', buf.Cells[4][0].Char)
		assert.Equal(t, '┘', buf.Cells[4][9].Char)
		assert.Equal(t, '─', buf.Cells[0][5].Char)
		assert.Equal(t, '│', buf.Cells[2][0].Char)
	})

	t.Run("should draw round border", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)
		border := style.BorderAll(style.BorderRound)

		buf.DrawBorder(0, 0, 10, 5, border)

		assert.Equal(t, '╭', buf.Cells[0][0].Char)
		assert.Equal(t, '╮', buf.Cells[0][9].Char)
		assert.Equal(t, '╰', buf.Cells[4][0].Char)
		assert.Equal(t, '╯', buf.Cells[4][9].Char)
	})

	t.Run("should draw double border", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)
		border := style.BorderAll(style.BorderDouble)

		buf.DrawBorder(0, 0, 10, 5, border)

		assert.Equal(t, '╔', buf.Cells[0][0].Char)
		assert.Equal(t, '╗', buf.Cells[0][9].Char)
		assert.Equal(t, '╚', buf.Cells[4][0].Char)
		assert.Equal(t, '╝', buf.Cells[4][9].Char)
		assert.Equal(t, '═', buf.Cells[0][5].Char)
		assert.Equal(t, '║', buf.Cells[2][0].Char)
	})

	t.Run("should draw bold border", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)
		border := style.BorderAll(style.BorderBold)

		buf.DrawBorder(0, 0, 10, 5, border)

		assert.Equal(t, '┏', buf.Cells[0][0].Char)
		assert.Equal(t, '┓', buf.Cells[0][9].Char)
		assert.Equal(t, '┗', buf.Cells[4][0].Char)
		assert.Equal(t, '┛', buf.Cells[4][9].Char)
		assert.Equal(t, '━', buf.Cells[0][5].Char)
		assert.Equal(t, '┃', buf.Cells[2][0].Char)
	})

	t.Run("should apply border color", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)
		border := style.BorderAllWithColor(style.BorderSingle, style.Color("#FF0000"))

		buf.DrawBorder(0, 0, 10, 5, border)

		assert.Equal(t, style.Color("#FF0000"), buf.Cells[0][0].Foreground)
		assert.Equal(t, style.Color("#FF0000"), buf.Cells[0][5].Foreground)
		assert.Equal(t, style.Color("#FF0000"), buf.Cells[2][0].Foreground)
	})

	t.Run("should not draw border with no style", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)
		border := style.Border{}

		buf.DrawBorder(0, 0, 10, 5, border)

		assert.Equal(t, ' ', buf.Cells[0][0].Char)
	})

	t.Run("should not draw border when width less than 2", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)
		border := style.BorderAll(style.BorderSingle)

		buf.DrawBorder(0, 0, 1, 5, border)

		assert.Equal(t, ' ', buf.Cells[0][0].Char)
	})

	t.Run("should not draw border when height less than 2", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)
		border := style.BorderAll(style.BorderSingle)

		buf.DrawBorder(0, 0, 10, 1, border)

		assert.Equal(t, ' ', buf.Cells[0][0].Char)
	})

	t.Run("should draw minimum 2x2 border", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)
		border := style.BorderAll(style.BorderSingle)

		buf.DrawBorder(0, 0, 2, 2, border)

		assert.Equal(t, '┌', buf.Cells[0][0].Char)
		assert.Equal(t, '┐', buf.Cells[0][1].Char)
		assert.Equal(t, '└', buf.Cells[1][0].Char)
		assert.Equal(t, '┘', buf.Cells[1][1].Char)
	})

	t.Run("should draw border at offset position", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(15, 10)
		border := style.BorderAll(style.BorderSingle)

		buf.DrawBorder(3, 2, 8, 5, border)

		assert.Equal(t, '┌', buf.Cells[2][3].Char)
		assert.Equal(t, '┐', buf.Cells[2][10].Char)
		assert.Equal(t, '└', buf.Cells[6][3].Char)
		assert.Equal(t, '┘', buf.Cells[6][10].Char)
		assert.Equal(t, ' ', buf.Cells[0][0].Char)
	})

	t.Run("should draw border with only top and bottom", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)
		edge := style.BorderEdge{Style: style.BorderSingle}
		border := style.Border{Top: edge, Bottom: edge}

		buf.DrawBorder(0, 0, 10, 5, border)

		assert.Equal(t, '─', buf.Cells[0][0].Char)
		assert.Equal(t, '─', buf.Cells[4][0].Char)
		assert.Equal(t, ' ', buf.Cells[2][0].Char)
	})

	t.Run("should draw border with only left and right", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)
		edge := style.BorderEdge{Style: style.BorderSingle}
		border := style.Border{Left: edge, Right: edge}

		buf.DrawBorder(0, 0, 10, 5, border)

		assert.Equal(t, '│', buf.Cells[0][0].Char)
		assert.Equal(t, '│', buf.Cells[0][9].Char)
		assert.Equal(t, ' ', buf.Cells[0][5].Char)
	})
}

func TestBuffer_DrawBorder_WithTitle(t *testing.T) {
	t.Parallel()

	t.Run("should render left-aligned title on top border", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 5)
		border := style.BorderAllWithTitle(style.BorderSingle, style.Color("#FFFFFF"), "Title")

		buf.DrawBorder(0, 0, 20, 5, border)

		assert.Equal(t, '┌', buf.Cells[0][0].Char)
		assert.Equal(t, '─', buf.Cells[0][1].Char)
		assert.Equal(t, ' ', buf.Cells[0][2].Char)
		assert.Equal(t, 'T', buf.Cells[0][3].Char)
		assert.Equal(t, 'i', buf.Cells[0][4].Char)
		assert.Equal(t, 't', buf.Cells[0][5].Char)
		assert.Equal(t, 'l', buf.Cells[0][6].Char)
		assert.Equal(t, 'e', buf.Cells[0][7].Char)
		assert.Equal(t, ' ', buf.Cells[0][8].Char)
		assert.Equal(t, '─', buf.Cells[0][9].Char)
	})

	t.Run("should render right-aligned text on border", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 5)
		topEdge := style.BorderEdge{Style: style.BorderSingle}
		topEdge = topEdge.AddText("Right", style.TextAlignRight)
		border := style.BorderAll(style.BorderSingle)
		border.Top = topEdge

		buf.DrawBorder(0, 0, 20, 5, border)

		assert.Equal(t, 't', buf.Cells[0][16].Char)
		assert.Equal(t, ' ', buf.Cells[0][17].Char)
		assert.Equal(t, '┐', buf.Cells[0][19].Char)
	})

	t.Run("should render center-aligned text on border", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(20, 5)
		topEdge := style.BorderEdge{Style: style.BorderSingle}
		topEdge = topEdge.AddText("Hi", style.TextAlignCenter)
		border := style.BorderAll(style.BorderSingle)
		border.Top = topEdge

		buf.DrawBorder(0, 0, 20, 5, border)

		foundH := false
		for x := 1; x < 19; x++ {
			if buf.Cells[0][x].Char == 'H' {
				foundH = true
				assert.Equal(t, 'i', buf.Cells[0][x+1].Char)
				break
			}
		}
		assert.True(t, foundH, "should find 'Hi' text in border")
	})
}

func Test_truncateRunes(t *testing.T) {
	t.Parallel()

	t.Run("should return unchanged when fits", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunes(runes, 10)

		assert.Equal(t, []rune("hello"), result)
	})

	t.Run("should return unchanged when exact fit", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunes(runes, 5)

		assert.Equal(t, []rune("hello"), result)
	})

	t.Run("should truncate with ellipsis when too long", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello world")

		result := truncateRunes(runes, 8)

		assert.Equal(t, []rune("hello..."), result)
	})

	t.Run("should truncate without ellipsis when maxLen <= 3", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunes(runes, 3)

		assert.Equal(t, []rune("hel"), result)
	})

	t.Run("should truncate to 2 chars without ellipsis", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunes(runes, 2)

		assert.Equal(t, []rune("he"), result)
	})

	t.Run("should truncate to 1 char", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunes(runes, 1)

		assert.Equal(t, []rune("h"), result)
	})

	t.Run("should handle unicode", func(t *testing.T) {
		t.Parallel()

		runes := []rune("日本語テスト")

		result := truncateRunes(runes, 5)

		assert.Equal(t, []rune("日本..."), result)
	})
}

func Test_truncateRunesFromStart(t *testing.T) {
	t.Parallel()

	t.Run("should return unchanged when fits", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunesFromStart(runes, 10)

		assert.Equal(t, []rune("hello"), result)
	})

	t.Run("should return unchanged when exact fit", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunesFromStart(runes, 5)

		assert.Equal(t, []rune("hello"), result)
	})

	t.Run("should truncate from start with ellipsis", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello world")

		result := truncateRunesFromStart(runes, 8)

		assert.Equal(t, []rune("...world"), result)
	})

	t.Run("should truncate without ellipsis when maxLen <= 3", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunesFromStart(runes, 3)

		assert.Equal(t, []rune("llo"), result)
	})

	t.Run("should truncate to last 2 chars", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunesFromStart(runes, 2)

		assert.Equal(t, []rune("lo"), result)
	})

	t.Run("should truncate to last 1 char", func(t *testing.T) {
		t.Parallel()

		runes := []rune("hello")

		result := truncateRunesFromStart(runes, 1)

		assert.Equal(t, []rune("o"), result)
	})

	t.Run("should handle unicode", func(t *testing.T) {
		t.Parallel()

		runes := []rune("日本語テスト")

		result := truncateRunesFromStart(runes, 5)

		assert.Equal(t, []rune("...スト"), result)
	})
}
