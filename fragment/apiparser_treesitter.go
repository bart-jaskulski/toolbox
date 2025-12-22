package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tsjavascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	tsphp "github.com/tree-sitter/tree-sitter-php/bindings/go"
	tstypescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

type tsParser struct {
	name  string
	langs map[string]*tree_sitter.Language
	exts  []string
}

func newJavaScriptParser() apiParser {
	return &tsParser{
		name: "javascript",
		langs: map[string]*tree_sitter.Language{
			".js":  tree_sitter.NewLanguage(tsjavascript.Language()),
			".jsx": tree_sitter.NewLanguage(tsjavascript.Language()),
		},
		exts: []string{".js", ".jsx"},
	}
}

func newTypeScriptParser() apiParser {
	return &tsParser{
		name: "typescript",
		langs: map[string]*tree_sitter.Language{
			".ts":  tree_sitter.NewLanguage(tstypescript.LanguageTypescript()),
			".tsx": tree_sitter.NewLanguage(tstypescript.LanguageTSX()),
		},
		exts: []string{".ts", ".tsx"},
	}
}

func newPhpParser() apiParser {
	return &tsParser{
		name: "php",
		langs: map[string]*tree_sitter.Language{
			".php": tree_sitter.NewLanguage(tsphp.LanguagePHP()),
		},
		exts: []string{".php"},
	}
}

func (p *tsParser) Language() string {
	return p.name
}

func (p *tsParser) Extensions() []string {
	return p.exts
}

func (p *tsParser) Parse(path string, src []byte) (*ApiFile, error) {
	ext := strings.ToLower(filepath.Ext(path))
	lang := p.langs[ext]
	if lang == nil {
		for _, candidate := range p.langs {
			lang = candidate
			break
		}
	}
	if lang == nil {
		return nil, fmt.Errorf("no language available for %s", path)
	}
	parser := tree_sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(lang); err != nil {
		return nil, fmt.Errorf("set language: %w", err)
	}
	tree := parser.Parse(src, nil)
	defer tree.Close()
	root := tree.RootNode()

	var symbols []ApiSymbol
	switch p.name {
	case "php":
		symbols = collectPHPSymbols(root, src)
	default:
		symbols = collectJSSymbols(root, src)
	}

	return &ApiFile{
		Path:     path,
		Language: p.name,
		Symbols:  symbols,
	}, nil
}

func collectJSSymbols(root *tree_sitter.Node, src []byte) []ApiSymbol {
	var symbols []ApiSymbol
	for i := uint(0); i < root.NamedChildCount(); i++ {
		child := root.NamedChild(i)
		if child == nil {
			continue
		}
		symbols = append(symbols, symbolsFromJSNode(child, src)...)
	}
	return symbols
}

func symbolsFromJSNode(node *tree_sitter.Node, src []byte) []ApiSymbol {
	if node == nil {
		return nil
	}
	switch node.Kind() {
	case "export_statement":
		return symbolsFromJSExport(node, src)
	case "function_declaration":
		return []ApiSymbol{symbolFromNamedNode("function", node, src)}
	case "class_declaration":
		return []ApiSymbol{classSymbolFromJS(node, src)}
	case "interface_declaration":
		return []ApiSymbol{symbolFromNamedNode("interface", node, src)}
	case "type_alias_declaration":
		return []ApiSymbol{symbolFromNamedNode("type", node, src)}
	case "enum_declaration":
		return []ApiSymbol{symbolFromNamedNode("enum", node, src)}
	case "lexical_declaration", "variable_declaration":
		return variableSymbolsFromJS(node, src)
	default:
		return nil
	}
}

func symbolsFromJSExport(node *tree_sitter.Node, src []byte) []ApiSymbol {
	for i := uint(0); i < node.NamedChildCount(); i++ {
		child := node.NamedChild(i)
		if child == nil {
			continue
		}
		if symbols := symbolsFromJSNode(child, src); len(symbols) > 0 {
			return symbols
		}
	}
	return nil
}

func classSymbolFromJS(node *tree_sitter.Node, src []byte) ApiSymbol {
	symbol := symbolFromNamedNode("class", node, src)
	body := node.ChildByFieldName("body")
	if body == nil {
		return symbol
	}
	for i := uint(0); i < body.NamedChildCount(); i++ {
		child := body.NamedChild(i)
		if child == nil {
			continue
		}
		switch child.Kind() {
		case "method_definition":
			symbol.Children = append(symbol.Children, symbolFromNamedNode("method", child, src))
		case "public_field_definition", "property_definition":
			symbol.Children = append(symbol.Children, symbolFromNamedNode("property", child, src))
		}
	}
	return symbol
}

func variableSymbolsFromJS(node *tree_sitter.Node, src []byte) []ApiSymbol {
	prefix := prefixFromDeclaration(src, node.StartByte(), node.EndByte())
	doc := sanitizeDoc(extractDocComment(src, node.StartByte()))
	var symbols []ApiSymbol
	for i := uint(0); i < node.NamedChildCount(); i++ {
		child := node.NamedChild(i)
		if child == nil || child.Kind() != "variable_declarator" {
			continue
		}
		name := nodeFieldText(child, src, "name")
		signature := variableSignature(prefix, name)
		symbols = append(symbols, ApiSymbol{
			Name:      name,
			Kind:      "variable",
			Signature: signature,
			Doc:       doc,
		})
	}
	return symbols
}

func collectPHPSymbols(root *tree_sitter.Node, src []byte) []ApiSymbol {
	var symbols []ApiSymbol
	for i := uint(0); i < root.NamedChildCount(); i++ {
		child := root.NamedChild(i)
		if child == nil {
			continue
		}
		switch child.Kind() {
		case "function_definition":
			symbols = append(symbols, symbolFromNamedNode("function", child, src))
		case "class_declaration":
			symbols = append(symbols, classSymbolFromPHP(child, src))
		case "interface_declaration":
			symbols = append(symbols, symbolFromNamedNode("interface", child, src))
		case "trait_declaration":
			symbols = append(symbols, symbolFromNamedNode("trait", child, src))
		}
	}
	return symbols
}

func classSymbolFromPHP(node *tree_sitter.Node, src []byte) ApiSymbol {
	symbol := symbolFromNamedNode("class", node, src)
	body := node.ChildByFieldName("body")
	if body == nil {
		return symbol
	}
	for i := uint(0); i < body.NamedChildCount(); i++ {
		child := body.NamedChild(i)
		if child == nil {
			continue
		}
		switch child.Kind() {
		case "method_declaration":
			symbol.Children = append(symbol.Children, symbolFromNamedNode("method", child, src))
		case "property_declaration":
			symbol.Children = append(symbol.Children, phpPropertySymbols(child, src)...)
		}
	}
	return symbol
}

func phpPropertySymbols(node *tree_sitter.Node, src []byte) []ApiSymbol {
	var symbols []ApiSymbol
	doc := sanitizeDoc(extractDocComment(src, node.StartByte()))
	for i := uint(0); i < node.NamedChildCount(); i++ {
		child := node.NamedChild(i)
		if child == nil {
			continue
		}
		if child.Kind() == "property_element" {
			name := nodeFieldText(child, src, "name")
			signature := variableSignature("property", name)
			symbols = append(symbols, ApiSymbol{
				Name:      name,
				Kind:      "property",
				Signature: signature,
				Doc:       doc,
			})
		}
	}
	return symbols
}

func symbolFromNamedNode(kind string, node *tree_sitter.Node, src []byte) ApiSymbol {
	name := nodeFieldText(node, src, "name")
	doc := sanitizeDoc(extractDocComment(src, node.StartByte()))
	signature := signatureFromNode(src, node.StartByte(), node.EndByte())
	signature = truncateSignature(signature, 200)
	if name == "" {
		name = debugSignature(kind, name, signature)
	}
	return ApiSymbol{
		Name:      name,
		Kind:      kind,
		Signature: signature,
		Doc:       doc,
	}
}

func nodeFieldText(node *tree_sitter.Node, src []byte, field string) string {
	child := node.ChildByFieldName(field)
	if child == nil {
		return ""
	}
	return strings.TrimSpace(string(src[child.StartByte():child.EndByte()]))
}

func extractDocComment(src []byte, start uint) string {
	if start == 0 {
		return ""
	}
	idx := int(start)
	if idx > len(src) {
		idx = len(src)
	}
	for idx > 0 && isSpace(src[idx-1]) {
		idx--
	}
	if idx >= 2 && src[idx-2] == '*' && src[idx-1] == '/' {
		commentStart := bytes.LastIndex(src[:idx-1], []byte("/*"))
		if commentStart >= 0 && commentStart+2 < idx-1 && src[commentStart+2] == '*' {
			return sanitizeDoc(cleanDocBlock(src[commentStart+3 : idx-2]))
		}
	}
	return sanitizeDoc(cleanLineDoc(src[:idx]))
}

func cleanDocBlock(block []byte) string {
	lines := strings.Split(string(block), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)
		lines[i] = line
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func cleanLineDoc(src []byte) string {
	lines := strings.Split(string(src), "\n")
	var docLines []string
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			if len(docLines) == 0 {
				continue
			}
			break
		}
		if !strings.HasPrefix(line, "//") {
			break
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "//"))
		docLines = append([]string{line}, docLines...)
	}
	return strings.TrimSpace(strings.Join(docLines, "\n"))
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}
