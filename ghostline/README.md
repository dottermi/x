# `./ghostline 👻`

Minimal readline API with ghost text suggestions

## `features`

- Ghost text suggestions (dimmed inline completions)
- Tab to accept suggestions
- Command history with arrow keys
- Multiline editing (Ctrl+J)
- Emacs-style keybindings

## `install`

```bash
go get github.com/termisquad/x/ghostline@latest
```

## `usage`

```go
package main

import (
    "fmt"
    "github.com/termisquad/x/ghostline"
)

func main() {
    suggestions := []string{"help", "history", "exit"}
    input := ghostline.NewInput(suggestions, nil, nil)

    for {
        line, err := input.Readline(">>> ")

        if err == ghostline.ErrInterrupted {
            continue
        }

        if err == ghostline.ErrEOF {
            break
        }

        if err != nil {
            break
        }

        input.AddHistory(line)
        fmt.Println("You typed:", line)
    }
}
```

## `keybindings`

| Key      | Action              |
| -------- | ------------------- |
| Tab      | Accept suggestion   |
| Enter    | Submit              |
| Ctrl+J   | New line            |
| Ctrl+C   | Interrupt           |
| Ctrl+D   | EOF (exit)          |
| ← →      | Move cursor         |
| ↑ ↓      | History / Move line |
| Ctrl+A   | Beginning of line   |
| Ctrl+E   | End of line         |
| Ctrl+K   | Kill to end of line |
| Ctrl+U   | Kill to beginning   |
| Ctrl+W   | Delete word         |
| Home/End | Line boundaries     |
| Delete   | Delete char         |

## `example`

```
go run ./examples/main.go
```

## `license`

MIT
