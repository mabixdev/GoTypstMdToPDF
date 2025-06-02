# Extracting the PDF Service Without Frontend

This guide shows you different ways to use just the PDF conversion functionality without the web frontend.

## 🎯 Option 1: Use as API Service (Recommended)

The easiest approach is to run the existing service and make HTTP requests to it.

### Start the Service
```bash
# Full service (with web UI)
make run

# API-only mode (no web UI)
make run-api
# or
./bin/md-pdf-service -api-only
```

### Use the API
```bash
# Convert markdown to PDF
curl -X POST http://localhost:3000/api/convert-to-pdf \
  -H "Content-Type: application/json" \
  -d '{
    "markdownContent": "# Test\n\nMath: $E = mc^2$",
    "options": {"filename": "test.pdf"}
  }' \
  --output test.pdf
```

**Advantages:**
- ✅ No code changes needed
- ✅ Works with any programming language
- ✅ Can be deployed as microservice
- ✅ Same proven functionality

**Use cases:**
- Integration with other applications
- Microservice architecture
- Remote PDF generation
- Multiple client applications

## 🖥️ Option 2: Command Line Interface

Build and use the CLI version for direct file conversion.

### Build and Use
```bash
# Build CLI version
make build-cli

# Convert files
./bin/md-pdf-cli -input test.md -output exam.pdf
./bin/md-pdf-cli -input test.md -template custom-template.typ
```

### CLI Help
```bash
./bin/md-pdf-cli -help
```

**Advantages:**
- ✅ Simple file-to-file conversion
- ✅ Can be used in scripts/automation
- ✅ No server needed
- ✅ Lightweight single binary

**Use cases:**
- Batch processing
- Shell scripts
- CI/CD pipelines
- Local file conversion

## 📦 Option 3: Go Library Package

Import the core functionality into your Go application.

### Library Usage
```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/mabixdev/TypstPDFService/pkg/mdpdf"
)

func main() {
    // Simple conversion
    markdown := "# Test\n\nMath: $E = mc^2$"
    ctx := context.Background()
    pdfBytes, err := mdpdf.QuickConvert(ctx, markdown)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Generated PDF: %d bytes\n", len(pdfBytes))
    
    // File conversion
    err = mdpdf.QuickConvertFile(ctx, "input.md", "output.pdf")
    if err != nil {
        log.Fatal(err)
    }
    
    // Advanced usage with custom options
    opts := &mdpdf.Options{
        TemplatePath: "custom-template.typ",
        MaxFileSize:  10 * 1024 * 1024, // 10MB
    }
    
    converter, err := mdpdf.NewConverter(opts)
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background() // or context.WithTimeout for timeout handling
    pdfBytes, err = converter.ConvertFromString(ctx, markdown)
    if err != nil {
        log.Fatal(err)
    }
}
```

**Advantages:**
- ✅ Direct integration in Go apps
- ✅ Type safety and compile-time checks
- ✅ Best performance (no HTTP overhead)
- ✅ Full control over configuration

**Use cases:**
- Go web applications
- Go CLI tools
- Embedded in larger Go services
- Custom processing pipelines

## 🔧 Build Instructions

### Prerequisites
```bash
# Install Go 1.21+
go version

# Clone repository
git clone <your-repo>
cd TypstPDFService
```

### Build Different Versions
```bash
# API-only service
make build-api

# CLI version
make build-cli

# Full service (with web UI)
make build

# All versions
make build-all
```

### Run Examples
```bash
# API service
make run-api

# CLI conversion
make run-cli

# Full service
make run
```

## 🐳 Docker Deployment

### API-Only Container
```bash
# Build API-only image
make docker-build-api

# Run API-only container
make docker-run-api
```

### Custom Dockerfile for API-only
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -ldflags="-X main.ApiOnly=true" -o md-pdf-service .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/md-pdf-service .
COPY exam-template.typ ./
EXPOSE 3000
CMD ["./md-pdf-service"]
```

## 📝 Integration Examples

### Python Client
```python
import requests

def convert_markdown_to_pdf(markdown_content, filename="document.pdf"):
    url = "http://localhost:3000/api/convert-to-pdf"
    payload = {
        "markdownContent": markdown_content,
        "options": {"filename": filename}
    }
    
    response = requests.post(url, json=payload)
    
    if response.status_code == 200:
        with open(filename, "wb") as f:
            f.write(response.content)
        return True
    else:
        print(f"Error: {response.text}")
        return False

# Usage
markdown = "# Test\n\nMath: $E = mc^2$"
convert_markdown_to_pdf(markdown, "test.pdf")
```

### Node.js Client
```javascript
const fetch = require('node-fetch');
const fs = require('fs');

async function convertToPDF(markdown, filename = 'output.pdf') {
    const response = await fetch('http://localhost:3000/api/convert-to-pdf', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            markdownContent: markdown,
            options: { filename }
        })
    });
    
    if (response.ok) {
        const buffer = await response.buffer();
        fs.writeFileSync(filename, buffer);
        console.log(`PDF saved as ${filename}`);
    } else {
        console.error('Conversion failed:', await response.text());
    }
}
```

### Shell Script
```bash
#!/bin/bash

convert_markdown() {
    local input_file="$1"
    local output_file="$2"
    
    if [ ! -f "$input_file" ]; then
        echo "Error: Input file not found: $input_file"
        return 1
    fi
    
    # Read markdown content
    local markdown_content=$(cat "$input_file")
    
    # Convert using API
    curl -s -X POST http://localhost:3000/api/convert-to-pdf \
        -H "Content-Type: application/json" \
        -d "$(jq -n --arg content "$markdown_content" --arg filename "$output_file" \
            '{markdownContent: $content, options: {filename: $filename}}')" \
        --output "$output_file"
    
    if [ $? -eq 0 ]; then
        echo "✅ Converted $input_file to $output_file"
    else
        echo "❌ Conversion failed"
        return 1
    fi
}

# Usage
convert_markdown "test.md" "output.pdf"
```

## 🚀 Performance Considerations

| Method | Startup Time | Memory Usage | Throughput | Best For |
|--------|-------------|--------------|------------|----------|
| API Service | ~1s | ~50MB | High | Multiple clients |
| CLI | ~0.1s | ~30MB | Medium | Single conversions |
| Go Library | ~0.01s | ~25MB | Highest | Embedded usage |

## 🔍 Troubleshooting

### Common Issues

**Template not found:**
```bash
# Ensure exam-template.typ exists
ls -la exam-template.typ

# Or specify custom template
./bin/md-pdf-cli -input test.md -template /path/to/template.typ
```

**API connection failed:**
```bash
# Check if service is running
curl http://localhost:3000/health

# Check port availability
netstat -tlnp | grep 3000
```

**Large file conversion fails:**
```bash
# Increase timeout
export TIMEOUT_DURATION=60s

# Increase file size limit  
export MAX_FILE_SIZE=104857600  # 100MB
```

## 📖 Next Steps

1. **Choose your integration method** based on your use case
2. **Test with sample files** to ensure compatibility
3. **Configure templates** for your document format needs
4. **Set up monitoring** for production deployments
5. **Scale horizontally** by running multiple API instances

For detailed API documentation, see: [examples/api-usage.md](examples/api-usage.md) 