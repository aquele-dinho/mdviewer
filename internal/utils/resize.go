package utils

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"golang.org/x/image/draw"
)

// ResizeImage resizes an image to the specified width, maintaining aspect ratio
func ResizeImage(imageData []byte, targetWidth int) ([]byte, error) {
	// Decode image
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Get original dimensions
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// If already smaller than target, return original
	if origWidth <= targetWidth {
		return imageData, nil
	}

	// Calculate new height maintaining aspect ratio
	aspectRatio := float64(origHeight) / float64(origWidth)
	targetHeight := int(float64(targetWidth) * aspectRatio)

	// Create new image with target dimensions
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Resize using bilinear interpolation
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	// Encode to PNG (use PNG for lossless quality)
	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return nil, fmt.Errorf("failed to encode resized image: %w", err)
	}

	// If resizing resulted in larger file, return original
	if buf.Len() > len(imageData) && format == "png" {
		return imageData, nil
	}

	return buf.Bytes(), nil
}
