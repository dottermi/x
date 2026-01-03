package combinator

// Between matches a parser surrounded by opening and closing delimiters.
// Returns only the middle parser's result; delimiters are discarded.
//
// Parameters:
//   - open: parser for the opening delimiter
//   - close: parser for the closing delimiter
//   - p: the main content parser
//
// Example:
//
//	quoted := Between(Char('"'), Char('"'), Many(NoneOf("\"")))
func Between(open, close, p Parser) Parser {
	return Map(Seq(open, p, close), func(v any) any {
		return v.([]any)[1]
	})
}

// Parens matches a parser surrounded by parentheses.
// Returns the inner parser's result.
// Example:
//
//	expression := Parens(Choice(Integer(), AddExpr()))
//	result := Parse(expression, "(1 + 2)")
//	// result.Value == ... // result of the inner expression
func Parens(p Parser) Parser {
	return Between(Char('('), Char(')'), p)
}

// Braces matches a parser surrounded by curly braces.
// Returns the inner parser's result.
// Example:
//
//	object := Braces(SepBy(KeyValuePair(), Char(',')))
//	result := Parse(object, "{key1: value1, key2: value2}")
//	// result.Value == []any{...} // slice of key-value pairs
func Braces(p Parser) Parser {
	return Between(Char('{'), Char('}'), p)
}

// Brackets matches a parser surrounded by square brackets.
// Returns the inner parser's result.
//
// Example:
//
//	array := Brackets(SepBy(Integer(), Char(',')))
//	result := Parse(array, "[1,2,3]")
//	// result.Value == []any{int64(1), int64(2), int64(3)}
func Brackets(p Parser) Parser {
	return Between(Char('['), Char(']'), p)
}

// Angles matches a parser surrounded by angle brackets.
// Returns the inner parser's result.
//
// Example:
//
//	tag := Angles(Ident())
//	result := Parse(tag, "<html>")
//	// result.Value == "html"
func Angles(p Parser) Parser {
	return Between(Char('<'), Char('>'), p)
}

// SepBy matches zero or more occurrences of a parser separated by a delimiter.
// Returns a slice of the matched values (separators are discarded).
// Always succeeds (returns empty slice if no matches).
//
// Example:
//
//	items := SepBy(Integer(), Char(','))
//	result := Parse(items, "1,2,3")
//	// result.Value == []any{int64(1), int64(2), int64(3)}
func SepBy(p, sep Parser) Parser {
	return Choice(SepBy1(p, sep), Map(EOF(), func(_ any) any { return []any{} }))
}

// SepBy1 matches one or more occurrences of a parser separated by a delimiter.
// Returns a slice of the matched values (separators are discarded).
// Fails if no matches are found.
//
// Example:
//
//	items := SepBy1(Ident(), Char(','))
//	result := Parse(items, "a,b,c")
//	// result.Value == []any{"a", "b", "c"}
func SepBy1(p, sep Parser) Parser {
	rest := Many(Right(sep, p))

	return Map(Seq(p, rest), func(v any) any {
		parts := v.([]any)
		first := parts[0]
		restItems := parts[1].([]any)
		return append([]any{first}, restItems...)
	})
}

// EndBy matches zero or more occurrences of a parser, each followed by a terminator.
// Returns a slice of the matched values (terminators are discarded).
//
// Example:
//
//	statements := EndBy(Ident(), Char(';'))
//	result := Parse(statements, "a;b;c;")
//	// result.Value == []any{"a", "b", "c"}
func EndBy(p, end Parser) Parser {
	return Many(Left(p, end))
}

// EndBy1 matches one or more occurrences of a parser, each followed by a terminator.
// Returns a slice of the matched values (terminators are discarded).
// Fails if no matches are found.
// Example:
//
//	statements := EndBy1(Ident(), Char(';'))
//	result := Parse(statements, "a;b;c;")
//	// result.Value == []any{"a", "b", "c"}
func EndBy1(p, end Parser) Parser {
	return Many1(Left(p, end))
}
