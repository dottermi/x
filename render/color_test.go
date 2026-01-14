package render

import "testing"

// BenchmarkColorFGCode measures foreground ANSI code generation.
func BenchmarkColorFGCode(b *testing.B) {
	b.Run("Set", func(b *testing.B) {
		c := RGB(128, 64, 255)
		b.ReportAllocs()
		for b.Loop() {
			_ = c.FGCode()
		}
	})

	b.Run("Default", func(b *testing.B) {
		c := Default()
		b.ReportAllocs()
		for b.Loop() {
			_ = c.FGCode()
		}
	})
}

// BenchmarkColorBGCode measures background ANSI code generation.
func BenchmarkColorBGCode(b *testing.B) {
	b.Run("Set", func(b *testing.B) {
		c := RGB(128, 64, 255)
		b.ReportAllocs()
		for b.Loop() {
			_ = c.BGCode()
		}
	})

	b.Run("Default", func(b *testing.B) {
		c := Default()
		b.ReportAllocs()
		for b.Loop() {
			_ = c.BGCode()
		}
	})
}

// BenchmarkColorRGB measures RGB color creation.
func BenchmarkColorRGB(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = RGB(128, 64, 255)
	}
}

// BenchmarkCellEqual measures cell equality comparison.
func BenchmarkCellEqual(b *testing.B) {
	b.Run("Equal", func(b *testing.B) {
		c1 := Cell{Char: 'A', FG: RGB(255, 0, 0), Bold: true}
		c2 := Cell{Char: 'A', FG: RGB(255, 0, 0), Bold: true}
		b.ReportAllocs()
		for b.Loop() {
			_ = c1.Equal(c2)
		}
	})

	b.Run("Different_Char", func(b *testing.B) {
		c1 := Cell{Char: 'A', FG: RGB(255, 0, 0)}
		c2 := Cell{Char: 'B', FG: RGB(255, 0, 0)}
		b.ReportAllocs()
		for b.Loop() {
			_ = c1.Equal(c2)
		}
	})

	b.Run("Different_Color", func(b *testing.B) {
		c1 := Cell{Char: 'A', FG: RGB(255, 0, 0)}
		c2 := Cell{Char: 'A', FG: RGB(0, 255, 0)}
		b.ReportAllocs()
		for b.Loop() {
			_ = c1.Equal(c2)
		}
	})

	b.Run("Different_Attrs", func(b *testing.B) {
		c1 := Cell{Char: 'A', Bold: true}
		c2 := Cell{Char: 'A', Italic: true}
		b.ReportAllocs()
		for b.Loop() {
			_ = c1.Equal(c2)
		}
	})
}

// BenchmarkEmptyCell measures empty cell creation.
func BenchmarkEmptyCell(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = EmptyCell()
	}
}
