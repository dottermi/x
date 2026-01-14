package style

import (
	"strconv"
	"strings"

	"github.com/dottermi/x/render"
)

// Color represents a 24-bit RGB color as a hex string.
// Use format "#RRGGBB" or "RRGGBB".
//
// Example:
//
//	red := style.Color("#FF0000")
//	blue := style.Color("#0000FF")
type Color string

// IsSet returns true if the color is defined.
func (c Color) IsSet() bool {
	return c != ""
}

// hex returns the numeric value of the color.
func (c Color) hex() uint32 {
	s := strings.TrimPrefix(string(c), "#")
	val, _ := strconv.ParseUint(s, 16, 32)
	return uint32(val)
}

// R returns the red component (0-255).
func (c Color) R() uint8 {
	return uint8((c.hex() >> 16) & 0xFF)
}

// G returns the green component (0-255).
func (c Color) G() uint8 {
	return uint8((c.hex() >> 8) & 0xFF)
}

// B returns the blue component (0-255).
func (c Color) B() uint8 {
	return uint8(c.hex() & 0xFF)
}

// ToRender converts a style.Color to a render.Color.
// Returns render.Default() if the color is not set.
func (c Color) ToRender() render.Color {
	if !c.IsSet() {
		return render.Default()
	}
	return render.RGB(c.R(), c.G(), c.B())
}
