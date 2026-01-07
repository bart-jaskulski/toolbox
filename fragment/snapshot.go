package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ProjectSnapshot holds normalized project data for formatters.
type ProjectSnapshot struct {
	Name     string
	Files    []FileEntry
	Packages []ProjectPackage
	Tree     *DirNode
	Metadata *ProjectMetadata
	API      *ApiIndex
	ApiOnly  bool
}

// DirNode represents a directory or file node.
type DirNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path,omitempty"`
	IsDir    bool       `json:"is_dir"`
	Children []*DirNode `json:"children,omitempty"`
}

func buildSnapshot(cfg *Config, filesToInclude map[string]string) (*ProjectSnapshot, int, error) {
	projectName := filepath.Base(cfg.ProjectRoot)

	relPaths := make([]string, 0, len(filesToInclude))
	for relPath := range filesToInclude {
		relPaths = append(relPaths, relPath)
	}
	sort.Strings(relPaths)

	files := make([]FileEntry, 0, len(relPaths))
	for _, relPath := range relPaths {
		absPath := filesToInclude[relPath]
		contentBytes, err := os.ReadFile(absPath)
		if err != nil {
			log.Printf("Warning: Could not read file '%s'. Content will be marked as failed. Error: %v", absPath, err)
			contentBytes = []byte(fmt.Sprintf("FAILED_TO_READ_FILE: %v", err))
		} else if len(contentBytes) == 0 {
			log.Printf("Info: File '%s' is empty.", absPath)
		}

		files = append(files, FileEntry{
			Path:     relPath,
			Content:  contentBytes,
			Language: languageForPath(relPath),
		})
	}

	snapshot := &ProjectSnapshot{
		Name:     projectName,
		Files:    files,
		Packages: cfg.ExtractedPackages,
		Metadata: cfg.Metadata,
	}
	if !cfg.NoTree {
		snapshot.Tree = buildDirTree(relPaths)
	}

	return snapshot, len(files), nil
}

type dirBuilder struct {
	Name     string
	Path     string
	IsDir    bool
	Children map[string]*dirBuilder
}

func buildDirTree(relPaths []string) *DirNode {
	root := &dirBuilder{
		Name:     ".",
		Path:     "",
		IsDir:    true,
		Children: make(map[string]*dirBuilder),
	}

	for _, relPath := range relPaths {
		parts := strings.Split(filepath.ToSlash(relPath), "/")
		current := root
		currentPath := ""
		for i, part := range parts {
			if part == "" {
				continue
			}
			if currentPath == "" {
				currentPath = part
			} else {
				currentPath = currentPath + "/" + part
			}
			isDir := i < len(parts)-1
			child := current.Children[part]
			if child == nil {
				child = &dirBuilder{
					Name:     part,
					Path:     currentPath,
					IsDir:    isDir,
					Children: make(map[string]*dirBuilder),
				}
				current.Children[part] = child
			}
			current = child
		}
	}

	return finalizeTree(root)
}

func finalizeTree(node *dirBuilder) *DirNode {
	if node == nil {
		return nil
	}
	children := make([]*DirNode, 0, len(node.Children))
	names := make([]string, 0, len(node.Children))
	for name := range node.Children {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		children = append(children, finalizeTree(node.Children[name]))
	}
	return &DirNode{
		Name:     node.Name,
		Path:     node.Path,
		IsDir:    node.IsDir,
		Children: children,
	}
}
