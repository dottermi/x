package render

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColor(t *testing.T) {
	t.Parallel()
	t.Run("RGB creates set color", func(t *testing.T) {
		t.Parallel()
		c := RGB(255, 128, 0)

		assert.True(t, c.IsSet())
		assert.Equal(t, uint8(255), c.R)
		assert.Equal(t, uint8(128), c.G)
		assert.Equal(t, uint8(0), c.B)
	})

	t.Run("Default creates unset color", func(t *testing.T) {
		t.Parallel()
		c := Default()

		assert.False(t, c.IsSet())
	})

	t.Run("FGCode returns correct ANSI", func(t *testing.T) {
		t.Parallel()
		c := RGB(255, 0, 0)

		assert.Equal(t, "\x1b[38;2;255;0;0m", c.FGCode())
	})

	t.Run("FGCode returns default for unset", func(t *testing.T) {
		t.Parallel()
		c := Default()

		assert.Equal(t, "\x1b[39m", c.FGCode())
	})

	t.Run("BGCode returns correct ANSI", func(t *testing.T) {
		t.Parallel()
		c := RGB(0, 255, 0)

		assert.Equal(t, "\x1b[48;2;0;255;0m", c.BGCode())
	})

	t.Run("BGCode returns default for unset", func(t *testing.T) {
		t.Parallel()
		c := Default()

		assert.Equal(t, "\x1b[49m", c.BGCode())
	})
}

func TestCell(t *testing.T) {
	t.Parallel()
	t.Run("EmptyCell returns space with defaults", func(t *testing.T) {
		t.Parallel()
		c := EmptyCell()

		assert.Equal(t, ' ', c.Char)
		assert.False(t, c.FG.IsSet())
		assert.False(t, c.BG.IsSet())
	})

	t.Run("Equal returns true for identical cells", func(t *testing.T) {
		t.Parallel()
		c1 := Cell{Char: 'A', FG: RGB(255, 0, 0), Bold: true}
		c2 := Cell{Char: 'A', FG: RGB(255, 0, 0), Bold: true}

		assert.True(t, c1.Equal(c2))
	})

	t.Run("Equal returns false for different cells", func(t *testing.T) {
		t.Parallel()
		c1 := Cell{Char: 'A', FG: RGB(255, 0, 0), Bold: true}
		c2 := Cell{Char: 'A', FG: RGB(255, 0, 0), Bold: false}

		assert.False(t, c1.Equal(c2))
	})

	t.Run("Equal checks all attributes", func(t *testing.T) {
		t.Parallel()
		base := Cell{Char: 'X'}

		assert.False(t, base.Equal(Cell{Char: 'Y'}))
		assert.False(t, base.Equal(Cell{Char: 'X', Bold: true}))
		assert.False(t, base.Equal(Cell{Char: 'X', Dim: true}))
		assert.False(t, base.Equal(Cell{Char: 'X', Italic: true}))
		assert.False(t, base.Equal(Cell{Char: 'X', Underline: true}))
		assert.False(t, base.Equal(Cell{Char: 'X', Strike: true}))
		assert.False(t, base.Equal(Cell{Char: 'X', Reverse: true}))
	})
}

func TestBuffer(t *testing.T) {
	t.Parallel()
	t.Run("NewBuffer creates correct dimensions", func(t *testing.T) {
		t.Parallel()
		buf := NewBuffer(80, 24)

		assert.Equal(t, 80, buf.Width)
		assert.Equal(t, 24, buf.Height)
		assert.Len(t, buf.Cells, 80*24)
	})

	t.Run("NewBuffer fills with empty cells", func(t *testing.T) {
		t.Parallel()
		buf := NewBuffer(5, 5)

		for i := range buf.Cells {
			assert.Equal(t, ' ', buf.Cells[i].Char)
		}
	})

	t.Run("Set and Get work correctly", func(t *testing.T) {
		t.Parallel()
		buf := NewBuffer(10, 10)
		cell := Cell{Char: 'X', FG: RGB(255, 0, 0)}

		buf.Set(5, 5, cell)

		assert.True(t, buf.Get(5, 5).Equal(cell))
	})

	t.Run("Set ignores out of bounds", func(t *testing.T) {
		t.Parallel()
		buf := NewBuffer(10, 10)

		assert.NotPanics(t, func() {
			buf.Set(-1, 0, Cell{Char: 'X'})
			buf.Set(100, 0, Cell{Char: 'X'})
			buf.Set(0, -1, Cell{Char: 'X'})
			buf.Set(0, 100, Cell{Char: 'X'})
		})
	})

	t.Run("Get returns EmptyCell for out of bounds", func(t *testing.T) {
		t.Parallel()
		buf := NewBuffer(10, 10)

		assert.True(t, buf.Get(-1, 0).Equal(EmptyCell()))
		assert.True(t, buf.Get(100, 0).Equal(EmptyCell()))
	})

	t.Run("Fill sets all cells", func(t *testing.T) {
		t.Parallel()
		buf := NewBuffer(5, 5)
		cell := Cell{Char: 'X', Bold: true}

		buf.Fill(cell)

		for i := range buf.Cells {
			assert.True(t, buf.Cells[i].Equal(cell))
		}
	})

	t.Run("Clone creates independent copy", func(t *testing.T) {
		t.Parallel()
		buf := NewBuffer(10, 10)
		buf.Set(0, 0, Cell{Char: 'A'})

		clone := buf.Clone()
		clone.Set(0, 0, Cell{Char: 'B'})

		assert.Equal(t, 'A', buf.Get(0, 0).Char)
		assert.Equal(t, 'B', clone.Get(0, 0).Char)
	})

	t.Run("Diff returns changed cells", func(t *testing.T) {
		t.Parallel()
		old := NewBuffer(10, 10)
		next := NewBuffer(10, 10)
		next.Set(0, 0, Cell{Char: 'X'})
		next.Set(5, 5, Cell{Char: 'Y'})

		changes := old.Diff(next)

		assert.Len(t, changes, 2)
	})

	t.Run("Diff returns empty for identical buffers", func(t *testing.T) {
		t.Parallel()
		buf1 := NewBuffer(10, 10)
		buf2 := buf1.Clone()

		changes := buf1.Diff(buf2)

		assert.Empty(t, changes)
	})

	t.Run("Diff includes correct positions", func(t *testing.T) {
		t.Parallel()
		old := NewBuffer(10, 10)
		next := NewBuffer(10, 10)
		next.Set(3, 7, Cell{Char: 'Z'})

		changes := old.Diff(next)

		require.Len(t, changes, 1)
		assert.Equal(t, 3, changes[0].X)
		assert.Equal(t, 7, changes[0].Y)
		assert.Equal(t, 'Z', changes[0].Cell.Char)
	})
}

func TestTerminal(t *testing.T) {
	t.Parallel()
	t.Run("NewTerminal creates with dimensions", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(80, 24)

		assert.Equal(t, 80, term.width)
		assert.Equal(t, 24, term.height)
		assert.NotNil(t, term.current)
	})

	t.Run("Render returns empty for no changes", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(10, 10)
		buf := NewBuffer(10, 10)

		output := term.Render(buf)

		assert.Empty(t, output)
	})

	t.Run("Render includes MoveCursor", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(10, 10)
		buf := NewBuffer(10, 10)
		buf.Set(5, 5, Cell{Char: 'X'})

		output := term.Render(buf)

		assert.Contains(t, output, "\x1b[6;6H")
		assert.Contains(t, output, "X")
	})

	t.Run("Render includes color codes", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(10, 10)
		buf := NewBuffer(10, 10)
		buf.Set(0, 0, Cell{Char: 'R', FG: RGB(255, 0, 0)})

		output := term.Render(buf)

		assert.Contains(t, output, "\x1b[38;2;255;0;0m")
	})

	t.Run("Render includes background color", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(10, 10)
		buf := NewBuffer(10, 10)
		buf.Set(0, 0, Cell{Char: 'B', BG: RGB(0, 0, 255)})

		output := term.Render(buf)

		assert.Contains(t, output, "\x1b[48;2;0;0;255m")
	})

	t.Run("Render includes attribute codes", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(10, 10)
		buf := NewBuffer(10, 10)
		buf.Set(0, 0, Cell{Char: 'B', Bold: true})

		output := term.Render(buf)

		assert.Contains(t, output, BoldOn)
	})

	t.Run("Render ends with Reset", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(10, 10)
		buf := NewBuffer(10, 10)
		buf.Set(0, 0, Cell{Char: 'X'})

		output := term.Render(buf)

		assert.True(t, strings.HasSuffix(output, Reset))
	})

	t.Run("RenderFull renders entire buffer", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(3, 2)
		buf := NewBuffer(3, 2)
		buf.Set(0, 0, Cell{Char: 'A'})
		buf.Set(1, 0, Cell{Char: 'B'})
		buf.Set(2, 0, Cell{Char: 'C'})
		buf.Set(0, 1, Cell{Char: 'D'})
		buf.Set(1, 1, Cell{Char: 'E'})
		buf.Set(2, 1, Cell{Char: 'F'})

		output := term.RenderFull(buf)

		assert.Contains(t, output, "A")
		assert.Contains(t, output, "B")
		assert.Contains(t, output, "C")
		assert.Contains(t, output, "D")
		assert.Contains(t, output, "E")
		assert.Contains(t, output, "F")
	})

	t.Run("RenderFull starts with cursor home", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(5, 5)
		buf := NewBuffer(5, 5)

		output := term.RenderFull(buf)

		assert.True(t, strings.HasPrefix(output, MoveCursor(0, 0)))
	})

	t.Run("Clear returns clear screen codes", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(10, 10)

		output := term.Clear()

		assert.Contains(t, output, ClearScreen)
		assert.Contains(t, output, CursorHome)
	})

	t.Run("Clear resets internal buffer", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(10, 10)
		buf := NewBuffer(10, 10)
		buf.Set(0, 0, Cell{Char: 'X'})
		term.Render(buf)

		term.Clear()

		// After clear, rendering the same buffer should show changes again
		buf2 := NewBuffer(10, 10)
		buf2.Set(0, 0, Cell{Char: 'X'})
		output := term.Render(buf2)
		assert.Contains(t, output, "X")
	})

	t.Run("Resize updates dimensions", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(10, 10)

		term.Resize(20, 15)

		assert.Equal(t, 20, term.width)
		assert.Equal(t, 15, term.height)
	})

	t.Run("Resize clears buffer", func(t *testing.T) {
		t.Parallel()
		term := NewTerminal(10, 10)
		buf := NewBuffer(10, 10)
		buf.Set(0, 0, Cell{Char: 'X'})
		term.Render(buf)

		term.Resize(10, 10)

		// After resize, should detect change again
		buf2 := NewBuffer(10, 10)
		buf2.Set(0, 0, Cell{Char: 'X'})
		output := term.Render(buf2)
		assert.Contains(t, output, "X")
	})
}

func TestANSICodes(t *testing.T) {
	t.Parallel()
	t.Run("MoveCursor generates correct code", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "\x1b[1;1H", MoveCursor(0, 0))
		assert.Equal(t, "\x1b[5;10H", MoveCursor(9, 4))
		assert.Equal(t, "\x1b[25;80H", MoveCursor(79, 24))
	})

	t.Run("HideCursor returns correct code", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, CursorHide, HideCursor())
	})

	t.Run("ShowCursor returns correct code", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, CursorShow, ShowCursor())
	})

	t.Run("Constants have correct values", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "\x1b[0m", Reset)
		assert.Equal(t, "\x1b[2J", ClearScreen)
		assert.Equal(t, "\x1b[H", CursorHome)
		assert.Equal(t, "\x1b[?25l", CursorHide)
		assert.Equal(t, "\x1b[?25h", CursorShow)
	})

	t.Run("Attribute codes are correct", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "\x1b[1m", BoldOn)
		assert.Equal(t, "\x1b[22m", BoldOff)
		assert.Equal(t, "\x1b[2m", DimOn)
		assert.Equal(t, "\x1b[3m", ItalicOn)
		assert.Equal(t, "\x1b[23m", ItalicOff)
		assert.Equal(t, "\x1b[4m", UnderlineOn)
		assert.Equal(t, "\x1b[24m", UnderlineOff)
		assert.Equal(t, "\x1b[7m", ReverseOn)
		assert.Equal(t, "\x1b[27m", ReverseOff)
		assert.Equal(t, "\x1b[9m", StrikeOn)
		assert.Equal(t, "\x1b[29m", StrikeOff)
	})
}
