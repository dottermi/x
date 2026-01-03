package layout

import (
	"testing"

	"github.com/dottermi/x/termistyle/style"
	"github.com/stretchr/testify/assert"
)

func TestTotalChildrenWidth(t *testing.T) {
	t.Parallel()

	t.Run("should return zero when children slice is empty", func(t *testing.T) {
		t.Parallel()

		children := []*Box{}

		result := totalChildrenWidth(children)

		assert.Equal(t, 0, result)
	})

	t.Run("should return width when single child has no margins", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width: 20,
			},
		}
		children := []*Box{child}

		result := totalChildrenWidth(children)

		assert.Equal(t, 20, result)
	})

	t.Run("should include horizontal margins for single child", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width: 15,
				Margin: style.Spacing{
					Left:  3,
					Right: 5,
				},
			},
		}
		children := []*Box{child}

		result := totalChildrenWidth(children)

		assert.Equal(t, 23, result) // 15 + 3 + 5
	})

	t.Run("should sum widths for multiple children without margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width: 10,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width: 15,
			},
		}
		child3 := &Box{
			Style: style.Style{
				Width: 8,
			},
		}
		children := []*Box{child1, child2, child3}

		result := totalChildrenWidth(children)

		assert.Equal(t, 33, result) // 10 + 15 + 8
	})

	t.Run("should sum widths and margins for multiple children with margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width: 10,
				Margin: style.Spacing{
					Left:  2,
					Right: 3,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width: 15,
				Margin: style.Spacing{
					Left:  1,
					Right: 4,
				},
			},
		}
		children := []*Box{child1, child2}

		result := totalChildrenWidth(children)

		assert.Equal(t, 35, result) // (10 + 2 + 3) + (15 + 1 + 4)
	})

	t.Run("should ignore vertical margins in width calculation", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width: 20,
				Margin: style.Spacing{
					Top:    10,
					Bottom: 10,
					Left:   5,
					Right:  5,
				},
			},
		}
		children := []*Box{child}

		result := totalChildrenWidth(children)

		assert.Equal(t, 30, result) // 20 + 5 + 5 (vertical margins ignored)
	})

	t.Run("should handle mix of children with different margin configurations", func(t *testing.T) {
		t.Parallel()

		childNoMargin := &Box{
			Style: style.Style{
				Width: 10,
			},
		}
		childLeftOnly := &Box{
			Style: style.Style{
				Width: 10,
				Margin: style.Spacing{
					Left: 5,
				},
			},
		}
		childRightOnly := &Box{
			Style: style.Style{
				Width: 10,
				Margin: style.Spacing{
					Right: 3,
				},
			},
		}
		childBothMargins := &Box{
			Style: style.Style{
				Width:  10,
				Margin: style.SpacingAll(2),
			},
		}
		children := []*Box{childNoMargin, childLeftOnly, childRightOnly, childBothMargins}

		result := totalChildrenWidth(children)

		// 10 + (10 + 5) + (10 + 3) + (10 + 2 + 2) = 52
		assert.Equal(t, 52, result)
	})
}

func TestBox_calculateFlexRow(t *testing.T) {
	t.Parallel()

	t.Run("should position single child at correct X with left margin", func(t *testing.T) {
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
			W: 40,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child},
		}

		box.calculateFlexRow(0, 0, 40, 20)

		assert.Equal(t, 7, child.X) // 0 (parent X) + 0 (startX) + 7 (margin left)
	})

	t.Run("should position single child at correct Y with top margin", func(t *testing.T) {
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
			W: 40,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child},
		}

		box.calculateFlexRow(0, 0, 40, 20)

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
			W: 40,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child},
		}

		box.calculateFlexRow(2, 1, 36, 18)

		assert.Equal(t, 17, child.X) // 10 (parent X) + 2 (startX) + 5 (margin left)
		assert.Equal(t, 24, child.Y) // 20 (parent Y) + 1 (startY) + 3 (margin top)
	})

	t.Run("should position multiple children correctly with margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  2,
					Right: 3,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  8,
				Height: 4,
				Margin: style.Spacing{
					Left:  1,
					Right: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexRow(0, 0, 50, 20)

		assert.Equal(t, 2, child1.X) // 0 + 2 (margin left)
		// child2.X = 0 + (10 + 2 + 3) + 1 = 16 (x advanced by child1 outer width, then margin left)
		assert.Equal(t, 16, child2.X)
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
				Width:  8,
				Height: 4,
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				Gap:           5,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexRow(0, 0, 50, 20)

		assert.Equal(t, 0, child1.X)
		assert.Equal(t, 15, child2.X) // 10 (child1 width) + 5 (gap)
	})

	t.Run("should combine margin and gap correctly", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  2,
					Right: 3,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  8,
				Height: 4,
				Margin: style.Spacing{
					Left:  1,
					Right: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 60,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				Gap:           4,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexRow(0, 0, 60, 20)

		assert.Equal(t, 2, child1.X) // 0 + 2 (margin left)
		// child2.X = 0 + (10 + 2 + 3) + 4 (gap) + 1 (margin left) = 20
		assert.Equal(t, 20, child2.X)
	})

	t.Run("should stretch height and account for vertical margins", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width: 10,
				// Height: 0 triggers stretch
				Margin: style.Spacing{
					Top:    3,
					Bottom: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child},
		}

		box.calculateFlexRow(0, 0, 40, 20)

		// Stretched height = innerH - vertical margins = 20 - 3 - 2 = 15
		assert.Equal(t, 15, child.Style.Height)
		assert.Equal(t, 3, child.Y) // top margin applied
	})

	t.Run("should handle FlexGrow with remaining space after margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Left:  2,
					Right: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Left:  3,
					Right: 3,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 60,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexRow(0, 0, 60, 20)

		// Total margins = 2 + 2 + 3 + 3 = 10
		// Available for content = 60 - 10 = 50
		// Base width = 10 + 10 = 20
		// Extra space = 50 - 20 = 30, split equally = 15 each
		assert.Equal(t, 25, child1.Style.Width) // 10 + 15
		assert.Equal(t, 25, child2.Style.Width) // 10 + 15
	})

	t.Run("should position with JustifyCenter and margins", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  2,
					Right: 3,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:        style.Flex,
				FlexDirection:  style.Row,
				JustifyContent: style.JustifyCenter,
			},
			Children: []*Box{child},
		}

		box.calculateFlexRow(0, 0, 50, 20)

		// Total width including margins = 10 + 2 + 3 = 15
		// Center offset = (50 - 15) / 2 = 17
		// child.X = 0 + 17 + 2 (margin left) = 19
		assert.Equal(t, 19, child.X)
	})

	t.Run("should position with JustifyEnd and margins", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  2,
					Right: 3,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:        style.Flex,
				FlexDirection:  style.Row,
				JustifyContent: style.JustifyEnd,
			},
			Children: []*Box{child},
		}

		box.calculateFlexRow(0, 0, 50, 20)

		// Total width including margins = 10 + 2 + 3 = 15
		// End offset = 50 - 15 = 35
		// child.X = 0 + 35 + 2 (margin left) = 37
		assert.Equal(t, 37, child.X)
	})

	t.Run("should position with JustifyBetween and margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  1,
					Right: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  2,
					Right: 1,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:        style.Flex,
				FlexDirection:  style.Row,
				JustifyContent: style.JustifyBetween,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexRow(0, 0, 50, 20)

		// Total children width = (10 + 1 + 2) + (10 + 2 + 1) = 26
		// Spacing = (50 - 26) / (2 - 1) = 24
		assert.Equal(t, 1, child1.X) // 0 + 1 (margin left)
		// child2.X = 0 + 13 + 24 (spacing) + 2 (margin left) = 39
		assert.Equal(t, 39, child2.X)
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
				Width:  8,
				Height: 4,
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child1, absoluteChild, child2},
		}

		box.calculateFlexRow(0, 0, 50, 20)

		assert.Equal(t, 0, child1.X)
		// child2 should be positioned as if absoluteChild doesn't exist
		assert.Equal(t, 10, child2.X) // 10 (child1 width)
	})

	t.Run("should return early when no children exist", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{},
		}

		box.calculateFlexRow(0, 0, 40, 20)

		assert.Empty(t, box.Children)
	})

	t.Run("should handle child with uniform margins using SpacingAll", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.SpacingAll(3),
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child},
		}

		box.calculateFlexRow(0, 0, 40, 20)

		assert.Equal(t, 3, child.X) // margin left
		assert.Equal(t, 3, child.Y) // margin top
	})

	t.Run("should clamp stretched height to zero when margins exceed inner height", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width: 10,
				// Height: 0 triggers stretch
				Margin: style.Spacing{
					Top:    15,
					Bottom: 10,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child},
		}

		box.calculateFlexRow(0, 0, 40, 20)

		// Stretched height = 20 - 15 - 10 = -5, clamped to 0
		assert.Equal(t, 0, child.Style.Height)
	})
}

func TestSplitIntoRowLines(t *testing.T) {
	t.Parallel()

	t.Run("should return single line when single child fits in width", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
			},
		}
		children := []*Box{child}

		lines := splitIntoRowLines(children, 50, 0)

		assert.Len(t, lines, 1)
		assert.Len(t, lines[0], 1)
		assert.Same(t, child, lines[0][0])
	})

	t.Run("should return single line when multiple children fit in width", func(t *testing.T) {
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
				Height: 5,
			},
		}
		child3 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
			},
		}
		children := []*Box{child1, child2, child3}

		lines := splitIntoRowLines(children, 50, 0)

		// 15 + 20 + 10 = 45, fits in 50
		assert.Len(t, lines, 1)
		assert.Len(t, lines[0], 3)
	})

	t.Run("should wrap children to multiple lines due to width", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
			},
		}
		child3 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
			},
		}
		children := []*Box{child1, child2, child3}

		lines := splitIntoRowLines(children, 45, 0)

		// 20 + 20 = 40 fits, then 20 wraps to new line
		assert.Len(t, lines, 2)
		assert.Len(t, lines[0], 2)
		assert.Len(t, lines[1], 1)
		assert.Same(t, child1, lines[0][0])
		assert.Same(t, child2, lines[0][1])
		assert.Same(t, child3, lines[1][0])
	})

	t.Run("should wrap children to multiple lines due to margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  15,
				Height: 5,
				Margin: style.Spacing{
					Left:  5,
					Right: 5,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  15,
				Height: 5,
				Margin: style.Spacing{
					Left:  5,
					Right: 5,
				},
			},
		}
		children := []*Box{child1, child2}

		// Without margins: 15 + 15 = 30, fits in 50
		// With margins: (15 + 5 + 5) + (15 + 5 + 5) = 50, exactly fits
		// Let's use a smaller width to trigger wrap
		lines := splitIntoRowLines(children, 40, 0)

		// child1 outer width = 25, child2 outer width = 25, total = 50 > 40
		assert.Len(t, lines, 2)
		assert.Len(t, lines[0], 1)
		assert.Len(t, lines[1], 1)
		assert.Same(t, child1, lines[0][0])
		assert.Same(t, child2, lines[1][0])
	})

	t.Run("should consider gap in line breaking", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
			},
		}
		children := []*Box{child1, child2}

		// Without gap: 20 + 20 = 40, fits in 45
		// With gap of 10: 20 + 10 + 20 = 50 > 45, should wrap
		lines := splitIntoRowLines(children, 45, 10)

		assert.Len(t, lines, 2)
		assert.Len(t, lines[0], 1)
		assert.Len(t, lines[1], 1)
	})

	t.Run("should return empty slice when children is empty", func(t *testing.T) {
		t.Parallel()

		children := []*Box{}

		lines := splitIntoRowLines(children, 50, 0)

		assert.Empty(t, lines)
	})

	t.Run("should handle child with zero width by using minimum of 1", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  0,
				Height: 5,
			},
		}
		children := []*Box{child}

		lines := splitIntoRowLines(children, 50, 0)

		assert.Len(t, lines, 1)
		assert.Len(t, lines[0], 1)
	})

	t.Run("should place oversized single child on its own line", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  100,
				Height: 5,
			},
		}
		children := []*Box{child}

		lines := splitIntoRowLines(children, 50, 0)

		// Even though child is wider than container, it gets its own line
		assert.Len(t, lines, 1)
		assert.Len(t, lines[0], 1)
	})

	t.Run("should combine margins and gap when calculating line breaks", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  2,
					Right: 3,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  3,
					Right: 2,
				},
			},
		}
		children := []*Box{child1, child2}

		// child1 outer = 10 + 2 + 3 = 15
		// child2 outer = 10 + 3 + 2 = 15
		// gap = 5
		// total = 15 + 5 + 15 = 35, fits in 40
		lines := splitIntoRowLines(children, 40, 5)

		assert.Len(t, lines, 1)
		assert.Len(t, lines[0], 2)
	})

	t.Run("should wrap when margins and gap exceed container width", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  2,
					Right: 3,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  3,
					Right: 2,
				},
			},
		}
		children := []*Box{child1, child2}

		// child1 outer = 15, child2 outer = 15, gap = 5, total = 35 > 30
		lines := splitIntoRowLines(children, 30, 5)

		assert.Len(t, lines, 2)
		assert.Len(t, lines[0], 1)
		assert.Len(t, lines[1], 1)
	})
}

func TestBox_calculateFlexRowWrap(t *testing.T) {
	t.Parallel()

	t.Run("should position single line children with margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left: 3,
					Top:  2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left: 4,
					Top:  1,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 60,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexRowWrap(0, 0, 60, 20)

		assert.Equal(t, 3, child1.X)  // 0 + 0 + 3 (margin left)
		assert.Equal(t, 2, child1.Y)  // 0 + 0 + 2 (margin top)
		assert.Equal(t, 17, child2.X) // 0 + (10 + 3 + 0) + 4 (margin left) = 17
		assert.Equal(t, 1, child2.Y)  // 0 + 0 + 1 (margin top)
	})

	t.Run("should position multiple lines with margins correctly", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
				Margin: style.Spacing{
					Left: 2,
					Top:  1,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 6,
				Margin: style.Spacing{
					Left:   3,
					Top:    2,
					Bottom: 2,
				},
			},
		}
		child3 := &Box{
			Style: style.Style{
				Width:  15,
				Height: 4,
				Margin: style.Spacing{
					Left: 1,
					Top:  3,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child1, child2, child3},
		}

		// child1 outer = 22, child2 outer = 23, total = 45 < 50
		// child3 outer = 16, adding to line: 45 + 16 = 61 > 50, wraps
		box.calculateFlexRowWrap(0, 0, 50, 30)

		// First line: child1 and child2
		assert.Equal(t, 2, child1.X)  // margin left
		assert.Equal(t, 1, child1.Y)  // margin top
		assert.Equal(t, 25, child2.X) // 22 (child1 outer) + 3 (margin left)
		assert.Equal(t, 2, child2.Y)  // margin top

		// Second line: child3
		// Line height of first line = max(5+1+0, 6+2+2) = max(6, 10) = 10
		assert.Equal(t, 1, child3.X)  // margin left
		assert.Equal(t, 13, child3.Y) // 10 (line height) + 0 (gap) + 3 (margin top)
	})

	t.Run("should calculate line heights including vertical margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
				Margin: style.Spacing{
					Top:    3,
					Bottom: 4,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 8,
				Margin: style.Spacing{
					Top:    1,
					Bottom: 1,
				},
			},
		}
		child3 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 45,
			H: 40,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child1, child2, child3},
		}

		// child1 outer width = 20, child2 outer width = 20, total = 40 < 45
		// child3 outer width = 20, adding: 40 + 20 = 60 > 45, wraps
		box.calculateFlexRowWrap(0, 0, 45, 40)

		// First line outer heights: child1 = 5+3+4 = 12, child2 = 8+1+1 = 10
		// Line height = max(12, 10) = 12
		// Second line starts at Y = 12 (first line height)
		assert.Equal(t, 12, child3.Y) // second line Y position
	})

	t.Run("should cause earlier line break when margins are added", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
			},
		}
		childWithMargin := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
				Margin: style.Spacing{
					Left:  5,
					Right: 5,
				},
			},
		}
		boxWithoutMargins := &Box{
			X: 0,
			Y: 0,
			W: 70,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child1, child2},
		}
		boxWithMargins := &Box{
			X: 0,
			Y: 0,
			W: 70,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{
				{
					Style: style.Style{
						Width:  20,
						Height: 5,
						Margin: style.Spacing{
							Left:  5,
							Right: 5,
						},
					},
				},
				{
					Style: style.Style{
						Width:  20,
						Height: 5,
						Margin: style.Spacing{
							Left:  5,
							Right: 5,
						},
					},
				},
				childWithMargin,
			},
		}

		boxWithoutMargins.calculateFlexRowWrap(0, 0, 70, 20)
		boxWithMargins.calculateFlexRowWrap(0, 0, 70, 20)

		// Without margins: 20 + 20 = 40 < 70, fits in one line
		assert.Equal(t, 0, child1.Y)
		assert.Equal(t, 0, child2.Y)

		// With margins: (20+10) + (20+10) = 60 < 70, first two fit
		// Third child: 60 + (20+10) = 90 > 70, wraps
		assert.Equal(t, 5, childWithMargin.Y) // wrapped to second line
	})

	t.Run("should apply JustifyCenter across wrapped lines", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  2,
					Right: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  2,
					Right: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:        style.Flex,
				FlexDirection:  style.Row,
				FlexWrap:       style.Wrap,
				JustifyContent: style.JustifyCenter,
			},
			Children: []*Box{child1, child2},
		}

		// child1 outer = 14, child2 outer = 14, total = 28
		// Center offset = (50 - 28) / 2 = 11
		box.calculateFlexRowWrap(0, 0, 50, 20)

		assert.Equal(t, 13, child1.X) // 11 + 2 (margin left)
		assert.Equal(t, 27, child2.X) // 11 + 14 + 2 (margin left)
	})

	t.Run("should apply JustifyEnd across wrapped lines", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  1,
					Right: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  2,
					Right: 1,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:        style.Flex,
				FlexDirection:  style.Row,
				FlexWrap:       style.Wrap,
				JustifyContent: style.JustifyEnd,
			},
			Children: []*Box{child1, child2},
		}

		// child1 outer = 13, child2 outer = 13, total = 26
		// End offset = 50 - 26 = 24
		box.calculateFlexRowWrap(0, 0, 50, 20)

		assert.Equal(t, 25, child1.X) // 24 + 1 (margin left)
		assert.Equal(t, 39, child2.X) // 24 + 13 + 2 (margin left)
	})

	t.Run("should apply JustifyBetween across wrapped lines", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  1,
					Right: 1,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  1,
					Right: 1,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:        style.Flex,
				FlexDirection:  style.Row,
				FlexWrap:       style.Wrap,
				JustifyContent: style.JustifyBetween,
			},
			Children: []*Box{child1, child2},
		}

		// child1 outer = 12, child2 outer = 12, total = 24
		// Spacing = (50 - 24) / (2 - 1) = 26
		box.calculateFlexRowWrap(0, 0, 50, 20)

		assert.Equal(t, 1, child1.X)  // 0 + 1 (margin left)
		assert.Equal(t, 39, child2.X) // 0 + 12 + 26 + 1 (margin left)
	})

	t.Run("should return early when no children exist", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{},
		}

		box.calculateFlexRowWrap(0, 0, 50, 20)

		assert.Empty(t, box.Children)
	})

	t.Run("should apply gap between lines", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  30,
				Height: 5,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  30,
				Height: 5,
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				FlexWrap:      style.Wrap,
				Gap:           3,
			},
			Children: []*Box{child1, child2},
		}

		// child1 = 30, child2 = 30, gap = 3
		// 30 + 3 + 30 = 63 > 40, wraps
		box.calculateFlexRowWrap(0, 0, 40, 30)

		assert.Equal(t, 0, child1.Y)
		assert.Equal(t, 8, child2.Y) // 5 (line height) + 3 (gap)
	})

	t.Run("should handle parent offset X and Y correctly", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left: 2,
					Top:  3,
				},
			},
		}
		box := &Box{
			X: 10,
			Y: 20,
			W: 50,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child},
		}

		box.calculateFlexRowWrap(5, 2, 40, 26)

		assert.Equal(t, 17, child.X) // 10 (parent X) + 5 (startX) + 2 (margin left)
		assert.Equal(t, 25, child.Y) // 20 (parent Y) + 2 (startY) + 3 (margin top)
	})

	t.Run("should handle child height stretch with vertical margins", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  20,
				Height: 0, // triggers stretch behavior based on line height
				Margin: style.Spacing{
					Top:    2,
					Bottom: 3,
				},
			},
		}
		sibling := &Box{
			Style: style.Style{
				Width:  20,
				Height: 10,
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child, sibling},
		}

		// Both fit on one line, line height = max(0+2+3, 10+0+0) = max(5, 10) = 10
		box.calculateFlexRowWrap(0, 0, 50, 30)

		// Child height is adjusted based on line height and margins
		assert.Equal(t, 2, child.Y) // margin top
	})

	t.Run("should skip absolute positioned children in wrapping", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
			},
		}
		absoluteChild := &Box{
			Style: style.Style{
				Position: style.Absolute,
				Width:    100,
				Height:   100,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child1, absoluteChild, child2},
		}

		box.calculateFlexRowWrap(0, 0, 50, 20)

		// Absolute child should not affect layout of normal flow children
		assert.Equal(t, 0, child1.X)
		assert.Equal(t, 20, child2.X) // positioned right after child1
		assert.Equal(t, 0, child1.Y)
		assert.Equal(t, 0, child2.Y) // same line
	})
}

func TestBox_distributeFlexSpaceRow(t *testing.T) {
	t.Parallel()

	t.Run("should reduce available grow space by horizontal margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Left:  5,
					Right: 5,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Left:  5,
					Right: 5,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 80,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child1, child2},
		}

		box.distributeFlexSpaceRow([]*Box{child1, child2}, 80)

		// Total margins = (5 + 5) + (5 + 5) = 20
		// Available for content = 80 - 20 = 60
		// Base width = 10 + 10 = 20
		// Extra space = 60 - 20 = 40, split equally = 20 each
		assert.Equal(t, 30, child1.Style.Width) // 10 + 20
		assert.Equal(t, 30, child2.Style.Width) // 10 + 20
	})

	t.Run("should calculate shrink correctly with margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  30,
				Height: 5,
				Margin: style.Spacing{
					Left:  5,
					Right: 5,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  30,
				Height: 5,
				Margin: style.Spacing{
					Left:  5,
					Right: 5,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 60,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child1, child2},
		}

		box.distributeFlexSpaceRow([]*Box{child1, child2}, 60)

		// Total margins = 20
		// Available for content = 60 - 20 = 40
		// Base width = 30 + 30 = 60
		// Overflow = 60 - 40 = 20, shrink proportionally
		// Default shrink = 1, each child shrinks by 10
		assert.Equal(t, 20, child1.Style.Width) // 30 - 10
		assert.Equal(t, 20, child2.Style.Width) // 30 - 10
	})

	t.Run("should not distribute when no grow and space is available", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				// FlexGrow: 0 (default)
				Margin: style.Spacing{
					Left:  2,
					Right: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child},
		}

		box.distributeFlexSpaceRow([]*Box{child}, 50)

		assert.Equal(t, 10, child.Style.Width) // unchanged
	})

	t.Run("should handle gap in space calculation", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Left:  2,
					Right: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Left:  2,
					Right: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 50,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
				Gap:           4,
			},
			Children: []*Box{child1, child2},
		}

		box.distributeFlexSpaceRow([]*Box{child1, child2}, 50)

		// Total margins = (2 + 2) + (2 + 2) = 8
		// Gap = 4
		// Available for content = 50 - 8 = 42
		// Base width = 10 + 10 = 20
		// Gap subtracted from extra: 42 - 20 - 4 = 18, split equally = 9 each
		assert.Equal(t, 19, child1.Style.Width) // 10 + 9
		assert.Equal(t, 19, child2.Style.Width) // 10 + 9
	})

	t.Run("should distribute grow proportionally based on FlexGrow values", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Left:  2,
					Right: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 2, // grows twice as much
				Margin: style.Spacing{
					Left:  2,
					Right: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 62,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{child1, child2},
		}

		box.distributeFlexSpaceRow([]*Box{child1, child2}, 62)

		// Total margins = 8
		// Available = 62 - 8 = 54
		// Base width = 20
		// Extra = 54 - 20 = 34
		// child1 gets 34 * (1/3) = 11, child2 gets 34 * (2/3) = 22 (+1 remainder)
		assert.Equal(t, 21, child1.Style.Width) // 10 + 11
		assert.Equal(t, 33, child2.Style.Width) // 10 + 23 (22 + 1 remainder)
	})

	t.Run("should handle mix of grow and non-grow children with margins", func(t *testing.T) {
		t.Parallel()

		fixedChild := &Box{
			Style: style.Style{
				Width:  20,
				Height: 5,
				// FlexGrow: 0 (default)
				Margin: style.Spacing{
					Left:  3,
					Right: 3,
				},
			},
		}
		growChild := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Left:  2,
					Right: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 60,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Row,
			},
			Children: []*Box{fixedChild, growChild},
		}

		box.distributeFlexSpaceRow([]*Box{fixedChild, growChild}, 60)

		assert.Equal(t, 20, fixedChild.Style.Width) // unchanged
		// Total margins = 6 + 4 = 10
		// Available = 60 - 10 = 50
		// Fixed child takes 20, grow child base = 10, total base = 30
		// Extra = 50 - 30 = 20, all goes to growChild
		assert.Equal(t, 30, growChild.Style.Width) // 10 + 20
	})
}
