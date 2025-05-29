package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/h2non/bimg"
	"github.com/spf13/cobra"

	"github.com/mingeme/imgood/internal/config"
	"github.com/mingeme/imgood/internal/s3"
)

var (
	copySourceKey     string
	copyTargetKey     string
	copyConvertFormat string
	copyQuality       int
	copyResize        string
)

var copyCmd = &cobra.Command{
	Use:     "cp",
	Aliases: []string{"copy"},
	Short:   "Copy an object in S3 with optional format conversion",
	Long: `Copy an object in S3 with optional format conversion and resizing.
	
Example:
  imgood cp -s source.jpg -t target.webp -f webp -q 80 -r 800,600`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate required parameters
		if copySourceKey == "" {
			fmt.Println("Error: Source key is required")
			cmd.Help()
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
		exists, err := s3Client.ObjectExists(copySourceKey)
		if err != nil {
			fmt.Printf("Error checking source object: %s\n", err)
			os.Exit(1)
		}
		if !exists {
			fmt.Printf("Error: Source object does not exist: %s\n", copySourceKey)
			os.Exit(1)
		}

		// Set default target key if not provided
		if copyTargetKey == "" {
			ext := filepath.Ext(copySourceKey)
			baseName := strings.TrimSuffix(copySourceKey, ext)

			// If format conversion is requested, change the extension
			if copyConvertFormat != "" {
				copyTargetKey = baseName + "." + copyConvertFormat
			} else {
				copyTargetKey = baseName + "-copy" + ext
			}
		}

		// Check if source and target are the same
		if copySourceKey == copyTargetKey {
			fmt.Println("Error: Source and target keys cannot be the same")
			os.Exit(1)
		}

		// Check if target already exists
		exists, err = s3Client.ObjectExists(copyTargetKey)
		if err != nil {
			fmt.Printf("Error checking target object: %s\n", err)
			os.Exit(1)
		}
		if exists {
			fmt.Printf("Error: Target object already exists: %s\n", copyTargetKey)
			os.Exit(1)
		}

		// Download the source object
		fmt.Printf("Downloading object: %s\n", copySourceKey)
		imageData, err := s3Client.GetObject(copySourceKey)
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
		if copyConvertFormat != "" || copyResize != "" {
			// Determine target format
			var targetFormat bimg.ImageType
			switch strings.ToLower(copyConvertFormat) {
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
				fmt.Printf("Unsupported format: %s. Using original format.\n", copyConvertFormat)
				targetFormat = imageType
			}

			// Create options for processing
			options := bimg.Options{
				Quality: copyQuality,
				Type:    targetFormat,
			}

			// Set width and height if provided
			if copyResize != "" {
				dimensions := strings.Split(copyResize, ",")
				if len(dimensions) == 2 {
					width, err := strconv.Atoi(strings.TrimSpace(dimensions[0]))
					if err == nil && width > 0 {
						options.Width = width
					}
					
					height, err := strconv.Atoi(strings.TrimSpace(dimensions[1]))
					if err == nil && height > 0 {
						options.Height = height
					}
				}
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
		fmt.Printf("Uploading to: %s\n", copyTargetKey)
		err = s3Client.UploadFile(copyTargetKey, outputData)
		if err != nil {
			fmt.Printf("Error uploading to S3: %s\n", err)
			os.Exit(1)
		}

		// Get and display the file URL
		s3URL := s3Client.GetFileURL(copyTargetKey)
		fmt.Printf("Successfully copied to: %s\n", s3URL)
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// Define command line flags for copy operation
	copyCmd.Flags().StringVarP(&copySourceKey, "source", "s", "", "Source S3 object key to copy (required)")
	copyCmd.Flags().StringVarP(&copyTargetKey, "target", "t", "", "Target S3 object key (destination)")
	copyCmd.Flags().StringVarP(&copyConvertFormat, "format", "f", "", "Convert to format (webp, jpeg, png)")
	copyCmd.Flags().IntVarP(&copyQuality, "quality", "q", 80, "Quality of the converted image (1-100)")
	copyCmd.Flags().StringVarP(&copyResize, "resize", "r", "", "Resize image to width,height (e.g., '800,600'). Use 0 for any dimension to maintain aspect ratio")

	// Mark required flags
	copyCmd.MarkFlagRequired("source")

	// Add shell completion for flags
	_ = copyCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"webp", "jpeg", "jpg", "png"}, cobra.ShellCompDirectiveNoFileComp
	})
}
