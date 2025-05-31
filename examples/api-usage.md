# Using the PDF Service as API-Only

## Start the Service
```bash
./bin/md-pdf-service
```

## API Usage Examples

### 1. Using curl

```bash
# Convert markdown to PDF
curl -X POST http://localhost:3000/api/convert-to-pdf \
  -H "Content-Type: application/json" \
  -d '{
    "markdownContent": "# Test Document\n\nMath: $E = mc^2$\n\n$$\\int_0^\\infty e^{-x^2} dx = \\frac{\\sqrt{\\pi}}{2}$$",
    "options": {"filename": "test.pdf"}
  }' \
  --output test.pdf
```

### 2. Using Python

```python
import requests
import json

# Your markdown content
markdown_content = """
# Chemistry Exam

## Question 1
Calculate the pH of a 0.100 M $\\text{H}_2\\text{SO}_3$ solution.

$$x = \\frac{-b \\pm \\sqrt{b^2 - 4ac}}{2a}$$
"""

# API request
url = "http://localhost:3000/api/convert-to-pdf"
payload = {
    "markdownContent": markdown_content,
    "options": {"filename": "chemistry_exam.pdf"}
}

response = requests.post(url, json=payload)

if response.status_code == 200:
    with open("chemistry_exam.pdf", "wb") as f:
        f.write(response.content)
    print("PDF generated successfully!")
else:
    print(f"Error: {response.status_code} - {response.text}")
```

### 3. Using JavaScript (Node.js)

```javascript
const fetch = require('node-fetch');
const fs = require('fs');

async function convertToPDF(markdownContent, filename = 'document.pdf') {
    try {
        const response = await fetch('http://localhost:3000/api/convert-to-pdf', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                markdownContent,
                options: { filename }
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${await response.text()}`);
        }

        const buffer = await response.buffer();
        fs.writeFileSync(filename, buffer);
        console.log(`PDF saved as ${filename}`);
    } catch (error) {
        console.error('Conversion failed:', error.message);
    }
}

// Usage
const markdown = `
# My Document
This is a test with math: $x^2 + y^2 = z^2$
`;

convertToPDF(markdown, 'output.pdf');
```

### 4. Using Go (as a client)

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

type ConvertRequest struct {
    MarkdownContent string                 `json:"markdownContent"`
    Options         map[string]interface{} `json:"options"`
}

func convertToPDF(markdownContent, filename string) error {
    req := ConvertRequest{
        MarkdownContent: markdownContent,
        Options: map[string]interface{}{
            "filename": filename,
        },
    }

    jsonData, err := json.Marshal(req)
    if err != nil {
        return err
    }

    resp, err := http.Post("http://localhost:3000/api/convert-to-pdf", 
        "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
    }

    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.Copy(file, resp.Body)
    return err
}

func main() {
    markdown := `# Test Document
    
Math example: $E = mc^2$

$$\int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2}$$`

    if err := convertToPDF(markdown, "test.pdf"); err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println("PDF generated successfully!")
    }
}
```

### 5. Health Check

```bash
curl http://localhost:3000/health
```

### 6. Service Stats

```bash
curl http://localhost:3000/api/stats
``` 