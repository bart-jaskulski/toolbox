package main

import (
	"path/filepath"
	"testing"
)

func TestInputsUnderRoot(t *testing.T) {
	root := filepath.Join(string(filepath.Separator), "repo")
	inRoot := []string{
		root,
		filepath.Join(root, "subdir"),
	}
	outside := []string{
		root,
		filepath.Join(string(filepath.Separator), "other"),
	}

	if !inputsUnderRoot(root, inRoot) {
		t.Fatalf("expected inputs under root")
	}
	if inputsUnderRoot(root, outside) {
		t.Fatalf("expected inputs outside root to return false")
	}
}
