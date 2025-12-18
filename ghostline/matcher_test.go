package ghostline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuzzyScore(t *testing.T) {
	t.Parallel()

	t.Run("no match returns -1", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, -1, fuzzyScore("xyz", "hello"))
	})

	t.Run("exact match scores high", func(t *testing.T) {
		t.Parallel()
		score := fuzzyScore("git", "git")
		assert.Positive(t, score)
	})

	t.Run("consecutive matches score higher", func(t *testing.T) {
		t.Parallel()
		consecutive := fuzzyScore("hel", "hello")
		sparse := fuzzyScore("hlo", "hello")
		assert.Greater(t, consecutive, sparse)
	})

	t.Run("match at start scores higher", func(t *testing.T) {
		t.Parallel()
		atStart := fuzzyScore("git", "git checkout")
		inMiddle := fuzzyScore("che", "git checkout")
		assert.Greater(t, atStart, inMiddle)
	})

	t.Run("shorter text scores higher", func(t *testing.T) {
		t.Parallel()
		short := fuzzyScore("git", "git")
		long := fuzzyScore("git", "git checkout branch")
		assert.Greater(t, short, long)
	})

	t.Run("match after separator scores bonus", func(t *testing.T) {
		t.Parallel()
		// "gc" in "git checkout" - c is after space
		withSep := fuzzyScore("gc", "git checkout")
		// "gc" in "gitconfig" - c is not after separator
		noSep := fuzzyScore("gc", "gitconfig")
		assert.Greater(t, withSep, noSep)
	})
}
