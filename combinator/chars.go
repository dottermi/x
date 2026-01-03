package combinator

import "unicode"

// Digit matches a single Unicode digit character (0-9 and other Unicode digits).
// Returns the matched rune.
func Digit() Parser {
	return Label(Satisfy(unicode.IsDigit), "digit")
}

// Letter matches a single Unicode letter character.
// Returns the matched rune.
func Letter() Parser {
	return Label(Satisfy(unicode.IsLetter), "letter")
}

// Space matches a single Unicode whitespace character (space, tab, newline, etc.).
// Returns the matched rune.
func Space() Parser {
	return Label(Satisfy(unicode.IsSpace), "whitespace")
}

// Spaces matches zero or more whitespace characters.
// Returns the matched whitespace as a string.
// Always succeeds (returns empty string if no whitespace).
func Spaces() Parser {
	return Map(Many(Space()), func(v any) any {
		runes := v.([]any) //nolint:errcheck,forcetypeassert // type is guaranteed by parser
		result := make([]rune, len(runes))
		for i, r := range runes {
			result[i] = r.(rune) //nolint:errcheck,forcetypeassert // type is guaranteed by parser
		}
		return string(result)
	})
}

// Spaces1 matches one or more whitespace characters.
// Returns the matched whitespace as a string.
// Fails if no whitespace is present.
func Spaces1() Parser {
	return Map(Many1(Space()), func(v any) any {
		runes := v.([]any) //nolint:errcheck,forcetypeassert // type is guaranteed by parser
		result := make([]rune, len(runes))
		for i, r := range runes {
			result[i] = r.(rune) //nolint:errcheck,forcetypeassert // type is guaranteed by parser
		}
		return string(result)
	})
}

// Alpha matches a single ASCII letter (a-z, A-Z).
// Returns the matched rune.
// Use [Letter] for full Unicode letter support.
func Alpha() Parser {
	return Label(Satisfy(func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
	}), "letter")
}

// AlphaNum matches a single Unicode letter or digit.
// Returns the matched rune.
func AlphaNum() Parser {
	return Label(Satisfy(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	}), "alphanumeric")
}

// Lower matches a single Unicode lowercase letter.
// Returns the matched rune.
func Lower() Parser {
	return Label(Satisfy(unicode.IsLower), "lowercase letter")
}

// Upper matches a single Unicode uppercase letter.
// Returns the matched rune.
func Upper() Parser {
	return Label(Satisfy(unicode.IsUpper), "uppercase letter")
}

// Newline matches a single newline character ('\n').
// Returns the matched rune.
func Newline() Parser {
	return Label(Char('\n'), "newline")
}

// Tab matches a single tab character ('\t').
// Returns the matched rune.
func Tab() Parser {
	return Label(Char('\t'), "tab")
}

// CRLF matches the Windows line ending sequence "\r\n".
// Returns the matched string.
func CRLF() Parser {
	return Label(String("\r\n"), "CRLF")
}

// EndOfLine matches either Unix ('\n') or Windows ("\r\n") line endings.
// Returns the matched string.
func EndOfLine() Parser {
	return Label(Choice(CRLF(), Newline()), "end of line")
}

// HexDigit matches a single hexadecimal digit (0-9, a-f, A-F).
// Returns the matched rune.
func HexDigit() Parser {
	return Label(Satisfy(func(r rune) bool {
		return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
	}), "hex digit")
}

// OctDigit matches a single octal digit (0-7).
// Returns the matched rune.
func OctDigit() Parser {
	return Label(Satisfy(func(r rune) bool {
		return r >= '0' && r <= '7'
	}), "octal digit")
}

// BinDigit matches a single binary digit ('0' or '1').
// Returns the matched rune.
func BinDigit() Parser {
	return Label(OneOf("01"), "binary digit")
}
