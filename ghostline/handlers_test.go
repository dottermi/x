package ghostline

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestInput(buffer string, cursorPos int) *Input {
	return &Input{
		buffer:    []rune(buffer),
		cursorPos: cursorPos,
		history:   NewHistory(),
		out:       &bytes.Buffer{},
	}
}

func TestLineBoundaries(t *testing.T) {
	t.Parallel()

	t.Run("findLineStart at beginning", func(t *testing.T) {
		t.Parallel()
		input := newTestInput("hello", 0)

		start := input.findLineStart()

		assert.Equal(t, 0, start)
	})

	t.Run("findLineStart in middle of line", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello", 3)

		start := input.findLineStart()

		assert.Equal(t, 0, start)
	})

	t.Run("findLineStart on second line", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello\nworld", 8)

		start := input.findLineStart()

		assert.Equal(t, 6, start)
	})

	t.Run("findLineEnd at end", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello", 5)

		end := input.findLineEnd()

		assert.Equal(t, 5, end)
	})

	t.Run("findLineEnd in middle of line", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello", 2)

		end := input.findLineEnd()

		assert.Equal(t, 5, end)
	})

	t.Run("findLineEnd before newline", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello\nworld", 3)

		end := input.findLineEnd()

		assert.Equal(t, 5, end)
	})
}

func TestCursorPosition(t *testing.T) {
	t.Parallel()

	t.Run("getCursorPosition at start", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello", 0)

		line, col := input.getCursorPosition()

		assert.Equal(t, 0, line)
		assert.Equal(t, 0, col)
	})

	t.Run("getCursorPosition in middle", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello", 3)

		line, col := input.getCursorPosition()

		assert.Equal(t, 0, line)
		assert.Equal(t, 3, col)
	})

	t.Run("getCursorPosition on second line", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello\nworld", 8)

		line, col := input.getCursorPosition()

		assert.Equal(t, 1, line)
		assert.Equal(t, 2, col)
	})

	t.Run("getCursorPosition at newline", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello\nworld", 6)

		line, col := input.getCursorPosition()

		assert.Equal(t, 1, line)
		assert.Equal(t, 0, col)
	})

	t.Run("getCursorPosition multiple lines", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("a\nb\nc\nd", 6)

		line, col := input.getCursorPosition()

		assert.Equal(t, 3, line)
		assert.Equal(t, 0, col)
	})

	t.Run("getCursorPosition with emoji", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("🚀hello", 4) // emoji (1 rune) + "hel" (3 runes)

		line, col := input.getCursorPosition()

		assert.Equal(t, 0, line)
		assert.Equal(t, 5, col) // emoji=2 + hel=3
	})

	t.Run("getCursorPosition with CJK", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("中文ok", 3) // 2 CJK + "o" (3 runes)

		line, col := input.getCursorPosition()

		assert.Equal(t, 0, line)
		assert.Equal(t, 5, col) // 中=2 + 文=2 + o=1
	})

	t.Run("getCursorPosition mixed unicode", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("a🚀b", 3) // a + emoji + b (3 runes)

		line, col := input.getCursorPosition()

		assert.Equal(t, 0, line)
		assert.Equal(t, 4, col) // a=1 + 🚀=2 + b=1
	})
}

func TestCtrlHandlers(t *testing.T) {
	t.Parallel()

	t.Run("handleCtrlA moves to line start", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello\nworld", 8)

		handleCtrlA(input, nil)

		assert.Equal(t, 6, input.cursorPos)
	})

	t.Run("handleCtrlE moves to line end", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello\nworld", 6)

		handleCtrlE(input, nil)

		assert.Equal(t, 11, input.cursorPos)
	})

	t.Run("handleCtrlK kills to end of line", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello world", 5)

		handleCtrlK(input, nil)

		assert.Equal(t, "hello", string(input.buffer))
		assert.Equal(t, 5, input.cursorPos)
	})

	t.Run("handleCtrlK preserves other lines", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello\nworld", 2)

		handleCtrlK(input, nil)

		assert.Equal(t, "he\nworld", string(input.buffer))
	})

	t.Run("handleCtrlU kills from line start", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello world", 5)

		handleCtrlU(input, nil)

		assert.Equal(t, " world", string(input.buffer))
		assert.Equal(t, 0, input.cursorPos)
	})

	t.Run("handleCtrlU preserves other lines", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello\nworld", 8)

		handleCtrlU(input, nil)

		assert.Equal(t, "hello\nrld", string(input.buffer))
		assert.Equal(t, 6, input.cursorPos)
	})

	t.Run("handleCtrlW deletes word", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello world", 11)

		handleCtrlW(input, nil)

		assert.Equal(t, "hello ", string(input.buffer))
		assert.Equal(t, 6, input.cursorPos)
	})

	t.Run("handleCtrlW skips trailing spaces", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello   ", 8)

		handleCtrlW(input, nil)

		assert.Empty(t, string(input.buffer))
		assert.Equal(t, 0, input.cursorPos)
	})

	t.Run("handleCtrlW at start does nothing", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello", 0)

		handleCtrlW(input, nil)

		assert.Equal(t, "hello", string(input.buffer))
		assert.Equal(t, 0, input.cursorPos)
	})
}

func TestBackspace(t *testing.T) {
	t.Parallel()

	t.Run("deletes character before cursor", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello", 3)

		handleBackspace(input, nil)

		assert.Equal(t, "helo", string(input.buffer))
		assert.Equal(t, 2, input.cursorPos)
	})

	t.Run("at start does nothing", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello", 0)

		handleBackspace(input, nil)

		assert.Equal(t, "hello", string(input.buffer))
		assert.Equal(t, 0, input.cursorPos)
	})

	t.Run("deletes newline", func(t *testing.T) {
		t.Parallel()

		input := newTestInput("hello\nworld", 6)

		handleBackspace(input, nil)

		assert.Equal(t, "helloworld", string(input.buffer))
		assert.Equal(t, 5, input.cursorPos)
	})
}

func TestNewInput(t *testing.T) {
	t.Parallel()

	t.Run("creates input with suggestions", func(t *testing.T) {
		t.Parallel()

		suggestions := []string{"a", "b", "c"}

		input := NewInput(suggestions, nil, nil)

		assert.NotNil(t, input)
		assert.Equal(t, suggestions, input.suggestions)
	})

	t.Run("creates input with history", func(t *testing.T) {
		t.Parallel()

		input := NewInput(nil, nil, nil)

		assert.NotNil(t, input.history)
	})

	t.Run("creates input with handlers", func(t *testing.T) {
		t.Parallel()

		input := NewInput(nil, nil, nil)

		assert.NotNil(t, input.handlers)
		assert.NotEmpty(t, input.handlers)
	})
}

func TestAddHistory(t *testing.T) {
	t.Parallel()

	t.Run("adds entry to history", func(t *testing.T) {
		t.Parallel()

		input := NewInput(nil, nil, nil)

		input.AddHistory("first")
		input.AddHistory("second")

		assert.Equal(t, 2, input.history.Len())
	})
}
