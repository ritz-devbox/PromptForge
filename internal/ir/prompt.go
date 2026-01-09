package ir

// PromptIR represents the machine-enforceable prompt contract.
// This is the authoritative artifact produced by compilation.
type PromptIR struct {
	// SystemRole defines the role and context for the LLM.
	SystemRole string `json:"system_role"`

	// Rules contains a list of strict behavioral constraints.
	Rules []Rule `json:"rules"`

	// InputSchema defines the expected input structure.
	InputSchema Schema `json:"input_schema"`

	// OutputSchema defines the expected output structure.
	OutputSchema Schema `json:"output_schema"`

	// FailureModes describes how the system should handle errors.
	FailureModes []FailureMode `json:"failure_modes"`
}

// Rule represents a single behavioral constraint.
type Rule struct {
	// ID is a unique identifier for the rule.
	ID string `json:"id"`

	// Description is a human-readable description of the rule.
	Description string `json:"description"`

	// Condition is an optional condition that must be met for the rule to apply.
	Condition string `json:"condition,omitempty"`
}

// Schema defines the structure of input or output data.
type Schema struct {
	// Type is the root type (e.g., "object", "array", "string").
	Type string `json:"type"`

	// Properties defines fields for object types.
	Properties map[string]Property `json:"properties,omitempty"`

	// Required lists required property names for object types.
	Required []string `json:"required,omitempty"`

	// Items defines the schema for array element types.
	Items *Schema `json:"items,omitempty"`
}

// Property defines a single field in a schema.
type Property struct {
	// Type is the property type (e.g., "string", "number", "boolean", "object", "array").
	Type string `json:"type"`

	// Description is a human-readable description of the property.
	Description string `json:"description,omitempty"`

	// Enum lists allowed values for this property.
	Enum []interface{} `json:"enum,omitempty"`

	// Properties defines nested properties for object types.
	Properties map[string]Property `json:"properties,omitempty"`

	// Items defines the schema for array element types.
	Items *Schema `json:"items,omitempty"`
}

// FailureMode describes how to handle a specific failure scenario.
type FailureMode struct {
	// ID is a unique identifier for the failure mode.
	ID string `json:"id"`

	// Condition describes when this failure mode applies.
	Condition string `json:"condition"`

	// Response describes how the system should respond to this failure.
	Response string `json:"response"`
}


