package combinator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLexeme(t *testing.T) {
	t.Run("should consume trailing whitespace", func(t *testing.T) {
		result := Parse(Lexeme(Integer()), "42   ")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
		assert.Equal(t, 5, result.State.Pos)
	})

	t.Run("should work without trailing whitespace", func(t *testing.T) {
		result := Parse(Lexeme(Integer()), "42")
		assert.True(t, result.OK)
	})

	t.Run("should consume tabs and newlines", func(t *testing.T) {
		result := Parse(Lexeme(Integer()), "42\t\n ")
		assert.True(t, result.OK)
		assert.Equal(t, 5, result.State.Pos)
	})
}

func TestSymbol(t *testing.T) {
	t.Run("should match string with trailing whitespace", func(t *testing.T) {
		result := Parse(Symbol("if"), "if   ")
		assert.True(t, result.OK)
		assert.Equal(t, "if", result.Value)
		assert.Equal(t, 5, result.State.Pos)
	})

	t.Run("should fail on partial match", func(t *testing.T) {
		assert.False(t, Parse(Symbol("if"), "in").OK)
	})
}

func TestToken(t *testing.T) {
	t.Run("should match identifier with trailing whitespace", func(t *testing.T) {
		result := Parse(Token(), "myVar   ")
		assert.True(t, result.OK)
		assert.Equal(t, "myVar", result.Value)
	})

	t.Run("should fail on digit-starting input", func(t *testing.T) {
		assert.False(t, Parse(Token(), "123abc").OK)
	})
}

func TestIntToken(t *testing.T) {
	t.Run("should match integer with trailing whitespace", func(t *testing.T) {
		result := Parse(IntToken(), "123   ")
		assert.True(t, result.OK)
		assert.Equal(t, int64(123), result.Value)
		assert.Equal(t, 6, result.State.Pos)
	})

	t.Run("should handle negative", func(t *testing.T) {
		result := Parse(IntToken(), "-42  ")
		assert.True(t, result.OK)
		assert.Equal(t, int64(-42), result.Value)
	})
}

func TestFloatToken(t *testing.T) {
	t.Run("should match float with trailing whitespace", func(t *testing.T) {
		result := Parse(FloatToken(), "3.14   ")
		assert.True(t, result.OK)
		assert.Equal(t, float64(3.14), result.Value)
		assert.Equal(t, 7, result.State.Pos)
	})
}

func TestStringToken(t *testing.T) {
	t.Run("should match string literal with trailing whitespace", func(t *testing.T) {
		result := Parse(StringToken(), `"hello"   `)
		assert.True(t, result.OK)
		assert.Equal(t, "hello", result.Value)
	})
}

func TestChainL1(t *testing.T) {
	addOp := Map(Char('+'), func(_ any) any {
		return func(a, b any) any { return a.(int64) + b.(int64) }
	})

	t.Run("should parse left-associative expression", func(t *testing.T) {
		result := Parse(ChainL1(Integer(), addOp), "1+2+3")
		assert.True(t, result.OK)
		assert.Equal(t, int64(6), result.Value)
	})

	t.Run("should handle single operand", func(t *testing.T) {
		result := Parse(ChainL1(Integer(), addOp), "42")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})

	t.Run("should fail on empty input", func(t *testing.T) {
		assert.False(t, Parse(ChainL1(Integer(), addOp), "").OK)
	})

	t.Run("should handle subtraction left-associatively", func(t *testing.T) {
		subOp := Map(Char('-'), func(_ any) any {
			return func(a, b any) any { return a.(int64) - b.(int64) }
		})
		result := Parse(ChainL1(Integer(), subOp), "10-3-2")
		assert.True(t, result.OK)
		assert.Equal(t, int64(5), result.Value) // ((10-3)-2) = 5
	})
}

func TestChainR1(t *testing.T) {
	powOp := Map(Char('^'), func(_ any) any {
		return func(a, b any) any {
			base, exp := a.(int64), b.(int64)
			result := int64(1)
			for i := int64(0); i < exp; i++ {
				result *= base
			}
			return result
		}
	})

	t.Run("should parse right-associative expression", func(t *testing.T) {
		result := Parse(ChainR1(Integer(), powOp), "2^3^2")
		assert.True(t, result.OK)
		assert.Equal(t, int64(512), result.Value) // 2^(3^2) = 2^9 = 512
	})

	t.Run("should handle single operand", func(t *testing.T) {
		result := Parse(ChainR1(Integer(), powOp), "42")
		assert.True(t, result.OK)
		assert.Equal(t, int64(42), result.Value)
	})

	t.Run("should fail on empty input", func(t *testing.T) {
		assert.False(t, Parse(ChainR1(Integer(), powOp), "").OK)
	})
}

func TestIntegration_JSONLikeArray(t *testing.T) {
	t.Run("should parse simple integer array", func(t *testing.T) {
		parser := Brackets(SepBy(Integer(), Seq(Char(','), Spaces())))
		result := Parse(parser, "[1, 2, 3]")
		require.True(t, result.OK)
		assert.Len(t, result.Value.([]any), 3)
	})

	t.Run("should parse empty array", func(t *testing.T) {
		emptyArray := Map(String(""), func(_ any) any { return []any{} })
		arrayContent := Choice(SepBy1(Integer(), Char(',')), emptyArray)
		parser := Brackets(arrayContent)
		result := Parse(parser, "[]")
		require.True(t, result.OK)
		assert.Empty(t, result.Value.([]any))
	})
}

func TestIntegration_SimpleExpression(t *testing.T) {
	t.Run("should parse arithmetic expression with precedence", func(t *testing.T) {
		addOp := Map(Lexeme(Char('+')), func(_ any) any {
			return func(a, b any) any { return a.(int64) + b.(int64) }
		})
		mulOp := Map(Lexeme(Char('*')), func(_ any) any {
			return func(a, b any) any { return a.(int64) * b.(int64) }
		})

		factor := IntToken()
		term := ChainL1(factor, mulOp)
		expr := ChainL1(term, addOp)

		result := Parse(expr, "2 + 3 * 4")
		assert.True(t, result.OK)
		assert.Equal(t, int64(14), result.Value) // 2 + (3*4) = 14
	})
}

func TestIntegration_FunctionCall(t *testing.T) {
	t.Run("should parse function call syntax", func(t *testing.T) {
		args := Parens(SepBy(Integer(), Seq(Char(','), Spaces())))
		funcCall := Seq(Ident(), args)
		result := Parse(funcCall, "sum(1, 2, 3)")
		require.True(t, result.OK)
		parts := result.Value.([]any)
		assert.Equal(t, "sum", parts[0])
		assert.Len(t, parts[1].([]any), 3)
	})
}

func TestIntegration_MultilineTracking(t *testing.T) {
	t.Run("should track line and column", func(t *testing.T) {
		parser := Seq(String("line1"), Newline(), String("line2"))
		result := Parse(parser, "line1\nline2")
		assert.True(t, result.OK)
		assert.Equal(t, 2, result.State.Line)
		assert.Equal(t, 6, result.State.Col)
	})
}
