// xmltypes.go
package main

import "encoding/xml"

// Project represents the root element of the XML output.
type Project struct {
	XMLName  xml.Name         `xml:"project"`
	Name     string           `xml:"name,attr"`
	Metadata *ProjectMetadata `xml:"metadata,omitempty"` // For non-package metadata
	Files    *Files           `xml:"files,omitempty"`
	Packages *PackagesHolder  `xml:"packages,omitempty"` // Unified packages list
	Tree     *Tree            `xml:"tree,omitempty"`
	API      *ApiIndex        `xml:"api,omitempty"`
}

// PackagesHolder wraps the list of packages to ensure the <packages> element.
type PackagesHolder struct {
	XMLName     xml.Name         `xml:"packages"`
	PackageList []ProjectPackage `xml:"package"`
}

// ProjectMetadata holds all non-package-specific extracted metadata for XML output.
type ProjectMetadata struct {
	XMLName xml.Name `xml:"metadata" json:"-"`
	// This struct can hold other types of metadata in the future,
	// not related to package dependencies.
}

// Files represents the container for multiple file elements. (Remains the same)
type Files struct {
	XMLName xml.Name `xml:"files"`
	Files   []File   `xml:"file"`
}

// File represents a single file entry in the XML. (Remains the same)
type File struct {
	XMLName  xml.Name `xml:"file"`
	Path     string   `xml:"path,attr"`
	Language string   `xml:"language,attr,omitempty"`
	Content  []byte   `xml:",cdata"`
}

// Tree represents a directory tree overview.
type Tree struct {
	XMLName xml.Name   `xml:"tree"`
	Nodes   []TreeNode `xml:"node"`
}

// TreeNode represents a directory or file node in the tree.
type TreeNode struct {
	XMLName  xml.Name   `xml:"node"`
	Name     string     `xml:"name,attr"`
	Path     string     `xml:"path,attr,omitempty"`
	Kind     string     `xml:"kind,attr"`
	Children []TreeNode `xml:"node,omitempty"`
}
