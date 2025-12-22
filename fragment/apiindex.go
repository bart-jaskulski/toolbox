package main

import "encoding/xml"

// ApiIndex represents a language-agnostic API index.
type ApiIndex struct {
	XMLName xml.Name  `xml:"api" json:"-"`
	Files   []ApiFile `xml:"file" json:"files"`
}

// ApiFile holds the API symbols for a single source file.
type ApiFile struct {
	Path     string      `xml:"path,attr" json:"path"`
	Language string      `xml:"language,attr" json:"language"`
	Symbols  []ApiSymbol `xml:"symbol" json:"symbols"`
}

// ApiSymbol represents a single API symbol.
type ApiSymbol struct {
	Name      string      `xml:"name,attr" json:"name"`
	Kind      string      `xml:"kind,attr" json:"kind"`
	Signature string      `xml:"signature,omitempty" json:"signature,omitempty"`
	Doc       string      `xml:"doc,omitempty" json:"doc,omitempty"`
	Children  []ApiSymbol `xml:"symbol,omitempty" json:"children,omitempty"`
}
