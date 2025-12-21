package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type jsonFormatter struct{}

func (jsonFormatter) Name() string {
	return outputFormatJSON
}

func (jsonFormatter) Write(w io.Writer, snapshot *ProjectSnapshot) error {
	out := jsonProject{
		Name:     snapshot.Name,
		Files:    make([]jsonFile, 0, len(snapshot.Files)),
		Packages: snapshot.Packages,
		Metadata: snapshot.Metadata,
	}
	if snapshot.Tree != nil {
		out.Tree = snapshot.Tree.Children
	}

	for _, entry := range snapshot.Files {
		out.Files = append(out.Files, jsonFile{
			Path:    entry.Path,
			Content: string(entry.Content),
		})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(out); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

type jsonProject struct {
	Name     string           `json:"name"`
	Files    []jsonFile       `json:"files"`
	Packages []ProjectPackage `json:"packages,omitempty"`
	Tree     []*DirNode       `json:"tree,omitempty"`
	Metadata *ProjectMetadata `json:"metadata,omitempty"`
}

type jsonFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}
