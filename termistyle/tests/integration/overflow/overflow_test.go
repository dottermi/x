package overflow_test

import (
	"testing"

	"github.com/dottermi/termitest/assert"
	ts "github.com/dottermi/x/termistyle"
	. "github.com/dottermi/x/termistyle/tests/integration/helper"
)

func TestGolden_Overflow_Hidden(t *testing.T) {
	// Container with overflow hidden - child extends beyond bounds
	root := ts.NewBox(ts.Style{
		Width:      40,
		Height:     15,
		Display:    ts.Flex,
		Background: ColorBg,
	})

	// Container with overflow hidden
	container := ts.NewBox(ts.Style{
		Width:      25,
		Height:     10,
		Display:    ts.Flex,
		Overflow:   ts.OverflowHidden, // Clips children
		Background: ColorSurface,
		Border:     ts.BorderAllWithColor(ts.BorderSingle, ColorAccent),
	})

	// Child that extends beyond container bounds
	child := ts.NewBox(ts.Style{
		Width:      35, // Wider than container (25)
		Height:     15, // Taller than container (10)
		Background: ColorCyan,
		Border:     ts.BorderAllWithColor(ts.BorderSingle, ColorCyan2),
	})
	container.AddChild(child)
	root.AddChild(container)

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}

func TestGolden_Overflow_Visible(t *testing.T) {
	// Container with overflow visible (default) - child extends beyond bounds
	root := ts.NewBox(ts.Style{
		Width:      40,
		Height:     15,
		Display:    ts.Flex,
		Background: ColorBg,
		Padding:    ts.Spacing{Left: 2, Top: 2},
	})

	// Container with overflow visible (default)
	container := ts.NewBox(ts.Style{
		Width:      20,
		Height:     8,
		Display:    ts.Flex,
		Overflow:   ts.OverflowVisible, // Does not clip
		Background: ColorSurface,
		Border:     ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
	})

	// Child that extends beyond container bounds
	child := ts.NewBox(ts.Style{
		Width:      30, // Wider than container (20)
		Height:     12, // Taller than container (8)
		Background: ColorBlue,
		Border:     ts.BorderAllWithColor(ts.BorderSingle, ColorCyan),
	})
	container.AddChild(child)
	root.AddChild(container)

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}
