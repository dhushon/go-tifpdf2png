// Package tifpdf2png provides utilities for converting TIFF and PDF files to PNG images.
package tifpdf2png

import (
	"fmt"
	"image"
	"log/slog"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gen2brain/go-fitz"
)

// ConvertPdfToPngWithImageDetails converts PDF to PNG and returns ImageDetail slice
func ConvertPdfToPngWithImageDetails(pdfFilename string, destpath string, prefix string) ([]*ImageDetail, error) {
	doc, err := fitz.New(pdfFilename)
	if err != nil {
		slog.Error("ConvertPdfToPngWithImageDetails: load PDF file", "error", err)
		return nil, err
	}
	defer func() {
		if err := doc.Close(); err != nil {
			slog.Warn("Failed to close PDF document", "error", err)
		}
	}()

	pageCount := doc.NumPage()
	if pageCount == 0 {
		slog.Error("ConvertPdfToPngWithImageDetails: no pages found in PDF file")
		return nil, fmt.Errorf("no pages found in PDF file")
	}

	if prefix == "" {
		prefix = time.Now().Format("20060102-150405-")
	}
	if destpath[len(destpath)-1:] != "/" {
		destpath = destpath + "/"
	}

	var imageDetails []*ImageDetail

	for pageNum := 0; pageNum < pageCount; pageNum++ {
		img, err := doc.Image(pageNum)
		if err != nil {
			slog.Error("ConvertPdfToPngWithImageDetails: render page",
				"page", pageNum,
				"error", err)
			return nil, err
		}

		if img == nil {
			slog.Warn("ConvertPdfToPngWithImageDetails: nil image for page", "page", pageNum)
			continue
		}

		var croppedFrame image.Image
		var cropInfo CropInfo

		croppedFrame, cropInfo = cropToContentWithInfo(img)
		whiteBackgroundFrame := convertToWhiteBackground(croppedFrame)

		outputFilename := prefix + strconv.Itoa(pageNum) + ".png"
		outputFilepath := destpath + outputFilename
		err = saveImageAsPng(whiteBackgroundFrame, outputFilepath)
		if err != nil {
			return nil, err
		}

		slog.Debug("Saved PDF page with crop info",
			"filename", outputFilepath,
			"page", pageNum,
			"offsetX", cropInfo.OffsetX,
			"offsetY", cropInfo.OffsetY,
			"originalSize", fmt.Sprintf("%dx%d", cropInfo.OriginalWidth, cropInfo.OriginalHeight),
			"croppedSize", fmt.Sprintf("%dx%d", cropInfo.CroppedWidth, cropInfo.CroppedHeight))

		var cropDetail *CropDetail
		var imageWidth, imageHeight int

		if cropInfo.CroppedWidth != cropInfo.OriginalWidth || cropInfo.CroppedHeight != cropInfo.OriginalHeight {
			cropDetail = &CropDetail{
				OffsetX:        cropInfo.OffsetX,
				OffsetY:        cropInfo.OffsetY,
				OriginalWidth:  cropInfo.OriginalWidth,
				OriginalHeight: cropInfo.OriginalHeight,
				CroppedWidth:   cropInfo.CroppedWidth,
				CroppedHeight:  cropInfo.CroppedHeight,
			}
			imageWidth = cropInfo.CroppedWidth
			imageHeight = cropInfo.CroppedHeight
		} else {
			cropDetail = nil
			imageWidth = cropInfo.OriginalWidth
			imageHeight = cropInfo.OriginalHeight
		}

		imageDetail := &ImageDetail{
			ActualType: "png",
			Page:       pageNum + 1,
			Pages:      pageCount,
			URL:        filepath.Join(destpath, outputFilename),
			Width:      imageWidth,
			Height:     imageHeight,
			Format:     "png",
			Quality:    95.0,
			CropDetail: cropDetail,
		}

		imageDetails = append(imageDetails, imageDetail)
	}

	return imageDetails, nil
}

// ConvertPdfToPng provides a simplified interface that returns only filenames
func ConvertPdfToPng(pdfFilename string, destpath string, prefix string) (*[]string, error) {
	imageDetails, err := ConvertPdfToPngWithImageDetails(pdfFilename, destpath, prefix)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, len(imageDetails))
	for i, detail := range imageDetails {
		filenames[i] = filepath.Base(detail.URL)
	}

	return &filenames, nil
}

// ConvertPdfToPngWithCropInfo provides backward compatibility
func ConvertPdfToPngWithCropInfo(pdfFilename string, destpath string, prefix string) (*ConversionResult, error) {
	imageDetails, err := ConvertPdfToPngWithImageDetails(pdfFilename, destpath, prefix)
	if err != nil {
		return nil, err
	}

	result := &ConversionResult{
		Filenames: make([]string, len(imageDetails)),
		CropInfos: make([]CropInfo, len(imageDetails)),
	}

	for i, detail := range imageDetails {
		result.Filenames[i] = filepath.Base(detail.URL)
		if detail.CropDetail != nil {
			result.CropInfos[i] = CropInfo{
				OffsetX:        detail.CropDetail.OffsetX,
				OffsetY:        detail.CropDetail.OffsetY,
				OriginalWidth:  detail.CropDetail.OriginalWidth,
				OriginalHeight: detail.CropDetail.OriginalHeight,
				CroppedWidth:   detail.CropDetail.CroppedWidth,
				CroppedHeight:  detail.CropDetail.CroppedHeight,
			}
		}
	}

	return result, nil
}
