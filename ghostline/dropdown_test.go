package ghostline

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	t.Run("returns zero for no matches", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("xyz"),
			suggestions: []string{"hello"},
		}

		assert.Equal(t, 0, input.currentMatchIndex())
	})
}

func TestGetPrevNextMatches(t *testing.T) {
	t.Parallel()

	t.Run("returns prev and next matches", func(t *testing.T) {
		t.Parallel()

		// Order by score: help, hero, hello
		input := &Input{
			buffer:      []rune("he"),
			suggestions: []string{"hello", "help", "hero"},
			matchIndex:  1, // current is "hero"
		}

		prev, next := input.getPrevNextMatches()

		assert.Equal(t, "help", prev)
		assert.Equal(t, "hello", next)
	})

	t.Run("wraps around at start", func(t *testing.T) {
		t.Parallel()

		// Order by score: help, hero, hello
		input := &Input{
			buffer:      []rune("he"),
			suggestions: []string{"hello", "help", "hero"},
			matchIndex:  0, // current is "help"
		}

		prev, next := input.getPrevNextMatches()

		assert.Equal(t, "hello", prev) // wraps to last
		assert.Equal(t, "hero", next)
	})

	t.Run("wraps around at end", func(t *testing.T) {
		t.Parallel()

		// Order by score: help, hero, hello
		input := &Input{
			buffer:      []rune("he"),
			suggestions: []string{"hello", "help", "hero"},
			matchIndex:  2, // current is "hello"
		}

		prev, next := input.getPrevNextMatches()

		assert.Equal(t, "hero", prev)
		assert.Equal(t, "help", next) // wraps to first
	})

	t.Run("returns empty for single match", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("wor"),
			suggestions: []string{"hello", "world"},
			matchIndex:  0,
		}

		prev, next := input.getPrevNextMatches()

		assert.Empty(t, prev)
		assert.Empty(t, next)
	})

	t.Run("returns empty for no matches", func(t *testing.T) {
		t.Parallel()

		input := &Input{
			buffer:      []rune("xyz"),
			suggestions: []string{"hello"},
		}

		prev, next := input.getPrevNextMatches()

		assert.Empty(t, prev)
		assert.Empty(t, next)
	})
}

func TestRenderDropdown(t *testing.T) {
	t.Parallel()

	t.Run("renders dropdown with multiple matches", func(t *testing.T) {
		t.Parallel()

		// Order by score: help, hero, hello
		var buf bytes.Buffer
		input := &Input{
			out:         &buf,
			buffer:      []rune("he"),
			suggestions: []string{"hello", "help", "hero"},
			matchIndex:  0, // current is "help"
		}

		input.renderDropdown()

		output := buf.String()
		assert.Contains(t, output, "[1/3")
		assert.Contains(t, output, "hello") // prev (wraps to last)
		assert.Contains(t, output, "hero")  // next
	})

	t.Run("does not render for single match", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		input := &Input{
			out:         &buf,
			buffer:      []rune("wor"),
			suggestions: []string{"hello", "world"},
			matchIndex:  0,
		}

		input.renderDropdown()

		assert.Empty(t, buf.String())
	})

	t.Run("does not render for no matches", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		input := &Input{
			out:         &buf,
			buffer:      []rune("xyz"),
			suggestions: []string{"hello"},
		}

		input.renderDropdown()

		assert.Empty(t, buf.String())
	})

	t.Run("updates with matchIndex changes", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		input := &Input{
			out:         &buf,
			buffer:      []rune("he"),
			suggestions: []string{"hello", "help"},
			matchIndex:  1,
		}

		input.renderDropdown()

		output := buf.String()
		assert.Contains(t, output, "[2/2")
	})
}
