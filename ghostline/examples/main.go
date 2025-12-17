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

	fmt.Println("Ghostline Demo")
	fmt.Println("Type and see suggestions. Tab to accept. Ctrl+C to exit.")
	fmt.Println()

	for {
		line, ok := input.Readline(">>> ")
		if !ok {
			fmt.Println("\nBye!")
			break
		}

		if line == "exit" {
			fmt.Println("Bye!")
			break
		}

		fmt.Printf("You typed: %q\n", line)
	}
}
