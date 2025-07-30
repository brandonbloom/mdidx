# mdidx

A tool that builds an index from a markdown file.

## Purpose

The purpose is to take a large markdown file and make a smaller one from that. 
The primary use case is to minimize the amount of context used when loading a 
large corpus of context into an LLM.

## How it works

For this first version, the only mechanism to do that is to extract all of the 
headers and build a full table of contents. Each entry in the table of contents 
is annotated with the section's line range.

## Example

**Input:**

```markdown
# Title

Some content here

## Subtitle

Some other content here.

# Footer

The end
```

**Output:**

```
source: ./original-big-markdown-file.md
---
1-8: Title
  5-8: Subtitle
9-11: Footer
```

Where each line (after the ---) has the format:

`{indent}{start}-{end}: {title}`

- Indent is two spaces for each header level beyond 1.
- Start is the first line that the header appears on.
- End is the last line before the next header of the same or lower level.
- Title is the text content of the header with any line breaks replaced by spaces.

Before the "---" is yaml content, and can be used for metadata.

## Usage

```bash
mdidx ./original-big-markdown-file.md
```

Which is the same as:

```bash
mdidx < ./original-big-markdown-file.md > ./original-big-markdown-file.mdidx
```

The tool can also be used as:

```bash
mdidx -o output.mdidx ./original-big-markdown-file.md
```

Where `-o` is short for `--output`.

## Installation

### From source

```bash
go install github.com/brandonbloom/mdidx@latest
```

### Local build

```bash
go build -o mdidx
```

## License

MIT License - see [LICENSE](LICENSE) file for details.