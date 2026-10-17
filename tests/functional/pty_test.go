package functional_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/creack/pty"
)

const fzshellPTYTimeout = 15 * time.Second

// runFzshellWithPTY gives fzshell a controlling terminal while keeping its
// stdout and stderr captured separately. The replacement flow uses fzf's
// select-one behavior, which initializes the terminal even for a single match.
func runFzshellWithPTY(t *testing.T, args []string, env map[string]string) CommandResult {
	t.Helper()

	if fzshellBinary == "" {
		t.Fatal("fzshell binary not built")
	}
	if repoRoot == "" {
		t.Fatal("repository root not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), fzshellPTYTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, fzshellBinary, args...)
	cmd.Dir = repoRoot

	envVars := environmentMap(os.Environ())
	defaults := map[string]string{
		"FZSHELL_DEBUG":   "1",
		"XDG_RUNTIME_DIR": t.TempDir(),
		"TERM":            "xterm-256color",
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

	tty, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 24, Cols: 80})
	if errors.Is(err, pty.ErrUnsupported) {
		return runFzshell(t, args, "", env)
	}
	if err != nil {
		t.Fatalf("failed to start fzshell with a pseudo-terminal: %v", err)
	}

	drained := make(chan struct{})
	go func() {
		defer close(drained)
		_, _ = io.Copy(io.Discard, tty)
	}()

	err = cmd.Wait()
	_ = tty.Close()
	<-drained

	if ctx.Err() != nil {
		t.Fatalf("fzshell timed out after %s (stderr: %q)", fzshellPTYTimeout, stderrBuf.String())
	}

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
