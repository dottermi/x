// Package combinator provides parser combinators for building recursive descent parsers.
//
// Parser combinators are small, composable functions that can be combined to create
// complex parsers from simple building blocks. This approach offers several advantages
// over traditional parsing methods:
//
//   - Declarative: parsers read like grammar specifications
//   - Composable: small parsers combine into larger ones
//   - Type-safe: compile-time checking of parser composition
//   - No code generation: parsers are regular Go functions
//
// The library is organized into layers of increasing abstraction:
//
//   - Primitives: [Char], [String], [Satisfy], [Any], [EOF]
//   - Character classes: [Digit], [Letter], [Space], [AlphaNum]
//   - Combinators: [Seq], [Choice], [Many], [Many1], [Opt], [Map]
//   - Token parsers: [Ident], [Integer], [StringLit], [Between]
//
// # Basic Usage
//
// Parse a single character:
//
//	result := combinator.Parse(combinator.Char('a'), "abc")
//	// result.Value == 'a'
//
// Combine parsers to match patterns:
//
//	digit := combinator.Digit()
//	number := combinator.Many1(digit)
//	result := combinator.Parse(number, "123")
//
// # Building Grammars
//
// Use [Seq] for sequences and [Choice] for alternatives:
//
//	// Match "true" or "false"
//	boolean := combinator.Choice(
//		combinator.String("true"),
//		combinator.String("false"),
//	)
//
// Transform results with [Map]:
//
//	integer := combinator.Map(combinator.Integer(), func(v any) any {
//		return v.(int64) * 2
//	})
//
// # Recursive Grammars
//
// Use [Rule] and [Ref] for recursive definitions:
//
//	var expr combinator.Rule
//	expr = func() combinator.Parser {
//		return combinator.Choice(
//			combinator.Integer(),
//			combinator.Parens(combinator.Ref(&expr)),
//		)
//	}
package combinator

// State represents the current position and context within the input being parsed.
// Tracks line and column numbers for meaningful error messages.
//
// Create a new State with [NewState] rather than constructing directly.
type State struct {
	Input []rune // Input stores the complete input as runes for Unicode support.
	Pos   int    // Pos is the current byte position in Input.
	Line  int    // Line is the current line number (1-indexed).
	Col   int    // Col is the current column number (1-indexed).
}

// NewState creates a parser state initialized at the beginning of the input string.
// Converts the input to runes for proper Unicode handling.
//
// Example:
//
//	state := NewState("hello")
//	fmt.Println(state.Current()) // 'h'
func NewState(input string) State {
	return State{
		Input: []rune(input),
		Pos:   0,
		Line:  1,
		Col:   1,
	}
}

// Current returns the rune at the current position, or 0 if at end of input.
func (s State) Current() rune {
	if s.Pos >= len(s.Input) {
		return 0
	}
	return s.Input[s.Pos]
}

// IsEOF reports whether the parser has reached the end of input.
func (s State) IsEOF() bool {
	return s.Pos >= len(s.Input)
}

// Advance returns a new State moved forward by one rune.
// Updates line and column tracking when encountering newlines.
// Returns the same state unchanged if already at EOF.
func (s State) Advance() State {
	if s.IsEOF() {
		return s
	}

	next := State{
		Input: s.Input,
		Pos:   s.Pos + 1,
		Line:  s.Line,
		Col:   s.Col + 1,
	}

	if s.Current() == '\n' {
		next.Line++
		next.Col = 1
	}

	return next
}

// AdvanceN returns a new State moved forward by n runes.
// Equivalent to calling Advance n times.
func (s State) AdvanceN(n int) State {
	current := s
	for range n {
		current = current.Advance()
	}
	return current
}

// Result holds the outcome of applying a parser to input.
// Check OK to determine success before accessing Value or Err.
type Result struct {
	OK    bool  // OK is true when parsing succeeded.
	Value any   // Value contains the parsed result on success.
	State State // State is the parser position after this result.
	Err   error // Err contains the error message on failure.
}

// Parser is a function that consumes input and produces a result.
// All parsing primitives and combinators in this package are Parsers.
//
// A Parser receives the current State, attempts to match, and returns a Result.
// On success, the Result contains the matched value and advanced State.
// On failure, the Result contains an error and the original State.
type Parser func(State) Result

// Rule enables recursive grammar definitions by deferring parser construction.
// Use with [Ref] to create self-referential parsers.
//
// Example:
//
//	var expr Rule
//	expr = func() Parser {
//		return Choice(Integer(), Parens(Ref(&expr)))
//	}
type Rule func() Parser

// Success constructs a successful parse result with the given value and updated state.
// Used internally by parsers; most users should use the higher-level combinators.
func Success(value any, state State) Result {
	return Result{
		OK:    true,
		Value: value,
		State: state,
	}
}

// Failure constructs a failed parse result with the given error and state.
// Used internally by parsers; most users should use the higher-level combinators.
func Failure(err error, state State) Result {
	return Result{
		OK:    false,
		Err:   err,
		State: state,
	}
}

// Parse runs a parser on the input string and returns the result.
// This is the main entry point for using parsers.
//
// Example:
//
//	result := Parse(Integer(), "42")
//	if result.OK {
//		fmt.Println(result.Value) // 42
//	}
func Parse(p Parser, input string) Result {
	return p(NewState(input))
}

// Run executes a parser, updates the state on success, and returns the typed value.
// This helper reduces boilerplate when writing sequential parsers manually.
//
// Example:
//
//	state := NewState("hello world")
//	name, err := Run[string](Ident(), &state)
//	if err != nil {
//		return Failure(err, state)
//	}
//	// state is now advanced past the identifier
func Run[T any](p Parser, s *State) (T, error) {
	var zero T
	res := p(*s)
	if !res.OK {
		return zero, res.Err
	}
	*s = res.State
	if v, ok := res.Value.(T); ok {
		return v, nil
	}
	return zero, nil
}
