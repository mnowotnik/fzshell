package functional_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	fzshellBinary string
	repoRoot      string
)

func TestMain(m *testing.M) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get working directory: %v\n", err)
		os.Exit(1)
	}

	repoRoot = filepath.Dir(filepath.Dir(cwd))

	tempDir, err := os.MkdirTemp("", "fzshell-build-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp dir: %v\n", err)
		os.Exit(1)
	}

	binaryPath := filepath.Join(tempDir, "fzshell")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/fzshell")
	buildCmd.Dir = repoRoot
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	buildCmd.Env = os.Environ()

	if err := buildCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "go build failed: %v\n", err)
		os.Exit(1)
	}

	fzshellBinary = binaryPath

	exitCode := m.Run()

	os.RemoveAll(tempDir)

	os.Exit(exitCode)
}

type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func runFzshell(t *testing.T, args []string, stdin string, env map[string]string) CommandResult {
	t.Helper()

	if fzshellBinary == "" {
		t.Fatal("fzshell binary not built")
	}

	if repoRoot == "" {
		t.Fatal("repository root not configured")
	}

	cmd := exec.Command(fzshellBinary, args...)
	cmd.Dir = repoRoot

	runtimeDir := t.TempDir()

	envVars := environmentMap(os.Environ())

	defaults := map[string]string{
		"FZSHELL_DEBUG":   "1",
		"XDG_RUNTIME_DIR": runtimeDir,
	}

	for key, value := range defaults {
		if env != nil {
			if _, ok := env[key]; ok {
				continue
			}
		}
		envVars[key] = value
	}

	for key, value := range env {
		envVars[key] = value
	}

	cmd.Env = mapToEnvList(envVars)

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run fzshell: %v", err)
		}
	}

	return CommandResult{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: exitCode,
	}
}

func environmentMap(env []string) map[string]string {
	result := make(map[string]string, len(env))
	for _, entry := range env {
		if entry == "" {
			continue
		}

		parts := strings.SplitN(entry, "=", 2)
		key := parts[0]
		value := ""
		if len(parts) == 2 {
			value = parts[1]
		}
		result[key] = value
	}
	return result
}

func mapToEnvList(env map[string]string) []string {
	result := make([]string, 0, len(env))
	for key, value := range env {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return result
}
