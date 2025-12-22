package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func buildApiSnapshot(cfg *Config, filesToInclude map[string]string) (*ProjectSnapshot, int, error) {
	projectName := filepath.Base(cfg.ProjectRoot)

	relPaths := make([]string, 0, len(filesToInclude))
	for relPath := range filesToInclude {
		relPaths = append(relPaths, relPath)
	}
	sort.Strings(relPaths)

	apiIndex := &ApiIndex{}
	for _, relPath := range relPaths {
		absPath := filesToInclude[relPath]
		ext := strings.ToLower(filepath.Ext(relPath))
		parser := apiParserForExt(ext)
		if parser == nil {
			continue
		}

		src, err := os.ReadFile(absPath)
		if err != nil {
			log.Printf("Warning: Could not read file '%s' for API index: %v", absPath, err)
			continue
		}

		apiFile, err := parser.Parse(relPath, src)
		if err != nil {
			log.Printf("Warning: Could not parse API from '%s': %v", relPath, err)
			continue
		}
		if apiFile != nil && len(apiFile.Symbols) > 0 {
			apiIndex.Files = append(apiIndex.Files, *apiFile)
		}
	}

	snapshot := &ProjectSnapshot{
		Name:     projectName,
		Packages: cfg.ExtractedPackages,
		Tree:     nil,
		Metadata: cfg.Metadata,
		API:      apiIndex,
		ApiOnly:  true,
	}

	if !cfg.NoTree {
		snapshot.Tree = buildDirTree(relPaths)
	}

	return snapshot, len(apiIndex.Files), nil
}

type apiParser interface {
	Language() string
	Extensions() []string
	Parse(path string, src []byte) (*ApiFile, error)
}

var apiParsers = []apiParser{
	newJavaScriptParser(),
	newTypeScriptParser(),
	newPhpParser(),
}

func apiParserForExt(ext string) apiParser {
	for _, parser := range apiParsers {
		for _, candidate := range parser.Extensions() {
			if ext == candidate {
				return parser
			}
		}
	}
	return nil
}

func sanitizeDoc(doc string) string {
	doc = strings.ReplaceAll(doc, "\r\n", "\n")
	doc = strings.ReplaceAll(doc, "\r", "\n")
	lines := strings.Split(doc, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func signatureFromHeader(header string) string {
	header = strings.TrimSpace(header)
	header = strings.TrimSuffix(header, "{")
	header = strings.TrimSpace(header)
	if strings.HasSuffix(header, ";") {
		header = strings.TrimSuffix(header, ";")
	}
	return strings.TrimSpace(header)
}

func truncateSignature(sig string, limit int) string {
	if len(sig) <= limit {
		return sig
	}
	return sig[:limit] + "..."
}

func signatureFromNode(src []byte, start, end uint) string {
	if int(start) >= len(src) || int(end) > len(src) || start >= end {
		return ""
	}
	header := src[start:end]
	for i, b := range header {
		if b == '{' || b == ';' || b == '\n' {
			header = header[:i]
			break
		}
	}
	return signatureFromHeader(string(header))
}

func variableSignature(prefix, name string) string {
	if name == "" {
		return strings.TrimSpace(prefix)
	}
	if prefix == "" {
		return name
	}
	return strings.TrimSpace(prefix + " " + name)
}

func prefixFromDeclaration(src []byte, start, end uint) string {
	if int(start) >= len(src) || int(end) > len(src) || start >= end {
		return ""
	}
	header := src[start:end]
	for i, b := range header {
		if b == ' ' || b == '\n' || b == '\t' {
			return string(header[:i])
		}
	}
	return string(header)
}

func debugSignature(kind, name, sig string) string {
	if sig != "" {
		return sig
	}
	return fmt.Sprintf("%s %s", kind, name)
}
