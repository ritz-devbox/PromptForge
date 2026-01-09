package commands

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCompile_SuccessfulCompilation tests successful compilation with existing plan.md.
func TestCompile_SuccessfulCompilation(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create promptforge directory
	promptforgeDir := filepath.Join(tmpDir, "promptforge")
	if err := os.MkdirAll(promptforgeDir, 0755); err != nil {
		t.Fatalf("Failed to create promptforge directory: %v", err)
	}

	// Create plan.md
	planPath := filepath.Join(promptforgeDir, "plan.md")
	planContent := []byte("# Prompt Plan\n\n## Goal\nTest goal")
	if err := os.WriteFile(planPath, planContent, 0644); err != nil {
		t.Fatalf("Failed to create plan.md: %v", err)
	}

	// Run compile
	err = Compile()
	if err != nil {
		t.Fatalf("Compile() failed: %v", err)
	}

	// Verify prompt.ir.json was created
	irPath := filepath.Join(tmpDir, "prompt.ir.json")
	if _, err := os.Stat(irPath); os.IsNotExist(err) {
		t.Error("prompt.ir.json was not created")
	}
}

// TestCompile_MissingPlanMD tests that Compile fails when plan.md does not exist.
func TestCompile_MissingPlanMD(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Run compile without creating plan.md
	err = Compile()
	if err == nil {
		t.Error("Compile() should fail when plan.md does not exist")
	}
}


