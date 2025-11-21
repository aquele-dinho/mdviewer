package pdf

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// ChromeDPExporter uses headless Chrome to generate PDFs
type ChromeDPExporter struct {
	timeout time.Duration
}

// NewChromeDPExporter creates a new ChromeDP-based PDF exporter
func NewChromeDPExporter() *ChromeDPExporter {
	return &ChromeDPExporter{
		timeout: 30 * time.Second,
	}
}

// GeneratePDF generates a PDF from HTML content using headless Chrome
func (e *ChromeDPExporter) GeneratePDF(htmlContent string) ([]byte, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	// Create chromedp context
	allocCtx, allocCancel := chromedp.NewContext(ctx)
	defer allocCancel()

	var pdfBuffer []byte

	// Generate PDF
	err := chromedp.Run(allocCtx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Get the frame tree to set document content
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			// Set the HTML content
			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Generate PDF with print options
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(false).
				WithPaperWidth(8.5).   // US Letter width in inches
				WithPaperHeight(11.0). // US Letter height in inches
				WithMarginTop(0.4).
				WithMarginBottom(0.4).
				WithMarginLeft(0.4).
				WithMarginRight(0.4).
				Do(ctx)

			if err != nil {
				return err
			}

			pdfBuffer = buf
			return nil
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfBuffer, nil
}

// SetTimeout sets the timeout for PDF generation
func (e *ChromeDPExporter) SetTimeout(timeout time.Duration) {
	e.timeout = timeout
}
