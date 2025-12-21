// internal/extractors/npm.go
package extractors

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    // "sort" // If sorting packages later
)

// packageJSON represents the relevant parts of package.json for internal parsing.
type packageJSON struct {
    Dependencies     map[string]string `json:"dependencies"`
    DevDependencies  map[string]string `json:"devDependencies"`
    PeerDependencies map[string]string `json:"peerDependencies"`
}

// NpmExtractor handles package.json
type NpmExtractor struct{}

// init registers this extractor.
func init() {
    Register("npm", NpmExtractor{})
}

func (e NpmExtractor) FileName() string {
    return "package.json"
}

func (e NpmExtractor) Extract(projectRoot string) (*ExtractedPackagesResult, error) {
	filePath := filepath.Join(projectRoot, e.FileName())
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Not found is ok
		}
		return nil, fmt.Errorf("reading %s: %w", e.FileName(), err)
	}

	var pj packageJSON
	if err := json.Unmarshal(data, &pj); err != nil {
		return nil, fmt.Errorf("parsing %s: JSON decoding error: %w", e.FileName(), err)
	}

	var groups []ScopedPackageGroup
	if len(pj.Dependencies) > 0 {
		groups = append(groups, ScopedPackageGroup{
			ScopeName: "dependencies",
			Packages:  mapToInternalPackages(pj.Dependencies),
		})
	}
	if len(pj.DevDependencies) > 0 {
		groups = append(groups, ScopedPackageGroup{
			ScopeName: "devDependencies",
			Packages:  mapToInternalPackages(pj.DevDependencies),
		})
	}
	if len(pj.PeerDependencies) > 0 {
		groups = append(groups, ScopedPackageGroup{
			ScopeName: "peerDependencies",
			Packages:  mapToInternalPackages(pj.PeerDependencies),
		})
	}

	if len(groups) == 0 {
		return nil, nil // No packages found, not an error
	}

	return &ExtractedPackagesResult{
		PackageType: "npm",
		Groups:      groups,
	}, nil
}

// mapToInternalPackages is now a shared utility in extractor.go
