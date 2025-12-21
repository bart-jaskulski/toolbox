// xmlgen.go
package main

import (
    "bufio"
    "encoding/xml"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sort" // Uncomment if sorting files
)

// generateXML creates the final XML output file, including metadata.
func generateXML(cfg *Config, filesToInclude map[string]string) (int, error) {
    log.Println("Generating XML output...")
    processedFilesCount := 0

    projectName := filepath.Base(cfg.ProjectRoot)

    project := Project{
        Name:     projectName,
        Metadata: cfg.Metadata, // Assign non-package metadata from config (likely nil for now)
        Files: Files{
            Files: make([]File, 0, len(filesToInclude)),
        },
    }

    // Populate unified packages section
    if len(cfg.ExtractedPackages) > 0 {
        project.Packages = &PackagesHolder{
            PackageList: cfg.ExtractedPackages,
        }
    }

    // Sort files by relative path for consistent output
    relPaths := make([]string, 0, len(filesToInclude))
    for relPath := range filesToInclude {
        relPaths = append(relPaths, relPath)
    }
    sort.Strings(relPaths)

    for _, relPath := range relPaths {
        absPath := filesToInclude[relPath]
        contentBytes, err := os.ReadFile(absPath)
        if err != nil {
            log.Printf("Warning: Could not read file '%s'. Content will be marked as failed. Error: %v", absPath, err)
            contentBytes = []byte(fmt.Sprintf(" FAILED_TO_READ_FILE: %v ", err))
        } else if len(contentBytes) == 0 {
            log.Printf("Info: File '%s' is empty.", absPath)
        }

        project.Files.Files = append(project.Files.Files, File{
            Path:    relPath,
            Content: contentBytes,
        })
        processedFilesCount++
    }

    outFile, err := os.Create(cfg.OutputFile)
    if err != nil {
        return 0, fmt.Errorf("failed to create output file '%s': %w", cfg.OutputFile, err)
    }
    defer outFile.Close()

    writer := bufio.NewWriter(outFile)

    _, err = writer.WriteString(xml.Header)
    if err != nil {
        return 0, fmt.Errorf("failed to write XML header: %w", err)
    }

    encoder := xml.NewEncoder(writer)
    encoder.Indent("", "  ")
    err = encoder.Encode(project) // Encode the project struct including metadata
    if err != nil {
        return 0, fmt.Errorf("failed to encode XML: %w", err)
    }

    _, _ = writer.WriteString("\n")

    err = writer.Flush()
    if err != nil {
        return 0, fmt.Errorf("failed to flush XML writer: %w", err)
    }

    log.Printf("Successfully wrote %d files to %s", processedFilesCount, cfg.OutputFile)
    return processedFilesCount, nil
}
