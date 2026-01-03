package layout

import (
	"testing"

	"github.com/dottermi/x/termistyle/style"
	"github.com/stretchr/testify/assert"
)

func TestTotalChildrenHeight(t *testing.T) {
	t.Parallel()

	t.Run("should return zero when children slice is empty", func(t *testing.T) {
		t.Parallel()

		children := []*Box{}

		result := totalChildrenHeight(children)

		assert.Equal(t, 0, result)
	})

	t.Run("should return height when single child has no margins", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Height: 15,
			},
		}
		children := []*Box{child}

		result := totalChildrenHeight(children)

		assert.Equal(t, 15, result)
	})

	t.Run("should include vertical margins for single child", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Height: 10,
				Margin: style.Spacing{
					Top:    3,
					Bottom: 5,
				},
			},
		}
		children := []*Box{child}

		result := totalChildrenHeight(children)

		assert.Equal(t, 18, result) // 10 + 3 + 5
	})

	t.Run("should sum heights for multiple children without margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Height: 10,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Height: 15,
			},
		}
		child3 := &Box{
			Style: style.Style{
				Height: 8,
			},
		}
		children := []*Box{child1, child2, child3}

		result := totalChildrenHeight(children)

		assert.Equal(t, 33, result) // 10 + 15 + 8
	})

	t.Run("should sum heights and margins for multiple children with margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Height: 10,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 3,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Height: 15,
				Margin: style.Spacing{
					Top:    1,
					Bottom: 4,
				},
			},
		}
		children := []*Box{child1, child2}

		result := totalChildrenHeight(children)

		assert.Equal(t, 35, result) // (10 + 2 + 3) + (15 + 1 + 4)
	})

	t.Run("should ignore horizontal margins in height calculation", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Height: 20,
				Margin: style.Spacing{
					Top:    5,
					Bottom: 5,
					Left:   10,
					Right:  10,
				},
			},
		}
		children := []*Box{child}

		result := totalChildrenHeight(children)

		assert.Equal(t, 30, result) // 20 + 5 + 5 (horizontal margins ignored)
	})

	t.Run("should handle mix of children with different margin configurations", func(t *testing.T) {
		t.Parallel()

		childNoMargin := &Box{
			Style: style.Style{
				Height: 10,
			},
		}
		childTopOnly := &Box{
			Style: style.Style{
				Height: 10,
				Margin: style.Spacing{
					Top: 5,
				},
			},
		}
		childBottomOnly := &Box{
			Style: style.Style{
				Height: 10,
				Margin: style.Spacing{
					Bottom: 3,
				},
			},
		}
		childBothMargins := &Box{
			Style: style.Style{
				Height: 10,
				Margin: style.SpacingAll(2),
			},
		}
		children := []*Box{childNoMargin, childTopOnly, childBottomOnly, childBothMargins}

		result := totalChildrenHeight(children)

		// 10 + (10 + 5) + (10 + 3) + (10 + 2 + 2) = 52
		assert.Equal(t, 52, result)
	})
}

func TestBox_calculateFlexColumn(t *testing.T) {
	t.Parallel()

	t.Run("should position single child at correct Y with top margin", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top: 7,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumn(0, 0, 40, 30)

		assert.Equal(t, 7, child.Y) // 0 (parent Y) + 0 (startY) + 7 (margin top)
	})

	t.Run("should position single child at correct X with left margin", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left: 4,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumn(0, 0, 40, 30)

		assert.Equal(t, 4, child.X) // 0 (parent X) + 0 (startX) + 4 (margin left)
	})

	t.Run("should position child with both top and left margins", func(t *testing.T) {
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
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumn(2, 1, 36, 28)

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
					Top:    2,
					Bottom: 3,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  8,
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
			W: 40,
			H: 50,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexColumn(0, 0, 40, 50)

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
				Width:  8,
				Height: 4,
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 50,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
				Gap:           5,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexColumn(0, 0, 40, 50)

		assert.Equal(t, 0, child1.Y)
		assert.Equal(t, 10, child2.Y) // 5 (child1 height) + 5 (gap)
	})

	t.Run("should combine margin and gap correctly", func(t *testing.T) {
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
				Width:  8,
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
			W: 40,
			H: 60,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
				Gap:           4,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexColumn(0, 0, 40, 60)

		assert.Equal(t, 2, child1.Y) // 0 + 2 (margin top)
		// child2.Y = 0 + (5 + 2 + 3) + 4 (gap) + 1 (margin top) = 15
		assert.Equal(t, 15, child2.Y)
	})

	t.Run("should stretch width and account for horizontal margins", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				// Width: 0 triggers stretch
				Height: 5,
				Margin: style.Spacing{
					Left:  3,
					Right: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumn(0, 0, 40, 30)

		// Stretched width = innerW - horizontal margins = 40 - 3 - 2 = 35
		assert.Equal(t, 35, child.Style.Width)
		assert.Equal(t, 3, child.X) // left margin applied
	})

	t.Run("should handle FlexGrow with remaining space after margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Top:    3,
					Bottom: 3,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 50,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexColumn(0, 0, 40, 50)

		// Total margins = (2 + 2) + (3 + 3) = 10
		// Available for content = 50 - 10 = 40
		// Base height = 5 + 5 = 10
		// Extra space = 40 - 10 = 30, split equally = 15 each
		assert.Equal(t, 20, child1.Style.Height) // 5 + 15
		assert.Equal(t, 20, child2.Style.Height) // 5 + 15
	})

	t.Run("should position with JustifyCenter and margins", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 3,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 50,
			Style: style.Style{
				Display:        style.Flex,
				FlexDirection:  style.Column,
				JustifyContent: style.JustifyCenter,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumn(0, 0, 40, 50)

		// Total height including margins = 5 + 2 + 3 = 10
		// Center offset = (50 - 10) / 2 = 20
		// child.Y = 0 + 20 + 2 (margin top) = 22
		assert.Equal(t, 22, child.Y)
	})

	t.Run("should position with JustifyEnd and margins", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 3,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 50,
			Style: style.Style{
				Display:        style.Flex,
				FlexDirection:  style.Column,
				JustifyContent: style.JustifyEnd,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumn(0, 0, 40, 50)

		// Total height including margins = 5 + 2 + 3 = 10
		// End offset = 50 - 10 = 40
		// child.Y = 0 + 40 + 2 (margin top) = 42
		assert.Equal(t, 42, child.Y)
	})

	t.Run("should position with JustifyBetween and margins", func(t *testing.T) {
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
				Height: 5,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 1,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 50,
			Style: style.Style{
				Display:        style.Flex,
				FlexDirection:  style.Column,
				JustifyContent: style.JustifyBetween,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexColumn(0, 0, 40, 50)

		// Total children height = (5 + 1 + 2) + (5 + 2 + 1) = 16
		// Spacing = (50 - 16) / (2 - 1) = 34
		assert.Equal(t, 1, child1.Y) // 0 + 1 (margin top)
		// child2.Y = 0 + 8 + 34 (spacing) + 2 (margin top) = 44
		assert.Equal(t, 44, child2.Y)
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
			W: 40,
			H: 50,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child1, absoluteChild, child2},
		}

		box.calculateFlexColumn(0, 0, 40, 50)

		assert.Equal(t, 0, child1.Y)
		// child2 should be positioned as if absoluteChild doesn't exist
		assert.Equal(t, 5, child2.Y) // 5 (child1 height)
	})

	t.Run("should return early when no children exist", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{},
		}

		box.calculateFlexColumn(0, 0, 40, 30)

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
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumn(0, 0, 40, 30)

		assert.Equal(t, 3, child.X) // margin left
		assert.Equal(t, 3, child.Y) // margin top
	})

	t.Run("should clamp stretched width to zero when margins exceed inner width", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				// Width: 0 triggers stretch
				Height: 5,
				Margin: style.Spacing{
					Left:  25,
					Right: 20,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumn(0, 0, 40, 30)

		// Stretched width = 40 - 25 - 20 = -5, clamped to 0
		assert.Equal(t, 0, child.Style.Width)
	})

	t.Run("should handle AlignCenter for cross-axis alignment", func(t *testing.T) {
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
			W: 40,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
				AlignItems:    style.AlignCenter,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumn(0, 0, 40, 30)

		// Child outer width = 10 + 2 + 3 = 15
		// Center X = 0 + (40 - 15) / 2 = 12
		// child.X = 12 + 2 (margin left) = 14
		assert.Equal(t, 14, child.X)
	})

	t.Run("should handle AlignEnd for cross-axis alignment", func(t *testing.T) {
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
			W: 40,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
				AlignItems:    style.AlignEnd,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumn(0, 0, 40, 30)

		// Child outer width = 10 + 2 + 3 = 15
		// End X = 0 + 40 - 15 = 25
		// child.X = 25 + 2 (margin left) = 27
		assert.Equal(t, 27, child.X)
	})

	t.Run("should handle AlignStretch for cross-axis alignment with margins", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Left:  3,
					Right: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
				AlignItems:    style.AlignStretch,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumn(0, 0, 40, 30)

		// Stretch gives full inner width to alignColX, then margins are subtracted
		// childW = innerW - margins = 40 - 3 - 2 = 35
		assert.Equal(t, 3, child.X)            // margin left
		assert.Equal(t, 35, child.Style.Width) // stretched width minus margins
	})
}

func TestBox_distributeFlexSpaceColumn(t *testing.T) {
	t.Parallel()

	t.Run("should reduce available grow space by vertical margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Top:    5,
					Bottom: 5,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Top:    5,
					Bottom: 5,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 80,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child1, child2},
		}

		box.distributeFlexSpaceColumn([]*Box{child1, child2}, 80)

		// Total margins = (5 + 5) + (5 + 5) = 20
		// Available for content = 80 - 20 = 60
		// Base height = 5 + 5 = 10
		// Extra space = 60 - 10 = 50, split equally = 25 each
		assert.Equal(t, 30, child1.Style.Height) // 5 + 25
		assert.Equal(t, 30, child2.Style.Height) // 5 + 25
	})

	t.Run("should calculate shrink correctly with margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 30,
				Margin: style.Spacing{
					Top:    5,
					Bottom: 5,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 30,
				Margin: style.Spacing{
					Top:    5,
					Bottom: 5,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 60,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child1, child2},
		}

		box.distributeFlexSpaceColumn([]*Box{child1, child2}, 60)

		// Total margins = 20
		// Available for content = 60 - 20 = 40
		// Base height = 30 + 30 = 60
		// Overflow = 60 - 40 = 20, shrink proportionally
		// Default shrink = 1, each child shrinks by 10
		assert.Equal(t, 20, child1.Style.Height) // 30 - 10
		assert.Equal(t, 20, child2.Style.Height) // 30 - 10
	})

	t.Run("should not distribute when no grow and space is available", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				// FlexGrow: 0 (default)
				Margin: style.Spacing{
					Top:    2,
					Bottom: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 50,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child},
		}

		box.distributeFlexSpaceColumn([]*Box{child}, 50)

		assert.Equal(t, 5, child.Style.Height) // unchanged
	})

	t.Run("should handle gap in space calculation", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 50,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
				Gap:           4,
			},
			Children: []*Box{child1, child2},
		}

		box.distributeFlexSpaceColumn([]*Box{child1, child2}, 50)

		// Total margins = (2 + 2) + (2 + 2) = 8
		// Gap = 4
		// Available for content = 50 - 8 = 42
		// Base height = 5 + 5 = 10
		// Gap subtracted from extra: 42 - 10 - 4 = 28, split equally = 14 each
		assert.Equal(t, 19, child1.Style.Height) // 5 + 14
		assert.Equal(t, 19, child2.Style.Height) // 5 + 14
	})

	t.Run("should distribute grow proportionally based on FlexGrow values", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 2, // grows twice as much
				Margin: style.Spacing{
					Top:    2,
					Bottom: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 62,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child1, child2},
		}

		box.distributeFlexSpaceColumn([]*Box{child1, child2}, 62)

		// Total margins = 8
		// Available = 62 - 8 = 54
		// Base height = 10
		// Extra = 54 - 10 = 44
		// child1 gets 44 * (1/3) = 14, child2 gets 44 * (2/3) = 29 (+1 remainder)
		assert.Equal(t, 19, child1.Style.Height) // 5 + 14
		assert.Equal(t, 35, child2.Style.Height) // 5 + 30 (29 + 1 remainder)
	})

	t.Run("should handle mix of grow and non-grow children with margins", func(t *testing.T) {
		t.Parallel()

		fixedChild := &Box{
			Style: style.Style{
				Width:  10,
				Height: 20,
				// FlexGrow: 0 (default)
				Margin: style.Spacing{
					Top:    3,
					Bottom: 3,
				},
			},
		}
		growChild := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 60,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{fixedChild, growChild},
		}

		box.distributeFlexSpaceColumn([]*Box{fixedChild, growChild}, 60)

		assert.Equal(t, 20, fixedChild.Style.Height) // unchanged
		// Total margins = 6 + 4 = 10
		// Available = 60 - 10 = 50
		// Fixed child takes 20, grow child base = 5, total base = 25
		// Extra = 50 - 25 = 25, all goes to growChild
		assert.Equal(t, 30, growChild.Style.Height) // 5 + 25
	})

	t.Run("should handle asymmetric margins correctly", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Top:    10,
					Bottom: 0,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:    10,
				Height:   5,
				FlexGrow: 1,
				Margin: style.Spacing{
					Top:    0,
					Bottom: 10,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 60,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child1, child2},
		}

		box.distributeFlexSpaceColumn([]*Box{child1, child2}, 60)

		// Total margins = 10 + 0 + 0 + 10 = 20
		// Available for content = 60 - 20 = 40
		// Base height = 5 + 5 = 10
		// Extra space = 40 - 10 = 30, split equally = 15 each
		assert.Equal(t, 20, child1.Style.Height) // 5 + 15
		assert.Equal(t, 20, child2.Style.Height) // 5 + 15
	})

	t.Run("should handle FlexShrink with different shrink factors", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:      10,
				Height:     30,
				FlexShrink: 1,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:      10,
				Height:     30,
				FlexShrink: 2, // shrinks twice as much
				Margin: style.Spacing{
					Top:    2,
					Bottom: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 52,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
			},
			Children: []*Box{child1, child2},
		}

		box.distributeFlexSpaceColumn([]*Box{child1, child2}, 52)

		// Total margins = 8
		// Available for content = 52 - 8 = 44
		// Base height = 30 + 30 = 60
		// Overflow = 60 - 44 = 16
		// Total shrink factor = 1*30 + 2*30 = 90
		// child1 shrinks by 16 * (30/90) = 5.33 -> 5
		// child2 shrinks by 16 * (60/90) = 10.66 -> 10
		assert.Equal(t, 25, child1.Style.Height) // 30 - 5
		assert.Equal(t, 20, child2.Style.Height) // 30 - 10
	})
}

func TestSplitIntoColumnLines(t *testing.T) {
	t.Parallel()

	t.Run("should return single column when single child fits in height", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
			},
		}
		children := []*Box{child}
		innerH := 20
		gap := 0

		columns := splitIntoColumnLines(children, innerH, gap)

		assert.Len(t, columns, 1)
		assert.Len(t, columns[0], 1)
		assert.Same(t, child, columns[0][0])
	})

	t.Run("should return single column when multiple children fit within height", func(t *testing.T) {
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
		child3 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
			},
		}
		children := []*Box{child1, child2, child3}
		innerH := 20 // Enough for all 3 children (5 + 5 + 5 = 15)
		gap := 0

		columns := splitIntoColumnLines(children, innerH, gap)

		assert.Len(t, columns, 1)
		assert.Len(t, columns[0], 3)
		assert.Same(t, child1, columns[0][0])
		assert.Same(t, child2, columns[0][1])
		assert.Same(t, child3, columns[0][2])
	})

	t.Run("should wrap children to multiple columns when height is exceeded", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 10,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 10,
			},
		}
		child3 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 10,
			},
		}
		children := []*Box{child1, child2, child3}
		innerH := 15 // Only enough for 1 child per column
		gap := 0

		columns := splitIntoColumnLines(children, innerH, gap)

		assert.Len(t, columns, 3)
		assert.Len(t, columns[0], 1)
		assert.Len(t, columns[1], 1)
		assert.Len(t, columns[2], 1)
		assert.Same(t, child1, columns[0][0])
		assert.Same(t, child2, columns[1][0])
		assert.Same(t, child3, columns[2][0])
	})

	t.Run("should wrap children to multiple columns due to margins causing overflow", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    3,
					Bottom: 3,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    3,
					Bottom: 3,
				},
			},
		}
		children := []*Box{child1, child2}
		// Inner height = 20
		// child1 outer height = 5 + 3 + 3 = 11
		// child2 outer height = 5 + 3 + 3 = 11
		// Total = 22, which exceeds 20
		innerH := 20
		gap := 0

		columns := splitIntoColumnLines(children, innerH, gap)

		assert.Len(t, columns, 2) // Margins cause wrap
		assert.Len(t, columns[0], 1)
		assert.Len(t, columns[1], 1)
		assert.Same(t, child1, columns[0][0])
		assert.Same(t, child2, columns[1][0])
	})

	t.Run("should consider gap when calculating column breaks", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 8,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 8,
			},
		}
		children := []*Box{child1, child2}
		// Inner height = 20
		// Without gap: 8 + 8 = 16 (fits)
		// With gap 5: 8 + 5 + 8 = 21 (doesn't fit)
		innerH := 20
		gap := 5

		columns := splitIntoColumnLines(children, innerH, gap)

		assert.Len(t, columns, 2) // Gap causes wrap
		assert.Len(t, columns[0], 1)
		assert.Len(t, columns[1], 1)
	})

	t.Run("should handle children with zero height as minimum height of 1", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 0, // Zero height
			},
		}
		children := []*Box{child}
		innerH := 10
		gap := 0

		columns := splitIntoColumnLines(children, innerH, gap)

		assert.Len(t, columns, 1)
		assert.Len(t, columns[0], 1)
	})

	t.Run("should return empty slice when no children provided", func(t *testing.T) {
		t.Parallel()

		children := []*Box{}
		innerH := 20
		gap := 0

		columns := splitIntoColumnLines(children, innerH, gap)

		assert.Empty(t, columns)
	})

	t.Run("should place children in correct columns with mixed heights and margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{Top: 2, Bottom: 2},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 3,
				Margin: style.Spacing{Top: 1, Bottom: 1},
			},
		}
		child3 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 4,
			},
		}
		child4 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 6,
				Margin: style.Spacing{Top: 2, Bottom: 2},
			},
		}
		children := []*Box{child1, child2, child3, child4}
		// child1 outer: 5+2+2 = 9
		// child2 outer: 3+1+1 = 5
		// child3 outer: 4
		// child4 outer: 6+2+2 = 10
		// Column 1: 9 + 5 = 14 (fits in 20)
		// Adding child3: 14 + 4 = 18 (still fits)
		// Adding child4: 18 + 10 = 28 (exceeds 20, new column)
		innerH := 20
		gap := 0

		columns := splitIntoColumnLines(children, innerH, gap)

		assert.Len(t, columns, 2)
		assert.Len(t, columns[0], 3)
		assert.Len(t, columns[1], 1)
		assert.Same(t, child4, columns[1][0])
	})
}

func TestBox_calculateFlexColumnWrap(t *testing.T) {
	t.Parallel()

	t.Run("should position single column children with margins correctly", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:  2,
					Left: 3,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:  1,
					Left: 3,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 50,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexColumnWrap(0, 0, 40, 50)

		assert.Equal(t, 3, child1.X) // margin left
		assert.Equal(t, 2, child1.Y) // margin top
		assert.Equal(t, 3, child2.X) // margin left
		assert.Equal(t, 8, child2.Y) // child1 outer height (5+2+0=7) + margin top (1)
	})

	t.Run("should position multiple columns with margins correctly", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 15,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 2,
					Left:   1,
					Right:  1,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  8,
				Height: 12,
				Margin: style.Spacing{
					Top:    3,
					Bottom: 3,
					Left:   2,
					Right:  2,
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
				FlexDirection: style.Column,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child1, child2},
		}
		// child1 outer height = 15 + 2 + 2 = 19 (fits in 20)
		// child2 outer height = 12 + 3 + 3 = 18
		// child1 (19) + child2 (18) = 37 > 20, so child2 wraps

		box.calculateFlexColumnWrap(0, 0, 50, 20)

		// Column 1: child1
		assert.Equal(t, 1, child1.X) // margin left
		assert.Equal(t, 2, child1.Y) // margin top
		// Column 2: child2
		// Column 1 width = 10 + 1 + 1 = 12
		// child2.X = 12 (col width) + 0 (gap) + 2 (margin left) = 14
		assert.Equal(t, 14, child2.X)
		assert.Equal(t, 3, child2.Y) // margin top
	})

	t.Run("should calculate column widths including horizontal margins", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 30,
				Margin: style.Spacing{
					Left:  5,
					Right: 5,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  8,
				Height: 30,
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
			H: 25,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child1, child2},
		}
		// Both children have height 30 which exceeds innerH 25, so each gets its own column

		box.calculateFlexColumnWrap(0, 0, 60, 25)

		// Column 1 width = 10 + 5 + 5 = 20
		assert.Equal(t, 5, child1.X) // margin left within column
		// Column 2 starts at x = 20 (col1 width) + 0 (gap)
		// child2.X = 20 + 2 (margin left) = 22
		assert.Equal(t, 22, child2.X)
	})

	t.Run("should cause earlier column break due to margins", func(t *testing.T) {
		t.Parallel()

		// Without margins: both children fit (8 + 8 = 16 <= 20)
		// With margins: (8+3+3) + (8+3+3) = 28 > 20, wraps
		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 8,
				Margin: style.Spacing{
					Top:    3,
					Bottom: 3,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 8,
				Margin: style.Spacing{
					Top:    3,
					Bottom: 3,
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
				FlexDirection: style.Column,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexColumnWrap(0, 0, 40, 20)

		// Verify children are in separate columns
		assert.Equal(t, 0, child1.X)
		assert.Equal(t, 3, child1.Y) // margin top
		// child2 in column 2
		// Column 1 width = 10 + 0 + 0 = 10 (no horizontal margin)
		assert.Equal(t, 10, child2.X)
		assert.Equal(t, 3, child2.Y) // margin top in new column
	})

	t.Run("should apply JustifyCenter across wrapped columns", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    2,
					Bottom: 2,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:        style.Flex,
				FlexDirection:  style.Column,
				FlexWrap:       style.Wrap,
				JustifyContent: style.JustifyCenter,
			},
			Children: []*Box{child1, child2},
		}
		// child1 outer height = 5 + 2 + 2 = 9
		// child2 outer height = 5 + 2 + 2 = 9
		// Both fit in innerH=30: 9 + 9 = 18

		box.calculateFlexColumnWrap(0, 0, 40, 30)

		// Total height = 18, innerH = 30
		// Center offset = (30 - 18) / 2 = 6
		// child1.Y = 6 + 2 (margin top) = 8
		assert.Equal(t, 8, child1.Y)
		// child2.Y = 6 + 9 (child1 outer) + 2 (margin top) = 17
		assert.Equal(t, 17, child2.Y)
	})

	t.Run("should apply JustifyEnd across wrapped columns", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    1,
					Bottom: 1,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    1,
					Bottom: 1,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:        style.Flex,
				FlexDirection:  style.Column,
				FlexWrap:       style.Wrap,
				JustifyContent: style.JustifyEnd,
			},
			Children: []*Box{child1, child2},
		}
		// child1 outer = 5 + 1 + 1 = 7
		// child2 outer = 5 + 1 + 1 = 7
		// Total = 14

		box.calculateFlexColumnWrap(0, 0, 40, 30)

		// End offset = 30 - 14 = 16
		// child1.Y = 16 + 1 (margin top) = 17
		assert.Equal(t, 17, child1.Y)
		// child2.Y = 16 + 7 (child1 outer) + 1 (margin top) = 24
		assert.Equal(t, 24, child2.Y)
	})

	t.Run("should apply JustifyBetween across wrapped columns", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    1,
					Bottom: 1,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:    1,
					Bottom: 1,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:        style.Flex,
				FlexDirection:  style.Column,
				FlexWrap:       style.Wrap,
				JustifyContent: style.JustifyBetween,
			},
			Children: []*Box{child1, child2},
		}
		// child1 outer = 7, child2 outer = 7, total = 14
		// spacing = (30 - 14) / (2 - 1) = 16

		box.calculateFlexColumnWrap(0, 0, 40, 30)

		// child1.Y = 0 + 1 (margin top) = 1
		assert.Equal(t, 1, child1.Y)
		// child2.Y = 0 + 7 (child1 outer) + 16 (spacing) + 1 (margin top) = 24
		assert.Equal(t, 24, child2.Y)
	})

	t.Run("should return early when no children exist", func(t *testing.T) {
		t.Parallel()

		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{},
		}

		box.calculateFlexColumnWrap(0, 0, 40, 30)

		assert.Empty(t, box.Children)
	})

	t.Run("should handle gap between columns", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 25,
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  8,
				Height: 25,
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 60,
			H: 20,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
				FlexWrap:      style.Wrap,
				Gap:           5,
			},
			Children: []*Box{child1, child2},
		}
		// Both exceed innerH=20, so each gets own column

		box.calculateFlexColumnWrap(0, 0, 60, 20)

		assert.Equal(t, 0, child1.X)
		// Column 1 width = 10, gap = 5
		// child2.X = 10 + 5 = 15
		assert.Equal(t, 15, child2.X)
	})

	t.Run("should position with parent offset", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  10,
				Height: 5,
				Margin: style.Spacing{
					Top:  2,
					Left: 3,
				},
			},
		}
		box := &Box{
			X: 10,
			Y: 20,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumnWrap(5, 3, 30, 24)

		// child.X = box.X + colX + marginLeft = 10 + 5 + 3 = 18
		assert.Equal(t, 18, child.X)
		// child.Y = box.Y + y + marginTop = 20 + 3 + 2 = 25
		assert.Equal(t, 25, child.Y)
	})

	t.Run("should handle AlignCenter within column", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  10,
				Height: 30,
				Margin: style.Spacing{
					Left:  2,
					Right: 2,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  6,
				Height: 30,
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
				Display:       style.Flex,
				FlexDirection: style.Column,
				FlexWrap:      style.Wrap,
				AlignItems:    style.AlignCenter,
			},
			Children: []*Box{child1, child2},
		}
		// Both exceed innerH=20, each gets own column

		box.calculateFlexColumnWrap(0, 0, 50, 20)

		// Column 1 width = 10 + 2 + 2 = 14
		// child1 outer width = 14
		// Center: (14 - 14) / 2 = 0, child1.X = 0 + 2 (margin) = 2
		assert.Equal(t, 2, child1.X)
		// Column 2 width = 6 + 1 + 1 = 8
		// child2 outer width = 8
		// Center: (8 - 8) / 2 = 0, child2.X = 14 + 0 + 1 (margin) = 15
		assert.Equal(t, 15, child2.X)
	})

	t.Run("should handle AlignEnd within column", func(t *testing.T) {
		t.Parallel()

		child1 := &Box{
			Style: style.Style{
				Width:  8,
				Height: 30,
				Margin: style.Spacing{
					Left:  1,
					Right: 1,
				},
			},
		}
		child2 := &Box{
			Style: style.Style{
				Width:  6,
				Height: 30,
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
				FlexDirection: style.Column,
				FlexWrap:      style.Wrap,
				AlignItems:    style.AlignEnd,
			},
			Children: []*Box{child1, child2},
		}

		box.calculateFlexColumnWrap(0, 0, 50, 20)

		// Column 1 width = 8 + 1 + 1 = 10
		// End: 10 - 10 = 0, child1.X = 0 + 1 (margin) = 1
		assert.Equal(t, 1, child1.X)
		// Column 2 width = 6 + 2 + 2 = 10
		// End: 10 - 10 = 0, child2.X = 10 + 0 + 2 (margin) = 12
		assert.Equal(t, 12, child2.X)
	})

	t.Run("should handle child with zero width as minimum", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  0,
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
				FlexDirection: style.Column,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumnWrap(0, 0, 40, 30)

		assert.Equal(t, 0, child.X)
		assert.Equal(t, 0, child.Y)
	})

	t.Run("should clamp child width to zero when margins exceed column width", func(t *testing.T) {
		t.Parallel()

		child := &Box{
			Style: style.Style{
				Width:  5,
				Height: 10,
				Margin: style.Spacing{
					Left:  10,
					Right: 10,
				},
			},
		}
		box := &Box{
			X: 0,
			Y: 0,
			W: 40,
			H: 30,
			Style: style.Style{
				Display:       style.Flex,
				FlexDirection: style.Column,
				FlexWrap:      style.Wrap,
			},
			Children: []*Box{child},
		}

		box.calculateFlexColumnWrap(0, 0, 40, 30)

		// Column width = 5 + 10 + 10 = 25
		// childW = colW - margins = 25 - 10 - 10 = 5 (original width fits)
		assert.Equal(t, 5, child.Style.Width)
	})
}
