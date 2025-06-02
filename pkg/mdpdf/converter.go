// Package mdpdf provides markdown to PDF conversion functionality using Typst
package mdpdf

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/francescoalemanno/gotypst"
)

// Converter handles markdown to PDF conversions
type Converter struct {
	templatePath string
	options      *Options
}

// Options configures the conversion process
type Options struct {
	// TemplatePath is the path to the Typst template file
	TemplatePath string
	// MaxFileSize limits the input markdown size (default: 50MB)
	MaxFileSize int64
	// Timeout sets the maximum conversion time (default: 30s)
	Timeout time.Duration
}

// DefaultOptions returns sensible default options
func DefaultOptions() *Options {
	return &Options{
		TemplatePath: "exam-template.typ",
		MaxFileSize:  50 * 1024 * 1024, // 50MB
		Timeout:      30 * time.Second,
	}
}

// NewConverter creates a new converter with the given options
func NewConverter(opts *Options) (*Converter, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	// Validate template exists
	if _, err := os.Stat(opts.TemplatePath); err != nil {
		return nil, fmt.Errorf("template file not found: %w", err)
	}

	return &Converter{
		templatePath: opts.TemplatePath,
		options:      opts,
	}, nil
}

// ConvertFromString converts markdown string to PDF bytes
func (c *Converter) ConvertFromString(ctx context.Context, markdownContent string) ([]byte, error) {
	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Validate input size
	if int64(len(markdownContent)) > c.options.MaxFileSize {
		return nil, fmt.Errorf("content exceeds maximum size limit (%d bytes)", c.options.MaxFileSize)
	}

	// Read template
	templateContent, err := os.ReadFile(c.templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	// Check context again before processing
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Replace placeholder
	templateStr := string(templateContent)
	if !strings.Contains(templateStr, "{{Placeholder Markdown}}") {
		return nil, fmt.Errorf("template must contain {{Placeholder Markdown}} placeholder")
	}

	typstContent := strings.Replace(templateStr, "{{Placeholder Markdown}}", markdownContent, 1)

	// Convert to PDF with context handling
	// Since gotypst.PDF doesn't support context, we'll use a goroutine with timeout
	type result struct {
		pdfBytes []byte
		err      error
	}

	resultChan := make(chan result, 1)
	go func() {
		pdfBytes, err := gotypst.PDF([]byte(typstContent))
		resultChan <- result{pdfBytes: pdfBytes, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		if res.err != nil {
			return nil, fmt.Errorf("typst compilation failed: %w", res.err)
		}

		if len(res.pdfBytes) == 0 {
			return nil, fmt.Errorf("generated PDF is empty")
		}

		return res.pdfBytes, nil
	}
}

// ConvertFromFile converts markdown file to PDF bytes
func (c *Converter) ConvertFromFile(ctx context.Context, inputPath string) ([]byte, error) {
	markdownContent, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	return c.ConvertFromString(ctx, string(markdownContent))
}

// ConvertFromFileToFile converts markdown file to PDF file
func (c *Converter) ConvertFromFileToFile(ctx context.Context, inputPath, outputPath string) error {
	pdfBytes, err := c.ConvertFromFile(ctx, inputPath)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, pdfBytes, 0644); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	return nil
}

// ConvertFromStringToFile converts markdown string to PDF file
func (c *Converter) ConvertFromStringToFile(ctx context.Context, markdownContent, outputPath string) error {
	pdfBytes, err := c.ConvertFromString(ctx, markdownContent)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, pdfBytes, 0644); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	return nil
}

// GetTemplateContent returns the current template content
func (c *Converter) GetTemplateContent() (string, error) {
	content, err := os.ReadFile(c.templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}
	return string(content), nil
}

// ValidateTemplate checks if the template is valid
func (c *Converter) ValidateTemplate() error {
	content, err := c.GetTemplateContent()
	if err != nil {
		return err
	}

	if !strings.Contains(content, "{{Placeholder Markdown}}") {
		return fmt.Errorf("template must contain {{Placeholder Markdown}} placeholder")
	}

	// Test compilation with minimal content
	testContent := strings.Replace(content, "{{Placeholder Markdown}}", "# Test", 1)
	_, err = gotypst.PDF([]byte(testContent))
	if err != nil {
		return fmt.Errorf("template compilation test failed: %w", err)
	}

	return nil
}

// Simple convenience functions

// QuickConvert provides a simple one-line conversion
func QuickConvert(ctx context.Context, markdownContent string) ([]byte, error) {
	converter, err := NewConverter(DefaultOptions())
	if err != nil {
		return nil, err
	}
	return converter.ConvertFromString(ctx, markdownContent)
}

// QuickConvertFile provides a simple file-to-file conversion
func QuickConvertFile(ctx context.Context, inputPath, outputPath string) error {
	converter, err := NewConverter(DefaultOptions())
	if err != nil {
		return err
	}
	return converter.ConvertFromFileToFile(ctx, inputPath, outputPath)
}
