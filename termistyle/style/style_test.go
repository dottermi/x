package style

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColor_IsSet(t *testing.T) {
	t.Parallel()

	t.Run("should return true when color has value with hash prefix", func(t *testing.T) {
		t.Parallel()

		color := Color("#FF0000")

		result := color.IsSet()

		assert.True(t, result)
	})

	t.Run("should return true when color has value without hash prefix", func(t *testing.T) {
		t.Parallel()

		color := Color("00FF00")

		result := color.IsSet()

		assert.True(t, result)
	})

	t.Run("should return false when color is empty string", func(t *testing.T) {
		t.Parallel()

		color := Color("")

		result := color.IsSet()

		assert.False(t, result)
	})
}

func TestColor_R(t *testing.T) {
	t.Parallel()

	t.Run("should return red component from hex color with hash prefix", func(t *testing.T) {
		t.Parallel()

		color := Color("#FF0000")

		result := color.R()

		assert.Equal(t, uint8(255), result)
	})

	t.Run("should return red component from hex color without hash prefix", func(t *testing.T) {
		t.Parallel()

		color := Color("AA5500")

		result := color.R()

		assert.Equal(t, uint8(170), result)
	})

	t.Run("should return zero when red component is zero", func(t *testing.T) {
		t.Parallel()

		color := Color("#00FFFF")

		result := color.R()

		assert.Equal(t, uint8(0), result)
	})

	t.Run("should return correct value for mixed color", func(t *testing.T) {
		t.Parallel()

		color := Color("#123456")

		result := color.R()

		assert.Equal(t, uint8(0x12), result)
	})
}

func TestColor_G(t *testing.T) {
	t.Parallel()

	t.Run("should return green component from hex color with hash prefix", func(t *testing.T) {
		t.Parallel()

		color := Color("#00FF00")

		result := color.G()

		assert.Equal(t, uint8(255), result)
	})

	t.Run("should return green component from hex color without hash prefix", func(t *testing.T) {
		t.Parallel()

		color := Color("00AA00")

		result := color.G()

		assert.Equal(t, uint8(170), result)
	})

	t.Run("should return zero when green component is zero", func(t *testing.T) {
		t.Parallel()

		color := Color("#FF00FF")

		result := color.G()

		assert.Equal(t, uint8(0), result)
	})

	t.Run("should return correct value for mixed color", func(t *testing.T) {
		t.Parallel()

		color := Color("#123456")

		result := color.G()

		assert.Equal(t, uint8(0x34), result)
	})
}

func TestColor_B(t *testing.T) {
	t.Parallel()

	t.Run("should return blue component from hex color with hash prefix", func(t *testing.T) {
		t.Parallel()

		color := Color("#0000FF")

		result := color.B()

		assert.Equal(t, uint8(255), result)
	})

	t.Run("should return blue component from hex color without hash prefix", func(t *testing.T) {
		t.Parallel()

		color := Color("0000AA")

		result := color.B()

		assert.Equal(t, uint8(170), result)
	})

	t.Run("should return zero when blue component is zero", func(t *testing.T) {
		t.Parallel()

		color := Color("#FFFF00")

		result := color.B()

		assert.Equal(t, uint8(0), result)
	})

	t.Run("should return correct value for mixed color", func(t *testing.T) {
		t.Parallel()

		color := Color("#123456")

		result := color.B()

		assert.Equal(t, uint8(0x56), result)
	})
}

func TestSpacingAll(t *testing.T) {
	t.Parallel()

	t.Run("should create uniform spacing with positive value", func(t *testing.T) {
		t.Parallel()

		spacing := SpacingAll(5)

		assert.Equal(t, 5, spacing.Top)
		assert.Equal(t, 5, spacing.Right)
		assert.Equal(t, 5, spacing.Bottom)
		assert.Equal(t, 5, spacing.Left)
	})

	t.Run("should create uniform spacing with zero value", func(t *testing.T) {
		t.Parallel()

		spacing := SpacingAll(0)

		assert.Equal(t, 0, spacing.Top)
		assert.Equal(t, 0, spacing.Right)
		assert.Equal(t, 0, spacing.Bottom)
		assert.Equal(t, 0, spacing.Left)
	})

	t.Run("should create uniform spacing with large value", func(t *testing.T) {
		t.Parallel()

		spacing := SpacingAll(100)

		assert.Equal(t, 100, spacing.Top)
		assert.Equal(t, 100, spacing.Right)
		assert.Equal(t, 100, spacing.Bottom)
		assert.Equal(t, 100, spacing.Left)
	})
}

func TestSpacingXY(t *testing.T) {
	t.Parallel()

	t.Run("should create spacing with different horizontal and vertical values", func(t *testing.T) {
		t.Parallel()

		spacing := SpacingXY(4, 2)

		assert.Equal(t, 2, spacing.Top)
		assert.Equal(t, 4, spacing.Right)
		assert.Equal(t, 2, spacing.Bottom)
		assert.Equal(t, 4, spacing.Left)
	})

	t.Run("should create spacing with zero horizontal value", func(t *testing.T) {
		t.Parallel()

		spacing := SpacingXY(0, 3)

		assert.Equal(t, 3, spacing.Top)
		assert.Equal(t, 0, spacing.Right)
		assert.Equal(t, 3, spacing.Bottom)
		assert.Equal(t, 0, spacing.Left)
	})

	t.Run("should create spacing with zero vertical value", func(t *testing.T) {
		t.Parallel()

		spacing := SpacingXY(5, 0)

		assert.Equal(t, 0, spacing.Top)
		assert.Equal(t, 5, spacing.Right)
		assert.Equal(t, 0, spacing.Bottom)
		assert.Equal(t, 5, spacing.Left)
	})

	t.Run("should create spacing with equal horizontal and vertical values", func(t *testing.T) {
		t.Parallel()

		spacing := SpacingXY(7, 7)

		assert.Equal(t, 7, spacing.Top)
		assert.Equal(t, 7, spacing.Right)
		assert.Equal(t, 7, spacing.Bottom)
		assert.Equal(t, 7, spacing.Left)
	})
}

func TestFontWeight_IsBold(t *testing.T) {
	t.Parallel()

	t.Run("should return true for weight 700 (bold)", func(t *testing.T) {
		t.Parallel()

		weight := WeightBold

		result := weight.IsBold()

		assert.True(t, result)
	})

	t.Run("should return true for weight exactly 600", func(t *testing.T) {
		t.Parallel()

		weight := FontWeight(600)

		result := weight.IsBold()

		assert.True(t, result)
	})

	t.Run("should return true for weight 900 (black)", func(t *testing.T) {
		t.Parallel()

		weight := FontWeight(900)

		result := weight.IsBold()

		assert.True(t, result)
	})

	t.Run("should return false for weight 400 (normal)", func(t *testing.T) {
		t.Parallel()

		weight := WeightNormal

		result := weight.IsBold()

		assert.False(t, result)
	})

	t.Run("should return false for weight 599", func(t *testing.T) {
		t.Parallel()

		weight := FontWeight(599)

		result := weight.IsBold()

		assert.False(t, result)
	})

	t.Run("should return false for weight 100 (thin)", func(t *testing.T) {
		t.Parallel()

		weight := FontWeight(100)

		result := weight.IsBold()

		assert.False(t, result)
	})
}

func TestFontStyle_IsItalic(t *testing.T) {
	t.Parallel()

	t.Run("should return true for StyleItalic", func(t *testing.T) {
		t.Parallel()

		style := StyleItalic

		result := style.IsItalic()

		assert.True(t, result)
	})

	t.Run("should return true for StyleOblique", func(t *testing.T) {
		t.Parallel()

		style := StyleOblique

		result := style.IsItalic()

		assert.True(t, result)
	})

	t.Run("should return false for StyleNormal", func(t *testing.T) {
		t.Parallel()

		style := StyleNormal

		result := style.IsItalic()

		assert.False(t, result)
	})
}

func TestTextDecoration_HasUnderline(t *testing.T) {
	t.Parallel()

	t.Run("should return true when underline decoration is set", func(t *testing.T) {
		t.Parallel()

		decoration := DecorationUnderline

		result := decoration.HasUnderline()

		assert.True(t, result)
	})

	t.Run("should return true when underline and line-through are both set", func(t *testing.T) {
		t.Parallel()

		decoration := DecorationUnderline | DecorationLineThrough

		result := decoration.HasUnderline()

		assert.True(t, result)
	})

	t.Run("should return false when only line-through is set", func(t *testing.T) {
		t.Parallel()

		decoration := DecorationLineThrough

		result := decoration.HasUnderline()

		assert.False(t, result)
	})

	t.Run("should return false when decoration is none", func(t *testing.T) {
		t.Parallel()

		decoration := DecorationNone

		result := decoration.HasUnderline()

		assert.False(t, result)
	})
}

func TestTextDecoration_HasLineThrough(t *testing.T) {
	t.Parallel()

	t.Run("should return true when line-through decoration is set", func(t *testing.T) {
		t.Parallel()

		decoration := DecorationLineThrough

		result := decoration.HasLineThrough()

		assert.True(t, result)
	})

	t.Run("should return true when underline and line-through are both set", func(t *testing.T) {
		t.Parallel()

		decoration := DecorationUnderline | DecorationLineThrough

		result := decoration.HasLineThrough()

		assert.True(t, result)
	})

	t.Run("should return false when only underline is set", func(t *testing.T) {
		t.Parallel()

		decoration := DecorationUnderline

		result := decoration.HasLineThrough()

		assert.False(t, result)
	})

	t.Run("should return false when decoration is none", func(t *testing.T) {
		t.Parallel()

		decoration := DecorationNone

		result := decoration.HasLineThrough()

		assert.False(t, result)
	})
}

func TestTextTransform_Apply(t *testing.T) {
	t.Parallel()

	t.Run("should transform text to uppercase", func(t *testing.T) {
		t.Parallel()

		transform := TransformUppercase
		text := "hello world"

		result := transform.Apply(text)

		assert.Equal(t, "HELLO WORLD", result)
	})

	t.Run("should transform text to lowercase", func(t *testing.T) {
		t.Parallel()

		transform := TransformLowercase
		text := "HELLO WORLD"

		result := transform.Apply(text)

		assert.Equal(t, "hello world", result)
	})

	t.Run("should capitalize first letter of each word", func(t *testing.T) {
		t.Parallel()

		transform := TransformCapitalize
		text := "hello world"

		result := transform.Apply(text)

		assert.Equal(t, "Hello World", result)
	})

	t.Run("should return unchanged text when transform is none", func(t *testing.T) {
		t.Parallel()

		transform := TransformNone
		text := "Hello World"

		result := transform.Apply(text)

		assert.Equal(t, "Hello World", result)
	})

	t.Run("should handle empty string", func(t *testing.T) {
		t.Parallel()

		transform := TransformUppercase
		text := ""

		result := transform.Apply(text)

		assert.Empty(t, result)
	})

	t.Run("should capitalize and lowercase remaining letters in words", func(t *testing.T) {
		t.Parallel()

		transform := TransformCapitalize
		text := "hELLO wORLD"

		result := transform.Apply(text)

		assert.Equal(t, "Hello World", result)
	})

	t.Run("should handle multiple spaces in capitalize", func(t *testing.T) {
		t.Parallel()

		transform := TransformCapitalize
		text := "hello   world"

		result := transform.Apply(text)

		assert.Equal(t, "Hello   World", result)
	})

	t.Run("should handle unicode characters in uppercase", func(t *testing.T) {
		t.Parallel()

		transform := TransformUppercase
		text := "cafe"

		result := transform.Apply(text)

		assert.Equal(t, "CAFE", result)
	})
}

func TestBorderEdge_IsSet(t *testing.T) {
	t.Parallel()

	t.Run("should return true when style is single", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderSingle}

		result := edge.IsSet()

		assert.True(t, result)
	})

	t.Run("should return true when style is double", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderDouble}

		result := edge.IsSet()

		assert.True(t, result)
	})

	t.Run("should return true when style is round", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderRound}

		result := edge.IsSet()

		assert.True(t, result)
	})

	t.Run("should return true when style is bold", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderBold}

		result := edge.IsSet()

		assert.True(t, result)
	})

	t.Run("should return false when style is none", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderNone}

		result := edge.IsSet()

		assert.False(t, result)
	})

	t.Run("should return false when style is empty string", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: ""}

		result := edge.IsSet()

		assert.False(t, result)
	})

	t.Run("should return false for zero value BorderEdge", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{}

		result := edge.IsSet()

		assert.False(t, result)
	})
}

func TestBorderEdge_AddText(t *testing.T) {
	t.Parallel()

	t.Run("should add text with left alignment", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderSingle}

		result := edge.AddText("Title", TextAlignLeft)

		assert.Len(t, result.Texts, 1)
		assert.Equal(t, "Title", result.Texts[0].Text)
		assert.Equal(t, TextAlignLeft, result.Texts[0].Align)
	})

	t.Run("should add text with center alignment", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderSingle}

		result := edge.AddText("Centered", TextAlignCenter)

		assert.Len(t, result.Texts, 1)
		assert.Equal(t, "Centered", result.Texts[0].Text)
		assert.Equal(t, TextAlignCenter, result.Texts[0].Align)
	})

	t.Run("should add text with right alignment", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderSingle}

		result := edge.AddText("Right", TextAlignRight)

		assert.Len(t, result.Texts, 1)
		assert.Equal(t, "Right", result.Texts[0].Text)
		assert.Equal(t, TextAlignRight, result.Texts[0].Align)
	})

	t.Run("should append multiple texts", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderSingle}

		result := edge.AddText("First", TextAlignLeft).AddText("Second", TextAlignRight)

		assert.Len(t, result.Texts, 2)
		assert.Equal(t, "First", result.Texts[0].Text)
		assert.Equal(t, "Second", result.Texts[1].Text)
	})

	t.Run("should preserve original edge style and color", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderDouble, Color: "#FF0000"}

		result := edge.AddText("Title", TextAlignLeft)

		assert.Equal(t, BorderDouble, result.Style)
		assert.Equal(t, Color("#FF0000"), result.Color)
	})
}

func TestBorderEdge_AddTextWithStyle(t *testing.T) {
	t.Parallel()

	t.Run("should add styled text with color", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderSingle}
		textStyle := Style{Foreground: "#00FF00"}

		result := edge.AddTextWithStyle("Styled", TextAlignLeft, textStyle)

		assert.Len(t, result.Texts, 1)
		assert.Equal(t, "Styled", result.Texts[0].Text)
		assert.Equal(t, TextAlignLeft, result.Texts[0].Align)
		assert.Equal(t, Color("#00FF00"), result.Texts[0].Style.Foreground)
	})

	t.Run("should add styled text with font weight", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderSingle}
		textStyle := Style{FontWeight: WeightBold}

		result := edge.AddTextWithStyle("Bold", TextAlignCenter, textStyle)

		assert.Equal(t, WeightBold, result.Texts[0].Style.FontWeight)
	})

	t.Run("should chain multiple styled texts", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderSingle}
		style1 := Style{Foreground: "#FF0000"}
		style2 := Style{Foreground: "#0000FF"}

		result := edge.
			AddTextWithStyle("Red", TextAlignLeft, style1).
			AddTextWithStyle("Blue", TextAlignRight, style2)

		assert.Len(t, result.Texts, 2)
		assert.Equal(t, Color("#FF0000"), result.Texts[0].Style.Foreground)
		assert.Equal(t, Color("#0000FF"), result.Texts[1].Style.Foreground)
	})

	t.Run("should preserve original edge properties", func(t *testing.T) {
		t.Parallel()

		edge := BorderEdge{Style: BorderBold, Color: "#AABBCC"}
		textStyle := Style{Foreground: "#112233"}

		result := edge.AddTextWithStyle("Text", TextAlignLeft, textStyle)

		assert.Equal(t, BorderBold, result.Style)
		assert.Equal(t, Color("#AABBCC"), result.Color)
	})
}

func TestBorder_HasAny(t *testing.T) {
	t.Parallel()

	t.Run("should return true when top border is set", func(t *testing.T) {
		t.Parallel()

		border := Border{Top: BorderEdge{Style: BorderSingle}}

		result := border.HasAny()

		assert.True(t, result)
	})

	t.Run("should return true when right border is set", func(t *testing.T) {
		t.Parallel()

		border := Border{Right: BorderEdge{Style: BorderSingle}}

		result := border.HasAny()

		assert.True(t, result)
	})

	t.Run("should return true when bottom border is set", func(t *testing.T) {
		t.Parallel()

		border := Border{Bottom: BorderEdge{Style: BorderSingle}}

		result := border.HasAny()

		assert.True(t, result)
	})

	t.Run("should return true when left border is set", func(t *testing.T) {
		t.Parallel()

		border := Border{Left: BorderEdge{Style: BorderSingle}}

		result := border.HasAny()

		assert.True(t, result)
	})

	t.Run("should return true when all borders are set", func(t *testing.T) {
		t.Parallel()

		border := BorderAll(BorderSingle)

		result := border.HasAny()

		assert.True(t, result)
	})

	t.Run("should return false when no borders are set", func(t *testing.T) {
		t.Parallel()

		border := Border{}

		result := border.HasAny()

		assert.False(t, result)
	})

	t.Run("should return false when all borders are none", func(t *testing.T) {
		t.Parallel()

		border := Border{
			Top:    BorderEdge{Style: BorderNone},
			Right:  BorderEdge{Style: BorderNone},
			Bottom: BorderEdge{Style: BorderNone},
			Left:   BorderEdge{Style: BorderNone},
		}

		result := border.HasAny()

		assert.False(t, result)
	})
}

func TestBorderAll(t *testing.T) {
	t.Parallel()

	t.Run("should create border with single style on all sides", func(t *testing.T) {
		t.Parallel()

		border := BorderAll(BorderSingle)

		assert.Equal(t, BorderSingle, border.Top.Style)
		assert.Equal(t, BorderSingle, border.Right.Style)
		assert.Equal(t, BorderSingle, border.Bottom.Style)
		assert.Equal(t, BorderSingle, border.Left.Style)
	})

	t.Run("should create border with double style on all sides", func(t *testing.T) {
		t.Parallel()

		border := BorderAll(BorderDouble)

		assert.Equal(t, BorderDouble, border.Top.Style)
		assert.Equal(t, BorderDouble, border.Right.Style)
		assert.Equal(t, BorderDouble, border.Bottom.Style)
		assert.Equal(t, BorderDouble, border.Left.Style)
	})

	t.Run("should create border with round style on all sides", func(t *testing.T) {
		t.Parallel()

		border := BorderAll(BorderRound)

		assert.Equal(t, BorderRound, border.Top.Style)
		assert.Equal(t, BorderRound, border.Right.Style)
		assert.Equal(t, BorderRound, border.Bottom.Style)
		assert.Equal(t, BorderRound, border.Left.Style)
	})

	t.Run("should create border with bold style on all sides", func(t *testing.T) {
		t.Parallel()

		border := BorderAll(BorderBold)

		assert.Equal(t, BorderBold, border.Top.Style)
		assert.Equal(t, BorderBold, border.Right.Style)
		assert.Equal(t, BorderBold, border.Bottom.Style)
		assert.Equal(t, BorderBold, border.Left.Style)
	})

	t.Run("should not set color when using BorderAll", func(t *testing.T) {
		t.Parallel()

		border := BorderAll(BorderSingle)

		assert.False(t, border.Top.Color.IsSet())
		assert.False(t, border.Right.Color.IsSet())
		assert.False(t, border.Bottom.Color.IsSet())
		assert.False(t, border.Left.Color.IsSet())
	})
}

func TestBorderAllWithColor(t *testing.T) {
	t.Parallel()

	t.Run("should create border with style and color on all sides", func(t *testing.T) {
		t.Parallel()

		border := BorderAllWithColor(BorderSingle, "#FF0000")

		assert.Equal(t, BorderSingle, border.Top.Style)
		assert.Equal(t, Color("#FF0000"), border.Top.Color)
		assert.Equal(t, BorderSingle, border.Right.Style)
		assert.Equal(t, Color("#FF0000"), border.Right.Color)
		assert.Equal(t, BorderSingle, border.Bottom.Style)
		assert.Equal(t, Color("#FF0000"), border.Bottom.Color)
		assert.Equal(t, BorderSingle, border.Left.Style)
		assert.Equal(t, Color("#FF0000"), border.Left.Color)
	})

	t.Run("should work with different border styles", func(t *testing.T) {
		t.Parallel()

		border := BorderAllWithColor(BorderDouble, "#00FF00")

		assert.Equal(t, BorderDouble, border.Top.Style)
		assert.Equal(t, Color("#00FF00"), border.Top.Color)
	})

	t.Run("should work with color without hash prefix", func(t *testing.T) {
		t.Parallel()

		border := BorderAllWithColor(BorderRound, "0000FF")

		assert.Equal(t, Color("0000FF"), border.Top.Color)
		assert.True(t, border.Top.Color.IsSet())
	})
}

func TestBorderXY(t *testing.T) {
	t.Parallel()

	t.Run("should create border with different horizontal and vertical edges", func(t *testing.T) {
		t.Parallel()

		horizontal := BorderEdge{Style: BorderSingle, Color: "#FF0000"}
		vertical := BorderEdge{Style: BorderDouble, Color: "#00FF00"}

		border := BorderXY(horizontal, vertical)

		assert.Equal(t, BorderDouble, border.Top.Style)
		assert.Equal(t, Color("#00FF00"), border.Top.Color)
		assert.Equal(t, BorderSingle, border.Right.Style)
		assert.Equal(t, Color("#FF0000"), border.Right.Color)
		assert.Equal(t, BorderDouble, border.Bottom.Style)
		assert.Equal(t, Color("#00FF00"), border.Bottom.Color)
		assert.Equal(t, BorderSingle, border.Left.Style)
		assert.Equal(t, Color("#FF0000"), border.Left.Color)
	})

	t.Run("should apply horizontal edge to left and right", func(t *testing.T) {
		t.Parallel()

		horizontal := BorderEdge{Style: BorderBold}
		vertical := BorderEdge{Style: BorderNone}

		border := BorderXY(horizontal, vertical)

		assert.True(t, border.Left.IsSet())
		assert.True(t, border.Right.IsSet())
		assert.False(t, border.Top.IsSet())
		assert.False(t, border.Bottom.IsSet())
	})

	t.Run("should apply vertical edge to top and bottom", func(t *testing.T) {
		t.Parallel()

		horizontal := BorderEdge{Style: BorderNone}
		vertical := BorderEdge{Style: BorderRound}

		border := BorderXY(horizontal, vertical)

		assert.False(t, border.Left.IsSet())
		assert.False(t, border.Right.IsSet())
		assert.True(t, border.Top.IsSet())
		assert.True(t, border.Bottom.IsSet())
	})
}

func TestBorderTRBL(t *testing.T) {
	t.Parallel()

	t.Run("should create border with explicit edges for each side", func(t *testing.T) {
		t.Parallel()

		top := BorderEdge{Style: BorderSingle, Color: "#FF0000"}
		right := BorderEdge{Style: BorderDouble, Color: "#00FF00"}
		bottom := BorderEdge{Style: BorderRound, Color: "#0000FF"}
		left := BorderEdge{Style: BorderBold, Color: "#FFFF00"}

		border := BorderTRBL(top, right, bottom, left)

		assert.Equal(t, BorderSingle, border.Top.Style)
		assert.Equal(t, Color("#FF0000"), border.Top.Color)
		assert.Equal(t, BorderDouble, border.Right.Style)
		assert.Equal(t, Color("#00FF00"), border.Right.Color)
		assert.Equal(t, BorderRound, border.Bottom.Style)
		assert.Equal(t, Color("#0000FF"), border.Bottom.Color)
		assert.Equal(t, BorderBold, border.Left.Style)
		assert.Equal(t, Color("#FFFF00"), border.Left.Color)
	})

	t.Run("should allow some edges to be unset", func(t *testing.T) {
		t.Parallel()

		top := BorderEdge{Style: BorderSingle}
		right := BorderEdge{}
		bottom := BorderEdge{Style: BorderSingle}
		left := BorderEdge{}

		border := BorderTRBL(top, right, bottom, left)

		assert.True(t, border.Top.IsSet())
		assert.False(t, border.Right.IsSet())
		assert.True(t, border.Bottom.IsSet())
		assert.False(t, border.Left.IsSet())
	})

	t.Run("should preserve text on edges", func(t *testing.T) {
		t.Parallel()

		top := BorderEdge{Style: BorderSingle}.AddText("Title", TextAlignCenter)
		right := BorderEdge{Style: BorderSingle}
		bottom := BorderEdge{Style: BorderSingle}.AddText("Footer", TextAlignLeft)
		left := BorderEdge{Style: BorderSingle}

		border := BorderTRBL(top, right, bottom, left)

		assert.Len(t, border.Top.Texts, 1)
		assert.Equal(t, "Title", border.Top.Texts[0].Text)
		assert.Len(t, border.Bottom.Texts, 1)
		assert.Equal(t, "Footer", border.Bottom.Texts[0].Text)
	})
}

func TestBorderAllWithTitle(t *testing.T) {
	t.Parallel()

	t.Run("should create border with title on top edge", func(t *testing.T) {
		t.Parallel()

		border := BorderAllWithTitle(BorderSingle, "#FFFFFF", "My Title")

		assert.Len(t, border.Top.Texts, 1)
		assert.Equal(t, "My Title", border.Top.Texts[0].Text)
		assert.Equal(t, TextAlignLeft, border.Top.Texts[0].Align)
	})

	t.Run("should set style on all edges", func(t *testing.T) {
		t.Parallel()

		border := BorderAllWithTitle(BorderDouble, "#FF0000", "Title")

		assert.Equal(t, BorderDouble, border.Top.Style)
		assert.Equal(t, BorderDouble, border.Right.Style)
		assert.Equal(t, BorderDouble, border.Bottom.Style)
		assert.Equal(t, BorderDouble, border.Left.Style)
	})

	t.Run("should set color on all edges", func(t *testing.T) {
		t.Parallel()

		border := BorderAllWithTitle(BorderRound, "#00FF00", "Title")

		assert.Equal(t, Color("#00FF00"), border.Top.Color)
		assert.Equal(t, Color("#00FF00"), border.Right.Color)
		assert.Equal(t, Color("#00FF00"), border.Bottom.Color)
		assert.Equal(t, Color("#00FF00"), border.Left.Color)
	})

	t.Run("should not add text to other edges", func(t *testing.T) {
		t.Parallel()

		border := BorderAllWithTitle(BorderBold, "#0000FF", "Title")

		assert.Empty(t, border.Right.Texts)
		assert.Empty(t, border.Bottom.Texts)
		assert.Empty(t, border.Left.Texts)
	})

	t.Run("should handle empty title", func(t *testing.T) {
		t.Parallel()

		border := BorderAllWithTitle(BorderSingle, "#FFFFFF", "")

		assert.Len(t, border.Top.Texts, 1)
		assert.Empty(t, border.Top.Texts[0].Text)
	})
}
