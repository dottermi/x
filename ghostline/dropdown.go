package ghostline

import "fmt"

// renderDropdown displays the match counter with prev/next hints.
func (i *Input) renderDropdown() {
	matchCount := i.countMatches()
	if matchCount > 1 {
		currentIdx := i.currentMatchIndex()
		prevMatch, nextMatch := i.getPrevNextMatches()
		_, _ = fmt.Fprintf(i.out, "\033[38;2;75;85;99m [%d/%d • ↑ %s  ↓ %s]\033[0m",
			currentIdx, matchCount, prevMatch, nextMatch)
	}
}

// countMatches returns the number of suggestions that match the last word.
func (i *Input) countMatches() int {
	return len(i.getMatches())
}

// currentMatchIndex returns the 1-based index of the current match for display.
func (i *Input) currentMatchIndex() int {
	matches := i.getMatches()
	if len(matches) == 0 {
		return 0
	}
	return (i.matchIndex % len(matches)) + 1
}

// getPrevNextMatches returns the previous and next match names for display hints.
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
