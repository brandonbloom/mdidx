package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const defaultComment = `This is a markdown index file created by mdidx. It contains a table of contents with line ranges from the source file.

Format: {indent}{start}-{end}: {title}
- Each line shows the header title and its line range in the source
- Indent is 2 spaces per header level (## = 2 spaces, ### = 4 spaces, etc.)
- Line ranges are inclusive and 1-based
- To find content for a section, look at those line numbers in the source file

Use this index to understand document structure and locate specific content without reading the entire file.`

func main() {
	var outputFile string
	var showHelp bool
	var addPreamble bool
	var noPreamble bool

	flag.StringVar(&outputFile, "o", "", "Output file (default: input.mdidx for files, stdout for stdin)")
	flag.StringVar(&outputFile, "output", "", "Output file (default: input.mdidx for files, stdout for stdin)")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&addPreamble, "preamble", false, "Add explanatory preamble to help LLM agents understand the format")
	flag.BoolVar(&noPreamble, "no-preamble", false, "Disable the default preamble (default preamble is added for file input, not for stdin)")
	flag.Parse()

	if showHelp {
		showHelpMessage()
		return
	}

	args := flag.Args()
	var input io.Reader
	var sourcePath string
	var defaultOutput string

	if len(args) == 0 {
		// Read from stdin
		input = os.Stdin
		sourcePath = ""
		defaultOutput = "" // stdout
	} else if len(args) == 1 {
		// Read from file
		filename := args[0]
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filename, err)
			os.Exit(1)
		}
		defer file.Close()

		input = file
		sourcePath = filename
		defaultOutput = filename + ".mdidx"
		if strings.HasSuffix(filename, ".md") {
			defaultOutput = strings.TrimSuffix(filename, ".md") + ".mdidx"
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: too many arguments\n")
		showUsage()
		os.Exit(1)
	}

	// Read input
	content, err := io.ReadAll(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Parse markdown
	parser := NewParser(content)
	headers, warnings, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing markdown: %v\n", err)
		os.Exit(1)
	}

	// Print warnings to stderr
	exitCode := 0
	for _, warning := range warnings {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
		exitCode = 1
	}

	// Generate index
	generator := NewIndexGenerator(sourcePath, headers, warnings)

	// Determine if preamble should be added
	shouldAddPreamble := false
	if addPreamble {
		shouldAddPreamble = true
	} else if !noPreamble {
		// Default behavior: add preamble for file input, not for stdin
		shouldAddPreamble = (sourcePath != "")
	}

	if shouldAddPreamble {
		generator.SetComment(defaultComment)
	}
	output := generator.Generate()

	// Write output
	if outputFile != "" {
		// Write to specified file
		err := os.WriteFile(outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file %s: %v\n", outputFile, err)
			os.Exit(1)
		}
	} else if defaultOutput != "" {
		// Write to default file
		err := os.WriteFile(defaultOutput, []byte(output), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file %s: %v\n", defaultOutput, err)
			os.Exit(1)
		}
	} else {
		// Write to stdout
		fmt.Print(output)
	}

	os.Exit(exitCode)
}

func showHelpMessage() {
	fmt.Printf(`mdidx - Build an index from a markdown file

Usage:
  mdidx [options] [file]
  mdidx [options] < input.md > output.mdidx

Arguments:
  file          Input markdown file (if not provided, reads from stdin)

Options:
  -o, --output     Output file (default: input.mdidx for files, stdout for stdin)
  --preamble       Add explanatory preamble to help LLM agents understand the format
  --no-preamble    Disable the default preamble (default: preamble added for file input, not for stdin)
  --help           Show this help message

Examples:
  mdidx document.md                    # Creates document.mdidx with preamble (default for files)
  mdidx --no-preamble document.md      # Creates document.mdidx without preamble
  mdidx -o index.mdidx document.md     # Creates index.mdidx with preamble
  mdidx < document.md > index.mdidx    # Read from stdin, write to stdout (no preamble by default)
  echo "# Test" | mdidx --preamble     # Force preamble for stdin input
`)
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: mdidx [options] [file]\n")
	fmt.Fprintf(os.Stderr, "Run 'mdidx --help' for more information.\n")
}
