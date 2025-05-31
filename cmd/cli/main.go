package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/francescoalemanno/gotypst"
)

func main() {
	var (
		inputFile    = flag.String("input", "", "Input markdown file (required)")
		outputFile   = flag.String("output", "", "Output PDF file (optional, defaults to input.pdf)")
		templateFile = flag.String("template", "exam-template.typ", "Template file path")
		help         = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	if *help || *inputFile == "" {
		showHelp()
		return
	}

	// Determine output file
	output := *outputFile
	if output == "" {
		ext := filepath.Ext(*inputFile)
		output = strings.TrimSuffix(*inputFile, ext) + ".pdf"
	}

	// Convert
	if err := convertMarkdownToPDF(*inputFile, output, *templateFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ PDF generated successfully: %s\n", output)
}

func showHelp() {
	fmt.Println("Markdown to PDF Converter (CLI)")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  md-pdf-cli -input <markdown-file> [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -input <file>      Input markdown file (required)")
	fmt.Println("  -output <file>     Output PDF file (optional)")
	fmt.Println("  -template <file>   Template file path (default: exam-template.typ)")
	fmt.Println("  -help              Show this help")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  md-pdf-cli -input test.md")
	fmt.Println("  md-pdf-cli -input test.md -output my-exam.pdf")
	fmt.Println("  md-pdf-cli -input test.md -template custom-template.typ")
}

func convertMarkdownToPDF(inputFile, outputFile, templateFile string) error {
	// Read markdown content
	markdownContent, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Read template
	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Replace placeholder with markdown content
	templateStr := string(templateContent)
	if !strings.Contains(templateStr, "{{Placeholder Markdown}}") {
		return fmt.Errorf("template must contain {{Placeholder Markdown}} placeholder")
	}

	typstContent := strings.Replace(templateStr, "{{Placeholder Markdown}}", string(markdownContent), 1)

	// Convert to PDF
	fmt.Printf("🔄 Converting %s to PDF...\n", inputFile)
	startTime := time.Now()

	pdfBytes, err := gotypst.PDF([]byte(typstContent))
	duration := time.Since(startTime)

	if err != nil {
		return fmt.Errorf("typst compilation failed: %w", err)
	}

	if len(pdfBytes) == 0 {
		return fmt.Errorf("generated PDF is empty")
	}

	// Write PDF file
	if err := os.WriteFile(outputFile, pdfBytes, 0644); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	fmt.Printf("📄 Generated %d bytes in %v\n", len(pdfBytes), duration)
	return nil
}
