# Developing mdidx

This document provides guidance on working with the mdidx project.

## Project Structure

- `parser.go` - Core parsing logic using goldmark AST visitor to extract headers
- `generator.go` - Index generation with proper indentation and YAML metadata
- `main.go` - CLI interface with argument parsing and I/O handling
- `parser_test.go` - Unit tests for parsing and generation logic
- `test.sh` - Integration test suite for CLI functionality

## Requirements

- Go with goldmark dependency (`github.com/yuin/goldmark`)
- Uses goldmark AST parsing (no ad-hoc markdown parsing)
- Before wiring it up to the CLI interface, the parsing and index generation 
  is fully tested with Go unit tests

## Building

```bash
go build -o mdidx
```

## Testing

### Unit Tests

Run the Go unit tests for core parsing and generation logic:

```bash
go test -v
```

### Integration Tests

Run the complete CLI integration test suite:

```bash
./test.sh
```

The integration tests cover:
- Help flag functionality
- Simple file processing with source field and proper indentation
- Custom output file (`-o` flag)
- Stdin/stdout processing
- Error handling for nonexistent files and invalid arguments

## Context Files

The project includes context documentation:

- `context/goldmark-ast.md` - Documentation about the goldmark AST package, 
  including node types, tree traversal, and key interfaces used for parsing

## Implementation Notes

- Uses goldmark's AST walker to visit heading nodes and track their positions
- Calculates line ranges by finding the next header of same/higher level
- Includes intermediate content (subheaders) within each section's range
- Generates YAML frontmatter with source filename when input is a file
- Handles stdin → stdout vs filename → .mdidx file automatically
- Prints warnings to stderr for parsing issues, continues processing, 
  exits non-zero if warnings occurred

## Known TODOs

- Exclude trailing blank lines from section ranges (currently includes them 
  for simplicity)

## Error Handling

The tool implements error handling with warnings to stderr and non-zero exit 
codes. It makes an attempt to process the entire file, recovering from issues. 
If there are any warnings printed, it returns a non-zero error code.