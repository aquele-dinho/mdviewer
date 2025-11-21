#!/bin/bash
set -e

VERSION=${1:-"0.1.0"}
OUTPUT_DIR="./dist"

echo "Building mdviewer v${VERSION} for multiple platforms..."

# Clean and create output directory
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Build for macOS (Intel)
echo "Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "$OUTPUT_DIR/mdviewer-darwin-amd64" ./cmd/mdviewer

# Build for macOS (Apple Silicon)
echo "Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "$OUTPUT_DIR/mdviewer-darwin-arm64" ./cmd/mdviewer

# Build for Linux (AMD64)
echo "Building for Linux (AMD64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$OUTPUT_DIR/mdviewer-linux-amd64" ./cmd/mdviewer

# Build for Linux (ARM64)
echo "Building for Linux (ARM64)..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o "$OUTPUT_DIR/mdviewer-linux-arm64" ./cmd/mdviewer

# Build for Windows (AMD64)
echo "Building for Windows (AMD64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "$OUTPUT_DIR/mdviewer-windows-amd64.exe" ./cmd/mdviewer

# Display results
echo ""
echo "Build complete! Binaries created in $OUTPUT_DIR:"
ls -lh "$OUTPUT_DIR"

# Calculate sizes
echo ""
echo "Binary sizes:"
du -h "$OUTPUT_DIR"/*

echo ""
echo "Note: All binaries require Chrome/Chromium to be installed on target systems"
echo "for local Mermaid diagram rendering. Without it, URL mode will be used as fallback."
I