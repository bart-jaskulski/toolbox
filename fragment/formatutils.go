package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	outputFormatXML      = "xml"
	outputFormatMarkdown = "markdown"
	outputFormatJSON     = "json"
)

func normalizeOutputFormat(raw string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	switch normalized {
	case "", outputFormatXML:
		return outputFormatXML, nil
	case "md", outputFormatMarkdown:
		return outputFormatMarkdown, nil
	case outputFormatJSON:
		return outputFormatJSON, nil
	default:
		return "", fmt.Errorf("unsupported format %q (use 'xml', 'markdown', or 'json')", raw)
	}
}

func inferFormatFromOutput(outputFile string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(outputFile)))
	switch ext {
	case ".xml":
		return outputFormatXML, true
	case ".md", ".markdown":
		return outputFormatMarkdown, true
	case ".json":
		return outputFormatJSON, true
	default:
		return "", false
	}
}

func defaultOutputFile(format string) string {
	if format == outputFormatMarkdown {
		return "output.md"
	}
	if format == outputFormatJSON {
		return "output.json"
	}
	return "output.xml"
}

func formatDisplayName(format string) string {
	switch format {
	case outputFormatXML:
		return "XML"
	case outputFormatMarkdown:
		return "Markdown"
	case outputFormatJSON:
		return "JSON"
	default:
		return format
	}
}
