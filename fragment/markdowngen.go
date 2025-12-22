package main

import (
	"fmt"
	"io"
	"strings"
)

type markdownFormatter struct{}

func (markdownFormatter) Name() string {
	return outputFormatMarkdown
}

func (markdownFormatter) Write(w io.Writer, snapshot *ProjectSnapshot) error {
	if _, err := io.WriteString(w, "# Project: "+snapshot.Name+"\n\n"); err != nil {
		return fmt.Errorf("failed to write markdown header: %w", err)
	}

	if len(snapshot.Packages) > 0 {
		if _, err := io.WriteString(w, "## Packages\n"); err != nil {
			return fmt.Errorf("failed to write packages header: %w", err)
		}
		if _, err := io.WriteString(w, "| type | scope | name | version |\n| --- | --- | --- | --- |\n"); err != nil {
			return fmt.Errorf("failed to write packages table header: %w", err)
		}
		for _, pkg := range snapshot.Packages {
			scope := pkg.Scope
			if scope == "" {
				scope = "-"
			}
			if _, err := fmt.Fprintf(
				w,
				"| %s | %s | %s | %s |\n",
				escapeMarkdownTable(pkg.Type),
				escapeMarkdownTable(scope),
				escapeMarkdownTable(pkg.Name),
				escapeMarkdownTable(pkg.Version),
			); err != nil {
				return fmt.Errorf("failed to write packages table: %w", err)
			}
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			return fmt.Errorf("failed to write packages separator: %w", err)
		}
	}

	if snapshot.Tree != nil && len(snapshot.Tree.Children) > 0 {
		if _, err := io.WriteString(w, "## Tree\n"); err != nil {
			return fmt.Errorf("failed to write tree header: %w", err)
		}
		if _, err := io.WriteString(w, "```\n"); err != nil {
			return fmt.Errorf("failed to open tree block: %w", err)
		}
		if err := writeMarkdownTree(w, snapshot.Tree.Children, "", true); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "```\n"); err != nil {
			return fmt.Errorf("failed to close tree block: %w", err)
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			return fmt.Errorf("failed to write tree separator: %w", err)
		}
	}

	if _, err := io.WriteString(w, "## Files\n\n"); err != nil {
		return fmt.Errorf("failed to write files header: %w", err)
	}

	for _, entry := range snapshot.Files {
		if _, err := io.WriteString(w, "### "+markdownInlineCode(entry.Path)+"\n"); err != nil {
			return fmt.Errorf("failed to write file header: %w", err)
		}

		fence := markdownFence(entry.Content)
		if _, err := io.WriteString(w, fence+"\n"); err != nil {
			return fmt.Errorf("failed to write code fence: %w", err)
		}

		if len(entry.Content) > 0 {
			if _, err := w.Write(entry.Content); err != nil {
				return fmt.Errorf("failed to write file contents: %w", err)
			}
		}
		if len(entry.Content) == 0 || entry.Content[len(entry.Content)-1] != '\n' {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return fmt.Errorf("failed to write content newline: %w", err)
			}
		}

		if _, err := io.WriteString(w, fence+"\n\n"); err != nil {
			return fmt.Errorf("failed to close code fence: %w", err)
		}
	}

	return nil
}

func escapeMarkdownTable(value string) string {
	escaped := strings.ReplaceAll(value, "|", "\\|")
	escaped = strings.ReplaceAll(escaped, "\r", " ")
	escaped = strings.ReplaceAll(escaped, "\n", " ")
	return escaped
}

func markdownInlineCode(value string) string {
	maxRun := 0
	run := 0
	for _, r := range value {
		if r == '`' {
			run++
			if run > maxRun {
				maxRun = run
			}
		} else {
			run = 0
		}
	}
	fenceLen := 1
	if maxRun+1 > fenceLen {
		fenceLen = maxRun + 1
	}
	fence := strings.Repeat("`", fenceLen)

	padLeft := strings.HasPrefix(value, "`") || strings.HasPrefix(value, " ")
	padRight := strings.HasSuffix(value, "`") || strings.HasSuffix(value, " ")
	leftPad := ""
	rightPad := ""
	if padLeft {
		leftPad = " "
	}
	if padRight {
		rightPad = " "
	}

	return fence + leftPad + value + rightPad + fence
}

func markdownFence(content []byte) string {
	maxRun := 0
	run := 0
	for _, b := range content {
		if b == '`' {
			run++
			if run > maxRun {
				maxRun = run
			}
		} else {
			run = 0
		}
	}
	fenceLen := 3
	if maxRun+1 > fenceLen {
		fenceLen = maxRun + 1
	}
	return strings.Repeat("`", fenceLen)
}

func writeMarkdownTree(w io.Writer, nodes []*DirNode, prefix string, isRoot bool) error {
	for i, node := range nodes {
		if node == nil {
			continue
		}
		isLast := i == len(nodes)-1
		branch := "|-- "
		nextPrefix := prefix + "|   "
		if isLast {
			branch = "`-- "
			nextPrefix = prefix + "    "
		}
		if isRoot {
			branch = ""
			nextPrefix = ""
		}

		suffix := ""
		if node.IsDir {
			suffix = "/"
		}
		if _, err := io.WriteString(w, prefix+branch+node.Name+suffix+"\n"); err != nil {
			return fmt.Errorf("failed to write tree node: %w", err)
		}
		if len(node.Children) > 0 {
			if err := writeMarkdownTree(w, node.Children, nextPrefix, false); err != nil {
				return err
			}
		}
	}
	return nil
}
