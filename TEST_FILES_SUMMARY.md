# Test Files Summary

## Overview
Created comprehensive test suite for `go-tifpdf2png` with automatic test file downloading from public sources.

## Test Files

### Downloaded Automatically (on first test run)

1. **UTM2GTIF.TIF** (TIFF)
   - Source: https://download.osgeo.org/geotiff/samples/spot/chicago/UTM2GTIF.TIF
   - Type: GeoTIFF image
   - Size: ~500KB
   - Used for: TIFF conversion testing

2. **dummy.pdf** (PDF)
   - Source: https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf
   - Type: Simple PDF document
   - Pages: 1
   - Used for: PDF conversion testing

## Test Suite (`processing_test.go`)

### Tests Included

1. **TestConvertTiffToPngWithImageDetails**
   - Downloads TIFF test file
   - Converts to PNG
   - Verifies metadata (dimensions, page count, format)
   - Checks output file existence

2. **TestConvertPdfToPngWithImageDetails**
   - Downloads PDF test file
   - Converts to PNG
   - Verifies metadata
   - Checks output file existence

3. **TestConvertTiffToPngWithCropping**
   - Tests automatic cropping functionality
   - Logs crop details when applied
   - Reports cropping statistics

4. **TestConvertPdfToPngWithEmptyPrefix**
   - Tests auto-generated filename feature
   - Verifies timestamp-based naming

5. **TestInvalidFile**
   - Tests error handling for non-existent files
   - Ensures proper error reporting

### Features

- **Automatic Download**: Test files are downloaded once and cached in `./testdata`
- **Skip on Failure**: Tests skip gracefully if downloads fail
- **Temporary Output**: Uses `t.TempDir()` for clean test isolation
- **Comprehensive Validation**: Checks file existence, metadata, and dimensions
- **Detailed Logging**: Provides verbose output for debugging

## Running Tests

```bash
# Run all tests
go test -v

# Run specific test
go test -v -run TestConvertTiffToPng

# Run with coverage
go test -v -cover

# Clean test cache and re-run
go clean -testcache && go test -v
```

## Test Results

All tests passing ✓
- TIFF conversion: Working
- PDF conversion: Working
- Cropping detection: Working
- Error handling: Working
- Auto-naming: Working

## Notes

- Test files are downloaded from reliable public sources
- Files are cached in `./testdata` after first download
- Tests use web GET approach - no manual file management needed
- GeoTIFF file works well for testing (standard TIFF format)
