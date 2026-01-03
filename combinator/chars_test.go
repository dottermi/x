package combinator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest // tests share parser state
func TestDigit(t *testing.T) {
	t.Run("should match digit", func(t *testing.T) {
		result := Parse(Digit(), "5abc")
		assert.True(t, result.OK)
		assert.Equal(t, '5', result.Value)
	})

	t.Run("should fail on non-digit", func(t *testing.T) {
		result := Parse(Digit(), "abc")
		assert.False(t, result.OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestLetter(t *testing.T) {
	t.Run("should match lowercase", func(t *testing.T) {
		result := Parse(Letter(), "abc")
		assert.True(t, result.OK)
		assert.Equal(t, 'a', result.Value)
	})

	t.Run("should match uppercase", func(t *testing.T) {
		result := Parse(Letter(), "ABC")
		assert.True(t, result.OK)
		assert.Equal(t, 'A', result.Value)
	})

	t.Run("should fail on digit", func(t *testing.T) {
		assert.False(t, Parse(Letter(), "123").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestSpace(t *testing.T) {
	t.Run("should match space", func(t *testing.T) {
		result := Parse(Space(), " abc")
		assert.True(t, result.OK)
		assert.Equal(t, ' ', result.Value)
	})

	t.Run("should match tab", func(t *testing.T) {
		result := Parse(Space(), "\tabc")
		assert.True(t, result.OK)
	})

	t.Run("should match newline", func(t *testing.T) {
		result := Parse(Space(), "\nabc")
		assert.True(t, result.OK)
	})

	t.Run("should fail on non-whitespace", func(t *testing.T) {
		assert.False(t, Parse(Space(), "abc").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestSpaces(t *testing.T) {
	t.Run("should match multiple spaces", func(t *testing.T) {
		result := Parse(Spaces(), "   abc")
		assert.True(t, result.OK)
		assert.Equal(t, "   ", result.Value)
	})

	t.Run("should succeed with empty on no whitespace", func(t *testing.T) {
		result := Parse(Spaces(), "abc")
		assert.True(t, result.OK)
		assert.Empty(t, result.Value)
	})

	t.Run("should match mixed whitespace", func(t *testing.T) {
		result := Parse(Spaces(), " \t\n abc")
		assert.True(t, result.OK)
		assert.Equal(t, " \t\n ", result.Value)
	})
}

//nolint:paralleltest // tests share parser state
func TestSpaces1(t *testing.T) {
	t.Run("should match one or more", func(t *testing.T) {
		result := Parse(Spaces1(), "  abc")
		assert.True(t, result.OK)
		assert.Equal(t, "  ", result.Value)
	})

	t.Run("should fail on no whitespace", func(t *testing.T) {
		assert.False(t, Parse(Spaces1(), "abc").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestAlpha(t *testing.T) {
	t.Run("should match ASCII letters", func(t *testing.T) {
		assert.True(t, Parse(Alpha(), "abc").OK)
		assert.True(t, Parse(Alpha(), "ABC").OK)
	})

	t.Run("should fail on digit", func(t *testing.T) {
		assert.False(t, Parse(Alpha(), "123").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestAlphaNum(t *testing.T) {
	t.Run("should match letter", func(t *testing.T) {
		assert.True(t, Parse(AlphaNum(), "abc").OK)
	})

	t.Run("should match digit", func(t *testing.T) {
		assert.True(t, Parse(AlphaNum(), "123").OK)
	})

	t.Run("should fail on special character", func(t *testing.T) {
		assert.False(t, Parse(AlphaNum(), "!@#").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestLower(t *testing.T) {
	t.Run("should match lowercase", func(t *testing.T) {
		assert.True(t, Parse(Lower(), "abc").OK)
	})

	t.Run("should fail on uppercase", func(t *testing.T) {
		assert.False(t, Parse(Lower(), "ABC").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestUpper(t *testing.T) {
	t.Run("should match uppercase", func(t *testing.T) {
		assert.True(t, Parse(Upper(), "ABC").OK)
	})

	t.Run("should fail on lowercase", func(t *testing.T) {
		assert.False(t, Parse(Upper(), "abc").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestNewline(t *testing.T) {
	t.Run("should match newline", func(t *testing.T) {
		result := Parse(Newline(), "\nabc")
		assert.True(t, result.OK)
		assert.Equal(t, '\n', result.Value)
	})

	t.Run("should fail on other", func(t *testing.T) {
		assert.False(t, Parse(Newline(), "abc").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestTab(t *testing.T) {
	t.Run("should match tab", func(t *testing.T) {
		result := Parse(Tab(), "\tabc")
		assert.True(t, result.OK)
		assert.Equal(t, '\t', result.Value)
	})

	t.Run("should fail on other", func(t *testing.T) {
		assert.False(t, Parse(Tab(), "abc").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestCRLF(t *testing.T) {
	t.Run("should match CRLF", func(t *testing.T) {
		result := Parse(CRLF(), "\r\nabc")
		assert.True(t, result.OK)
		assert.Equal(t, "\r\n", result.Value)
	})

	t.Run("should fail on just LF", func(t *testing.T) {
		assert.False(t, Parse(CRLF(), "\nabc").OK)
	})

	t.Run("should fail on just CR", func(t *testing.T) {
		assert.False(t, Parse(CRLF(), "\rabc").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestEndOfLine(t *testing.T) {
	t.Run("should match Unix ending", func(t *testing.T) {
		assert.True(t, Parse(EndOfLine(), "\nabc").OK)
	})

	t.Run("should match Windows ending", func(t *testing.T) {
		result := Parse(EndOfLine(), "\r\nabc")
		assert.True(t, result.OK)
		assert.Equal(t, "\r\n", result.Value)
	})
}

//nolint:paralleltest // tests share parser state
func TestHexDigit(t *testing.T) {
	t.Run("should match 0-9", func(t *testing.T) {
		assert.True(t, Parse(HexDigit(), "9").OK)
	})

	t.Run("should match a-f", func(t *testing.T) {
		assert.True(t, Parse(HexDigit(), "f").OK)
	})

	t.Run("should match A-F", func(t *testing.T) {
		assert.True(t, Parse(HexDigit(), "F").OK)
	})

	t.Run("should fail on g", func(t *testing.T) {
		assert.False(t, Parse(HexDigit(), "g").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestOctDigit(t *testing.T) {
	t.Run("should match 0-7", func(t *testing.T) {
		assert.True(t, Parse(OctDigit(), "7").OK)
	})

	t.Run("should fail on 8", func(t *testing.T) {
		assert.False(t, Parse(OctDigit(), "8").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestBinDigit(t *testing.T) {
	t.Run("should match 0 and 1", func(t *testing.T) {
		assert.True(t, Parse(BinDigit(), "0").OK)
		assert.True(t, Parse(BinDigit(), "1").OK)
	})

	t.Run("should fail on 2", func(t *testing.T) {
		assert.False(t, Parse(BinDigit(), "2").OK)
	})
}
