// internal/extractors/extractor.go
package extractors

// ScopedPackageGroup holds a list of packages for a specific scope (e.g., "dependencies", "require-dev").
type ScopedPackageGroup struct {
	ScopeName string    // e.g., "require", "dependencies", "devDependencies"
	Packages  []Package // The list of packages for this scope
}

// ExtractedPackagesResult is a standardized structure for extractors to return package information.
type ExtractedPackagesResult struct {
	PackageType string               // e.g., "composer", "npm", "gomod"
	Groups      []ScopedPackageGroup // A list of package groups (e.g., require, require-dev)
}

// MetadataExtractor defines the contract for extracting metadata from a specific file.
type MetadataExtractor interface {
	// FileName returns the name of the metadata file this extractor handles (e.g., "composer.json").
	FileName() string
	// Extract attempts to parse the file (located at projectRoot + FileName())
	// and returns an ExtractedPackagesResult containing the package type and scoped packages, or nil.
	// It should return nil error if the file simply doesn't exist or contains no relevant packages.
	// It returns a non-nil error for parsing issues or other failures.
	Extract(projectRoot string) (*ExtractedPackagesResult, error)
}

// Package holds common dependency information extracted by various extractors.
// This is an internal representation used by extractors. The main package
// might have its own struct for XML generation.
type Package struct {
    Name    string
    Version string
}

// availableExtractors holds all known metadata extractors.
// Populated by init() functions in other files within this package.
var availableExtractors = make(map[string]MetadataExtractor)

// Register adds an extractor to the registry. Typically called from init().
func Register(name string, extractor MetadataExtractor) {
    if _, exists := availableExtractors[name]; exists {
        // Handle duplicate registration if necessary (e.g., log warning)
        return
    }
    availableExtractors[name] = extractor
}

// GetAllExtractors returns a slice of all registered extractors.
func GetAllExtractors() []MetadataExtractor {
    list := make([]MetadataExtractor, 0, len(availableExtractors))
    for _, extractor := range availableExtractors {
        list = append(list, extractor)
    }
    // Optional: Sort extractors by filename for consistent checking order
    // sort.Slice(list, func(i, j int) bool { return list[i].FileName() < list[j].FileName() })
    return list
}

// mapToInternalPackages converts map[string]string to []Package.
// This is a helper function for use within the extractors package.
func mapToInternalPackages(depMap map[string]string) []Package {
	if len(depMap) == 0 {
		return nil
	}
	packages := make([]Package, 0, len(depMap))
	for name, version := range depMap {
		packages = append(packages, Package{Name: name, Version: version})
	}
	return packages
}
