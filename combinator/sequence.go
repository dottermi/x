package combinator

import "fmt"

// Pair holds two values of potentially different types.
type Pair[A, B any] struct {
	First  A
	Second B
}

// Triple holds three values of potentially different types.
type Triple[A, B, C any] struct {
	First  A
	Second B
	Third  C
}

// Seq2 runs two parsers in sequence and returns a Pair of results.
// Fails immediately if either parser fails.
//
// Example:
//
//	// Match "a" followed by "b"
//	ab := Seq2(Char('a'), Char('b'))
//	result := Parse(ab, "ab")
//	// result.Value == Pair[rune, rune]{First: 'a', Second: 'b'}
func Seq2[A, B any](p1 Parser[A], p2 Parser[B]) Parser[Pair[A, B]] {
	return func(state State) Result[Pair[A, B]] {
		r1 := p1(state)
		if !r1.OK {
			return Failure[Pair[A, B]](r1.Err, r1.State)
		}
		r2 := p2(r1.State)
		if !r2.OK {
			return Failure[Pair[A, B]](r2.Err, r2.State)
		}
		return Success(Pair[A, B]{First: r1.Value, Second: r2.Value}, r2.State)
	}
}

// Seq3 runs three parsers in sequence and returns a Triple of results.
// Fails immediately if any parser fails.
func Seq3[A, B, C any](p1 Parser[A], p2 Parser[B], p3 Parser[C]) Parser[Triple[A, B, C]] {
	return func(state State) Result[Triple[A, B, C]] {
		r1 := p1(state)
		if !r1.OK {
			return Failure[Triple[A, B, C]](r1.Err, r1.State)
		}
		r2 := p2(r1.State)
		if !r2.OK {
			return Failure[Triple[A, B, C]](r2.Err, r2.State)
		}
		r3 := p3(r2.State)
		if !r3.OK {
			return Failure[Triple[A, B, C]](r3.Err, r3.State)
		}
		return Success(Triple[A, B, C]{First: r1.Value, Second: r2.Value, Third: r3.Value}, r3.State)
	}
}

// Choice tries parsers in order and returns the first successful result.
// All parsers must return the same type.
// Fails only if all alternatives fail, returning the last error.
//
// Example:
//
//	// Match "true" or "false"
//	boolean := Choice(String("true"), String("false"))
//	result := Parse(boolean, "true")
//	// result.Value == "true"
func Choice[T any](parsers ...Parser[T]) Parser[T] {
	return func(state State) Result[T] {
		var lastErr error

		for _, p := range parsers {
			r := p(state)
			if r.OK {
				return r
			}
			lastErr = r.Err
		}

		if lastErr != nil {
			return Failure[T](lastErr, state)
		}
		return Failure[T](fmt.Errorf("no alternatives matched at line %d, col %d", state.Line, state.Col), state)
	}
}

// Many matches zero or more occurrences of a parser.
// Returns a slice of all matched values.
// Always succeeds (returns empty slice if no matches).
//
// Example:
//
//	digits := Many(Digit())
//	result := Parse(digits, "123abc")
//	// result.Value == []rune{'1', '2', '3'}
func Many[T any](p Parser[T]) Parser[[]T] {
	return func(state State) Result[[]T] {
		var values []T
		current := state

		for {
			r := p(current)
			if !r.OK {
				break
			}
			values = append(values, r.Value)

			// Avoid infinite loop if parser doesn't consume input
			if r.State.Pos == current.Pos {
				break
			}
			current = r.State
		}

		return Success(values, current)
	}
}

// Many1 matches one or more occurrences of a parser.
// Returns a slice of all matched values.
// Fails if no matches are found.
//
// Example:
//
//	digits := Many1(Digit())
//	result := Parse(digits, "123")
//	// result.Value == []rune{'1', '2', '3'}
func Many1[T any](p Parser[T]) Parser[[]T] {
	return func(state State) Result[[]T] {
		first := p(state)
		if !first.OK {
			return Failure[[]T](first.Err, first.State)
		}

		rest := Many(p)(first.State)
		values := append([]T{first.Value}, rest.Value...)

		return Success(values, rest.State)
	}
}

// Opt makes a parser optional, returning a pointer (nil on failure).
// Always succeeds.
//
// Example:
//
//	sign := Opt(Char('-'))
//	result := Parse(sign, "42")  // result.Value == nil
//	result = Parse(sign, "-42") // result.Value == ptr to '-'
func Opt[T any](p Parser[T]) Parser[*T] {
	return func(state State) Result[*T] {
		r := p(state)
		if r.OK {
			return Success(&r.Value, r.State)
		}
		return Success[*T](nil, state)
	}
}

// And sequences two parsers, returning only the second parser's result.
// Both parsers must succeed.
func And[A, B any](p1 Parser[A], p2 Parser[B]) Parser[B] {
	return func(state State) Result[B] {
		r1 := p1(state)
		if !r1.OK {
			return Failure[B](r1.Err, r1.State)
		}
		return p2(r1.State)
	}
}

// Left sequences two parsers, returning only the first parser's result.
// Both parsers must succeed; the second parser's value is discarded.
//
// Example:
//
//	// Match a number followed by a semicolon, keep only the number
//	num := Left(Integer(), Char(';'))
func Left[A, B any](p1 Parser[A], p2 Parser[B]) Parser[A] {
	return func(state State) Result[A] {
		r1 := p1(state)
		if !r1.OK {
			return r1
		}
		r2 := p2(r1.State)
		if !r2.OK {
			return Failure[A](r2.Err, r2.State)
		}
		return Success(r1.Value, r2.State)
	}
}

// Right sequences two parsers, returning only the second parser's result.
// Both parsers must succeed; the first parser's value is discarded.
//
// Example:
//
//	// Skip whitespace before a number
//	num := Right(Spaces(), Integer())
func Right[A, B any](p1 Parser[A], p2 Parser[B]) Parser[B] {
	return func(state State) Result[B] {
		r1 := p1(state)
		if !r1.OK {
			return Failure[B](r1.Err, r1.State)
		}
		return p2(r1.State)
	}
}

// Count matches exactly n occurrences of a parser.
// Returns a slice of all matched values.
// Fails if fewer than n matches are found.
//
// Example:
//
//	hex := Count(2, HexDigit()) // match exactly 2 hex digits
func Count[T any](n int, p Parser[T]) Parser[[]T] {
	return func(state State) Result[[]T] {
		var values []T
		current := state

		for range n {
			r := p(current)
			if !r.OK {
				return Failure[[]T](r.Err, r.State)
			}
			values = append(values, r.Value)
			current = r.State
		}

		return Success(values, current)
	}
}
