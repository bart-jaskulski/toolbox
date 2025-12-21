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
	fmt.Printf("File processing complete. Added %d files to %s.\n", fileCount, formatDisplayName(cfg.OutputFormat))
	if cfg.Metadata != nil {
		fmt.Println("Project metadata was included.") // Add a note if metadata was added
	}
	fmt.Printf("Output saved to %s\n", cfg.OutputFile)
	if !cfg.IncludeBinary {
		fmt.Println("(Likely binary files were skipped unless --include-binary was used)")
	}
	fmt.Println("-------------------------------------------------")
}
