package ghostline

import (
	"slices"
	"strings"
)

// breakChars defines word boundary characters for suggestion matching.
// Inspired by rustyline's default break characters.
const breakChars = " \t\n\"'`@$><;|&{}()[],.:;"

// extractLastWord returns the final word from text for suggestion matching.
// Scans backward from the end until a break character is found.
// Returns the entire text if no break character exists.
//
// Example:
//
//	extractLastWord("git comm")  // returns "comm"
//	extractLastWord("hello")     // returns "hello"
func extractLastWord(text string) string {
	lastIdx := strings.LastIndexAny(text, breakChars)
	if lastIdx == -1 {
		return text
	}
	return text[lastIdx+1:]
}

// scoredMatch pairs a suggestion with its relevance score for sorting.
// Prefix matches are prioritized over fuzzy matches.
type scoredMatch struct {
	text     string // the suggestion text
	score    int    // fuzzy match score (higher is better)
	isPrefix bool   // true if suggestion starts with the typed word
}

// getMatches returns all suggestions matching the last word in the buffer.
// Combines prefix matching and fuzzy matching, sorted by relevance.
// Prefix matches appear first, followed by fuzzy matches sorted by score.
// Returns nil if the buffer is empty or no matches are found.
func (i *Input) getMatches() []string {
	if len(i.buffer) == 0 {
		return nil
	}

	lastWord := extractLastWord(string(i.buffer))
	if lastWord == "" {
		return nil
	}

	var matches []scoredMatch

	for _, s := range i.suggestions {
		isPrefix := strings.HasPrefix(strings.ToLower(s), strings.ToLower(lastWord))
		score := fuzzyScore(lastWord, s)

		if isPrefix || score >= 0 {
			matches = append(matches, scoredMatch{
				text:     s,
				score:    score,
				isPrefix: isPrefix,
			})
		}
	}

	if len(matches) == 0 {
		return nil
	}

	slices.SortFunc(matches, func(a, b scoredMatch) int {
		if a.isPrefix != b.isPrefix {
			if a.isPrefix {
				return -1
			}
			return 1
		}
		if b.score != a.score {
			return b.score - a.score
		}
		return strings.Compare(strings.ToLower(a.text), strings.ToLower(b.text))
	})

	result := make([]string, len(matches))
	for idx, m := range matches {
		result[idx] = m.text
	}

	return result
}

// findMatch returns the currently selected suggestion for Tab completion.
// Cycles through matches based on matchIndex, wrapping at the end.
// Returns empty string if no matches exist.
func (i *Input) findMatch() string {
	matches := i.getMatches()
	if len(matches) == 0 {
		return ""
	}
	idx := i.matchIndex % len(matches)
	return matches[idx]
}

// lastWordStart returns the buffer index where the last word begins.
// Used to determine how much text to replace when accepting a suggestion.
func (i *Input) lastWordStart() int {
	text := string(i.buffer)
	lastWord := extractLastWord(text)
	return len(i.buffer) - len([]rune(lastWord))
}

// findGhost returns the ghost text portion to display after the cursor.
// Ghost text is the untyped suffix of the current suggestion.
// Returns empty string if no suggestion matches or cursor is mid-buffer.
//
// Example: If user typed "hel" and suggestion is "hello", returns "lo".
func (i *Input) findGhost() string {
	text := string(i.buffer)
	lastWord := extractLastWord(text)
	if lastWord == "" {
		return ""
	}

	match := i.findMatch()
	if match == "" || len(match) <= len(lastWord) {
		return ""
	}
	return match[len(lastWord):]
}
