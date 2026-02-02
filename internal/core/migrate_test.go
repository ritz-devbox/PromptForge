package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateIR_AddsVersion(t *testing.T) {
	tmpDir := t.TempDir()

	irPayload := map[string]interface{}{
		"system_role": "Test role",
		"rules": []map[string]interface{}{
			{"id": "rule-1", "description": "Test rule"},
		},
		"input_schema": map[string]interface{}{
			"type": "object",
		},
		"output_schema": map[string]interface{}{
			"type": "object",
		},
		"failure_modes": []map[string]interface{}{
			{"id": "fm-1", "condition": "Test", "response": "Test"},
		},
	}

	data, err := json.MarshalIndent(irPayload, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal IR payload: %v", err)
	}

	irPath := filepath.Join(tmpDir, "prompt.ir.json")
	if err := os.WriteFile(irPath, data, 0644); err != nil {
		t.Fatalf("Failed to write prompt.ir.json: %v", err)
	}

	if err := MigrateIR(tmpDir); err != nil {
		t.Fatalf("MigrateIR() failed: %v", err)
	}

	updated, err := os.ReadFile(irPath)
	if err != nil {
		t.Fatalf("Failed to read migrated prompt.ir.json: %v", err)
	}
	if !containsVersion(updated) {
		t.Fatal("Expected migrated IR to include version")
	}
}

func TestMigrateIR_UnsupportedVersion(t *testing.T) {
	tmpDir := t.TempDir()

	irPayload := map[string]interface{}{
		"version":     "9.9",
		"system_role": "Test role",
		"rules": []map[string]interface{}{
			{"id": "rule-1", "description": "Test rule"},
		},
		"input_schema": map[string]interface{}{
			"type": "object",
		},
		"output_schema": map[string]interface{}{
			"type": "object",
		},
		"failure_modes": []map[string]interface{}{
			{"id": "fm-1", "condition": "Test", "response": "Test"},
		},
	}

	data, err := json.MarshalIndent(irPayload, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal IR payload: %v", err)
	}

	irPath := filepath.Join(tmpDir, "prompt.ir.json")
	if err := os.WriteFile(irPath, data, 0644); err != nil {
		t.Fatalf("Failed to write prompt.ir.json: %v", err)
	}

	if err := MigrateIR(tmpDir); err == nil {
		t.Fatal("Expected MigrateIR() to fail for unsupported version")
	}
}

func containsVersion(data []byte) bool {
	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return false
	}
	_, ok := payload["version"]
	return ok
}
