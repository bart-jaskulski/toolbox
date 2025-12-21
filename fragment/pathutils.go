package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// findRealpathCmd checks for realpath or grealpath and returns the command name.
func findRealpathCmd() string {
	if _, err := exec.LookPath("realpath"); err == nil {
		return "realpath"
	}
	if _, err := exec.LookPath("grealpath"); err == nil {
		log.Println("Info: Using 'grealpath' (found Homebrew coreutils on macOS?).")
		return "grealpath" // Found Homebrew coreutils on macOS
	}
	log.Println("Warning: 'realpath' command not found (part of coreutils). Path handling may be less robust (using filepath.Abs).")
	return "" // Not found
}

// determineProjectRoot finds the project root (Git or PWD) and returns it and whether it's a Git repo.
func determineProjectRoot() (projectRoot string, isGitRepo bool, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", false, fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Default to PWD
	projectRoot = pwd
	isGitRepo = false

	// Check for Git
	gitPath, gitErr := exec.LookPath("git")
	if gitErr != nil {
		log.Println("Info: 'git' command not found. Cannot check for Git repository or use .gitignore.")
	} else {
		// Check if inside a work tree
		_, errGitCheck := runCommand(pwd, gitPath, "rev-parse", "--is-inside-work-tree")
		isGitRepo = (errGitCheck == nil)

		if isGitRepo {
			gitRoot, errGitRoot := runCommand(pwd, gitPath, "rev-parse", "--show-toplevel")
			if errGitRoot != nil {
				isGitRepo = false // Correct state on error
				log.Printf("Warning: Could not determine Git root despite being inside a work tree (%v). Using PWD.", errGitRoot)
				// projectRoot remains pwd
			} else {
				projectRoot = gitRoot // Use Git root
				log.Printf("Detected Git repository. Project root: %s", projectRoot)
			}
		} else {
			log.Printf("Not inside a Git repository. Using current directory as project root: %s", projectRoot)
		}
	}

	// Ensure projectRoot is absolute and clean
	projectRoot, err = filepath.Abs(projectRoot)
	if err != nil {
		return "", isGitRepo, fmt.Errorf("failed to get absolute path for project root '%s': %w", projectRoot, err)
	}
	projectRoot = filepath.Clean(projectRoot)
	return projectRoot, isGitRepo, nil
}

// resolvePath resolves a given path to an absolute path, preferring realpath if available.
// Accepts the realpath command name (or "" if unavailable) as an argument.
func resolvePath(path string, realpathCmd string) (string, error) {
	if realpathCmd != "" {
		absPath, err := runCommand("", realpathCmd, "-m", path) // Use the provided command
		if err != nil {
			log.Printf("Warning: %s failed for '%s': %v. Falling back to filepath.Abs.", realpathCmd, path, err)
			// Fallback if realpath fails
		} else {
			return filepath.Clean(absPath), nil
		}
	}
	// Fallback if realpath command is not available or failed
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("filepath.Abs failed for '%s': %w", path, err)
	}
	return filepath.Clean(abs), nil
}

// resolveAndValidatePaths resolves output and input paths, validates inputs.
// Returns absolute output path, slice of absolute input paths, and error.
func resolveAndValidatePaths(outputFile string, inputDirs []string, realpathCmd string) (string, []string, error) {
	var err error
	// Resolve output file path
	absoluteOutputFile, err := resolvePath(outputFile, realpathCmd)
	if err != nil {
		return "", nil, fmt.Errorf("failed to resolve output file path '%s': %w", outputFile, err)
	}
	log.Printf("Output file target (absolute): %s", absoluteOutputFile)

	// Resolve and validate input directory paths
	absoluteInputDirs := make([]string, 0, len(inputDirs))
	hasInvalidInput := false
	seenInputDirs := make(map[string]bool) // Avoid duplicates

	for _, dir := range inputDirs {
		absDir, err := resolvePath(dir, realpathCmd)
		if err != nil {
			log.Printf("Error: Could not resolve path for input directory '%s': %v", dir, err)
			hasInvalidInput = true
			continue
		}

		if seenInputDirs[absDir] {
			log.Printf("Info: Skipping duplicate input directory '%s' (resolved to '%s')", dir, absDir)
			continue
		}

		info, err := os.Stat(absDir)
		if err != nil {
			log.Printf("Error: Cannot access input directory '%s' (resolved to '%s'): %v", dir, absDir, err)
			hasInvalidInput = true
			continue
		}
		if !info.IsDir() {
			log.Printf("Error: Input path '%s' (resolved to '%s') is not a directory.", dir, absDir)
			hasInvalidInput = true
			continue
		}
		absoluteInputDirs = append(absoluteInputDirs, absDir)
		seenInputDirs[absDir] = true
	}

	if hasInvalidInput {
		return "", nil, errors.New("aborting due to invalid input directories")
	}
	if len(absoluteInputDirs) == 0 {
		return "", nil, errors.New("no valid input directories remaining after validation")
	}
	return absoluteOutputFile, absoluteInputDirs, nil
}
