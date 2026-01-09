package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/promptforge/promptforge/internal/core"
)

// Compile reads promptforge/plan.md and produces prompt.ir.json in the repository root.
// This is a thin CLI wrapper around core.CompileProject.
func Compile() error {
	// Get current working directory for project root
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// CLI-specific: output to repository root
	outputPath := filepath.Join(projectDir, "prompt.ir.json")

	// Call core compilation logic
	ir, err := core.CompileProject(projectDir, outputPath)
	if err != nil {
		return err
	}

	// CLI-specific: print success message
	planPath := filepath.Join(projectDir, "promptforge", "plan.md")
	fmt.Printf("Compiled %s to %s\n", planPath, outputPath)
	_ = ir // IR is returned but not used in CLI
	return nil
}

