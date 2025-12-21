// xmltypes.go
package main

import "encoding/xml"

// Project represents the root element of the XML output.
type Project struct {
    XMLName  xml.Name         `xml:"project"`
    Name     string           `xml:"name,attr"`
    Metadata *ProjectMetadata `xml:"metadata,omitempty"` // For non-package metadata
    Files    Files            `xml:"files"`
    Packages *PackagesHolder  `xml:"packages,omitempty"` // Unified packages list
}

// PackagesHolder wraps the list of packages to ensure the <packages> element.
type PackagesHolder struct {
    XMLName     xml.Name     `xml:"packages"`
    PackageList []XmlPackage `xml:"package"`
}

// ProjectMetadata holds all non-package-specific extracted metadata for XML output.
type ProjectMetadata struct {
    XMLName xml.Name `xml:"metadata"`
    // This struct can hold other types of metadata in the future,
    // not related to package dependencies.
}

// XmlPackage represents a single dependency entry for XML output.
// It now includes type and scope attributes.
type XmlPackage struct {
    XMLName xml.Name `xml:"package"`
    Type    string   `xml:"type,attr"`            // e.g., "composer", "npm"
    Scope   string   `xml:"scope,attr,omitempty"` // e.g., "require", "dependencies", "require-dev"
    Name    string   `xml:"name,attr"`
    Version string   `xml:"version,attr"`
}

// Files represents the container for multiple file elements. (Remains the same)
type Files struct {
    XMLName xml.Name `xml:"files"`
    Files   []File   `xml:"file"`
}

// File represents a single file entry in the XML. (Remains the same)
type File struct {
    XMLName xml.Name `xml:"file"`
    Path    string   `xml:"path,attr"`
    Content []byte   `xml:",cdata"`
}
