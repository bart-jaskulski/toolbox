package main

import (
	"path/filepath"
	"strings"
)

func languageForPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		base := strings.ToLower(filepath.Base(path))
		switch base {
		case "dockerfile":
			return "dockerfile"
		case "makefile":
			return "makefile"
		}
		return ""
	}

	if lang, ok := languageByExtension[ext]; ok {
		return lang
	}

	return strings.TrimPrefix(ext, ".")
}

var languageByExtension = map[string]string{
	".bash":     "bash",
	".c":        "c",
	".cc":       "cpp",
	".cfg":      "ini",
	".conf":     "conf",
	".cpp":      "cpp",
	".cs":       "csharp",
	".css":      "css",
	".cxx":      "cpp",
	".go":       "go",
	".h":        "c",
	".hpp":      "cpp",
	".htm":      "html",
	".html":     "html",
	".ini":      "ini",
	".java":     "java",
	".js":       "javascript",
	".json":     "json",
	".jsx":      "jsx",
	".md":       "markdown",
	".markdown": "markdown",
	".php":      "php",
	".py":       "python",
	".rb":       "ruby",
	".rs":       "rust",
	".sass":     "sass",
	".scss":     "scss",
	".sh":       "bash",
	".sql":      "sql",
	".toml":     "toml",
	".ts":       "typescript",
	".tsx":      "tsx",
	".txt":      "text",
	".xml":      "xml",
	".yaml":     "yaml",
	".yml":      "yaml",
	".zsh":      "zsh",
}
