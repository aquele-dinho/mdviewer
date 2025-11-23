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
	basePath string // Base directory for resolving relative image paths
}

// NewSimpleViewer creates a new simple viewer
func NewSimpleViewer(r *renderer.Renderer) *SimpleViewer {
	return &SimpleViewer{
		renderer: r,
		basePath: ".",
	}
}

// SetBasePath sets the base directory for resolving relative image paths
func (v *SimpleViewer) SetBasePath(path string) {
	v.basePath = path
}

// View renders and displays markdown content
func (v *SimpleViewer) View(content []byte) error {
	opts := v.renderer.GetOptions()

	// Always use inline content rendering to handle both images and mermaid.
	// The function will detect if there are any special blocks to handle.
	if err := v.renderWithInlineContent(content, opts); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: inline content rendering failed: %v\n", err)
	}

	return nil
}

// ViewFile reads and displays a markdown file
func (v *SimpleViewer) ViewFile(path string) error {
	// Set base path to the directory containing the markdown file
	if absPath, err := filepath.Abs(path); err == nil {
		v.basePath = filepath.Dir(absPath)
	}

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

// renderWithInlineContent renders the markdown in segments separated by
// special content blocks (images, mermaid) so they appear inline at their original positions.
func (v *SimpleViewer) renderWithInlineContent(content []byte, opts renderer.RenderOptions) error {
	// First preprocess links so we detect images after Obsidian syntax conversion
	text := v.renderer.PreprocessLinks(string(content))
	blocks := renderer.DetectContentBlocks(text)
	if len(blocks) == 0 {
		// No special content: just render normally with preprocessed content.
		rendered, err := v.renderer.RenderBytes([]byte(text))
		if err != nil {
			return fmt.Errorf("failed to render content: %w", err)
		}
		fmt.Print(rendered)
		return nil
	}

	lines := strings.Split(text, "\n")

	// Only create mermaid compiler if we have mermaid blocks
	var compiler *mermaid.Compiler
	hasMermaid := false
	for _, block := range blocks {
		if block.Type == renderer.BlockTypeMermaid {
			hasMermaid = true
			break
		}
	}

	if hasMermaid {
		var err error
		compiler, err = mermaid.NewCompiler()
		if err != nil {
			// Fall back to plain rendered markdown if compiler cannot be created.
			rendered, rErr := v.renderer.RenderBytes([]byte(text))
			if rErr != nil {
				return fmt.Errorf("failed to render content without mermaid: %w", rErr)
			}
			fmt.Print(rendered)
			return fmt.Errorf("failed to create mermaid compiler: %w", err)
		}
		defer compiler.Close()
	}

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

		// Render the content block inline based on its type.
		if err := v.renderSingleContentBlock(compiler, block, i, opts); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to render content block %d: %v\n", i+1, err)
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

// renderSingleContentBlock dispatches rendering based on block type
func (v *SimpleViewer) renderSingleContentBlock(compiler *mermaid.Compiler, block renderer.ContentBlock, index int, opts renderer.RenderOptions) error {
	switch block.Type {
	case renderer.BlockTypeMermaid:
		return v.renderSingleMermaidBlock(compiler, *block.Mermaid, index, opts)
	case renderer.BlockTypeImage:
		return v.renderSingleImageBlock(*block.Image, index)
	default:
		return fmt.Errorf("unknown block type: %v", block.Type)
	}
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

// renderSingleImageBlock renders a single image block inline
func (v *SimpleViewer) renderSingleImageBlock(block renderer.ImageBlock, index int) error {
	// Check if terminal supports inline images
	if !utils.SupportsInlineImages() {
		// Fallback: let Glamour render the image as text
		markdownImg := fmt.Sprintf("![%s](%s)", block.AltText, block.Path)
		rendered, err := v.renderer.RenderBytes([]byte(markdownImg))
		if err != nil {
			return fmt.Errorf("failed to render image placeholder: %w", err)
		}
		fmt.Print(rendered)
		return nil
	}

	// Resolve image path
	imgPath := block.Path
	if !filepath.IsAbs(imgPath) {
		imgPath = filepath.Join(v.basePath, imgPath)
	}

	// Load image file
	imageData, err := os.ReadFile(imgPath)
	if err != nil {
		// File not found or error reading - fall back to text representation
		markdownImg := fmt.Sprintf("![%s](%s)", block.AltText, block.Path)
		rendered, err := v.renderer.RenderBytes([]byte(markdownImg))
		if err != nil {
			return fmt.Errorf("failed to render image placeholder: %w", err)
		}
		fmt.Print(rendered)
		return nil
	}

	// Resize image if width specified
	if block.Width > 0 {
		resized, err := utils.ResizeImage(imageData, block.Width)
		if err != nil {
			// If resize fails, use original
			fmt.Fprintf(os.Stderr, "Warning: failed to resize image: %v\n", err)
		} else {
			imageData = resized
		}
	}

	// Display inline image
	if block.AltText != "" {
		fmt.Printf("🖼️  %s:\n", block.AltText)
	} else {
		fmt.Printf("🖼️  Image:\n")
	}

	protocol := utils.DetectImageProtocol()
	if err := utils.DisplayInlineImage(imageData, protocol); err != nil {
		return fmt.Errorf("failed to display inline image: %w", err)
	}
	fmt.Println()

	return nil
}
