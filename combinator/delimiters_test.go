package combinator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // tests share parser state
func TestBetween(t *testing.T) {
	t.Run("should match content between delimiters", func(t *testing.T) {
		result := Parse(Between(Char('['), Char(']'), Integer()), "[42]")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})

	t.Run("should fail on missing opening", func(t *testing.T) {
		assert.False(t, Parse(Between(Char('['), Char(']'), Integer()), "42]").OK)
	})

	t.Run("should fail on missing closing", func(t *testing.T) {
		assert.False(t, Parse(Between(Char('['), Char(']'), Integer()), "[42").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestParens(t *testing.T) {
	t.Run("should match content in parentheses", func(t *testing.T) {
		result := Parse(Parens(Integer()), "(42)")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})

	t.Run("should fail on mismatched", func(t *testing.T) {
		assert.False(t, Parse(Parens(Integer()), "(42]").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestBraces(t *testing.T) {
	t.Run("should match content in braces", func(t *testing.T) {
		result := Parse(Braces(Integer()), "{42}")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})
}

//nolint:paralleltest // tests share parser state
func TestBrackets(t *testing.T) {
	t.Run("should match content in brackets", func(t *testing.T) {
		result := Parse(Brackets(Integer()), "[42]")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})
}

//nolint:paralleltest // tests share parser state
func TestAngles(t *testing.T) {
	t.Run("should match content in angle brackets", func(t *testing.T) {
		result := Parse(Angles(Ident()), "<html>")
		assert.True(t, result.OK)
		assert.Equal(t, "html", result.Value)
	})
}

//nolint:paralleltest // tests share parser state
func TestSepBy(t *testing.T) {
	t.Run("should match items separated by delimiter", func(t *testing.T) {
		result := Parse(SepBy(Integer(), Char(',')), "1,2,3")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 3)
		assert.Equal(t, int64(1), result.Value[0])
		assert.Equal(t, int64(2), result.Value[1])
		assert.Equal(t, int64(3), result.Value[2])
	})

	t.Run("should match single item", func(t *testing.T) {
		result := Parse(SepBy(Integer(), Char(',')), "42")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 1)
	})

	t.Run("should succeed with empty input", func(t *testing.T) {
		result := Parse(SepBy(Integer(), Char(',')), "")
		require.True(t, result.OK)
		assert.Empty(t, result.Value)
	})
}

//nolint:paralleltest // tests share parser state
func TestSepBy1(t *testing.T) {
	t.Run("should match items separated by delimiter", func(t *testing.T) {
		result := Parse(SepBy1(Ident(), Char(',')), "a,b,c")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 3)
	})

	t.Run("should match single item", func(t *testing.T) {
		result := Parse(SepBy1(Ident(), Char(',')), "x")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 1)
		assert.Equal(t, "x", result.Value[0])
	})

	t.Run("should fail on empty input", func(t *testing.T) {
		assert.False(t, Parse(SepBy1(Integer(), Char(',')), "").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestEndBy(t *testing.T) {
	t.Run("should match items followed by terminator", func(t *testing.T) {
		result := Parse(EndBy(Ident(), Char(';')), "a;b;c;")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 3)
		assert.Equal(t, "a", result.Value[0])
		assert.Equal(t, "b", result.Value[1])
		assert.Equal(t, "c", result.Value[2])
	})

	t.Run("should succeed with empty input", func(t *testing.T) {
		result := Parse(EndBy(Ident(), Char(';')), "")
		require.True(t, result.OK)
		assert.Empty(t, result.Value)
	})

	t.Run("should stop when terminator missing", func(t *testing.T) {
		result := Parse(EndBy(Ident(), Char(';')), "a;b")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 1)
	})
}

//nolint:paralleltest // tests share parser state
func TestEndBy1(t *testing.T) {
	t.Run("should match one or more items", func(t *testing.T) {
		result := Parse(EndBy1(Ident(), Char(';')), "a;b;")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 2)
	})

	t.Run("should fail on empty input", func(t *testing.T) {
		assert.False(t, Parse(EndBy1(Ident(), Char(';')), "").OK)
	})

	t.Run("should fail when no terminator", func(t *testing.T) {
		assert.False(t, Parse(EndBy1(Ident(), Char(';')), "a").OK)
	})
}
