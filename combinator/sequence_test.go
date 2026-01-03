package combinator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // tests share parser state
func TestSeq(t *testing.T) {
	t.Run("should match sequence", func(t *testing.T) {
		result := Parse(Seq(Char('a'), Char('b'), Char('c')), "abc")
		require.True(t, result.OK)
		values := result.Value.([]any) //nolint:errcheck,forcetypeassert // test assertion
		assert.Len(t, values, 3)
		assert.Equal(t, 'a', values[0])
		assert.Equal(t, 'b', values[1])
		assert.Equal(t, 'c', values[2])
	})

	t.Run("should fail if any fails", func(t *testing.T) {
		assert.False(t, Parse(Seq(Char('a'), Char('b'), Char('c')), "axc").OK)
	})

	t.Run("should handle empty sequence", func(t *testing.T) {
		result := Parse(Seq(), "abc")
		require.True(t, result.OK)
		assert.Empty(t, result.Value.([]any)) //nolint:errcheck,forcetypeassert // test assertion
	})

	t.Run("should advance state", func(t *testing.T) {
		result := Parse(Seq(Char('a'), Char('b')), "abc")
		assert.Equal(t, 2, result.State.Pos)
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

	t.Run("should handle empty choice", func(t *testing.T) {
		result := Parse(Choice(), "abc")
		assert.False(t, result.OK)
		assert.Contains(t, result.Err.Error(), "no alternatives")
	})
}

//nolint:paralleltest // tests share parser state
func TestMany(t *testing.T) {
	t.Run("should match zero occurrences", func(t *testing.T) {
		result := Parse(Many(Digit()), "abc")
		require.True(t, result.OK)
		assert.Empty(t, result.Value.([]any)) //nolint:errcheck,forcetypeassert // test assertion
	})

	t.Run("should match multiple occurrences", func(t *testing.T) {
		result := Parse(Many(Digit()), "123abc")
		require.True(t, result.OK)
		values := result.Value.([]any) //nolint:errcheck,forcetypeassert // test assertion
		assert.Len(t, values, 3)
	})

	t.Run("should stop at first non-match", func(t *testing.T) {
		result := Parse(Many(Digit()), "12abc")
		require.True(t, result.OK)
		assert.Len(t, result.Value.([]any), 2) //nolint:errcheck,forcetypeassert // test assertion
		assert.Equal(t, 2, result.State.Pos)
	})
}

//nolint:paralleltest // tests share parser state
func TestMany1(t *testing.T) {
	t.Run("should match one occurrence", func(t *testing.T) {
		result := Parse(Many1(Digit()), "1abc")
		require.True(t, result.OK)
		assert.Len(t, result.Value.([]any), 1) //nolint:errcheck,forcetypeassert // test assertion
	})

	t.Run("should match multiple occurrences", func(t *testing.T) {
		result := Parse(Many1(Digit()), "123abc")
		require.True(t, result.OK)
		assert.Len(t, result.Value.([]any), 3) //nolint:errcheck,forcetypeassert // test assertion
	})

	t.Run("should fail on zero occurrences", func(t *testing.T) {
		assert.False(t, Parse(Many1(Digit()), "abc").OK)
	})
}

//nolint:paralleltest // tests share parser state
func TestOpt(t *testing.T) {
	t.Run("should return value when matched", func(t *testing.T) {
		result := Parse(Opt(Char('-')), "-42")
		assert.True(t, result.OK)
		assert.Equal(t, '-', result.Value)
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
		values := result.Value.([]any) //nolint:errcheck,forcetypeassert // test assertion
		assert.Len(t, values, 3)
	})

	t.Run("should fail if fewer than n", func(t *testing.T) {
		assert.False(t, Parse(Count(5, Digit()), "123").OK)
	})

	t.Run("should handle zero count", func(t *testing.T) {
		result := Parse(Count(0, Digit()), "abc")
		require.True(t, result.OK)
		assert.Empty(t, result.Value.([]any)) //nolint:errcheck,forcetypeassert // test assertion
	})
}
