package position_test

import (
	"testing"

	"github.com/dottermi/termitest/assert"
	ts "github.com/dottermi/x/termistyle"
	. "github.com/dottermi/x/termistyle/tests/integration/helper"
)

func TestGolden_Position_Absolute(t *testing.T) {
	// Container with absolute positioned children at specific offsets
	root := ts.NewBox(ts.Style{
		Width:      40,
		Height:     15,
		Display:    ts.Flex,
		Background: ColorDark,
		Border:     ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
	})

	// Box at top-left (offset 2, 1)
	box1 := ts.NewBox(ts.Style{
		Width:      10,
		Height:     3,
		Position:   ts.Absolute,
		Background: ColorCyan,
	})
	box1.X = 2
	box1.Y = 1
	root.AddChild(box1)

	// Box at center-ish (offset 15, 5)
	box2 := ts.NewBox(ts.Style{
		Width:      12,
		Height:     4,
		Position:   ts.Absolute,
		Background: ColorBlue,
	})
	box2.X = 15
	box2.Y = 5
	root.AddChild(box2)

	// Box at bottom-right (offset 25, 10)
	box3 := ts.NewBox(ts.Style{
		Width:      10,
		Height:     3,
		Position:   ts.Absolute,
		Background: ColorDark2,
	})
	box3.X = 25
	box3.Y = 10
	root.AddChild(box3)

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}

func TestGolden_Position_ZIndex(t *testing.T) {
	// Container with overlapping boxes at different z-index levels
	root := ts.NewBox(ts.Style{
		Width:      35,
		Height:     15,
		Display:    ts.Flex,
		Background: ColorDark,
		Border:     ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
	})

	// Bottom layer (z-index 0) - largest, at back
	bottom := ts.NewBox(ts.Style{
		Width:      20,
		Height:     10,
		Position:   ts.Absolute,
		ZIndex:     0,
		Background: ColorDark2,
		Border:     ts.BorderAllWithColor(ts.BorderSingle, ColorBlue),
	})
	bottom.X = 2
	bottom.Y = 2
	root.AddChild(bottom)

	// Middle layer (z-index 1) - overlaps bottom
	middle := ts.NewBox(ts.Style{
		Width:      15,
		Height:     8,
		Position:   ts.Absolute,
		ZIndex:     1,
		Background: ColorBlue,
		Border:     ts.BorderAllWithColor(ts.BorderSingle, ColorCyan),
	})
	middle.X = 8
	middle.Y = 4
	root.AddChild(middle)

	// Top layer (z-index 2) - overlaps both
	top := ts.NewBox(ts.Style{
		Width:      12,
		Height:     5,
		Position:   ts.Absolute,
		ZIndex:     2,
		Background: ColorCyan,
		Border:     ts.BorderAllWithColor(ts.BorderSingle, ColorCyan2),
	})
	top.X = 15
	top.Y = 6
	root.AddChild(top)

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}

func TestGolden_Position_Complete(t *testing.T) {
	// Simulate a modal overlay over content
	root := ts.NewBox(ts.Style{
		Width:      50,
		Height:     20,
		Display:    ts.Flex,
		Background: ColorBg,
	})

	// Background content (normal flow)
	content := ts.NewBox(ts.Style{
		Width:         48,
		Height:        18,
		Display:       ts.Flex,
		FlexDirection: ts.Column,
		Background:    ColorSurface,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Padding:       ts.SpacingAll(1),
		Gap:           1,
	})
	content.AddChild(ts.Text("Dashboard", ts.Style{Foreground: ColorText, FontWeight: ts.WeightBold}))
	content.AddChild(ts.Text("Welcome back!", ts.Style{Foreground: ColorMuted}))
	content.AddChild(ts.NewBox(ts.Style{Width: 44, Height: 5, Background: ColorDark2}))
	content.AddChild(ts.NewBox(ts.Style{Width: 44, Height: 5, Background: ColorDark2}))
	root.AddChild(content)

	// Modal overlay (absolute, high z-index)
	modal := ts.NewBox(ts.Style{
		Width:      30,
		Height:     10,
		Position:   ts.Absolute,
		ZIndex:     10,
		Background: ColorDark,
		Border:     ts.BorderAllWithColor(ts.BorderDouble, ColorAccent),
		Padding:    ts.SpacingAll(1),
	})
	modal.X = 10
	modal.Y = 5

	// Modal title
	title := ts.NewBox(ts.Style{
		Width:   26,
		Height:  1,
		Display: ts.Flex,
	})
	title.AddChild(ts.Text("Confirm Action", ts.Style{Foreground: ColorCyan2, FontWeight: ts.WeightBold}))
	modal.AddChild(title)

	// Modal message
	msg := ts.NewBox(ts.Style{
		Width:   26,
		Height:  3,
		Display: ts.Flex,
	})
	msg.AddChild(ts.Text("Are you sure?", ts.Style{Foreground: ColorText}))
	modal.AddChild(msg)

	// Modal buttons row
	buttons := ts.NewBox(ts.Style{
		Width:          26,
		Height:         3,
		Display:        ts.Flex,
		FlexDirection:  ts.Row,
		JustifyContent: ts.JustifyEnd,
		Gap:            2,
	})
	buttons.AddChild(ts.NewBox(ts.Style{
		Width:      8,
		Height:     1,
		Background: ColorBlue,
	}))
	buttons.AddChild(ts.NewBox(ts.Style{
		Width:      8,
		Height:     1,
		Background: ColorCyan,
	}))
	modal.AddChild(buttons)

	root.AddChild(modal)

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}
