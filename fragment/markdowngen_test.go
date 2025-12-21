package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteMarkdownTree(t *testing.T) {
	nodes := []*DirNode{
		{
			Name:  "a",
			IsDir: true,
			Children: []*DirNode{
				{Name: "b.txt", IsDir: false},
				{Name: "c.txt", IsDir: false},
			},
		},
		{Name: "z.txt", IsDir: false},
	}

	var buf bytes.Buffer
	if err := writeMarkdownTree(&buf, nodes, "", true); err != nil {
		t.Fatalf("writeMarkdownTree: %v", err)
	}

	got := buf.String()
	want := strings.Join([]string{
		"a/",
		"|-- b.txt",
		"`-- c.txt",
		"z.txt",
		"",
	}, "\n")

	if got != want {
		t.Fatalf("unexpected tree output:\n%s", got)
	}
}

func TestMarkdownFormatterNoTree(t *testing.T) {
	snapshot := &ProjectSnapshot{
		Name: "demo",
		Files: []FileEntry{
			{Path: "a.txt", Content: []byte("hello")},
		},
	}

	var buf bytes.Buffer
	if err := (markdownFormatter{}).Write(&buf, snapshot); err != nil {
		t.Fatalf("markdown formatter: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "## Tree") {
		t.Fatalf("did not expect tree section when tree is nil")
	}
	if !strings.Contains(out, "### `a.txt`") {
		t.Fatalf("expected file heading")
	}
}
