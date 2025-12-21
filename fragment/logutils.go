// logutils.go
package main

import (
	"fmt"
	"log"
)

// printProcessingInfo logs the initial processing setup.
func printProcessingInfo(cfg *Config) {
	log.Println("Processing files from:")
	for _, absDir := range cfg.AbsoluteInputDirs {
		log.Printf("  - %s", absDir)
	}
	// Combined excludes are logged verbosely in runConcatenation
	// We can just mention that exclusions are active if any exist
	if len(cfg.CombinedExcludes) > 0 {
		log.Println("Applying exclusion patterns (see verbose output for details).")
	}

	if !cfg.IncludeBinary {
		log.Println("Heuristic binary file detection is enabled (skipping files with null bytes).")
	} else {
		log.Println("Including all files (--include-binary specified).")
	}
}

// printCompletionInfo prints the final summary to stdout.
// (No changes needed here)
func printCompletionInfo(cfg *Config, fileCount int) {
	fmt.Println("-------------------------------------------------")
	fmt.Printf("Wrote %d files to %s (%s)\n", fileCount, cfg.OutputFile, formatDisplayName(cfg.OutputFormat))
	fmt.Println("-------------------------------------------------")
}
