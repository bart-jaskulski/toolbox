package main

import "testing"

func TestLanguageForPath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{path: "Main.java", want: "java"},
		{path: "script.sh", want: "bash"},
		{path: "README.md", want: "markdown"},
		{path: "Dockerfile", want: "dockerfile"},
		{path: "Makefile", want: "makefile"},
		{path: "notes.abc", want: "abc"},
		{path: "noext", want: ""},
	}

	for _, tt := range tests {
		if got := languageForPath(tt.path); got != tt.want {
			t.Fatalf("languageForPath(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}
