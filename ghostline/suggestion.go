package ghostline

import "strings"

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

// getMatches returns all suggestions that match the last word.
func (i *Input) getMatches() []string {
	if len(i.buffer) == 0 {
		return nil
	}

	text := string(i.buffer)
	lastWord := extractLastWord(text)
	if lastWord == "" {
		return nil
	}

	var matches []string
	lastWordLower := strings.ToLower(lastWord)
	for _, s := range i.suggestions {
		sLower := strings.ToLower(s)
		if strings.HasPrefix(sLower, lastWordLower) {
			matches = append(matches, s)
		}
	}
	return matches
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
