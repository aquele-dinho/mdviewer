package renderer

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/aquele_dinho/mdviewer/internal/mermaid"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// HTMLRenderer converts markdown to HTML for PDF generation
type HTMLRenderer struct {
	md goldmark.Markdown
}

// NewHTMLRenderer creates a new HTML renderer
func NewHTMLRenderer() *HTMLRenderer {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,        // GitHub Flavored Markdown
			extension.Typographer, // Smart quotes, dashes
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(), // Allow raw HTML (needed for embedded Mermaid SVGs)
		),
	)

	return &HTMLRenderer{
		md: md,
	}
}

// RenderToHTML converts markdown content to HTML
func (r *HTMLRenderer) RenderToHTML(markdown string) (string, error) {
	// First, process mermaid diagrams and replace with rendered SVGs
	processed, err := r.processMermaidDiagrams(markdown)
	if err != nil {
		return "", fmt.Errorf("failed to process mermaid diagrams: %w", err)
	}

	var buf bytes.Buffer
	
	if err := r.md.Convert([]byte(processed), &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown to HTML: %w", err)
	}

	// Wrap in a complete HTML document with styling
	html := r.wrapHTML(buf.String())
	return html, nil
}

// wrapHTML wraps the HTML content in a complete document with CSS
func (r *HTMLRenderer) wrapHTML(content string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Markdown Document</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        h1, h2, h3, h4, h5, h6 {
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }
        h1 { font-size: 2em; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
        h2 { font-size: 1.5em; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
        h3 { font-size: 1.25em; }
        code {
            background-color: #f6f8fa;
            padding: 0.2em 0.4em;
            margin: 0;
            font-size: 85%%;
            border-radius: 3px;
            font-family: 'SF Mono', Monaco, Consolas, monospace;
        }
        pre {
            background-color: #f6f8fa;
            padding: 16px;
            overflow: auto;
            border-radius: 6px;
        }
        pre code {
            background-color: transparent;
            padding: 0;
        }
        blockquote {
            padding: 0 1em;
            color: #6a737d;
            border-left: 0.25em solid #dfe2e5;
            margin: 0;
        }
        table {
            border-collapse: collapse;
            width: 100%%;
        }
        table th, table td {
            padding: 6px 13px;
            border: 1px solid #dfe2e5;
        }
        table tr:nth-child(2n) {
            background-color: #f6f8fa;
        }
        img {
            max-width: 100%%;
        }
        hr {
            border: 0;
            border-top: 1px solid #eaecef;
            margin: 24px 0;
        }
        a {
            color: #0366d6;
            text-decoration: none;
        }
			a:hover {
				text-decoration: underline;
			}
			.mermaid-diagram {
				margin: 20px 0;
				text-align: center;
			}
		</style>
	</head>
	<body>
	%s
	</body>
	</html>`, content)
}

// processMermaidDiagrams detects mermaid code blocks and replaces them with rendered SVGs
func (r *HTMLRenderer) processMermaidDiagrams(markdown string) (string, error) {
	// Detect mermaid blocks
	mermaidBlocks := DetectMermaidBlocks(markdown)
	if len(mermaidBlocks) == 0 {
		return markdown, nil
	}

	// Create mermaid compiler
	compiler, err := mermaid.NewCompiler()
	if err != nil {
		// If compiler creation fails, return original markdown (mermaid will show as code blocks)
		return markdown, nil
	}
	defer compiler.Close()

	// Process markdown by replacing mermaid blocks with rendered SVGs
	result := markdown
	
	// Process blocks in reverse order to maintain string indices
	for i := len(mermaidBlocks) - 1; i >= 0; i-- {
		block := mermaidBlocks[i]
		
		// Render diagram to SVG
		svgResult, err := compiler.Render(block.Content)
		if err != nil || svgResult.Error != nil {
			// If rendering fails, keep the code block
			continue
		}
		
		// Create HTML with embedded SVG
		svgHTML := fmt.Sprintf(`<div class="mermaid-diagram">%s</div>`, svgResult.SVG)
		
		// Find the mermaid code block in the markdown
		// Pattern: ```mermaid\n...content...\n```
		pattern := regexp.MustCompile("(?s)```mermaid\\s*\\n" + regexp.QuoteMeta(block.Content) + "\\n```")
		
		// Replace the mermaid code block with the SVG
		result = pattern.ReplaceAllString(result, svgHTML)
	}

	return result, nil
}
