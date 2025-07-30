package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var outputFile string
	var showHelp bool

	flag.StringVar(&outputFile, "o", "", "Output file (default: input.mdidx for files, stdout for stdin)")
	flag.StringVar(&outputFile, "output", "", "Output file (default: input.mdidx for files, stdout for stdin)")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
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
  -o, --output  Output file (default: input.mdidx for files, stdout for stdin)
  --help        Show this help message

Examples:
  mdidx document.md                    # Creates document.mdidx
  mdidx -o index.mdidx document.md     # Creates index.mdidx
  mdidx < document.md > index.mdidx    # Read from stdin, write to stdout
`)
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: mdidx [options] [file]\n")
	fmt.Fprintf(os.Stderr, "Run 'mdidx --help' for more information.\n")
}
