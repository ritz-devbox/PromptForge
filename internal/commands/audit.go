package commands

import (
	"fmt"
	"os"

	"github.com/promptforge/promptforge/internal/core"
)

// Audit validates prompt.ir.json and schema integrity.
func Audit() error {
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	issues, err := core.AuditProject(projectDir)
	if err != nil {
		return err
	}

	var errorCount int
	for _, issue := range issues {
		fmt.Printf("%s: %s\n", issue.Severity, issue.Message)
		if issue.Severity == "error" {
			errorCount++
		}
	}

	if errorCount > 0 {
		return fmt.Errorf("audit failed with %d error(s)", errorCount)
	}

	if len(issues) == 0 {
		fmt.Println("Audit passed with no issues")
	}

	return nil
}
