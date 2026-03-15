package bindings

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	toolspec "github.com/hydrocode-de/tool-spec-go"
)

func TestGeneratorsIncludeLoggingHelpers(t *testing.T) {
	spec := toolspec.ToolSpec{
		Name: "sample_tool",
		Parameters: map[string]toolspec.ParameterSpec{
			"foo": {ToolType: "string"},
		},
		Data: map[string]toolspec.DataSpec{
			"bar": {Description: "input dataset"},
		},
	}

	testCases := []struct {
		target    string
		output    string
		fragments []string
	}{
		{
			target: "python",
			output: "parameters.py",
			fragments: []string{
				"def get_run_context()",
				"def get_logger()",
				"class GotapLogger:",
				"GOTAP_LOG_FILE",
			},
		},
		{
			target: "r",
			output: "parameters.R",
			fragments: []string{
				"get_run_context <- function()",
				"get_logger <- function()",
				"GOTAP_LOG_FILE",
			},
		},
		{
			target: "node-js",
			output: "parameters.js",
			fragments: []string{
				"export function getRunContext()",
				"export function getLogger()",
				"GOTAP_LOG_FILE",
			},
		},
		{
			target: "node-ts",
			output: "parameters.ts",
			fragments: []string{
				"export function getRunContext(): RunContext",
				"export function getLogger(): GotapLogger",
				"GOTAP_LOG_FILE",
			},
		},
		{
			target: "matlab",
			output: "get_parameters.m",
			fragments: []string{
				"function context = get_run_context()",
				"function logger = get_logger()",
				"GOTAP_LOG_FILE",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.target, func(t *testing.T) {
			generator, ok := GetGenerator(tc.target)
			if !ok {
				t.Fatalf("generator %q not found", tc.target)
			}

			outputPath := filepath.Join(t.TempDir(), tc.output)
			if err := generator.Generate(spec, outputPath); err != nil {
				t.Fatalf("Generate returned error: %v", err)
			}

			contentBytes, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatalf("ReadFile returned error: %v", err)
			}
			content := string(contentBytes)

			for _, fragment := range tc.fragments {
				if !strings.Contains(content, fragment) {
					t.Fatalf("generated file missing fragment %q", fragment)
				}
			}
		})
	}
}
