package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/promptforge/promptforge/internal/core"
	"github.com/promptforge/promptforge/internal/ir"
)

// Compile reads promptforge/plan.md and produces prompt.ir.json in the repository root.
// This is a thin CLI wrapper around core.CompileProject.
func Compile(explain bool) error {
	// Get current working directory for project root
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// CLI-specific: output to repository root
	outputPath := filepath.Join(projectDir, "prompt.ir.json")

	// Call core compilation logic
	var compiledIR *ir.PromptIR
	if explain {
		explainPath := filepath.Join(projectDir, "prompt.ir.explain.json")
		compiledIR, err = core.CompileProjectWithExplain(projectDir, outputPath, explainPath)
		if err != nil {
			return err
		}
		fmt.Printf("Wrote explain report to %s\n", explainPath)
	} else {
		compiledIR, err = core.CompileProject(projectDir, outputPath)
		if err != nil {
			return err
		}
	}

	// CLI-specific: print success message
	planPath := filepath.Join(projectDir, "promptforge", "plan.md")
	fmt.Printf("Compiled %s to %s\n", planPath, outputPath)
	_ = compiledIR // IR is returned but not used in CLI
	return nil
}
