package renderer

import "sort"

// BlockType identifies the type of content block
type BlockType int

const (
	BlockTypeMermaid BlockType = iota
	BlockTypeImage
)

// ContentBlock represents any block of special content (Mermaid or Image)
type ContentBlock struct {
	Type      BlockType
	Mermaid   *MermaidBlock // Non-nil if Type == BlockTypeMermaid
	Image     *ImageBlock   // Non-nil if Type == BlockTypeImage
	StartLine int           // 1-indexed line number
	EndLine   int           // 1-indexed line number (exclusive for next segment)
}

// DetectContentBlocks finds all special content blocks (images and mermaid) in markdown
func DetectContentBlocks(content string) []ContentBlock {
	var blocks []ContentBlock
	
	// Detect Mermaid blocks
	mermaidBlocks := DetectMermaidBlocks(content)
	for i := range mermaidBlocks {
		blocks = append(blocks, ContentBlock{
			Type:      BlockTypeMermaid,
			Mermaid:   &mermaidBlocks[i],
			StartLine: mermaidBlocks[i].StartLine,
			EndLine:   mermaidBlocks[i].EndLine,
		})
	}
	
	// Detect image blocks
	imageBlocks := DetectImageBlocks(content)
	for i := range imageBlocks {
		blocks = append(blocks, ContentBlock{
			Type:      BlockTypeImage,
			Image:     &imageBlocks[i],
			StartLine: imageBlocks[i].StartLine,
			EndLine:   imageBlocks[i].EndLine,
		})
	}
	
	// Sort by start line
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].StartLine < blocks[j].StartLine
	})
	
	return blocks
}
