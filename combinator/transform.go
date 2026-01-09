package combinator

import "fmt"

// Map transforms the result of a parser using the provided function.
// The function receives the parsed value and returns a new value.
//
// Example:
//
//	upper := Map(Letter(), func(r rune) rune {
//		return unicode.ToUpper(r)
//	})
func Map[T, U any](p Parser[T], fn func(T) U) Parser[U] {
	return func(state State) Result[U] {
		r := p(state)
		if !r.OK {
			return Failure[U](r.Err, r.State)
		}
		return Success(fn(r.Value), r.State)
	}
}

// MapErr transforms the error of a failed parser using the provided function.
// Useful for adding context to error messages.
func MapErr[T any](p Parser[T], fn func(error) error) Parser[T] {
	return func(state State) Result[T] {
		r := p(state)
		if r.OK {
			return r
		}
		return Failure[T](fn(r.Err), r.State)
	}
}

// Label adds a descriptive name to a parser's error message.
// Wraps the original error with "expected <label>: <original error>".
//
// Example:
//
//	digit := Label(Range('0', '9'), "digit")
func Label[T any](p Parser[T], label string) Parser[T] {
	return MapErr(p, func(err error) error {
		return fmt.Errorf("expected %s: %w", label, err)
	})
}

// Skip runs a parser but discards its result, returning struct{}.
// The parser must still succeed.
func Skip[T any](p Parser[T]) Parser[struct{}] {
	return Map(p, func(_ T) struct{} { return struct{}{} })
}

// SkipMany matches zero or more occurrences, discarding all results.
// Always succeeds, returning struct{}.
func SkipMany[T any](p Parser[T]) Parser[struct{}] {
	return Skip(Many(p))
}

// SkipMany1 matches one or more occurrences, discarding all results.
// Fails if no matches are found.
func SkipMany1[T any](p Parser[T]) Parser[struct{}] {
	return Skip(Many1(p))
}

// Not inverts a parser's result: success becomes failure and vice versa.
// Does not consume input on success.
// Useful for negative lookahead patterns.
//
// Example:
//
//	// Match identifier that is not a keyword
//	notKeyword := And(Not(String("if")), Ident())
func Not[T any](p Parser[T]) Parser[struct{}] {
	return func(state State) Result[struct{}] {
		r := p(state)
		if r.OK {
			return Failure[struct{}](fmt.Errorf("unexpected match at line %d, col %d", state.Line, state.Col), state)
		}
		return Success(struct{}{}, state)
	}
}

// LookAhead tests a parser without consuming any input.
// Returns the matched value but leaves the position unchanged.
// Useful for conditional parsing based on what comes next.
func LookAhead[T any](p Parser[T]) Parser[T] {
	return func(state State) Result[T] {
		r := p(state)
		if r.OK {
			return Success(r.Value, state)
		}
		return r
	}
}

// Lazy defers parser evaluation by accepting a function that returns a parser.
// Enables mutual recursion between parsers.
func Lazy[T any](f func() Parser[T]) Parser[T] {
	return func(state State) Result[T] {
		return f()(state)
	}
}

// Ref creates a parser from a Rule pointer, enabling recursive grammar definitions.
// The Rule is evaluated each time the parser runs.
//
// Example:
//
//	var expr Rule[int64]
//	expr = func() Parser[int64] {
//		return Choice(Integer(), Parens(Ref(&expr)))
//	}
//	result := Parse(Ref(&expr), "((42))")
func Ref[T any](r *Rule[T]) Parser[T] {
	return func(state State) Result[T] {
		return (*r)()(state)
	}
}
