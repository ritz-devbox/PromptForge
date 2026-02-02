package templates

import "fmt"

type Template struct {
	Name        string
	Description string
	Plan        string
}

var registry = []Template{
	{
		Name:        "api-guardrails",
		Description: "API assistant with strict input/output and error handling",
		Plan: `# Prompt Plan

## Goal
Design a reliable API assistant that validates inputs and returns structured JSON responses.

## Constraints
- Output must be valid JSON
- Validate all required fields
- Reject ambiguous requests with a clear error

## Out of Scope
- Payments
- Authentication flows
`,
	},
	{
		Name:        "data-extraction",
		Description: "Extract structured fields from unstructured text",
		Plan: `# Prompt Plan

## Goal
Extract structured fields from unstructured text into a strict JSON schema.

## Constraints
- Do not infer missing values
- Return a validation error if required fields are missing
- Preserve original wording in extracted fields

## Out of Scope
- Summarization beyond field extraction
- Rewriting or paraphrasing the input
`,
	},
	{
		Name:        "code-review",
		Description: "Review code for correctness, security, and tests",
		Plan: `# Prompt Plan

## Goal
Provide a concise, high-signal code review focused on correctness, security, and tests.

## Constraints
- Prioritize critical issues over style
- Provide actionable remediation steps
- Call out missing tests

## Out of Scope
- Large refactors without a clear bug
- Formatting-only feedback
`,
	},
}

// List returns all available templates.
func List() []Template {
	return registry
}

// Get returns a template by name.
func Get(name string) (Template, error) {
	for _, tmpl := range registry {
		if tmpl.Name == name {
			return tmpl, nil
		}
	}
	return Template{}, fmt.Errorf("unknown template: %s", name)
}
