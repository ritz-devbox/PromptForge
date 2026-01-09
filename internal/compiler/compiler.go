package compiler

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/promptforge/promptforge/internal/ir"
	"github.com/promptforge/promptforge/internal/parser"
)

// Compile transforms human-readable plan content into a PromptIR.
// This is the core compilation step: intent → IR.
//
// The compiler parses plan.md and generates a PromptIR that reflects
// the user's intent expressed in Goal, Constraints, and Out of Scope sections.
func Compile(planContent []byte) (*ir.PromptIR, error) {
	// Verify planContent is not empty
	if len(planContent) == 0 {
		return nil, fmt.Errorf("plan.md is empty")
	}

	// Parse the plan.md content
	plan, err := parser.ParsePlan(planContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse plan: %w", err)
	}

	// Generate system_role from Goal
	systemRole := generateSystemRole(plan.Goal)

	// Always include baseline rules to ensure deterministic behavior
	baselineRules := []ir.Rule{
		{
			ID:          "output-json",
			Description: "Output must be valid JSON",
		},
		{
			ID:          "no-explanations",
			Description: "Do not include explanations unless explicitly requested",
		},
		{
			ID:          "no-inference",
			Description: "Do not infer missing values - fail if required data is missing",
		},
		{
			ID:          "fail-ambiguity",
			Description: "Fail on ambiguity - request clarification if intent is unclear",
		},
	}

	// Track existing rule IDs to prevent collisions
	existingRuleIDs := make(map[string]bool)
	for _, rule := range baselineRules {
		existingRuleIDs[rule.ID] = true
	}

	// Generate rules from Constraints (with collision detection)
	rules := generateRules(plan.Constraints, existingRuleIDs)

	// Combine baseline rules with user-defined constraints
	allRules := append(baselineRules, rules...)

	// Always include baseline failure modes
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

	// Track existing failure mode IDs to prevent collisions
	existingFailureModeIDs := make(map[string]bool)
	for _, fm := range baselineFailureModes {
		existingFailureModeIDs[fm.ID] = true
	}

	// Generate failure_modes from Out of Scope (with collision detection)
	failureModes := generateFailureModes(plan.OutOfScope, existingFailureModeIDs)

	allFailureModes := append(baselineFailureModes, failureModes...)

	// Generate schemas (basic structure for now, can be enhanced)
	inputSchema, outputSchema := generateSchemas(plan)

	return &ir.PromptIR{
		SystemRole:   systemRole,
		Rules:        allRules,
		InputSchema:  inputSchema,
		OutputSchema: outputSchema,
		FailureModes: allFailureModes,
	}, nil
}

// generateSystemRole creates a system role description from the Goal
func generateSystemRole(goal string) string {
	if goal == "" {
		return "You are a deterministic assistant that follows strict rules and schemas."
	}

	// Format the goal as a system role
	goal = strings.TrimSpace(goal)
	
	// Capitalize first letter if needed
	if len(goal) > 0 && goal[0] >= 'a' && goal[0] <= 'z' {
		goal = strings.ToUpper(string(goal[0])) + goal[1:]
	}

	// Ensure it ends with proper punctuation
	if !strings.HasSuffix(goal, ".") && !strings.HasSuffix(goal, "!") && !strings.HasSuffix(goal, "?") {
		goal += "."
	}

	return fmt.Sprintf("You are an assistant designed to: %s You must follow all specified rules and constraints strictly.", goal)
}

// generateRules converts constraints into Rule objects.
// existingIDs tracks IDs already in use to prevent collisions.
func generateRules(constraints []string, existingIDs map[string]bool) []ir.Rule {
	if len(constraints) == 0 {
		return []ir.Rule{}
	}

	rules := make([]ir.Rule, 0, len(constraints))
	for i, constraint := range constraints {
		constraint = strings.TrimSpace(constraint)
		if constraint == "" {
			continue
		}

		// Generate a unique rule ID from the constraint
		ruleID := generateUniqueRuleID(constraint, i, existingIDs)
		existingIDs[ruleID] = true

		rules = append(rules, ir.Rule{
			ID:          ruleID,
			Description: constraint,
		})
	}

	return rules
}

// generateUniqueRuleID creates a unique ID for a rule based on its content.
// If the generated ID collides with existingIDs, appends a suffix to make it unique.
func generateUniqueRuleID(description string, index int, existingIDs map[string]bool) string {
	baseID := generateRuleID(description, index)
	
	// If ID doesn't exist, return it
	if !existingIDs[baseID] {
		return baseID
	}
	
	// If collision, append suffix
	suffix := 1
	for {
		candidateID := fmt.Sprintf("%s-%d", baseID, suffix)
		if !existingIDs[candidateID] {
			return candidateID
		}
		suffix++
		// Safety: prevent infinite loop
		if suffix > 1000 {
			return fmt.Sprintf("constraint-%d", index)
		}
	}
}

// generateRuleID creates a base ID for a rule based on its content.
func generateRuleID(description string, index int) string {
	// Extract key words from description
	// First remove all punctuation and special characters, then split
	cleaned := strings.ToLower(description)
	// Remove punctuation and special characters (keep only letters, numbers, spaces)
	cleaned = regexp.MustCompile(`[^a-z0-9\s]`).ReplaceAllString(cleaned, " ")
	// Remove extra spaces
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	words := strings.Fields(cleaned)
	
	// Take first few meaningful words
	var keyWords []string
	for _, word := range words {
		// Skip common words
		if word == "the" || word == "a" || word == "an" || word == "is" || 
		   word == "are" || word == "must" || word == "should" || word == "can" ||
		   word == "will" || word == "be" || word == "to" || word == "of" ||
		   word == "not" || word == "only" {
			continue
		}
		if len(keyWords) < 3 {
			keyWords = append(keyWords, word)
		}
	}

	// Create ID from key words
	if len(keyWords) > 0 {
		id := strings.Join(keyWords, "-")
		// Clean up ID - ensure only alphanumeric and hyphens
		id = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(id, "")
		id = strings.Trim(id, "-")
		if len(id) > 30 {
			id = id[:30]
		}
		if id != "" {
			return fmt.Sprintf("constraint-%s", id)
		}
	}

	return fmt.Sprintf("constraint-%d", index)
}

// generateFailureModes converts out of scope items into FailureMode objects.
// existingIDs tracks IDs already in use to prevent collisions.
func generateFailureModes(outOfScope []string, existingIDs map[string]bool) []ir.FailureMode {
	if len(outOfScope) == 0 {
		return []ir.FailureMode{}
	}

	failureModes := make([]ir.FailureMode, 0, len(outOfScope))
	for i, item := range outOfScope {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		// Generate unique failure mode ID
		fmID := generateUniqueFailureModeID(item, i, existingIDs)
		existingIDs[fmID] = true

		// Format the condition and response
		condition := fmt.Sprintf("Request involves: %s", item)
		response := fmt.Sprintf("Return error indicating that %s is out of scope and cannot be handled", item)

		failureModes = append(failureModes, ir.FailureMode{
			ID:        fmID,
			Condition: condition,
			Response:  response,
		})
	}

	return failureModes
}

// generateUniqueFailureModeID creates a unique ID for a failure mode.
// If the generated ID collides with existingIDs, appends a suffix to make it unique.
func generateUniqueFailureModeID(description string, index int, existingIDs map[string]bool) string {
	baseID := generateFailureModeID(description, index)
	
	// If ID doesn't exist, return it
	if !existingIDs[baseID] {
		return baseID
	}
	
	// If collision, append suffix
	suffix := 1
	for {
		candidateID := fmt.Sprintf("%s-%d", baseID, suffix)
		if !existingIDs[candidateID] {
			return candidateID
		}
		suffix++
		// Safety: prevent infinite loop
		if suffix > 1000 {
			return fmt.Sprintf("out-of-scope-%d", index)
		}
	}
}

// generateFailureModeID creates a base ID for a failure mode.
func generateFailureModeID(description string, index int) string {
	// Remove "Does not" prefix and clean
	cleaned := regexp.MustCompile(`(?i)^does\s+not\s+`).ReplaceAllString(description, "")
	cleaned = regexp.MustCompile(`[^\w\s]`).ReplaceAllString(strings.ToLower(cleaned), " ")
	words := strings.Fields(cleaned)
	
	var keyWords []string
	for _, word := range words {
		if word == "handle" || word == "support" || word == "debug" {
			continue
		}
		if len(keyWords) < 2 {
			keyWords = append(keyWords, word)
		}
	}

	if len(keyWords) > 0 {
		id := strings.Join(keyWords, "-")
		// Clean up ID - ensure only alphanumeric and hyphens
		id = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(id, "")
		id = strings.Trim(id, "-")
		if len(id) > 25 {
			id = id[:25]
		}
		if id != "" {
			return fmt.Sprintf("out-of-scope-%s", id)
		}
	}

	return fmt.Sprintf("out-of-scope-%d", index)
}

// generateSchemas creates basic input and output schemas
// This is a simplified version - can be enhanced to extract schema info from plan
func generateSchemas(plan *parser.Plan) (ir.Schema, ir.Schema) {
	// For now, return basic object schemas
	// Future enhancement: parse plan content to extract schema requirements
	return ir.Schema{
			Type:       "object",
			Properties: map[string]ir.Property{},
			Required:   []string{},
		}, ir.Schema{
			Type:       "object",
			Properties: map[string]ir.Property{},
			Required:   []string{},
		}
}

// ValidateIR ensures all required fields are present and valid.
// Fails compilation if validation fails.
func ValidateIR(promptIR *ir.PromptIR) error {
	if promptIR == nil {
		return fmt.Errorf("prompt IR is nil")
	}

	// Validate system_role is present and non-empty
	if promptIR.SystemRole == "" {
		return fmt.Errorf("system_role is required and cannot be empty")
	}

	// Validate rules array is not empty
	if len(promptIR.Rules) == 0 {
		return fmt.Errorf("rules array cannot be empty")
	}

	// Validate each rule has required fields
	for i, rule := range promptIR.Rules {
		if rule.ID == "" {
			return fmt.Errorf("rules[%d].id is required and cannot be empty", i)
		}
		if rule.Description == "" {
			return fmt.Errorf("rules[%d].description is required and cannot be empty", i)
		}
	}

	// Validate input_schema has type
	if promptIR.InputSchema.Type == "" {
		return fmt.Errorf("input_schema.type is required and cannot be empty")
	}

	// Validate output_schema has type
	if promptIR.OutputSchema.Type == "" {
		return fmt.Errorf("output_schema.type is required and cannot be empty")
	}

	// Validate failure_modes array is not empty
	if len(promptIR.FailureModes) == 0 {
		return fmt.Errorf("failure_modes array cannot be empty")
	}

	// Validate each failure mode has required fields
	for i, fm := range promptIR.FailureModes {
		if fm.ID == "" {
			return fmt.Errorf("failure_modes[%d].id is required and cannot be empty", i)
		}
		if fm.Condition == "" {
			return fmt.Errorf("failure_modes[%d].condition is required and cannot be empty", i)
		}
		if fm.Response == "" {
			return fmt.Errorf("failure_modes[%d].response is required and cannot be empty", i)
		}
	}

	return nil
}

// WriteIR writes the PromptIR to a JSON file.
// Returns detailed error messages for file system issues.
func WriteIR(promptIR *ir.PromptIR, outputPath string) error {
	// Validate IR before writing - fail compilation if validation fails
	if err := ValidateIR(promptIR); err != nil {
		return fmt.Errorf("IR validation failed: %w", err)
	}

	// Validate output path
	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	data, err := json.MarshalIndent(promptIR, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal IR to JSON: %w", err)
	}

	// Write file with better error handling
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied: cannot write to %s", outputPath)
		}
		// Check for disk full (platform-specific)
		if err.Error() == "no space left on device" || err.Error() == "not enough space" {
			return fmt.Errorf("disk full: cannot write to %s", outputPath)
		}
		return fmt.Errorf("failed to write IR file to %s: %w", outputPath, err)
	}

	return nil
}

