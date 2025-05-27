package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/h2non/bimg"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// initConfig initializes the configuration from config file and environment variables
func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	viper.AddConfigPath(".")
	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(home, ".imgood"))
	}

	viper.SetEnvPrefix("IMGOOD")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	return nil
}

// printUsage prints the usage instructions for the tool
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  imgood [command] [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  up\tUpload an image to S3 with optional compression")
	fmt.Println("\nFor more information about a command, run:")
	fmt.Println("  imgood [command] --help")
}

// printUploadUsage prints the usage instructions for the upload command
func printUploadUsage() {
	fmt.Println("Usage:")
	fmt.Println("  imgood up [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -i, --input string\tPath to the input image file (required)")
	fmt.Println("  -k, --key string\tS3 object key (path in bucket), defaults to filename")
	fmt.Println("  -c, --compress\t\tCompress image before uploading")
	fmt.Println("  -q, --quality int\tQuality of the compressed image (1-100) (default 80)")
	fmt.Println("  -w, --width int\t\tWidth of the output image (0 for original)")
	fmt.Println("  -h, --height int\tHeight of the output image (0 for original)")
	fmt.Println("\nS3 Configuration:")
	fmt.Println("  Configure S3 settings in config.toml or using environment variables:")
	fmt.Println("  - IMGOOD_S3_BUCKET\t\tS3 bucket name")
	fmt.Println("  - IMGOOD_S3_ENDPOINT\t\tS3 endpoint URL (for non-AWS S3 services)")
	fmt.Println("  - IMGOOD_S3_REGION\t\tAWS region")
	fmt.Println("  - IMGOOD_S3_ACCESS_KEY\tAWS access key ID")
	fmt.Println("  - IMGOOD_S3_SECRET_KEY\tAWS secret access key")
}

func main() {
	if err := initConfig(); err != nil {
		fmt.Printf("Warning: %s\n", err)
	}

	// Check if no arguments provided
	if len(os.Args) == 1 {
		printUsage()
		return
	}

	// Handle subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "up":
			// Check if help flag is provided
			if len(os.Args) > 2 && (os.Args[2] == "--help" || os.Args[2] == "-h") {
				printUploadUsage()
				return
			}
			os.Args = append(os.Args[:1], os.Args[2:]...)
			uploadCommand()
			return
		case "--help", "-h":
			printUsage()
			return
		default:
			fmt.Printf("Unknown command: %s\n\n", os.Args[1])
			printUsage()
			os.Exit(1)
		}
	}
}

func uploadCommand() {
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

	// Get S3 configuration from Viper
	bucket := viper.GetString("s3.bucket")
	endpoint := viper.GetString("s3.endpoint")
	region := viper.GetString("s3.region")
	accessKey := viper.GetString("s3.access_key")
	secretKey := viper.GetString("s3.secret_key")

	if bucket == "" {
		fmt.Println("Error: S3 bucket name is required")
		fmt.Println("Set it in config.toml file or IMGOOD_S3_BUCKET environment variable")
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

	// Compress the image if requested
	imageData := buffer
	if *compress {
		// Determine original image format
		imageType := bimg.DetermineImageType(buffer)
		originalFormat := bimg.ImageTypeName(imageType)
		fmt.Printf("Original format: %s\n", originalFormat)

		// Create options for processing
		options := bimg.Options{
			Quality: *quality,
			Type:    bimg.WEBP, // Convert to WebP format for better compression
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

		// Use the compressed image for upload
		imageData = newImage
		fmt.Printf("Compressed image: %d bytes (%.2f%% of original)\n", len(newImage), float64(len(newImage))/float64(len(buffer))*100)
	}

	// Set up S3 client
	cfg, err := configureAWS(region, accessKey, secretKey)
	if err != nil {
		fmt.Printf("Error configuring AWS: %s\n", err)
		os.Exit(1)
	}

	// Create S3 client with options
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	// Set default key if not provided
	if *key == "" {
		baseName := filepath.Base(*inputPath)

		// If we're compressing and converting to WebP, change the extension
		if *compress {
			// Remove original extension and add .webp
			baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName)) + ".webp"
		}

		*key = baseName
	}

	// Upload to S3
	ctx := context.Background()
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(*key),
		Body:   bytes.NewReader(imageData),
	})

	if err != nil {
		fmt.Printf("Error uploading to S3: %s\n", err)
		os.Exit(1)
	}

	// Construct the S3 URL
	var s3URL string
	if endpoint != "" {
		// For custom S3 endpoints
		s3URL = fmt.Sprintf("%s/%s/%s", endpoint, bucket, *key)
	} else {
		// For AWS S3
		s3URL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, *key)
	}

	fmt.Printf("Successfully uploaded to S3: %s\n", s3URL)
}

// configureAWS sets up the AWS configuration with the provided credentials and region
func configureAWS(region, accessKey, secretKey string) (aws.Config, error) {
	configOptions := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	// Add credentials if provided
	if accessKey != "" && secretKey != "" {
		configOptions = append(configOptions, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		))
	}

	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background(), configOptions...)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
