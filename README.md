# go-tifpdf2png

A Go library and CLI tool for converting TIFF and PDF files to PNG images with intelligent cropping and background normalization.

## Features

- **Multi-format Support**: Converts both TIFF (.tif, .tiff) and PDF (.pdf) files to PNG images
- **Automatic Format Detection**: Intelligently routes to the appropriate converter based on file extension
- **Intelligent Cropping**: Automatically detects and crops to content boundaries
- **Background Detection**: Detects and normalizes dark/light backgrounds for optimal contrast
- **Detailed Metadata**: Returns comprehensive image details including dimensions, crop information, and page counts
- **Clean Output**: Status messages to stderr, JSON data to stdout for easy piping

## Installation

### As a Library

```bash
go get github.com/dhushon/go-tifpdf2png
```

### As a CLI Tool

```bash
go install github.com/dhushon/go-tifpdf2png/cmd/converttifpdf@latest
```

## Usage

### CLI Tool

```bash
# Convert a TIFF file
converttifpdf document.tif

# Convert a PDF file
converttifpdf invoice.pdf

# Capture JSON output
converttifpdf payment.pdf > metadata.json

# Save both output and errors
converttifpdf report.tif > metadata.json 2> conversion.log
```

### As a Library

```go
package main

import (
    "fmt"
    "github.com/dhushon/go-tifpdf2png"
)

func main() {
    // Convert TIFF to PNG
    tiffDetails, err := tifpdf2png.ConvertTiffToPngWithImageDetails(
        "document.tif",
        "/output/path",
        "document-page-",
    )
    if err != nil {
        panic(err)
    }

    // Convert PDF to PNG
    pdfDetails, err := tifpdf2png.ConvertPdfToPngWithImageDetails(
        "invoice.pdf",
        "/output/path",
        "invoice-page-",
    )
    if err != nil {
        panic(err)
    }

    fmt.Printf("Converted %d TIFF pages and %d PDF pages\n", 
        len(tiffDetails), len(pdfDetails))
}
```

## API Documentation

### Main Functions

#### `ConvertTiffToPngWithImageDetails(tiffFilename, destpath, prefix string) ([]*ImageDetail, error)`

Converts a TIFF file to PNG images with detailed metadata.

**Parameters:**
- `tiffFilename`: Path to the input TIFF file
- `destpath`: Directory for output PNG files
- `prefix`: Filename prefix for output files (empty string for auto-generation)

**Returns:** Slice of `ImageDetail` structures containing metadata for each page

#### `ConvertPdfToPngWithImageDetails(pdfFilename, destpath, prefix string) ([]*ImageDetail, error)`

Converts a PDF file to PNG images with detailed metadata.

**Parameters:**
- `pdfFilename`: Path to the input PDF file
- `destpath`: Directory for output PNG files  
- `prefix`: Filename prefix for output files (empty string for auto-generation)

**Returns:** Slice of `ImageDetail` structures containing metadata for each page

### Data Structures

#### `ImageDetail`

```go
type ImageDetail struct {
    ActualType string      // The actual type of the image (e.g., "png")
    Page       int         // Page number (1-based)
    Pages      int         // Total number of pages
    URL        string      // Path to the output PNG file
    Width      int         // Width of the image in pixels
    Height     int         // Height of the image in pixels
    Format     string      // Image format (e.g., "png")
    Quality    float64     // Quality metric (0-100)
    CropDetail *CropDetail // Crop information if cropping occurred
}
```

#### `CropDetail`

```go
type CropDetail struct {
    OffsetX        int // X offset of the crop from original
    OffsetY        int // Y offset of the crop from original
    OriginalWidth  int // Original width before cropping
    OriginalHeight int // Original height before cropping
    CroppedWidth   int // Width after cropping
    CroppedHeight  int // Height after cropping
}
```

## Image Processing Pipeline

Both TIFF and PDF files undergo the same processing:

1. **Rendering**: 
   - TIFF: Direct multi-frame extraction
   - PDF: Page rendering at 150 DPI

2. **Cropping**: Automatic detection and removal of empty margins using content boundary analysis

3. **Background Normalization**: 
   - Analyzes corner and edge pixels to detect background type
   - Converts to white background with black content for optimal contrast
   - Handles both light and dark source backgrounds

4. **Metadata Generation**: 
   - Records original and cropped dimensions
   - Tracks crop offsets for coordinate mapping
   - Includes page numbering and total page count

## Output Format

### PNG Files
Generated with the naming pattern: `<prefix><page-number>.png`

Example: `document-page-0.png`, `document-page-1.png`, etc.

### JSON Metadata

```json
[
  {
    "actual_type": "png",
    "page": 1,
    "pages": 3,
    "url": "/path/to/output/document-page-0.png",
    "width": 2550,
    "height": 3300,
    "format": "png",
    "quality": 95,
    "crop_detail": {
      "offset_x": 100,
      "offset_y": 150,
      "original_width": 2650,
      "original_height": 3400,
      "cropped_width": 2550,
      "cropped_height": 3300
    }
  }
]
```

## Dependencies

- `github.com/gen2brain/go-fitz` - PDF rendering
- `github.com/dhushon/tiff` - TIFF decoding
- `github.com/disintegration/imaging` - Image processing

## License

Copyright © 2025 GH Inner

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

Originally developed as part of the [NGPA-UX](https://github.com/gh-inner/ngpa-ux) project for healthcare payment document processing. PDF conversion adapted from [platform-commons](https://github.com/gh-inner/platform-commons).
