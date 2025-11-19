package tifpdf2png

// ImageDetail contains detailed information about a converted image page
type ImageDetail struct {
	ActualType string      `json:"actual_type"`           // The actual type of the image (e.g., "png")
	Page       int         `json:"page"`                  // Page number (1-based)
	Pages      int         `json:"pages"`                 // Total number of pages
	URL        string      `json:"url"`                   // Path to the output PNG file
	Width      int         `json:"width"`                 // Width of the image in pixels
	Height     int         `json:"height"`                // Height of the image in pixels
	Format     string      `json:"format"`                // Image format (e.g., "png")
	Quality    float64     `json:"quality"`               // Quality metric (0-100)
	CropDetail *CropDetail `json:"crop_detail,omitempty"` // Crop information if cropping occurred
}

// CropDetail contains information about how an image was cropped
type CropDetail struct {
	OffsetX        int `json:"offset_x"`        // X offset of the crop from original
	OffsetY        int `json:"offset_y"`        // Y offset of the crop from original
	OriginalWidth  int `json:"original_width"`  // Original width before cropping
	OriginalHeight int `json:"original_height"` // Original height before cropping
	CroppedWidth   int `json:"cropped_width"`   // Width after cropping
	CroppedHeight  int `json:"cropped_height"`  // Height after cropping
}

// CropInfo is an internal structure used during image processing
type CropInfo struct {
	OffsetX        int `json:"offset_x"`
	OffsetY        int `json:"offset_y"`
	OriginalWidth  int `json:"original_width"`
	OriginalHeight int `json:"original_height"`
	CroppedWidth   int `json:"cropped_width"`
	CroppedHeight  int `json:"cropped_height"`
}

// ConversionResult contains the results of document to PNG conversion (legacy structure)
type ConversionResult struct {
	Filenames []string   `json:"filenames"`  // List of generated PNG filenames
	CropInfos []CropInfo `json:"crop_infos"` // Crop information for each page
}
