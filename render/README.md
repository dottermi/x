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
