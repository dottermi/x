package main

import (
	"fmt"

	"github.com/dottermi/x/ghostline"
)

// arrowWithStyle returns a colored prompt arrow (cyan).
func arrowWithStyle() string {
	return "\033[36m>>>\033[0m "
}

func main() {
	suggestions := []string{
		"Class(",
		"bold", "italic", "underline", "strikethrough",
		"uppercase", "lowercase", "capitalize",
		"text-left", "text-center", "text-right",
		"text-red-500", "text-green-500", "text-blue-500", "text-yellow-500",
		"text-white", "text-black",
		"bg-red-500", "bg-green-500", "bg-blue-500", "bg-yellow-500",
		"bg-white", "bg-black",
		"border", "border-rounded", "border-double",
		"p-1", "p-2", "p-4", "px-2", "py-1",
		"m-1", "m-2", "mx-2", "my-1",
		"w-20", "w-40", "w-60",
		"flex", "flex-row", "flex-col", "gap-1", "gap-2",
		":help", ":exit",
	}

	input := ghostline.NewInput(suggestions, nil, nil)

	fmt.Println("👻 Ghostline Demo")
	fmt.Println("Tab=accept • ↑↓=history • Ctrl+J=newline • Ctrl+D=exit")
	fmt.Println()

	for {
		line, err := input.Readline(arrowWithStyle())

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

		if line == ":exit" {
			fmt.Println("Bye!")
			break
		}

		input.AddHistory(line)
		fmt.Printf("You typed: %q\n", line)
	}
}
