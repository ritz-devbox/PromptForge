package ir

import "encoding/json"

const PromptIRSchemaID = "https://promptforge.dev/schemas/prompt-ir.schema.json"

// PromptIRSchemaJSON returns the JSON Schema for PromptIR, formatted for writing to disk.
func PromptIRSchemaJSON() ([]byte, error) {
	return json.MarshalIndent(promptIRSchema(), "", "  ")
}

func promptIRSchema() map[string]interface{} {
	return map[string]interface{}{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$id":     PromptIRSchemaID,
		"title":   "PromptForge Prompt IR",
		"type":    "object",
		"required": []string{
			"version",
			"system_role",
			"rules",
			"input_schema",
			"output_schema",
			"failure_modes",
		},
		"additionalProperties": false,
		"properties": map[string]interface{}{
			"version": map[string]interface{}{
				"type":      "string",
				"minLength": 1,
			},
			"system_role": map[string]interface{}{
				"type":      "string",
				"minLength": 1,
			},
			"rules": map[string]interface{}{
				"type":     "array",
				"minItems": 1,
				"items": map[string]interface{}{
					"$ref": "#/$defs/rule",
				},
			},
			"input_schema": map[string]interface{}{
				"$ref": "#/$defs/schema",
			},
			"output_schema": map[string]interface{}{
				"$ref": "#/$defs/schema",
			},
			"failure_modes": map[string]interface{}{
				"type":     "array",
				"minItems": 1,
				"items": map[string]interface{}{
					"$ref": "#/$defs/failure_mode",
				},
			},
		},
		"$defs": map[string]interface{}{
			"rule": map[string]interface{}{
				"type": "object",
				"required": []string{
					"id",
					"description",
				},
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":      "string",
						"minLength": 1,
					},
					"description": map[string]interface{}{
						"type":      "string",
						"minLength": 1,
					},
					"condition": map[string]interface{}{
						"type": "string",
					},
				},
			},
			"schema": map[string]interface{}{
				"type": "object",
				"required": []string{
					"type",
				},
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"type": map[string]interface{}{
						"type":      "string",
						"minLength": 1,
					},
					"properties": map[string]interface{}{
						"type": "object",
						"additionalProperties": map[string]interface{}{
							"$ref": "#/$defs/property",
						},
					},
					"required": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"items": map[string]interface{}{
						"$ref": "#/$defs/schema",
					},
				},
			},
			"property": map[string]interface{}{
				"type": "object",
				"required": []string{
					"type",
				},
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"type": map[string]interface{}{
						"type":      "string",
						"minLength": 1,
					},
					"description": map[string]interface{}{
						"type": "string",
					},
					"enum": map[string]interface{}{
						"type": "array",
					},
					"properties": map[string]interface{}{
						"type": "object",
						"additionalProperties": map[string]interface{}{
							"$ref": "#/$defs/property",
						},
					},
					"items": map[string]interface{}{
						"$ref": "#/$defs/schema",
					},
				},
			},
			"failure_mode": map[string]interface{}{
				"type": "object",
				"required": []string{
					"id",
					"condition",
					"response",
				},
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":      "string",
						"minLength": 1,
					},
					"condition": map[string]interface{}{
						"type":      "string",
						"minLength": 1,
					},
					"response": map[string]interface{}{
						"type":      "string",
						"minLength": 1,
					},
				},
			},
		},
	}
}
