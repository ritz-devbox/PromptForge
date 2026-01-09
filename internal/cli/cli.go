package cli

import (
	"fmt"
	"os"

	"github.com/promptforge/promptforge/internal/commands"
)

// Execute runs the CLI application.
func Execute() error {
	if len(os.Args) < 2 {
		printHelp()
		return fmt.Errorf("expected command: init or compile")
	}

	command := os.Args[1]

	// Handle help flags
	if command == "--help" || command == "-h" || command == "help" {
		printHelp()
		return nil
	}

	switch command {
	case "init":
		// Accept optional description as remaining arguments
		var description string
		if len(os.Args) > 2 {
			// Join all remaining arguments as the description
			description = ""
			for i := 2; i < len(os.Args); i++ {
				if i > 2 {
					description += " "
				}
				description += os.Args[i]
			}
		}
		return commands.Init(description)
	case "compile":
		return commands.Compile()
	default:
		printHelp()
		return fmt.Errorf("unknown command: %s (expected: init or compile)", command)
	}
}

// printHelp displays usage information.
func printHelp() {
	fmt.Println("PromptForge - Compile human intent into machine-enforceable prompt contracts")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  promptforge init [description]    Initialize a new project")
	fmt.Println("  promptforge compile                Compile plan.md to prompt.ir.json")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init      Create promptforge/plan.md with your rough idea")
	fmt.Println("            If description is provided, it populates the Goal section")
	fmt.Println()
	fmt.Println("  compile   Read promptforge/plan.md and generate prompt.ir.json")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  promptforge init \"I want a chatbot assistant\"")
	fmt.Println("  promptforge compile")
	fmt.Println()
	fmt.Println("For more information, see: https://github.com/promptforge/promptforge")
}


