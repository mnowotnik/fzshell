package functional_test

import (
	"path/filepath"
	"testing"
)

func TestReturnAllOutputsAllItems(t *testing.T) {
	configPath := filepath.Join(repoRoot, "tests", "functional", "fixtures", "simple.yaml")

	result := runFzshell(t, []string{"--config", configPath, "--all", "simple"}, "", nil)

	if result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr: %q)", result.ExitCode, result.Stderr)
	}

	const expectedOutput = "alpha\nbeta\n"
	if result.Stdout != expectedOutput {
		t.Fatalf("unexpected stdout: got %q, want %q", result.Stdout, expectedOutput)
	}

	if result.Stderr != "" {
		t.Fatalf("expected empty stderr, got %q", result.Stderr)
	}
}

func TestReplacementRendersLineBuffer(t *testing.T) {
	configPath := filepath.Join(repoRoot, "tests", "functional", "fixtures", "replacement.yaml")

	result := runFzshell(t, []string{"--config", configPath, "ticket 42/build"}, "", nil)

	if result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr: %q)", result.ExitCode, result.Stderr)
	}

	const expectedOutput = "ticket 42/build -> 42|build|command:42:build"
	if result.Stdout != expectedOutput {
		t.Fatalf("unexpected stdout: got %q, want %q", result.Stdout, expectedOutput)
	}

	if result.Stderr != "" {
		t.Fatalf("expected empty stderr, got %q", result.Stderr)
	}
}
