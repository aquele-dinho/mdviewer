package renderer

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// ImageBlock represents a markdown image found in the content
type ImageBlock struct {
	AltText   string // Alt text from ![alt](path)
	Path      string // Image file path or URL
	Width     int    // Optional width hint in pixels (from Obsidian |width syntax)
	StartLine int    // Line number where image appears (1-indexed)
	EndLine   int    // Same as StartLine for images (single line)
}

// Supported image formats for inline display
var supportedImageFormats = []string{".png", ".jpg", ".jpeg", ".gif", ".webp"}

// Image detection regex: ![alt text](path)
var imageRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)

// DetectImageBlocks finds all markdown images in the content
func DetectImageBlocks(content string) []ImageBlock {
	var blocks []ImageBlock
	lines := strings.Split(content, "\n")
	
	for lineNum, line := range lines {
		matches := imageRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) < 3 {
				continue
			}
			
			altText := match[1]
			path := strings.TrimSpace(match[2])
			
			// Skip URLs (http://, https://, etc.)
			if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
				continue
			}
			
			// Check if it's a supported format
			if !isSupportedImageFormat(path) {
				continue
			}
			
			// Parse width from alt text if present (format: "name|width=400")
			width := 0
			cleanAlt := altText
			if strings.Contains(altText, "|width=") {
				parts := strings.SplitN(altText, "|width=", 2)
				if len(parts) == 2 {
					cleanAlt = parts[0]
					if w, err := strconv.Atoi(parts[1]); err == nil {
						width = w
					}
				}
			}
			
			blocks = append(blocks, ImageBlock{
				AltText:   cleanAlt,
				Path:      path,
				Width:     width,
				StartLine: lineNum + 1, // 1-indexed
				EndLine:   lineNum + 1,
			})
		}
	}
	
	return blocks
}

// isSupportedImageFormat checks if the file has a supported image extension
func isSupportedImageFormat(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, supported := range supportedImageFormats {
		if ext == supported {
			return true
		}
	}
	return false
}
