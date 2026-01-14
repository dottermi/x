# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
go build ./...              # Build all packages
go test ./...               # Run all tests
go test -run TestBuffer     # Run specific test
go test -v ./...            # Verbose test output
```

## Architecture

This is a terminal rendering library that implements double-buffering with diff optimization, inspired by Ratatui's immediate-mode pattern.

### Core Types

- **Color** (`color.go`): RGB colors with `Set` flag to distinguish explicit colors from terminal defaults
- **Cell** (`cell.go`): Single character position with char, FG/BG colors, and text attributes (bold, italic, etc.)
- **Buffer** (`buffer.go`): 2D grid of cells in row-major order (`index = y*Width + x`)
- **Terminal** (`terminal.go`): Double-buffered renderer that tracks previous frame and outputs only changed cells

### Rendering Flow

1. Create `Terminal` with dimensions
2. Each frame: create new `Buffer`, draw cells into it
3. Call `term.Render(buf)` which diffs against previous frame
4. Only changed cells produce ANSI output

### Key Design Decisions

- `Color.Set` field distinguishes `RGB(0,0,0)` (black) from `Default()` (terminal default)
- Out-of-bounds `Buffer.Set()` silently clips; `Get()` returns `EmptyCell()`
- `Terminal.RenderTo()` writes directly to `io.Writer` to avoid string allocation
- Style tracking in `cellStyle` struct minimizes redundant ANSI codes between adjacent cells

## Plan Mode

- Make the plan extremely concise. Sacrifice grammar for the sake of concision.
- At the end of each plan, give me a list of unresolved questions to answer, if any.
