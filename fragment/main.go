// main.go
package main

import (
	"bufio"
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
	ExtractedPackages []ProjectPackage // Holds unified list of extracted package information
	DefaultExcludes   []string         // Default patterns to exclude
	CombinedExcludes  []string         // User + Default patterns
	NoIgnore          bool             // Do not respect .gitignore (Git mode only)
	NoTree            bool             // Do not generate directory tree overview
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
		Usage:     "Concatenates files and project metadata into XML, Markdown, or JSON",
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
				Usage:   "Output format: xml, markdown, or json",
			},
			&cli.StringSliceFlag{
				Name:    "exclude",
				Aliases: []string{"x"},
				Usage:   "Add a glob `PATTERN` to exclude (can be repeated)",
			},
			&cli.BoolFlag{
				Name:  "no-ignore",
				Value: false,
				Usage: "Do not respect .gitignore when scanning (Git mode only)",
			},
			&cli.BoolFlag{
				Name:  "no-tree",
				Value: false,
				Usage: "Do not generate directory tree overview",
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
Concatenates files from specified directories (default: '.') into a single XML, Markdown, or JSON file.
Includes project metadata extracted from composer.json or package.json if found at the root.
Project root is Git repo root (if applicable) or PWD. Paths relative to root.

Respects .gitignore (if Git). Applies command-line exclusions. Excludes output file.
Common lock files (composer.lock, package-lock.json, etc.) are excluded by default.
Binary files are detected heuristically and skipped unless --include-binary is used.
Use --no-ignore to bypass .gitignore when scanning a Git repository.
If --format is omitted, the format is inferred from the output file extension when possible.
Use --no-tree to skip generating the directory tree overview.

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
	rawFormat := cmd.String("format")
	outputFile := cmd.String("output")
	formatSet := cmd.IsSet("format")
	outputSet := cmd.IsSet("output")

	outputFormat := outputFormatXML
	var err error
	if formatSet {
		normalized, err := normalizeOutputFormat(rawFormat)
		if err != nil {
			return err
		}
		outputFormat = normalized
	} else if outputSet {
		if inferred, ok := inferFormatFromOutput(outputFile); ok {
			outputFormat = inferred
		}
	}

	// --- Create Config from flags ---
	cfg := Config{
		OutputFile:      outputFile,
		OutputFormat:    outputFormat,
		ExcludePatterns: cmd.StringSlice("exclude"), // User patterns
		NoIgnore:        cmd.Bool("no-ignore"),
		NoTree:          cmd.Bool("no-tree"),
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

	if !outputSet {
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
	useGit := cfg.IsGitRepo && inputsUnderRoot(cfg.ProjectRoot, cfg.AbsoluteInputDirs)
	if cfg.IsGitRepo && !useGit {
		log.Println("Info: Input directories are outside the Git root; using directory walk instead of git ls-files.")
	}
	if useGit {
		err = gatherFilesGit(&cfg, filesToInclude)
	} else {
		err = gatherFilesWalk(&cfg, filesToInclude)
	}
	if err != nil {
		return fmt.Errorf("gathering files: %w", err)
	}

	// --- Build Snapshot ---
	snapshot, fileCount, err := buildSnapshot(&cfg, filesToInclude)
	if err != nil {
		return fmt.Errorf("building snapshot: %w", err)
	}

	// --- Generate Output ---
	formatter, err := getFormatter(cfg.OutputFormat)
	if err != nil {
		return err
	}

	outFile, err := os.Create(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", cfg.OutputFile, err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	if err := formatter.Write(writer, snapshot); err != nil {
		return fmt.Errorf("generating %s: %w", cfg.OutputFormat, err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flushing output: %w", err)
	}

	// --- Print Final Summary ---
	printCompletionInfo(&cfg, fileCount)
	return nil
}
