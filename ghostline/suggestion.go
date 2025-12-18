package ghostline

import (
	"sort"
	"strings"
)

// breakChars defines word boundary characters inspired by rustyline.
// These characters separate "words" for suggestion matching.
const breakChars = " \t\n\"'`@$><;|&{}()[],.:;"

// extractLastWord returns the last "word" from text by scanning
// backward until a break character is found.
func extractLastWord(text string) string {
	lastIdx := strings.LastIndexAny(text, breakChars)
	if lastIdx == -1 {
		return text
	}
	return text[lastIdx+1:]
}

// scoredMatch holds a suggestion and its match score.
type scoredMatch struct {
	text  string
	score int
}

// getMatches returns all suggestions that match the last word.
// Prefix matches come first, then fuzzy matches, all sorted by score.
func (i *Input) getMatches() []string {
	if len(i.buffer) == 0 {
		return nil
	}

	text := string(i.buffer)
	lastWord := extractLastWord(text)
	if lastWord == "" {
		return nil
	}

	var prefixMatches, fuzzyMatches []scoredMatch

	for _, s := range i.suggestions {
		if prefixMatch(lastWord, s) {
			prefixMatches = append(prefixMatches, scoredMatch{s, fuzzyScore(lastWord, s)})
			continue
		}

		if score := fuzzyScore(lastWord, s); score >= 0 {
			fuzzyMatches = append(fuzzyMatches, scoredMatch{s, score})
		}
	}

	// Sort by score descending
	sort.Slice(prefixMatches, func(i, j int) bool {
		return prefixMatches[i].score > prefixMatches[j].score
	})
	sort.Slice(fuzzyMatches, func(i, j int) bool {
		return fuzzyMatches[i].score > fuzzyMatches[j].score
	})

	// No matches
	if len(prefixMatches) == 0 && len(fuzzyMatches) == 0 {
		return nil
	}

	// Extract strings
	result := make([]string, 0, len(prefixMatches)+len(fuzzyMatches))
	for _, m := range prefixMatches {
		result = append(result, m.text)
	}
	for _, m := range fuzzyMatches {
		result = append(result, m.text)
	}

	return result
}

// findMatch returns the current matching suggestion based on matchIndex.
// Used by handleTab to replace the typed word with the correct case.
func (i *Input) findMatch() string {
	matches := i.getMatches()
	if len(matches) == 0 {
		return ""
	}
	idx := i.matchIndex % len(matches)
	return matches[idx]
}

// lastWordStart returns the position where the last word starts in the buffer.
func (i *Input) lastWordStart() int {
	text := string(i.buffer)
	lastWord := extractLastWord(text)
	return len(i.buffer) - len([]rune(lastWord))
}

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
