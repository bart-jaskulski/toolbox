package main

import (
	"fmt"
	"strings"
)

const (
	outputFormatXML      = "xml"
	outputFormatMarkdown = "markdown"
)

func normalizeOutputFormat(raw string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	switch normalized {
	case "", outputFormatXML:
		return outputFormatXML, nil
	case "md", outputFormatMarkdown:
		return outputFormatMarkdown, nil
	default:
		return "", fmt.Errorf("unsupported format %q (use 'xml' or 'markdown')", raw)
	}
}

func defaultOutputFile(format string) string {
	if format == outputFormatMarkdown {
		return "output.md"
	}
	return "output.xml"
}

func formatDisplayName(format string) string {
	switch format {
	case outputFormatXML:
		return "XML"
	case outputFormatMarkdown:
		return "Markdown"
	default:
		return format
	}
}
