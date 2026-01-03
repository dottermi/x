package combinator

import "fmt"

// Map transforms the result of a parser using the provided function.
// The function receives the parsed value and returns a new value.
//
// Example:
//
//	upper := Map(Letter(), func(v any) any {
//		return unicode.ToUpper(v.(rune))
//	})
func Map(p Parser, fn func(any) any) Parser {
	return func(state State) Result {
		r := p(state)
		if !r.OK {
			return r
		}
		return Success(fn(r.Value), r.State)
	}
}

// MapErr transforms the error of a failed parser using the provided function.
// Useful for adding context to error messages.
func MapErr(p Parser, fn func(error) error) Parser {
	return func(state State) Result {
		r := p(state)
		if r.OK {
			return r
		}
		return Failure(fn(r.Err), r.State)
	}
}

// Label adds a descriptive name to a parser's error message.
// Wraps the original error with "expected <label>: <original error>".
//
// Example:
//
//	digit := Label(Range('0', '9'), "digit")
func Label(p Parser, label string) Parser {
	return MapErr(p, func(err error) error {
		return fmt.Errorf("expected %s: %w", label, err)
	})
}

// Skip runs a parser but discards its result, returning nil.
// The parser must still succeed.
func Skip(p Parser) Parser {
	return Map(p, func(_ any) any { return nil })
}

// SkipMany matches zero or more occurrences, discarding all results.
// Always succeeds, returning nil.
func SkipMany(p Parser) Parser {
	return Skip(Many(p))
}

// SkipMany1 matches one or more occurrences, discarding all results.
// Fails if no matches are found.
func SkipMany1(p Parser) Parser {
	return Skip(Many1(p))
}

// Not inverts a parser's result: success becomes failure and vice versa.
// Does not consume input on success.
// Useful for negative lookahead patterns.
//
// Example:
//
//	// Match identifier that is not a keyword
//	notKeyword := Seq(Not(String("if")), Ident())
func Not(p Parser) Parser {
	return func(state State) Result {
		r := p(state)
		if r.OK {
			return Failure(fmt.Errorf("unexpected match at line %d, col %d", state.Line, state.Col), state)
		}
		return Success(nil, state)
	}
}

// LookAhead tests a parser without consuming any input.
// Returns the matched value but leaves the position unchanged.
// Useful for conditional parsing based on what comes next.
func LookAhead(p Parser) Parser {
	return func(state State) Result {
		r := p(state)
		if r.OK {
			return Success(r.Value, state)
		}
		return r
	}
}

// Lazy defers parser evaluation by accepting a pointer to a Parser.
// Enables mutual recursion between parsers.
// For self-referential grammars, prefer using [Rule] and [Ref].
func Lazy(p *Parser) Parser {
	return func(state State) Result {
		return (*p)(state)
	}
}

// Ref creates a parser from a Rule pointer, enabling recursive grammar definitions.
// The Rule is evaluated each time the parser runs.
//
// Example:
//
//	var expr Rule
//	expr = func() Parser {
//		return Choice(Integer(), Parens(Ref(&expr)))
//	}
//	result := Parse(Ref(&expr), "((42))")
func Ref(r *Rule) Parser {
	return func(state State) Result {
		return (*r)()(state)
	}
}
