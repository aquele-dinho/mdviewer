package pdf

import (
	"fmt"

	"github.com/aquele_dinho/mdviewer/internal/renderer"
	"github.com/aquele_dinho/mdviewer/internal/utils"
)

// Exporter handles PDF export from markdown
type Exporter struct {
	htmlRenderer *renderer.HTMLRenderer
	pdfGenerator *ChromeDPExporter
}

// NewExporter creates a new PDF exporter
func NewExporter() *Exporter {
	return &Exporter{
		htmlRenderer: renderer.NewHTMLRenderer(),
		pdfGenerator: NewChromeDPExporter(),
	}
}

// ExportToPDF converts markdown content to PDF and saves it to a file
func (e *Exporter) ExportToPDF(markdown string, outputPath string) error {
	// Convert markdown to HTML
	html, err := e.htmlRenderer.RenderToHTML(markdown)
	if err != nil {
		return fmt.Errorf("failed to render HTML: %w", err)
	}

	// Generate PDF from HTML
	pdfBytes, err := e.pdfGenerator.GeneratePDF(html)
	if err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Save PDF to file
	if err := utils.WriteFile(outputPath, pdfBytes); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	return nil
}

// ExportFileToPDF reads a markdown file and exports it to PDF
func (e *Exporter) ExportFileToPDF(inputPath string, outputPath string) error {
	// Read markdown file
	content, err := utils.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	return e.ExportToPDF(string(content), outputPath)
}
