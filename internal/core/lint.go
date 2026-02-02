package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/promptforge/promptforge/internal/linter"
)

// LintProject reads plan.md and returns lint diagnostics.
func LintProject(projectDir string) ([]linter.Diagnostic, error) {
	if projectDir == "" {
		return nil, fmt.Errorf("project directory cannot be empty")
	}

	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("project directory does not exist: %s", projectDir)
	}

	planPath := filepath.Join(projectDir, "promptforge", "plan.md")
	if _, err := os.Stat(planPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("plan.md not found at %s. Run 'promptforge init' first", planPath)
	}

	planContent, err := os.ReadFile(planPath)
	if err != nil {
		if os.IsPermission(err) {
			return nil, fmt.Errorf("permission denied: cannot read plan.md at %s", planPath)
		}
		return nil, fmt.Errorf("failed to read plan.md at %s: %w", planPath, err)
	}

	return linter.LintPlan(planContent), nil
}
