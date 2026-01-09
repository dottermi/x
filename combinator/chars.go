package combinator

import "unicode"

// Digit matches a single Unicode digit character (0-9 and other Unicode digits).
// Returns the matched rune.
func Digit() Parser[rune] {
	return Label(Satisfy(unicode.IsDigit), "digit")
}

// Letter matches a single Unicode letter character.
// Returns the matched rune.
func Letter() Parser[rune] {
	return Label(Satisfy(unicode.IsLetter), "letter")
}

// Space matches a single Unicode whitespace character (space, tab, newline, etc.).
// Returns the matched rune.
func Space() Parser[rune] {
	return Label(Satisfy(unicode.IsSpace), "whitespace")
}

// Spaces matches zero or more whitespace characters.
// Returns the matched whitespace as a string.
// Always succeeds (returns empty string if no whitespace).
func Spaces() Parser[string] {
	return Map(Many(Space()), func(runes []rune) string {
		return string(runes)
	})
}

// Spaces1 matches one or more whitespace characters.
// Returns the matched whitespace as a string.
// Fails if no whitespace is present.
func Spaces1() Parser[string] {
	return Map(Many1(Space()), func(runes []rune) string {
		return string(runes)
	})
}

// Alpha matches a single ASCII letter (a-z, A-Z).
// Returns the matched rune.
// Use [Letter] for full Unicode letter support.
func Alpha() Parser[rune] {
	return Label(Satisfy(func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
	}), "letter")
}

// AlphaNum matches a single Unicode letter or digit.
// Returns the matched rune.
func AlphaNum() Parser[rune] {
	return Label(Satisfy(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	}), "alphanumeric")
}

// Lower matches a single Unicode lowercase letter.
// Returns the matched rune.
func Lower() Parser[rune] {
	return Label(Satisfy(unicode.IsLower), "lowercase letter")
}

// Upper matches a single Unicode uppercase letter.
// Returns the matched rune.
func Upper() Parser[rune] {
	return Label(Satisfy(unicode.IsUpper), "uppercase letter")
}

// Newline matches a single newline character ('\n').
// Returns the matched rune.
func Newline() Parser[rune] {
	return Label(Char('\n'), "newline")
}

// Tab matches a single tab character ('\t').
// Returns the matched rune.
func Tab() Parser[rune] {
	return Label(Char('\t'), "tab")
}

// CRLF matches the Windows line ending sequence "\r\n".
// Returns the matched string.
func CRLF() Parser[string] {
	return Label(String("\r\n"), "CRLF")
}

// EndOfLine matches either Unix ('\n') or Windows ("\r\n") line endings.
// Returns the matched string.
func EndOfLine() Parser[string] {
	return Label(Choice(CRLF(), Map(Newline(), func(r rune) string { return string(r) })), "end of line")
}

// HexDigit matches a single hexadecimal digit (0-9, a-f, A-F).
// Returns the matched rune.
func HexDigit() Parser[rune] {
	return Label(Satisfy(func(r rune) bool {
		return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
	}), "hex digit")
}

// OctDigit matches a single octal digit (0-7).
// Returns the matched rune.
func OctDigit() Parser[rune] {
	return Label(Satisfy(func(r rune) bool {
		return r >= '0' && r <= '7'
	}), "octal digit")
}

// BinDigit matches a single binary digit ('0' or '1').
// Returns the matched rune.
func BinDigit() Parser[rune] {
	return Label(OneOf("01"), "binary digit")
}
