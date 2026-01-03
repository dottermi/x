package combinator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Run("should transform successful result", func(t *testing.T) {
		parser := Map(Integer(), func(v any) any { return v.(int64) * 2 })
		result := Parse(parser, "21")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})

	t.Run("should propagate failure", func(t *testing.T) {
		parser := Map(Integer(), func(v any) any { return v.(int64) * 2 })
		assert.False(t, Parse(parser, "abc").OK)
	})

	t.Run("should allow type conversion", func(t *testing.T) {
		parser := Map(Many(Digit()), func(v any) any {
			runes := v.([]any)
			s := ""
			for _, r := range runes {
				s += string(r.(rune))
			}
			return s
		})
		result := Parse(parser, "123")
		assert.True(t, result.OK)
		assert.Equal(t, "123", result.Value)
	})
}

func TestMapErr(t *testing.T) {
	t.Run("should transform error on failure", func(t *testing.T) {
		parser := MapErr(Digit(), func(err error) error { return assert.AnError })
		result := Parse(parser, "abc")
		assert.False(t, result.OK)
		assert.Equal(t, assert.AnError, result.Err)
	})

	t.Run("should not affect success", func(t *testing.T) {
		parser := MapErr(Digit(), func(err error) error { return assert.AnError })
		result := Parse(parser, "123")
		assert.True(t, result.OK)
	})
}

func TestLabel(t *testing.T) {
	t.Run("should add label to error", func(t *testing.T) {
		result := Parse(Label(Range('0', '9'), "digit"), "abc")
		assert.False(t, result.OK)
		assert.Contains(t, result.Err.Error(), "expected digit")
	})

	t.Run("should not affect success", func(t *testing.T) {
		result := Parse(Label(Range('0', '9'), "digit"), "5")
		assert.True(t, result.OK)
	})
}

func TestSkip(t *testing.T) {
	t.Run("should return nil on success", func(t *testing.T) {
		result := Parse(Skip(String("hello")), "hello world")
		assert.True(t, result.OK)
		assert.Nil(t, result.Value)
	})

	t.Run("should consume input", func(t *testing.T) {
		result := Parse(Skip(String("hello")), "hello world")
		assert.Equal(t, 5, result.State.Pos)
	})

	t.Run("should propagate failure", func(t *testing.T) {
		assert.False(t, Parse(Skip(String("hello")), "world").OK)
	})
}

func TestSkipMany(t *testing.T) {
	t.Run("should skip zero occurrences", func(t *testing.T) {
		result := Parse(SkipMany(Digit()), "abc")
		assert.True(t, result.OK)
		assert.Nil(t, result.Value)
		assert.Equal(t, 0, result.State.Pos)
	})

	t.Run("should skip multiple occurrences", func(t *testing.T) {
		result := Parse(SkipMany(Digit()), "123abc")
		assert.True(t, result.OK)
		assert.Equal(t, 3, result.State.Pos)
	})
}

func TestSkipMany1(t *testing.T) {
	t.Run("should skip one or more", func(t *testing.T) {
		result := Parse(SkipMany1(Digit()), "123abc")
		assert.True(t, result.OK)
		assert.Equal(t, 3, result.State.Pos)
	})

	t.Run("should fail on zero occurrences", func(t *testing.T) {
		assert.False(t, Parse(SkipMany1(Digit()), "abc").OK)
	})
}

func TestNot(t *testing.T) {
	t.Run("should succeed when parser fails", func(t *testing.T) {
		result := Parse(Not(String("if")), "else")
		assert.True(t, result.OK)
		assert.Nil(t, result.Value)
		assert.Equal(t, 0, result.State.Pos)
	})

	t.Run("should fail when parser succeeds", func(t *testing.T) {
		result := Parse(Not(String("if")), "if")
		assert.False(t, result.OK)
		assert.Contains(t, result.Err.Error(), "unexpected match")
	})
}

func TestLookAhead(t *testing.T) {
	t.Run("should return value without consuming", func(t *testing.T) {
		result := Parse(LookAhead(String("hello")), "hello world")
		assert.True(t, result.OK)
		assert.Equal(t, "hello", result.Value)
		assert.Equal(t, 0, result.State.Pos)
	})

	t.Run("should propagate failure", func(t *testing.T) {
		assert.False(t, Parse(LookAhead(String("hello")), "world").OK)
	})
}

func TestLazy(t *testing.T) {
	t.Run("should defer parser evaluation", func(t *testing.T) {
		var p Parser
		lazy := Lazy(&p)
		p = Char('x')
		result := Parse(lazy, "x")
		assert.True(t, result.OK)
		assert.Equal(t, 'x', result.Value)
	})
}

func TestRef(t *testing.T) {
	t.Run("should enable recursive grammars", func(t *testing.T) {
		var expr Rule
		expr = func() Parser {
			return Choice(Integer(), Parens(Ref(&expr)))
		}
		result := Parse(Ref(&expr), "((42))")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})

	t.Run("should handle deeply nested", func(t *testing.T) {
		var expr Rule
		expr = func() Parser {
			return Choice(Integer(), Parens(Ref(&expr)))
		}
		result := Parse(Ref(&expr), "(((((1)))))")
		assert.True(t, result.OK)
		assert.Equal(t, int64(1), result.Value)
	})
}
