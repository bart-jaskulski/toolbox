package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"
)

func TestXMLFormatter(t *testing.T) {
	snapshot := &ProjectSnapshot{
		Name: "demo",
		Files: []FileEntry{
			{Path: "a.txt", Content: []byte("hello"), Language: "text"},
		},
		Packages: []ProjectPackage{
			{Type: "npm", Scope: "dependencies", Name: "react", Version: "18.2.0"},
		},
		Tree: &DirNode{
			Name:  ".",
			IsDir: true,
			Children: []*DirNode{
				{Name: "a.txt", Path: "a.txt", IsDir: false},
			},
		},
	}

	var buf bytes.Buffer
	if err := (xmlFormatter{}).Write(&buf, snapshot); err != nil {
		t.Fatalf("xml formatter: %v", err)
	}

	var project Project
	if err := xml.Unmarshal(buf.Bytes(), &project); err != nil {
		t.Fatalf("xml unmarshal: %v", err)
	}
	if project.Name != "demo" {
		t.Fatalf("expected project name demo, got %q", project.Name)
	}
	if project.Packages == nil || len(project.Packages.PackageList) != 1 {
		t.Fatalf("expected 1 package")
	}
	if project.Files == nil || len(project.Files.Files) != 1 {
		t.Fatalf("expected 1 file")
	}
	if project.Files.Files[0].Language != "text" {
		t.Fatalf("expected file language text, got %q", project.Files.Files[0].Language)
	}
	if project.Tree == nil || len(project.Tree.Nodes) != 1 {
		t.Fatalf("expected 1 tree node")
	}
}

func TestJSONFormatter(t *testing.T) {
	snapshot := &ProjectSnapshot{
		Name: "demo",
		Files: []FileEntry{
			{Path: "a.txt", Content: []byte("hello"), Language: "text"},
		},
		Packages: []ProjectPackage{
			{Type: "npm", Scope: "dependencies", Name: "react", Version: "18.2.0"},
		},
		Tree: &DirNode{
			Name:  ".",
			IsDir: true,
			Children: []*DirNode{
				{Name: "a.txt", Path: "a.txt", IsDir: false},
			},
		},
	}

	var buf bytes.Buffer
	if err := (jsonFormatter{}).Write(&buf, snapshot); err != nil {
		t.Fatalf("json formatter: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if out["name"] != "demo" {
		t.Fatalf("expected name demo, got %v", out["name"])
	}
	if files, ok := out["files"].([]any); !ok || len(files) != 1 {
		t.Fatalf("expected 1 file in json output")
	}
	if files, ok := out["files"].([]any); ok && len(files) == 1 {
		fileObj, ok := files[0].(map[string]any)
		if !ok {
			t.Fatalf("expected file to be object")
		}
		if fileObj["language"] != "text" {
			t.Fatalf("expected file language text, got %v", fileObj["language"])
		}
	}
	if tree, ok := out["tree"].([]any); !ok || len(tree) != 1 {
		t.Fatalf("expected tree in json output")
	}
}

func TestJSONFormatterNoTree(t *testing.T) {
	snapshot := &ProjectSnapshot{
		Name: "demo",
		Files: []FileEntry{
			{Path: "a.txt", Content: []byte("hello"), Language: "text"},
		},
	}

	var buf bytes.Buffer
	if err := (jsonFormatter{}).Write(&buf, snapshot); err != nil {
		t.Fatalf("json formatter: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if _, ok := out["tree"]; ok {
		t.Fatalf("expected tree to be omitted when not present")
	}
}

func TestXMLFormatterNoTree(t *testing.T) {
	snapshot := &ProjectSnapshot{
		Name: "demo",
		Files: []FileEntry{
			{Path: "a.txt", Content: []byte("hello"), Language: "text"},
		},
	}

	var buf bytes.Buffer
	if err := (xmlFormatter{}).Write(&buf, snapshot); err != nil {
		t.Fatalf("xml formatter: %v", err)
	}

	var project Project
	if err := xml.Unmarshal(buf.Bytes(), &project); err != nil {
		t.Fatalf("xml unmarshal: %v", err)
	}
	if project.Tree != nil {
		t.Fatalf("expected tree to be omitted when not present")
	}
}

func TestMarkdownFormatterTree(t *testing.T) {
	snapshot := &ProjectSnapshot{
		Name: "demo",
		Files: []FileEntry{
			{Path: "a.txt", Content: []byte("hello"), Language: "text"},
		},
		Tree: &DirNode{
			Name:  ".",
			IsDir: true,
			Children: []*DirNode{
				{Name: "a.txt", Path: "a.txt", IsDir: false},
			},
		},
	}

	var buf bytes.Buffer
	if err := (markdownFormatter{}).Write(&buf, snapshot); err != nil {
		t.Fatalf("markdown formatter: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "## Tree") {
		t.Fatalf("expected tree header")
	}
	if !strings.Contains(out, "```") {
		t.Fatalf("expected tree fenced block")
	}
	if !strings.Contains(out, "a.txt") {
		t.Fatalf("expected tree output to include file")
	}
	if !strings.Contains(out, "### `a.txt`") {
		t.Fatalf("expected file heading")
	}
}
