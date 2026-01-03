package combinator

import (
	"strconv"
	"strings"
)

// Ident matches a programming language identifier.
// Accepts a letter or underscore followed by any combination of letters, digits, or underscores.
// Returns the identifier as a string.
//
// Example:
//
//	result := Parse(Ident(), "myVar123")
//	// result.Value == "myVar123"
func Ident() Parser {
	first := Choice(Letter(), Char('_'))
	rest := Many(Choice(AlphaNum(), Char('_')))

	return Map(Seq(first, rest), func(v any) any {
		parts := v.([]any)
		firstChar := parts[0].(rune)
		restChars := parts[1].([]any)

		var sb strings.Builder
		sb.WriteRune(firstChar)
		for _, r := range restChars {
			sb.WriteRune(r.(rune))
		}
		return sb.String()
	})
}

// Keyword matches a specific keyword that is not followed by alphanumeric characters.
// Prevents matching "if" in "iffy" by requiring a word boundary.
//
// Example:
//
//	result := Parse(Keyword("if"), "if (x)")  // succeeds
//	result = Parse(Keyword("if"), "iffy")    // fails
func Keyword(kw string) Parser {
	return Left(String(kw), Not(AlphaNum()))
}

// Integer matches an optionally negative integer and returns it as int64.
// Accepts an optional leading minus sign followed by one or more digits.
//
// Example:
//
//	result := Parse(Integer(), "-42")
//	// result.Value == int64(-42)
func Integer() Parser {
	sign := Opt(Char('-'))
	digits := Many1(Digit())

	return Map(Seq(sign, digits), func(v any) any {
		parts := v.([]any)
		neg := parts[0]
		digitRunes := parts[1].([]any)

		var sb strings.Builder
		if neg != nil {
			sb.WriteRune(neg.(rune))
		}
		for _, r := range digitRunes {
			sb.WriteRune(r.(rune))
		}

		n, _ := strconv.ParseInt(sb.String(), 10, 64)
		return n
	})
}

// Float matches a decimal number and returns it as float64.
// Accepts an optional sign, integer part, and optional decimal part.
//
// Example:
//
//	result := Parse(Float(), "-3.14")
//	// result.Value == float64(-3.14)
func Float() Parser {
	sign := Opt(Char('-'))
	intPart := Many1(Digit())
	decPart := Opt(Seq(Char('.'), Many1(Digit())))

	return Map(Seq(sign, intPart, decPart), func(v any) any {
		parts := v.([]any)
		neg := parts[0]
		intDigits := parts[1].([]any)
		dec := parts[2]

		var sb strings.Builder
		if neg != nil {
			sb.WriteRune(neg.(rune))
		}
		for _, r := range intDigits {
			sb.WriteRune(r.(rune))
		}

		if dec != nil {
			decParts := dec.([]any)
			sb.WriteRune('.')
			for _, r := range decParts[1].([]any) {
				sb.WriteRune(r.(rune))
			}
		}

		f, _ := strconv.ParseFloat(sb.String(), 64)
		return f
	})
}

// StringLit matches a double-quoted string with basic escape support.
// Returns the string contents without the surrounding quotes.
// Supports backslash escapes (e.g., \" and \\).
//
// Example:
//
//	result := Parse(StringLit(), `"hello \"world\""`)
//	// result.Value == `hello "world"`
func StringLit() Parser {
	quote := Char('"')
	escaped := Right(Char('\\'), Any())
	regular := Satisfy(func(r rune) bool {
		return r != '"' && r != '\\'
	})
	content := Many(Choice(escaped, regular))

	return Map(Seq(quote, content, quote), func(v any) any {
		parts := v.([]any)
		chars := parts[1].([]any)

		var sb strings.Builder
		for _, r := range chars {
			sb.WriteRune(r.(rune))
		}
		return sb.String()
	})
}

// CharLit matches a single-quoted character literal with escape support.
// Returns the character as a rune.
//
// Example:
//
//	result := Parse(CharLit(), `'a'`)
//	// result.Value == 'a'
func CharLit() Parser {
	quote := Char('\'')
	escaped := Right(Char('\\'), Any())
	regular := Satisfy(func(r rune) bool {
		return r != '\'' && r != '\\'
	})

	return Map(Seq(quote, Choice(escaped, regular), quote), func(v any) any {
		parts := v.([]any)
		return parts[1].(rune)
	})
}
