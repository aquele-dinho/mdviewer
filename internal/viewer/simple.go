package viewer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aquele_dinho/mdviewer/internal/mermaid"
	"github.com/aquele_dinho/mdviewer/internal/renderer"
	"github.com/aquele_dinho/mdviewer/internal/utils"
)

// SimpleViewer displays markdown content directly to stdout
type SimpleViewer struct {
	renderer *renderer.Renderer
}

// NewSimpleViewer creates a new simple viewer
func NewSimpleViewer(r *renderer.Renderer) *SimpleViewer {
	return &SimpleViewer{
		renderer: r,
	}
}

// View renders and displays markdown content
func (v *SimpleViewer) View(content []byte) error {
	rendered, err := v.renderer.RenderBytes(content)
	if err != nil {
		return fmt.Errorf("failed to render content: %w", err)
	}

	// Output to stdout
	fmt.Print(rendered)
	
	// Handle mermaid diagrams if enabled
	if err := v.renderMermaidDiagrams(content); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: mermaid rendering failed: %v\n", err)
	}
	
	return nil
}

// ViewFile reads and displays a markdown file
func (v *SimpleViewer) ViewFile(path string) error {
	content, err := utils.ReadFile(path)
	if err != nil {
		return err
	}

	return v.View(content)
}

// ViewStdin reads from stdin and displays the content
func (v *SimpleViewer) ViewStdin() error {
	content, err := utils.ReadFile("-")
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	return v.View(content)
}

// renderMermaidDiagrams handles local mermaid diagram rendering
func (v *SimpleViewer) renderMermaidDiagrams(content []byte) error {
	// Get render options from renderer
	opts := v.renderer.GetOptions()
	
	// Skip if mermaid is disabled or mode is URL
	if opts.NoMermaid || opts.MermaidMode == "url" {
		return nil
	}
	
	// Detect mermaid blocks
	mermaidBlocks := renderer.DetectMermaidBlocks(string(content))
	if len(mermaidBlocks) == 0 {
		return nil
	}
	
	// Create mermaid compiler
	compiler, err := mermaid.NewCompiler()
	if err != nil {
		return fmt.Errorf("failed to create mermaid compiler: %w", err)
	}
	defer compiler.Close()
	
	// Render each diagram
	for i, block := range mermaidBlocks {
		result, err := compiler.Render(block.Content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to render diagram %d: %v\n", i+1, err)
			continue
		}
		
		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "Warning: diagram %d rendering error: %v\n", i+1, result.Error)
			continue
		}
		
		// Handle based on mode
		switch opts.MermaidMode {
		case "terminal":
			// Check if terminal supports inline images
			if utils.SupportsInlineImages() {
				// Render to PNG for inline display
				width := result.Width
				height := result.Height
				if width == 0 {
					width = 800
				}
				if height == 0 {
					height = 600
				}
				
				pngData, err := compiler.RenderToPNG(block.Content, width, height)
				if err == nil {
					// Display inline
					fmt.Printf("\n📊 Mermaid Diagram (%s):\n", block.Type)
					protocol := utils.DetectImageProtocol()
					if err := utils.DisplayInlineImage(pngData, protocol); err == nil {
						fmt.Println()
					} else {
						fmt.Fprintf(os.Stderr, "Warning: failed to display inline: %v\n", err)
					}
				} else {
					fmt.Fprintf(os.Stderr, "Warning: failed to render PNG for inline display: %v\n", err)
				}
			} else {
				// Fallback to ASCII preview
				preview := mermaid.GenerateASCIIPreview(block.Type, result.Width, result.Height)
				fmt.Print("\n" + preview)
			}
			
			// Save to file only if --keep-mermaid-files is set
			if opts.KeepMermaidFiles {
				filename := fmt.Sprintf("diagram-%d.svg", i+1)
				outputPath := filepath.Join(opts.MermaidOutDir, filename)
				if err := mermaid.SaveSVGToFile(result.SVG, outputPath); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to save SVG: %v\n", err)
				} else {
					fmt.Printf("  💾 Saved to: %s\n\n", outputPath)
				}
			}
			
		case "svg":
			// Just save to file
			filename := fmt.Sprintf("diagram-%d.svg", i+1)
			outputPath := filepath.Join(opts.MermaidOutDir, filename)
			if err := mermaid.SaveSVGToFile(result.SVG, outputPath); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to save SVG: %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "Saved diagram to %s\n", outputPath)
			}
			
		case "png":
			// Render to PNG and save
			// Use dimensions from SVG or defaults
			width := result.Width
			height := result.Height
			if width == 0 {
				width = 1200
			}
			if height == 0 {
				height = 800
			}
			
			pngData, err := compiler.RenderToPNG(block.Content, width, height)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to render PNG for diagram %d: %v\n", i+1, err)
				continue
			}
			
			filename := fmt.Sprintf("diagram-%d.png", i+1)
			outputPath := filepath.Join(opts.MermaidOutDir, filename)
			if err := mermaid.SavePNGToFile(pngData, outputPath); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to save PNG: %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "Saved diagram to %s\n", outputPath)
			}
		}
	}
	
	return nil
}
