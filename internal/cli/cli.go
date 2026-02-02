package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/promptforge/promptforge/internal/commands"
)

// Execute runs the CLI application.
func Execute() error {
	if len(os.Args) < 2 {
		printHelp()
		return fmt.Errorf("expected command: init, compile, lint, templates, migrate, or audit")
	}

	command := os.Args[1]

	// Handle help flags
	if command == "--help" || command == "-h" || command == "help" {
		printHelp()
		return nil
	}

	switch command {
	case "init":
		var description string
		var templateName string
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			switch arg {
			case "--template":
				if i+1 >= len(os.Args) {
					return fmt.Errorf("missing value for --template")
				}
				templateName = os.Args[i+1]
				i++
			default:
				if strings.HasPrefix(arg, "-") {
					return fmt.Errorf("unknown flag for init: %s", arg)
				}
				if description != "" {
					description += " "
				}
				description += arg
			}
		}
		return commands.Init(description, templateName)
	case "compile":
		explain := false
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			switch arg {
			case "--explain":
				explain = true
			default:
				return fmt.Errorf("unknown flag for compile: %s", arg)
			}
		}
		return commands.Compile(explain)
	case "lint":
		return commands.Lint()
	case "templates":
		return commands.ListTemplates()
	case "migrate":
		return commands.Migrate()
	case "audit":
		return commands.Audit()
	default:
		printHelp()
		return fmt.Errorf("unknown command: %s (expected: init, compile, lint, templates, migrate, or audit)", command)
	}
}

// printHelp displays usage information.
func printHelp() {
	fmt.Println("PromptForge - Compile human intent into machine-enforceable prompt contracts")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  promptforge init [description]    Initialize a new project")
	fmt.Println("  promptforge compile                Compile plan.md to prompt.ir.json")
	fmt.Println("  promptforge lint                   Lint plan.md and report issues")
	fmt.Println("  promptforge templates              List available templates")
	fmt.Println("  promptforge migrate                Upgrade prompt.ir.json to current version")
	fmt.Println("  promptforge audit                  Validate prompt.ir.json integrity")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init      Create promptforge/plan.md with your rough idea")
	fmt.Println("            If description is provided, it populates the Goal section")
	fmt.Println("            Use --template <name> to start from a preset")
	fmt.Println()
	fmt.Println("  compile   Read promptforge/plan.md and generate prompt.ir.json")
	fmt.Println("            Use --explain to write prompt.ir.explain.json")
	fmt.Println("  lint      Analyze promptforge/plan.md and print diagnostics")
	fmt.Println("  templates List built-in plan templates")
	fmt.Println("  migrate   Upgrade prompt.ir.json to the latest IR format")
	fmt.Println("  audit     Validate prompt.ir.json and schema compatibility")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  promptforge init \"I want a chatbot assistant\"")
	fmt.Println("  promptforge compile")
	fmt.Println("  promptforge lint")
	fmt.Println("  promptforge templates")
	fmt.Println("  promptforge migrate")
	fmt.Println("  promptforge audit")
	fmt.Println()
	fmt.Println("For more information, see: https://github.com/promptforge/promptforge")
}
