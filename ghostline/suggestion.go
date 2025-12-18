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

func (i *Input) findGhost() string {
	if len(i.buffer) == 0 {
		return ""
	}

	text := string(i.buffer)
	lastWord := extractLastWord(text)
	if lastWord == "" {
		return ""
	}

	for _, s := range i.suggestions {
		if strings.HasPrefix(s, lastWord) && len(s) > len(lastWord) {
			return s[len(lastWord):]
		}
	}
	return ""
}
