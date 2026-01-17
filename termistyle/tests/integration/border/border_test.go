package border_test

import (
	"testing"

	"github.com/dottermi/termitest/assert"
	"github.com/dottermi/x/render"
	"github.com/dottermi/x/termistyle/draw"
	"github.com/dottermi/x/termistyle/style"
	. "github.com/dottermi/x/termistyle/tests/integration/helper"
)

func renderBuffer(buf *render.Buffer) string {
	term := render.NewTerminal(buf.Width, buf.Height)
	return term.RenderFull(buf)
}

func TestGolden_Border_Styles(t *testing.T) {
	buf := render.NewBuffer(50, 20)

	// Single
	draw.DrawBorder(buf, 0, 0, 20, 5, style.BorderAllWithColor(style.BorderSingle, ColorBorder))
	draw.DrawText(buf, 2, 2, "Single", ColorMuted, ColorBg)

	// Double
	draw.DrawBorder(buf, 25, 0, 20, 5, style.BorderAllWithColor(style.BorderDouble, ColorBorder))
	draw.DrawText(buf, 27, 2, "Double", ColorMuted, ColorBg)

	// Round
	draw.DrawBorder(buf, 0, 7, 20, 5, style.BorderAllWithColor(style.BorderRound, ColorBorder))
	draw.DrawText(buf, 2, 9, "Round", ColorMuted, ColorBg)

	// Bold
	draw.DrawBorder(buf, 25, 7, 20, 5, style.BorderAllWithColor(style.BorderBold, ColorBorder))
	draw.DrawText(buf, 27, 9, "Bold", ColorMuted, ColorBg)

	assert.Golden(t, renderBuffer(buf))
}

func TestGolden_Border_Colors(t *testing.T) {
	buf := render.NewBuffer(50, 14)

	// Cyan border
	draw.DrawBorder(buf, 0, 0, 20, 5, style.BorderAllWithColor(style.BorderSingle, ColorCyan))
	draw.DrawText(buf, 2, 2, "Cyan", ColorCyan, ColorBg)

	// Cyan2 border
	draw.DrawBorder(buf, 25, 0, 20, 5, style.BorderAllWithColor(style.BorderSingle, ColorCyan2))
	draw.DrawText(buf, 27, 2, "Cyan2", ColorCyan2, ColorBg)

	// Blue border
	draw.DrawBorder(buf, 0, 7, 20, 5, style.BorderAllWithColor(style.BorderSingle, ColorBlue))
	draw.DrawText(buf, 2, 9, "Blue", ColorBlue, ColorBg)

	// Blue2 border
	draw.DrawBorder(buf, 25, 7, 20, 5, style.BorderAllWithColor(style.BorderSingle, ColorBlue2))
	draw.DrawText(buf, 27, 9, "Blue2", ColorBlue2, ColorBg)

	assert.Golden(t, renderBuffer(buf))
}

func TestGolden_Border_Titles(t *testing.T) {
	buf := render.NewBuffer(50, 28)

	// Top: Title left
	draw.DrawBorder(buf, 0, 0, 45, 5, style.BorderAllWithTitle(style.BorderSingle, ColorAccent, "Top Left"))

	// Top: Title center
	border := style.BorderAllWithColor(style.BorderSingle, ColorAccent)
	border.Top = border.Top.AddText("Top Center", style.TextAlignCenter)
	draw.DrawBorder(buf, 0, 6, 45, 5, border)

	// Top: Title right
	border2 := style.BorderAllWithColor(style.BorderSingle, ColorAccent)
	border2.Top = border2.Top.AddText("Top Right", style.TextAlignRight)
	draw.DrawBorder(buf, 0, 12, 45, 5, border2)

	// Bottom: Title left
	border3 := style.BorderAllWithColor(style.BorderSingle, ColorAccent)
	border3.Bottom = border3.Bottom.AddText("Bottom Left", style.TextAlignLeft)
	draw.DrawBorder(buf, 0, 18, 45, 5, border3)

	// Bottom: Title center + right
	border4 := style.BorderAllWithColor(style.BorderSingle, ColorAccent)
	border4.Bottom = border4.Bottom.AddText("Center", style.TextAlignCenter)
	border4.Bottom = border4.Bottom.AddText("Right", style.TextAlignRight)
	draw.DrawBorder(buf, 0, 24, 45, 4, border4)

	assert.Golden(t, renderBuffer(buf))
}

func TestGolden_Border_Partial(t *testing.T) {
	buf := render.NewBuffer(50, 20)

	// Top only
	border1 := style.Border{Top: style.BorderEdge{Style: style.BorderSingle, Color: ColorBorder}}
	draw.DrawBorder(buf, 0, 0, 20, 5, border1)
	draw.DrawText(buf, 2, 2, "Top only", ColorMuted, ColorBg)

	// Left and Right only
	border2 := style.Border{
		Left:  style.BorderEdge{Style: style.BorderSingle, Color: ColorBorder},
		Right: style.BorderEdge{Style: style.BorderSingle, Color: ColorBorder},
	}
	draw.DrawBorder(buf, 25, 0, 20, 5, border2)
	draw.DrawText(buf, 27, 2, "Left+Right", ColorMuted, ColorBg)

	// Top and Bottom only
	border3 := style.Border{
		Top:    style.BorderEdge{Style: style.BorderSingle, Color: ColorBorder},
		Bottom: style.BorderEdge{Style: style.BorderSingle, Color: ColorBorder},
	}
	draw.DrawBorder(buf, 0, 7, 20, 5, border3)
	draw.DrawText(buf, 2, 9, "Top+Bottom", ColorMuted, ColorBg)

	// Three sides (no bottom)
	border4 := style.Border{
		Top:   style.BorderEdge{Style: style.BorderSingle, Color: ColorBorder},
		Left:  style.BorderEdge{Style: style.BorderSingle, Color: ColorBorder},
		Right: style.BorderEdge{Style: style.BorderSingle, Color: ColorBorder},
	}
	draw.DrawBorder(buf, 25, 7, 20, 5, border4)
	draw.DrawText(buf, 27, 9, "No bottom", ColorMuted, ColorBg)

	// Bottom only
	border5 := style.Border{Bottom: style.BorderEdge{Style: style.BorderSingle, Color: ColorBorder}}
	draw.DrawBorder(buf, 0, 14, 20, 5, border5)
	draw.DrawText(buf, 2, 16, "Bottom only", ColorMuted, ColorBg)

	// L-shape (Top + Left)
	border6 := style.Border{
		Top:  style.BorderEdge{Style: style.BorderSingle, Color: ColorBorder},
		Left: style.BorderEdge{Style: style.BorderSingle, Color: ColorBorder},
	}
	draw.DrawBorder(buf, 25, 14, 20, 5, border6)
	draw.DrawText(buf, 27, 16, "Top+Left", ColorMuted, ColorBg)

	assert.Golden(t, renderBuffer(buf))
}

func TestGolden_Border_Complete(t *testing.T) {
	buf := render.NewBuffer(60, 30)

	// Box 1: User Profile - Round accent border
	border1 := style.BorderAllWithColor(style.BorderRound, ColorAccent)
	border1.Top = border1.Top.AddText("User Profile", style.TextAlignCenter)
	border1.Bottom = border1.Bottom.AddText("Status: Online", style.TextAlignRight)
	draw.DrawBorder(buf, 0, 0, 28, 8, border1)
	draw.DrawStyledText(buf, 2, 2, "Name: John Doe", style.Style{
		Foreground: ColorText,
		FontWeight: style.WeightBold,
	})
	draw.DrawStyledText(buf, 2, 3, "Email: john@example.com", style.Style{
		Foreground: ColorMuted,
	})
	draw.DrawStyledText(buf, 2, 5, "Premium Member", style.Style{
		Foreground:     ColorCyan2,
		FontStyle:      style.StyleItalic,
		TextDecoration: style.DecorationUnderline,
	})

	// Box 2: Menu - Double cyan border
	border2 := style.BorderAllWithColor(style.BorderDouble, ColorCyan)
	border2.Top = border2.Top.AddText("Menu", style.TextAlignLeft)
	border2.Top = border2.Top.AddText("[X]", style.TextAlignRight)
	draw.DrawBorder(buf, 30, 0, 28, 8, border2)
	draw.DrawText(buf, 32, 2, "1. Dashboard", ColorText, ColorBg)
	draw.DrawText(buf, 32, 3, "2. Settings", ColorText, ColorBg)
	draw.DrawText(buf, 32, 4, "3. Help", ColorText, ColorBg)
	draw.DrawStyledText(buf, 32, 6, "Select option...", style.Style{
		Foreground: ColorMuted,
		FontStyle:  style.StyleItalic,
	})

	// Box 3: Notice - Bold blue2 border
	border3 := style.BorderAllWithColor(style.BorderBold, ColorBlue2)
	border3.Top = border3.Top.AddText("Notice", style.TextAlignCenter)
	draw.DrawBorder(buf, 0, 10, 58, 7, border3)
	draw.DrawStyledText(buf, 2, 12, "INFO:", style.Style{
		Foreground: ColorCyan,
		FontWeight: style.WeightBold,
	})
	draw.DrawStyledText(buf, 9, 12, "System status check complete", style.Style{
		Foreground: ColorText,
	})
	draw.DrawStyledText(buf, 2, 14, "CPU: 45%", style.Style{Foreground: ColorCyan2})
	draw.DrawStyledText(buf, 15, 14, "RAM: 67%", style.Style{Foreground: ColorCyan})
	draw.DrawStyledText(buf, 28, 14, "Disk: 32%", style.Style{Foreground: ColorAccent})

	// Box 4: Note - Partial border (no right)
	border4 := style.Border{
		Top:    style.BorderEdge{Style: style.BorderSingle, Color: ColorBlue},
		Bottom: style.BorderEdge{Style: style.BorderSingle, Color: ColorBlue},
		Left:   style.BorderEdge{Style: style.BorderSingle, Color: ColorBlue},
	}
	border4.Top = border4.Top.AddText("Note", style.TextAlignLeft)
	draw.DrawBorder(buf, 0, 19, 28, 9, border4)
	draw.DrawStyledTextInBox(buf, 2, 21, "This is a longer text that should wrap nicely within the box boundaries.", 24, 5, style.Style{
		Foreground: ColorMuted,
		TextWrap:   style.WrapWord,
	})

	// Box 5: Stats - Round border
	border5 := style.BorderAllWithColor(style.BorderRound, ColorBorder)
	border5.Top = border5.Top.AddText("Stats", style.TextAlignCenter)
	border5.Bottom = border5.Bottom.AddText("Updated: Now", style.TextAlignCenter)
	draw.DrawBorder(buf, 30, 19, 28, 9, border5)
	draw.DrawStyledText(buf, 32, 21, "Downloads", style.Style{
		Foreground:    ColorText,
		TextTransform: style.TransformUppercase,
	})
	draw.DrawStyledText(buf, 45, 21, "1,234", style.Style{
		Foreground: ColorCyan2,
		FontWeight: style.WeightBold,
	})
	draw.DrawStyledText(buf, 32, 23, "Uploads", style.Style{
		Foreground:    ColorText,
		TextTransform: style.TransformUppercase,
	})
	draw.DrawStyledText(buf, 45, 23, "567", style.Style{
		Foreground: ColorCyan2,
		FontWeight: style.WeightBold,
	})
	draw.DrawStyledText(buf, 32, 25, "Errors", style.Style{
		Foreground:    ColorText,
		TextTransform: style.TransformUppercase,
	})
	draw.DrawStyledText(buf, 45, 25, "3", style.Style{
		Foreground: ColorCyan,
		FontWeight: style.WeightBold,
	})

	assert.Golden(t, renderBuffer(buf))
}
