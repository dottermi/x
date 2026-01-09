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
func Ident() Parser[string] {
	first := Choice(Letter(), Char('_'))
	rest := Many(Choice(AlphaNum(), Char('_')))

	return Map(Seq2(first, rest), func(p Pair[rune, []rune]) string {
		var sb strings.Builder
		sb.WriteRune(p.First)
		for _, r := range p.Second {
			sb.WriteRune(r)
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
func Keyword(kw string) Parser[string] {
	return Left(String(kw), Not(AlphaNum()))
}

// Integer matches an optionally negative integer and returns it as int64.
// Accepts an optional leading minus sign followed by one or more digits.
//
// Example:
//
//	result := Parse(Integer(), "-42")
//	// result.Value == int64(-42)
func Integer() Parser[int64] {
	sign := Opt(Char('-'))
	digits := Many1(Digit())

	return Map(Seq2(sign, digits), func(p Pair[*rune, []rune]) int64 {
		var sb strings.Builder
		if p.First != nil {
			sb.WriteRune(*p.First)
		}
		for _, r := range p.Second {
			sb.WriteRune(r)
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
func Float() Parser[float64] {
	sign := Opt(Char('-'))
	intPart := Many1(Digit())
	decPart := Opt(Seq2(Char('.'), Many1(Digit())))

	return Map(Seq3(sign, intPart, decPart), func(t Triple[*rune, []rune, *Pair[rune, []rune]]) float64 {
		var sb strings.Builder
		if t.First != nil {
			sb.WriteRune(*t.First)
		}
		for _, r := range t.Second {
			sb.WriteRune(r)
		}
		if t.Third != nil {
			sb.WriteRune('.')
			for _, r := range t.Third.Second {
				sb.WriteRune(r)
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
func StringLit() Parser[string] {
	quote := Char('"')
	escaped := Right(Char('\\'), Any())
	regular := Satisfy(func(r rune) bool {
		return r != '"' && r != '\\'
	})
	content := Many(Choice(escaped, regular))

	return Map(Seq3(quote, content, quote), func(t Triple[rune, []rune, rune]) string {
		var sb strings.Builder
		for _, r := range t.Second {
			sb.WriteRune(r)
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
func CharLit() Parser[rune] {
	quote := Char('\'')
	escaped := Right(Char('\\'), Any())
	regular := Satisfy(func(r rune) bool {
		return r != '\'' && r != '\\'
	})

	return Map(Seq3(quote, Choice(escaped, regular), quote), func(t Triple[rune, rune, rune]) rune {
		return t.Second
	})
}
