package combinator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest // tests share parser state
func TestChar(t *testing.T) {
	t.Run("should match expected character", func(t *testing.T) {
		result := Parse(Char('a'), "abc")
		assert.True(t, result.OK)
		assert.Equal(t, 'a', result.Value)
		assert.Equal(t, 1, result.State.Pos)
	})

	t.Run("should fail on different character", func(t *testing.T) {
		result := Parse(Char('a'), "xyz")
		assert.False(t, result.OK)
		assert.Contains(t, result.Err.Error(), "expected 'a'")
	})

	t.Run("should fail on EOF", func(t *testing.T) {
		result := Parse(Char('a'), "")
		assert.False(t, result.OK)
		assert.Contains(t, result.Err.Error(), "unexpected EOF")
	})

	t.Run("should match special characters", func(t *testing.T) {
		result := Parse(Char('\n'), "\nabc")
		assert.True(t, result.OK)
		assert.Equal(t, '\n', result.Value)
	})
}

//nolint:paralleltest // tests share parser state
func TestString(t *testing.T) {
	t.Run("should match exact string", func(t *testing.T) {
		result := Parse(String("hello"), "hello world")
		assert.True(t, result.OK)
		assert.Equal(t, "hello", result.Value)
		assert.Equal(t, 5, result.State.Pos)
	})

	t.Run("should fail on partial match", func(t *testing.T) {
		result := Parse(String("hello"), "help")
		assert.False(t, result.OK)
	})

	t.Run("should fail on EOF during match", func(t *testing.T) {
		result := Parse(String("hello"), "hel")
		assert.False(t, result.OK)
	})

	t.Run("should match empty string", func(t *testing.T) {
		result := Parse(String(""), "anything")
		assert.True(t, result.OK)
		assert.Empty(t, result.Value)
	})
}

//nolint:paralleltest // tests share parser state
func TestSatisfy(t *testing.T) {
	t.Run("should match when predicate is true", func(t *testing.T) {
		result := Parse(Satisfy(func(r rune) bool { return r >= 'a' && r <= 'z' }), "abc")
		assert.True(t, result.OK)
		assert.Equal(t, 'a', result.Value)
	})

	t.Run("should fail when predicate is false", func(t *testing.T) {
		result := Parse(Satisfy(func(r rune) bool { return r >= '0' && r <= '9' }), "abc")
		assert.False(t, result.OK)
	})

	t.Run("should fail on EOF", func(t *testing.T) {
		result := Parse(Satisfy(func(r rune) bool { return true }), "")
		assert.False(t, result.OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestAny(t *testing.T) {
	t.Run("should match any character", func(t *testing.T) {
		result := Parse(Any(), "x")
		assert.True(t, result.OK)
		assert.Equal(t, 'x', result.Value)
	})

	t.Run("should fail on EOF", func(t *testing.T) {
		result := Parse(Any(), "")
		assert.False(t, result.OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestEOF(t *testing.T) {
	t.Run("should succeed at end of input", func(t *testing.T) {
		result := Parse(EOF(), "")
		assert.True(t, result.OK)
	})

	t.Run("should fail when input remains", func(t *testing.T) {
		result := Parse(EOF(), "remaining")
		assert.False(t, result.OK)
		assert.Contains(t, result.Err.Error(), "expected EOF")
	})
}

//nolint:paralleltest // tests share parser state
func TestOneOf(t *testing.T) {
	t.Run("should match character in set", func(t *testing.T) {
		result := Parse(OneOf("aeiou"), "apple")
		assert.True(t, result.OK)
		assert.Equal(t, 'a', result.Value)
	})

	t.Run("should fail for character not in set", func(t *testing.T) {
		result := Parse(OneOf("aeiou"), "xyz")
		assert.False(t, result.OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestNoneOf(t *testing.T) {
	t.Run("should match character not in set", func(t *testing.T) {
		result := Parse(NoneOf("aeiou"), "xyz")
		assert.True(t, result.OK)
		assert.Equal(t, 'x', result.Value)
	})

	t.Run("should fail for character in set", func(t *testing.T) {
		result := Parse(NoneOf("aeiou"), "apple")
		assert.False(t, result.OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestRange(t *testing.T) {
	t.Run("should match character in range", func(t *testing.T) {
		result := Parse(Range('a', 'z'), "m")
		assert.True(t, result.OK)
		assert.Equal(t, 'm', result.Value)
	})

	t.Run("should match bounds", func(t *testing.T) {
		assert.True(t, Parse(Range('a', 'z'), "a").OK)
		assert.True(t, Parse(Range('a', 'z'), "z").OK)
	})

	t.Run("should fail outside range", func(t *testing.T) {
		result := Parse(Range('a', 'z'), "A")
		assert.False(t, result.OK)
	})
}
