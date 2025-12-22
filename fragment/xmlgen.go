// xmlgen.go
package main

import (
	"encoding/xml"
	"fmt"
	"io"
)

type xmlFormatter struct{}

func (xmlFormatter) Name() string {
	return outputFormatXML
}

func (xmlFormatter) Write(w io.Writer, snapshot *ProjectSnapshot) error {
	project := Project{
		Name:     snapshot.Name,
		Metadata: snapshot.Metadata,
	}

	if len(snapshot.Packages) > 0 {
		project.Packages = &PackagesHolder{
			PackageList: snapshot.Packages,
		}
	}

	if snapshot.Tree != nil {
		project.Tree = &Tree{
			Nodes: toXMLTreeNodes(snapshot.Tree.Children),
		}
	}

	if snapshot.API != nil {
		project.API = snapshot.API
	}

	if !snapshot.ApiOnly && len(snapshot.Files) > 0 {
		project.Files = &Files{
			Files: make([]File, 0, len(snapshot.Files)),
		}
		for _, entry := range snapshot.Files {
			project.Files.Files = append(project.Files.Files, File{
				Path:    entry.Path,
				Content: entry.Content,
			})
		}
	}

	if _, err := io.WriteString(w, xml.Header); err != nil {
		return fmt.Errorf("failed to write XML header: %w", err)
	}

	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	if err := encoder.Encode(project); err != nil {
		return fmt.Errorf("failed to encode XML: %w", err)
	}

	if _, err := io.WriteString(w, "\n"); err != nil {
		return fmt.Errorf("failed to write XML trailer: %w", err)
	}
	return nil
}

func toXMLTreeNodes(nodes []*DirNode) []TreeNode {
	if len(nodes) == 0 {
		return nil
	}
	out := make([]TreeNode, 0, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}
		kind := "file"
		if node.IsDir {
			kind = "dir"
		}
		out = append(out, TreeNode{
			Name:     node.Name,
			Path:     node.Path,
			Kind:     kind,
			Children: toXMLTreeNodes(node.Children),
		})
	}
	return out
}
