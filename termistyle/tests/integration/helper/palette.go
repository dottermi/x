// Package helper provides shared utilities for integration tests.
package helper

import "github.com/dottermi/x/termistyle/style"

// Terminal Blues palette by PropFeds - Blue tones only
const (
	ColorBg      = style.Color("#0f0f1b") // Darkest background
	ColorDark    = style.Color("#1a1a2e") // Dark blue
	ColorDark2   = style.Color("#252542") // Slightly lighter
	ColorBlue    = style.Color("#3a4a6a") // Medium blue
	ColorBlue2   = style.Color("#4a5a7a") // Lighter medium blue
	ColorCyan    = style.Color("#5aacac") // Cyan highlight
	ColorCyan2   = style.Color("#7acaca") // Light cyan
	ColorText    = style.Color("#c8c8d0") // Primary text
	ColorMuted   = style.Color("#6a6a7a") // Muted text
	ColorBorder  = style.Color("#3a4a6a") // Border color
	ColorAccent  = style.Color("#5aacac") // Accent color
	ColorSurface = style.Color("#1a1a2e") // Surface/card background
)
