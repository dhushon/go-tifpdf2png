package tifpdf2png

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log/slog"
	"os"

	"github.com/disintegration/imaging"
)

// cropToContentWithInfo crops an image to its content boundaries and returns crop information
func cropToContentWithInfo(img image.Image) (image.Image, CropInfo) {
	bounds := img.Bounds()
	minX, minY, maxX, maxY := bounds.Max.X, bounds.Max.Y, bounds.Min.X, bounds.Min.Y

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	cropInfo := CropInfo{
		OffsetX:        minX,
		OffsetY:        minY,
		OriginalWidth:  bounds.Dx(),
		OriginalHeight: bounds.Dy(),
		CroppedWidth:   maxX + 1 - minX,
		CroppedHeight:  maxY + 1 - minY,
	}

	croppedImg := imaging.Crop(img, image.Rect(minX, minY, maxX+1, maxY+1))
	return croppedImg, cropInfo
}

// convertToWhiteBackground ensures the image has a white background with black content
func convertToWhiteBackground(src image.Image) image.Image {
	bounds := src.Bounds()

	darkPixelCount := 0
	lightPixelCount := 0
	sampledPixels := 0

	sampleRegions := []image.Rectangle{
		{bounds.Min, image.Point{bounds.Min.X + 50, bounds.Min.Y + 50}},
		{image.Point{bounds.Max.X - 50, bounds.Min.Y}, image.Point{bounds.Max.X, bounds.Min.Y + 50}},
		{image.Point{bounds.Min.X, bounds.Max.Y - 50}, image.Point{bounds.Min.X + 50, bounds.Max.Y}},
		{image.Point{bounds.Max.X - 50, bounds.Max.Y - 50}, bounds.Max},
	}

	for _, region := range sampleRegions {
		for y := region.Min.Y; y < region.Max.Y && y < bounds.Max.Y; y++ {
			for x := region.Min.X; x < region.Max.X && x < bounds.Max.X; x++ {
				if x < bounds.Min.X || y < bounds.Min.Y {
					continue
				}

				r, g, b, a := src.At(x, y).RGBA()
				if uint8(a>>8) == 0 {
					continue
				}

				r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
				luminance := 0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8)

				sampledPixels++
				if luminance < 128 {
					darkPixelCount++
				} else {
					lightPixelCount++
				}
			}
		}
	}

	if sampledPixels < 100 {
		step := 20
		for y := bounds.Min.Y; y < bounds.Max.Y; y += step {
			for x := bounds.Min.X; x < bounds.Max.X; x += step {
				r, g, b, a := src.At(x, y).RGBA()
				if uint8(a>>8) == 0 {
					continue
				}

				r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
				luminance := 0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8)

				sampledPixels++
				if luminance < 128 {
					darkPixelCount++
				} else {
					lightPixelCount++
				}
			}
		}
	}

	shouldInvert := float64(darkPixelCount)/float64(sampledPixels) > 0.5

	slog.Debug("Background analysis",
		"darkPixels", darkPixelCount,
		"lightPixels", lightPixelCount,
		"sampledPixels", sampledPixels,
		"darkRatio", float64(darkPixelCount)/float64(sampledPixels),
		"shouldInvert", shouldInvert)

	dst := image.NewRGBA(bounds)
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	draw.Draw(dst, bounds, &image.Uniform{white}, image.Point{}, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)

			if a8 == 0 {
				dst.Set(x, y, white)
				continue
			}

			luminance := 0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8)

			if shouldInvert {
				if luminance < 128 {
					dst.Set(x, y, white)
				} else {
					dst.Set(x, y, black)
				}
			} else {
				if luminance < 128 {
					dst.Set(x, y, black)
				} else {
					dst.Set(x, y, white)
				}
			}
		}
	}

	return dst
}

// saveImageAsPng saves an image to a PNG file
func saveImageAsPng(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
