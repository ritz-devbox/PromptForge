package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInitializeProject_Success tests successful project initialization.
func TestInitializeProject_Success(t *testing.T) {
	tmpDir := t.TempDir()

	err := InitializeProject(tmpDir, "")
	if err != nil {
		t.Fatalf("InitializeProject() failed: %v", err)
	}

	// Verify promptforge directory was created
	promptforgeDir := filepath.Join(tmpDir, "promptforge")
	if _, err := os.Stat(promptforgeDir); os.IsNotExist(err) {
		t.Error("promptforge directory was not created")
	}

	// Verify plan.md was created
	planPath := filepath.Join(promptforgeDir, "plan.md")
	if _, err := os.Stat(planPath); os.IsNotExist(err) {
		t.Error("plan.md was not created")
	}
}

// TestInitializeProject_PlanAlreadyExists tests that initialization fails when plan.md already exists.
func TestInitializeProject_PlanAlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create plan.md first
	promptforgeDir := filepath.Join(tmpDir, "promptforge")
	if err := os.MkdirAll(promptforgeDir, 0755); err != nil {
		t.Fatalf("Failed to create promptforge directory: %v", err)
	}
	planPath := filepath.Join(promptforgeDir, "plan.md")
	if err := os.WriteFile(planPath, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to create existing plan.md: %v", err)
	}

	// Try to initialize - should fail
	err := InitializeProject(tmpDir, "")
	if err == nil {
		t.Error("InitializeProject() should fail when plan.md already exists")
	}
}

// TestCompileProject_Success tests successful compilation.
func TestCompileProject_Success(t *testing.T) {
	tmpDir := t.TempDir()

	// Create promptforge directory and plan.md
	promptforgeDir := filepath.Join(tmpDir, "promptforge")
	if err := os.MkdirAll(promptforgeDir, 0755); err != nil {
		t.Fatalf("Failed to create promptforge directory: %v", err)
	}
	planPath := filepath.Join(promptforgeDir, "plan.md")
	planContent := []byte("# Prompt Plan\n\n## Goal\nTest goal")
	if err := os.WriteFile(planPath, planContent, 0644); err != nil {
		t.Fatalf("Failed to create plan.md: %v", err)
	}

	// Compile
	outputPath := filepath.Join(tmpDir, "prompt.ir.json")
	ir, err := CompileProject(tmpDir, outputPath)
	if err != nil {
		t.Fatalf("CompileProject() failed: %v", err)
	}

	// Verify IR was returned
	if ir == nil {
		t.Error("CompileProject() should return non-nil IR")
	}

	// Verify output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("prompt.ir.json was not created")
	}
}

// TestInitializeProject_WithDescription tests initialization with a description.
func TestInitializeProject_WithDescription(t *testing.T) {
	tmpDir := t.TempDir()

	description := "I want a code review assistant that checks for security vulnerabilities"
	err := InitializeProject(tmpDir, description)
	if err != nil {
		t.Fatalf("InitializeProject() failed: %v", err)
	}

	// Verify plan.md was created
	planPath := filepath.Join(tmpDir, "promptforge", "plan.md")
	content, err := os.ReadFile(planPath)
	if err != nil {
		t.Fatalf("Failed to read plan.md: %v", err)
	}

	// Verify description is in the Goal section
	contentStr := string(content)
	if !strings.Contains(contentStr, description) {
		t.Errorf("plan.md should contain the description in Goal section. Got: %s", contentStr)
	}
}

// TestCompileProject_MissingPlan tests that compilation fails when plan.md does not exist.
func TestCompileProject_MissingPlan(t *testing.T) {
	tmpDir := t.TempDir()

	outputPath := filepath.Join(tmpDir, "prompt.ir.json")
	ir, err := CompileProject(tmpDir, outputPath)
	if err == nil {
		t.Error("CompileProject() should fail when plan.md does not exist")
	}
	if ir != nil {
		t.Error("CompileProject() should return nil IR on error")
	}
}


