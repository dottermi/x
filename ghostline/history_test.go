package ghostline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistory(t *testing.T) {
	t.Parallel()

	t.Run("NewHistory creates empty history", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()

		assert.NotNil(t, h)
		assert.Equal(t, 0, h.Len())
	})

	t.Run("Add stores entries", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()

		h.Add("first")
		h.Add("second")
		h.Add("third")

		assert.Equal(t, 3, h.Len())
	})

	t.Run("Add ignores empty strings", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()

		h.Add("")
		h.Add("   ")
		h.Add("\t")

		assert.Equal(t, 0, h.Len())
	})

	t.Run("Add ignores consecutive duplicates", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()

		h.Add("same")
		h.Add("same")
		h.Add("same")

		assert.Equal(t, 1, h.Len())
	})

	t.Run("Add allows non-consecutive duplicates", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()

		h.Add("first")
		h.Add("second")
		h.Add("first")

		assert.Equal(t, 3, h.Len())
	})

	t.Run("Add trims whitespace", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()

		h.Add("  hello  ")
		h.Add("hello")

		assert.Equal(t, 1, h.Len())
	})
}

func TestHistoryNavigation(t *testing.T) {
	t.Parallel()

	t.Run("Previous returns entries from newest to oldest", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()
		h.Add("first")
		h.Add("second")
		h.Add("third")
		h.Reset("")

		entry, ok := h.Previous("")
		assert.True(t, ok)
		assert.Equal(t, "third", entry)

		entry, ok = h.Previous("")
		assert.True(t, ok)
		assert.Equal(t, "second", entry)

		entry, ok = h.Previous("")
		assert.True(t, ok)
		assert.Equal(t, "first", entry)
	})

	t.Run("Previous returns false at oldest entry", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()
		h.Add("only")
		h.Reset("")

		h.Previous("")
		_, ok := h.Previous("")

		assert.False(t, ok)
	})

	t.Run("Previous saves current input", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()
		h.Add("history")
		h.Reset("")

		h.Previous("typed")
		entry, _ := h.Next()

		assert.Equal(t, "typed", entry)
	})

	t.Run("Next returns entries from oldest to newest", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()
		h.Add("first")
		h.Add("second")
		h.Reset("")

		h.Previous("")
		h.Previous("")

		entry, ok := h.Next()
		assert.True(t, ok)
		assert.Equal(t, "second", entry)
	})

	t.Run("Next returns current input at end", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()
		h.Add("history")
		h.Reset("")

		h.Previous("current")
		entry, ok := h.Next()

		assert.True(t, ok)
		assert.Equal(t, "current", entry)
	})

	t.Run("Next returns false when not navigating", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()
		h.Add("entry")
		h.Reset("current")

		entry, ok := h.Next()

		assert.False(t, ok)
		assert.Equal(t, "current", entry)
	})

	t.Run("Reset clears navigation state", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()
		h.Add("first")
		h.Add("second")
		h.Reset("")

		h.Previous("")
		h.Reset("new")

		entry, ok := h.Previous("new")
		assert.True(t, ok)
		assert.Equal(t, "second", entry)
	})

	t.Run("Previous on empty history returns current", func(t *testing.T) {
		t.Parallel()

		h := NewHistory()
		h.Reset("")

		entry, ok := h.Previous("typed")

		assert.False(t, ok)
		assert.Equal(t, "typed", entry)
	})
}
