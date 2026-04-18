package exec

import (
	"fmt"
	"os"
	os_exec "os/exec"
	"strings"
)

// Runner holds configuration for executing a subprocess with injected secrets.
type Runner struct {
	Env    map[string]string
	Stdout *os.File
	Stderr *os.File
}

// NewRunner creates a Runner with the provided secret map.
func NewRunner(secrets map[string]string) *Runner {
	return &Runner{
		Env:    secrets,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run executes the given command with the merged environment.
// Secrets are injected as environment variables without writing to disk.
func (r *Runner) Run(command string, args []string) error {
	if command == "" {
		return fmt.Errorf("exec: command must not be empty")
	}

	cmd := os_exec.Command(command, args...)
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr

	// Start with the current process environment.
	baseEnv := os.Environ()

	// Append secrets, overriding any existing keys.
	overrideKeys := make(map[string]bool, len(r.Env))
	for k := range r.Env {
		overrideKeys[strings.ToUpper(k)] = true
	}

	filtered := make([]string, 0, len(baseEnv))
	for _, e := range baseEnv {
		parts := strings.SplitN(e, "=", 2)
		if !overrideKeys[strings.ToUpper(parts[0])] {
			filtered = append(filtered, e)
		}
	}

	for k, v := range r.Env {
		filtered = append(filtered, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Env = filtered

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("exec: command failed: %w", err)
	}
	return nil
}
