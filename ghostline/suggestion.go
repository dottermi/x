package ghostline

import (
	"slices"
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
	text     string
	score    int
	isPrefix bool
}

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
