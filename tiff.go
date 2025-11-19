package tifpdf2png

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	tiff "github.com/dhushon/tiff"
)

// ConvertTiffToPngWithImageDetails converts TIFF to PNG and returns ImageDetail slice
func ConvertTiffToPngWithImageDetails(tiffFilename string, destpath string, prefix string) ([]*ImageDetail, error) {
	file, err := os.Open(tiffFilename)
	if err != nil {
		slog.Error("ConvertTiffToPngWithImageDetails: load file", "error", err)
		return nil, err
	}
	defer file.Close()

	frames, _, err := tiff.DecodeAll(file)
	if err != nil {
		slog.Error("ConvertTiffToPngWithImageDetails: decode tiff", "error", err)
		return nil, err
	}

	if len(frames) == 0 {
		slog.Error("ConvertTiffToPngWithImageDetails: no images found in TIFF file")
		return nil, fmt.Errorf("no images found in TIFF file")
	}

	if prefix == "" {
		prefix = time.Now().Format("20060102-150405-")
	}
	if destpath[len(destpath)-1:] != "/" {
		destpath = destpath + "/"
	}

	var imageDetails []*ImageDetail

	for i, img := range frames {
		if len(img) == 0 {
			slog.Warn("ConvertTiffToPngWithImageDetails: empty image frame", "frameIndex", i)
			continue
		}
		if img[0] == nil {
			slog.Warn("ConvertTiffToPngWithImageDetails: nil image frame", "frameIndex", i)
			continue
		}

		croppedFrame, cropInfo := cropToContentWithInfo(img[0])
		whiteBackgroundFrame := convertToWhiteBackground(croppedFrame)

		outputFilename := prefix + strconv.Itoa(i) + ".png"
		outputFilepath := destpath + outputFilename
		err = saveImageAsPng(whiteBackgroundFrame, outputFilepath)
		if err != nil {
			return nil, err
		}

		slog.Debug("Saved file with crop info",
			"filename", outputFilepath,
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
			Page:       i + 1,
			Pages:      len(frames),
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

// ConvertTiffToPng provides a simplified interface that returns only filenames
func ConvertTiffToPng(tiffFilename string, destpath string, prefix string) (*[]string, error) {
	imageDetails, err := ConvertTiffToPngWithImageDetails(tiffFilename, destpath, prefix)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, len(imageDetails))
	for i, detail := range imageDetails {
		filenames[i] = filepath.Base(detail.URL)
	}

	return &filenames, nil
}

// ConvertTiffToPngWithCropInfo provides backward compatibility
func ConvertTiffToPngWithCropInfo(tiffFilename string, destpath string, prefix string) (*ConversionResult, error) {
	imageDetails, err := ConvertTiffToPngWithImageDetails(tiffFilename, destpath, prefix)
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
