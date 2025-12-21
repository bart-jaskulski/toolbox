// files.go
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// gatherFilesGit uses 'git ls-files' to find candidate files.
func gatherFilesGit(cfg *Config, filesToInclude map[string]string) error {
	if cfg.NoIgnore {
		log.Println("Scanning Git repository for files (using git ls-files, ignoring .gitignore)...")
	} else {
		log.Println("Scanning Git repository for files (using git ls-files, respecting .gitignore)...")
	}
	args := []string{"ls-files", "-c", "-o", "-z"}
	if !cfg.NoIgnore {
		args = append(args, "--exclude-standard")
	}
	cmdOutput, err := runCommand(cfg.ProjectRoot, "git", args...)
	// ... (error handling as before) ...
	if err != nil {
		// ... handle git errors or fallback ...
		if strings.Contains(err.Error(), "not a git repository") {
			log.Println("Warning: 'git ls-files' failed (not a git repository?). Falling back to directory walk.")
			return gatherFilesWalk(cfg, filesToInclude)
		}
		return fmt.Errorf("git ls-files command failed: %w", err)
	}
	if cmdOutput == "" {
		log.Println("Info: 'git ls-files' returned no files.")
		return nil
	}

	gitFiles := strings.Split(cmdOutput, "\x00")

	for _, relPath := range gitFiles {
		// ... (path joining and cleaning as before) ...
		if relPath == "" {
			continue
		}
		absPath := filepath.Join(cfg.ProjectRoot, relPath)
		absPath = filepath.Clean(absPath)

		if absPath == cfg.AbsoluteOutputFile {
			continue
		} // Skip output

		// ... (check if in input dirs as before) ...
		isInInput := false
		for _, inputRoot := range cfg.AbsoluteInputDirs {
			if absPath == inputRoot || strings.HasPrefix(absPath, inputRoot+string(filepath.Separator)) {
				info, errStat := os.Stat(absPath)
				if errStat == nil && !info.IsDir() {
					isInInput = true
					break
				}
			}
		}
		if !isInInput {
			continue
		}

		// *** Use CombinedExcludes ***
		if checkExclusion(relPath, cfg.CombinedExcludes) {
			continue
		}

		// ... (check binary as before) ...
		if !cfg.IncludeBinary {
			isBin, err := isLikelyBinary(absPath)
			if err != nil {
				log.Printf("Warning: Could not check file type for '%s': %v. Skipping.", absPath, err)
				continue
			}
			if isBin {
				log.Printf("    Skipping likely binary file: %s", relPath)
				continue
			}
		}

		filesToInclude[relPath] = absPath
	}
	return nil
}

// gatherFilesWalk uses filepath.WalkDir for non-Git scenarios.
func gatherFilesWalk(cfg *Config, filesToInclude map[string]string) error {
	log.Println("Scanning specified directories using directory walk...")
	processedAbsPaths := make(map[string]bool)

	for _, inputRoot := range cfg.AbsoluteInputDirs {
		log.Printf("  Scanning under: %s", inputRoot)
		err := filepath.WalkDir(inputRoot, func(absPath string, d fs.DirEntry, walkErr error) error {
			// ... (error handling for walkErr as before) ...
			if walkErr != nil {
				log.Printf("Warning: Error accessing path %q during walk: %v", absPath, walkErr)
				if errors.Is(walkErr, fs.ErrPermission) {
					return nil
				}
				return nil // Skip item but continue walk
			}

			// Calculate relative path early for directory exclusion check
			relPath, errRel := filepath.Rel(cfg.ProjectRoot, absPath)
			if errRel != nil {
				// Use absolute path for logging if relative fails, but might affect exclusion matching
				log.Printf("Warning: Could not determine relative path for %s from %s: %v. Exclusion checks might be affected.", absPath, cfg.ProjectRoot, errRel)
				relPath = absPath // Fallback, less ideal for matching
			}

			if d.IsDir() {
				// Check if directory itself is excluded using CombinedExcludes
				// Add trailing slash for directory-specific patterns like 'node_modules/'
				dirRelPathForCheck := relPath
				if !strings.HasSuffix(dirRelPathForCheck, "/") {
					dirRelPathForCheck += "/"
				}
				// Check both 'dir/' and 'dir' patterns
				// *** Use CombinedExcludes ***
				if checkExclusion(dirRelPathForCheck, cfg.CombinedExcludes) || checkExclusion(relPath, cfg.CombinedExcludes) {
					log.Printf("    Skipping excluded directory: %s", relPath)
					return filepath.SkipDir
				}
				return nil // Continue walking into directory
			}

			// Process Files
			absPath = filepath.Clean(absPath)
			if absPath == cfg.AbsoluteOutputFile {
				return nil
			} // Skip output
			if processedAbsPaths[absPath] {
				return nil
			} // De-duplicate
			processedAbsPaths[absPath] = true

			// We already calculated relPath above

			// *** Use CombinedExcludes ***
			if checkExclusion(relPath, cfg.CombinedExcludes) {
				return nil
			}

			// ... (check binary as before) ...
			if !cfg.IncludeBinary {
				isBin, errCheck := isLikelyBinary(absPath)
				if errCheck != nil {
					log.Printf("Warning: Could not check file type for '%s': %v. Skipping.", absPath, errCheck)
					return nil
				}
				if isBin {
					log.Printf("    Skipping likely binary file: %s", relPath)
					return nil
				}
			}

			filesToInclude[relPath] = absPath
			return nil
		}) // End WalkDir func

		if err != nil {
			log.Printf("Error during directory walk for '%s': %v", inputRoot, err)
		}
	} // End loop over input directories
	return nil
}

// checkExclusion checks if a path matches any pattern in the provided list.
// (No changes needed in this function itself, it already takes the list as argument)
func checkExclusion(relPath string, excludePatterns []string) bool {
	// ... implementation remains the same ...
	isExcluded := false
	var matchingPattern string
	relPath = filepath.ToSlash(relPath)

	for _, pattern := range excludePatterns {
		pattern = filepath.ToSlash(pattern)
		isSimpleName := !strings.ContainsAny(pattern, "*/?[]")

		if strings.HasPrefix(pattern, "*.") && !strings.Contains(pattern, "/") {
			suffix := pattern[1:]
			if strings.HasSuffix(relPath, suffix) {
				isExcluded = true
				matchingPattern = pattern
				break
			}
		} else if strings.HasSuffix(pattern, "/") {
			if strings.HasPrefix(relPath, pattern) {
				isExcluded = true
				matchingPattern = pattern
				break
			}
		} else if isSimpleName {
			// Match exact filename or if path starts with 'name/'
			// Use filepath.Base for exact filename match robustness
			if filepath.Base(relPath) == pattern || strings.HasPrefix(relPath, pattern+"/") {
				isExcluded = true
				matchingPattern = pattern
				break
			}
		} else {
			match, _ := filepath.Match(pattern, relPath)
			if match {
				isExcluded = true
				matchingPattern = pattern
				break
			}
			matchBase, _ := filepath.Match(pattern, filepath.Base(relPath))
			if matchBase {
				isExcluded = true
				matchingPattern = pattern
				break
			}
		}
	}
	if isExcluded {
		log.Printf("    Excluding file matching pattern '%s': %s", matchingPattern, relPath)
		return true
	}
	return false
}

// isLikelyBinary checks the first chunk of a file for null bytes.
// (No changes needed in this function itself)
func isLikelyBinary(absPath string) (bool, error) {
	// ... implementation remains the same ...
	file, err := os.Open(absPath)
	if err != nil {
		return false, fmt.Errorf("failed to open file %s: %w", absPath, err)
	}
	defer file.Close()
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("failed to read file %s: %w", absPath, err)
	}
	if n == 0 {
		return false, nil
	}
	if bytes.Contains(buffer[:n], []byte{0}) {
		return true, nil
	}
	return false, nil
}
