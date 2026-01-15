# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
make test-fmt

# Run a single test
go test -run TestName ./path/to/package

# Lint
make lint

# Format code
make fmt

# Build
go build ./...
```

## Architecture

termistyle is a CSS-like styling library for terminal UIs. It uses a render pipeline:

```
Style → Box Tree → Layout (Yoga) → Buffer → ANSI Output
```

### Key Components

**Two Modules**: This workspace contains two Go modules:
- `github.com/dottermi/x/termistyle` - The styling/layout library
- `github.com/dottermi/x/render` - Low-level terminal rendering with double-buffering

**Layout Engine**: Uses `github.com/kjk/flex` (Yoga port) for flexbox calculations. The `layout/yoga.go` file bridges termistyle styles to flex nodes.

**Render Pipeline**:
1. `layout.Box.Calculate()` - Computes X/Y/W/H positions using Yoga
2. `termistyle.Draw()` - Renders box tree to `render.Buffer`
3. `render.Terminal.Render()` - Converts buffer to ANSI escape sequences with diff optimization

**Type Re-exports**: `termistyle.go` is a facade that re-exports types from sub-packages (`style/`, `layout/`, `draw/`) for convenient single-import usage.

### Package Responsibilities

- `style/` - Style definitions: colors, borders, text properties, spacing
- `layout/` - Box tree structure and Yoga integration for layout calculation
- `draw/` - Text and border rendering to buffer
- `x/render/` - Cell/Buffer/Terminal types, ANSI generation, double-buffering

### Color Types

Two color systems exist:
- `style.Color` - Hex string "#RRGGBB" used in style definitions
- `render.Color` - RGB struct used by the renderer

Convert with `style.Color.ToRender()`.

### Buffer Access

The render buffer uses flat row-major storage. Access cells via methods:
```go
buf.Get(x, y)      // returns Cell
buf.Set(x, y, c)   // sets Cell
buf.SetClipped(x, y, c, clip)  // sets with clipping bounds
```

## Terminal Blues Color Palette

Use this palette for examples and tests. Based on "Terminal Blues" by PropFeds.

| Name | Hex | Usage |
|------|-----|-------|
| Background | `#0f0f1b` | Main background |
| Dark Blue | `#1a1a2e` | Dark areas, secondary bg |
| Medium Blue | `#3a4a6a` | Walls, blocks, borders |
| Cyan | `#5aacac` | Water, highlights, links |
| Dark Green | `#2a5a2a` | Dark elements |
| Green | `#4a8a4a` | Success, positive |
| Bright Green | `#6aba6a` | Bright highlights |
| Orange/Gold | `#d4a020` | Player, focus, warnings |
| Red | `#c83737` | Errors, health, danger |
| Pink | `#c85a8a` | UI accents, special |
| Brown | `#8a6a3a` | Earth tones, secondary |
| White | `#c8c8d0` | Primary text |
| Gray | `#6a6a7a` | Secondary text, disabled |

```go
// Example usage
const (
    ColorBg      = style.Color("#0f0f1b")
    ColorDark    = style.Color("#1a1a2e")
    ColorBlue    = style.Color("#3a4a6a")
    ColorCyan    = style.Color("#5aacac")
    ColorGreen   = style.Color("#4a8a4a")
    ColorOrange  = style.Color("#d4a020")
    ColorRed     = style.Color("#c83737")
    ColorPink    = style.Color("#c85a8a")
    ColorText    = style.Color("#c8c8d0")
    ColorMuted   = style.Color("#6a6a7a")
)
```
