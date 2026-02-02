package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/promptforge/promptforge/internal/core"
	"github.com/promptforge/promptforge/internal/linter"
)

// Lint reads promptforge/plan.md and prints diagnostics.
func Lint() error {
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	diagnostics, err := core.LintProject(projectDir)
	if err != nil {
		return err
	}

	planPath := filepath.Join(projectDir, "promptforge", "plan.md")
	var errorCount int
	for _, diag := range diagnostics {
		fmt.Printf("%s:%d:%d: %s %s %s\n", planPath, diag.Line, diag.Column, diag.Severity, diag.Code, diag.Message)
		if diag.Severity == linter.SeverityError {
			errorCount++
		}
	}

	if errorCount > 0 {
		return fmt.Errorf("lint failed with %d error(s)", errorCount)
	}

	return nil
}
