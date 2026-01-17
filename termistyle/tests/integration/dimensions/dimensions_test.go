package dimensions_test

import (
	"testing"

	"github.com/dottermi/termitest/assert"
	ts "github.com/dottermi/x/termistyle"
	. "github.com/dottermi/x/termistyle/tests/integration/helper"
)

func TestGolden_Dimensions_MinMax(t *testing.T) {
	// Row 1: MaxWidth constraint - child wants 30, but max is 15
	row1 := ts.NewBox(ts.Style{
		Width:         38,
		Height:        5,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		AlignItems:    ts.AlignCenter,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Padding:       ts.SpacingAll(1),
	})
	row1.AddChild(ts.NewBox(ts.Style{
		Width:      30, // Wants 30
		MaxWidth:   15, // But limited to 15
		Height:     2,
		Background: ColorCyan,
	}))

	// Row 2: MinWidth constraint - container is small, but child has min
	row2 := ts.NewBox(ts.Style{
		Width:         38,
		Height:        5,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		AlignItems:    ts.AlignCenter,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Padding:       ts.SpacingAll(1),
	})
	row2.AddChild(ts.NewBox(ts.Style{
		Width:      5,  // Wants 5
		MinWidth:   20, // But minimum is 20
		Height:     2,
		Background: ColorBlue,
	}))

	// Row 3: MaxHeight constraint
	row3 := ts.NewBox(ts.Style{
		Width:         38,
		Height:        8,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		AlignItems:    ts.AlignStretch,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Padding:       ts.SpacingAll(1),
	})
	row3.AddChild(ts.NewBox(ts.Style{
		Width:     10,
		MaxHeight: 3, // Limited height even with stretch
		Background: ColorDark2,
	}))
	row3.AddChild(ts.NewBox(ts.Style{
		Width:      10,
		Background: ColorAccent, // Will stretch to full height
	}))

	// Row 4: MinHeight constraint
	row4 := ts.NewBox(ts.Style{
		Width:         38,
		Height:        5,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		AlignItems:    ts.AlignStart,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Padding:       ts.SpacingAll(1),
	})
	row4.AddChild(ts.NewBox(ts.Style{
		Width:     10,
		Height:    1,
		MinHeight: 2, // Min overrides small height
		Background: ColorCyan2,
	}))

	root := ts.NewBox(ts.Style{
		Width:         40,
		Height:        28,
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

func TestGolden_Dimensions_AspectRatio(t *testing.T) {
	// Row 1: Width given, height calculated from aspect ratio (2:1)
	row1 := ts.NewBox(ts.Style{
		Width:         38,
		Height:        8,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		AlignItems:    ts.AlignStart,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Padding:       ts.SpacingAll(1),
		Gap:           2,
	})
	// Width 20, aspect 2:1 -> height 10
	row1.AddChild(ts.NewBox(ts.Style{
		Width:       20,
		AspectRatio: 2.0, // Width/Height = 2, so height = 10
		Background:  ColorCyan,
	}))

	// Row 2: Height given, width calculated from aspect ratio
	row2 := ts.NewBox(ts.Style{
		Width:         38,
		Height:        8,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		AlignItems:    ts.AlignStart,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Padding:       ts.SpacingAll(1),
		Gap:           2,
	})
	// Height 4, aspect 2:1 -> width 8
	row2.AddChild(ts.NewBox(ts.Style{
		Height:      4,
		AspectRatio: 2.0, // Width/Height = 2, so width = 8
		Background:  ColorBlue,
	}))
	// Height 4, aspect 0.5:1 -> width 2
	row2.AddChild(ts.NewBox(ts.Style{
		Height:      4,
		AspectRatio: 0.5, // Width/Height = 0.5, so width = 2
		Background:  ColorDark2,
	}))

	// Row 3: Square aspect ratio (1:1)
	row3 := ts.NewBox(ts.Style{
		Width:         38,
		Height:        8,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		AlignItems:    ts.AlignStart,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Padding:       ts.SpacingAll(1),
		Gap:           2,
	})
	row3.AddChild(ts.NewBox(ts.Style{
		Width:       6,
		AspectRatio: 1.0, // Square
		Background:  ColorAccent,
	}))
	row3.AddChild(ts.NewBox(ts.Style{
		Height:      5,
		AspectRatio: 1.0, // Square
		Background:  ColorCyan2,
	}))

	root := ts.NewBox(ts.Style{
		Width:         40,
		Height:        28,
		Display:       ts.Flex,
		FlexDirection: ts.Column,
		Gap:           1,
	})
	root.AddChild(row1)
	root.AddChild(row2)
	root.AddChild(row3)

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}
