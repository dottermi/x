package ghostline

import "strings"

// prefixMatch checks if text starts with pattern (case-insensitive).
func prefixMatch(pattern, text string) bool {
	return strings.HasPrefix(
		strings.ToLower(text),
		strings.ToLower(pattern),
	)
}

// fuzzyScore returns a score for how well pattern matches text.
// Higher score = better match. Returns -1 if no match.
func fuzzyScore(pattern, text string) int {
	pattern = strings.ToLower(pattern)
	text = strings.ToLower(text)

	if len(pattern) == 0 {
		return 0
	}

	score := 0
	pi := 0
	consecutive := 0
	prevMatchIdx := -1

	for i, char := range text {
		if pi < len(pattern) && char == rune(pattern[pi]) {
			pi++
			consecutive++

			// Bonus: match at start
			if i == 0 {
				score += 10
			}

			// Bonus: consecutive matches
			if consecutive > 1 {
				score += 5 * consecutive
			}

			// Bonus: match after separator (space, -, _)
			if i > 0 {
				prev := rune(text[i-1])
				if prev == ' ' || prev == '-' || prev == '_' {
					score += 8
				}
			}

			// Penalty: gap between matches
			if prevMatchIdx >= 0 {
				gap := i - prevMatchIdx - 1
				score -= gap
			}

			prevMatchIdx = i
		} else {
			consecutive = 0
		}
	}

	if pi < len(pattern) {
		return -1 // no match
	}

	// Bonus: shorter text = more relevant
	score += 50 - len(text)

	return score
}
