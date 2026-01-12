package layout

import (
	"testing"

	"github.com/dottermi/x/termistyle/style"
	"github.com/stretchr/testify/assert"
)

func TestBox_MinMaxDimensions(t *testing.T) {
	t.Parallel()

	t.Run("should constrain width to MaxWidth", func(t *testing.T) {
		t.Parallel()

		box := NewBox(style.Style{
			Width:    100,
			MaxWidth: 50,
			Height:   20,
		})

		box.Calculate()

		assert.Equal(t, 50, box.W)
	})

	t.Run("should constrain height to MaxHeight", func(t *testing.T) {
		t.Parallel()

		box := NewBox(style.Style{
			Width:     100,
			Height:    100,
			MaxHeight: 30,
		})

		box.Calculate()

		assert.Equal(t, 30, box.H)
	})

	t.Run("should expand width to MinWidth", func(t *testing.T) {
		t.Parallel()

		box := NewBox(style.Style{
			Width:    10,
			MinWidth: 50,
			Height:   20,
		})

		box.Calculate()

		assert.Equal(t, 50, box.W)
	})

	t.Run("should expand height to MinHeight", func(t *testing.T) {
		t.Parallel()

		box := NewBox(style.Style{
			Width:     100,
			Height:    10,
			MinHeight: 50,
		})

		box.Calculate()

		assert.Equal(t, 50, box.H)
	})
}

func TestBox_AspectRatio(t *testing.T) {
	t.Parallel()

	t.Run("should calculate height from width when aspect ratio is set", func(t *testing.T) {
		t.Parallel()

		child := NewBox(style.Style{
			Width:       100,
			AspectRatio: 2.0, // width:height = 2:1
		})

		container := NewBox(style.Style{
			Width:   200,
			Height:  200,
			Display: style.Flex,
		})
		container.AddChild(child)

		container.Calculate()

		assert.Equal(t, 100, child.W)
		assert.Equal(t, 50, child.H) // 100 / 2 = 50
	})

	t.Run("should calculate width from height when aspect ratio is set in flex context", func(t *testing.T) {
		t.Parallel()

		child := NewBox(style.Style{
			Height:      50,
			AspectRatio: 2.0, // width:height = 2:1
		})

		container := NewBox(style.Style{
			Width:         200,
			Height:        200,
			Display:       style.Flex,
			FlexDirection: style.Column,
		})
		container.AddChild(child)

		container.Calculate()

		assert.Equal(t, 100, child.W) // 50 * 2 = 100
		assert.Equal(t, 50, child.H)
	})
}

func TestBox_AlignSelf(t *testing.T) {
	t.Parallel()

	t.Run("should override parent AlignItems for specific child", func(t *testing.T) {
		t.Parallel()

		child1 := NewBox(style.Style{
			Width:  20,
			Height: 10,
		})
		child2 := NewBox(style.Style{
			Width:     20,
			Height:    10,
			AlignSelf: style.AlignEnd, // Override to end
		})

		container := NewBox(style.Style{
			Width:         100,
			Height:        50,
			Display:       style.Flex,
			FlexDirection: style.Row,
			AlignItems:    style.AlignStart, // Default is start
		})
		container.AddChild(child1)
		container.AddChild(child2)

		container.Calculate()

		// child1 should be at top (AlignStart)
		assert.Equal(t, 0, child1.Y)
		// child2 should be at bottom (AlignEnd): 50 - 10 = 40
		assert.Equal(t, 40, child2.Y)
	})
}

func TestBox_FlexLayout(t *testing.T) {
	t.Parallel()

	t.Run("should position children in row layout", func(t *testing.T) {
		t.Parallel()

		child1 := NewBox(style.Style{
			Width:  20,
			Height: 10,
		})
		child2 := NewBox(style.Style{
			Width:  30,
			Height: 10,
		})

		container := NewBox(style.Style{
			Width:         100,
			Height:        20,
			Display:       style.Flex,
			FlexDirection: style.Row,
		})
		container.AddChild(child1)
		container.AddChild(child2)

		container.Calculate()

		assert.Equal(t, 0, child1.X)
		assert.Equal(t, 20, child2.X)
	})

	t.Run("should position children in column layout", func(t *testing.T) {
		t.Parallel()

		child1 := NewBox(style.Style{
			Width:  20,
			Height: 10,
		})
		child2 := NewBox(style.Style{
			Width:  20,
			Height: 15,
		})

		container := NewBox(style.Style{
			Width:         40,
			Height:        50,
			Display:       style.Flex,
			FlexDirection: style.Column,
		})
		container.AddChild(child1)
		container.AddChild(child2)

		container.Calculate()

		assert.Equal(t, 0, child1.Y)
		assert.Equal(t, 10, child2.Y)
	})

	t.Run("should apply JustifyCenter", func(t *testing.T) {
		t.Parallel()

		child := NewBox(style.Style{
			Width:  20,
			Height: 10,
		})

		container := NewBox(style.Style{
			Width:          100,
			Height:         20,
			Display:        style.Flex,
			FlexDirection:  style.Row,
			JustifyContent: style.JustifyCenter,
		})
		container.AddChild(child)

		container.Calculate()

		// Child should be centered: (100 - 20) / 2 = 40
		assert.Equal(t, 40, child.X)
	})

	t.Run("should apply AlignCenter", func(t *testing.T) {
		t.Parallel()

		child := NewBox(style.Style{
			Width:  20,
			Height: 10,
		})

		container := NewBox(style.Style{
			Width:         100,
			Height:        50,
			Display:       style.Flex,
			FlexDirection: style.Row,
			AlignItems:    style.AlignCenter,
		})
		container.AddChild(child)

		container.Calculate()

		// Child should be vertically centered: (50 - 10) / 2 = 20
		assert.Equal(t, 20, child.Y)
	})

	t.Run("should distribute space with FlexGrow", func(t *testing.T) {
		t.Parallel()

		child1 := NewBox(style.Style{
			Width:    20,
			Height:   10,
			FlexGrow: 1,
		})
		child2 := NewBox(style.Style{
			Width:    20,
			Height:   10,
			FlexGrow: 1,
		})

		container := NewBox(style.Style{
			Width:         100,
			Height:        20,
			Display:       style.Flex,
			FlexDirection: style.Row,
		})
		container.AddChild(child1)
		container.AddChild(child2)

		container.Calculate()

		// Each child should get 50 (100 / 2)
		assert.Equal(t, 50, child1.W)
		assert.Equal(t, 50, child2.W)
	})
}
