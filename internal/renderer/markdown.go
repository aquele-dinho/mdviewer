package renderer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
)

// RenderOptions contains configuration for rendering markdown
type RenderOptions struct {
	Style            string // Style name: "dark", "light", "auto"
	Width            int    // Terminal width for wrapping
	NoMermaid        bool   // Skip mermaid diagram detection
	MermaidMode      string // Mermaid rendering mode: "terminal", "svg", "url"
	MermaidOutDir    string // Output directory for SVG files
	KeepMermaidFiles bool   // Save mermaid diagram files to disk
}

// Renderer handles markdown rendering
type Renderer struct {
	options RenderOptions
	glamour *glamour.TermRenderer
}

// NewRenderer creates a new markdown renderer
func NewRenderer(opts RenderOptions) (*Renderer, error) {
	// Set defaults
	if opts.Width == 0 {
		opts.Width = 80
	}
	if opts.Style == "" {
		opts.Style = "auto"
	}

	// Create glamour renderer
	glamourOpts := []glamour.TermRendererOption{
		glamour.WithWordWrap(opts.Width),
	}

	// Handle style selection
	switch opts.Style {
	case "auto":
		glamourOpts = append(glamourOpts, glamour.WithAutoStyle())
	case "dark", "light":
		glamourOpts = append(glamourOpts, glamour.WithStylePath(opts.Style))
	case "notty", "clean":
		// Use custom clean style without hash prefixes
		glamourOpts = append(glamourOpts, glamour.WithStylesFromJSONBytes([]byte(CustomStyle)))
	default:
		// Try to load custom style file
		glamourOpts = append(glamourOpts, glamour.WithStylePath(opts.Style))
	}

	glamourRenderer, err := glamour.NewTermRenderer(glamourOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create glamour renderer: %w", err)
	}

	return &Renderer{
		options: opts,
		glamour: glamourRenderer,
	}, nil
}

// PreprocessLinks exposes the link preprocessing function
func (r *Renderer) PreprocessLinks(content string) string {
	return PreprocessLinks(content)
}

// Render renders markdown content to ANSI-styled terminal output
func (r *Renderer) Render(content string) (string, error) {
	// First preprocess links (Markdown + Obsidian-style) so Glamour can
	// render them as normal links/images.
	content = PreprocessLinks(content)

	// Check for mermaid diagrams if enabled
	if !r.options.NoMermaid {
		content = r.processMermaid(content)
	}

	// Render with glamour
	rendered, err := r.glamour.Render(content)
	if err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	return rendered, nil
}

// RenderBytes renders markdown bytes to ANSI-styled terminal output
func (r *Renderer) RenderBytes(content []byte) (string, error) {
	return r.Render(string(content))
}

// GetOptions returns the renderer's options
func (r *Renderer) GetOptions() RenderOptions {
	return r.options
}

// processMermaid adds visual indicators for mermaid diagrams
func (r *Renderer) processMermaid(content string) string {
	// Detect mermaid code blocks
	mermaidBlocks := DetectMermaidBlocks(content)
	
	if len(mermaidBlocks) == 0 {
		return content
	}

	// Only add URL indicators in URL mode
	// In terminal/svg/png modes, the mermaid code block will be rendered normally
	// and the local rendering will happen after markdown rendering
	if r.options.MermaidMode != "url" {
		return content
	}

	// URL mode: Add visual indicators with URLs before mermaid blocks
	lines := strings.Split(content, "\n")
	offset := 0
	
	for _, block := range mermaidBlocks {
		// Generate viewing URLs
		liveURL := GenerateMermaidLiveURL(block)
		imageURL := GetMermaidInkURL(block)
		
		// Create the indicator with clickable URLs
		indicator := fmt.Sprintf("\n> 📊 **Mermaid Diagram** (%s)\n> \n> 🔗 View: <%s>\n> 📷 Image: <%s>\n", 
			block.Type, liveURL, imageURL)
		
		// Calculate the actual line index (0-based)
		insertIndex := block.StartLine - 1 + offset
		
		if insertIndex >= 0 && insertIndex < len(lines) {
			// Insert the indicator before the mermaid code block
			lines = append(lines[:insertIndex], append([]string{indicator}, lines[insertIndex:]...)...)
			offset++ // Adjust offset for next insertion
		}
	}

	return strings.Join(lines, "\n")
}
