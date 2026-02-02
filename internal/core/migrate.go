package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/promptforge/promptforge/internal/compiler"
	"github.com/promptforge/promptforge/internal/ir"
)

// MigrateIR upgrades prompt.ir.json to the current IR version.
func MigrateIR(projectDir string) error {
	if projectDir == "" {
		return fmt.Errorf("project directory cannot be empty")
	}

	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", projectDir)
	}

	irPath := filepath.Join(projectDir, "prompt.ir.json")
	data, err := os.ReadFile(irPath)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied: cannot read prompt.ir.json at %s", irPath)
		}
		return fmt.Errorf("failed to read prompt.ir.json at %s: %w", irPath, err)
	}

	var promptIR ir.PromptIR
	if err := json.Unmarshal(data, &promptIR); err != nil {
		return fmt.Errorf("failed to parse prompt.ir.json: %w", err)
	}

	migrated, err := migrateToCurrent(&promptIR)
	if err != nil {
		return err
	}

	if migrated {
		if err := compiler.WriteIR(&promptIR, irPath); err != nil {
			return err
		}
	}

	schemaPath := filepath.Join(projectDir, "prompt.ir.schema.json")
	if err := compiler.WriteIRSchema(schemaPath); err != nil {
		return err
	}

	return nil
}

func migrateToCurrent(promptIR *ir.PromptIR) (bool, error) {
	if promptIR == nil {
		return false, fmt.Errorf("prompt IR is nil")
	}

	if promptIR.Version == "" || promptIR.Version == "0" {
		promptIR.Version = ir.CurrentVersion
		return true, nil
	}

	if promptIR.Version != ir.CurrentVersion {
		return false, fmt.Errorf("unsupported IR version: %s", promptIR.Version)
	}

	return false, nil
}
