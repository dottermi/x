package style

import (
	"strings"
	"unicode"
)

// FontWeight represents the weight (boldness) of text.
// Values follow CSS font-weight: 100 (thin) to 900 (black).
// Normal is 400, Bold is 700. Values >= 600 render as bold in terminals.
type FontWeight int

// FontWeight values.
const (
	WeightNormal FontWeight = 400
	WeightBold   FontWeight = 700
)

// TextTransform represents text transformation (case changes).
type TextTransform int

// TextTransform values.
const (
	TransformNone TextTransform = iota
	TransformUppercase
	TransformLowercase
	TransformCapitalize
)

// FontStyle represents the style of text (normal, italic, oblique).
type FontStyle int

// FontStyle values.
const (
	StyleNormal FontStyle = iota
	StyleItalic
	StyleOblique // Rendered as italic in terminals
)

// TextDecoration represents text decorations (underline, line-through).
// Values can be combined using bitwise OR.
type TextDecoration uint8

// TextDecoration values.
const (
	DecorationNone        TextDecoration = 0
	DecorationUnderline   TextDecoration = 1 << 0
	DecorationLineThrough TextDecoration = 1 << 1
)

// TextAlign represents horizontal text alignment within a container.
type TextAlign int

// TextAlign values.
const (
	// TextAlignLeft aligns text to the left (default).
	TextAlignLeft TextAlign = iota
	// TextAlignCenter centers text horizontally.
	TextAlignCenter
	// TextAlignRight aligns text to the right.
	TextAlignRight
)

// TextWrap controls how text wraps within a container.
type TextWrap int

// TextWrap values.
const (
	// WrapNone disables wrapping; text may overflow.
	WrapNone TextWrap = iota
	// WrapWord wraps at word boundaries.
	WrapWord
	// WrapChar wraps at character boundaries.
	WrapChar
)

// TextOverflow controls how overflowed text is handled.
type TextOverflow int

// TextOverflow values.
const (
	// TextOverflowClip clips the text at the boundary (default).
	TextOverflowClip TextOverflow = iota
	// TextOverflowEllipsis truncates with "..." when text overflows.
	TextOverflowEllipsis
)

// IsBold returns true if the weight should render as bold (>= 600).
func (w FontWeight) IsBold() bool {
	return w >= 600
}

// IsItalic returns true if the style should render as italic.
func (s FontStyle) IsItalic() bool {
	return s == StyleItalic || s == StyleOblique
}

// HasUnderline returns true if underline decoration is set.
func (d TextDecoration) HasUnderline() bool {
	return d&DecorationUnderline != 0
}

// HasLineThrough returns true if line-through decoration is set.
func (d TextDecoration) HasLineThrough() bool {
	return d&DecorationLineThrough != 0
}

// Apply transforms the given text according to the transform type.
func (t TextTransform) Apply(text string) string {
	if t == TransformUppercase {
		return strings.ToUpper(text)
	}

	if t == TransformLowercase {
		return strings.ToLower(text)
	}

	if t == TransformCapitalize {
		return capitalize(text)
	}

	return text
}

// capitalize capitalizes the first letter of each word.
func capitalize(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	capitalizeNext := true

	for _, r := range s {
		switch {
		case unicode.IsSpace(r):
			capitalizeNext = true
			result.WriteRune(r)
		case capitalizeNext:
			result.WriteRune(unicode.ToUpper(r))
			capitalizeNext = false
		default:
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}
