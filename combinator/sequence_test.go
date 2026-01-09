package combinator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // tests share parser state
func TestSeq2(t *testing.T) {
	t.Run("should match sequence of two", func(t *testing.T) {
		result := Parse(Seq2(Char('a'), Char('b')), "abc")
		require.True(t, result.OK)
		assert.Equal(t, 'a', result.Value.First)
		assert.Equal(t, 'b', result.Value.Second)
	})

	t.Run("should fail if first fails", func(t *testing.T) {
		assert.False(t, Parse(Seq2(Char('a'), Char('b')), "xbc").OK)
	})

	t.Run("should fail if second fails", func(t *testing.T) {
		assert.False(t, Parse(Seq2(Char('a'), Char('b')), "axc").OK)
	})

	t.Run("should advance state", func(t *testing.T) {
		result := Parse(Seq2(Char('a'), Char('b')), "abc")
		assert.Equal(t, 2, result.State.Pos)
	})
}

//nolint:paralleltest // tests share parser state
func TestSeq3(t *testing.T) {
	t.Run("should match sequence of three", func(t *testing.T) {
		result := Parse(Seq3(Char('a'), Char('b'), Char('c')), "abcd")
		require.True(t, result.OK)
		assert.Equal(t, 'a', result.Value.First)
		assert.Equal(t, 'b', result.Value.Second)
		assert.Equal(t, 'c', result.Value.Third)
	})

	t.Run("should fail if any fails", func(t *testing.T) {
		assert.False(t, Parse(Seq3(Char('a'), Char('b'), Char('c')), "axc").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestChoice(t *testing.T) {
	t.Run("should return first success", func(t *testing.T) {
		result := Parse(Choice(String("true"), String("false")), "true")
		assert.True(t, result.OK)
		assert.Equal(t, "true", result.Value)
	})

	t.Run("should try alternatives", func(t *testing.T) {
		result := Parse(Choice(String("true"), String("false")), "false")
		assert.True(t, result.OK)
		assert.Equal(t, "false", result.Value)
	})

	t.Run("should fail when none match", func(t *testing.T) {
		assert.False(t, Parse(Choice(String("true"), String("false")), "null").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestMany(t *testing.T) {
	t.Run("should match zero occurrences", func(t *testing.T) {
		result := Parse(Many(Digit()), "abc")
		require.True(t, result.OK)
		assert.Empty(t, result.Value)
	})

	t.Run("should match multiple occurrences", func(t *testing.T) {
		result := Parse(Many(Digit()), "123abc")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 3)
		assert.Equal(t, '1', result.Value[0])
		assert.Equal(t, '2', result.Value[1])
		assert.Equal(t, '3', result.Value[2])
	})

	t.Run("should stop at first non-match", func(t *testing.T) {
		result := Parse(Many(Digit()), "12abc")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 2)
		assert.Equal(t, 2, result.State.Pos)
	})
}

//nolint:paralleltest // tests share parser state
func TestMany1(t *testing.T) {
	t.Run("should match one occurrence", func(t *testing.T) {
		result := Parse(Many1(Digit()), "1abc")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 1)
	})

	t.Run("should match multiple occurrences", func(t *testing.T) {
		result := Parse(Many1(Digit()), "123abc")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 3)
	})

	t.Run("should fail on zero occurrences", func(t *testing.T) {
		assert.False(t, Parse(Many1(Digit()), "abc").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestOpt(t *testing.T) {
	t.Run("should return pointer when matched", func(t *testing.T) {
		result := Parse(Opt(Char('-')), "-42")
		assert.True(t, result.OK)
		require.NotNil(t, result.Value)
		assert.Equal(t, '-', *result.Value)
	})

	t.Run("should return nil when not matched", func(t *testing.T) {
		result := Parse(Opt(Char('-')), "42")
		assert.True(t, result.OK)
		assert.Nil(t, result.Value)
		assert.Equal(t, 0, result.State.Pos)
	})
}

//nolint:paralleltest // tests share parser state
func TestAnd(t *testing.T) {
	t.Run("should return second result", func(t *testing.T) {
		result := Parse(And(Char('a'), Char('b')), "ab")
		assert.True(t, result.OK)
		assert.Equal(t, 'b', result.Value)
	})

	t.Run("should fail if first fails", func(t *testing.T) {
		assert.False(t, Parse(And(Char('a'), Char('b')), "xb").OK)
	})

	t.Run("should fail if second fails", func(t *testing.T) {
		assert.False(t, Parse(And(Char('a'), Char('b')), "ax").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestLeft(t *testing.T) {
	t.Run("should return first result", func(t *testing.T) {
		result := Parse(Left(Integer(), Char(';')), "42;")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})

	t.Run("should consume both inputs", func(t *testing.T) {
		result := Parse(Left(Char('a'), Char('b')), "abc")
		assert.True(t, result.OK)
		assert.Equal(t, 2, result.State.Pos)
	})

	t.Run("should fail if second fails", func(t *testing.T) {
		assert.False(t, Parse(Left(Char('a'), Char('b')), "ax").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestRight(t *testing.T) {
	t.Run("should return second result", func(t *testing.T) {
		result := Parse(Right(Spaces(), Integer()), "  42")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})

	t.Run("should consume both inputs", func(t *testing.T) {
		result := Parse(Right(Char('a'), Char('b')), "abc")
		assert.True(t, result.OK)
		assert.Equal(t, 2, result.State.Pos)
	})
}

//nolint:paralleltest // tests share parser state
func TestCount(t *testing.T) {
	t.Run("should match exactly n occurrences", func(t *testing.T) {
		result := Parse(Count(3, HexDigit()), "abc123")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 3)
	})

	t.Run("should fail if fewer than n", func(t *testing.T) {
		assert.False(t, Parse(Count(5, Digit()), "123").OK)
	})

	t.Run("should handle zero count", func(t *testing.T) {
		result := Parse(Count(0, Digit()), "abc")
		require.True(t, result.OK)
		assert.Empty(t, result.Value)
	})
}
