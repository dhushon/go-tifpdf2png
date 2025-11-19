// Package main provides a CLI tool for converting TIFF and PDF files to PNG images.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gh-inner/go-tifpdf2png"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s --version\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nConverts a TIFF or PDF file to PNG images in the current working directory\n")
		fmt.Fprintf(os.Stderr, "and outputs ImageDetails to stdout as JSON.\n")
		fmt.Fprintf(os.Stderr, "\nSupported formats: .tif, .tiff, .pdf\n")
		os.Exit(1)
	}

	// Handle version flag
	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Printf("go-tifpdf2png convert %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", date)
		os.Exit(0)
	}

	inputFile := os.Args[1]

	// Check if file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: File '%s' does not exist\n", inputFile)
		os.Exit(1)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Extract base filename without extension for prefix
	baseName := filepath.Base(inputFile)
	ext := filepath.Ext(baseName)
	prefix := baseName[:len(baseName)-len(ext)] + "-page-"

	// Detect file type and route to appropriate converter
	extLower := strings.ToLower(ext)
	var imageDetails []*tifpdf2png.ImageDetail

	switch extLower {
	case ".pdf":
		imageDetails, err = tifpdf2png.ConvertPdfToPngWithImageDetails(inputFile, cwd, prefix)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting PDF to PNG: %v\n", err)
			os.Exit(1)
		}
	case ".tif", ".tiff":
		imageDetails, err = tifpdf2png.ConvertTiffToPngWithImageDetails(inputFile, cwd, prefix)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting TIFF to PNG: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unsupported file type '%s'\n", ext)
		fmt.Fprintf(os.Stderr, "Supported formats: .tif, .tiff, .pdf\n")
		os.Exit(1)
	}

	// Output ImageDetails as JSON to stdout
	output, err := json.MarshalIndent(imageDetails, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling ImageDetails to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))

	// Print summary to stderr so it doesn't interfere with JSON output
	fmt.Fprintf(os.Stderr, "\n✓ Converted %d page(s) from %s to PNG\n", len(imageDetails), inputFile)
	fmt.Fprintf(os.Stderr, "✓ Output files in: %s\n", cwd)
}
