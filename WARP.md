# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

`mdviewer` is a cross-platform terminal markdown viewer with local Mermaid diagram rendering, built in Go. It renders markdown files with beautiful ANSI styling directly in the terminal and can export to PDF.

**Key Features:**
- Terminal markdown rendering with multiple themes (via Glamour)
- Local Mermaid diagram rendering using headless Chrome (chromedp)
- Stdin support for piping markdown content
- PDF export with rendered diagrams
- Cross-platform support (macOS, Linux, Windows)

## Prerequisites

- **Go 1.21+** (currently using Go 1.25.0)
- **Chrome/Chromium** installed on the system for Mermaid rendering and PDF export (chromedp requirement)
  - Without Chrome, the tool falls back to URL mode for Mermaid diagrams

## Development Commands

### Building
```bash
# Build the binary
go build -o mdviewer ./cmd/mdviewer

# Build for all platforms
./build.sh

# Build with specific version
./build.sh 1.0.0
```

### Running
```bash
# Run directly with go
go run ./cmd/mdviewer <file>

# Run the compiled binary
./mdviewer README.md

# View with specific options
./mdviewer test.md --style dark --mermaid-mode terminal

# Save Mermaid diagrams to disk (terminal mode only)
./mdviewer test.md --keep-mermaid-files
./mdviewer test.md -k --mermaid-output-dir=./diagrams
```

### Testing
```bash
# Run all tests
go test ./...

# Verbose test output
go test -v ./...

# Test a specific package
go test ./internal/renderer
```

**Note:** Currently no test files exist in the codebase.

### Dependency Management
```bash
# Download dependencies
go mod download

# Tidy up dependencies
go mod tidy

# Update all dependencies
go get -u ./...
```

## Architecture

### Entry Point
- **`cmd/mdviewer/main.go`**: CLI entry point using Cobra framework
  - Defines command-line flags and options
  - Handles input from files or stdin
  - Routes to viewer or PDF export based on flags

### Core Components

#### 1. Renderer (`internal/renderer/`)
Handles markdown rendering with multiple output formats:
- **`markdown.go`**: Main renderer using Glamour for ANSI terminal output
- **`html.go`**: HTML renderer using Goldmark (for PDF export)
- **`mermaid.go`**: Mermaid diagram detection and URL generation
- **`styles.go`**: Custom style definitions (clean style without hash prefixes)

**Key Pattern**: The renderer processes markdown in two stages:
1. Detect and process Mermaid blocks (if enabled)
2. Render to target format (ANSI or HTML)

#### 2. Mermaid Compiler (`internal/mermaid/`)
Local Mermaid rendering engine using headless Chrome:
- **`compiler.go`**: Core chromedp-based rendering logic
  - Creates fresh browser contexts per render (prevents memory leaks)
  - Renders diagrams to SVG or PNG
  - 30-second timeout per diagram
- **`embed.go`**: Embeds mermaid.min.js using go:embed
- **`svg.go`**: SVG utilities (cleaning, dimension extraction)

**Important**: Each render creates a fresh chromedp context and cancels it when done to prevent resource leaks.

#### 3. PDF Exporter (`internal/pdf/`)
PDF generation using chromedp:
- **`exporter.go`**: High-level export interface
- **`chromedp.go`**: Low-level chromedp PDF generation
- Converts markdown → HTML → PDF with rendered Mermaid diagrams

#### 4. Viewer (`internal/viewer/`)
Display logic for terminal output:
- **`simple.go`**: Simple viewer that reads files or stdin and displays rendered output

#### 5. Utilities (`internal/utils/`)
Helper functions:
- **`terminal.go`**: Terminal width detection
- **`termimg.go`**: Terminal inline image support detection (iTerm2, Kitty protocols)
- **`file.go`**: File I/O utilities
- **`browser.go`**: Browser launch utilities for Mermaid URLs

### Key Technology Stack
- **Glamour**: Terminal markdown rendering with ANSI styling
- **Cobra**: CLI framework for command-line interface
- **Goldmark**: Markdown parser (for HTML rendering)
- **Chromedp**: Headless Chrome automation (for Mermaid and PDF)

### Data Flow

1. **Terminal Viewing**:
   ```
   Input (file/stdin) → Renderer (Glamour) → Mermaid Detector → 
   Local Mermaid Renderer (chromedp) → Terminal Display
   ```

2. **PDF Export**:
   ```
   Input → HTML Renderer (Goldmark) → Mermaid Renderer (chromedp) → 
   PDF Generator (chromedp) → File Output
   ```

## Mermaid Rendering Modes

The tool supports multiple rendering modes (via `--mermaid-mode`):

- **`terminal`** (default): Display inline images if terminal supports it, otherwise ASCII preview. **Memory-only by default** - no files created unless `--keep-mermaid-files` is used.
- **`svg`**: Export SVG files to disk (temp directory by default)
- **`png`**: Export PNG files to disk (temp directory by default)
- **`url`**: Show clickable URLs to mermaid.live and mermaid.ink (fallback when Chrome unavailable)

### File Management

- **Default behavior**: Terminal mode works entirely in memory (no files created)
- **`--keep-mermaid-files` / `-k`**: Save diagram files to disk in terminal mode
- **Default output directory**: System temp directory (`os.TempDir()`)
- **Custom output directory**: Use `--mermaid-output-dir=/path/to/dir`
- **SVG/PNG export modes**: Always save files (this is their purpose)

## Important Code Patterns

### Chromedp Context Management
Always create fresh contexts and defer cancellation:
```go
ctx, cancel := chromedp.NewContext(context.Background())
defer cancel()
```

### Style Selection
The renderer supports multiple style modes:
- Built-in: `auto`, `dark`, `light`, `clean` (custom without hash prefixes)
- Custom: Path to a Glamour style JSON file

### Mermaid Block Detection
Uses regex to find mermaid code blocks and extracts:
- Diagram type
- Source code
- Line numbers (for insertion of indicators)

## Release Process

Releases are automated via GitHub Actions (`.github/workflows/release.yml`):
1. Push a version tag: `git tag v1.0.0 && git push origin v1.0.0`
2. GitHub Actions builds binaries for all platforms
3. Creates GitHub release with binaries and checksums

**Supported Platforms**:
- macOS (Intel and Apple Silicon)
- Linux (AMD64 and ARM64)
- Windows (AMD64)

## Common Development Tasks

### Adding a New CLI Flag
1. Add flag variable in `cmd/mdviewer/main.go` (global vars)
2. Register flag in `init()` function
3. Pass flag value through `RenderOptions` or use directly
4. Update help text in Cobra command definition

### Adding a New Rendering Mode
1. Define mode constant in `internal/mermaid/compiler.go`
2. Update mode handling in `internal/renderer/markdown.go`
3. Add flag documentation in `cmd/mdviewer/main.go`

### Debugging Chromedp Issues
- Check Chrome/Chromium installation
- Increase timeout values if diagrams are timing out
- Use chromedp debug logging (not currently enabled)
- Verify mermaid.min.js is properly embedded

## Testing Approach

When writing tests:
- Mock chromedp for Mermaid rendering tests (avoid real browser dependency)
- Test markdown rendering with known inputs
- Test style selection and width detection
- Use table-driven tests for multiple scenarios

## Dependencies Update Strategy

- Keep Go toolchain up to date (currently on 1.25.0)
- Update Glamour for terminal rendering improvements
- Update chromedp for Chrome compatibility
- Pin Goldmark and other parsers to stable versions
