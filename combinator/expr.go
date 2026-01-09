package combinator

// Lexeme wraps a parser to consume trailing whitespace after matching.
// Useful for building tokenizers that ignore whitespace between tokens.
//
// Example:
//
//	num := Lexeme(Integer())
//	result := Parse(Seq2(num, num), "42   17")
//	// Parses both numbers, ignoring whitespace between them
func Lexeme[T any](p Parser[T]) Parser[T] {
	return Left(p, Spaces())
}

// Symbol matches a string and consumes trailing whitespace.
// Convenience wrapper around Lexeme(String(s)).
//
// Example:
//
//	result := Parse(Symbol("if"), "if   ")
//	// result.Value == "if"
func Symbol(s string) Parser[string] {
	return Lexeme(String(s))
}

// Token matches an identifier and consumes trailing whitespace.
// Convenience wrapper around Lexeme(Ident()).
//
// Example:
//
//	result := Parse(Token(), "variableName   ")
//	// result.Value == "variableName"
func Token() Parser[string] {
	return Lexeme(Ident())
}

// IntToken matches an integer and consumes trailing whitespace.
// Convenience wrapper around Lexeme(Integer()).
//
// Example:
//
//	result := Parse(IntToken(), "123   ")
//	// result.Value == int64(123)
func IntToken() Parser[int64] {
	return Lexeme(Integer())
}

// FloatToken matches a float and consumes trailing whitespace.
// Convenience wrapper around Lexeme(Float()).
//
// Example:
//
//	result := Parse(FloatToken(), "3.14   ")
//	// result.Value == 3.14
func FloatToken() Parser[float64] {
	return Lexeme(Float())
}

// StringToken matches a string literal and consumes trailing whitespace.
// Convenience wrapper around Lexeme(StringLit()).
//
// Example:
//
//	result := Parse(StringToken(), `"hello world"   `)
//	// result.Value == "hello world"
func StringToken() Parser[string] {
	return Lexeme(StringLit())
}

// ChainL1 parses left-associative binary expressions like "1 + 2 + 3".
// The op parser must return a function that combines two values.
//
// Parameters:
//   - p: parser for the operands
//   - op: parser for the operator (must return a binary function)
//
// Example:
//
//	addOp := Map(Char('+'), func(_ rune) func(int64, int64) int64 {
//		return func(a, b int64) int64 { return a + b }
//	})
//	expr := ChainL1(Integer(), addOp)
//	result := Parse(expr, "1+2+3")
//	// result.Value == int64(6), computed as ((1+2)+3)
func ChainL1[T any](p Parser[T], op Parser[func(T, T) T]) Parser[T] {
	return func(state State) Result[T] {
		r := p(state)
		if !r.OK {
			return r
		}

		acc := r.Value
		current := r.State

		for {
			opResult := op(current)
			if !opResult.OK {
				break
			}

			nextResult := p(opResult.State)
			if !nextResult.OK {
				break
			}

			acc = opResult.Value(acc, nextResult.Value)
			current = nextResult.State
		}

		return Success(acc, current)
	}
}

// ChainR1 parses right-associative binary expressions like "2 ^ 3 ^ 4".
// The op parser must return a function that combines two values.
//
// Parameters:
//   - p: parser for the operands
//   - op: parser for the operator (must return a binary function)
//
// Example:
//
//	// For exponentiation: 2^3^4 = 2^(3^4)
//	powOp := Map(Char('^'), func(_ rune) func(float64, float64) float64 {
//		return func(a, b float64) float64 { return math.Pow(a, b) }
//	})
//	expr := ChainR1(Float(), powOp)
func ChainR1[T any](p Parser[T], op Parser[func(T, T) T]) Parser[T] {
	return func(state State) Result[T] {
		r := p(state)
		if !r.OK {
			return r
		}

		opResult := op(r.State)
		if !opResult.OK {
			return r
		}

		restResult := ChainR1(p, op)(opResult.State)
		if !restResult.OK {
			return r
		}

		return Success(opResult.Value(r.Value, restResult.Value), restResult.State)
	}
}
