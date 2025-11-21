package mermaid

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

// Compiler renders Mermaid diagrams to SVG using chromedp
type Compiler struct {
	// Empty struct - we'll create fresh contexts per render
}

// NewCompiler creates a new Mermaid compiler with chromedp
func NewCompiler() (*Compiler, error) {
	return &Compiler{}, nil
}


// Render compiles a mermaid diagram to SVG
func (c *Compiler) Render(diagramCode string) (*SVGResult, error) {
	// Create fresh chromedp context for this render
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Create timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Escape the diagram code for JavaScript
	diagramJSON, err := json.Marshal(diagramCode)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal diagram code: %w", err)
	}

	// First load mermaid.js, then render
	var result map[string]interface{}
	err = chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		// Load mermaid.js
		chromedp.Evaluate(MermaidJS, nil),
		// Initialize mermaid
		chromedp.Evaluate(`
			mermaid.initialize({
				startOnLoad: false,
				theme: 'default',
				securityLevel: 'loose'
			});
			window.mermaidResult = null;
		`, nil),
		// Start rendering the diagram
		chromedp.Evaluate(fmt.Sprintf(`
			(async function() {
				try {
					const diagramCode = %s;
					const result = await mermaid.render('diagram-' + Date.now(), diagramCode);
					window.mermaidResult = {
						svg: result.svg,
						error: null
					};
				} catch (error) {
					window.mermaidResult = {
						svg: null,
						error: error.message || String(error)
					};
				}
			})();
		`, diagramJSON), nil),
		// Poll until result is ready
		chromedp.Poll(`window.mermaidResult`, &result, chromedp.WithPollingTimeout(20*time.Second)),
	)


	if err != nil {
		return &SVGResult{Error: fmt.Errorf("chromedp evaluation failed: %w", err)}, nil
	}

	// Check for errors in the result
	if errMsg, ok := result["error"].(string); ok && errMsg != "" {
		return &SVGResult{Error: fmt.Errorf("mermaid rendering error: %s", errMsg)}, nil
	}

	// Extract SVG
	svg, ok := result["svg"].(string)
	if !ok || svg == "" {
		return &SVGResult{Error: fmt.Errorf("no SVG returned from mermaid (result: %+v)", result)}, nil
	}

	// Clean and extract dimensions
	svg = CleanSVG(svg)
	width, height := ExtractSVGDimensions(svg)

	return &SVGResult{
		SVG:    svg,
		Width:  width,
		Height: height,
		Error:  nil,
	}, nil
}

// RenderToPNG renders a mermaid diagram to PNG bytes
func (c *Compiler) RenderToPNG(diagramCode string, width, height int) ([]byte, error) {
	// Create fresh chromedp context for this render
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Create timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Escape the diagram code for JavaScript
	diagramJSON, err := json.Marshal(diagramCode)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal diagram code: %w", err)
	}

	// First load mermaid.js, then render, then take screenshot
	var pngBytes []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		// Set viewport size
		chromedp.EmulateViewport(int64(width), int64(height)),
		// Load mermaid.js
		chromedp.Evaluate(MermaidJS, nil),
		// Initialize mermaid
		chromedp.Evaluate(`
			mermaid.initialize({
				startOnLoad: false,
				theme: 'default',
				securityLevel: 'loose'
			});
			window.mermaidResult = null;
		`, nil),
		// Render and inject into DOM
		chromedp.Evaluate(fmt.Sprintf(`
			(async function() {
				try {
					const diagramCode = %s;
					const result = await mermaid.render('diagram-' + Date.now(), diagramCode);
					// Inject SVG into body
					document.body.innerHTML = result.svg;
					window.mermaidResult = { success: true };
				} catch (error) {
					window.mermaidResult = { success: false, error: error.message };
				}
			})();
		`, diagramJSON), nil),
		// Wait for rendering
		chromedp.Poll(`window.mermaidResult`, nil, chromedp.WithPollingTimeout(20*time.Second)),
		// Take screenshot
		chromedp.FullScreenshot(&pngBytes, 100),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to render PNG: %w", err)
	}

	return pngBytes, nil
}

// Close shuts down the compiler and releases resources
func (c *Compiler) Close() {
	// Nothing to clean up - contexts are per-render
}

// RenderMode defines how diagrams should be rendered
type RenderMode string

const (
	RenderModeTerminal RenderMode = "terminal" // Display in terminal (ASCII preview + save)
	RenderModeSVG      RenderMode = "svg"      // Export to SVG files only
	RenderModePNG      RenderMode = "png"      // Export to PNG files only
	RenderModeURL      RenderMode = "url"      // Use external URLs (fallback)
)
