# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Ghostline is a minimal API for inline suggestions in Go. It provides readline functionality with ghost text — dimmed autocomplete suggestions that users accept with Tab.

## Build Commands

```bash
go build ./...     # Build
go test ./...      # Run tests
go vet ./...       # Static analysis
```

## Architecture

The library consists of a single `Input` struct that manages interactive input:

- **ghostline.go** - Core `Input` struct and `Readline()` loop with raw terminal mode
- **handlers.go** - Key handlers for navigation, editing, and history
- **history.go** - Command history with up/down navigation
- **suggestion.go** - `findGhost()` prefix matching against last word
- **render.go** - ANSI escape codes for multiline rendering and cursor positioning
- **terminal.go** - Raw mode using `golang.org/x/term`
- **errors.go** - `ErrInterrupted` (Ctrl+C) and `ErrEOF` (Ctrl+D)

## Key Behaviors

- Tab accepts ghost suggestion
- Enter submits, Ctrl+J adds newline
- Arrow keys for cursor movement and history
- Emacs keybindings (Ctrl+A/E/K/U/W)
- Dependency injection for `io.Reader`/`io.Writer` (testing)
