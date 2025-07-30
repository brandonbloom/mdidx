package main

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Header struct {
	Level     int
	Title     string
	StartLine int
	EndLine   int
}

type Parser struct {
	source   []byte
	lines    []string
	headers  []Header
	warnings []string
}

func NewParser(source []byte) *Parser {
	lines := strings.Split(string(source), "\n")
	return &Parser{
		source: source,
		lines:  lines,
	}
}

func (p *Parser) Parse() ([]Header, []string, error) {
	md := goldmark.New()
	reader := text.NewReader(p.source)
	doc := md.Parser().Parse(reader)

	p.extractHeaders(doc, reader)
	p.calculateEndLines()

	return p.headers, p.warnings, nil
}

func (p *Parser) extractHeaders(node ast.Node, reader text.Reader) {
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if heading, ok := n.(*ast.Heading); ok {
			lines := heading.Lines()
			var startPos int
			if lines.Len() > 0 {
				startPos = lines.At(0).Start
			}
			startLine := p.getLineNumber(startPos)
			title := p.extractHeadingText(heading, reader)

			header := Header{
				Level:     heading.Level,
				Title:     title,
				StartLine: startLine,
			}

			p.headers = append(p.headers, header)
		}

		return ast.WalkContinue, nil
	})
}

func (p *Parser) extractHeadingText(heading *ast.Heading, reader text.Reader) string {
	text := heading.Text(reader.Source())
	// Replace line breaks with spaces as specified
	return strings.ReplaceAll(strings.TrimSpace(string(text)), "\n", " ")
}

func (p *Parser) getLineNumber(pos int) int {
	lineStart := 0
	for lineNum := 0; lineNum < len(p.lines); lineNum++ {
		lineEnd := lineStart + len(p.lines[lineNum])
		if pos >= lineStart && pos <= lineEnd {
			return lineNum + 1 // 1-based line numbering
		}
		lineStart = lineEnd + 1 // +1 for the newline character
	}
	return 1 // fallback
}

func (p *Parser) calculateEndLines() {
	for i := range p.headers {
		p.headers[i].EndLine = p.findEndLine(i)
	}
}

func (p *Parser) findEndLine(headerIndex int) int {
	currentHeader := p.headers[headerIndex]

	// Find the next header at the same level or higher (lower number)
	nextBoundaryLine := len(p.lines) // Default to end of file

	for j := headerIndex + 1; j < len(p.headers); j++ {
		if p.headers[j].Level <= currentHeader.Level {
			nextBoundaryLine = p.headers[j].StartLine - 1
			break
		}
	}

	// Trim trailing empty lines from the section range
	for nextBoundaryLine > currentHeader.StartLine && strings.TrimSpace(p.lines[nextBoundaryLine-1]) == "" {
		nextBoundaryLine--
	}

	// Ensure we don't go before the header itself
	if nextBoundaryLine < currentHeader.StartLine {
		nextBoundaryLine = currentHeader.StartLine
	}

	return nextBoundaryLine
}
