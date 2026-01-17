package flex_test

import (
	"testing"

	"github.com/dottermi/termitest/assert"
	ts "github.com/dottermi/x/termistyle"
	. "github.com/dottermi/x/termistyle/tests/integration/helper"
)

func TestGolden_Flex_Grow(t *testing.T) {
	// Row 1: Equal grow (1:1:1)
	row1 := ts.NewBox(ts.Style{
		Width:         38,
		Height:        3,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
	})
	row1.AddChild(ts.NewBox(ts.Style{FlexGrow: 1, Height: 1, Background: ColorCyan}))
	row1.AddChild(ts.NewBox(ts.Style{FlexGrow: 1, Height: 1, Background: ColorBlue}))
	row1.AddChild(ts.NewBox(ts.Style{FlexGrow: 1, Height: 1, Background: ColorDark2}))

	// Row 2: Different grow ratios (1:2:1)
	row2 := ts.NewBox(ts.Style{
		Width:         38,
		Height:        3,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
	})
	row2.AddChild(ts.NewBox(ts.Style{FlexGrow: 1, Height: 1, Background: ColorCyan}))
	row2.AddChild(ts.NewBox(ts.Style{FlexGrow: 2, Height: 1, Background: ColorBlue}))
	row2.AddChild(ts.NewBox(ts.Style{FlexGrow: 1, Height: 1, Background: ColorDark2}))

	// Row 3: Mixed fixed + grow (fixed, grow, fixed)
	row3 := ts.NewBox(ts.Style{
		Width:         38,
		Height:        3,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
	})
	row3.AddChild(ts.NewBox(ts.Style{Width: 8, Height: 1, Background: ColorCyan}))
	row3.AddChild(ts.NewBox(ts.Style{FlexGrow: 1, Height: 1, Background: ColorBlue}))
	row3.AddChild(ts.NewBox(ts.Style{Width: 8, Height: 1, Background: ColorDark2}))

	// Row 4: One item grows (0:1:0)
	row4 := ts.NewBox(ts.Style{
		Width:         38,
		Height:        3,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
	})
	row4.AddChild(ts.NewBox(ts.Style{Width: 6, Height: 1, Background: ColorCyan}))
	row4.AddChild(ts.NewBox(ts.Style{FlexGrow: 1, Height: 1, Background: ColorAccent}))
	row4.AddChild(ts.NewBox(ts.Style{Width: 6, Height: 1, Background: ColorDark2}))

	root := ts.NewBox(ts.Style{
		Width:         40,
		Height:        16,
		Display:       ts.Flex,
		FlexDirection: ts.Column,
		Gap:           1,
	})
	root.AddChild(row1)
	root.AddChild(row2)
	root.AddChild(row3)
	root.AddChild(row4)

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}

func TestGolden_Flex_AlignSelf(t *testing.T) {
	// Container with AlignItems: Start, but children with different AlignSelf
	root := ts.NewBox(ts.Style{
		Width:         40,
		Height:        12,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		AlignItems:    ts.AlignStart, // Default alignment
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Gap:           1,
		Padding:       ts.SpacingAll(1),
	})

	// Item 1: Uses parent's AlignItems (Start)
	root.AddChild(ts.NewBox(ts.Style{
		Width:      6,
		Height:     3,
		Background: ColorCyan,
	}))

	// Item 2: AlignSelf Center
	root.AddChild(ts.NewBox(ts.Style{
		Width:      6,
		Height:     3,
		AlignSelf:  ts.AlignCenter,
		Background: ColorBlue,
	}))

	// Item 3: AlignSelf End
	root.AddChild(ts.NewBox(ts.Style{
		Width:      6,
		Height:     3,
		AlignSelf:  ts.AlignEnd,
		Background: ColorDark2,
	}))

	// Item 4: AlignSelf Stretch
	root.AddChild(ts.NewBox(ts.Style{
		Width:      6,
		AlignSelf:  ts.AlignStretch,
		Background: ColorAccent,
	}))

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}
