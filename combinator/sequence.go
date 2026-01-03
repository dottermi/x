package combinator

import "fmt"

// Seq runs parsers in sequence and returns a slice of all results.
// Fails immediately if any parser fails, returning that parser's error.
//
// Example:
//
//	// Match "a" followed by "b"
//	ab := Seq(Char('a'), Char('b'))
//	result := Parse(ab, "ab")
//	// result.Value == []any{'a', 'b'}
func Seq(parsers ...Parser) Parser {
	return func(state State) Result {
		var values []any
		current := state

		for _, p := range parsers {
			r := p(current)
			if !r.OK {
				return r
			}
			values = append(values, r.Value)
			current = r.State
		}

		return Success(values, current)
	}
}

// Choice tries parsers in order and returns the first successful result.
// Fails only if all alternatives fail, returning the last error.
//
// Example:
//
//	// Match "true" or "false"
//	boolean := Choice(String("true"), String("false"))
//	result := Parse(boolean, "true")
//	// result.Value == "true"
func Choice(parsers ...Parser) Parser {
	return func(state State) Result {
		var lastErr error

		for _, p := range parsers {
			r := p(state)
			if r.OK {
				return r
			}
			lastErr = r.Err
		}

		if lastErr != nil {
			return Failure(lastErr, state)
		}
		return Failure(fmt.Errorf("no alternatives matched at line %d, col %d", state.Line, state.Col), state)
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
//	// result.Value == []any{'1', '2', '3'}
func Many(p Parser) Parser {
	return func(state State) Result {
		var values []any
		current := state

		for {
			r := p(current)
			if !r.OK {
				break
			}
			values = append(values, r.Value)
			current = r.State

			// Evitar loop infinito se o parser não consumir nada
			if current.Pos == state.Pos {
				break
			}
			state = current
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
//	// result.Value == []any{'1', '2', '3'}
func Many1(p Parser) Parser {
	return func(state State) Result {
		first := p(state)
		if !first.OK {
			return first
		}

		rest := Many(p)(first.State)
		values := append([]any{first.Value}, rest.Value.([]any)...) //nolint:errcheck,forcetypeassert // type is guaranteed by Many

		return Success(values, rest.State)
	}
}

// Opt makes a parser optional, returning nil on failure instead of an error.
// Always succeeds.
//
// Example:
//
//	sign := Opt(Char('-'))
//	result := Parse(sign, "42")  // result.Value == nil
//	result = Parse(sign, "-42") // result.Value == '-'
func Opt(p Parser) Parser {
	return func(state State) Result {
		r := p(state)
		if r.OK {
			return r
		}
		return Success(nil, state)
	}
}

// And sequences two parsers, returning only the second parser's result.
// Both parsers must succeed.
func And(p1, p2 Parser) Parser {
	return func(state State) Result {
		r1 := p1(state)
		if !r1.OK {
			return r1
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
func Left(p1, p2 Parser) Parser {
	return func(state State) Result {
		r1 := p1(state)
		if !r1.OK {
			return r1
		}
		r2 := p2(r1.State)
		if !r2.OK {
			return r2
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
func Right(p1, p2 Parser) Parser {
	return func(state State) Result {
		r1 := p1(state)
		if !r1.OK {
			return r1
		}
		r2 := p2(r1.State)
		if !r2.OK {
			return r2
		}
		return Success(r2.Value, r2.State)
	}
}

// Count matches exactly n occurrences of a parser.
// Returns a slice of all matched values.
// Fails if fewer than n matches are found.
//
// Example:
//
//	hex := Count(2, HexDigit()) // match exactly 2 hex digits
func Count(n int, p Parser) Parser {
	return func(state State) Result {
		var values []any
		current := state

		for range n {
			r := p(current)
			if !r.OK {
				return r
			}
			values = append(values, r.Value)
			current = r.State
		}

		return Success(values, current)
	}
}
