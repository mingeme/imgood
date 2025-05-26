package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/bimg"
	"github.com/spf13/pflag"
)

func main() {
	// Define command line flags
	inputPath := pflag.StringP("input", "i", "", "Path to the input image file")
	quality := pflag.IntP("quality", "q", 80, "Quality of the compressed image (1-100)")
	width := pflag.IntP("width", "w", 0, "Width of the output image (0 for original)")
	height := pflag.IntP("height", "h", 0, "Height of the output image (0 for original)")
	pflag.Parse()

	// Validate input path
	if *inputPath == "" {
		fmt.Println("Error: Input path is required")
		pflag.Usage()
		os.Exit(1)
	}

	// Check if input file exists
	if _, err := os.Stat(*inputPath); os.IsNotExist(err) {
		fmt.Printf("Error: Input file does not exist: %s\n", *inputPath)
		os.Exit(1)
	}

	// Read the image
	buffer, err := bimg.Read(*inputPath)
	if err != nil {
		fmt.Printf("Error reading image: %s\n", err)
		os.Exit(1)
	}

	// Get original image info
	originalImage := bimg.NewImage(buffer)
	size, err := originalImage.Size()
	if err != nil {
		fmt.Printf("Error getting image size: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Original image: %dx%d, %d bytes\n", size.Width, size.Height, len(buffer))

	// Create options for processing
	options := bimg.Options{
		Quality: *quality,
	}

	// Set width and height if provided
	if *width > 0 {
		options.Width = *width
	}
	if *height > 0 {
		options.Height = *height
	}

	// Process the image
	newImage, err := originalImage.Process(options)
	if err != nil {
		fmt.Printf("Error processing image: %s\n", err)
		os.Exit(1)
	}

	// Use system's /tmp directory for temporary files
	tmpDir := os.TempDir()

	// Generate output filename
	fileName := filepath.Base(*inputPath)
	fileExt := filepath.Ext(fileName)
	fileNameWithoutExt := strings.TrimSuffix(fileName, fileExt)
	outputPath := filepath.Join(tmpDir, fmt.Sprintf("%s_compressed%s", fileNameWithoutExt, fileExt))

	// Save the processed image
	err = bimg.Write(outputPath, newImage)
	if err != nil {
		fmt.Printf("Error saving image: %s\n", err)
		os.Exit(1)
	}

	// Get compressed image info
	compressedSize := len(newImage)
	compressionRatio := float64(compressedSize) / float64(len(buffer)) * 100

	fmt.Printf("Compressed image: %s, %d bytes (%.2f%% of original)\n", outputPath, compressedSize, compressionRatio)
}
