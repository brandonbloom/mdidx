package main

import (
	"fmt"
	"strings"
)

type IndexGenerator struct {
	source   string
	comment  string
	headers  []Header
	warnings []string
}

func NewIndexGenerator(source string, headers []Header, warnings []string) *IndexGenerator {
	return &IndexGenerator{
		source:   source,
		headers:  headers,
		warnings: warnings,
	}
}

func (g *IndexGenerator) SetComment(comment string) {
	g.comment = comment
}

func (g *IndexGenerator) Generate() string {
	var output strings.Builder

	// Generate YAML frontmatter
	g.generateYAMLFrontmatter(&output)

	// Generate separator
	output.WriteString("---\n")

	// Generate table of contents
	for _, header := range g.headers {
		indent := strings.Repeat("  ", header.Level-1) // 2 spaces per level beyond 1
		line := fmt.Sprintf("%s%d-%d: %s\n", indent, header.StartLine, header.EndLine, header.Title)
		output.WriteString(line)
	}

	return output.String()
}

func (g *IndexGenerator) generateYAMLFrontmatter(output *strings.Builder) {
	hasMetadata := g.source != "" || g.comment != ""

	if !hasMetadata {
		return
	}

	if g.source != "" {
		output.WriteString(fmt.Sprintf("source: %s\n", g.source))
	}

	if g.comment != "" {
		output.WriteString("comment: >\n")
		// Indent the comment content by 2 spaces
		commentLines := strings.Split(strings.TrimSpace(g.comment), "\n")
		for _, line := range commentLines {
			output.WriteString(fmt.Sprintf("  %s\n", line))
		}
	}
}
