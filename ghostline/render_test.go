package ghostline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisibleWidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"plain text", ">>>", 3},
		{"with spaces", ">>> ", 4},
		{"cyan color", "\033[36m>>>\033[0m", 3},
		{"cyan with space", "\033[36m>>>\033[0m ", 4},
		{"red bold", "\033[1;31mERROR\033[0m", 5},
		{"multiple codes", "\033[1m\033[32mOK\033[0m", 2},
		{"empty string", "", 0},
		{"only ansi", "\033[36m\033[0m", 0},
		{"wide chars", "你好", 4},
		{"wide with ansi", "\033[31m你好\033[0m", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := visibleWidth(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
