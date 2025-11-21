package mermaid

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SVGResult represents a rendered SVG diagram
type SVGResult struct {
	SVG    string
	Width  int
	Height int
	Error  error
}

// ExtractSVGDimensions extracts width and height from SVG content
func ExtractSVGDimensions(svg string) (width, height int) {
	// Default dimensions
	width, height = 800, 600

	// Try to extract from viewBox
	viewBoxRegex := regexp.MustCompile(`viewBox="[^"]*\s+([0-9.]+)\s+([0-9.]+)"`)
	if matches := viewBoxRegex.FindStringSubmatch(svg); len(matches) == 3 {
		fmt.Sscanf(matches[1], "%d", &width)
		fmt.Sscanf(matches[2], "%d", &height)
		return
	}

	// Try to extract from width/height attributes
	widthRegex := regexp.MustCompile(`width="([0-9.]+)"`)
	heightRegex := regexp.MustCompile(`height="([0-9.]+)"`)

	if matches := widthRegex.FindStringSubmatch(svg); len(matches) == 2 {
		fmt.Sscanf(matches[1], "%d", &width)
	}
	if matches := heightRegex.FindStringSubmatch(svg); len(matches) == 2 {
		fmt.Sscanf(matches[1], "%d", &height)
	}

	return
}

// SaveSVGToFile saves SVG content to a file
func SaveSVGToFile(svg, outputPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write SVG to file
	if err := os.WriteFile(outputPath, []byte(svg), 0644); err != nil {
		return fmt.Errorf("failed to write SVG file: %w", err)
	}

	return nil
}

// SavePNGToFile saves PNG bytes to a file
func SavePNGToFile(pngData []byte, outputPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write PNG to file
	if err := os.WriteFile(outputPath, pngData, 0644); err != nil {
		return fmt.Errorf("failed to write PNG file: %w", err)
	}

	return nil
}

// CleanSVG removes unnecessary whitespace and prettifies SVG
func CleanSVG(svg string) string {
	// Remove excessive whitespace
	svg = strings.TrimSpace(svg)
	
	// Ensure it starts with <?xml or <svg
	if !strings.HasPrefix(svg, "<?xml") && !strings.HasPrefix(svg, "<svg") {
		return svg
	}

	return svg
}

// GenerateASCIIPreview generates a simple ASCII representation of diagram info
func GenerateASCIIPreview(diagramType string, width, height int) string {
	var preview strings.Builder
	
	// Top border
	preview.WriteString("┌" + strings.Repeat("─", 50) + "┐\n")
	
	// Content
	preview.WriteString(fmt.Sprintf("│ 📊 Mermaid Diagram: %-32s │\n", diagramType))
	preview.WriteString(fmt.Sprintf("│ 📐 Dimensions: %dx%d px %-24s │\n", width, height, ""))
	preview.WriteString("│ ✅ Rendered locally                              │\n")
	
	// Bottom border
	preview.WriteString("└" + strings.Repeat("─", 50) + "┘\n")
	
	return preview.String()
}
