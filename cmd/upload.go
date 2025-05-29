package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/h2non/bimg"
	"github.com/spf13/cobra"

	"github.com/mingeme/imgood/internal/config"
	"github.com/mingeme/imgood/internal/image"
	"github.com/mingeme/imgood/internal/s3"
)

var (
	uploadInputPath  string
	uploadKey        string
	uploadCompress   bool
	uploadQuality    int
	uploadResize     string
	uploadTimestamp  bool
	uploadKeepMetadata bool
	uploadNoRotate    bool
)

var uploadCmd = &cobra.Command{
	Use:     "up",
	Aliases: []string{"upload"},
	Short:   "Upload an image to S3 with optional compression",
	Long: `Upload an image to S3 with optional compression and resizing.
	
Example:
  imgood up -i image.jpg -c -q 80 -r 800,600`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate required parameters
		if uploadInputPath == "" {
			fmt.Println("Error: Input path is required")
			cmd.Help()
			os.Exit(1)
		}

		// Get S3 configuration
		s3Config := config.GetS3Config()

		// Check if input file exists
		if _, err := os.Stat(uploadInputPath); os.IsNotExist(err) {
			fmt.Printf("Error: Input file does not exist: %s\n", uploadInputPath)
			os.Exit(1)
		}

		// Create image processor
		processor, err := image.NewProcessor(uploadInputPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Get original image info
		width0, height0, size, format := processor.GetOriginalInfo()
		fmt.Printf("Original image: %dx%d, %d bytes, format: %s\n", width0, height0, size, format)

		// Process the image if compression is requested or if we need to handle EXIF orientation/metadata
		var imageData []byte
		if uploadCompress || !uploadKeepMetadata || !uploadNoRotate {
			// Process the image
			processOpts := image.ProcessOptions{
				Quality:      uploadQuality,
				Width:        0,
				Height:       0,
				Format:       bimg.WEBP, // Convert to WebP format for better compression
				KeepMetadata: uploadKeepMetadata,
				NoRotate:     uploadNoRotate,
			}
			
			// If not compressing but still processing for orientation/metadata, keep original format
			if !uploadCompress {
				imageType := bimg.DetermineImageType(processor.GetOriginalBuffer())
				processOpts.Format = imageType
			}
			
			// Parse resize parameter if provided
			if uploadResize != "" {
				dimensions := strings.Split(uploadResize, ",")
				if len(dimensions) == 2 {
					width, err := strconv.Atoi(strings.TrimSpace(dimensions[0]))
					if err == nil && width > 0 {
						processOpts.Width = width
					}
					
					height, err := strconv.Atoi(strings.TrimSpace(dimensions[1]))
					if err == nil && height > 0 {
						processOpts.Height = height
					}
				}
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
		if uploadKey == "" {
			uploadKey = image.GetOutputFilename(uploadInputPath, uploadCompress, bimg.WEBP, uploadTimestamp)
		}

		// Upload to S3
		err = s3Client.UploadFile(uploadKey, imageData)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Get and display the file URL
		s3URL := s3Client.GetFileURL(uploadKey)
		fmt.Printf("Successfully uploaded to S3: %s\n", s3URL)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Define command line flags for image processing
	uploadCmd.Flags().StringVarP(&uploadInputPath, "input", "i", "", "Path to the input image file (required)")
	uploadCmd.Flags().StringVarP(&uploadKey, "key", "k", "", "S3 object key (path in bucket)")
	uploadCmd.Flags().BoolVarP(&uploadCompress, "compress", "c", false, "Compress image before uploading")
	uploadCmd.Flags().IntVarP(&uploadQuality, "quality", "q", 80, "Quality of the compressed image (1-100)")
	uploadCmd.Flags().StringVarP(&uploadResize, "resize", "r", "", "Resize image to width,height (e.g., '800,600'). Use 0 for any dimension to maintain aspect ratio")
	uploadCmd.Flags().BoolVarP(&uploadTimestamp, "timestamp", "t", false, "Use timestamp as filename when key is not specified")
	uploadCmd.Flags().BoolVar(&uploadKeepMetadata, "keep-metadata", false, "Keep image metadata (EXIF, etc.)")
	uploadCmd.Flags().BoolVar(&uploadNoRotate, "no-rotate", false, "Disable automatic rotation based on EXIF orientation")

	// Mark required flags
	uploadCmd.MarkFlagRequired("input")

	// Add shell completion for flags
	_ = uploadCmd.RegisterFlagCompletionFunc("input", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveDefault
	})
}
