package commands

import (
	"fmt"
	"os"

	"github.com/promptforge/promptforge/internal/core"
)

// Migrate upgrades prompt.ir.json to the current IR version.
func Migrate() error {
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := core.MigrateIR(projectDir); err != nil {
		return err
	}

	fmt.Println("Migrated prompt.ir.json to current version")
	return nil
}
