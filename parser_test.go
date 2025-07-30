package main

import (
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Header
	}{
		{
			name: "simple headers",
			input: `# Title

Some content here

## Subtitle

Some other content here.

# Footer

The end`,
			expected: []Header{
				{Level: 1, Title: "Title", StartLine: 1, EndLine: 7},
				{Level: 2, Title: "Subtitle", StartLine: 5, EndLine: 7},
				{Level: 1, Title: "Footer", StartLine: 9, EndLine: 11},
			},
		},
		{
			name: "nested headers",
			input: `# Main

Content

## Sub1

Content

### SubSub1

Content

### SubSub2

Content

## Sub2

Content

# Another Main

Content`,
			expected: []Header{
				{Level: 1, Title: "Main", StartLine: 1, EndLine: 19},
				{Level: 2, Title: "Sub1", StartLine: 5, EndLine: 15},
				{Level: 3, Title: "SubSub1", StartLine: 9, EndLine: 11},
				{Level: 3, Title: "SubSub2", StartLine: 13, EndLine: 15},
				{Level: 2, Title: "Sub2", StartLine: 17, EndLine: 19},
				{Level: 1, Title: "Another Main", StartLine: 21, EndLine: 23},
			},
		},
		{
			name: "single line header only",
			input: `# This is a single line header

Some content`,
			expected: []Header{
				{Level: 1, Title: "This is a single line header", StartLine: 1, EndLine: 3},
			},
		},
		{
			name: "empty sections",
			input: `# Header1

## Header2

# Header3`,
			expected: []Header{
				{Level: 1, Title: "Header1", StartLine: 1, EndLine: 3},
				{Level: 2, Title: "Header2", StartLine: 3, EndLine: 3},
				{Level: 1, Title: "Header3", StartLine: 5, EndLine: 5},
			},
		},
		{
			name: "headers with trailing empty lines",
			input: `# Header1

Content here


## Header2

More content


`,
			expected: []Header{
				{Level: 1, Title: "Header1", StartLine: 1, EndLine: 8},
				{Level: 2, Title: "Header2", StartLine: 6, EndLine: 8},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser([]byte(tt.input))
			headers, warnings, err := parser.Parse()

			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}

			if len(warnings) > 0 {
				t.Errorf("Parse() warnings = %v", warnings)
			}

			if !reflect.DeepEqual(headers, tt.expected) {
				t.Errorf("Parse() headers mismatch\nGot:\n")
				for i, h := range headers {
					t.Errorf("  [%d] Level:%d, Title:%q, StartLine:%d, EndLine:%d", i, h.Level, h.Title, h.StartLine, h.EndLine)
				}
				t.Errorf("Expected:\n")
				for i, h := range tt.expected {
					t.Errorf("  [%d] Level:%d, Title:%q, StartLine:%d, EndLine:%d", i, h.Level, h.Title, h.StartLine, h.EndLine)
				}
			}
		})
	}
}

func TestIndexGenerator_Generate(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		headers  []Header
		expected string
	}{
		{
			name:   "simple index with source",
			source: "./original-big-markdown-file.md",
			headers: []Header{
				{Level: 1, Title: "Title", StartLine: 1, EndLine: 8},
				{Level: 2, Title: "Subtitle", StartLine: 5, EndLine: 8},
				{Level: 1, Title: "Footer", StartLine: 9, EndLine: 11},
			},
			expected: `source: ./original-big-markdown-file.md
---
1-8: Title
  5-8: Subtitle
9-11: Footer
`,
		},
		{
			name:   "index without source",
			source: "",
			headers: []Header{
				{Level: 1, Title: "Title", StartLine: 1, EndLine: 3},
				{Level: 2, Title: "Subtitle", StartLine: 2, EndLine: 3},
			},
			expected: `---
1-3: Title
  2-3: Subtitle
`,
		},
		{
			name:   "deeply nested headers",
			source: "test.md",
			headers: []Header{
				{Level: 1, Title: "H1", StartLine: 1, EndLine: 10},
				{Level: 2, Title: "H2", StartLine: 2, EndLine: 10},
				{Level: 3, Title: "H3", StartLine: 3, EndLine: 10},
				{Level: 4, Title: "H4", StartLine: 4, EndLine: 10},
			},
			expected: `source: test.md
---
1-10: H1
  2-10: H2
    3-10: H3
      4-10: H4
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := NewIndexGenerator(tt.source, tt.headers, nil)
			result := generator.Generate()

			if result != tt.expected {
				t.Errorf("Generate() mismatch\nGot:\n%s\nExpected:\n%s", result, tt.expected)
			}
		})
	}
}
