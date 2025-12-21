package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckExclusion(t *testing.T) {
	patterns := []string{
		"*.lock",
		"package-lock.json",
		"node_modules/",
		"vendor",
		"exact.txt",
		"src/*.go",
	}

	tests := []struct {
		path string
		want bool
	}{
		{path: "composer.lock", want: true},
		{path: "nested/package-lock.json", want: true},
		{path: "node_modules/react/index.js", want: true},
		{path: "vendor/github.com/foo", want: true},
		{path: "exact.txt", want: true},
		{path: "dir/exact.txt", want: true},
		{path: "src/main.go", want: true},
		{path: "src/main.ts", want: false},
	}

	for _, tt := range tests {
		if got := checkExclusion(tt.path, patterns); got != tt.want {
			t.Fatalf("checkExclusion(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestIsLikelyBinary(t *testing.T) {
	dir := t.TempDir()

	textPath := filepath.Join(dir, "text.txt")
	if err := os.WriteFile(textPath, []byte("hello\nworld\n"), 0o600); err != nil {
		t.Fatalf("write text file: %v", err)
	}
	isBin, err := isLikelyBinary(textPath)
	if err != nil {
		t.Fatalf("isLikelyBinary text: %v", err)
	}
	if isBin {
		t.Fatalf("expected text file to be non-binary")
	}

	binPath := filepath.Join(dir, "bin.dat")
	if err := os.WriteFile(binPath, []byte{0x00, 0x01, 0x02}, 0o600); err != nil {
		t.Fatalf("write binary file: %v", err)
	}
	isBin, err = isLikelyBinary(binPath)
	if err != nil {
		t.Fatalf("isLikelyBinary binary: %v", err)
	}
	if !isBin {
		t.Fatalf("expected binary file to be detected")
	}
}
