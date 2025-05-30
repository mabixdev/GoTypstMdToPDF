package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the PDF service
	service := NewPDFService()

	// Create Gin router
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Serve static files
	r.Use(static.Serve("/", static.LocalFile("./public", false)))

	// API routes
	api := r.Group("/api")
	{
		api.POST("/convert-to-pdf", service.ConvertToPDFHandler)
		api.POST("/convert-markdown-to-pdf", service.ConvertMarkdownToPDFHandler)
		api.GET("/stats", service.StatsHandler)
	}

	// Health check
	r.GET("/health", service.HealthHandler)

	// Get port from environment or use default
	port := getEnvOr("PORT", "3000")

	log.Printf("🚀 Markdown to PDF Service (Go) starting on port %s", port)
	log.Printf("📝 Access the web interface at http://localhost:%s", port)

	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnvOr(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Server configuration
type Config struct {
	Port            string
	TempDir         string
	SkeletonPath    string
	MaxFileSize     int64
	TimeoutDuration time.Duration
}

func LoadConfig() *Config {
	port := getEnvOr("PORT", "3000")
	tempDir := getEnvOr("TEMP_DIR", "./temp")
	skeletonPath := getEnvOr("SKELETON_PATH", "./exam-template.typ")

	maxFileSize := int64(50 * 1024 * 1024) // 50MB default
	if sizeStr := os.Getenv("MAX_FILE_SIZE"); sizeStr != "" {
		if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
			maxFileSize = size
		}
	}

	timeoutDuration := 30 * time.Second // 30s default
	if timeoutStr := os.Getenv("TIMEOUT_DURATION"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			timeoutDuration = timeout
		}
	}

	return &Config{
		Port:            port,
		TempDir:         tempDir,
		SkeletonPath:    skeletonPath,
		MaxFileSize:     maxFileSize,
		TimeoutDuration: timeoutDuration,
	}
}
