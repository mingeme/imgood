package imgood

import (
	"fmt"
	"os"

	"github.com/h2non/bimg"
	"github.com/spf13/pflag"
	
	"github.com/mingeme/imgood/internal/config"
	"github.com/mingeme/imgood/internal/image"
	"github.com/mingeme/imgood/internal/s3"
)

// ExecuteUpload runs the upload command
func ExecuteUpload() {
	// Define command line flags for image processing
	inputPath := pflag.StringP("input", "i", "", "Path to the input image file")
	key := pflag.StringP("key", "k", "", "S3 object key (path in bucket)")
	compress := pflag.BoolP("compress", "c", false, "Compress image before uploading")
	quality := pflag.IntP("quality", "q", 80, "Quality of the compressed image (1-100)")
	width := pflag.IntP("width", "w", 0, "Width of the output image (0 for original)")
	height := pflag.IntP("height", "h", 0, "Height of the output image (0 for original)")
	pflag.Parse()

	// Validate required parameters
	if *inputPath == "" {
		fmt.Println("Error: Input path is required")
		pflag.Usage()
		os.Exit(1)
	}

	// Get S3 configuration
	s3Config := config.GetS3Config()

	// Check if input file exists
	if _, err := os.Stat(*inputPath); os.IsNotExist(err) {
		fmt.Printf("Error: Input file does not exist: %s\n", *inputPath)
		os.Exit(1)
	}

	// Create image processor
	processor, err := image.NewProcessor(*inputPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get original image info
	width0, height0, size, format := processor.GetOriginalInfo()
	fmt.Printf("Original image: %dx%d, %d bytes, format: %s\n", width0, height0, size, format)

	// Process the image if compression is requested
	var imageData []byte
	if *compress {
		// Process the image
		processOpts := image.ProcessOptions{
			Quality: *quality,
			Width:   *width,
			Height:  *height,
			Format:  bimg.WEBP, // Convert to WebP format for better compression
		}
		
		newImage, err := processor.Process(processOpts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		
		imageData = newImage
		fmt.Printf("Compressed image: %d bytes (%.2f%% of original)\n", 
			len(newImage), float64(len(newImage))/float64(size)*100)
	} else {
		// Use original image
		imageData = processor.GetOriginalBuffer()
	}

	// Create S3 client
	s3Client, err := s3.NewClient(s3Config)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Println("Check your S3 configuration in config.toml or environment variables")
		os.Exit(1)
	}

	// Set default key if not provided
	if *key == "" {
		*key = image.GetOutputFilename(*inputPath, *compress, bimg.WEBP)
	}

	// Upload to S3
	err = s3Client.UploadFile(*key, imageData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get and display the file URL
	s3URL := s3Client.GetFileURL(*key)
	fmt.Printf("Successfully uploaded to S3: %s\n", s3URL)
}
