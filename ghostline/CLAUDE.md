# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Ghostline is a Go library that provides interactive readline functionality with "ghost text" (inline suggestions) support. It enables terminal applications to show autocomplete suggestions as dimmed text that users can accept with Tab.

## Build Commands

```bash
go build ./...     # Build
go test ./...      # Run tests
go vet ./...       # Static analysis
```

## Architecture

The library consists of a single `Input` struct that manages interactive input with ghost text suggestions:

- **ghostline.go** - Core `Input` struct and main `Readline()` loop that handles raw terminal mode and character-by-character input processing
- **suggestion.go** - `findGhost()` implements prefix-based suggestion matching against the last word typed
- **render.go** - `render()` handles terminal output using ANSI escape codes (dim text for ghost suggestions, cursor positioning)
- **terminal.go** - Raw mode enable/disable using `golang.org/x/term`
- **keys.go** - Key constants (Ctrl+C, Tab, Enter, Backspace, etc.)

## Key Behaviors

- Tab accepts the current ghost suggestion
- Ghost text appears as dimmed text after the cursor
- Suggestions match against the last word (space-delimited)
- Ctrl+C and Ctrl+D abort input
- The library uses dependency injection for `io.Reader`/`io.Writer` to support testing
