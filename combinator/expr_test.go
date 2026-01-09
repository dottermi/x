package combinator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // tests share parser state
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

//nolint:paralleltest // tests share parser state
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

//nolint:paralleltest // tests share parser state
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

//nolint:paralleltest // tests share parser state
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

//nolint:paralleltest // tests share parser state
func TestFloatToken(t *testing.T) {
	t.Run("should match float with trailing whitespace", func(t *testing.T) {
		result := Parse(FloatToken(), "3.14   ")
		assert.True(t, result.OK)
		assert.Equal(t, float64(3.14), result.Value)
		assert.Equal(t, 7, result.State.Pos)
	})
}

//nolint:paralleltest // tests share parser state
func TestStringToken(t *testing.T) {
	t.Run("should match string literal with trailing whitespace", func(t *testing.T) {
		result := Parse(StringToken(), `"hello"   `)
		assert.True(t, result.OK)
		assert.Equal(t, "hello", result.Value)
	})
}

//nolint:paralleltest // tests share parser state
func TestChainL1(t *testing.T) {
	addOp := Map(Char('+'), func(_ rune) func(int64, int64) int64 {
		return func(a, b int64) int64 { return a + b }
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
		subOp := Map(Char('-'), func(_ rune) func(int64, int64) int64 {
			return func(a, b int64) int64 { return a - b }
		})
		result := Parse(ChainL1(Integer(), subOp), "10-3-2")
		assert.True(t, result.OK)
		assert.Equal(t, int64(5), result.Value) // ((10-3)-2) = 5
	})
}

//nolint:paralleltest // tests share parser state
func TestChainR1(t *testing.T) {
	powOp := Map(Char('^'), func(_ rune) func(int64, int64) int64 {
		return func(base, exp int64) int64 {
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

//nolint:paralleltest // tests share parser state
func TestIntegration_JSONLikeArray(t *testing.T) {
	t.Run("should parse simple integer array", func(t *testing.T) {
		parser := Brackets(SepBy(Integer(), Seq2(Char(','), Spaces())))
		result := Parse(parser, "[1, 2, 3]")
		require.True(t, result.OK)
		assert.Len(t, result.Value, 3)
	})

	t.Run("should parse empty array", func(t *testing.T) {
		emptyArray := Map(String(""), func(_ string) []int64 { return []int64{} })
		arrayContent := Choice(SepBy1(Integer(), Char(',')), emptyArray)
		parser := Brackets(arrayContent)
		result := Parse(parser, "[]")
		require.True(t, result.OK)
		assert.Empty(t, result.Value)
	})
}

//nolint:paralleltest // tests share parser state
func TestIntegration_SimpleExpression(t *testing.T) {
	t.Run("should parse arithmetic expression with precedence", func(t *testing.T) {
		addOp := Map(Lexeme(Char('+')), func(_ rune) func(int64, int64) int64 {
			return func(a, b int64) int64 { return a + b }
		})
		mulOp := Map(Lexeme(Char('*')), func(_ rune) func(int64, int64) int64 {
			return func(a, b int64) int64 { return a * b }
		})

		factor := IntToken()
		term := ChainL1(factor, mulOp)
		expr := ChainL1(term, addOp)

		result := Parse(expr, "2 + 3 * 4")
		assert.True(t, result.OK)
		assert.Equal(t, int64(14), result.Value) // 2 + (3*4) = 14
	})
}

//nolint:paralleltest // tests share parser state
func TestIntegration_FunctionCall(t *testing.T) {
	t.Run("should parse function call syntax", func(t *testing.T) {
		args := Parens(SepBy(Integer(), Seq2(Char(','), Spaces())))
		funcCall := Seq2(Ident(), args)
		result := Parse(funcCall, "sum(1, 2, 3)")
		require.True(t, result.OK)
		assert.Equal(t, "sum", result.Value.First)
		assert.Len(t, result.Value.Second, 3)
	})
}

//nolint:paralleltest // tests share parser state
func TestIntegration_MultilineTracking(t *testing.T) {
	t.Run("should track line and column", func(t *testing.T) {
		parser := Seq3(String("line1"), Newline(), String("line2"))
		result := Parse(parser, "line1\nline2")
		assert.True(t, result.OK)
		assert.Equal(t, 2, result.State.Line)
		assert.Equal(t, 6, result.State.Col)
	})
}
