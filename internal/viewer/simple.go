package viewer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	opts := v.renderer.GetOptions()

	// For URL mode or when mermaid is disabled, keep existing simple pipeline.
	if opts.NoMermaid || opts.MermaidMode == "url" {
		rendered, err := v.renderer.RenderBytes(content)
		if err != nil {
			return fmt.Errorf("failed to render content: %w", err)
		}

		fmt.Print(rendered)
		return nil
	}

	// For terminal/svg/png modes, render markdown in segments between mermaid
	// blocks so that diagrams appear inline rather than appended at the end.
	if err := v.renderWithMermaidSegments(content, opts); err != nil {
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

// renderWithMermaidSegments renders the markdown in segments separated by
// mermaid blocks so that diagrams appear inline at their original positions.
func (v *SimpleViewer) renderWithMermaidSegments(content []byte, opts renderer.RenderOptions) error {
	text := string(content)
	blocks := renderer.DetectMermaidBlocks(text)
	if len(blocks) == 0 {
		// No mermaid: just render normally.
		rendered, err := v.renderer.RenderBytes(content)
		if err != nil {
			return fmt.Errorf("failed to render content: %w", err)
		}
		fmt.Print(rendered)
		return nil
	}

	lines := strings.Split(text, "\n")

	compiler, err := mermaid.NewCompiler()
	if err != nil {
		// Fall back to plain rendered markdown if compiler cannot be created.
		rendered, rErr := v.renderer.RenderBytes(content)
		if rErr != nil {
			return fmt.Errorf("failed to render content without mermaid: %w", rErr)
		}
		fmt.Print(rendered)
		return fmt.Errorf("failed to create mermaid compiler: %w", err)
	}
	defer compiler.Close()

	currLine := 0

	for i, block := range blocks {
		start := block.StartLine - 1
		end := block.EndLine // end is exclusive for the next segment
		if start < currLine {
			start = currLine
		}
		if start > len(lines) {
			start = len(lines)
		}
		if end > len(lines) {
			end = len(lines)
		}

		// Render the markdown segment before this block.
		if start > currLine {
			segment := strings.Join(lines[currLine:start], "\n")
			if strings.TrimSpace(segment) != "" {
				rendered, err := v.renderer.RenderBytes([]byte(segment))
				if err != nil {
					return fmt.Errorf("failed to render markdown segment: %w", err)
				}
				fmt.Print(rendered)
			}
		}

		// Render the mermaid block inline.
		if err := v.renderSingleMermaidBlock(compiler, block, i, opts); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to render mermaid diagram %d: %v\\n", i+1, err)
		}

		currLine = end
	}

	// Render any remaining markdown after the last block.
	if currLine < len(lines) {
		segment := strings.Join(lines[currLine:], "\n")
		if strings.TrimSpace(segment) != "" {
			rendered, err := v.renderer.RenderBytes([]byte(segment))
			if err != nil {
				return fmt.Errorf("failed to render trailing markdown segment: %w", err)
			}
			fmt.Print(rendered)
		}
	}

	return nil
}

// renderSingleMermaidBlock renders a single mermaid block according to the
// current options, writing output inline.
func (v *SimpleViewer) renderSingleMermaidBlock(compiler *mermaid.Compiler, block renderer.MermaidBlock, index int, opts renderer.RenderOptions) error {
	// First render to SVG to get dimensions and SVG content.
	result, err := compiler.Render(block.Content)
	if err != nil {
		return err
	}
	if result.Error != nil {
		return result.Error
	}

	switch opts.MermaidMode {
	case "terminal":
		// Terminal mode: either inline image or ASCII preview, with the code
		// fence fully replaced by the visualization.
		if utils.SupportsInlineImages() {
			// Render to PNG for inline display.
			width := result.Width
			height := result.Height
			if width == 0 {
				width = 800
			}
			if height == 0 {
				height = 600
			}

			pngData, err := compiler.RenderToPNG(block.Content, width, height)
			if err != nil {
				return fmt.Errorf("failed to render PNG for inline display: %w", err)
			}

			fmt.Printf("📊 Mermaid Diagram (%s):\n", block.Type)
			protocol := utils.DetectImageProtocol()
			if err := utils.DisplayInlineImage(pngData, protocol); err != nil {
				return fmt.Errorf("failed to display inline image: %w", err)
			}
			fmt.Println()
		} else {
			// Fallback: ASCII preview box.
			preview := mermaid.GenerateASCIIPreview(block.Type, result.Width, result.Height)
			fmt.Print(preview)
		}

		// Optionally save SVG when requested.
		if opts.KeepMermaidFiles {
			filename := fmt.Sprintf("diagram-%d.svg", index+1)
			outputPath := filepath.Join(opts.MermaidOutDir, filename)
			if err := mermaid.SaveSVGToFile(result.SVG, outputPath); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to save SVG: %v\n", err)
			} else {
				fmt.Printf("  💾 Saved to: %s\n", outputPath)
			}
		}

	case "svg":
		filename := fmt.Sprintf("diagram-%d.svg", index+1)
		outputPath := filepath.Join(opts.MermaidOutDir, filename)
		if err := mermaid.SaveSVGToFile(result.SVG, outputPath); err != nil {
			return fmt.Errorf("failed to save SVG: %w", err)
		}
		// Show the original mermaid fence followed by a clickable file path.
		code := fmt.Sprintf("```mermaid\n%s\n```\n", strings.TrimSpace(block.Content))
		if rendered, err := v.renderer.RenderBytes([]byte(code)); err == nil {
			fmt.Print(rendered)
		} else {
			fmt.Print(code)
		}
		fmt.Printf("📁 Mermaid diagram %d %s\n", index+1, outputPath)

	case "png":
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
			return fmt.Errorf("failed to render PNG: %w", err)
		}

		filename := fmt.Sprintf("diagram-%d.png", index+1)
		outputPath := filepath.Join(opts.MermaidOutDir, filename)
		if err := mermaid.SavePNGToFile(pngData, outputPath); err != nil {
			return fmt.Errorf("failed to save PNG: %w", err)
		}
		code := fmt.Sprintf("```mermaid\n%s\n```\n", strings.TrimSpace(block.Content))
		if rendered, err := v.renderer.RenderBytes([]byte(code)); err == nil {
			fmt.Print(rendered)
		} else {
			fmt.Print(code)
		}
		fmt.Printf("📁 Mermaid diagram %d %s\n", index+1, outputPath)
	}

	return nil
}
