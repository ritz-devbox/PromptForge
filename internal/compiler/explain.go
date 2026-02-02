package compiler

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/promptforge/promptforge/internal/ir"
	"github.com/promptforge/promptforge/internal/parser"
)

type ExplainReport struct {
	SystemRole   ExplainSystemRole    `json:"system_role"`
	Rules        []ExplainRule        `json:"rules"`
	InputSchema  ExplainSchema        `json:"input_schema"`
	OutputSchema ExplainSchema        `json:"output_schema"`
	FailureModes []ExplainFailureMode `json:"failure_modes"`
}

type ExplainSource struct {
	Type    string `json:"type"`
	Section string `json:"section,omitempty"`
	Line    int    `json:"line,omitempty"`
}

type ExplainSystemRole struct {
	Value  string        `json:"value"`
	Source ExplainSource `json:"source"`
}

type ExplainRule struct {
	ID          string        `json:"id"`
	Description string        `json:"description"`
	Source      ExplainSource `json:"source"`
}

type ExplainSchema struct {
	Source ExplainSource `json:"source"`
}

type ExplainFailureMode struct {
	ID        string        `json:"id"`
	Condition string        `json:"condition"`
	Response  string        `json:"response"`
	Source    ExplainSource `json:"source"`
}

// CompileWithExplain compiles plan content and returns an explain report.
func CompileWithExplain(planContent []byte) (*ir.PromptIR, *ExplainReport, error) {
	if len(planContent) == 0 {
		return nil, nil, fmt.Errorf("plan.md is empty")
	}

	plan, err := parser.ParsePlanWithLines(planContent)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse plan: %w", err)
	}

	systemRole := generateSystemRole(plan.Plan.Goal)
	baselineRules := []ir.Rule{
		{ID: "output-json", Description: "Output must be valid JSON"},
		{ID: "no-explanations", Description: "Do not include explanations unless explicitly requested"},
		{ID: "no-inference", Description: "Do not infer missing values - fail if required data is missing"},
		{ID: "fail-ambiguity", Description: "Fail on ambiguity - request clarification if intent is unclear"},
	}

	existingRuleIDs := make(map[string]bool)
	for _, rule := range baselineRules {
		existingRuleIDs[rule.ID] = true
	}

	rules, rulesExplain := generateRulesWithExplain(plan.Constraints, existingRuleIDs)
	allRules := append(baselineRules, rules...)

	baselineFailureModes := []ir.FailureMode{
		{
			ID:        "invalid-input",
			Condition: "Input does not match input_schema",
			Response:  "Return error indicating schema validation failure",
		},
		{
			ID:        "ambiguous-request",
			Condition: "Request cannot be unambiguously interpreted",
			Response:  "Return error indicating ambiguity and request clarification",
		},
		{
			ID:        "missing-required",
			Condition: "Required fields are missing from input",
			Response:  "Return error listing missing required fields",
		},
	}

	existingFailureModeIDs := make(map[string]bool)
	for _, fm := range baselineFailureModes {
		existingFailureModeIDs[fm.ID] = true
	}

	failureModes, failureExplain := generateFailureModesWithExplain(plan.OutOfScope, existingFailureModeIDs)
	allFailureModes := append(baselineFailureModes, failureModes...)

	inputSchema, outputSchema := generateSchemas(plan.Plan)

	report := &ExplainReport{
		SystemRole: ExplainSystemRole{
			Value: systemRole,
			Source: ExplainSource{
				Type:    "plan",
				Section: "Goal",
				Line:    plan.GoalLine,
			},
		},
		Rules: append([]ExplainRule{
			{
				ID:          baselineRules[0].ID,
				Description: baselineRules[0].Description,
				Source:      ExplainSource{Type: "baseline"},
			},
			{
				ID:          baselineRules[1].ID,
				Description: baselineRules[1].Description,
				Source:      ExplainSource{Type: "baseline"},
			},
			{
				ID:          baselineRules[2].ID,
				Description: baselineRules[2].Description,
				Source:      ExplainSource{Type: "baseline"},
			},
			{
				ID:          baselineRules[3].ID,
				Description: baselineRules[3].Description,
				Source:      ExplainSource{Type: "baseline"},
			},
		}, rulesExplain...),
		InputSchema: ExplainSchema{
			Source: ExplainSource{Type: "baseline"},
		},
		OutputSchema: ExplainSchema{
			Source: ExplainSource{Type: "baseline"},
		},
		FailureModes: append([]ExplainFailureMode{
			{
				ID:        baselineFailureModes[0].ID,
				Condition: baselineFailureModes[0].Condition,
				Response:  baselineFailureModes[0].Response,
				Source:    ExplainSource{Type: "baseline"},
			},
			{
				ID:        baselineFailureModes[1].ID,
				Condition: baselineFailureModes[1].Condition,
				Response:  baselineFailureModes[1].Response,
				Source:    ExplainSource{Type: "baseline"},
			},
			{
				ID:        baselineFailureModes[2].ID,
				Condition: baselineFailureModes[2].Condition,
				Response:  baselineFailureModes[2].Response,
				Source:    ExplainSource{Type: "baseline"},
			},
		}, failureExplain...),
	}

	return &ir.PromptIR{
		Version:      ir.CurrentVersion,
		SystemRole:   systemRole,
		Rules:        allRules,
		InputSchema:  inputSchema,
		OutputSchema: outputSchema,
		FailureModes: allFailureModes,
	}, report, nil
}

func generateRulesWithExplain(constraints []parser.PlanItem, existingIDs map[string]bool) ([]ir.Rule, []ExplainRule) {
	if len(constraints) == 0 {
		return []ir.Rule{}, []ExplainRule{}
	}

	rules := make([]ir.Rule, 0, len(constraints))
	explain := make([]ExplainRule, 0, len(constraints))

	for i, constraint := range constraints {
		description := strings.TrimSpace(constraint.Text)
		if description == "" {
			continue
		}

		ruleID := generateUniqueRuleID(description, i, existingIDs)
		existingIDs[ruleID] = true

		rules = append(rules, ir.Rule{
			ID:          ruleID,
			Description: description,
		})

		explain = append(explain, ExplainRule{
			ID:          ruleID,
			Description: description,
			Source: ExplainSource{
				Type:    "plan",
				Section: "Constraints",
				Line:    constraint.Line,
			},
		})
	}

	return rules, explain
}

func generateFailureModesWithExplain(items []parser.PlanItem, existingIDs map[string]bool) ([]ir.FailureMode, []ExplainFailureMode) {
	if len(items) == 0 {
		return []ir.FailureMode{}, []ExplainFailureMode{}
	}

	failureModes := make([]ir.FailureMode, 0, len(items))
	explain := make([]ExplainFailureMode, 0, len(items))

	for i, item := range items {
		text := strings.TrimSpace(item.Text)
		if text == "" {
			continue
		}

		fmID := generateUniqueFailureModeID(text, i, existingIDs)
		existingIDs[fmID] = true

		condition := fmt.Sprintf("Request involves: %s", text)
		response := fmt.Sprintf("Return error indicating that %s is out of scope and cannot be handled", text)

		failureModes = append(failureModes, ir.FailureMode{
			ID:        fmID,
			Condition: condition,
			Response:  response,
		})

		explain = append(explain, ExplainFailureMode{
			ID:        fmID,
			Condition: condition,
			Response:  response,
			Source: ExplainSource{
				Type:    "plan",
				Section: "Out of Scope",
				Line:    item.Line,
			},
		})
	}

	return failureModes, explain
}

// WriteExplainReport writes the explain report to a JSON file.
func WriteExplainReport(report *ExplainReport, outputPath string) error {
	if report == nil {
		return fmt.Errorf("explain report is nil")
	}
	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal explain report: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied: cannot write to %s", outputPath)
		}
		if err.Error() == "no space left on device" || err.Error() == "not enough space" {
			return fmt.Errorf("disk full: cannot write to %s", outputPath)
		}
		return fmt.Errorf("failed to write explain report to %s: %w", outputPath, err)
	}

	return nil
}
