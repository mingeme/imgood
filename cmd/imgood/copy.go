package imgood

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/bimg"
	"github.com/spf13/pflag"

	"github.com/mingeme/imgood/internal/config"
	"github.com/mingeme/imgood/internal/s3"
)

// ExecuteCopy runs the copy command
func ExecuteCopy() {
	// Define command line flags for copy operation
	sourceKey := pflag.StringP("source", "s", "", "Source S3 object key to copy")
	targetKey := pflag.StringP("target", "t", "", "Target S3 object key (destination)")
	convertFormat := pflag.StringP("format", "f", "webp", "Convert to format (webp, jpeg, png)")
	quality := pflag.IntP("quality", "q", 80, "Quality of the converted image (1-100)")
	width := pflag.IntP("width", "w", 0, "Width of the output image (0 for original)")
	height := pflag.IntP("height", "h", 0, "Height of the output image (0 for original)")
	pflag.Parse()

	// Validate required parameters
	if *sourceKey == "" {
		fmt.Println("Error: Source key is required")
		pflag.Usage()
		os.Exit(1)
	}

	// Get S3 configuration and create client
	s3Config := config.GetS3Config()
	s3Client, err := s3.NewClient(s3Config)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Println("Check your S3 configuration in config.toml or environment variables")
		os.Exit(1)
	}

	// Check if source object exists
	exists, err := s3Client.ObjectExists(*sourceKey)
	if err != nil {
		fmt.Printf("Error checking source object: %s\n", err)
		os.Exit(1)
	}
	if !exists {
		fmt.Printf("Error: Source object does not exist: %s\n", *sourceKey)
		os.Exit(1)
	}

	// Set default target key if not provided
	if *targetKey == "" {
		ext := filepath.Ext(*sourceKey)
		baseName := strings.TrimSuffix(*sourceKey, ext)

		// If format conversion is requested, change the extension
		if *convertFormat != "" {
			*targetKey = baseName + "." + *convertFormat
		} else {
			*targetKey = baseName + ext
		}
	}

	// Check if source and target are the same
	if *sourceKey == *targetKey {
		fmt.Println("Error: Source and target keys cannot be the same")
		os.Exit(1)
	}

	// Check if target already exists
	exists, err = s3Client.ObjectExists(*targetKey)
	if err != nil {
		fmt.Printf("Error checking target object: %s\n", err)
		os.Exit(1)
	}
	if exists {
		fmt.Printf("Error: Target object already exists: %s\n", *targetKey)
		os.Exit(1)
	}

	// Download the source object
	fmt.Printf("Downloading object: %s\n", *sourceKey)
	imageData, err := s3Client.GetObject(*sourceKey)
	if err != nil {
		fmt.Printf("Error downloading source object: %s\n", err)
		os.Exit(1)
	}

	// Get original image info
	originalImage := bimg.NewImage(imageData)
	size, err := originalImage.Size()
	if err != nil {
		fmt.Printf("Error getting image size: %s\n", err)
		os.Exit(1)
	}
	imageType := bimg.DetermineImageType(imageData)
	originalFormat := bimg.ImageTypeName(imageType)
	fmt.Printf("Original image: %dx%d, %d bytes, format: %s\n",
		size.Width, size.Height, len(imageData), originalFormat)

	// Process the image if format conversion or resizing is requested
	var outputData []byte
	if *convertFormat != "" || *width > 0 || *height > 0 {
		// Determine target format
		var targetFormat bimg.ImageType
		switch strings.ToLower(*convertFormat) {
		case "webp":
			targetFormat = bimg.WEBP
		case "jpeg", "jpg":
			targetFormat = bimg.JPEG
		case "png":
			targetFormat = bimg.PNG
		case "":
			// If no format specified but resizing is requested, keep original format
			targetFormat = imageType
		default:
			fmt.Printf("Unsupported format: %s. Using original format.\n", *convertFormat)
			targetFormat = imageType
		}

		// Create options for processing
		options := bimg.Options{
			Quality: *quality,
			Type:    targetFormat,
		}

		// Set width and height if provided
		if *width > 0 {
			options.Width = *width
		}
		if *height > 0 {
			options.Height = *height
		}

		// Process the image
		outputData, err = originalImage.Process(options)
		if err != nil {
			fmt.Printf("Error processing image: %s\n", err)
			os.Exit(1)
		}

		newFormat := bimg.ImageTypeName(targetFormat)
		fmt.Printf("Converted image: %d bytes, format: %s (%.2f%% of original)\n",
			len(outputData), newFormat, float64(len(outputData))/float64(len(imageData))*100)
	} else {
		// No conversion needed, use original data
		outputData = imageData
		fmt.Println("No conversion requested, copying original image")
	}

	// Upload to target key
	fmt.Printf("Uploading to: %s\n", *targetKey)
	err = s3Client.UploadFile(*targetKey, outputData)
	if err != nil {
		fmt.Printf("Error uploading to S3: %s\n", err)
		os.Exit(1)
	}

	// Get and display the file URL
	s3URL := s3Client.GetFileURL(*targetKey)
	fmt.Printf("Successfully copied to: %s\n", s3URL)
}
