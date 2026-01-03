package termistyle

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dottermi/x/termistyle/draw"
	"github.com/dottermi/x/termistyle/layout"
	"github.com/dottermi/x/termistyle/style"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Type aliases verification

func TestTypeAliases(t *testing.T) {
	t.Parallel()

	t.Run("should alias Style from style package", func(t *testing.T) {
		t.Parallel()

		var s Style
		s.Width = 10
		s.Height = 5

		assert.Equal(t, 10, s.Width)
		assert.Equal(t, 5, s.Height)
	})

	t.Run("should alias Color from style package", func(t *testing.T) {
		t.Parallel()

		var c Color = "#FF0000"

		assert.Equal(t, Color("#FF0000"), c)
	})

	t.Run("should alias Spacing from style package", func(t *testing.T) {
		t.Parallel()

		spacing := SpacingAll(5)

		assert.Equal(t, 5, spacing.Top)
		assert.Equal(t, 5, spacing.Right)
		assert.Equal(t, 5, spacing.Bottom)
		assert.Equal(t, 5, spacing.Left)
	})

	t.Run("should alias Box from layout package", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{Width: 10, Height: 5})

		assert.NotNil(t, box)
		assert.Equal(t, 10, box.Style.Width)
	})

	t.Run("should alias Buffer from draw package", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(10, 5)

		assert.NotNil(t, buf)
		assert.Equal(t, 10, buf.Width)
		assert.Equal(t, 5, buf.Height)
	})

	t.Run("should alias Cell from draw package", func(t *testing.T) {
		t.Parallel()

		var c Cell
		c.Char = 'X'
		c.Foreground = "#FFFFFF"

		assert.Equal(t, 'X', c.Char)
		assert.Equal(t, Color("#FFFFFF"), c.Foreground)
	})

	t.Run("should alias Border from style package", func(t *testing.T) {
		t.Parallel()

		border := BorderAll(BorderSingle)

		assert.True(t, border.Top.IsSet())
		assert.True(t, border.Right.IsSet())
		assert.True(t, border.Bottom.IsSet())
		assert.True(t, border.Left.IsSet())
	})

	t.Run("should alias FontWeight from style package", func(t *testing.T) {
		t.Parallel()

		w := WeightBold

		assert.Equal(t, WeightBold, w)
	})

	t.Run("should alias FontStyle from style package", func(t *testing.T) {
		t.Parallel()

		fs := StyleItalic

		assert.Equal(t, StyleItalic, fs)
	})

	t.Run("should alias TextDecoration from style package", func(t *testing.T) {
		t.Parallel()

		td := DecorationUnderline

		assert.Equal(t, DecorationUnderline, td)
	})

	t.Run("should alias TextTransform from style package", func(t *testing.T) {
		t.Parallel()

		tt := TransformUppercase

		assert.Equal(t, TransformUppercase, tt)
	})

	t.Run("should alias TextAlign from style package", func(t *testing.T) {
		t.Parallel()

		ta := TextAlignCenter

		assert.Equal(t, TextAlignCenter, ta)
	})

	t.Run("should alias TextWrap from style package", func(t *testing.T) {
		t.Parallel()

		tw := WrapWord

		assert.Equal(t, WrapWord, tw)
	})

	t.Run("should alias TextOverflow from style package", func(t *testing.T) {
		t.Parallel()

		to := TextOverflowEllipsis

		assert.Equal(t, TextOverflowEllipsis, to)
	})

	t.Run("should alias Overflow from style package", func(t *testing.T) {
		t.Parallel()

		o := OverflowHidden

		assert.Equal(t, OverflowHidden, o)
	})

	t.Run("should alias Justify from style package", func(t *testing.T) {
		t.Parallel()

		j := JustifyCenter

		assert.Equal(t, JustifyCenter, j)
	})

	t.Run("should alias Align from style package", func(t *testing.T) {
		t.Parallel()

		a := AlignCenter

		assert.Equal(t, AlignCenter, a)
	})

	t.Run("should alias FlexWrap from style package", func(t *testing.T) {
		t.Parallel()

		fw := Wrap

		assert.Equal(t, Wrap, fw)
	})

	t.Run("should alias BorderEdge from style package", func(t *testing.T) {
		t.Parallel()

		var be BorderEdge
		be.Style = BorderSingle
		be.Color = "#FF0000"

		assert.True(t, be.IsSet())
	})

	t.Run("should alias BorderText from style package", func(t *testing.T) {
		t.Parallel()

		var bt BorderText
		bt.Text = "Title"
		bt.Align = TextAlignLeft

		assert.Equal(t, "Title", bt.Text)
	})
}

// Constants verification

func TestConstants(t *testing.T) {
	t.Parallel()

	t.Run("should export display mode constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.Flex, Flex)
		assert.Equal(t, style.Block, Block)
	})

	t.Run("should export position constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.Relative, Relative)
		assert.Equal(t, style.Absolute, Absolute)
	})

	t.Run("should export flex direction constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.Row, Row)
		assert.Equal(t, style.Column, Column)
	})

	t.Run("should export justify content constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.JustifyStart, JustifyStart)
		assert.Equal(t, style.JustifyCenter, JustifyCenter)
		assert.Equal(t, style.JustifyEnd, JustifyEnd)
		assert.Equal(t, style.JustifyBetween, JustifyBetween)
		assert.Equal(t, style.JustifyAround, JustifyAround)
		assert.Equal(t, style.SpaceBetween, SpaceBetween)
		assert.Equal(t, style.SpaceAround, SpaceAround)
	})

	t.Run("should export align items constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.AlignStart, AlignStart)
		assert.Equal(t, style.AlignCenter, AlignCenter)
		assert.Equal(t, style.AlignEnd, AlignEnd)
		assert.Equal(t, style.AlignStretch, AlignStretch)
		assert.Equal(t, style.Stretch, Stretch)
	})

	t.Run("should export font weight constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.WeightNormal, WeightNormal)
		assert.Equal(t, style.WeightBold, WeightBold)
	})

	t.Run("should export font style constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.StyleNormal, StyleNormal)
		assert.Equal(t, style.StyleItalic, StyleItalic)
		assert.Equal(t, style.StyleOblique, StyleOblique)
	})

	t.Run("should export text decoration constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.DecorationNone, DecorationNone)
		assert.Equal(t, style.DecorationUnderline, DecorationUnderline)
		assert.Equal(t, style.DecorationLineThrough, DecorationLineThrough)
	})

	t.Run("should export text transform constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.TransformNone, TransformNone)
		assert.Equal(t, style.TransformUppercase, TransformUppercase)
		assert.Equal(t, style.TransformLowercase, TransformLowercase)
		assert.Equal(t, style.TransformCapitalize, TransformCapitalize)
	})

	t.Run("should export text align constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.TextAlignLeft, TextAlignLeft)
		assert.Equal(t, style.TextAlignCenter, TextAlignCenter)
		assert.Equal(t, style.TextAlignRight, TextAlignRight)
	})

	t.Run("should export text wrap constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.WrapNone, WrapNone)
		assert.Equal(t, style.WrapWord, WrapWord)
		assert.Equal(t, style.WrapChar, WrapChar)
	})

	t.Run("should export text overflow constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.TextOverflowClip, TextOverflowClip)
		assert.Equal(t, style.TextOverflowEllipsis, TextOverflowEllipsis)
	})

	t.Run("should export overflow constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.OverflowVisible, OverflowVisible)
		assert.Equal(t, style.OverflowHidden, OverflowHidden)
	})

	t.Run("should export flex wrap constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, style.NoWrap, NoWrap)
		assert.Equal(t, style.Wrap, Wrap)
	})
}

// Draw function tests

func TestDraw(t *testing.T) {
	t.Parallel()

	t.Run("should create buffer with box dimensions", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{Width: 20, Height: 10})

		buf := Draw(box)

		require.NotNil(t, buf)
		assert.Equal(t, 20, buf.Width)
		assert.Equal(t, 10, buf.Height)
	})

	t.Run("should fill buffer with spaces by default", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{Width: 5, Height: 3})

		buf := Draw(box)

		for y := 0; y < buf.Height; y++ {
			for x := 0; x < buf.Width; x++ {
				assert.Equal(t, ' ', buf.Cells[y][x].Char)
			}
		}
	})

	t.Run("should draw box with background color", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{
			Width:      10,
			Height:     5,
			Background: Color("#FF0000"),
		})

		buf := Draw(box)

		assert.Equal(t, Color("#FF0000"), buf.Cells[2][5].Background)
	})

	t.Run("should draw box with single border", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{
			Width:  10,
			Height: 5,
			Border: BorderAll(BorderSingle),
		})

		buf := Draw(box)

		assert.Equal(t, rune(0x250C), buf.Cells[0][0].Char)
		assert.Equal(t, rune(0x2510), buf.Cells[0][9].Char)
		assert.Equal(t, rune(0x2514), buf.Cells[4][0].Char)
		assert.Equal(t, rune(0x2518), buf.Cells[4][9].Char)
		assert.Equal(t, rune(0x2500), buf.Cells[0][5].Char)
		assert.Equal(t, rune(0x2502), buf.Cells[2][0].Char)
	})

	t.Run("should draw box with children", func(t *testing.T) {
		t.Parallel()

		parent := NewBox(Style{
			Width:   20,
			Height:  10,
			Display: Block,
		})
		child := NewBox(Style{
			Width:      10,
			Height:     3,
			Background: Color("#00FF00"),
		})
		parent.AddChild(child)

		buf := Draw(parent)

		assert.Equal(t, Color("#00FF00"), buf.Cells[0][0].Background)
		assert.Equal(t, Color("#00FF00"), buf.Cells[2][9].Background)
	})

	t.Run("should draw box with text content", func(t *testing.T) {
		t.Parallel()

		box := Text("Hello", Style{
			Width:      10,
			Height:     1,
			Foreground: Color("#FFFFFF"),
		})

		buf := Draw(box)

		assert.Equal(t, 'H', buf.Cells[0][0].Char)
		assert.Equal(t, 'e', buf.Cells[0][1].Char)
		assert.Equal(t, 'l', buf.Cells[0][2].Char)
		assert.Equal(t, 'l', buf.Cells[0][3].Char)
		assert.Equal(t, 'o', buf.Cells[0][4].Char)
	})

	t.Run("should draw nested boxes with flex layout", func(t *testing.T) {
		t.Parallel()

		parent := NewBox(Style{
			Width:         30,
			Height:        10,
			Display:       Flex,
			FlexDirection: Row,
		})
		child1 := NewBox(Style{Width: 10, Height: 5, Background: Color("#FF0000")})
		child2 := NewBox(Style{Width: 10, Height: 5, Background: Color("#00FF00")})
		parent.AddChild(child1)
		parent.AddChild(child2)

		buf := Draw(parent)

		assert.Equal(t, Color("#FF0000"), buf.Cells[0][0].Background)
		assert.Equal(t, Color("#00FF00"), buf.Cells[0][10].Background)
	})

	t.Run("should draw border with background inside", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{
			Width:      10,
			Height:     5,
			Border:     BorderAll(BorderSingle),
			Background: Color("#0000FF"),
		})

		buf := Draw(box)

		assert.Equal(t, rune(0x250C), buf.Cells[0][0].Char)
		assert.Equal(t, Color("#0000FF"), buf.Cells[2][5].Background)
		assert.NotEqual(t, Color("#0000FF"), buf.Cells[0][0].Background)
	})

	t.Run("should calculate layout before drawing", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{Width: 15, Height: 8})

		buf := Draw(box)

		assert.Equal(t, 15, box.W)
		assert.Equal(t, 8, box.H)
		assert.Equal(t, 0, box.X)
		assert.Equal(t, 0, box.Y)
		assert.NotNil(t, buf)
	})
}

// Print function tests

func TestPrint(t *testing.T) {
	t.Parallel()

	t.Run("should render box to writer", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{Width: 5, Height: 2})
		var buf bytes.Buffer

		Print(&buf, box)

		output := buf.String()
		assert.NotEmpty(t, output)
		assert.Contains(t, output, " ")
	})

	t.Run("should render box with content to writer", func(t *testing.T) {
		t.Parallel()

		box := Text("Hi", Style{Width: 5, Height: 1})
		var buf bytes.Buffer

		Print(&buf, box)

		output := buf.String()
		assert.Contains(t, output, "H")
		assert.Contains(t, output, "i")
	})

	t.Run("should include ANSI codes for colors", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{
			Width:      5,
			Height:     1,
			Background: Color("#FF0000"),
		})
		var buf bytes.Buffer

		Print(&buf, box)

		output := buf.String()
		assert.Contains(t, output, "\x1b[")
	})

	t.Run("should render multiple rows", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{Width: 3, Height: 3})
		var buf bytes.Buffer

		Print(&buf, box)

		output := buf.String()
		lines := strings.Split(output, "\n")
		assert.GreaterOrEqual(t, len(lines), 2)
	})
}

// Println function tests

func TestPrintln(t *testing.T) {
	t.Parallel()

	t.Run("should render box to writer with trailing newline", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{Width: 5, Height: 1})
		var buf bytes.Buffer

		Println(&buf, box)

		output := buf.String()
		assert.True(t, strings.HasSuffix(output, "\n"))
	})

	t.Run("should add exactly one newline after output", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{Width: 3, Height: 1})
		var printBuf, printlnBuf bytes.Buffer

		Print(&printBuf, box)
		Println(&printlnBuf, box)

		printOutput := printBuf.String()
		printlnOutput := printlnBuf.String()

		assert.Equal(t, printOutput+"\n", printlnOutput)
	})

	t.Run("should render box content before newline", func(t *testing.T) {
		t.Parallel()

		box := Text("Test", Style{Width: 10, Height: 1})
		var buf bytes.Buffer

		Println(&buf, box)

		output := buf.String()
		assert.Contains(t, output, "T")
		assert.Contains(t, output, "e")
		assert.Contains(t, output, "s")
		assert.Contains(t, output, "t")
		assert.True(t, strings.HasSuffix(output, "\n"))
	})
}

// borderOffsets function tests

func TestBorderOffsets(t *testing.T) {
	t.Parallel()

	t.Run("should return all zeros for no borders", func(t *testing.T) {
		t.Parallel()

		border := style.Border{}

		top, right, bottom, left := borderOffsets(border)

		assert.Equal(t, 0, top)
		assert.Equal(t, 0, right)
		assert.Equal(t, 0, bottom)
		assert.Equal(t, 0, left)
	})

	t.Run("should return all ones for all borders", func(t *testing.T) {
		t.Parallel()

		border := BorderAll(BorderSingle)

		top, right, bottom, left := borderOffsets(border)

		assert.Equal(t, 1, top)
		assert.Equal(t, 1, right)
		assert.Equal(t, 1, bottom)
		assert.Equal(t, 1, left)
	})

	t.Run("should return correct values for top border only", func(t *testing.T) {
		t.Parallel()

		border := style.Border{
			Top: style.BorderEdge{Style: style.BorderSingle},
		}

		top, right, bottom, left := borderOffsets(border)

		assert.Equal(t, 1, top)
		assert.Equal(t, 0, right)
		assert.Equal(t, 0, bottom)
		assert.Equal(t, 0, left)
	})

	t.Run("should return correct values for right border only", func(t *testing.T) {
		t.Parallel()

		border := style.Border{
			Right: style.BorderEdge{Style: style.BorderSingle},
		}

		top, right, bottom, left := borderOffsets(border)

		assert.Equal(t, 0, top)
		assert.Equal(t, 1, right)
		assert.Equal(t, 0, bottom)
		assert.Equal(t, 0, left)
	})

	t.Run("should return correct values for bottom border only", func(t *testing.T) {
		t.Parallel()

		border := style.Border{
			Bottom: style.BorderEdge{Style: style.BorderSingle},
		}

		top, right, bottom, left := borderOffsets(border)

		assert.Equal(t, 0, top)
		assert.Equal(t, 0, right)
		assert.Equal(t, 1, bottom)
		assert.Equal(t, 0, left)
	})

	t.Run("should return correct values for left border only", func(t *testing.T) {
		t.Parallel()

		border := style.Border{
			Left: style.BorderEdge{Style: style.BorderSingle},
		}

		top, right, bottom, left := borderOffsets(border)

		assert.Equal(t, 0, top)
		assert.Equal(t, 0, right)
		assert.Equal(t, 0, bottom)
		assert.Equal(t, 1, left)
	})

	t.Run("should return correct values for top and bottom borders", func(t *testing.T) {
		t.Parallel()

		border := style.Border{
			Top:    style.BorderEdge{Style: style.BorderSingle},
			Bottom: style.BorderEdge{Style: style.BorderSingle},
		}

		top, right, bottom, left := borderOffsets(border)

		assert.Equal(t, 1, top)
		assert.Equal(t, 0, right)
		assert.Equal(t, 1, bottom)
		assert.Equal(t, 0, left)
	})

	t.Run("should return correct values for left and right borders", func(t *testing.T) {
		t.Parallel()

		border := style.Border{
			Left:  style.BorderEdge{Style: style.BorderSingle},
			Right: style.BorderEdge{Style: style.BorderSingle},
		}

		top, right, bottom, left := borderOffsets(border)

		assert.Equal(t, 0, top)
		assert.Equal(t, 1, right)
		assert.Equal(t, 0, bottom)
		assert.Equal(t, 1, left)
	})

	t.Run("should handle different border styles", func(t *testing.T) {
		t.Parallel()

		border := style.Border{
			Top:    style.BorderEdge{Style: style.BorderSingle},
			Right:  style.BorderEdge{Style: style.BorderDouble},
			Bottom: style.BorderEdge{Style: style.BorderRound},
			Left:   style.BorderEdge{Style: style.BorderBold},
		}

		top, right, bottom, left := borderOffsets(border)

		assert.Equal(t, 1, top)
		assert.Equal(t, 1, right)
		assert.Equal(t, 1, bottom)
		assert.Equal(t, 1, left)
	})

	t.Run("should return zero for BorderNone style", func(t *testing.T) {
		t.Parallel()

		border := style.Border{
			Top: style.BorderEdge{Style: style.BorderNone},
		}

		top, right, bottom, left := borderOffsets(border)

		assert.Equal(t, 0, top)
		assert.Equal(t, 0, right)
		assert.Equal(t, 0, bottom)
		assert.Equal(t, 0, left)
	})
}

// sortByZIndex function tests

func TestSortByZIndex(t *testing.T) {
	t.Parallel()

	t.Run("should return same order when no z-index set", func(t *testing.T) {
		t.Parallel()

		child1 := &layout.Box{Style: style.Style{Width: 1}}
		child2 := &layout.Box{Style: style.Style{Width: 2}}
		child3 := &layout.Box{Style: style.Style{Width: 3}}
		children := []*layout.Box{child1, child2, child3}

		result := sortByZIndex(children)

		assert.Equal(t, child1, result[0])
		assert.Equal(t, child2, result[1])
		assert.Equal(t, child3, result[2])
	})

	t.Run("should sort children by z-index ascending", func(t *testing.T) {
		t.Parallel()

		child1 := &layout.Box{Style: style.Style{ZIndex: 3}}
		child2 := &layout.Box{Style: style.Style{ZIndex: 1}}
		child3 := &layout.Box{Style: style.Style{ZIndex: 2}}
		children := []*layout.Box{child1, child2, child3}

		result := sortByZIndex(children)

		assert.Equal(t, child2, result[0])
		assert.Equal(t, child3, result[1])
		assert.Equal(t, child1, result[2])
	})

	t.Run("should maintain original order for equal z-index (stable sort)", func(t *testing.T) {
		t.Parallel()

		child1 := &layout.Box{Style: style.Style{Width: 1, ZIndex: 1}}
		child2 := &layout.Box{Style: style.Style{Width: 2, ZIndex: 1}}
		child3 := &layout.Box{Style: style.Style{Width: 3, ZIndex: 1}}
		children := []*layout.Box{child1, child2, child3}

		result := sortByZIndex(children)

		assert.Equal(t, child1, result[0])
		assert.Equal(t, child2, result[1])
		assert.Equal(t, child3, result[2])
	})

	t.Run("should return single child unchanged", func(t *testing.T) {
		t.Parallel()

		child := &layout.Box{Style: style.Style{ZIndex: 5}}
		children := []*layout.Box{child}

		result := sortByZIndex(children)

		require.Len(t, result, 1)
		assert.Equal(t, child, result[0])
	})

	t.Run("should return empty slice unchanged", func(t *testing.T) {
		t.Parallel()

		children := []*layout.Box{}

		result := sortByZIndex(children)

		assert.Empty(t, result)
	})

	t.Run("should return nil unchanged", func(t *testing.T) {
		t.Parallel()

		var children []*layout.Box

		result := sortByZIndex(children)

		assert.Nil(t, result)
	})

	t.Run("should handle mixed zero and non-zero z-index", func(t *testing.T) {
		t.Parallel()

		child1 := &layout.Box{Style: style.Style{Width: 1, ZIndex: 0}}
		child2 := &layout.Box{Style: style.Style{Width: 2, ZIndex: 2}}
		child3 := &layout.Box{Style: style.Style{Width: 3, ZIndex: 0}}
		child4 := &layout.Box{Style: style.Style{Width: 4, ZIndex: 1}}
		children := []*layout.Box{child1, child2, child3, child4}

		result := sortByZIndex(children)

		assert.Equal(t, 0, result[0].Style.ZIndex)
		assert.Equal(t, 0, result[1].Style.ZIndex)
		assert.Equal(t, 1, result[2].Style.ZIndex)
		assert.Equal(t, 2, result[3].Style.ZIndex)
	})

	t.Run("should handle negative z-index values", func(t *testing.T) {
		t.Parallel()

		child1 := &layout.Box{Style: style.Style{Width: 1, ZIndex: 0}}
		child2 := &layout.Box{Style: style.Style{Width: 2, ZIndex: -1}}
		child3 := &layout.Box{Style: style.Style{Width: 3, ZIndex: 1}}
		children := []*layout.Box{child1, child2, child3}

		result := sortByZIndex(children)

		assert.Equal(t, child2, result[0])
		assert.Equal(t, child1, result[1])
		assert.Equal(t, child3, result[2])
	})

	t.Run("should not modify original slice", func(t *testing.T) {
		t.Parallel()

		child1 := &layout.Box{Style: style.Style{ZIndex: 3}}
		child2 := &layout.Box{Style: style.Style{ZIndex: 1}}
		child3 := &layout.Box{Style: style.Style{ZIndex: 2}}
		children := []*layout.Box{child1, child2, child3}

		_ = sortByZIndex(children)

		assert.Equal(t, child1, children[0])
		assert.Equal(t, child2, children[1])
		assert.Equal(t, child3, children[2])
	})
}

// intersectClipRect function tests

func TestIntersectClipRect(t *testing.T) {
	t.Parallel()

	t.Run("should return intersection of overlapping rects", func(t *testing.T) {
		t.Parallel()

		a := draw.ClipRect{X: 0, Y: 0, W: 20, H: 20}
		b := draw.ClipRect{X: 10, Y: 10, W: 20, H: 20}

		result := intersectClipRect(a, b)

		assert.Equal(t, 10, result.X)
		assert.Equal(t, 10, result.Y)
		assert.Equal(t, 10, result.W)
		assert.Equal(t, 10, result.H)
	})

	t.Run("should return zero size for non-overlapping rects horizontally", func(t *testing.T) {
		t.Parallel()

		a := draw.ClipRect{X: 0, Y: 0, W: 10, H: 10}
		b := draw.ClipRect{X: 20, Y: 0, W: 10, H: 10}

		result := intersectClipRect(a, b)

		assert.Equal(t, 0, result.W)
	})

	t.Run("should return zero size for non-overlapping rects vertically", func(t *testing.T) {
		t.Parallel()

		a := draw.ClipRect{X: 0, Y: 0, W: 10, H: 10}
		b := draw.ClipRect{X: 0, Y: 20, W: 10, H: 10}

		result := intersectClipRect(a, b)

		assert.Equal(t, 0, result.H)
	})

	t.Run("should return contained rect when one contains the other", func(t *testing.T) {
		t.Parallel()

		outer := draw.ClipRect{X: 0, Y: 0, W: 100, H: 100}
		inner := draw.ClipRect{X: 20, Y: 20, W: 30, H: 30}

		result := intersectClipRect(outer, inner)

		assert.Equal(t, 20, result.X)
		assert.Equal(t, 20, result.Y)
		assert.Equal(t, 30, result.W)
		assert.Equal(t, 30, result.H)
	})

	t.Run("should return same rect when intersecting with itself", func(t *testing.T) {
		t.Parallel()

		rect := draw.ClipRect{X: 10, Y: 20, W: 30, H: 40}

		result := intersectClipRect(rect, rect)

		assert.Equal(t, rect, result)
	})

	t.Run("should handle adjacent rects (no overlap)", func(t *testing.T) {
		t.Parallel()

		a := draw.ClipRect{X: 0, Y: 0, W: 10, H: 10}
		b := draw.ClipRect{X: 10, Y: 0, W: 10, H: 10}

		result := intersectClipRect(a, b)

		assert.Equal(t, 0, result.W)
	})

	t.Run("should handle partial vertical overlap", func(t *testing.T) {
		t.Parallel()

		a := draw.ClipRect{X: 0, Y: 0, W: 20, H: 10}
		b := draw.ClipRect{X: 5, Y: 5, W: 20, H: 10}

		result := intersectClipRect(a, b)

		assert.Equal(t, 5, result.X)
		assert.Equal(t, 5, result.Y)
		assert.Equal(t, 15, result.W)
		assert.Equal(t, 5, result.H)
	})

	t.Run("should handle partial horizontal overlap", func(t *testing.T) {
		t.Parallel()

		a := draw.ClipRect{X: 0, Y: 0, W: 10, H: 20}
		b := draw.ClipRect{X: 5, Y: 5, W: 10, H: 20}

		result := intersectClipRect(a, b)

		assert.Equal(t, 5, result.X)
		assert.Equal(t, 5, result.Y)
		assert.Equal(t, 5, result.W)
		assert.Equal(t, 15, result.H)
	})

	t.Run("should handle zero-size input rect", func(t *testing.T) {
		t.Parallel()

		a := draw.ClipRect{X: 0, Y: 0, W: 0, H: 0}
		b := draw.ClipRect{X: 0, Y: 0, W: 10, H: 10}

		result := intersectClipRect(a, b)

		assert.Equal(t, 0, result.W)
		assert.Equal(t, 0, result.H)
	})

	t.Run("should handle rects at origin", func(t *testing.T) {
		t.Parallel()

		a := draw.ClipRect{X: 0, Y: 0, W: 50, H: 50}
		b := draw.ClipRect{X: 0, Y: 0, W: 30, H: 30}

		result := intersectClipRect(a, b)

		assert.Equal(t, 0, result.X)
		assert.Equal(t, 0, result.Y)
		assert.Equal(t, 30, result.W)
		assert.Equal(t, 30, result.H)
	})

	t.Run("should be commutative", func(t *testing.T) {
		t.Parallel()

		a := draw.ClipRect{X: 5, Y: 10, W: 20, H: 15}
		b := draw.ClipRect{X: 15, Y: 5, W: 20, H: 25}

		result1 := intersectClipRect(a, b)
		result2 := intersectClipRect(b, a)

		assert.Equal(t, result1, result2)
	})

	t.Run("should handle single pixel overlap", func(t *testing.T) {
		t.Parallel()

		a := draw.ClipRect{X: 0, Y: 0, W: 10, H: 10}
		b := draw.ClipRect{X: 9, Y: 9, W: 10, H: 10}

		result := intersectClipRect(a, b)

		assert.Equal(t, 9, result.X)
		assert.Equal(t, 9, result.Y)
		assert.Equal(t, 1, result.W)
		assert.Equal(t, 1, result.H)
	})
}

// Constructor function tests

func TestConstructorFunctions(t *testing.T) {
	t.Parallel()

	t.Run("should create box with NewBox", func(t *testing.T) {
		t.Parallel()

		box := NewBox(Style{Width: 20, Height: 10})

		require.NotNil(t, box)
		assert.Equal(t, 20, box.Style.Width)
		assert.Equal(t, 10, box.Style.Height)
	})

	t.Run("should create text element with Text", func(t *testing.T) {
		t.Parallel()

		text := Text("Hello World", Style{Foreground: Color("#FFFFFF")})

		require.NotNil(t, text)
		assert.Equal(t, "Hello World", text.Content)
		assert.Equal(t, Color("#FFFFFF"), text.Style.Foreground)
	})

	t.Run("should auto-calculate text width", func(t *testing.T) {
		t.Parallel()

		text := Text("Hello", Style{})

		assert.Equal(t, 5, text.Style.Width)
		assert.Equal(t, 1, text.Style.Height)
	})

	t.Run("should create buffer with NewBuffer", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(80, 24)

		require.NotNil(t, buf)
		assert.Equal(t, 80, buf.Width)
		assert.Equal(t, 24, buf.Height)
	})

	t.Run("should render buffer to string with Render", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(5, 1)

		result := Render(buf)

		assert.NotEmpty(t, result)
	})

	t.Run("should render buffer to writer with RenderTo", func(t *testing.T) {
		t.Parallel()

		buf := NewBuffer(5, 1)
		var w bytes.Buffer

		RenderTo(buf, &w)

		assert.NotEmpty(t, w.String())
	})
}

// Border helper functions tests

func TestBorderHelperFunctions(t *testing.T) {
	t.Parallel()

	t.Run("should create border with BorderAll", func(t *testing.T) {
		t.Parallel()

		border := BorderAll(BorderSingle)

		assert.True(t, border.Top.IsSet())
		assert.True(t, border.Right.IsSet())
		assert.True(t, border.Bottom.IsSet())
		assert.True(t, border.Left.IsSet())
		assert.Equal(t, style.BorderSingle, border.Top.Style)
	})

	t.Run("should create border with color using BorderAllWithColor", func(t *testing.T) {
		t.Parallel()

		border := BorderAllWithColor(BorderDouble, Color("#FF0000"))

		assert.Equal(t, style.BorderDouble, border.Top.Style)
		assert.Equal(t, Color("#FF0000"), border.Top.Color)
	})

	t.Run("should create border with title using BorderAllWithTitle", func(t *testing.T) {
		t.Parallel()

		border := BorderAllWithTitle(BorderRound, Color("#00FF00"), "My Title")

		assert.Equal(t, style.BorderRound, border.Top.Style)
		require.Len(t, border.Top.Texts, 1)
		assert.Equal(t, "My Title", border.Top.Texts[0].Text)
	})

	t.Run("should create border with BorderXY", func(t *testing.T) {
		t.Parallel()

		xEdge := style.BorderEdge{Style: style.BorderSingle}
		yEdge := style.BorderEdge{Style: style.BorderDouble}

		border := BorderXY(xEdge, yEdge)

		assert.Equal(t, style.BorderDouble, border.Top.Style)
		assert.Equal(t, style.BorderSingle, border.Right.Style)
		assert.Equal(t, style.BorderDouble, border.Bottom.Style)
		assert.Equal(t, style.BorderSingle, border.Left.Style)
	})

	t.Run("should create border with BorderTRBL", func(t *testing.T) {
		t.Parallel()

		top := style.BorderEdge{Style: style.BorderSingle}
		right := style.BorderEdge{Style: style.BorderDouble}
		bottom := style.BorderEdge{Style: style.BorderRound}
		left := style.BorderEdge{Style: style.BorderBold}

		border := BorderTRBL(top, right, bottom, left)

		assert.Equal(t, style.BorderSingle, border.Top.Style)
		assert.Equal(t, style.BorderDouble, border.Right.Style)
		assert.Equal(t, style.BorderRound, border.Bottom.Style)
		assert.Equal(t, style.BorderBold, border.Left.Style)
	})
}

// Integration tests

func TestIntegration(t *testing.T) {
	t.Parallel()

	t.Run("should render complete box tree to string", func(t *testing.T) {
		t.Parallel()

		container := NewBox(Style{
			Width:         40,
			Height:        10,
			Display:       Flex,
			FlexDirection: Row,
			Border:        BorderAll(BorderSingle),
			Background:    Color("#1a1a1a"),
		})

		child1 := Text("Hello", Style{
			Width:      10,
			Height:     1,
			Foreground: Color("#FFFFFF"),
		})
		child2 := Text("World", Style{
			Width:      10,
			Height:     1,
			Foreground: Color("#00FF00"),
		})

		container.AddChild(child1)
		container.AddChild(child2)

		buf := Draw(container)
		output := Render(buf)

		assert.NotEmpty(t, output)
		assert.Contains(t, output, "H")
		assert.Contains(t, output, "W")
	})

	t.Run("should handle deeply nested boxes", func(t *testing.T) {
		t.Parallel()

		root := NewBox(Style{Width: 50, Height: 20})
		level1 := NewBox(Style{Width: 40, Height: 15})
		level2 := NewBox(Style{Width: 30, Height: 10})
		level3 := Text("Deep", Style{})

		level2.AddChild(level3)
		level1.AddChild(level2)
		root.AddChild(level1)

		buf := Draw(root)

		assert.Equal(t, 50, buf.Width)
		assert.Equal(t, 20, buf.Height)
		assert.Contains(t, Render(buf), "D")
	})

	t.Run("should render box with z-index ordering", func(t *testing.T) {
		t.Parallel()

		container := NewBox(Style{Width: 20, Height: 10})
		bottom := NewBox(Style{
			Width:      10,
			Height:     5,
			ZIndex:     0,
			Position:   Absolute,
			Background: Color("#FF0000"),
		})
		top := NewBox(Style{
			Width:      10,
			Height:     5,
			ZIndex:     1,
			Position:   Absolute,
			Background: Color("#00FF00"),
		})

		container.AddChild(top)
		container.AddChild(bottom)

		buf := Draw(container)

		assert.Equal(t, Color("#00FF00"), buf.Cells[0][0].Background)
	})
}
