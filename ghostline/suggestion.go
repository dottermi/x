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

// findMatch returns the full matching suggestion for the last word.
// Used by handleTab to replace the typed word with the correct case.
func (i *Input) findMatch() string {
	if len(i.buffer) == 0 {
		return ""
	}

	text := string(i.buffer)
	lastWord := extractLastWord(text)
	if lastWord == "" {
		return ""
	}

	lastWordLower := strings.ToLower(lastWord)
	for _, s := range i.suggestions {
		sLower := strings.ToLower(s)
		if strings.HasPrefix(sLower, lastWordLower) && len(s) >= len(lastWord) {
			return s
		}
	}
	return ""
}

// lastWordStart returns the position where the last word starts in the buffer.
func (i *Input) lastWordStart() int {
	text := string(i.buffer)
	lastWord := extractLastWord(text)
	return len(i.buffer) - len([]rune(lastWord))
}

// countMatches returns the number of suggestions that match the last word.
func (i *Input) countMatches() int {
	if len(i.buffer) == 0 {
		return 0
	}

	text := string(i.buffer)
	lastWord := extractLastWord(text)
	if lastWord == "" {
		return 0
	}

	count := 0
	lastWordLower := strings.ToLower(lastWord)
	for _, s := range i.suggestions {
		sLower := strings.ToLower(s)
		if strings.HasPrefix(sLower, lastWordLower) {
			count++
		}
	}
	return count
}

func (i *Input) findGhost() string {
	if len(i.buffer) == 0 {
		return ""
	}

	text := string(i.buffer)
	lastWord := extractLastWord(text)
	if lastWord == "" {
		return ""
	}

	lastWordLower := strings.ToLower(lastWord)
	for _, s := range i.suggestions {
		sLower := strings.ToLower(s)
		if strings.HasPrefix(sLower, lastWordLower) && len(s) > len(lastWord) {
			return s[len(lastWord):]
		}
	}
	return ""
}
