package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// TerminalImageProtocol represents different terminal inline image protocols
type TerminalImageProtocol string

const (
	ProtocolITerm2 TerminalImageProtocol = "iterm2"
	ProtocolKitty  TerminalImageProtocol = "kitty"
	ProtocolSixel  TerminalImageProtocol = "sixel"
	ProtocolNone   TerminalImageProtocol = "none"
)

// DetectImageProtocol detects which inline image protocol the terminal supports
func DetectImageProtocol() TerminalImageProtocol {
	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")
	wtSession := os.Getenv("WT_SESSION")
	
	// Check for iTerm2
	if termProgram == "iTerm.app" {
		return ProtocolITerm2
	}
	
	// Check for Warp (supports iTerm2 protocol)
	if termProgram == "WarpTerminal" {
		return ProtocolITerm2
	}
	
	// Check for Kitty
	if term == "xterm-kitty" || termProgram == "kitty" {
		return ProtocolKitty
	}
	
	// Check for Windows Terminal (supports Sixel in v1.22+)
	if wtSession != "" {
		return ProtocolSixel
	}
	
	// VSCode terminal (supports iTerm2 protocol)
	if termProgram == "vscode" {
		return ProtocolITerm2
	}
	
	return ProtocolNone
}

// DisplayInlineImage displays an image inline using the terminal's supported protocol
func DisplayInlineImage(imageData []byte, protocol TerminalImageProtocol) error {
	return DisplayInlineImageTo(os.Stdout, imageData, protocol)
}

// DisplayInlineImageTo writes an inline image to the given writer using the terminal's supported protocol
func DisplayInlineImageTo(w io.Writer, imageData []byte, protocol TerminalImageProtocol) error {
	switch protocol {
	case ProtocolITerm2:
		return writeITerm2Image(w, imageData)
	case ProtocolKitty:
		return writeKittyImage(w, imageData)
	case ProtocolSixel:
		return writeITerm2Image(w, imageData) // Sixel falls back to iTerm2
	default:
		return fmt.Errorf("inline images not supported in this terminal")
	}
}

// writeITerm2Image writes an image using iTerm2's inline image protocol to the given writer
func writeITerm2Image(w io.Writer, imageData []byte) error {
	encoded := base64.StdEncoding.EncodeToString(imageData)
	_, err := fmt.Fprintf(w, "\033]1337;File=inline=1:%s\a\n", encoded)
	return err
}

// writeKittyImage writes an image using Kitty's graphics protocol to the given writer
func writeKittyImage(w io.Writer, imageData []byte) error {
	encoded := base64.StdEncoding.EncodeToString(imageData)
	_, err := fmt.Fprintf(w, "\033_Ga=T,f=100;%s\033\\\n", encoded)
	return err
}

// SupportsInlineImages returns true if the current terminal supports inline images
func SupportsInlineImages() bool {
	return DetectImageProtocol() != ProtocolNone
}
