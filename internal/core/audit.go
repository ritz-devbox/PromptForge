package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/promptforge/promptforge/internal/compiler"
	"github.com/promptforge/promptforge/internal/ir"
)

type AuditIssue struct {
	Severity string
	Message  string
}

// AuditProject validates prompt.ir.json and schema integrity.
func AuditProject(projectDir string) ([]AuditIssue, error) {
	if projectDir == "" {
		return nil, fmt.Errorf("project directory cannot be empty")
	}

	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("project directory does not exist: %s", projectDir)
	}

	var issues []AuditIssue

	irPath := filepath.Join(projectDir, "prompt.ir.json")
	data, err := os.ReadFile(irPath)
	if err != nil {
		if os.IsPermission(err) {
			return nil, fmt.Errorf("permission denied: cannot read prompt.ir.json at %s", irPath)
		}
		return nil, fmt.Errorf("failed to read prompt.ir.json at %s: %w", irPath, err)
	}

	var promptIR ir.PromptIR
	if err := json.Unmarshal(data, &promptIR); err != nil {
		return nil, fmt.Errorf("failed to parse prompt.ir.json: %w", err)
	}

	if err := compiler.ValidateIR(&promptIR); err != nil {
		issues = append(issues, AuditIssue{
			Severity: "error",
			Message:  fmt.Sprintf("IR validation failed: %s", err.Error()),
		})
	}

	if promptIR.Version != ir.CurrentVersion {
		issues = append(issues, AuditIssue{
			Severity: "error",
			Message:  fmt.Sprintf("IR version %s does not match current %s. Run 'promptforge migrate'.", promptIR.Version, ir.CurrentVersion),
		})
	}

	schemaPath := filepath.Join(projectDir, "prompt.ir.schema.json")
	schemaOnDisk, err := os.ReadFile(schemaPath)
	if err != nil {
		if os.IsNotExist(err) {
			issues = append(issues, AuditIssue{
				Severity: "warn",
				Message:  "prompt.ir.schema.json is missing",
			})
		} else if os.IsPermission(err) {
			return nil, fmt.Errorf("permission denied: cannot read prompt.ir.schema.json at %s", schemaPath)
		} else {
			return nil, fmt.Errorf("failed to read prompt.ir.schema.json at %s: %w", schemaPath, err)
		}
	} else {
		currentSchema, err := ir.PromptIRSchemaJSON()
		if err != nil {
			return nil, fmt.Errorf("failed to build current schema: %w", err)
		}
		if !bytes.Equal(schemaOnDisk, currentSchema) {
			issues = append(issues, AuditIssue{
				Severity: "warn",
				Message:  "prompt.ir.schema.json does not match the current schema. Run 'promptforge compile' to refresh.",
			})
		}
	}

	return issues, nil
}
