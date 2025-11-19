# Test Data Files

This directory contains test files for the go-tifpdf2png library.

## Files

### dummy.pdf

A simple single-page PDF file from W3C test resources.

- Source: [https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf]
- Pages: 1
- Used in: PDF conversion tests

### sample.tif (to be generated)

A multi-page TIFF file for testing TIFF conversion.

## Generating Test Files

If you need to create your own multi-page TIFF for testing, you can use ImageMagick:

```bash
# Create a simple 2-page TIFF
convert -size 200x200 xc:white -font helvetica -pointsize 24 \
    -draw "text 50,100 'Page 1'" page1.tif
convert -size 200x200 xc:white -font helvetica -pointsize 24 \
    -draw "text 50,100 'Page 2'" page2.tif
convert page1.tif page2.tif testdata/sample.tif
rm page1.tif page2.tif
```

Or download from public sources:

- TIFF Library test suite: [http://www.libtiff.org/images.html]
- TIFF samples: [https://github.com/rordenlab/dcm2niix/tree/master/TIFF]
