package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGatherFilesGitRespectsIgnore(t *testing.T) {
	if _, err := runCommand("", "git", "--version"); err != nil {
		t.Skip("git not available")
	}

	root := t.TempDir()
	if _, err := runCommand(root, "git", "init"); err != nil {
		t.Fatalf("git init: %v", err)
	}

	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("*.log\nignored_dir/\n"), 0o600); err != nil {
		t.Fatalf("write .gitignore: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "keep.txt"), []byte("keep\n"), 0o600); err != nil {
		t.Fatalf("write keep.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "ignored.log"), []byte("ignore\n"), 0o600); err != nil {
		t.Fatalf("write ignored.log: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "ignored_dir"), 0o700); err != nil {
		t.Fatalf("mkdir ignored_dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "ignored_dir", "file.txt"), []byte("ignore\n"), 0o600); err != nil {
		t.Fatalf("write ignored_dir/file.txt: %v", err)
	}

	cfg := Config{
		ProjectRoot:        root,
		AbsoluteInputDirs:  []string{root},
		AbsoluteOutputFile: filepath.Join(root, "output.txt"),
		CombinedExcludes:   nil,
		IncludeBinary:      true,
		NoIgnore:           false,
	}
	files := make(map[string]string)

	if err := gatherFilesGit(&cfg, files); err != nil {
		t.Fatalf("gatherFilesGit: %v", err)
	}

	if _, ok := files["keep.txt"]; !ok {
		t.Fatalf("expected keep.txt to be included")
	}
	if _, ok := files["ignored.log"]; ok {
		t.Fatalf("expected ignored.log to be excluded")
	}
	if _, ok := files[filepath.Join("ignored_dir", "file.txt")]; ok {
		t.Fatalf("expected ignored_dir/file.txt to be excluded")
	}
}

func TestGatherFilesGitNoIgnore(t *testing.T) {
	if _, err := runCommand("", "git", "--version"); err != nil {
		t.Skip("git not available")
	}

	root := t.TempDir()
	if _, err := runCommand(root, "git", "init"); err != nil {
		t.Fatalf("git init: %v", err)
	}

	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("*.log\n"), 0o600); err != nil {
		t.Fatalf("write .gitignore: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "keep.txt"), []byte("keep\n"), 0o600); err != nil {
		t.Fatalf("write keep.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "ignored.log"), []byte("ignore\n"), 0o600); err != nil {
		t.Fatalf("write ignored.log: %v", err)
	}

	cfg := Config{
		ProjectRoot:        root,
		AbsoluteInputDirs:  []string{root},
		AbsoluteOutputFile: filepath.Join(root, "output.txt"),
		CombinedExcludes:   nil,
		IncludeBinary:      true,
		NoIgnore:           true,
	}
	files := make(map[string]string)

	if err := gatherFilesGit(&cfg, files); err != nil {
		t.Fatalf("gatherFilesGit: %v", err)
	}

	if _, ok := files["keep.txt"]; !ok {
		t.Fatalf("expected keep.txt to be included")
	}
	if _, ok := files["ignored.log"]; !ok {
		t.Fatalf("expected ignored.log to be included when --no-ignore is set")
	}
}
