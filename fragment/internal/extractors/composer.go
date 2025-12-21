// internal/extractors/composer.go
package extractors

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    // "sort" // If sorting packages later
)

// composerJSON represents the relevant parts of composer.json for internal parsing.
type composerJSON struct {
    Require    map[string]string `json:"require"`
    RequireDev map[string]string `json:"require-dev"`
}

// ComposerExtractor handles composer.json
type ComposerExtractor struct{}

// init registers this extractor with the registry.
func init() {
    Register("composer", ComposerExtractor{})
}

func (e ComposerExtractor) FileName() string {
    return "composer.json"
}

func (e ComposerExtractor) Extract(projectRoot string) (*ExtractedPackagesResult, error) {
	filePath := filepath.Join(projectRoot, e.FileName())
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // File not found is not an error for extraction attempt
		}
		return nil, fmt.Errorf("reading %s: %w", e.FileName(), err)
	}

	var cj composerJSON
	if err := json.Unmarshal(data, &cj); err != nil {
		return nil, fmt.Errorf("parsing %s: JSON decoding error: %w", e.FileName(), err)
	}

	var groups []ScopedPackageGroup
	if len(cj.Require) > 0 {
		groups = append(groups, ScopedPackageGroup{
			ScopeName: "require",
			Packages:  mapToInternalPackages(cj.Require),
		})
	}
	if len(cj.RequireDev) > 0 {
		groups = append(groups, ScopedPackageGroup{
			ScopeName: "require-dev",
			Packages:  mapToInternalPackages(cj.RequireDev),
		})
	}

	if len(groups) == 0 {
		return nil, nil // No packages found, not an error
	}

	return &ExtractedPackagesResult{
		PackageType: "composer",
		Groups:      groups,
	}, nil
}

// mapToInternalPackages is now a shared utility in extractor.go
