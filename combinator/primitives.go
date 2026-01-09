package combinator

import (
	"fmt"
	"slices"
)

// Char matches a single specific character and returns it as a rune.
// Fails with an error message showing the expected and actual characters.
//
// Example:
//
//	result := Parse(Char('a'), "abc")
//	// result.Value == 'a'
func Char(r rune) Parser[rune] {
	return func(state State) Result[rune] {
		if state.IsEOF() {
			return Failure[rune](fmt.Errorf("unexpected EOF, expected '%c' at line %d, col %d", r, state.Line, state.Col), state)
		}

		if state.Current() != r {
			return Failure[rune](fmt.Errorf("expected '%c', got '%c' at line %d, col %d", r, state.Current(), state.Line, state.Col), state)
		}

		return Success(r, state.Advance())
	}
}

// String matches an exact string and returns it.
// Fails if the input does not start with the expected string.
//
// Example:
//
//	result := Parse(String("hello"), "hello world")
//	// result.Value == "hello"
func String(s string) Parser[string] {
	return func(state State) Result[string] {
		current := state

		for _, r := range s {
			if current.IsEOF() {
				return Failure[string](fmt.Errorf("unexpected EOF, expected '%s' at line %d, col %d", s, state.Line, state.Col), state)
			}

			if current.Current() != r {
				return Failure[string](fmt.Errorf("expected '%s' at line %d, col %d", s, state.Line, state.Col), state)
			}

			current = current.Advance()
		}

		return Success(s, current)
	}
}

// Satisfy matches a single character that satisfies the predicate function.
// Returns the matched rune on success.
//
// This is the fundamental building block for character-based parsers.
// Most character class parsers like [Digit] and [Letter] are built on Satisfy.
//
// Example:
//
//	vowel := Satisfy(func(r rune) bool {
//		return r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u'
//	})
func Satisfy(pred func(rune) bool) Parser[rune] {
	return func(state State) Result[rune] {
		if state.IsEOF() {
			return Failure[rune](fmt.Errorf("unexpected EOF at line %d, col %d", state.Line, state.Col), state)
		}

		r := state.Current()
		if !pred(r) {
			return Failure[rune](fmt.Errorf("unexpected '%c' at line %d, col %d", r, state.Line, state.Col), state)
		}

		return Success(r, state.Advance())
	}
}

// Any matches any single character and returns it as a rune.
// Fails only at end of input.
func Any() Parser[rune] {
	return func(state State) Result[rune] {
		if state.IsEOF() {
			return Failure[rune](fmt.Errorf("unexpected EOF at line %d, col %d", state.Line, state.Col), state)
		}

		return Success(state.Current(), state.Advance())
	}
}

// EOF matches the end of input and returns nil.
// Use to ensure the parser consumed all input.
//
// Example:
//
//	complete := Seq2(Integer(), EOF())
//	result := Parse(complete, "42")    // succeeds
//	result = Parse(complete, "42abc")  // fails
func EOF() Parser[struct{}] {
	return func(state State) Result[struct{}] {
		if !state.IsEOF() {
			return Failure[struct{}](fmt.Errorf("expected EOF, got '%c' at line %d, col %d", state.Current(), state.Line, state.Col), state)
		}

		return Success(struct{}{}, state)
	}
}

// OneOf matches any single character that appears in the provided string.
// Returns the matched rune.
//
// Example:
//
//	op := OneOf("+-*/")
//	result := Parse(op, "+") // result.Value == '+'
func OneOf(chars string) Parser[rune] {
	runes := []rune(chars)
	return Satisfy(func(r rune) bool {
		return slices.Contains(runes, r)
	})
}

// NoneOf matches any single character that does not appear in the provided string.
// Returns the matched rune.
//
// Example:
//
//	notQuote := NoneOf("\"")
//	content := Many(notQuote) // matches until a quote
func NoneOf(chars string) Parser[rune] {
	runes := []rune(chars)
	return Satisfy(func(r rune) bool {
		return !slices.Contains(runes, r)
	})
}

// Range matches any character within the inclusive rune range [from, to].
// Returns the matched rune.
//
// Parameters:
//   - from: lower bound of the range (inclusive)
//   - to: upper bound of the range (inclusive)
//
// Example:
//
//	lowercase := Range('a', 'z')
//	result := Parse(lowercase, "m") // result.Value == 'm'
func Range(from, to rune) Parser[rune] {
	return Satisfy(func(r rune) bool {
		return r >= from && r <= to
	})
}
