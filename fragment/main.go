// main.go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"slices" // Go 1.21+ needed for slices

	"github.com/urfave/cli/v3"
)

// Config holds the application's configuration and resolved state.
type Config struct {
	ProjectRoot        string
	AbsoluteOutputFile string
	AbsoluteInputDirs  []string
	IsGitRepo          bool
	RealpathCmd        string

	// From CLI flags/args
	OutputFile      string
	OutputFormat    string
	InputDirs       []string
	ExcludePatterns []string // User-provided patterns
	IncludeBinary   bool
	Verbose         bool

	// Added fields
	Metadata          *ProjectMetadata // Holds non-package metadata (if any in future)
	ExtractedPackages []XmlPackage     // Holds unified list of extracted package information
	DefaultExcludes   []string         // Default patterns to exclude
	CombinedExcludes  []string         // User + Default patterns
}

// List of common lock files to exclude by default.
var defaultExcludePatterns = []string{
	"composer.lock",
	"package-lock.json",
	"yarn.lock",
	"go.sum",
	"Gemfile.lock",
	"Pipfile.lock",
	"poetry.lock",
	// Add more if needed
}

func main() {
	// Configure logging - remove timestamps initially
	log.SetFlags(0)

	cmd := &cli.Command{
		Name:      "concat",
		Version:   "1.4.0", // Version updated
		Usage:     "Concatenates files and project metadata into XML or Markdown",
		UsageText: "concat [command options] [directory...]",
		Flags: []cli.Flag{
			// ... (other flags remain the same) ...
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "output.xml",
				Usage:   "Specify output file path",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "xml",
				Usage:   "Output format: xml or markdown",
			},
			&cli.StringSliceFlag{
				Name:    "exclude",
				Aliases: []string{"x"},
				Usage:   "Add a glob `PATTERN` to exclude (can be repeated)",
			},
			&cli.BoolFlag{
				Name:  "include-binary",
				Value: false,
				Usage: "Include all files (don't skip likely binary files)",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Value:   false,
				Usage:   "Enable verbose logging output",
			},
		},
		Action: runConcatenation,
		Description: `
Concatenates files from specified directories (default: '.') into a single XML or Markdown file.
Includes project metadata extracted from composer.json or package.json if found at the root.
Project root is Git repo root (if applicable) or PWD. Paths relative to root.

Respects .gitignore (if Git). Applies command-line exclusions. Excludes output file.
Common lock files (composer.lock, package-lock.json, etc.) are excluded by default.
Binary files are detected heuristically and skipped unless --include-binary is used.

Use -v or --verbose to see detailed processing steps.

Exclusion Patterns (-x): Uses Bash-style globbing relative to project root.
  '*.ext' matches suffixes recursively.
  'dir/' matches directory prefix recursively.
  'name' matches file 'name' or directory 'name/' recursively.
  Other patterns use standard globbing (e.g., 'src/*.js').
        `,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.SetOutput(os.Stderr)
		log.Fatalf("Error: %v", err)
	}
}

// runConcatenation is the main logic, executed by cli.Command.Action
func runConcatenation(ctx context.Context, cmd *cli.Command) error {
	outputFormat, err := normalizeOutputFormat(cmd.String("format"))
	if err != nil {
		return err
	}

	// --- Create Config from flags ---
	cfg := Config{
		OutputFile:      cmd.String("output"),
		OutputFormat:    outputFormat,
		ExcludePatterns: cmd.StringSlice("exclude"), // User patterns
		IncludeBinary:   cmd.Bool("include-binary"),
		Verbose:         cmd.Bool("verbose"),
		InputDirs:       cmd.Args().Slice(),
		DefaultExcludes: defaultExcludePatterns, // Assign defaults
	}

	// --- Configure Logging ---
	if !cfg.Verbose {
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(os.Stderr)
	}

	if len(cfg.InputDirs) == 0 {
		log.Println("Info: No input directories specified, defaulting to '.'")
		cfg.InputDirs = []string{"."}
	}

	if !cmd.IsSet("output") {
		cfg.OutputFile = defaultOutputFile(cfg.OutputFormat)
	}

	// --- Preparations ---
	cfg.RealpathCmd = findRealpathCmd()

	cfg.ProjectRoot, cfg.IsGitRepo, err = determineProjectRoot()
	if err != nil {
		return fmt.Errorf("determining project root: %w", err)
	}

	// --- Extract Metadata ---
	// Happens after project root is known
	// cfg.Metadata would be for non-package metadata if any such extractors are added.
	// For now, it will remain nil.
	cfg.ExtractedPackages = extractPackageData(&cfg) // Store extracted package info

	// --- Combine Exclusions ---
	// Start with defaults, then add user patterns, avoiding duplicates.
	cfg.CombinedExcludes = slices.Clone(cfg.DefaultExcludes)
	seenExcludes := make(map[string]bool, len(cfg.DefaultExcludes))
	for _, p := range cfg.DefaultExcludes {
		seenExcludes[p] = true
	}
	for _, p := range cfg.ExcludePatterns { // Iterate user patterns
		if !seenExcludes[p] {
			cfg.CombinedExcludes = append(cfg.CombinedExcludes, p)
			seenExcludes[p] = true
		}
	}
	// Log combined exclusions if verbose
	if cfg.Verbose {
		log.Println("Effective exclusion patterns:")
		for _, p := range cfg.CombinedExcludes {
			isDefault := slices.Contains(cfg.DefaultExcludes, p) && !slices.Contains(cfg.ExcludePatterns, p)
			tag := ""
			if isDefault {
				tag = " (default)"
			}
			log.Printf("  - %s%s", p, tag)
		}
	}

	// --- Resolve Paths ---
	cfg.AbsoluteOutputFile, cfg.AbsoluteInputDirs, err = resolveAndValidatePaths(
		cfg.OutputFile,
		cfg.InputDirs,
		cfg.RealpathCmd,
	)
	if err != nil {
		return fmt.Errorf("resolving paths: %w", err)
	}

	printProcessingInfo(&cfg) // Pass config

	// --- Gather Files ---
	filesToInclude := make(map[string]string)
	// Pass cfg.CombinedExcludes to gathering functions
	if cfg.IsGitRepo {
		err = gatherFilesGit(&cfg, filesToInclude)
	} else {
		err = gatherFilesWalk(&cfg, filesToInclude)
	}
	if err != nil {
		return fmt.Errorf("gathering files: %w", err)
	}

	// --- Generate Output ---
	var fileCount int
	switch cfg.OutputFormat {
	case outputFormatXML:
		fileCount, err = generateXML(&cfg, filesToInclude) // Pass config (contains metadata)
	case outputFormatMarkdown:
		fileCount, err = generateMarkdown(&cfg, filesToInclude)
	default:
		return fmt.Errorf("unsupported output format: %s", cfg.OutputFormat)
	}
	if err != nil {
		return fmt.Errorf("generating %s: %w", cfg.OutputFormat, err)
	}

	// --- Print Final Summary ---
	printCompletionInfo(&cfg, fileCount)
	return nil
}
