package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/promptforge/promptforge/internal/core"
)

// Init initializes a new PromptForge project by creating a promptforge/plan.md file.
// This is a thin CLI wrapper around core.InitializeProjectWithTemplate.
// If description is provided, it will be used to populate the Goal section.
func Init(description string, templateName string) error {
	// Get current working directory for project root
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Call core initialization logic
	if err := core.InitializeProjectWithTemplate(projectDir, description, templateName); err != nil {
		return err
	}

	// CLI-specific: print success message
	planPath := filepath.Join(projectDir, "promptforge", "plan.md")
	fmt.Printf("Initialized PromptForge project: created %s\n", planPath)
	return nil
}
