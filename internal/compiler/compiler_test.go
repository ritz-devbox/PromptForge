package compiler

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/promptforge/promptforge/internal/ir"
)

// TestCompile_SuccessfulIRGeneration tests that Compile generates a valid IR.
func TestCompile_SuccessfulIRGeneration(t *testing.T) {
	planContent := []byte("# Prompt Plan\n\n## Goal\nTest goal")
	
	ir, err := Compile(planContent)
	if err != nil {
		t.Fatalf("Compile() failed: %v", err)
	}

	// Validate the generated IR
	if err := ValidateIR(ir); err != nil {
		t.Fatalf("Generated IR failed validation: %v", err)
	}

	// Verify required fields are present
	if ir.SystemRole == "" {
		t.Error("SystemRole should not be empty")
	}
	if len(ir.Rules) == 0 {
		t.Error("Rules should not be empty")
	}
	if ir.InputSchema.Type == "" {
		t.Error("InputSchema.Type should not be empty")
	}
	if ir.OutputSchema.Type == "" {
		t.Error("OutputSchema.Type should not be empty")
	}
	if len(ir.FailureModes) == 0 {
		t.Error("FailureModes should not be empty")
	}
}

// TestCompile_EmptyPlanContent tests that Compile fails on empty plan content.
func TestCompile_EmptyPlanContent(t *testing.T) {
	planContent := []byte("")
	
	ir, err := Compile(planContent)
	if err == nil {
		t.Error("Compile() should fail on empty plan content")
	}
	if ir != nil {
		t.Error("Compile() should return nil IR on error")
	}
}

// TestCompile_IDCollisionDetection tests that duplicate rule IDs are handled.
func TestCompile_IDCollisionDetection(t *testing.T) {
	planContent := []byte(`# Prompt Plan

## Goal
Test goal

## Constraints
- Must be secure
- Must be secure
- Must be secure
`)

	ir, err := Compile(planContent)
	if err != nil {
		t.Fatalf("Compile() failed: %v", err)
	}

	// Check that all rule IDs are unique
	ruleIDs := make(map[string]bool)
	for _, rule := range ir.Rules {
		if ruleIDs[rule.ID] {
			t.Errorf("Duplicate rule ID found: %s", rule.ID)
		}
		ruleIDs[rule.ID] = true
	}
}

// TestCompile_VeryLongGoal tests that very long goals are handled.
func TestCompile_VeryLongGoal(t *testing.T) {
	longGoal := strings.Repeat("This is a very long goal. ", 100)
	planContent := []byte("# Prompt Plan\n\n## Goal\n" + longGoal)
	
	ir, err := Compile(planContent)
	if err != nil {
		t.Fatalf("Compile() failed with long goal: %v", err)
	}
	
	if len(ir.SystemRole) == 0 {
		t.Error("SystemRole should not be empty")
	}
}

// TestCompile_SpecialCharactersInConstraints tests handling of special characters.
func TestCompile_SpecialCharactersInConstraints(t *testing.T) {
	planContent := []byte(`# Prompt Plan

## Goal
Test goal

## Constraints
- Must handle "quotes" and 'apostrophes'
- Must support unicode: 中文 🚀
- Must handle special chars: @#$%^&*()
`)

	ir, err := Compile(planContent)
	if err != nil {
		t.Fatalf("Compile() failed with special characters: %v", err)
	}
	
	// All rules should have valid IDs (no special chars in IDs)
	for _, rule := range ir.Rules {
		if strings.ContainsAny(rule.ID, "@#$%^&*()") {
			t.Errorf("Rule ID contains special characters: %s", rule.ID)
		}
	}
}

// TestCompile_EmptyConstraints tests that empty constraints are handled.
func TestCompile_EmptyConstraints(t *testing.T) {
	planContent := []byte(`# Prompt Plan

## Goal
Test goal

## Constraints

## Out of Scope
`)

	ir, err := Compile(planContent)
	if err != nil {
		t.Fatalf("Compile() failed: %v", err)
	}
	
	// Should still have baseline rules
	if len(ir.Rules) < 4 {
		t.Errorf("Expected at least 4 baseline rules, got %d", len(ir.Rules))
	}
}

// TestCompile_ManyConstraints tests handling of many constraints.
func TestCompile_ManyConstraints(t *testing.T) {
	var constraints []string
	for i := 0; i < 50; i++ {
		constraints = append(constraints, "- Must handle constraint number "+string(rune('A'+i%26)))
	}
	
	planContent := []byte(`# Prompt Plan

## Goal
Test goal

## Constraints
` + strings.Join(constraints, "\n") + `

## Out of Scope
`)

	ir, err := Compile(planContent)
	if err != nil {
		t.Fatalf("Compile() failed with many constraints: %v", err)
	}
	
	// Should have baseline rules + some constraints (parser may filter duplicates)
	if len(ir.Rules) < 4 {
		t.Errorf("Expected at least 4 baseline rules, got %d", len(ir.Rules))
	}
	
	// All IDs should be unique
	ruleIDs := make(map[string]bool)
	for _, rule := range ir.Rules {
		if ruleIDs[rule.ID] {
			t.Errorf("Duplicate rule ID found: %s", rule.ID)
		}
		ruleIDs[rule.ID] = true
	}
}

// TestValidateIR_ValidIR tests that ValidateIR passes for a valid IR.
func TestValidateIR_ValidIR(t *testing.T) {
	validIR := &ir.PromptIR{
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{
				ID:          "rule-1",
				Description: "Test rule",
			},
		},
		InputSchema: ir.Schema{
			Type: "object",
		},
		OutputSchema: ir.Schema{
			Type: "object",
		},
		FailureModes: []ir.FailureMode{
			{
				ID:        "fm-1",
				Condition: "Test condition",
				Response:  "Test response",
			},
		},
	}

	if err := ValidateIR(validIR); err != nil {
		t.Fatalf("ValidateIR() should pass for valid IR: %v", err)
	}
}

// TestValidateIR_NilIR tests that ValidateIR fails on nil IR.
func TestValidateIR_NilIR(t *testing.T) {
	if err := ValidateIR(nil); err == nil {
		t.Error("ValidateIR() should fail on nil IR")
	}
}

// TestValidateIR_EmptySystemRole tests that ValidateIR fails when system_role is empty.
func TestValidateIR_EmptySystemRole(t *testing.T) {
	ir := &ir.PromptIR{
		SystemRole: "",
		Rules: []ir.Rule{
			{ID: "rule-1", Description: "Test"},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "Test", Response: "Test"},
		},
	}

	if err := ValidateIR(ir); err == nil {
		t.Error("ValidateIR() should fail when system_role is empty")
	}
}

// TestValidateIR_EmptyRules tests that ValidateIR fails when rules array is empty.
func TestValidateIR_EmptyRules(t *testing.T) {
	ir := &ir.PromptIR{
		SystemRole: "Test role",
		Rules:      []ir.Rule{},
		InputSchema: ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "Test", Response: "Test"},
		},
	}

	if err := ValidateIR(ir); err == nil {
		t.Error("ValidateIR() should fail when rules array is empty")
	}
}

// TestValidateIR_RuleMissingID tests that ValidateIR fails when a rule is missing ID.
func TestValidateIR_RuleMissingID(t *testing.T) {
	ir := &ir.PromptIR{
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{ID: "", Description: "Test"},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "Test", Response: "Test"},
		},
	}

	if err := ValidateIR(ir); err == nil {
		t.Error("ValidateIR() should fail when rule ID is empty")
	}
}

// TestValidateIR_RuleMissingDescription tests that ValidateIR fails when a rule is missing description.
func TestValidateIR_RuleMissingDescription(t *testing.T) {
	ir := &ir.PromptIR{
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{ID: "rule-1", Description: ""},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "Test", Response: "Test"},
		},
	}

	if err := ValidateIR(ir); err == nil {
		t.Error("ValidateIR() should fail when rule description is empty")
	}
}

// TestValidateIR_EmptyInputSchemaType tests that ValidateIR fails when input_schema.type is empty.
func TestValidateIR_EmptyInputSchemaType(t *testing.T) {
	ir := &ir.PromptIR{
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{ID: "rule-1", Description: "Test"},
		},
		InputSchema:  ir.Schema{Type: ""},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "Test", Response: "Test"},
		},
	}

	if err := ValidateIR(ir); err == nil {
		t.Error("ValidateIR() should fail when input_schema.type is empty")
	}
}

// TestValidateIR_EmptyOutputSchemaType tests that ValidateIR fails when output_schema.type is empty.
func TestValidateIR_EmptyOutputSchemaType(t *testing.T) {
	ir := &ir.PromptIR{
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{ID: "rule-1", Description: "Test"},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: ""},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "Test", Response: "Test"},
		},
	}

	if err := ValidateIR(ir); err == nil {
		t.Error("ValidateIR() should fail when output_schema.type is empty")
	}
}

// TestValidateIR_EmptyFailureModes tests that ValidateIR fails when failure_modes array is empty.
func TestValidateIR_EmptyFailureModes(t *testing.T) {
	ir := &ir.PromptIR{
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{ID: "rule-1", Description: "Test"},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{},
	}

	if err := ValidateIR(ir); err == nil {
		t.Error("ValidateIR() should fail when failure_modes array is empty")
	}
}

// TestValidateIR_FailureModeMissingID tests that ValidateIR fails when a failure mode is missing ID.
func TestValidateIR_FailureModeMissingID(t *testing.T) {
	ir := &ir.PromptIR{
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{ID: "rule-1", Description: "Test"},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "", Condition: "Test", Response: "Test"},
		},
	}

	if err := ValidateIR(ir); err == nil {
		t.Error("ValidateIR() should fail when failure mode ID is empty")
	}
}

// TestValidateIR_FailureModeMissingCondition tests that ValidateIR fails when a failure mode is missing condition.
func TestValidateIR_FailureModeMissingCondition(t *testing.T) {
	ir := &ir.PromptIR{
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{ID: "rule-1", Description: "Test"},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "", Response: "Test"},
		},
	}

	if err := ValidateIR(ir); err == nil {
		t.Error("ValidateIR() should fail when failure mode condition is empty")
	}
}

// TestValidateIR_FailureModeMissingResponse tests that ValidateIR fails when a failure mode is missing response.
func TestValidateIR_FailureModeMissingResponse(t *testing.T) {
	ir := &ir.PromptIR{
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{ID: "rule-1", Description: "Test"},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "Test", Response: ""},
		},
	}

	if err := ValidateIR(ir); err == nil {
		t.Error("ValidateIR() should fail when failure mode response is empty")
	}
}

// TestWriteIR_InvalidIR tests that WriteIR fails when given an invalid IR.
func TestWriteIR_InvalidIR(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "prompt.ir.json")

	invalidIR := &ir.PromptIR{
		SystemRole: "", // Invalid: empty system_role
		Rules: []ir.Rule{
			{ID: "rule-1", Description: "Test"},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "Test", Response: "Test"},
		},
	}

	err := WriteIR(invalidIR, outputPath)
	if err == nil {
		t.Error("WriteIR() should fail when IR is invalid")
	}
}

// TestWriteIR_EmptyPath tests that WriteIR fails with empty path.
func TestWriteIR_EmptyPath(t *testing.T) {
	validIR := &ir.PromptIR{
		SystemRole: "Test role",
		Rules: []ir.Rule{
			{ID: "rule-1", Description: "Test"},
		},
		InputSchema:  ir.Schema{Type: "object"},
		OutputSchema: ir.Schema{Type: "object"},
		FailureModes: []ir.FailureMode{
			{ID: "fm-1", Condition: "Test", Response: "Test"},
		},
	}

	err := WriteIR(validIR, "")
	if err == nil {
		t.Error("WriteIR() should fail with empty path")
	}
}
