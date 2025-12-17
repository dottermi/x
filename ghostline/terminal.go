package ghostline

import "golang.org/x/term"

// enableRawMode puts the terminal into raw mode for character-by-character input.
// Saves the original state for later restoration by disableRawMode.
func (i *Input) enableRawMode() error {
	state, err := term.MakeRaw(i.fd)
	if err != nil {
		return err
	}
	i.oldState = state
	return nil
}

// disableRawMode restores the terminal to its original state.
// Safe to call even if enableRawMode was never called.
func (i *Input) disableRawMode() {
	if i.oldState != nil {
		_ = term.Restore(i.fd, i.oldState)
	}
}
