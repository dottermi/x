# termistyle

A CSS-like styling library for terminal UIs in Go. Build beautiful terminal interfaces with flexbox layouts, 24-bit RGB colors, borders, and text styling.

## Features

- **Flexbox Layout** - Row/column direction, justify-content, align-items, gap, flex-grow/shrink, flex-wrap
- **24-bit RGB Colors** - Full color support for foreground and background
- **Borders** - Single, double, rounded, bold styles with per-side color control and titles
- **Text Styling** - Bold, italic, underline, strikethrough, alignment, wrapping, transforms
- **Auto-sizing** - Dimensions calculated from content when not specified
- **Absolute Positioning** - Position elements relative to container
- **Overflow Control** - Visible or hidden with clipping

## Installation

```bash
go get github.com/dottermi/x/termistyle
```

## Quick Start

```go
package main

import (
    "os"

    ts "github.com/dottermi/x/termistyle"
)

func main() {
    container := ts.NewBox(ts.Style{
        Width:      40,
        Height:     10,
        Display:    ts.Flex,
        Direction:  ts.Row,
        Justify:    ts.JustifyCenter,
        Align:      ts.AlignCenter,
        Background: ts.Color("#1a1a1a"),
        Border:     ts.BorderAllWithColor(ts.BorderRound, ts.Color("#FFFFFF")),
    })

    text := ts.Text("Hello, Terminal!", ts.Style{
        Foreground: ts.Color("#00FF00"),
        FontWeight: ts.WeightBold,
    })

    container.AddChild(text)
    ts.Println(os.Stdout, container)
}
```

## API Overview

### Creating Elements

```go
// Create a box with style
box := ts.NewBox(ts.Style{...})

// Create a text element (auto-calculates dimensions)
text := ts.Text("content", ts.Style{...})

// Add children to a box
box.AddChild(child)
```

### Rendering

```go
// Draw to buffer
buffer := ts.Draw(box)

// Render buffer to string
output := ts.Render(buffer)

// Render directly to writer
ts.RenderTo(buffer, os.Stdout)

// All-in-one: draw and print
ts.Print(os.Stdout, box)
ts.Println(os.Stdout, box)  // with newline
```

### Style Properties

```go
ts.Style{
    // Dimensions
    Width:     40,
    Height:    10,
    MinWidth:  20,
    MaxWidth:  100,
    MinHeight: 5,
    MaxHeight: 50,

    // Layout
    Display:   ts.Flex,       // Flex, Block
    Direction: ts.Row,        // Row, Column
    Justify:   ts.JustifyCenter,  // JustifyStart, JustifyCenter, JustifyEnd, JustifyBetween, JustifyAround
    Align:     ts.AlignCenter,    // AlignStart, AlignCenter, AlignEnd, AlignStretch
    FlexWrap:  ts.Wrap,       // NoWrap, Wrap
    Gap:       1,             // Gap between flex items
    FlexGrow:  1,             // Grow factor
    FlexShrink: 1,            // Shrink factor
    FlexBasis: 20,            // Initial size

    // Spacing
    Padding: ts.SpacingAll(2),    // or ts.SpacingXY(2, 1) or ts.Spacing{Top: 1, Right: 2, Bottom: 1, Left: 2}
    Margin:  ts.SpacingAll(1),

    // Colors
    Foreground: ts.Color("#FFFFFF"),
    Background: ts.Color("#1a1a1a"),

    // Border
    Border: ts.BorderAllWithColor(ts.BorderRound, ts.Color("#FF0000")),

    // Text
    TextAlign:      ts.TextAlignCenter,  // TextAlignLeft, TextAlignCenter, TextAlignRight
    TextWrap:       ts.WrapWord,         // WrapNone, WrapWord, WrapChar
    TextOverflow:   ts.TextOverflowEllipsis,  // TextOverflowClip, TextOverflowEllipsis
    FontWeight:     ts.WeightBold,       // WeightNormal, WeightBold
    FontStyle:      ts.StyleItalic,      // StyleNormal, StyleItalic, StyleOblique
    TextDecoration: ts.DecorationUnderline,   // DecorationNone, DecorationUnderline, DecorationLineThrough
    TextTransform:  ts.TransformUppercase,    // TransformNone, TransformUppercase, TransformLowercase, TransformCapitalize

    // Positioning
    Position: ts.Absolute,    // Relative, Absolute
    Top:      5,
    Left:     10,
    ZIndex:   1,

    // Overflow
    Overflow: ts.OverflowHidden,  // OverflowVisible, OverflowHidden
}
```

### Border Helpers

```go
// Same style on all sides
ts.BorderAll(ts.BorderSingle)

// Same style and color on all sides
ts.BorderAllWithColor(ts.BorderRound, ts.Color("#FF0000"))

// With title on top edge
ts.BorderAllWithTitle(ts.BorderDouble, ts.Color("#00FF00"), "Title")

// Different horizontal/vertical styles
ts.BorderXY(horizontalEdge, verticalEdge)

// Explicit per-side configuration
ts.BorderTRBL(top, right, bottom, left)
```

### Border with Text

```go
// Create a border edge with text on the left
topEdge := ts.BorderEdge{
    Style: ts.BorderRound,
    Color: ts.Color("#FFFFFF"),
}.AddText("Title", ts.TextAlignLeft)

// Multiple texts on the same edge (left, center, right)
topEdge = ts.BorderEdge{
    Style: ts.BorderSingle,
    Color: ts.Color("#888888"),
}.AddText("Left", ts.TextAlignLeft).
  AddText("Center", ts.TextAlignCenter).
  AddText("Right", ts.TextAlignRight)

// Text with custom styling
topEdge = ts.BorderEdge{
    Style: ts.BorderDouble,
    Color: ts.Color("#00FF00"),
}.AddTextWithStyle("Status: OK", ts.TextAlignLeft, ts.Style{
    Foreground: ts.Color("#00FF00"),
    FontWeight: ts.WeightBold,
})

// Use in a box
box := ts.NewBox(ts.Style{
    Width:  40,
    Height: 10,
    Border: ts.Border{
        Top:    topEdge,
        Right:  ts.BorderEdge{Style: ts.BorderRound, Color: ts.Color("#FFFFFF")},
        Bottom: ts.BorderEdge{Style: ts.BorderRound, Color: ts.Color("#FFFFFF")},
        Left:   ts.BorderEdge{Style: ts.BorderRound, Color: ts.Color("#FFFFFF")},
    },
})
```

### Hiding Borders

```go
// Hide specific sides using BorderNone
box := ts.NewBox(ts.Style{
    Width:  30,
    Height: 5,
    Border: ts.Border{
        Top:    ts.BorderEdge{Style: ts.BorderSingle, Color: ts.Color("#FFFFFF")},
        Right:  ts.BorderEdge{Style: ts.BorderSingle, Color: ts.Color("#FFFFFF")},
        Bottom: ts.BorderEdge{Style: ts.BorderNone}, // hidden
        Left:   ts.BorderEdge{Style: ts.BorderNone}, // hidden
    },
})

// Or use BorderXY to hide horizontal/vertical sides
box = ts.NewBox(ts.Style{
    Width:  30,
    Height: 5,
    Border: ts.BorderXY(
        ts.BorderEdge{Style: ts.BorderSingle, Color: ts.Color("#FFFFFF")}, // top/bottom
        ts.BorderEdge{Style: ts.BorderNone}, // left/right hidden
    ),
})
```

### Border Styles

- `ts.BorderNone` - No border (hidden)
- `ts.BorderSingle` - Single line (`─ │ ┐ ┘ └ ┌`)
- `ts.BorderDouble` - Double line (`═ ║ ╗ ╝ ╚ ╔`)
- `ts.BorderRound` - Rounded corners (`─ │ ╮ ╯ ╰ ╭`)

## Architecture

The library follows a render pipeline:

```
Style → Box Tree → Layout → Buffer → ANSI Output
```

1. Create box tree with `NewBox()` and `AddChild()`
2. `Draw()` computes layout positions (X/Y/W/H)
3. `Render()` converts buffer to ANSI escape sequences

### Package Structure

- `termistyle.go` - Public API facade
- `style/` - Style definitions, colors, borders, text properties
- `layout/` - Box tree and layout calculations (block, flex)
- `draw/` - 2D buffer, text rendering, border drawing
- `render/` - ANSI escape sequence generation

## Development

```bash
make lint          # Run golangci-lint
make test          # Run all tests
make test-fmt      # Run tests with formatted output and coverage
make fmt           # Format code
```

## License

MIT
