package linter

import (
	"strings"
	"testing"
)

func TestLintPlan_MissingGoal(t *testing.T) {
	content := []byte("# Prompt Plan\n\n## Constraints\n- Be strict\n")
	diags := LintPlan(content)

	if !hasCode(diags, "PF100") {
		t.Fatal("expected PF100 missing Goal diagnostic")
	}
}

func TestLintPlan_EmptyGoal(t *testing.T) {
	content := []byte("# Prompt Plan\n\n## Goal\n\n## Constraints\n- Be strict\n")
	diags := LintPlan(content)

	if !hasCode(diags, "PF101") {
		t.Fatal("expected PF101 empty Goal diagnostic")
	}
}

func TestLintPlan_DuplicateSection(t *testing.T) {
	content := []byte("# Prompt Plan\n\n## Goal\nOne\n\n## Goal\nTwo\n")
	diags := LintPlan(content)

	if !hasCode(diags, "PF102") {
		t.Fatal("expected PF102 duplicate section diagnostic")
	}
}

func TestLintPlan_UnknownSection(t *testing.T) {
	content := []byte("# Prompt Plan\n\n## Notes\nSomething\n")
	diags := LintPlan(content)

	if !hasCode(diags, "PF103") {
		t.Fatal("expected PF103 unknown section diagnostic")
	}
}

func TestLintPlan_ShortGoalWarning(t *testing.T) {
	content := []byte("# Prompt Plan\n\n## Goal\nShort\n")
	diags := LintPlan(content)

	if !hasCode(diags, "PF202") {
		t.Fatal("expected PF202 short Goal warning")
	}
}

func TestLintPlan_MissingOptionalSectionsWarnings(t *testing.T) {
	content := []byte("# Prompt Plan\n\n## Goal\nA clear goal statement.\n")
	diags := LintPlan(content)

	if !hasCode(diags, "PF200") {
		t.Fatal("expected PF200 constraints warning")
	}
	if !hasCode(diags, "PF201") {
		t.Fatal("expected PF201 out of scope warning")
	}
}

func TestLintPlan_VagueConstraintWarning(t *testing.T) {
	content := []byte(strings.Join([]string{
		"# Prompt Plan",
		"",
		"## Goal",
		"A detailed goal statement for testing.",
		"",
		"## Constraints",
		"- Use various approaches",
		"",
		"## Out of Scope",
		"- Nothing else",
	}, "\n"))

	diags := LintPlan(content)
	if !hasCode(diags, "PF203") {
		t.Fatal("expected PF203 vague constraint warning")
	}
}

func hasCode(diags []Diagnostic, code string) bool {
	for _, diag := range diags {
		if diag.Code == code {
			return true
		}
	}
	return false
}
