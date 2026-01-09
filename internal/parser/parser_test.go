package parser

import (
	"testing"
)

func TestParsePlan_CompletePlan(t *testing.T) {
	content := `# Prompt Plan

## Goal
I want a chatbot that helps users debug API errors

## Constraints
- Must provide code examples in Python and JavaScript
- Must validate API endpoints before suggesting fixes
- Must not execute code, only show examples

## Out of Scope
- Does not handle authentication issues
- Does not debug network connectivity problems
`

	plan, err := ParsePlan([]byte(content))
	if err != nil {
		t.Fatalf("ParsePlan() failed: %v", err)
	}

	if plan.Goal != "I want a chatbot that helps users debug API errors" {
		t.Errorf("Goal mismatch: got %q", plan.Goal)
	}

	if len(plan.Constraints) != 3 {
		t.Errorf("Expected 3 constraints, got %d", len(plan.Constraints))
	}

	if len(plan.OutOfScope) != 2 {
		t.Errorf("Expected 2 out of scope items, got %d", len(plan.OutOfScope))
	}
}

func TestParsePlan_WithComments(t *testing.T) {
	content := `# Prompt Plan

## Goal
Create a web application

## Constraints
<!-- List the limitations -->
- Must work on mobile devices
- Must be secure

## Out of Scope
<!-- What not to handle -->
- Does not support offline mode
`

	plan, err := ParsePlan([]byte(content))
	if err != nil {
		t.Fatalf("ParsePlan() failed: %v", err)
	}

	// Comments should be removed
	if len(plan.Constraints) != 2 {
		t.Errorf("Expected 2 constraints (comments removed), got %d", len(plan.Constraints))
	}
}

func TestParsePlan_EmptyContent(t *testing.T) {
	content := []byte("")
	
	plan, err := ParsePlan(content)
	if err == nil {
		t.Error("ParsePlan() should fail on empty content")
	}
	if plan != nil {
		t.Error("ParsePlan() should return nil plan on error")
	}
	if err.Error() != "plan content is empty" {
		t.Errorf("Expected 'plan content is empty', got %q", err.Error())
	}
}

func TestParsePlan_MissingGoal(t *testing.T) {
	content := `# Prompt Plan

## Constraints
- Some constraint

## Out of Scope
- Some item
`

	plan, err := ParsePlan([]byte(content))
	if err == nil {
		t.Error("ParsePlan() should fail when Goal section is missing")
	}
	if plan != nil {
		t.Error("ParsePlan() should return nil plan on error")
	}
}

func TestParsePlan_EmptyGoal(t *testing.T) {
	content := `# Prompt Plan

## Goal
<!-- Empty goal -->

## Constraints
- Some constraint
`

	plan, err := ParsePlan([]byte(content))
	if err == nil {
		t.Error("ParsePlan() should fail when Goal section is empty")
	}
	if plan != nil {
		t.Error("ParsePlan() should return nil plan on error")
	}
}

func TestParsePlan_WindowsLineEndings(t *testing.T) {
	content := "# Prompt Plan\r\n\r\n## Goal\r\nTest goal\r\n\r\n## Constraints\r\n- Constraint 1\r\n"
	
	plan, err := ParsePlan([]byte(content))
	if err != nil {
		t.Fatalf("ParsePlan() failed with Windows line endings: %v", err)
	}
	
	if plan.Goal != "Test goal" {
		t.Errorf("Goal mismatch: got %q", plan.Goal)
	}
}

func TestParsePlan_SpecialCharacters(t *testing.T) {
	content := `# Prompt Plan

## Goal
Test with "quotes" and 'apostrophes' and unicode: 测试 🚀

## Constraints
- Must handle "special" characters
- Must support unicode: 中文
`

	plan, err := ParsePlan([]byte(content))
	if err != nil {
		t.Fatalf("ParsePlan() failed with special characters: %v", err)
	}
	
	if plan.Goal == "" {
		t.Error("Goal should not be empty")
	}
	if len(plan.Constraints) != 2 {
		t.Errorf("Expected 2 constraints, got %d", len(plan.Constraints))
	}
}

func TestParsePlan_VeryLongContent(t *testing.T) {
	// Create a very long goal
	longGoal := "Test goal " + string(make([]byte, 10000))
	content := `# Prompt Plan

## Goal
` + longGoal + `

## Constraints
- Constraint 1
`

	plan, err := ParsePlan([]byte(content))
	if err != nil {
		t.Fatalf("ParsePlan() failed with long content: %v", err)
	}
	
	if len(plan.Goal) == 0 {
		t.Error("Goal should not be empty")
	}
}

func TestParsePlan_DuplicateConstraints(t *testing.T) {
	content := `# Prompt Plan

## Goal
Test goal

## Constraints
- Must be secure
- Must be secure
- Must be secure
`

	plan, err := ParsePlan([]byte(content))
	if err != nil {
		t.Fatalf("ParsePlan() failed: %v", err)
	}
	
	// Duplicates should be filtered out
	if len(plan.Constraints) != 1 {
		t.Errorf("Expected 1 unique constraint (duplicates filtered), got %d", len(plan.Constraints))
	}
}

func TestParsePlan_EmptySections(t *testing.T) {
	content := `# Prompt Plan

## Goal
Test goal

## Constraints

## Out of Scope
`

	plan, err := ParsePlan([]byte(content))
	if err != nil {
		t.Fatalf("ParsePlan() failed: %v", err)
	}
	
	// Empty sections should result in empty arrays
	if len(plan.Constraints) != 0 {
		t.Errorf("Expected 0 constraints, got %d", len(plan.Constraints))
	}
	if len(plan.OutOfScope) != 0 {
		t.Errorf("Expected 0 out of scope items, got %d", len(plan.OutOfScope))
	}
}

func TestParsePlan_MalformedMarkdown(t *testing.T) {
	content := `# Prompt Plan

##Goal
Test goal (missing space)

## Constraints
- Constraint 1
`

	plan, err := ParsePlan([]byte(content))
	if err == nil {
		t.Error("ParsePlan() should fail with malformed markdown (missing space in header)")
	}
	if plan != nil {
		t.Error("ParsePlan() should return nil plan on error")
	}
}

func TestParsePlan_MultipleSections(t *testing.T) {
	content := `# Prompt Plan

## Goal
First goal

## Goal
Second goal

## Constraints
- Constraint 1
`

	plan, err := ParsePlan([]byte(content))
	if err != nil {
		t.Fatalf("ParsePlan() failed: %v", err)
	}
	
	// Should use the first Goal section
	if plan.Goal != "First goal" {
		t.Errorf("Expected first goal, got %q", plan.Goal)
	}
}

func TestParseList_EmptyText(t *testing.T) {
	result := parseList("")
	if len(result) != 0 {
		t.Errorf("Expected empty list, got %d items", len(result))
	}
}

func TestParseList_OnlyComments(t *testing.T) {
	result := parseList("<!-- Comment -->\n<!-- Another comment -->")
	if len(result) != 0 {
		t.Errorf("Expected empty list (comments filtered), got %d items", len(result))
	}
}

func TestParseList_NumberedList(t *testing.T) {
	result := parseList("1. First item\n2. Second item\n3. Third item")
	if len(result) != 3 {
		t.Errorf("Expected 3 items, got %d", len(result))
	}
}

func TestParseList_MixedMarkers(t *testing.T) {
	result := parseList("- Bullet 1\n* Bullet 2\n• Bullet 3\n1. Numbered 1")
	if len(result) != 4 {
		t.Errorf("Expected 4 items, got %d", len(result))
	}
}
