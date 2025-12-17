package ghostline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
}
