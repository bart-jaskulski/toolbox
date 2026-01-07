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
		Packages: snapshot.Packages,
		Metadata: snapshot.Metadata,
		API:      snapshot.API,
	}
	if snapshot.Tree != nil {
		out.Tree = snapshot.Tree.Children
	}

	if !snapshot.ApiOnly {
		out.Files = make([]jsonFile, 0, len(snapshot.Files))
		for _, entry := range snapshot.Files {
			out.Files = append(out.Files, jsonFile{
				Path:     entry.Path,
				Language: entry.Language,
				Content:  string(entry.Content),
			})
		}
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
	Files    []jsonFile       `json:"files,omitempty"`
	Packages []ProjectPackage `json:"packages,omitempty"`
	Tree     []*DirNode       `json:"tree,omitempty"`
	Metadata *ProjectMetadata `json:"metadata,omitempty"`
	API      *ApiIndex        `json:"api,omitempty"`
}

type jsonFile struct {
	Path     string `json:"path"`
	Language string `json:"language,omitempty"`
	Content  string `json:"content"`
}
