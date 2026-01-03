package combinator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest // tests share parser state
func TestIdent(t *testing.T) {
	t.Run("should match simple identifier", func(t *testing.T) {
		result := Parse(Ident(), "myVar")
		assert.True(t, result.OK)
		assert.Equal(t, "myVar", result.Value)
	})

	t.Run("should match with underscore prefix", func(t *testing.T) {
		result := Parse(Ident(), "_private")
		assert.True(t, result.OK)
		assert.Equal(t, "_private", result.Value)
	})

	t.Run("should match with digits", func(t *testing.T) {
		result := Parse(Ident(), "var123")
		assert.True(t, result.OK)
		assert.Equal(t, "var123", result.Value)
	})

	t.Run("should fail when starting with digit", func(t *testing.T) {
		assert.False(t, Parse(Ident(), "123var").OK)
	})

	t.Run("should match single character", func(t *testing.T) {
		result := Parse(Ident(), "x")
		assert.True(t, result.OK)
		assert.Equal(t, "x", result.Value)
	})

	t.Run("should match single underscore", func(t *testing.T) {
		result := Parse(Ident(), "_")
		assert.True(t, result.OK)
		assert.Equal(t, "_", result.Value)
	})
}

//nolint:paralleltest // tests share parser state
func TestKeyword(t *testing.T) {
	t.Run("should match keyword followed by non-alphanumeric", func(t *testing.T) {
		result := Parse(Keyword("if"), "if (x)")
		assert.True(t, result.OK)
		assert.Equal(t, "if", result.Value)
	})

	t.Run("should match keyword at end of input", func(t *testing.T) {
		result := Parse(Keyword("if"), "if")
		assert.True(t, result.OK)
	})

	t.Run("should fail when keyword is prefix", func(t *testing.T) {
		assert.False(t, Parse(Keyword("if"), "iffy").OK)
	})

	t.Run("should fail when followed by digit", func(t *testing.T) {
		assert.False(t, Parse(Keyword("var"), "var123").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestInteger(t *testing.T) {
	t.Run("should match positive", func(t *testing.T) {
		result := Parse(Integer(), "42")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})

	t.Run("should match negative", func(t *testing.T) {
		result := Parse(Integer(), "-42")
		assert.True(t, result.OK)
		assert.Equal(t, int64(-42), result.Value)
	})

	t.Run("should match zero", func(t *testing.T) {
		result := Parse(Integer(), "0")
		assert.True(t, result.OK)
		assert.Equal(t, int64(0), result.Value)
	})

	t.Run("should fail on non-digit", func(t *testing.T) {
		assert.False(t, Parse(Integer(), "abc").OK)
	})

	t.Run("should stop at non-digit", func(t *testing.T) {
		result := Parse(Integer(), "123abc")
		assert.True(t, result.OK)
		assert.Equal(t, int64(123), result.Value)
		assert.Equal(t, 3, result.State.Pos)
	})
}

//nolint:paralleltest // tests share parser state
func TestFloat(t *testing.T) {
	t.Run("should match simple float", func(t *testing.T) {
		result := Parse(Float(), "3.14")
		assert.True(t, result.OK)
		assert.Equal(t, float64(3.14), result.Value)
	})

	t.Run("should match negative float", func(t *testing.T) {
		result := Parse(Float(), "-3.14")
		assert.True(t, result.OK)
		assert.Equal(t, float64(-3.14), result.Value)
	})

	t.Run("should match integer as float", func(t *testing.T) {
		result := Parse(Float(), "42")
		assert.True(t, result.OK)
		assert.Equal(t, float64(42), result.Value)
	})

	t.Run("should match with leading zero", func(t *testing.T) {
		result := Parse(Float(), "0.01")
		assert.True(t, result.OK)
		assert.Equal(t, float64(0.01), result.Value)
	})
}

//nolint:paralleltest // tests share parser state
func TestStringLit(t *testing.T) {
	t.Run("should match simple string", func(t *testing.T) {
		result := Parse(StringLit(), `"hello"`)
		assert.True(t, result.OK)
		assert.Equal(t, "hello", result.Value)
	})

	t.Run("should match empty string", func(t *testing.T) {
		result := Parse(StringLit(), `""`)
		assert.True(t, result.OK)
		assert.Empty(t, result.Value)
	})

	t.Run("should handle escaped quote", func(t *testing.T) {
		result := Parse(StringLit(), `"hello \"world\""`)
		assert.True(t, result.OK)
		assert.Equal(t, `hello "world"`, result.Value)
	})

	t.Run("should handle escaped backslash", func(t *testing.T) {
		result := Parse(StringLit(), `"path\\to\\file"`)
		assert.True(t, result.OK)
		assert.Equal(t, `path\to\file`, result.Value)
	})

	t.Run("should fail on unclosed string", func(t *testing.T) {
		assert.False(t, Parse(StringLit(), `"unclosed`).OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestCharLit(t *testing.T) {
	t.Run("should match simple character", func(t *testing.T) {
		result := Parse(CharLit(), `'a'`)
		assert.True(t, result.OK)
		assert.Equal(t, 'a', result.Value)
	})

	t.Run("should handle escaped quote", func(t *testing.T) {
		result := Parse(CharLit(), `'\''`)
		assert.True(t, result.OK)
		assert.Equal(t, '\'', result.Value)
	})

	t.Run("should handle escaped backslash", func(t *testing.T) {
		result := Parse(CharLit(), `'\\'`)
		assert.True(t, result.OK)
		assert.Equal(t, '\\', result.Value)
	})

	t.Run("should fail on empty char literal", func(t *testing.T) {
		assert.False(t, Parse(CharLit(), `''`).OK)
	})
}
