package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/francescoalemanno/gotypst"
	"github.com/gin-gonic/gin"
)

// PDFService handles PDF conversion operations
type PDFService struct {
	config        *Config
	activeJobs    map[string]*ConversionJob
	activeJobsMux sync.RWMutex
}

// ConversionJob represents an active conversion job
type ConversionJob struct {
	ID        string
	StartTime time.Time
	Context   context.Context
	Cancel    context.CancelFunc
}

// ConvertRequest represents the API request structure
type ConvertRequest struct {
	MarkdownContent string                 `json:"markdownContent"`
	TypstContent    string                 `json:"typstContent"`
	Options         map[string]interface{} `json:"options"`
}

// StatsResponse represents the stats API response
type StatsResponse struct {
	ActiveJobs int                      `json:"activeProcesses"`
	Jobs       []map[string]interface{} `json:"processes"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string        `json:"status"`
	Stats     StatsResponse `json:"stats"`
	Timestamp string        `json:"timestamp"`
	Message   string        `json:"message,omitempty"`
}

// NewPDFService creates a new PDF service instance
func NewPDFService() *PDFService {
	config := LoadConfig()

	// Ensure temp directory exists
	if err := os.MkdirAll(config.TempDir, 0755); err != nil {
		fmt.Printf("Warning: Could not create temp directory: %v\n", err)
	}

	return &PDFService{
		config:     config,
		activeJobs: make(map[string]*ConversionJob),
	}
}

// ConvertToPDFHandler handles the main conversion endpoint (supports both markdown and typst)
func (s *PDFService) ConvertToPDFHandler(c *gin.Context) {
	var req ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Determine conversion type
	if req.MarkdownContent != "" {
		s.convertMarkdownToPDF(c, req.MarkdownContent, req.Options)
	} else if req.TypstContent != "" {
		s.convertTypstToPDF(c, req.TypstContent, req.Options)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing markdownContent or typstContent in request body"})
	}
}

// ConvertMarkdownToPDFHandler handles dedicated markdown conversion
func (s *PDFService) ConvertMarkdownToPDFHandler(c *gin.Context) {
	var req ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	if req.MarkdownContent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing markdownContent in request body"})
		return
	}

	s.convertMarkdownToPDF(c, req.MarkdownContent, req.Options)
}

// convertMarkdownToPDF processes markdown using skeleton template
func (s *PDFService) convertMarkdownToPDF(c *gin.Context, markdownContent string, options map[string]interface{}) {
	// Validate content size
	if int64(len(markdownContent)) > s.config.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content exceeds maximum file size limit"})
		return
	}

	fmt.Printf("Starting markdown to PDF conversion for %d characters\n", len(markdownContent))

	// Read skeleton template
	skeletonContent, err := os.ReadFile(s.config.SkeletonPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read skeleton template: " + err.Error()})
		return
	}

	// Replace placeholder with markdown content
	skeletonStr := string(skeletonContent)
	if !strings.Contains(skeletonStr, "{{Placeholder Markdown}}") {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Skeleton template must contain {{Placeholder Markdown}} placeholder"})
		return
	}

	typstContent := strings.Replace(skeletonStr, "{{Placeholder Markdown}}", markdownContent, 1)

	// Convert using Typst
	s.convertTypstToPDF(c, typstContent, options)
}

// convertTypstToPDF converts Typst content to PDF using the simple gotypst API
func (s *PDFService) convertTypstToPDF(c *gin.Context, typstContent string, options map[string]interface{}) {
	// Create conversion job for tracking
	jobID := generateJobID()
	ctx, cancel := context.WithTimeout(context.Background(), s.config.TimeoutDuration)
	defer cancel()

	job := &ConversionJob{
		ID:        jobID,
		StartTime: time.Now(),
		Context:   ctx,
		Cancel:    cancel,
	}

	// Register job
	s.activeJobsMux.Lock()
	s.activeJobs[jobID] = job
	s.activeJobsMux.Unlock()

	// Cleanup job on completion
	defer func() {
		s.activeJobsMux.Lock()
		delete(s.activeJobs, jobID)
		s.activeJobsMux.Unlock()
	}()

	fmt.Printf("Starting Typst conversion for job %s (%d characters)\n", jobID, len(typstContent))

	// Convert using the simple gotypst API
	startTime := time.Now()
	pdfBytes, err := gotypst.PDF([]byte(typstContent))
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("Typst compilation failed for job %s after %v: %v\n", jobID, duration, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Typst compilation failed: " + err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	if len(pdfBytes) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Generated PDF is empty"})
		return
	}

	fmt.Printf("PDF generated successfully for job %s: %d bytes in %v\n", jobID, len(pdfBytes), duration)

	// Get filename from options
	filename := "document.pdf"
	if options != nil {
		if fn, ok := options["filename"].(string); ok && fn != "" {
			filename = fn
			if !strings.HasSuffix(filename, ".pdf") {
				filename += ".pdf"
			}
		}
	}

	// Set response headers
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	// Send PDF data
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// StatsHandler returns service statistics
func (s *PDFService) StatsHandler(c *gin.Context) {
	s.activeJobsMux.RLock()
	defer s.activeJobsMux.RUnlock()

	jobs := make([]map[string]interface{}, 0, len(s.activeJobs))
	for id, job := range s.activeJobs {
		jobs = append(jobs, map[string]interface{}{
			"id":       id,
			"duration": time.Since(job.StartTime).Milliseconds(),
			"pid":      fmt.Sprintf("go-%s", id[:8]),
		})
	}

	stats := StatsResponse{
		ActiveJobs: len(s.activeJobs),
		Jobs:       jobs,
	}

	c.JSON(http.StatusOK, stats)
}

// HealthHandler returns service health status
func (s *PDFService) HealthHandler(c *gin.Context) {
	// Test Typst compilation with simple content
	testTypst := `#set page(paper: "a4", margin: 1cm)
= Health Check
This is a test document to verify Typst compilation.`

	startTime := time.Now()
	pdfBytes, err := gotypst.PDF([]byte(testTypst))
	duration := time.Since(startTime)

	var response HealthResponse
	if err != nil || len(pdfBytes) == 0 {
		response = HealthResponse{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("Typst compilation test failed: %v", err),
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Get current stats
	s.activeJobsMux.RLock()
	stats := StatsResponse{
		ActiveJobs: len(s.activeJobs),
		Jobs:       make([]map[string]interface{}, 0),
	}
	s.activeJobsMux.RUnlock()

	response = HealthResponse{
		Status:    "healthy",
		Stats:     stats,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   fmt.Sprintf("Test compilation successful (%d bytes in %v)", len(pdfBytes), duration),
	}

	c.JSON(http.StatusOK, response)
}

// generateJobID creates a unique job identifier
func generateJobID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
