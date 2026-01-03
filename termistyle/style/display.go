package style

// Display controls how an element lays out its children.
type Display int

const (
	// Block stacks children vertically.
	Block Display = iota
	// Flex enables flexbox layout with configurable direction.
	Flex
	// None hides the element from rendering.
	None
)

// PositionType determines how an element is positioned.
type PositionType int

const (
	// Relative positions element in normal flow.
	Relative PositionType = iota
	// Absolute positions element relative to its container.
	Absolute
)

// FlexDirection sets the main axis for flex layout.
type FlexDirection int

const (
	// Row arranges children horizontally.
	Row FlexDirection = iota
	// Column arranges children vertically.
	Column
)

// Justify controls distribution of children along the main axis (justify-content).
type Justify int

const (
	// JustifyStart aligns children to the start (flex-start).
	JustifyStart Justify = iota
	// JustifyCenter centers children.
	JustifyCenter
	// JustifyEnd aligns children to the end (flex-end).
	JustifyEnd
	// JustifyBetween distributes with equal space between (space-between).
	JustifyBetween
	// JustifyAround distributes with equal space around (space-around).
	JustifyAround
)

// CSS-like aliases for Justify constants.
const (
	SpaceBetween = JustifyBetween
	SpaceAround  = JustifyAround
)

// Align controls positioning of children along the cross axis (align-items).
type Align int

const (
	// AlignStart positions children at cross-axis start (flex-start).
	AlignStart Align = iota
	// AlignCenter centers children on cross axis.
	AlignCenter
	// AlignEnd positions children at cross-axis end (flex-end).
	AlignEnd
	// AlignStretch expands children to fill cross axis.
	AlignStretch
)

// Stretch is a CSS-like alias for AlignStretch.
const Stretch = AlignStretch

// FlexWrap controls whether flex items wrap onto multiple lines.
type FlexWrap int

const (
	// NoWrap keeps all items on a single line (default).
	NoWrap FlexWrap = iota
	// Wrap allows items to wrap onto multiple lines.
	Wrap
)

// Overflow controls how content that exceeds the container is handled.
type Overflow int

const (
	// OverflowVisible allows content to render outside container bounds (default).
	OverflowVisible Overflow = iota
	// OverflowHidden clips content that exceeds container bounds.
	OverflowHidden
)
