package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func runRootForTest(t *testing.T, args ...string) (string, error) {
	t.Helper()

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	return out.String(), err
}

func TestRootVersionLongFlag(t *testing.T) {
	out, err := runRootForTest(t, "--version")
	if err != nil {
		t.Fatalf("expected no error for --version, got %v", err)
	}
	if !strings.Contains(out, rootCmd.Version) {
		t.Fatalf("expected output to contain version %q, got %q", rootCmd.Version, out)
	}
}

func TestRootVersionShortFlag(t *testing.T) {
	out, err := runRootForTest(t, "-v")
	if err != nil {
		t.Fatalf("expected no error for -v, got %v", err)
	}
	if !strings.Contains(out, rootCmd.Version) {
		t.Fatalf("expected output to contain version %q, got %q", rootCmd.Version, out)
	}
}
