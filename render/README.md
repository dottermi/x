# `./render`

Efficient terminal rendering with double-buffering and diff optimization.

## Install

```bash
go get github.com/dottermi/x/render
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/dottermi/x/render"
)

func main() {
    // Create terminal renderer
    term := render.NewTerminal(80, 24)

    // Create buffer and draw
    buf := render.NewBuffer(80, 24)
    buf.Set(0, 0, render.Cell{
        Char: 'H',
        FG:   render.RGB(255, 0, 0),
        Bold: true,
    })
    buf.Set(1, 0, render.Cell{Char: 'i', FG: render.RGB(255, 0, 0)})

    // Render (only changed cells)
    fmt.Print(term.Render(buf))
}
```

## How it works

```
Frame 1:          Frame 2:
┌──────────┐      ┌──────────┐
│ Hello    │      │ Hello    │
│ World    │  →   │ Gopher!  │  → Only "Gopher!" is sent
└──────────┘      └──────────┘
```

Instead of redrawing everything, `Terminal.Render()` compares with the previous frame and outputs only the ANSI codes needed to update changed cells.

## API

### Color

```go
render.RGB(r, g, b uint8) Color  // Create RGB color
render.Default() Color           // Terminal default color
color.IsSet() bool               // Check if explicitly set
color.FGCode() string            // ANSI foreground code
color.BGCode() string            // ANSI background code
```

### Cell

```go
render.EmptyCell() Cell          // Space with default colors
cell.Equal(other Cell) bool      // Compare cells
```

### Buffer

```go
render.NewBuffer(w, h int) *Buffer
buf.Set(x, y int, c Cell)        // Write cell (clips out-of-bounds)
buf.Get(x, y int) Cell           // Read cell
buf.Fill(c Cell)                 // Fill entire buffer
buf.Clone() *Buffer              // Deep copy
buf.Diff(new *Buffer) []CellChange  // Get changed cells
```

### Terminal

```go
render.NewTerminal(w, h int) *Terminal
term.Render(buf *Buffer) string      // Diff-optimized output
term.RenderTo(buf *Buffer, w io.Writer)
term.RenderFull(buf *Buffer) string  // Full render (no diff)
term.Clear() string                  // Clear screen
term.Resize(w, h int)                // Update dimensions
```

### ANSI Helpers

```go
render.MoveCursor(x, y int) string
render.HideCursor() string
render.ShowCursor() string
```

## Benchmark

```
Environment
Apple M4 Pro (14 cores: 10P + 4E)
48 GB RAM
macOS 26.0.1
Go 1.25.5 darwin/arm64
```

| Operation                        | Time    | Allocs |
| -------------------------------- | ------- | ------ |
| `Buffer.Set`                     | 1.9 ns  | 0      |
| `Buffer.Get`                     | 2.3 ns  | 0      |
| `NewBuffer(80x24)`               | 3.2 µs  | 2      |
| `Buffer.Fill(80x24)`             | 747 ns  | 0      |
| `Buffer.Clone(80x24)`            | 2.9 µs  | 2      |
| `Buffer.Diff (no changes)`       | 9.9 µs  | 0      |
| `Buffer.Diff (1% sparse)`        | 9.9 µs  | 6      |
| `Buffer.Diff (50%)`              | 14.7 µs | 11     |
| `Terminal.RenderTo (no changes)` | 12.4 µs | 2      |
| `Terminal.RenderTo (1% sparse)`  | 17.1 µs | 31     |
| `Terminal.RenderTo (50%)`        | 41.2 µs | 977    |
| `Terminal.RenderFullTo(80x24)`   | 40.0 µs | 1924   |
| `Color.RGB`                      | 1.6 ns  | 0      |
| `Color.FGCode (default)`         | 1.7 ns  | 0      |
| `Color.FGCode (set)`             | 22.9 ns | 1      |
| `EmptyCell`                      | 1.7 ns  | 0      |
| `Cell.Equal`                     | 4.2 ns  | 0      |

```bash
go test -bench=. -benchmem ./...
```
