// metadata.go
package main

import (
	"log"

	"github.com/bart-jaskulski/toolbox/fragment/internal/extractors"
)

// extractPackageData iterates through registered extractors and aggregates package information
// into a unified list of ProjectPackage structs.
func extractPackageData(cfg *Config) []ProjectPackage {
	log.Println("Attempting to extract package information using registered extractors...")
	var allPkgs []ProjectPackage
	foundAnyPackages := false

	// Get all registered extractors
	allExtractors := extractors.GetAllExtractors()

	for _, extractor := range allExtractors {
		fileName := extractor.FileName()
		log.Printf("  Checking for: %s", fileName)

		// Extract now returns *ExtractedPackagesResult or nil
		extractedData, err := extractor.Extract(cfg.ProjectRoot)

		if err != nil {
			log.Printf("Warning: Failed to process %s: %v", fileName, err)
			continue
		}

		if extractedData == nil { // File not found, or no packages in file
			continue
		}

		pkgType := extractedData.PackageType
		foundPackagesInThisFile := false
		for _, group := range extractedData.Groups {
			if len(group.Packages) > 0 {
				pkgs := mapInternalToProjectPackages(group.Packages, pkgType, group.ScopeName)
				allPkgs = append(allPkgs, pkgs...)
				foundPackagesInThisFile = true
			}
		}

		if foundPackagesInThisFile {
			log.Printf("    Found and parsed %s", fileName)
			foundAnyPackages = true
		} else {
			// This case might occur if an extractor successfully parses a file
			// but determines there are no relevant packages to report (e.g. all groups are empty).
			// log.Printf("    Parsed %s, but no relevant packages found.", fileName) // Optional logging
		}
	}

	if !foundAnyPackages {
		log.Println("  No supported package metadata files found or successfully parsed at project root.")
		return nil // Return nil if no packages were found
	}
	// Optional: Sort allPkgs here by Type, Scope, Name if desired for consistent output
	return allPkgs
}

// mapInternalToProjectPackages converts []extractors.Package to []main.ProjectPackage,
// assigning the given package type and scope.
func mapInternalToProjectPackages(internalPkgs []extractors.Package, pkgType string, pkgScope string) []ProjectPackage {
	if len(internalPkgs) == 0 {
		return nil
	}
	projectPkgs := make([]ProjectPackage, len(internalPkgs))
	for i, pkg := range internalPkgs {
		projectPkgs[i] = ProjectPackage{
			Type:    pkgType,
			Scope:   pkgScope,
			Name:    pkg.Name,
			Version: pkg.Version,
		}
	}
	// Optional: Sort xmlPkgs by Name here if needed and not sorted earlier
	return projectPkgs
}
