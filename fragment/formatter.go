package main

import (
	"fmt"
	"io"
)

// Formatter writes a project snapshot in a specific format.
type Formatter interface {
	Name() string
	Write(io.Writer, *ProjectSnapshot) error
}

func getFormatter(format string) (Formatter, error) {
	switch format {
	case outputFormatXML:
		return xmlFormatter{}, nil
	case outputFormatMarkdown:
		return markdownFormatter{}, nil
	case outputFormatJSON:
		return jsonFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}
