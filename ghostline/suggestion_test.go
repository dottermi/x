package ghostline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractLastWord(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple word", "hello", "hello"},
		{"after space", "git com", "com"},
		{"after open paren", "Class(", ""},
		{"after quote", `Class("tex`, "tex"},
		{"after comma", "a, b", "b"},
		{"after colon", "key: val", "val"},
		{"multiple delimiters", `fn(arg, "str`, "str"},
		{"empty string", "", ""},
		{"only delimiter", "(", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := extractLastWord(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindGhost(t *testing.T) {
	t.Parallel()

	t.Run("returns suffix for matching prefix", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("hel"),
			suggestions: []string{"hello", "help", "world"},
		}

		ghost := input.findGhost()

		assert.Equal(t, "lo", ghost)
	})

	t.Run("returns empty for no match", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("xyz"),
			suggestions: []string{"hello", "help", "world"},
		}

		ghost := input.findGhost()

		assert.Empty(t, ghost)
	})

	t.Run("returns empty for empty buffer", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune{},
			suggestions: []string{"hello", "help"},
		}

		ghost := input.findGhost()

		assert.Empty(t, ghost)
	})

	t.Run("returns empty when buffer ends with space", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("hello "),
			suggestions: []string{"hello", "hello world"},
		}

		ghost := input.findGhost()

		assert.Empty(t, ghost)
	})

	t.Run("matches last word only", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("git com"),
			suggestions: []string{"commit", "checkout"},
		}

		ghost := input.findGhost()

		assert.Equal(t, "mit", ghost)
	})

	t.Run("returns first match", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("he"),
			suggestions: []string{"hello", "help", "hero"},
		}

		ghost := input.findGhost()

		assert.Equal(t, "llo", ghost)
	})

	t.Run("exact match returns empty", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("hello"),
			suggestions: []string{"hello"},
		}

		ghost := input.findGhost()

		assert.Empty(t, ghost)
	})

	t.Run("case sensitive matching", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("HEL"),
			suggestions: []string{"hello", "HELLO"},
		}

		ghost := input.findGhost()

		assert.Equal(t, "LO", ghost)
	})

	t.Run("matches after open paren", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("Class(tex"),
			suggestions: []string{"text-red", "text-blue"},
		}

		ghost := input.findGhost()

		assert.Equal(t, "t-red", ghost)
	})

	t.Run("matches after quote", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune(`Class("tex`),
			suggestions: []string{"text-red"},
		}

		ghost := input.findGhost()

		assert.Equal(t, "t-red", ghost)
	})

	t.Run("matches after comma", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("a, hel"),
			suggestions: []string{"hello"},
		}

		ghost := input.findGhost()

		assert.Equal(t, "lo", ghost)
	})

	t.Run("no match after delimiter with no input", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("Class("),
			suggestions: []string{"text-red"},
		}

		ghost := input.findGhost()

		assert.Empty(t, ghost)
	})
}
