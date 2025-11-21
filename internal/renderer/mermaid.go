package renderer

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

// MermaidBlock represents a detected mermaid diagram
type MermaidBlock struct {
	Type      string // diagram type: flowchart, sequence, etc.
	Content   string // raw mermaid code
	StartLine int    // line number where block starts
	EndLine   int    // line number where block ends
}

// MermaidDiagramTypes maps mermaid diagram types to descriptions
var MermaidDiagramTypes = map[string]string{
	"graph":       "Flowchart",
	"flowchart":   "Flowchart",
	"sequenceDiagram": "Sequence Diagram",
	"classDiagram":    "Class Diagram",
	"stateDiagram":    "State Diagram",
	"erDiagram":       "ER Diagram",
	"gantt":          "Gantt Chart",
	"pie":            "Pie Chart",
	"gitgraph":       "Git Graph",
	"journey":        "User Journey",
}

// DetectMermaidBlocks scans markdown content for mermaid code blocks
func DetectMermaidBlocks(content string) []MermaidBlock {
	var blocks []MermaidBlock
	
	// Regular expression to match mermaid code blocks
	// Matches: ```mermaid\n...content...\n```
	mermaidRegex := regexp.MustCompile("(?s)```mermaid\\s*\\n(.*?)\\n```")
	
	matches := mermaidRegex.FindAllStringSubmatchIndex(content, -1)
	
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}
		
		// Extract mermaid content
		mermaidContent := content[match[2]:match[3]]
		
		// Determine diagram type
		diagramType := detectDiagramType(mermaidContent)
		
		// Calculate line numbers
		startLine := strings.Count(content[:match[0]], "\n") + 1
		endLine := strings.Count(content[:match[1]], "\n") + 1
		
		blocks = append(blocks, MermaidBlock{
			Type:      diagramType,
			Content:   mermaidContent,
			StartLine: startLine,
			EndLine:   endLine,
		})
	}
	
	return blocks
}

// detectDiagramType attempts to identify the mermaid diagram type
func detectDiagramType(content string) string {
	content = strings.TrimSpace(content)
	lines := strings.Split(content, "\n")
	
	if len(lines) == 0 {
		return "Unknown"
	}
	
	// Check first non-empty line
	firstLine := strings.TrimSpace(lines[0])
	
	for keyword, typeName := range MermaidDiagramTypes {
		if strings.HasPrefix(firstLine, keyword) {
			return typeName
		}
	}
	
	return "Unknown"
}

// ExportMermaidBlock exports a mermaid block to a .mmd file
func ExportMermaidBlock(block MermaidBlock, outputPath string) error {
	// Will be implemented with file utils
	return nil
}

// GenerateMermaidLiveURL creates a mermaid.live URL for viewing
func GenerateMermaidLiveURL(block MermaidBlock) string {
	// Create a proper mermaid diagram with type prefix if needed
	diagramCode := strings.TrimSpace(block.Content)
	
	// Base64 encode the mermaid content
	encoded := base64.StdEncoding.EncodeToString([]byte(diagramCode))
	
	// Create the mermaid.live URL with pako compression format
	return fmt.Sprintf("https://mermaid.live/edit#pako:%s", encoded)
}

// GetMermaidInkURL creates a mermaid.ink URL for inline SVG rendering
func GetMermaidInkURL(block MermaidBlock) string {
	diagramCode := strings.TrimSpace(block.Content)
	encoded := base64.StdEncoding.EncodeToString([]byte(diagramCode))
	return fmt.Sprintf("https://mermaid.ink/img/%s", encoded)
}
