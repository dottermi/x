package border_test

import (
	"testing"

	"github.com/dottermi/termitest/assert"
	"github.com/dottermi/x/render"
	"github.com/dottermi/x/termistyle/draw"
	"github.com/dottermi/x/termistyle/style"
)

func renderBuffer(buf *render.Buffer) string {
	term := render.NewTerminal(buf.Width, buf.Height)
	return term.RenderFull(buf)
}

func TestGolden_Border_Styles(t *testing.T) {
	buf := render.NewBuffer(50, 20)

	// Single
	draw.DrawBorder(buf, 0, 0, 20, 5, style.BorderAll(style.BorderSingle))

	// Double
	draw.DrawBorder(buf, 25, 0, 20, 5, style.BorderAll(style.BorderDouble))

	// Round
	draw.DrawBorder(buf, 0, 7, 20, 5, style.BorderAll(style.BorderRound))

	// Bold
	draw.DrawBorder(buf, 25, 7, 20, 5, style.BorderAll(style.BorderBold))

	// Labels
	draw.DrawText(buf, 2, 2, "Single", style.Color("#888888"), style.Color(""))
	draw.DrawText(buf, 27, 2, "Double", style.Color("#888888"), style.Color(""))
	draw.DrawText(buf, 2, 9, "Round", style.Color("#888888"), style.Color(""))
	draw.DrawText(buf, 27, 9, "Bold", style.Color("#888888"), style.Color(""))

	assert.Golden(t, renderBuffer(buf))
}

func TestGolden_Border_Colors(t *testing.T) {
	buf := render.NewBuffer(50, 14)

	// Red border
	draw.DrawBorder(buf, 0, 0, 20, 5, style.BorderAllWithColor(style.BorderSingle, style.Color("#FF0000")))
	draw.DrawText(buf, 2, 2, "Red", style.Color("#FF0000"), style.Color(""))

	// Green border
	draw.DrawBorder(buf, 25, 0, 20, 5, style.BorderAllWithColor(style.BorderSingle, style.Color("#00FF00")))
	draw.DrawText(buf, 27, 2, "Green", style.Color("#00FF00"), style.Color(""))

	// Blue border
	draw.DrawBorder(buf, 0, 7, 20, 5, style.BorderAllWithColor(style.BorderSingle, style.Color("#0000FF")))
	draw.DrawText(buf, 2, 9, "Blue", style.Color("#0000FF"), style.Color(""))

	// Yellow border
	draw.DrawBorder(buf, 25, 7, 20, 5, style.BorderAllWithColor(style.BorderSingle, style.Color("#FFFF00")))
	draw.DrawText(buf, 27, 9, "Yellow", style.Color("#FFFF00"), style.Color(""))

	assert.Golden(t, renderBuffer(buf))
}

func TestGolden_Border_Titles(t *testing.T) {
	buf := render.NewBuffer(50, 28)

	// Top: Title left
	draw.DrawBorder(buf, 0, 0, 45, 5, style.BorderAllWithTitle(style.BorderSingle, style.Color("#FFFFFF"), "Top Left"))

	// Top: Title center
	border := style.BorderAll(style.BorderSingle)
	border.Top = border.Top.AddText("Top Center", style.TextAlignCenter)
	draw.DrawBorder(buf, 0, 6, 45, 5, border)

	// Top: Title right
	border2 := style.BorderAll(style.BorderSingle)
	border2.Top = border2.Top.AddText("Top Right", style.TextAlignRight)
	draw.DrawBorder(buf, 0, 12, 45, 5, border2)

	// Bottom: Title left
	border3 := style.BorderAll(style.BorderSingle)
	border3.Bottom = border3.Bottom.AddText("Bottom Left", style.TextAlignLeft)
	draw.DrawBorder(buf, 0, 18, 45, 5, border3)

	// Bottom: Title center + right
	border4 := style.BorderAll(style.BorderSingle)
	border4.Bottom = border4.Bottom.AddText("Center", style.TextAlignCenter)
	border4.Bottom = border4.Bottom.AddText("Right", style.TextAlignRight)
	draw.DrawBorder(buf, 0, 24, 45, 4, border4)

	assert.Golden(t, renderBuffer(buf))
}

func TestGolden_Border_Partial(t *testing.T) {
	buf := render.NewBuffer(50, 20)

	// Top only
	border1 := style.Border{Top: style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#FFFFFF")}}
	draw.DrawBorder(buf, 0, 0, 20, 5, border1)
	draw.DrawText(buf, 2, 2, "Top only", style.Color("#888888"), style.Color(""))

	// Left and Right only
	border2 := style.Border{
		Left:  style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#FFFFFF")},
		Right: style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#FFFFFF")},
	}
	draw.DrawBorder(buf, 25, 0, 20, 5, border2)
	draw.DrawText(buf, 27, 2, "Left+Right", style.Color("#888888"), style.Color(""))

	// Top and Bottom only
	border3 := style.Border{
		Top:    style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#FFFFFF")},
		Bottom: style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#FFFFFF")},
	}
	draw.DrawBorder(buf, 0, 7, 20, 5, border3)
	draw.DrawText(buf, 2, 9, "Top+Bottom", style.Color("#888888"), style.Color(""))

	// Three sides (no bottom)
	border4 := style.Border{
		Top:   style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#FFFFFF")},
		Left:  style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#FFFFFF")},
		Right: style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#FFFFFF")},
	}
	draw.DrawBorder(buf, 25, 7, 20, 5, border4)
	draw.DrawText(buf, 27, 9, "No bottom", style.Color("#888888"), style.Color(""))

	// Bottom only
	border5 := style.Border{Bottom: style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#FFFFFF")}}
	draw.DrawBorder(buf, 0, 14, 20, 5, border5)
	draw.DrawText(buf, 2, 16, "Bottom only", style.Color("#888888"), style.Color(""))

	// L-shape (Top + Left)
	border6 := style.Border{
		Top:  style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#FFFFFF")},
		Left: style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#FFFFFF")},
	}
	draw.DrawBorder(buf, 25, 14, 20, 5, border6)
	draw.DrawText(buf, 27, 16, "Top+Left", style.Color("#888888"), style.Color(""))

	assert.Golden(t, renderBuffer(buf))
}

func TestGolden_Border_Complete(t *testing.T) {
	buf := render.NewBuffer(60, 30)

	// Box 1: Round border with title, colored, with styled text inside
	border1 := style.BorderAllWithColor(style.BorderRound, style.Color("#00FFFF"))
	border1.Top = border1.Top.AddText("User Profile", style.TextAlignCenter)
	border1.Bottom = border1.Bottom.AddText("Status: Online", style.TextAlignRight)
	draw.DrawBorder(buf, 0, 0, 28, 8, border1)
	draw.DrawStyledText(buf, 2, 2, "Name: John Doe", style.Style{
		Foreground: style.Color("#FFFFFF"),
		FontWeight: style.WeightBold,
	})
	draw.DrawStyledText(buf, 2, 3, "Email: john@example.com", style.Style{
		Foreground: style.Color("#AAAAAA"),
	})
	draw.DrawStyledText(buf, 2, 5, "Premium Member", style.Style{
		Foreground:     style.Color("#FFD700"),
		FontStyle:      style.StyleItalic,
		TextDecoration: style.DecorationUnderline,
	})

	// Box 2: Double border with multiple titles
	border2 := style.BorderAllWithColor(style.BorderDouble, style.Color("#FF00FF"))
	border2.Top = border2.Top.AddText("Menu", style.TextAlignLeft)
	border2.Top = border2.Top.AddText("[X]", style.TextAlignRight)
	draw.DrawBorder(buf, 30, 0, 28, 8, border2)
	draw.DrawText(buf, 32, 2, "1. Dashboard", style.Color("#FFFFFF"), style.Color(""))
	draw.DrawText(buf, 32, 3, "2. Settings", style.Color("#FFFFFF"), style.Color(""))
	draw.DrawText(buf, 32, 4, "3. Help", style.Color("#FFFFFF"), style.Color(""))
	draw.DrawStyledText(buf, 32, 6, "Select option...", style.Style{
		Foreground: style.Color("#666666"),
		FontStyle:  style.StyleItalic,
	})

	// Box 3: Bold border with colors and text styles
	border3 := style.BorderAllWithColor(style.BorderBold, style.Color("#FF0000"))
	border3.Top = border3.Top.AddText("Warning!", style.TextAlignCenter)
	draw.DrawBorder(buf, 0, 10, 58, 7, border3)
	draw.DrawStyledText(buf, 2, 12, "ALERT:", style.Style{
		Foreground: style.Color("#FF0000"),
		FontWeight: style.WeightBold,
	})
	draw.DrawStyledText(buf, 10, 12, "System resources running low", style.Style{
		Foreground: style.Color("#FFFF00"),
	})
	draw.DrawStyledText(buf, 2, 14, "CPU: 95%", style.Style{Foreground: style.Color("#FF6666")})
	draw.DrawStyledText(buf, 15, 14, "RAM: 87%", style.Style{Foreground: style.Color("#FFAA00")})
	draw.DrawStyledText(buf, 28, 14, "Disk: 45%", style.Style{Foreground: style.Color("#66FF66")})

	// Box 4: Single border with partial sides and text wrap
	border4 := style.Border{
		Top:    style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#00FF00")},
		Bottom: style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#00FF00")},
		Left:   style.BorderEdge{Style: style.BorderSingle, Color: style.Color("#00FF00")},
	}
	border4.Top = border4.Top.AddText("Note", style.TextAlignLeft)
	draw.DrawBorder(buf, 0, 19, 28, 9, border4)
	draw.DrawStyledTextInBox(buf, 2, 21, "This is a longer text that should wrap nicely within the box boundaries.", 24, 5, style.Style{
		Foreground: style.Color("#CCCCCC"),
		TextWrap:   style.WrapWord,
	})

	// Box 5: Mixed styles
	border5 := style.BorderAllWithColor(style.BorderRound, style.Color("#8888FF"))
	border5.Top = border5.Top.AddText("Stats", style.TextAlignCenter)
	border5.Bottom = border5.Bottom.AddText("Updated: Now", style.TextAlignCenter)
	draw.DrawBorder(buf, 30, 19, 28, 9, border5)
	draw.DrawStyledText(buf, 32, 21, "Downloads", style.Style{
		Foreground:    style.Color("#FFFFFF"),
		TextTransform: style.TransformUppercase,
	})
	draw.DrawStyledText(buf, 45, 21, "1,234", style.Style{
		Foreground: style.Color("#00FF00"),
		FontWeight: style.WeightBold,
	})
	draw.DrawStyledText(buf, 32, 23, "Uploads", style.Style{
		Foreground:    style.Color("#FFFFFF"),
		TextTransform: style.TransformUppercase,
	})
	draw.DrawStyledText(buf, 45, 23, "567", style.Style{
		Foreground: style.Color("#00FF00"),
		FontWeight: style.WeightBold,
	})
	draw.DrawStyledText(buf, 32, 25, "Errors", style.Style{
		Foreground:    style.Color("#FFFFFF"),
		TextTransform: style.TransformUppercase,
	})
	draw.DrawStyledText(buf, 45, 25, "3", style.Style{
		Foreground: style.Color("#FF0000"),
		FontWeight: style.WeightBold,
	})

	assert.Golden(t, renderBuffer(buf))
}
