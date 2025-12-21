package main

// ProjectPackage represents a single dependency entry.
// XML tags are used for XML output; JSON uses separate tags in the formatter.
type ProjectPackage struct {
	Type    string `xml:"type,attr" json:"type"`
	Scope   string `xml:"scope,attr,omitempty" json:"scope,omitempty"`
	Name    string `xml:"name,attr" json:"name"`
	Version string `xml:"version,attr" json:"version"`
}

// FileEntry holds a file's content and relative path.
type FileEntry struct {
	Path    string
	Content []byte
}
