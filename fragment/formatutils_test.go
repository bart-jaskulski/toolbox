package main

import "testing"

func TestNormalizeOutputFormat(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		want      string
		expectErr bool
	}{
		{name: "default", raw: "", want: outputFormatXML},
		{name: "xml", raw: "xml", want: outputFormatXML},
		{name: "markdown", raw: "markdown", want: outputFormatMarkdown},
		{name: "md alias", raw: "md", want: outputFormatMarkdown},
		{name: "json", raw: "json", want: outputFormatJSON},
		{name: "unknown", raw: "yaml", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeOutputFormat(tt.raw)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error for %q", tt.raw)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestInferFormatFromOutput(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		want     string
		expectOK bool
	}{
		{name: "xml", path: "output.xml", want: outputFormatXML, expectOK: true},
		{name: "md", path: "output.md", want: outputFormatMarkdown, expectOK: true},
		{name: "markdown", path: "output.markdown", want: outputFormatMarkdown, expectOK: true},
		{name: "json", path: "output.json", want: outputFormatJSON, expectOK: true},
		{name: "unknown", path: "output.txt", expectOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := inferFormatFromOutput(tt.path)
			if ok != tt.expectOK {
				t.Fatalf("expected ok=%v, got %v", tt.expectOK, ok)
			}
			if ok && got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
