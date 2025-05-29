package image

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/h2non/bimg"
)

// Processor handles image processing operations
type Processor struct {
	originalImage *bimg.Image
	buffer        []byte
	width         int
	height        int
}

// ProcessOptions contains options for image processing
type ProcessOptions struct {
	Quality int
	Width   int
	Height  int
	Format  bimg.ImageType
}

// NewProcessor creates a new image processor from a file
func NewProcessor(filePath string) (*Processor, error) {
	// Read the image
	buffer, err := bimg.Read(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading image: %w", err)
	}

	// Create image object
	originalImage := bimg.NewImage(buffer)

	// Get image size
	size, err := originalImage.Size()
	if err != nil {
		return nil, fmt.Errorf("error getting image size: %w", err)
	}

	return &Processor{
		originalImage: originalImage,
		buffer:        buffer,
		width:         size.Width,
		height:        size.Height,
	}, nil
}

// GetOriginalInfo returns information about the original image
func (p *Processor) GetOriginalInfo() (int, int, int, string) {
	imageType := bimg.DetermineImageType(p.buffer)
	originalFormat := bimg.ImageTypeName(imageType)
	return p.width, p.height, len(p.buffer), originalFormat
}

// GetOriginalBuffer returns the original image buffer
func (p *Processor) GetOriginalBuffer() []byte {
	return p.buffer
}

// Process compresses and optionally resizes the image
func (p *Processor) Process(opts ProcessOptions) ([]byte, error) {
	// Create options for processing
	options := bimg.Options{
		Quality: opts.Quality,
		Type:    opts.Format,
	}

	// Set width and height if provided
	if opts.Width > 0 {
		options.Width = opts.Width
	}
	if opts.Height > 0 {
		options.Height = opts.Height
	}

	// Process the image
	newImage, err := p.originalImage.Process(options)
	if err != nil {
		return nil, fmt.Errorf("error processing image: %w", err)
	}

	return newImage, nil
}

// GetOutputFilename returns an appropriate filename for the processed image
func GetOutputFilename(inputPath string, compress bool, format bimg.ImageType, useTimestamp bool) string {
	// Determine the appropriate extension
	var extension string
	if !compress {
		// Keep original extension if not compressing
		extension = filepath.Ext(inputPath)
	} else {
		extension = "." + strings.ToLower(bimg.ImageTypeName(format))
	}

	// Use timestamp or original filename as base
	if useTimestamp {
		// Generate timestamp-based filename with format: YYYYMMDDHHMMSSmmm.ext
		timestamp := time.Now().Format("20060102150405.000")
		// replace last dot to empty
		timestamp = strings.ReplaceAll(timestamp, ".", "")
		return timestamp + extension
	}

	baseName := filepath.Base(inputPath)
	// Use original filename
	if !compress {
		return baseName
	}

	// Remove original extension and add new one
	return strings.TrimSuffix(baseName, filepath.Ext(baseName)) + extension
}
