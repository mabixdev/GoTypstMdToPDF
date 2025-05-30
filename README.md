# Markdown to PDF Service with Typst

A powerful web service that converts Markdown content with embedded LaTeX mathematics into beautiful PDFs using Typst. Perfect for academic documents, exams, reports, and any content requiring mathematical notation.

## 🚀 Features

- **Markdown Support**: Full Markdown syntax support with headers, lists, formatting, and code blocks
- **LaTeX Mathematics**: Embedded LaTeX math expressions (`$inline$` and `$$display$$`)
- **Beautiful Typography**: Professional document formatting using Typst
- **Web Interface**: Easy-to-use web interface for instant conversions
- **Template System**: Customizable document templates via skeleton files
- **Real-time Processing**: Fast conversion with progress indicators
- **File Management**: Automatic temporary file cleanup and download handling

## 📋 Requirements

- Node.js (v14 or higher)
- [Typst CLI](https://github.com/typst/typst) (v0.12.0 or higher)
- npm or yarn

## 🛠️ Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/mabixdev/MdPDFServicewithTypst.git
   cd MdPDFServicewithTypst
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Install Typst CLI:**
   ```bash
   # macOS (via Homebrew)
   brew install typst
   
   # Or download from https://github.com/typst/typst/releases
   ```

4. **Create temp directory:**
   ```bash
   mkdir temp
   ```

## 🚀 Usage

### Starting the Server

```bash
npm start
# or
node server.js
```

The service will be available at `http://localhost:3000`

### Web Interface

1. Open `http://localhost:3000` in your browser
2. Paste your Markdown content with LaTeX math into the text area
3. Optionally specify a filename
4. Click "Generate PDF" to download your document

### Example Markdown Input

```markdown
# Chemistry Exam

## Question 1

Calculate the pH of a 0.100 M $\text{H}_2\text{SO}_3$ solution.

The acid dissociation constants are:
- $K_{a1} = 1.7 \times 10^{-2}$
- $K_{a2} = 6.4 \times 10^{-8}$

### Solution

For the quadratic equation $ax^2 + bx + c = 0$:

$$x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}$$

**Code Example:**
```python
import math

def calculate_ph(concentration):
    return -math.log10(concentration)
```
```

## 🔧 API Endpoints

### Convert Markdown to PDF

```bash
POST /api/convert-to-pdf
Content-Type: application/json

{
  "markdownContent": "# Your markdown here...",
  "options": {
    "filename": "document.pdf"
  }
}
```

### Health Check

```bash
GET /health
```

### Service Statistics

```bash
GET /api/stats
```

## ⚙️ Configuration

### Skeleton Template

The system uses `sceleton.typ` as a template. Customize it to change document formatting:

```typst
#import "@preview/mitex:0.2.4": mitex
#import "@preview/cmarker:0.1.1"

// Document setup
#set page(margin: 2cm)
#set text(font: "Arial", size: 11pt)

// Your custom styling here...

#cmarker.render(`
{{Placeholder Markdown}}
`, math: mitex)
```

### Environment Variables

- `PORT`: Server port (default: 3000)
- `TEMP_DIR`: Temporary files directory
- `TYPST_TIMEOUT`: Compilation timeout in ms

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Frontend  │────│   Node.js API   │────│   Typst CLI     │
│   (HTML/JS)     │    │   (Express)     │    │   (Compiler)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                        │                        │
         │                        │                        │
    User Input              Template System           PDF Output
   (Markdown +              (sceleton.typ +           (Beautiful
    LaTeX math)              user content)             document)
```

### Workflow

1. **User Input**: Markdown with embedded LaTeX via web interface
2. **Template Processing**: Content inserted into `sceleton.typ` template
3. **Temporary Files**: Generated `.typ` file saved to temp directory
4. **Typst Compilation**: CLI converts `.typ` to PDF
5. **File Delivery**: PDF downloaded, temporary files cleaned up

## 📁 Project Structure

```
MdPDFServicewithTypst/
├── public/                 # Frontend assets
│   ├── index.html         # Web interface
│   ├── script.js          # Frontend JavaScript
│   └── style.css          # Styling
├── services/              # Backend services
│   └── typst-pdf-service.js  # Core conversion logic
├── temp/                  # Temporary files (auto-created)
├── sceleton.typ          # Document template
├── server.js             # Express server
├── package.json          # Dependencies
└── README.md             # This file
```

## 🔒 Security Features

- Input validation and sanitization
- File size limits
- Temporary file cleanup
- Process timeout protection
- Safe template processing

## 🛠️ Development

### Running in Development Mode

```bash
# Install nodemon for auto-restart
npm install -g nodemon

# Start with auto-reload
nodemon server.js
```

### Testing the API

```bash
# Test with curl
curl -X POST http://localhost:3000/api/convert-to-pdf \
  -H "Content-Type: application/json" \
  -d '{
    "markdownContent": "# Test\n\nMath: $E = mc^2$",
    "options": {"filename": "test.pdf"}
  }' \
  --output test.pdf
```

## 📝 Examples

See the `test.md` file for a complete chemistry exam example with complex LaTeX mathematics.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is open source. Please check the license file for details.

## 🔗 Links

- [Typst Documentation](https://typst.app/docs/)
- [Typst CLI Repository](https://github.com/typst/typst)
- [LaTeX Math Reference](https://katex.org/docs/supported.html)

## 🐛 Troubleshooting

### Common Issues

1. **"Typst CLI not found"**
   - Install Typst CLI: `brew install typst` (macOS) or download from GitHub releases

2. **"Package requires newer Typst version"**
   - Update Typst: `brew upgrade typst`
   - Or modify `sceleton.typ` to use compatible package versions

3. **"Compilation failed"**
   - Check Markdown syntax
   - Ensure LaTeX math is properly escaped
   - Verify code blocks are properly formatted

4. **"Empty PDF generated"**
   - Check server logs for compilation errors
   - Verify template placeholder is present in `sceleton.typ`

---

Made with ❤️ using [Typst](https://typst.app/) and [Node.js](https://nodejs.org/) 