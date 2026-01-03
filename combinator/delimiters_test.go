package combinator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestBraces(t *testing.T) {
	t.Run("should match content in braces", func(t *testing.T) {
		result := Parse(Braces(Integer()), "{42}")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})
}

func TestBrackets(t *testing.T) {
	t.Run("should match content in brackets", func(t *testing.T) {
		result := Parse(Brackets(Integer()), "[42]")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})
}

func TestAngles(t *testing.T) {
	t.Run("should match content in angle brackets", func(t *testing.T) {
		result := Parse(Angles(Ident()), "<html>")
		assert.True(t, result.OK)
		assert.Equal(t, "html", result.Value)
	})
}

func TestSepBy(t *testing.T) {
	t.Run("should match items separated by delimiter", func(t *testing.T) {
		result := Parse(SepBy(Integer(), Char(',')), "1,2,3")
		require.True(t, result.OK)
		values := result.Value.([]any)
		assert.Len(t, values, 3)
		assert.Equal(t, int64(1), values[0])
		assert.Equal(t, int64(2), values[1])
		assert.Equal(t, int64(3), values[2])
	})

	t.Run("should match single item", func(t *testing.T) {
		result := Parse(SepBy(Integer(), Char(',')), "42")
		require.True(t, result.OK)
		values := result.Value.([]any)
		assert.Len(t, values, 1)
	})

	t.Run("should succeed with empty input", func(t *testing.T) {
		result := Parse(SepBy(Integer(), Char(',')), "")
		require.True(t, result.OK)
		assert.Empty(t, result.Value.([]any))
	})
}

func TestSepBy1(t *testing.T) {
	t.Run("should match items separated by delimiter", func(t *testing.T) {
		result := Parse(SepBy1(Ident(), Char(',')), "a,b,c")
		require.True(t, result.OK)
		assert.Len(t, result.Value.([]any), 3)
	})

	t.Run("should match single item", func(t *testing.T) {
		result := Parse(SepBy1(Ident(), Char(',')), "x")
		require.True(t, result.OK)
		values := result.Value.([]any)
		assert.Len(t, values, 1)
		assert.Equal(t, "x", values[0])
	})

	t.Run("should fail on empty input", func(t *testing.T) {
		assert.False(t, Parse(SepBy1(Integer(), Char(',')), "").OK)
	})
}

func TestEndBy(t *testing.T) {
	t.Run("should match items followed by terminator", func(t *testing.T) {
		result := Parse(EndBy(Ident(), Char(';')), "a;b;c;")
		require.True(t, result.OK)
		values := result.Value.([]any)
		assert.Len(t, values, 3)
		assert.Equal(t, "a", values[0])
		assert.Equal(t, "b", values[1])
		assert.Equal(t, "c", values[2])
	})

	t.Run("should succeed with empty input", func(t *testing.T) {
		result := Parse(EndBy(Ident(), Char(';')), "")
		require.True(t, result.OK)
		assert.Empty(t, result.Value.([]any))
	})

	t.Run("should stop when terminator missing", func(t *testing.T) {
		result := Parse(EndBy(Ident(), Char(';')), "a;b")
		require.True(t, result.OK)
		values := result.Value.([]any)
		assert.Len(t, values, 1)
	})
}

func TestEndBy1(t *testing.T) {
	t.Run("should match one or more items", func(t *testing.T) {
		result := Parse(EndBy1(Ident(), Char(';')), "a;b;")
		require.True(t, result.OK)
		assert.Len(t, result.Value.([]any), 2)
	})

	t.Run("should fail on empty input", func(t *testing.T) {
		assert.False(t, Parse(EndBy1(Ident(), Char(';')), "").OK)
	})

	t.Run("should fail when no terminator", func(t *testing.T) {
		assert.False(t, Parse(EndBy1(Ident(), Char(';')), "a").OK)
	})
}
