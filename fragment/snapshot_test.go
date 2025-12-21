package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildDirTree(t *testing.T) {
	relPaths := []string{
		"a/b.txt",
		"a/c/d.go",
		"z.txt",
	}

	tree := buildDirTree(relPaths)
	if tree == nil {
		t.Fatalf("expected tree")
	}
	if len(tree.Children) != 2 {
		t.Fatalf("expected 2 root children, got %d", len(tree.Children))
	}
	if tree.Children[0].Name != "a" || !tree.Children[0].IsDir {
		t.Fatalf("expected first child to be directory 'a'")
	}
	if tree.Children[1].Name != "z.txt" || tree.Children[1].IsDir {
		t.Fatalf("expected second child to be file 'z.txt'")
	}
}

func TestBuildSnapshotNoTree(t *testing.T) {
	dir := t.TempDir()
	fileA := filepath.Join(dir, "a.txt")
	fileB := filepath.Join(dir, "b.txt")

	if err := os.WriteFile(fileA, []byte("alpha\n"), 0o600); err != nil {
		t.Fatalf("write fileA: %v", err)
	}
	if err := os.WriteFile(fileB, []byte("beta\n"), 0o600); err != nil {
		t.Fatalf("write fileB: %v", err)
	}

	cfg := Config{
		ProjectRoot: dir,
		NoTree:      true,
	}
	files := map[string]string{
		"b.txt": fileB,
		"a.txt": fileA,
	}

	snapshot, count, err := buildSnapshot(&cfg, files)
	if err != nil {
		t.Fatalf("buildSnapshot error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 files, got %d", count)
	}
	if snapshot.Tree != nil {
		t.Fatalf("expected tree to be nil when NoTree is set")
	}
	if len(snapshot.Files) != 2 {
		t.Fatalf("expected 2 files in snapshot, got %d", len(snapshot.Files))
	}
	if snapshot.Files[0].Path != "a.txt" || snapshot.Files[1].Path != "b.txt" {
		t.Fatalf("expected files sorted by path, got %q then %q", snapshot.Files[0].Path, snapshot.Files[1].Path)
	}
}
