package ghostline

import "strings"

// History stores command history with up/down arrow navigation.
// Preserves the user's in-progress input when navigating through entries.
type History struct {
	entries []string // chronological list of past commands
	pos     int      // navigation position (-1 = current input, 0+ = history index)
	current string   // saves in-progress input during navigation
}

// NewHistory creates an empty history.
func NewHistory() *History {
	return &History{
		entries: []string{},
		pos:     -1,
	}
}

// Add appends an entry to history after trimming whitespace.
// Ignores empty strings and consecutive duplicates.
func (h *History) Add(entry string) {
	entry = strings.TrimSpace(entry)
	if entry == "" {
		return
	}

	// Don't add duplicates of the last entry
	if len(h.entries) > 0 && h.entries[len(h.entries)-1] == entry {
		return
	}

	h.entries = append(h.entries, entry)
	h.pos = -1
}

// Len returns the number of stored history entries.
func (h *History) Len() int {
	return len(h.entries)
}

// Reset prepares history for a new readline session.
// Stores the current input for restoration after navigation.
func (h *History) Reset(current string) {
	h.pos = -1
	h.current = current
}

// Previous navigates to an older history entry.
// Saves the current input on first call for later restoration.
// Returns the entry and true, or the current input and false if at oldest.
func (h *History) Previous(current string) (string, bool) {
	if len(h.entries) == 0 {
		return current, false
	}

	// Save current input when starting navigation
	if h.pos == -1 {
		h.current = current
		h.pos = len(h.entries) - 1
		return h.entries[h.pos], true
	}

	// Move to older entry
	if h.pos > 0 {
		h.pos--
		return h.entries[h.pos], true
	}

	// Already at oldest
	return h.entries[h.pos], false
}

// Next navigates to a newer history entry or restores the original input.
// Returns the entry and true, or the saved input and true when reaching the end.
func (h *History) Next() (string, bool) {
	if h.pos == -1 {
		return h.current, false
	}

	// Move to newer entry
	if h.pos < len(h.entries)-1 {
		h.pos++
		return h.entries[h.pos], true
	}

	// Back to current input
	h.pos = -1
	return h.current, true
}
