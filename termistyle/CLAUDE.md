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
