package core

import (
	"path/filepath"
	"testing"

	"github.com/promptforge/promptforge/internal/compiler"
	"github.com/promptforge/promptforge/internal/ir"
)

func TestAuditProject_Pass(t *testing.T) {
	tmpDir := t.TempDir()

	promptIR := &ir.PromptIR{
		Version:    ir.CurrentVersion,
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{ID: "rule-1", Description: "Test rule"},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "Test", Response: "Test"},
		},
	}

	irPath := filepath.Join(tmpDir, "prompt.ir.json")
	if err := compiler.WriteIR(promptIR, irPath); err != nil {
		t.Fatalf("WriteIR() failed: %v", err)
	}

	schemaPath := filepath.Join(tmpDir, "prompt.ir.schema.json")
	if err := compiler.WriteIRSchema(schemaPath); err != nil {
		t.Fatalf("WriteIRSchema() failed: %v", err)
	}

	issues, err := AuditProject(tmpDir)
	if err != nil {
		t.Fatalf("AuditProject() failed: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("Expected no audit issues, got %d", len(issues))
	}
}

func TestAuditProject_MissingSchemaWarns(t *testing.T) {
	tmpDir := t.TempDir()

	promptIR := &ir.PromptIR{
		Version:    ir.CurrentVersion,
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{ID: "rule-1", Description: "Test rule"},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "Test", Response: "Test"},
		},
	}

	irPath := filepath.Join(tmpDir, "prompt.ir.json")
	if err := compiler.WriteIR(promptIR, irPath); err != nil {
		t.Fatalf("WriteIR() failed: %v", err)
	}

	issues, err := AuditProject(tmpDir)
	if err != nil {
		t.Fatalf("AuditProject() failed: %v", err)
	}

	if len(issues) == 0 {
		t.Fatal("Expected warnings when schema is missing")
	}
}

func TestAuditProject_VersionMismatch(t *testing.T) {
	tmpDir := t.TempDir()

	promptIR := &ir.PromptIR{
		Version:    "0.9",
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{ID: "rule-1", Description: "Test rule"},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "Test", Response: "Test"},
		},
	}

	irPath := filepath.Join(tmpDir, "prompt.ir.json")
	if err := compiler.WriteIR(promptIR, irPath); err != nil {
		t.Fatalf("WriteIR() failed: %v", err)
	}

	issues, err := AuditProject(tmpDir)
	if err != nil {
		t.Fatalf("AuditProject() failed: %v", err)
	}

	if !containsAuditSeverity(issues, "error") {
		t.Fatal("Expected version mismatch to be an error")
	}
}

func containsAuditSeverity(issues []AuditIssue, severity string) bool {
	for _, issue := range issues {
		if issue.Severity == severity {
			return true
		}
	}
	return false
}
