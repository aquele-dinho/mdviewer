# Changelog

All notable changes to mdviewer will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2025-11-23

### Added
- **Inline image display**: Images now render directly in supported terminals (Warp, iTerm2, Kitty, Windows Terminal, VSCode, WezTerm)
- **Image resizing support**: Obsidian-style sizing syntax `![[image.png|400]]` resizes images to specified width while maintaining aspect ratio
- **Obsidian compatibility**: Full support for Obsidian wiki-links (`[[page]]`, `[[page|label]]`) and image embeds (`![[image.png]]`)
- **Multiple image format support**: PNG, JPEG, GIF, and WebP images
- **Smart path resolution**: Relative image paths are resolved from the markdown file's directory
- **Graceful fallback**: Terminals without inline image support display text placeholders
- Automatic protocol detection for iTerm2, Kitty, and Sixel image protocols

### Changed
- Refactored viewer to use unified content block detection for both images and Mermaid diagrams
- Images and Mermaid diagrams now render in document order (segmented rendering)
- Preprocessing now handles Obsidian syntax before content detection

### Technical
- Added `internal/renderer/images.go` for image block detection
- Added `internal/renderer/blocks.go` for unified content block management
- Added `internal/utils/resize.go` for image resizing with bilinear interpolation
- Extended `internal/utils/termimg.go` for inline image display
- Added `golang.org/x/image/draw` dependency for image processing

## [0.1.0] - 2024-XX-XX

### Added
- Initial release with markdown rendering
- Local Mermaid diagram rendering using chromedp
- Multiple rendering modes (terminal, SVG, PNG, URL)
- PDF export functionality
- Multiple color themes (auto, dark, light, clean)
- Stdin support
- Terminal width auto-detection
- Cross-platform support (macOS, Linux, Windows)

[0.2.0]: https://github.com/aquele_dinho/mdviewer/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/aquele_dinho/mdviewer/releases/tag/v0.1.0
