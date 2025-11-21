package main

import (
	"fmt"
	"os"

	"github.com/aquele_dinho/mdviewer/internal/pdf"
	"github.com/aquele_dinho/mdviewer/internal/renderer"
	"github.com/aquele_dinho/mdviewer/internal/utils"
	"github.com/aquele_dinho/mdviewer/internal/viewer"
	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags. Default is "dev" for local builds.
var version = "dev"

var (
	// Global flags
	style            string
	width            int
	noMermaid        bool
	openMermaid      bool
	exportPDF        string
	mermaidMode      string
	mermaidOutDir    string
	keepMermaidFiles bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "mdviewer [file]",
	Short: "A beautiful terminal markdown viewer with mermaid support",
	Long: `mdviewer is a cross-platform terminal markdown viewer that renders
markdown files with beautiful ANSI styling. It supports mermaid diagrams
and can export to PDF.

Examples:
  mdviewer README.md                    # View a markdown file
  cat file.md | mdviewer                # Read from stdin
  mdviewer file.md --style dark         # Use dark theme
  mdviewer file.md --export-pdf out.pdf # Export to PDF
`,
	Args: cobra.MaximumNArgs(1),
	RunE: runView,
}

func init() {
	// Wire version into Cobra (supports --version and in help output)
	rootCmd.Version = version

	// Add flags
	rootCmd.Flags().StringVarP(&style, "style", "s", "clean", "Color style: clean (default), auto, dark, light, or path to custom style")
	rootCmd.Flags().IntVarP(&width, "width", "w", 0, "Terminal width for word wrapping (0 = auto-detect)")
	rootCmd.Flags().BoolVar(&noMermaid, "no-mermaid", false, "Disable mermaid diagram detection")
	rootCmd.Flags().BoolVar(&openMermaid, "open-mermaid", false, "Open mermaid diagrams in browser automatically")
	rootCmd.Flags().StringVarP(&exportPDF, "export-pdf", "p", "", "Export to PDF file")
	rootCmd.Flags().StringVar(&mermaidMode, "mermaid-mode", "terminal", "Mermaid rendering mode: terminal (default), svg, png, url")
	rootCmd.Flags().StringVar(&mermaidOutDir, "mermaid-output-dir", os.TempDir(), "Directory for exported diagram files (default: system temp directory)")
	rootCmd.Flags().BoolVarP(&keepMermaidFiles, "keep-mermaid-files", "k", false, "Save Mermaid diagram files (SVG/PNG) to disk")
}

func runView(cmd *cobra.Command, args []string) error {
	// Determine input source
	var inputPath string
	if len(args) > 0 {
		inputPath = args[0]
	} else {
		// Check if stdin is piped
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			inputPath = "-" // stdin
		} else {
			return fmt.Errorf("no input file specified. Use 'mdviewer --help' for usage")
		}
	}

	// Auto-detect terminal width if not specified
	if width == 0 {
		width = utils.GetTerminalWidth()
	}

	// Create renderer
	rendererOpts := renderer.RenderOptions{
		Style:            style,
		Width:            width,
		NoMermaid:        noMermaid,
		MermaidMode:      mermaidMode,
		MermaidOutDir:    mermaidOutDir,
		KeepMermaidFiles: keepMermaidFiles,
	}

	mdRenderer, err := renderer.NewRenderer(rendererOpts)
	if err != nil {
		return fmt.Errorf("failed to create renderer: %w", err)
	}

	// Handle PDF export
	if exportPDF != "" {
		return exportToPDF(inputPath, exportPDF, mdRenderer)
	}

	// Handle mermaid diagram opening if requested
	if openMermaid && !noMermaid {
		// Read content to detect mermaid diagrams
		content, err := utils.ReadFile(inputPath)
		if err != nil {
			return fmt.Errorf("failed to read file for mermaid detection: %w", err)
		}
		
		// Detect mermaid blocks
		mermaidBlocks := renderer.DetectMermaidBlocks(string(content))
		
		if len(mermaidBlocks) > 0 {
			fmt.Fprintf(os.Stderr, "Opening %d mermaid diagram(s) in browser...\n", len(mermaidBlocks))
			
			for i, block := range mermaidBlocks {
				url := renderer.GenerateMermaidLiveURL(block)
				fmt.Fprintf(os.Stderr, "  %d. %s: %s\n", i+1, block.Type, url)
				
				if err := utils.OpenURL(url); err != nil {
					fmt.Fprintf(os.Stderr, "     Warning: failed to open URL: %v\n", err)
				}
			}
			fmt.Fprintln(os.Stderr)
		}
	}

	// Create viewer and display
	simpleViewer := viewer.NewSimpleViewer(mdRenderer)

	if inputPath == "-" {
		return simpleViewer.ViewStdin()
	}

	return simpleViewer.ViewFile(inputPath)
}

func exportToPDF(inputPath, outputPath string, r *renderer.Renderer) error {
	// Create PDF exporter
	exporter := pdf.NewExporter()

	// Export to PDF
	fmt.Fprintf(os.Stderr, "Generating PDF from %s...\n", inputPath)
	err := exporter.ExportFileToPDF(inputPath, outputPath)
	if err != nil {
		return fmt.Errorf("PDF export failed: %w", err)
	}

	fmt.Fprintf(os.Stderr, "PDF successfully exported to %s\n", outputPath)
	return nil
}
