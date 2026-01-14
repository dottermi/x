package render

import (
	"bytes"
	"testing"
)

// BenchmarkTerminalRenderTo measures diff-based rendering performance.
func BenchmarkTerminalRenderTo(b *testing.B) {
	cases := []struct {
		name       string
		changeRate float64
	}{
		{"no_changes", 0.0},
		{"single_cell", -1}, // special case: 1 cell
		{"sparse_1pct", 0.01},
		{"half_50pct", 0.50},
		{"full_100pct", 1.0},
	}

	const width, height = 80, 24

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			term := NewTerminal(width, height)
			buf := NewBuffer(width, height)

			// Apply changes based on rate
			if tc.changeRate == -1 {
				buf.Set(0, 0, Cell{Char: 'X'})
			} else {
				totalCells := width * height
				changeCells := int(float64(totalCells) * tc.changeRate)
				for i := 0; i < changeCells; i++ {
					x := i % width
					y := i / width
					buf.Set(x, y, Cell{Char: 'X', FG: RGB(255, 0, 0)})
				}
			}

			var out bytes.Buffer
			b.ReportAllocs()
			for b.Loop() {
				out.Reset()
				term.RenderTo(buf, &out)
				// Reset terminal state for next iteration
				term.current = NewBuffer(width, height)
			}
		})
	}
}

// BenchmarkTerminalRenderFullTo measures full rendering (no diff).
func BenchmarkTerminalRenderFullTo(b *testing.B) {
	sizes := []struct {
		name string
		w, h int
	}{
		{"10x10", 10, 10},
		{"80x24", 80, 24},
		{"200x50", 200, 50},
	}

	for _, s := range sizes {
		b.Run(s.name, func(b *testing.B) {
			term := NewTerminal(s.w, s.h)
			buf := NewBuffer(s.w, s.h)
			buf.Fill(Cell{Char: '#', FG: RGB(100, 150, 200)})

			var out bytes.Buffer
			b.ReportAllocs()
			for b.Loop() {
				out.Reset()
				term.RenderFullTo(buf, &out)
			}
		})
	}
}

// BenchmarkTerminalRenderStyleChanges measures rendering with frequent style changes.
func BenchmarkTerminalRenderStyleChanges(b *testing.B) {
	const width, height = 80, 24
	term := NewTerminal(width, height)
	buf := NewBuffer(width, height)

	// Alternating colors and styles
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cell := Cell{
				Char:   rune('A' + (x+y)%26),
				FG:     RGB(uint8((x*10)%256), uint8((y*10)%256), 128),
				BG:     RGB(uint8((y*5)%256), uint8((x*5)%256), 64),
				Bold:   (x+y)%2 == 0,
				Italic: (x+y)%3 == 0,
			}
			buf.Set(x, y, cell)
		}
	}

	var out bytes.Buffer
	b.ReportAllocs()
	for b.Loop() {
		out.Reset()
		term.RenderFullTo(buf, &out)
	}
}

// BenchmarkTerminalResize measures resize operation performance.
func BenchmarkTerminalResize(b *testing.B) {
	term := NewTerminal(80, 24)

	b.ReportAllocs()
	for b.Loop() {
		term.Resize(100, 30)
		term.Resize(80, 24)
	}
}

// BenchmarkTerminalClear measures clear operation performance.
func BenchmarkTerminalClear(b *testing.B) {
	term := NewTerminal(80, 24)

	b.ReportAllocs()
	for b.Loop() {
		_ = term.Clear()
	}
}
