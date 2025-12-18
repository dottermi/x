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

func TestFindMatch(t *testing.T) {
	t.Parallel()

	t.Run("returns full suggestion for matching prefix", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("hel"),
			suggestions: []string{"hello", "help"},
		}

		match := input.findMatch()

		assert.Equal(t, "hello", match)
	})

	t.Run("case insensitive returns suggestion with original case", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("class"),
			suggestions: []string{"Class", "CLASSIC"},
		}

		match := input.findMatch()

		assert.Equal(t, "Class", match)
	})

	t.Run("returns exact match", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("hello"),
			suggestions: []string{"hello"},
		}

		match := input.findMatch()

		assert.Equal(t, "hello", match)
	})

	t.Run("returns empty for no match", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("xyz"),
			suggestions: []string{"hello"},
		}

		match := input.findMatch()

		assert.Empty(t, match)
	})
}

func TestLastWordStart(t *testing.T) {
	t.Parallel()

	t.Run("single word", func(t *testing.T) {
		t.Parallel()

		input := &Input{buffer: []rune("hello")}

		start := input.lastWordStart()

		assert.Equal(t, 0, start)
	})

	t.Run("after space", func(t *testing.T) {
		t.Parallel()

		input := &Input{buffer: []rune("git com")}

		start := input.lastWordStart()

		assert.Equal(t, 4, start)
	})

	t.Run("after delimiter", func(t *testing.T) {
		t.Parallel()

		input := &Input{buffer: []rune("Class(tex")}

		start := input.lastWordStart()

		assert.Equal(t, 6, start)
	})
}

func TestGetMatches(t *testing.T) {
	t.Parallel()

	t.Run("returns all matching suggestions", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("he"),
			suggestions: []string{"hello", "help", "hero", "world"},
		}

		matches := input.getMatches()

		assert.Equal(t, []string{"hello", "help", "hero"}, matches)
	})

	t.Run("case insensitive", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("CLASS"),
			suggestions: []string{"Class", "classic", "other"},
		}

		matches := input.getMatches()

		assert.Equal(t, []string{"Class", "classic"}, matches)
	})

	t.Run("returns nil for no matches", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("xyz"),
			suggestions: []string{"hello"},
		}

		matches := input.getMatches()

		assert.Nil(t, matches)
	})
}

func TestCurrentMatchIndex(t *testing.T) {
	t.Parallel()

	t.Run("returns 1-based index", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("he"),
			suggestions: []string{"hello", "help", "hero"},
			matchIndex:  0,
		}

		assert.Equal(t, 1, input.currentMatchIndex())

		input.matchIndex = 1
		assert.Equal(t, 2, input.currentMatchIndex())

		input.matchIndex = 2
		assert.Equal(t, 3, input.currentMatchIndex())
	})

	t.Run("wraps around", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("he"),
			suggestions: []string{"hello", "help"},
			matchIndex:  5, // 5 % 2 = 1, so 1-based = 2
		}

		assert.Equal(t, 2, input.currentMatchIndex())
	})
}

func TestCountMatches(t *testing.T) {
	t.Parallel()

	t.Run("counts multiple matches", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("he"),
			suggestions: []string{"hello", "help", "hero", "world"},
		}

		count := input.countMatches()

		assert.Equal(t, 3, count)
	})

	t.Run("counts single match", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("wor"),
			suggestions: []string{"hello", "world"},
		}

		count := input.countMatches()

		assert.Equal(t, 1, count)
	})

	t.Run("case insensitive counting", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("class"),
			suggestions: []string{"Class", "CLASSIC", "classify", "other"},
		}

		count := input.countMatches()

		assert.Equal(t, 3, count)
	})

	t.Run("returns zero for no match", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("xyz"),
			suggestions: []string{"hello", "world"},
		}

		count := input.countMatches()

		assert.Equal(t, 0, count)
	})

	t.Run("returns zero for empty buffer", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune{},
			suggestions: []string{"hello"},
		}

		count := input.countMatches()

		assert.Equal(t, 0, count)
	})
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

	t.Run("case insensitive matching", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("HEL"),
			suggestions: []string{"hello", "HELLO"},
		}

		ghost := input.findGhost()

		assert.Equal(t, "lo", ghost)
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
