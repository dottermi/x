package ghostline

import "strings"

// fuzzyScore calculates how well pattern matches text using fuzzy matching.
// Returns a positive score for matches (higher is better) or -1 for no match.
//
// Scoring bonuses:
//   - Match at start of text: +10
//   - Consecutive character matches: +5 per consecutive char
//   - Match after word separator (space, hyphen, underscore): +8
//   - Shorter text: up to +50
//
// Scoring penalties:
//   - Gap between matched characters: -1 per gap character
//
// Example:
//
//	fuzzyScore("gc", "git-checkout")  // positive score (matches g...c)
//	fuzzyScore("xyz", "hello")        // returns -1 (no match)
func fuzzyScore(pattern, text string) int {
	if len(pattern) == 0 {
		return 0
	}

	pRunes := []rune(strings.ToLower(pattern))
	tRunes := []rune(strings.ToLower(text))

	score, pIdx, prevMatchIdx, consecutive := 0, 0, -1, 0

	for tIdx, char := range tRunes {
		if pIdx >= len(pRunes) || char != pRunes[pIdx] {
			consecutive = 0
			continue
		}

		score += calculateMatchBonus(tIdx, prevMatchIdx, consecutive, tRunes)

		consecutive++
		prevMatchIdx = tIdx
		pIdx++
	}

	if pIdx < len(pRunes) {
		return -1
	}

	return score + (50 - len(tRunes))
}

// calculateMatchBonus computes the bonus score for a matched character.
// Bonuses are awarded for matches at the start, consecutive matches,
// and matches after word separators. Penalties are applied for gaps.
func calculateMatchBonus(tIdx, prevMatchIdx, consecutive int, tRunes []rune) int {
	bonus := 0

	if tIdx == 0 {
		bonus += 10
	}

	if consecutive > 0 {
		bonus += 5 * (consecutive + 1)
	}

	if tIdx > 0 {
		prev := tRunes[tIdx-1]
		if prev == ' ' || prev == '-' || prev == '_' {
			bonus += 8
		}
	}

	if prevMatchIdx >= 0 {
		gap := tIdx - prevMatchIdx - 1
		bonus -= gap
	}

	return bonus
}
