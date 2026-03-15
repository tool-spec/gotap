package logging

import (
	"path/filepath"
	"testing"
)

func TestNewRunContext(t *testing.T) {
	outputFolder := filepath.Join("tmp", "out")

	ctx, err := NewRunContext(outputFolder)
	if err != nil {
		t.Fatalf("NewRunContext returned error: %v", err)
	}

	if ctx.RunID == "" {
		t.Fatal("expected non-empty run id")
	}
	if got, want := ctx.LogFile, filepath.Join(outputFolder, "_logs.jsonl"); got != want {
		t.Fatalf("log file mismatch: got %q want %q", got, want)
	}
	if ctx.LogFormat != LogFormatJSONL {
		t.Fatalf("log format mismatch: got %q", ctx.LogFormat)
	}
	if ctx.LogLevel != DefaultLevel {
		t.Fatalf("log level mismatch: got %q", ctx.LogLevel)
	}
}

func TestRunContextEnv(t *testing.T) {
	ctx := RunContext{
		RunID:     "abc123",
		LogFile:   "/tmp/out/_logs.jsonl",
		LogFormat: LogFormatJSONL,
		LogLevel:  DefaultLevel,
	}

	env := ctx.Env()
	if len(env) != 4 {
		t.Fatalf("unexpected env count: got %d", len(env))
	}

	expected := map[string]bool{
		EnvRunID + "=abc123":                 false,
		EnvLogFile + "=/tmp/out/_logs.jsonl": false,
		EnvLogFormat + "=" + LogFormatJSONL:  false,
		EnvLogLevel + "=" + DefaultLevel:     false,
	}

	for _, item := range env {
		if _, ok := expected[item]; ok {
			expected[item] = true
		}
	}

	for item, seen := range expected {
		if !seen {
			t.Fatalf("missing env item %q", item)
		}
	}
}
