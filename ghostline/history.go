package ghostline

import "strings"

// History stores command history and supports navigation.
type History struct {
	entries []string
	pos     int    // current position during navigation (-1 = new input)
	current string // saves current input when navigating
}

// NewHistory creates an empty history.
func NewHistory() *History {
	return &History{
		entries: []string{},
		pos:     -1,
	}
}

// Add appends a non-empty, non-duplicate entry to history.
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

// Len returns the number of entries.
func (h *History) Len() int {
	return len(h.entries)
}

// Reset resets navigation position for a new readline session.
func (h *History) Reset(current string) {
	h.pos = -1
	h.current = current
}

// Previous moves to the previous (older) entry.
// Returns the entry and true if available, or current input and false if at oldest.
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

// Next moves to the next (newer) entry.
// Returns the entry and true if available, or current input when back at newest.
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
