package mdpdf

import (
	"context"
	"testing"
	"time"
)

func getTestOptions() *Options {
	opts := DefaultOptions()
	opts.TemplatePath = "../../exam-template.typ" // Adjust path for test
	return opts
}

func TestConvertFromStringWithContext(t *testing.T) {
	converter, err := NewConverter(getTestOptions())
	if err != nil {
		t.Fatalf("Failed to create converter: %v", err)
	}

	// Test with normal context
	ctx := context.Background()
	markdownContent := "# Test\n\nThis is a test document."

	pdfBytes, err := converter.ConvertFromString(ctx, markdownContent)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("Generated PDF is empty")
	}
}

func TestConvertFromStringWithTimeout(t *testing.T) {
	converter, err := NewConverter(getTestOptions())
	if err != nil {
		t.Fatalf("Failed to create converter: %v", err)
	}

	// Test with timeout context (5 seconds should be enough for normal conversion)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	markdownContent := "# Test\n\nThis is a test document."

	pdfBytes, err := converter.ConvertFromString(ctx, markdownContent)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("Generated PDF is empty")
	}
}

func TestConvertFromStringWithCancelledContext(t *testing.T) {
	converter, err := NewConverter(getTestOptions())
	if err != nil {
		t.Fatalf("Failed to create converter: %v", err)
	}

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	markdownContent := "# Test\n\nThis is a test document."

	_, err = converter.ConvertFromString(ctx, markdownContent)
	if err == nil {
		t.Fatal("Expected error due to cancelled context")
	}

	if err != context.Canceled {
		t.Fatalf("Expected context.Canceled error, got: %v", err)
	}
}
