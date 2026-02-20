package io

import (
	"os"
	"path/filepath"
	"testing"
)

const testInputJSON = `{
  "demo": {
    "parameters": {},
    "data": {}
  }
}`

func TestReadInputFileFallsBackToInputsJSON(t *testing.T) {
	dir := t.TempDir()
	fallbackPath := filepath.Join(dir, "inputs.json")

	if err := os.WriteFile(fallbackPath, []byte(testInputJSON), 0o644); err != nil {
		t.Fatalf("failed to write fallback inputs.json: %v", err)
	}

	_, err := ReadInputFile(filepath.Join(dir, "input.json"))
	if err != nil {
		t.Fatalf("expected fallback from input.json to inputs.json, got error: %v", err)
	}
}

func TestReadInputFilePrefersInputJSONWhenPresent(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	fallbackPath := filepath.Join(dir, "inputs.json")

	if err := os.WriteFile(inputPath, []byte(testInputJSON), 0o644); err != nil {
		t.Fatalf("failed to write input.json: %v", err)
	}
	if err := os.WriteFile(fallbackPath, []byte("not-json"), 0o644); err != nil {
		t.Fatalf("failed to write inputs.json: %v", err)
	}

	_, err := ReadInputFile(inputPath)
	if err != nil {
		t.Fatalf("expected input.json to be used first, got error: %v", err)
	}
}
