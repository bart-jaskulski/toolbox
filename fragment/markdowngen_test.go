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
