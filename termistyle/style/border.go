package style

// BorderStyle identifies a border character set.
// Use with the Borders map to get actual characters.
type BorderStyle string

const (
	// BorderNone disables border rendering (no space occupied).
	BorderNone BorderStyle = "none"
	// BorderHidden renders invisible border (occupies space but uses spaces).
	BorderHidden BorderStyle = "hidden"
	// BorderSingle uses thin box-drawing characters.
	BorderSingle BorderStyle = "single"
	// BorderDouble uses double-line box-drawing characters.
	BorderDouble BorderStyle = "double"
	// BorderRound uses single lines with rounded corners.
	BorderRound BorderStyle = "round"
	// BorderBold uses thick box-drawing characters.
	BorderBold BorderStyle = "bold"
)

// BorderChars holds the Unicode characters for drawing a border.
// Each field represents a specific position or edge.
type BorderChars struct {
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
	Horizontal  rune
	Vertical    rune
}

// Borders maps border styles to their character sets.
// Use this to look up characters for a given BorderStyle.
//
// Example:
//
//	chars := style.Borders[style.BorderRound]
//	fmt.Printf("%c", chars.TopLeft)  // prints:
var Borders = map[BorderStyle]BorderChars{
	BorderHidden: {
		TopLeft:     ' ',
		TopRight:    ' ',
		BottomLeft:  ' ',
		BottomRight: ' ',
		Horizontal:  ' ',
		Vertical:    ' ',
	},
	BorderSingle: {
		TopLeft:     '┌',
		TopRight:    '┐',
		BottomLeft:  '└',
		BottomRight: '┘',
		Horizontal:  '─',
		Vertical:    '│',
	},
	BorderDouble: {
		TopLeft:     '╔',
		TopRight:    '╗',
		BottomLeft:  '╚',
		BottomRight: '╝',
		Horizontal:  '═',
		Vertical:    '║',
	},
	BorderRound: {
		TopLeft:     '╭',
		TopRight:    '╮',
		BottomLeft:  '╰',
		BottomRight: '╯',
		Horizontal:  '─',
		Vertical:    '│',
	},
	BorderBold: {
		TopLeft:     '┏',
		TopRight:    '┓',
		BottomLeft:  '┗',
		BottomRight: '┛',
		Horizontal:  '━',
		Vertical:    '┃',
	},
}

// BorderText represents a single text element on a horizontal border edge.
// Multiple BorderTexts can be placed on the same border (top or bottom).
type BorderText struct {
	Text  string    // The text content to display
	Align TextAlign // Position: left, center, or right
	Style Style     // Optional: text styling (color, bold, italic, etc.)
}

// BorderEdge defines a border edge with style and color.
// Use Texts to display text on horizontal edges (top/bottom).
type BorderEdge struct {
	Style BorderStyle
	Color Color
	Texts []BorderText // Multiple texts on this edge (only for horizontal edges)
}

// IsSet returns true if this border edge has a style.
func (e BorderEdge) IsSet() bool {
	return e.Style != "" && e.Style != BorderNone
}

// GetTexts returns all texts for this edge.
func (e BorderEdge) GetTexts() []BorderText {
	return e.Texts
}

// AddText appends a text element to the border edge.
// Returns a new BorderEdge for chaining.
func (e BorderEdge) AddText(text string, align TextAlign) BorderEdge {
	e.Texts = append(e.Texts, BorderText{Text: text, Align: align})
	return e
}

// AddTextWithStyle appends a styled text element to the border edge.
// Supports color, bold, italic, underline, etc.
// Returns a new BorderEdge for chaining.
func (e BorderEdge) AddTextWithStyle(text string, align TextAlign, style Style) BorderEdge {
	e.Texts = append(e.Texts, BorderText{Text: text, Align: align, Style: style})
	return e
}

// Border specifies border style and color for each side.
type Border struct {
	Top    BorderEdge
	Right  BorderEdge
	Bottom BorderEdge
	Left   BorderEdge
}

// HasAny returns true if any side has a border.
func (b Border) HasAny() bool {
	return b.Top.IsSet() || b.Right.IsSet() ||
		b.Bottom.IsSet() || b.Left.IsSet()
}

// BorderAll creates a Border with the same style on all sides.
func BorderAll(s BorderStyle) Border {
	edge := BorderEdge{Style: s}
	return Border{Top: edge, Right: edge, Bottom: edge, Left: edge}
}

// BorderAllWithColor creates a Border with the same style and color on all sides.
func BorderAllWithColor(s BorderStyle, c Color) Border {
	edge := BorderEdge{Style: s, Color: c}
	return Border{Top: edge, Right: edge, Bottom: edge, Left: edge}
}

// BorderXY creates a Border with horizontal (left/right) and vertical (top/bottom) styles.
func BorderXY(x, y BorderEdge) Border {
	return Border{Top: y, Right: x, Bottom: y, Left: x}
}

// BorderTRBL creates a Border with explicit top, right, bottom, left edges.
func BorderTRBL(top, right, bottom, left BorderEdge) Border {
	return Border{Top: top, Right: right, Bottom: bottom, Left: left}
}

// BorderAllWithTitle creates a Border with the same style on all sides
// and a title on the top edge (left-aligned).
func BorderAllWithTitle(s BorderStyle, c Color, title string) Border {
	edge := BorderEdge{Style: s, Color: c}
	return Border{
		Top:    edge.AddText(title, TextAlignLeft),
		Right:  edge,
		Bottom: edge,
		Left:   edge,
	}
}
