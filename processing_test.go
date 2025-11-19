package tifpdf2png

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

const (
	// Multi-page TIFF sample from libtiff.org test images
	multiPageTiffURL = "https://download.osgeo.org/geotiff/samples/spot/chicago/UTM2GTIF.TIF"

	// Multi-page PDF sample from PDF specification examples
	multiPagePDFURL = "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"
)

// downloadTestFile downloads a file from a URL to the specified path
func downloadTestFile(t *testing.T, url, destPath string) error {
	t.Helper()

	// Create testdata directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// Skip download if file already exists
	if _, err := os.Stat(destPath); err == nil {
		t.Logf("Test file already exists: %s", destPath)
		return nil
	}

	t.Logf("Downloading test file from %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return io.ErrUnexpectedEOF
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func TestConvertTiffToPngWithImageDetails(t *testing.T) {
	testTiffPath := filepath.Join("testdata", "UTM2GTIF.TIF")

	// Download test file
	if err := downloadTestFile(t, multiPageTiffURL, testTiffPath); err != nil {
		t.Skipf("Could not download test TIFF file: %v", err)
	}

	// Create temporary output directory
	outputDir := t.TempDir()

	// Run conversion
	details, err := ConvertTiffToPngWithImageDetails(
		testTiffPath,
		outputDir,
		"test-tiff-",
	)

	if err != nil {
		t.Fatalf("ConvertTiffToPngWithImageDetails failed: %v", err)
	}

	// Verify results
	if len(details) == 0 {
		t.Fatal("Expected at least one page, got none")
	}

	t.Logf("Successfully converted %d TIFF page(s)", len(details))

	for i, detail := range details {
		t.Logf("Page %d: %dx%d, file: %s", detail.Page, detail.Width, detail.Height, detail.URL)

		// Verify file exists
		if _, err := os.Stat(detail.URL); err != nil {
			t.Errorf("Output file does not exist: %s", detail.URL)
		}

		// Verify metadata
		if detail.Page != i+1 {
			t.Errorf("Expected page %d, got %d", i+1, detail.Page)
		}
		if detail.Pages != len(details) {
			t.Errorf("Expected total pages %d, got %d", len(details), detail.Pages)
		}
		if detail.Width <= 0 || detail.Height <= 0 {
			t.Errorf("Invalid dimensions: %dx%d", detail.Width, detail.Height)
		}
		if detail.Format != "png" {
			t.Errorf("Expected format 'png', got '%s'", detail.Format)
		}
	}
}

func TestConvertPdfToPngWithImageDetails(t *testing.T) {
	testPdfPath := filepath.Join("testdata", "dummy.pdf")

	// Download test file
	if err := downloadTestFile(t, multiPagePDFURL, testPdfPath); err != nil {
		t.Skipf("Could not download test PDF file: %v", err)
	}

	// Create temporary output directory
	outputDir := t.TempDir()

	// Run conversion
	details, err := ConvertPdfToPngWithImageDetails(
		testPdfPath,
		outputDir,
		"test-pdf-",
	)

	if err != nil {
		t.Fatalf("ConvertPdfToPngWithImageDetails failed: %v", err)
	}

	// Verify results
	if len(details) == 0 {
		t.Fatal("Expected at least one page, got none")
	}

	t.Logf("Successfully converted %d PDF page(s)", len(details))

	for i, detail := range details {
		t.Logf("Page %d: %dx%d, file: %s", detail.Page, detail.Width, detail.Height, detail.URL)

		// Verify file exists
		if _, err := os.Stat(detail.URL); err != nil {
			t.Errorf("Output file does not exist: %s", detail.URL)
		}

		// Verify metadata
		if detail.Page != i+1 {
			t.Errorf("Expected page %d, got %d", i+1, detail.Page)
		}
		if detail.Pages != len(details) {
			t.Errorf("Expected total pages %d, got %d", len(details), detail.Pages)
		}
		if detail.Width <= 0 || detail.Height <= 0 {
			t.Errorf("Invalid dimensions: %dx%d", detail.Width, detail.Height)
		}
		if detail.Format != "png" {
			t.Errorf("Expected format 'png', got '%s'", detail.Format)
		}
	}
}

func TestConvertTiffToPngWithCropping(t *testing.T) {
	testTiffPath := filepath.Join("testdata", "UTM2GTIF.TIF")

	// Download test file
	if err := downloadTestFile(t, multiPageTiffURL, testTiffPath); err != nil {
		t.Skipf("Could not download test TIFF file: %v", err)
	}

	// Create temporary output directory
	outputDir := t.TempDir()

	// Run conversion
	details, err := ConvertTiffToPngWithImageDetails(
		testTiffPath,
		outputDir,
		"",
	)

	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	// Check if any cropping occurred
	croppedPages := 0
	for _, detail := range details {
		if detail.CropDetail != nil {
			croppedPages++
			t.Logf("Page %d was cropped: %dx%d -> %dx%d (offset: %d,%d)",
				detail.Page,
				detail.CropDetail.OriginalWidth,
				detail.CropDetail.OriginalHeight,
				detail.CropDetail.CroppedWidth,
				detail.CropDetail.CroppedHeight,
				detail.CropDetail.OffsetX,
				detail.CropDetail.OffsetY,
			)
		}
	}

	t.Logf("%d out of %d pages were cropped", croppedPages, len(details))
}

func TestConvertPdfToPngWithEmptyPrefix(t *testing.T) {
	testPdfPath := filepath.Join("testdata", "dummy.pdf")

	// Download test file
	if err := downloadTestFile(t, multiPagePDFURL, testPdfPath); err != nil {
		t.Skipf("Could not download test PDF file: %v", err)
	}

	// Create temporary output directory
	outputDir := t.TempDir()

	// Run conversion with empty prefix
	details, err := ConvertPdfToPngWithImageDetails(
		testPdfPath,
		outputDir,
		"",
	)

	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if len(details) == 0 {
		t.Fatal("Expected at least one page")
	}

	// Verify that auto-generated filenames are used
	for _, detail := range details {
		if detail.URL == "" {
			t.Error("Expected non-empty URL in image detail")
		}
		t.Logf("Generated filename: %s", filepath.Base(detail.URL))
	}
}

func TestInvalidFile(t *testing.T) {
	outputDir := t.TempDir()

	// Test with non-existent file
	_, err := ConvertTiffToPngWithImageDetails(
		"nonexistent.tiff",
		outputDir,
		"test-",
	)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	// Test with non-existent PDF
	_, err = ConvertPdfToPngWithImageDetails(
		"nonexistent.pdf",
		outputDir,
		"test-",
	)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
