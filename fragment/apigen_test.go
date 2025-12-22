package main

import "testing"

func TestTypeScriptApiParser(t *testing.T) {
	parser := newTypeScriptParser()
	src := []byte(`/**
 * doc for foo
 */
export function foo(bar) { return bar }

export class Baz {
  /** doc for method */
  method(x) { return x }
}`)

	apiFile, err := parser.Parse("file.ts", src)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if apiFile == nil || len(apiFile.Symbols) == 0 {
		t.Fatalf("expected symbols from typescript parser")
	}
	found := false
	for _, symbol := range apiFile.Symbols {
		if symbol.Name == "foo" && symbol.Kind == "function" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected function foo in symbols")
	}
}

func TestPhpApiParser(t *testing.T) {
	parser := newPhpParser()
	src := []byte(`<?php
/** doc for foo */
function foo($bar) { return $bar; }

class Baz {
  /** doc for method */
  public function method($x) { return $x; }
}`)

	apiFile, err := parser.Parse("file.php", src)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if apiFile == nil || len(apiFile.Symbols) == 0 {
		t.Fatalf("expected symbols from php parser")
	}
	found := false
	for _, symbol := range apiFile.Symbols {
		if symbol.Name == "foo" && symbol.Kind == "function" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected function foo in symbols")
	}
}
