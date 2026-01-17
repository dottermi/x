package display_test

import (
	"testing"

	"github.com/dottermi/termitest/assert"
	ts "github.com/dottermi/x/termistyle"
	"github.com/dottermi/x/termistyle/style"
	. "github.com/dottermi/x/termistyle/tests/integration/helper"
)

func TestGolden_Display_None(t *testing.T) {
	// Container with some visible children and some hidden
	root := ts.NewBox(ts.Style{
		Width:         40,
		Height:        15,
		Display:       ts.Flex,
		FlexDirection: ts.Column,
		Background:    ColorBg,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Padding:       ts.SpacingAll(1),
		Gap:           1,
	})

	// Visible item 1
	root.AddChild(ts.NewBox(ts.Style{
		Width:      36,
		Height:     2,
		Display:    ts.Flex,
		Background: ColorCyan,
	}))

	// Hidden item (Display: None) - should not appear and not take space
	root.AddChild(ts.NewBox(ts.Style{
		Width:      36,
		Height:     2,
		Display:    style.None, // Hidden!
		Background: ColorAccent,
	}))

	// Visible item 2 - should be right after item 1, not after hidden space
	root.AddChild(ts.NewBox(ts.Style{
		Width:      36,
		Height:     2,
		Display:    ts.Flex,
		Background: ColorBlue,
	}))

	// Another hidden item
	root.AddChild(ts.NewBox(ts.Style{
		Width:      36,
		Height:     2,
		Display:    style.None, // Hidden!
		Background: ColorCyan2,
	}))

	// Visible item 3
	root.AddChild(ts.NewBox(ts.Style{
		Width:      36,
		Height:     2,
		Display:    ts.Flex,
		Background: ColorDark2,
	}))

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}
