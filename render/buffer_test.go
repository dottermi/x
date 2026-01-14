package render

import "testing"

// BenchmarkBufferClone measures buffer cloning performance across different sizes.
func BenchmarkBufferClone(b *testing.B) {
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
			buf := NewBuffer(s.w, s.h)
			buf.Fill(Cell{Char: 'X', FG: RGB(255, 0, 0)})
			b.ReportAllocs()
			for b.Loop() {
				_ = buf.Clone()
			}
		})
	}
}

// BenchmarkBufferDiff measures diff computation with varying change rates.
func BenchmarkBufferDiff(b *testing.B) {
	cases := []struct {
		name       string
		changeRate float64
	}{
		{"no_changes", 0.0},
		{"sparse_1pct", 0.01},
		{"half_50pct", 0.50},
		{"full_100pct", 1.0},
	}

	const width, height = 80, 24

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			old := NewBuffer(width, height)
			next := old.Clone()

			// Apply changes based on rate
			totalCells := width * height
			changeCells := int(float64(totalCells) * tc.changeRate)
			for i := 0; i < changeCells; i++ {
				x := i % width
				y := i / width
				next.Set(x, y, Cell{Char: 'X'})
			}

			b.ReportAllocs()
			for b.Loop() {
				_ = old.Diff(next)
			}
		})
	}
}

// BenchmarkBufferFill measures buffer fill performance across different sizes.
func BenchmarkBufferFill(b *testing.B) {
	sizes := []struct {
		name string
		w, h int
	}{
		{"10x10", 10, 10},
		{"80x24", 80, 24},
		{"200x50", 200, 50},
	}

	cell := Cell{Char: '#', FG: RGB(100, 150, 200), BG: RGB(10, 20, 30)}

	for _, s := range sizes {
		b.Run(s.name, func(b *testing.B) {
			buf := NewBuffer(s.w, s.h)
			b.ReportAllocs()
			for b.Loop() {
				buf.Fill(cell)
			}
		})
	}
}

// BenchmarkBufferSetGet measures individual cell operations.
func BenchmarkBufferSetGet(b *testing.B) {
	buf := NewBuffer(80, 24)
	cell := Cell{Char: 'A', FG: RGB(255, 255, 255)}

	b.Run("Set", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			buf.Set(40, 12, cell)
		}
	})

	b.Run("Get", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = buf.Get(40, 12)
		}
	})

	b.Run("Set_OutOfBounds", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			buf.Set(100, 100, cell)
		}
	})

	b.Run("Get_OutOfBounds", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = buf.Get(100, 100)
		}
	})
}

// BenchmarkNewBuffer measures buffer allocation performance.
func BenchmarkNewBuffer(b *testing.B) {
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
			b.ReportAllocs()
			for b.Loop() {
				_ = NewBuffer(s.w, s.h)
			}
		})
	}
}
