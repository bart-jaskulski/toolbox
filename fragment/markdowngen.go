package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func generateMarkdown(cfg *Config, filesToInclude map[string]string) (int, error) {
	log.Println("Generating Markdown output...")
	processedFilesCount := 0

	projectName := filepath.Base(cfg.ProjectRoot)

	relPaths := make([]string, 0, len(filesToInclude))
	for relPath := range filesToInclude {
		relPaths = append(relPaths, relPath)
	}
	sort.Strings(relPaths)

	outFile, err := os.Create(cfg.OutputFile)
	if err != nil {
		return 0, fmt.Errorf("failed to create output file '%s': %w", cfg.OutputFile, err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)

	_, err = writer.WriteString("# Project: " + projectName + "\n\n")
	if err != nil {
		return 0, fmt.Errorf("failed to write markdown header: %w", err)
	}

	if len(cfg.ExtractedPackages) > 0 {
		_, err = writer.WriteString("## Packages\n")
		if err != nil {
			return 0, fmt.Errorf("failed to write packages header: %w", err)
		}
		_, err = writer.WriteString("| type | scope | name | version |\n| --- | --- | --- | --- |\n")
		if err != nil {
			return 0, fmt.Errorf("failed to write packages table header: %w", err)
		}
		for _, pkg := range cfg.ExtractedPackages {
			scope := pkg.Scope
			if scope == "" {
				scope = "-"
			}
			_, err = fmt.Fprintf(
				writer,
				"| %s | %s | %s | %s |\n",
				escapeMarkdownTable(pkg.Type),
				escapeMarkdownTable(scope),
				escapeMarkdownTable(pkg.Name),
				escapeMarkdownTable(pkg.Version),
			)
			if err != nil {
				return 0, fmt.Errorf("failed to write packages table: %w", err)
			}
		}
		_, err = writer.WriteString("\n")
		if err != nil {
			return 0, fmt.Errorf("failed to write packages separator: %w", err)
		}
	}

	_, err = writer.WriteString("## Files\n\n")
	if err != nil {
		return 0, fmt.Errorf("failed to write files header: %w", err)
	}

	for _, relPath := range relPaths {
		absPath := filesToInclude[relPath]
		contentBytes, readErr := os.ReadFile(absPath)
		if readErr != nil {
			log.Printf("Warning: Could not read file '%s'. Content will be marked as failed. Error: %v", absPath, readErr)
			contentBytes = []byte(fmt.Sprintf("FAILED_TO_READ_FILE: %v", readErr))
		} else if len(contentBytes) == 0 {
			log.Printf("Info: File '%s' is empty.", absPath)
		}

		_, err = writer.WriteString("### `" + relPath + "`\n")
		if err != nil {
			return 0, fmt.Errorf("failed to write file header: %w", err)
		}

		fence := markdownFence(contentBytes)
		_, err = writer.WriteString(fence + "\n")
		if err != nil {
			return 0, fmt.Errorf("failed to write code fence: %w", err)
		}

		if len(contentBytes) > 0 {
			_, err = writer.Write(contentBytes)
			if err != nil {
				return 0, fmt.Errorf("failed to write file contents: %w", err)
			}
		}
		if len(contentBytes) == 0 || contentBytes[len(contentBytes)-1] != '\n' {
			_, err = writer.WriteString("\n")
			if err != nil {
				return 0, fmt.Errorf("failed to write content newline: %w", err)
			}
		}

		_, err = writer.WriteString(fence + "\n\n")
		if err != nil {
			return 0, fmt.Errorf("failed to close code fence: %w", err)
		}

		processedFilesCount++
	}

	err = writer.Flush()
	if err != nil {
		return 0, fmt.Errorf("failed to flush Markdown writer: %w", err)
	}

	log.Printf("Successfully wrote %d files to %s", processedFilesCount, cfg.OutputFile)
	return processedFilesCount, nil
}

func escapeMarkdownTable(value string) string {
	escaped := strings.ReplaceAll(value, "|", "\\|")
	escaped = strings.ReplaceAll(escaped, "\r", " ")
	escaped = strings.ReplaceAll(escaped, "\n", " ")
	return escaped
}

func markdownFence(content []byte) string {
	maxRun := 0
	run := 0
	for _, b := range content {
		if b == '`' {
			run++
			if run > maxRun {
				maxRun = run
			}
		} else {
			run = 0
		}
	}
	fenceLen := 3
	if maxRun+1 > fenceLen {
		fenceLen = maxRun + 1
	}
	return strings.Repeat("`", fenceLen)
}
