package ghostline

import "strings"

func (i *Input) findGhost() string {
	if len(i.buffer) == 0 {
		return ""
	}

	text := string(i.buffer)

	parts := strings.Fields(text)
	if len(parts) == 0 {
		return ""
	}
	lastWord := parts[len(parts)-1]

	if strings.HasSuffix(text, " ") {
		return ""
	}

	for _, s := range i.suggestions {
		if strings.HasPrefix(s, lastWord) && len(s) > len(lastWord) {
			return s[len(lastWord):]
		}
	}
	return ""
}
