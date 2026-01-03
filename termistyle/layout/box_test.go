package layout

import (
	"testing"

	"github.com/dottermi/x/termistyle/style"
	"github.com/stretchr/testify/assert"
)

func TestBox_outerWidth(t *testing.T) {
	t.Parallel()

	t.Run("should return width when margin is zero", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			W:     10,
			Style: style.Style{},
		}

		result := box.outerWidth()

		assert.Equal(t, 10, result)
	})

	t.Run("should return width plus margins when using uniform margin", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			W: 10,
			Style: style.Style{
				Margin: style.SpacingAll(5),
			},
		}

		result := box.outerWidth()

		assert.Equal(t, 20, result) // 10 + 5 (left) + 5 (right)
	})

	t.Run("should return width plus left and right margins when using asymmetric margins", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			W: 10,
			Style: style.Style{
				Margin: style.Spacing{
					Top:    2,
					Right:  3,
					Bottom: 4,
					Left:   7,
				},
			},
		}

		result := box.outerWidth()

		assert.Equal(t, 20, result) // 10 + 7 (left) + 3 (right)
	})

	t.Run("should return only width when top and bottom margins are set but left and right are zero", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			W: 15,
			Style: style.Style{
				Margin: style.Spacing{
					Top:    10,
					Right:  0,
					Bottom: 10,
					Left:   0,
				},
			},
		}

		result := box.outerWidth()

		assert.Equal(t, 15, result)
	})
}

func TestBox_outerHeight(t *testing.T) {
	t.Parallel()

	t.Run("should return height when margin is zero", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			H:     8,
			Style: style.Style{},
		}

		result := box.outerHeight()

		assert.Equal(t, 8, result)
	})

	t.Run("should return height plus margins when using uniform margin", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			H: 8,
			Style: style.Style{
				Margin: style.SpacingAll(3),
			},
		}

		result := box.outerHeight()

		assert.Equal(t, 14, result) // 8 + 3 (top) + 3 (bottom)
	})

	t.Run("should return height plus top and bottom margins when using asymmetric margins", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			H: 8,
			Style: style.Style{
				Margin: style.Spacing{
					Top:    5,
					Right:  1,
					Bottom: 2,
					Left:   1,
				},
			},
		}

		result := box.outerHeight()

		assert.Equal(t, 15, result) // 8 + 5 (top) + 2 (bottom)
	})

	t.Run("should return only height when left and right margins are set but top and bottom are zero", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			H: 12,
			Style: style.Style{
				Margin: style.Spacing{
					Top:    0,
					Right:  10,
					Bottom: 0,
					Left:   10,
				},
			},
		}

		result := box.outerHeight()

		assert.Equal(t, 12, result)
	})
}

func TestBox_marginLeft(t *testing.T) {
	t.Parallel()

	t.Run("should return zero when margin is not set", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			Style: style.Style{},
		}

		result := box.marginLeft()

		assert.Equal(t, 0, result)
	})

	t.Run("should return uniform value when using SpacingAll", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			Style: style.Style{
				Margin: style.SpacingAll(4),
			},
		}

		result := box.marginLeft()

		assert.Equal(t, 4, result)
	})

	t.Run("should return left margin value when using asymmetric margins", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			Style: style.Style{
				Margin: style.Spacing{
					Top:    1,
					Right:  2,
					Bottom: 3,
					Left:   9,
				},
			},
		}

		result := box.marginLeft()

		assert.Equal(t, 9, result)
	})

	t.Run("should return zero when only other margins are set", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			Style: style.Style{
				Margin: style.Spacing{
					Top:    5,
					Right:  5,
					Bottom: 5,
					Left:   0,
				},
			},
		}

		result := box.marginLeft()

		assert.Equal(t, 0, result)
	})
}

func TestBox_marginTop(t *testing.T) {
	t.Parallel()

	t.Run("should return zero when margin is not set", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			Style: style.Style{},
		}

		result := box.marginTop()

		assert.Equal(t, 0, result)
	})

	t.Run("should return uniform value when using SpacingAll", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			Style: style.Style{
				Margin: style.SpacingAll(6),
			},
		}

		result := box.marginTop()

		assert.Equal(t, 6, result)
	})

	t.Run("should return top margin value when using asymmetric margins", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			Style: style.Style{
				Margin: style.Spacing{
					Top:    11,
					Right:  2,
					Bottom: 3,
					Left:   4,
				},
			},
		}

		result := box.marginTop()

		assert.Equal(t, 11, result)
	})

	t.Run("should return zero when only other margins are set", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			Style: style.Style{
				Margin: style.Spacing{
					Top:    0,
					Right:  5,
					Bottom: 5,
					Left:   5,
				},
			},
		}

		result := box.marginTop()

		assert.Equal(t, 0, result)
	})
}

func TestBox_calculateAutoSize(t *testing.T) {
	t.Parallel()

	t.Run("should return zero size when box has no children and no content", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			Style: style.Style{},
		}

		w, h := box.calculateAutoSize()

		assert.Equal(t, 0, w)
		assert.Equal(t, 0, h)
	})

	t.Run("should return content size when box has text content", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			Style:   style.Style{},
			Content: "Hello",
		}

		w, h := box.calculateAutoSize()

		assert.Equal(t, 5, w) // len("Hello")
		assert.Equal(t, 1, h)
	})

	t.Run("should return child dimensions when single child has no margin", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  20,
				Height: 10,
			},
		}
		box := &Box{
			Style:    style.Style{},
			Children: []*Box{child},
		}

		w, h := box.calculateAutoSize()

		assert.Equal(t, 20, w)
		assert.Equal(t, 10, h)
	})

	t.Run("should include margins when single child has margin", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  20,
				Height: 10,
				Margin: style.SpacingAll(5),
			},
		}
		box := &Box{
			Style:    style.Style{},
			Children: []*Box{child},
		}

		w, h := box.calculateAutoSize()

		assert.Equal(t, 30, w) // 20 + 5 (left) + 5 (right)
		assert.Equal(t, 20, h) // 10 + 5 (top) + 5 (bottom)
	})

	t.Run("should sum heights for multiple children in block display", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  15,
				Height: 5,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 8,
			},
		}
		box := &Box{
			Style: style.Style{
				Display: style.Block,
			},
			Children: []*Box{child1, child2},
		}

		w, h := box.calculateAutoSize()

		assert.Equal(t, 20, w) // max width
		assert.Equal(t, 13, h) // 5 + 8 sum heights
	})

	t.Run("should sum heights for multiple children in flex column display", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  15,
				Height: 5,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 8,
			},
		}
		box := &Box{
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child1, child2},
		}

		w, h := box.calculateAutoSize()

		assert.Equal(t, 20, w) // max width
		assert.Equal(t, 13, h) // 5 + 8 sum heights
	})

	t.Run("should sum widths for multiple children in flex row display", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  15,
				Height: 5,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 8,
			},
		}
		box := &Box{
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child1, child2},
		}

		w, h := box.calculateAutoSize()

		assert.Equal(t, 35, w) // 15 + 20 sum widths
		assert.Equal(t, 8, h)  // max height
	})

	t.Run("should include margins for multiple children with different margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    2,
					Right:  3,
					Bottom: 2,
					Left:   1,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  8,
				Height: 4,
				Margin: style.Spacing{
					Top:    1,
					Right:  1,
					Bottom: 3,
					Left:   4,
				},
			},
		}
		box := &Box{
			Style: style.Style{
				Display: style.Block,
			},
			Children: []*Box{child1, child2},
		}

		w, h := box.calculateAutoSize()

		// child1 outer: 10 + 1 + 3 = 14, child2 outer: 8 + 4 + 1 = 13, max = 14
		assert.Equal(t, 14, w)
		// child1 outer: 5 + 2 + 2 = 9, child2 outer: 4 + 1 + 3 = 8, sum = 17
		assert.Equal(t, 17, h)
	})

	t.Run("should include gap in height calculation for block display", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
			},
		}
		box := &Box{
			Style: style.Style{
				Display: style.Block,
				Gap:     3,
			},
			Children: []*Box{child1, child2},
		}

		w, h := box.calculateAutoSize()

		assert.Equal(t, 10, w)
		assert.Equal(t, 13, h) // 5 + 3 (gap) + 5
	})

	t.Run("should include gap in width calculation for flex row display", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
			},
		}
		box := &Box{
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				Gap:           4,
			},
			Children: []*Box{child1, child2},
		}

		w, h := box.calculateAutoSize()

		assert.Equal(t, 24, w) // 10 + 4 (gap) + 10
		assert.Equal(t, 5, h)
	})

	t.Run("should combine margins and gap for multiple children", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.SpacingAll(2),
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.SpacingAll(2),
			},
		}
		box := &Box{
			Style: style.Style{
				Display: style.Block,
				Gap:     3,
			},
			Children: []*Box{child1, child2},
		}

		w, h := box.calculateAutoSize()

		// child outer width: 10 + 2 + 2 = 14
		assert.Equal(t, 14, w)
		// child1 outer: 5 + 2 + 2 = 9, child2 outer: 9, gap: 3, total: 9 + 3 + 9 = 21
		assert.Equal(t, 21, h)
	})
}

func TestBox_calculateBlock(t *testing.T) {
	t.Parallel()

	t.Run("should position child at correct X with left margin", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left: 7,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 30,
			H: 20,
			Style: style.Style{
				Display: style.Block,
			},
			Children: []*Box{child},
		}

		box.calculateBlock(0, 0)

		assert.Equal(t, 7, child.X) // 0 (parent X) + 0 (startX) + 7 (margin left)
		assert.Equal(t, 0, child.Y)
	})

	t.Run("should position child at correct Y with top margin", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top: 4,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 30,
			H: 20,
			Style: style.Style{
				Display: style.Block,
			},
			Children: []*Box{child},
		}

		box.calculateBlock(0, 0)

		assert.Equal(t, 0, child.X)
		assert.Equal(t, 4, child.Y) // 0 (parent Y) + 0 (startY) + 4 (margin top)
	})

	t.Run("should position child with both left and top margins", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:  3,
					Left: 5,
				},
			},
		}
		box := &Box{
			X: 10,
			Y: 20,
			W: 30,
			H: 20,
			Style: style.Style{
				Display: style.Block,
			},
			Children: []*Box{child},
		}

		box.calculateBlock(2, 1)

		assert.Equal(t, 17, child.X) // 10 (parent X) + 2 (startX) + 5 (margin left)
		assert.Equal(t, 24, child.Y) // 20 (parent Y) + 1 (startY) + 3 (margin top)
	})

	t.Run("should stack multiple children vertically respecting margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 3,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 4,
				Margin: style.Spacing{
					Top:    1,
					Bottom: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 30,
			H: 30,
			Style: style.Style{
				Display: style.Block,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateBlock(0, 0)

		assert.Equal(t, 2, child1.Y) // 0 + 2 (margin top)
		// child2.Y = 0 + (5 + 2 + 3) + 1 = 11 (y advanced by child1 outer height, then margin top)
		assert.Equal(t, 11, child2.Y)
	})

	t.Run("should apply gap between children", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 4,
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 30,
			H: 30,
			Style: style.Style{
				Display: style.Block,
				Gap:     3,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateBlock(0, 0)

		assert.Equal(t, 0, child1.Y)
		assert.Equal(t, 8, child2.Y) // 5 (child1 height) + 3 (gap)
	})

	t.Run("should combine margin and gap correctly", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    1,
					Bottom: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 4,
				Margin: style.Spacing{
					Top:    3,
					Bottom: 1,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 30,
			H: 30,
			Style: style.Style{
				Display: style.Block,
				Gap:     2,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateBlock(0, 0)

		assert.Equal(t, 1, child1.Y) // 0 + 1 (margin top)
		// child2.Y = y (after child1) + margin top
		// y after child1 = 0 + (5 + 1 + 2) + 2 = 10 (outer height + gap)
		// child2.Y = 10 + 3 (margin top) = 13
		assert.Equal(t, 13, child2.Y)
	})

	t.Run("should skip absolute positioned children in normal flow", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
			},
		}
		absoluteChild := &Box{
			Style: style.Style{
				Position: style.Absolute,
				Width:    20,
				Height:   10,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 4,
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 30,
			H: 30,
			Style: style.Style{
				Display: style.Block,
			},
			Children: []*Box{child1, absoluteChild, child2},
		}

		box.calculateBlock(0, 0)

		assert.Equal(t, 0, child1.Y)
		// child2 should be positioned as if absoluteChild doesn't exist
		assert.Equal(t, 5, child2.Y) // 5 (child1 height)
	})

	t.Run("should position children with uniform margins using SpacingAll", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.SpacingAll(3),
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 4,
				Margin: style.SpacingAll(2),
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 30,
			H: 30,
			Style: style.Style{
				Display: style.Block,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateBlock(0, 0)

		assert.Equal(t, 3, child1.X) // margin left
		assert.Equal(t, 3, child1.Y) // margin top
		assert.Equal(t, 2, child2.X) // margin left
		// child2.Y = 0 + (5 + 3 + 3) + 2 = 13 (y advanced by child1 outer height, then margin top)
		assert.Equal(t, 13, child2.Y)
	})
}
