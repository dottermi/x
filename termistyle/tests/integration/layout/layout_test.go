package layout_test

import (
	"testing"

	"github.com/dottermi/termitest/assert"
	ts "github.com/dottermi/x/termistyle"
	. "github.com/dottermi/x/termistyle/tests/integration/helper"
)

func TestGolden_Layout_Direction(t *testing.T) {
	// Row direction
	row := ts.NewBox(ts.Style{
		Width:         25,
		Height:        7,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Gap:           1,
	})
	row.AddChild(ts.NewBox(ts.Style{Width: 5, Height: 3, Background: ColorCyan}))
	row.AddChild(ts.NewBox(ts.Style{Width: 5, Height: 3, Background: ColorBlue}))
	row.AddChild(ts.NewBox(ts.Style{Width: 5, Height: 3, Background: ColorDark2}))

	// Column direction
	col := ts.NewBox(ts.Style{
		Width:         12,
		Height:        14,
		Display:       ts.Flex,
		FlexDirection: ts.Column,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Gap:           1,
	})
	col.AddChild(ts.NewBox(ts.Style{Width: 8, Height: 3, Background: ColorCyan}))
	col.AddChild(ts.NewBox(ts.Style{Width: 8, Height: 3, Background: ColorBlue}))
	col.AddChild(ts.NewBox(ts.Style{Width: 8, Height: 3, Background: ColorDark2}))

	// Container
	root := ts.NewBox(ts.Style{
		Width:         40,
		Height:        16,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		Gap:           2,
	})
	root.AddChild(row)
	root.AddChild(col)

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}

func TestGolden_Layout_Justify(t *testing.T) {
	makeRow := func(justify ts.Justify) *ts.Box {
		box := ts.NewBox(ts.Style{
			Width:          38,
			Height:         3,
			Display:        ts.Flex,
			FlexDirection:  ts.Row,
			JustifyContent: justify,
			Border:         ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		})
		box.AddChild(ts.NewBox(ts.Style{Width: 4, Height: 1, Background: ColorCyan}))
		box.AddChild(ts.NewBox(ts.Style{Width: 4, Height: 1, Background: ColorBlue}))
		box.AddChild(ts.NewBox(ts.Style{Width: 4, Height: 1, Background: ColorDark2}))
		return box
	}

	root := ts.NewBox(ts.Style{
		Width:         40,
		Height:        20,
		Display:       ts.Flex,
		FlexDirection: ts.Column,
		Gap:           1,
	})
	root.AddChild(makeRow(ts.JustifyStart))
	root.AddChild(makeRow(ts.JustifyCenter))
	root.AddChild(makeRow(ts.JustifyEnd))
	root.AddChild(makeRow(ts.JustifyBetween))
	root.AddChild(makeRow(ts.JustifyAround))

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}

func TestGolden_Layout_Align(t *testing.T) {
	makeRow := func(align ts.Align) *ts.Box {
		box := ts.NewBox(ts.Style{
			Width:         18,
			Height:        7,
			Display:       ts.Flex,
			FlexDirection: ts.Row,
			AlignItems:    align,
			Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
			Gap:           1,
		})
		box.AddChild(ts.NewBox(ts.Style{Width: 3, Height: 2, Background: ColorCyan}))
		box.AddChild(ts.NewBox(ts.Style{Width: 3, Height: 3, Background: ColorBlue}))
		box.AddChild(ts.NewBox(ts.Style{Width: 3, Height: 1, Background: ColorDark2}))
		return box
	}

	root := ts.NewBox(ts.Style{
		Width:         40,
		Height:        18,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		FlexWrap:      ts.Wrap,
		Gap:           2,
	})
	root.AddChild(makeRow(ts.AlignStart))
	root.AddChild(makeRow(ts.AlignCenter))
	root.AddChild(makeRow(ts.AlignEnd))
	root.AddChild(makeRow(ts.AlignStretch))

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}

func TestGolden_Layout_Padding(t *testing.T) {
	// Box with uniform padding
	uniform := ts.NewBox(ts.Style{
		Width:      20,
		Height:     8,
		Display:    ts.Flex,
		Background: ColorSurface,
		Padding:    ts.SpacingAll(2),
		Border:     ts.BorderAllWithColor(ts.BorderSingle, ColorCyan),
	})
	uniform.AddChild(ts.NewBox(ts.Style{Width: 14, Height: 4, Background: ColorBlue}))

	// Box with asymmetric padding
	asymmetric := ts.NewBox(ts.Style{
		Width:      20,
		Height:     8,
		Display:    ts.Flex,
		Background: ColorSurface,
		Padding:    ts.Spacing{Top: 1, Right: 4, Bottom: 1, Left: 2},
		Border:     ts.BorderAllWithColor(ts.BorderSingle, ColorAccent),
	})
	asymmetric.AddChild(ts.NewBox(ts.Style{Width: 12, Height: 4, Background: ColorBlue2}))

	root := ts.NewBox(ts.Style{
		Width:         44,
		Height:        10,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		Gap:           2,
	})
	root.AddChild(uniform)
	root.AddChild(asymmetric)

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}

func TestGolden_Layout_Margin(t *testing.T) {
	root := ts.NewBox(ts.Style{
		Width:         30,
		Height:        12,
		Display:       ts.Flex,
		FlexDirection: ts.Column,
		Background:    ColorSurface,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
	})

	// Child with top margin
	root.AddChild(ts.NewBox(ts.Style{
		Width:      10,
		Height:     2,
		Background: ColorCyan,
		Margin:     ts.Spacing{Top: 1, Left: 2},
	}))

	// Child with all margins
	root.AddChild(ts.NewBox(ts.Style{
		Width:      10,
		Height:     2,
		Background: ColorBlue,
		Margin:     ts.SpacingAll(1),
	}))

	// Child with left margin (auto-like effect)
	root.AddChild(ts.NewBox(ts.Style{
		Width:      10,
		Height:     2,
		Background: ColorDark2,
		Margin:     ts.Spacing{Left: 10},
	}))

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}

func TestGolden_Layout_Nested(t *testing.T) {
	// Header
	header := ts.NewBox(ts.Style{
		Width:          38,
		Height:         3,
		Display:        ts.Flex,
		JustifyContent: ts.JustifyBetween,
		AlignItems:     ts.AlignCenter,
		Background:     ColorSurface,
		Border:         ts.BorderAllWithColor(ts.BorderSingle, ColorAccent),
	})
	header.AddChild(ts.Text("Logo", ts.Style{Foreground: ColorText, FontWeight: ts.WeightBold}))
	header.AddChild(ts.Text("Menu", ts.Style{Foreground: ColorMuted}))

	// Sidebar
	sidebar := ts.NewBox(ts.Style{
		Width:         10,
		Height:        10,
		Display:       ts.Flex,
		FlexDirection: ts.Column,
		Background:    ColorSurface,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Gap:           1,
		Padding:       ts.Spacing{Top: 1, Left: 1},
	})
	sidebar.AddChild(ts.Text("Home", ts.Style{Foreground: ColorCyan}))
	sidebar.AddChild(ts.Text("About", ts.Style{Foreground: ColorText}))
	sidebar.AddChild(ts.Text("Help", ts.Style{Foreground: ColorText}))

	// Content
	content := ts.NewBox(ts.Style{
		Width:      26,
		Height:     10,
		Display:    ts.Flex,
		Background: ColorSurface,
		Border:     ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
		Padding:    ts.SpacingAll(1),
	})
	content.AddChild(ts.Text("Main Content", ts.Style{Foreground: ColorText}))

	// Body (sidebar + content)
	body := ts.NewBox(ts.Style{
		Width:         38,
		Height:        10,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		Gap:           1,
	})
	body.AddChild(sidebar)
	body.AddChild(content)

	// Root container
	root := ts.NewBox(ts.Style{
		Width:         40,
		Height:        14,
		Display:       ts.Flex,
		FlexDirection: ts.Column,
		Gap:           1,
	})
	root.AddChild(header)
	root.AddChild(body)

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}

func TestGolden_Layout_FlexWrap(t *testing.T) {
	root := ts.NewBox(ts.Style{
		Width:         30,
		Height:        12,
		Display:       ts.Flex,
		FlexDirection: ts.Row,
		FlexWrap:      ts.Wrap,
		Gap:           1,
		Border:        ts.BorderAllWithColor(ts.BorderSingle, ColorBorder),
	})

	// Add multiple items that will wrap - all blue tones
	colors := []ts.Color{ColorCyan, ColorCyan2, ColorBlue, ColorBlue2, ColorDark2, ColorAccent}
	for _, c := range colors {
		root.AddChild(ts.NewBox(ts.Style{
			Width:      8,
			Height:     3,
			Background: c,
		}))
	}

	buf := ts.Draw(root)
	assert.Golden(t, ts.Render(buf))
}
