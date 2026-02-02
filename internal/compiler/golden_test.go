package compiler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/promptforge/promptforge/internal/ir"
)

func TestCompile_GoldenSimple(t *testing.T) {
	planPath := filepath.Join("testdata", "simple_plan.md")
	expectedPath := filepath.Join("testdata", "simple_prompt.ir.json")

	planContent, err := os.ReadFile(planPath)
	if err != nil {
		t.Fatalf("Failed to read plan fixture: %v", err)
	}

	compiled, err := Compile(planContent)
	if err != nil {
		t.Fatalf("Compile() failed: %v", err)
	}

	if compiled.Version == "" {
		t.Fatal("Compiled IR should include a version")
	}

	expectedBytes, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read expected fixture: %v", err)
	}

	var expectedIR ir.PromptIR
	if err := json.Unmarshal(expectedBytes, &expectedIR); err != nil {
		t.Fatalf("Failed to unmarshal expected fixture: %v", err)
	}

	if err := ValidateIR(&expectedIR); err != nil {
		t.Fatalf("Expected fixture failed validation: %v", err)
	}

	var expectedPayload interface{}
	if err := json.Unmarshal(expectedBytes, &expectedPayload); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	compiledBytes, err := json.Marshal(compiled)
	if err != nil {
		t.Fatalf("Failed to marshal compiled IR: %v", err)
	}

	var compiledPayload interface{}
	if err := json.Unmarshal(compiledBytes, &compiledPayload); err != nil {
		t.Fatalf("Failed to unmarshal compiled JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedPayload, compiledPayload) {
		t.Fatalf("Golden mismatch. Update %s if this change is expected.", expectedPath)
	}
}
