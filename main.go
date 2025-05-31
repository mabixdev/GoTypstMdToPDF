package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// Build-time variable for API-only mode
var ApiOnly string

func main() {
	// Command line flags
	var apiOnly = flag.Bool("api-only", false, "Run in API-only mode (no web UI)")
	flag.Parse()

	// Check if API-only mode is enabled via build flag or command line
	isApiOnly := *apiOnly || ApiOnly == "true"

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

	// Serve static files only if not in API-only mode
	if !isApiOnly {
		// Check if public directory exists
		if _, err := os.Stat("./public"); err == nil {
			r.Use(static.Serve("/", static.LocalFile("./public", false)))
			log.Println("📁 Serving static files from ./public")
		} else {
			log.Println("⚠️  No public directory found, running in API-only mode")
			isApiOnly = true
		}
	}

	// API routes
	api := r.Group("/api")
	{
		api.POST("/convert-to-pdf", service.ConvertToPDFHandler)
		api.POST("/convert-markdown-to-pdf", service.ConvertMarkdownToPDFHandler)
		api.GET("/stats", service.StatsHandler)
	}

	// Health check
	r.GET("/health", service.HealthHandler)

	// Root endpoint for API-only mode
	if isApiOnly {
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"service": "Markdown to PDF Service",
				"version": "1.0.0",
				"mode":    "API-only",
				"endpoints": gin.H{
					"convert":    "POST /api/convert-to-pdf",
					"convert-md": "POST /api/convert-markdown-to-pdf",
					"health":     "GET /health",
					"stats":      "GET /api/stats",
				},
				"docs": "https://github.com/mabixdev/TypstPDFService#api-endpoints",
			})
		})
	}

	// Get port from environment or use default
	port := getEnvOr("PORT", "3000")

	if isApiOnly {
		log.Printf("🚀 Markdown to PDF Service (API-only) starting on port %s", port)
		log.Printf("📡 API endpoints available at http://localhost:%s/api/", port)
		log.Printf("💡 Use POST /api/convert-to-pdf to convert markdown to PDF")
	} else {
		log.Printf("🚀 Markdown to PDF Service (Go) starting on port %s", port)
		log.Printf("📝 Access the web interface at http://localhost:%s", port)
	}

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
