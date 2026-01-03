package combinator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest // tests share parser state
func TestNewState(t *testing.T) {
	t.Run("should create state with correct initial position", func(t *testing.T) {
		state := NewState("hello")

		assert.Equal(t, 0, state.Pos)
		assert.Equal(t, 1, state.Line)
		assert.Equal(t, 1, state.Col)
	})

	t.Run("should convert input to runes", func(t *testing.T) {
		state := NewState("hello")
		assert.Equal(t, []rune{'h', 'e', 'l', 'l', 'o'}, state.Input)
	})

	t.Run("should handle unicode", func(t *testing.T) {
		state := NewState("café")
		assert.Len(t, state.Input, 4)
	})

	t.Run("should handle empty input", func(t *testing.T) {
		state := NewState("")
		assert.Empty(t, state.Input)
	})
}

//nolint:paralleltest // tests share parser state
func TestState_Current(t *testing.T) {
	t.Run("should return first rune", func(t *testing.T) {
		state := NewState("abc")
		assert.Equal(t, 'a', state.Current())
	})

	t.Run("should return zero at EOF", func(t *testing.T) {
		state := NewState("")
		assert.Equal(t, rune(0), state.Current())
	})
}

//nolint:paralleltest // tests share parser state
func TestState_IsEOF(t *testing.T) {
	t.Run("should return true for empty input", func(t *testing.T) {
		assert.True(t, NewState("").IsEOF())
	})

	t.Run("should return false at beginning", func(t *testing.T) {
		assert.False(t, NewState("hello").IsEOF())
	})

	t.Run("should return true at end", func(t *testing.T) {
		state := NewState("ab").Advance().Advance()
		assert.True(t, state.IsEOF())
	})
}

//nolint:paralleltest // tests share parser state
func TestState_Advance(t *testing.T) {
	t.Run("should increment position", func(t *testing.T) {
		next := NewState("abc").Advance()
		assert.Equal(t, 1, next.Pos)
		assert.Equal(t, 'b', next.Current())
	})

	t.Run("should increment column", func(t *testing.T) {
		next := NewState("ab").Advance()
		assert.Equal(t, 1, next.Line)
		assert.Equal(t, 2, next.Col)
	})

	t.Run("should handle newline", func(t *testing.T) {
		next := NewState("a\nb").Advance().Advance()
		assert.Equal(t, 2, next.Line)
		assert.Equal(t, 1, next.Col)
	})

	t.Run("should return same state at EOF", func(t *testing.T) {
		state := NewState("")
		next := state.Advance()
		assert.Equal(t, state.Pos, next.Pos)
	})
}

//nolint:paralleltest // tests share parser state
func TestState_AdvanceN(t *testing.T) {
	t.Run("should advance by n positions", func(t *testing.T) {
		next := NewState("abcdef").AdvanceN(3)
		assert.Equal(t, 3, next.Pos)
		assert.Equal(t, 'd', next.Current())
	})

	t.Run("should handle zero", func(t *testing.T) {
		next := NewState("abc").AdvanceN(0)
		assert.Equal(t, 0, next.Pos)
	})

	t.Run("should track across newlines", func(t *testing.T) {
		next := NewState("ab\ncd").AdvanceN(4)
		assert.Equal(t, 2, next.Line)
		assert.Equal(t, 2, next.Col)
	})
}

//nolint:paralleltest // tests share parser state
func TestSuccess(t *testing.T) {
	t.Run("should create successful result", func(t *testing.T) {
		state := NewState("test")
		result := Success("value", state)

		assert.True(t, result.OK)
		assert.Equal(t, "value", result.Value)
		assert.NoError(t, result.Err)
	})
}

//nolint:paralleltest // tests share parser state
func TestFailure(t *testing.T) {
	t.Run("should create failed result", func(t *testing.T) {
		state := NewState("test")
		result := Failure(assert.AnError, state)

		assert.False(t, result.OK)
		assert.Nil(t, result.Value)
		assert.Equal(t, assert.AnError, result.Err)
	})
}

//nolint:paralleltest // tests share parser state
func TestParse(t *testing.T) {
	t.Run("should run parser on input", func(t *testing.T) {
		result := Parse(Char('a'), "abc")
		assert.True(t, result.OK)
		assert.Equal(t, 'a', result.Value)
	})

	t.Run("should handle empty input", func(t *testing.T) {
		result := Parse(EOF(), "")
		assert.True(t, result.OK)
	})
}
