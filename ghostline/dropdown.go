package ghostline

import "fmt"

// renderDropdown displays match navigation hints when multiple suggestions exist.
// Shows format: [current/total - prev  next] in dimmed text.
// Only renders when there are 2 or more matching suggestions.
func (i *Input) renderDropdown() {
	matchCount := i.countMatches()
	if matchCount > 1 {
		currentIdx := i.currentMatchIndex()
		prevMatch, nextMatch := i.getPrevNextMatches()
		_, _ = fmt.Fprintf(i.out, "\033[38;2;75;85;99m [%d/%d • ↑ %s  ↓ %s]\033[0m",
			currentIdx, matchCount, prevMatch, nextMatch)
	}
}

// countMatches returns how many suggestions match the current input.
func (i *Input) countMatches() int {
	return len(i.getMatches())
}

// currentMatchIndex returns the 1-based position of the selected match.
// Returns 0 if no matches exist.
func (i *Input) currentMatchIndex() int {
	matches := i.getMatches()
	if len(matches) == 0 {
		return 0
	}
	return (i.matchIndex % len(matches)) + 1
}

// getPrevNextMatches returns adjacent suggestion names for navigation hints.
// Returns empty strings if fewer than 2 matches exist.
func (i *Input) getPrevNextMatches() (prev, next string) {
	matches := i.getMatches()
	if len(matches) < 2 {
		return "", ""
	}

	idx := i.matchIndex % len(matches)
	prevIdx := (idx - 1 + len(matches)) % len(matches)
	nextIdx := (idx + 1) % len(matches)

	return matches[prevIdx], matches[nextIdx]
}
