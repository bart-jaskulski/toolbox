package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// runCommand executes an external command and returns its stdout or an error.
// (No changes needed in this function itself)
func runCommand(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir // Set working directory if provided
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Log the command being run if verbose? Could be added via config.
	// if cfg.Verbose { log.Printf("Running command: %s %s [in %s]", name, strings.Join(args, " "), dir) }

	err := cmd.Run()
	if err != nil {
		stderrStr := strings.TrimSpace(stderr.String())
		// Provide a more structured error message
		baseMsg := fmt.Sprintf("command '%s %s' failed in directory '%s': %v", name, strings.Join(args, " "), dir, err)
		if stderrStr != "" {
			// Return stderr along with the error message
			return "", fmt.Errorf("%s\nstderr: %s", baseMsg, stderrStr)
		}
		// Use errors.New or fmt.Errorf consistently
		return "", fmt.Errorf("%s", baseMsg)
	}
	return strings.TrimSpace(stdout.String()), nil
}
