package utils

import (
	"encoding/base64"
	"fmt"
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
	switch protocol {
	case ProtocolITerm2:
		return displayITerm2Image(imageData)
	case ProtocolKitty:
		return displayKittyImage(imageData)
	case ProtocolSixel:
		return displaySixelImage(imageData)
	default:
		return fmt.Errorf("inline images not supported in this terminal")
	}
}

// displayITerm2Image displays an image using iTerm2's inline image protocol
// Protocol: ESC ] 1337 ; File=inline=1:<base64> BEL
func displayITerm2Image(imageData []byte) error {
	encoded := base64.StdEncoding.EncodeToString(imageData)
	
	// iTerm2 protocol: ESC ] 1337 ; File=inline=1:<base64> BEL
	fmt.Printf("\033]1337;File=inline=1:%s\a\n", encoded)
	
	return nil
}

// displayKittyImage displays an image using Kitty's graphics protocol
// Simplified version - just transmit and display
func displayKittyImage(imageData []byte) error {
	encoded := base64.StdEncoding.EncodeToString(imageData)
	
	// Kitty protocol: ESC _G<control data>;base64_data ESC \
	// a=T means transmit and display, f=100 means PNG format
	fmt.Printf("\033_Ga=T,f=100;%s\033\\\n", encoded)
	
	return nil
}

// displaySixelImage displays an image using Sixel protocol
// Note: This is a simplified implementation that outputs PNG as base64
// A full implementation would convert to actual Sixel format
func displaySixelImage(imageData []byte) error {
	// For Sixel, we need to convert PNG to Sixel format
	// This is complex, so for now we fall back to iTerm2 protocol
	// which Windows Terminal also supports as of recent versions
	return displayITerm2Image(imageData)
}

// SupportsInlineImages returns true if the current terminal supports inline images
func SupportsInlineImages() bool {
	return DetectImageProtocol() != ProtocolNone
}
