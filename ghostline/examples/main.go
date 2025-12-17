package main

import (
	"fmt"

	"github.com/termisquad/x/ghostline"
)

func main() {
	suggestions := []string{
		"help", "hello", "history", "exit",
		"clear", "config", "commit", "checkout",
	}

	input := ghostline.NewInput(suggestions, nil, nil)

	fmt.Println("👻 Ghostline Demo")
	fmt.Println("Tab=accept • ↑↓=history • Ctrl+J=newline • Ctrl+D=exit")
	fmt.Println()

	for {
		line, err := input.Readline(">>> ")

		if err == ghostline.ErrInterrupted {
			fmt.Println("^C")
			continue
		}

		if err == ghostline.ErrEOF {
			fmt.Println("\nBye!")
			break
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}

		if line == "exit" {
			fmt.Println("Bye!")
			break
		}

		input.AddHistory(line)
		fmt.Printf("You typed: %q\n", line)
	}
}
