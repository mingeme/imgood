package main

import (
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

func main() {
	if err := initConfig(); err != nil {
		fmt.Printf("Warning: %s\n", err)
	}
	if len(os.Args) > 1 && os.Args[1] == "up" {
		os.Args = append(os.Args[:1], os.Args[2:]...)
		uploadCommand()
		return
	}

	compressCommand()
}

func compressCommand() {
	// Define command line flags for compression
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

func uploadCommand() {
	// Define command line flags for S3 upload
	inputPath := pflag.StringP("input", "i", "", "Path to the input image file")

	// Use config values as defaults for flags
	bucket := pflag.StringP("bucket", "b", viper.GetString("s3.bucket"), "S3 bucket name")
	key := pflag.StringP("key", "k", "", "S3 object key (path in bucket)")
	endpoint := pflag.StringP("endpoint", "e", viper.GetString("s3.endpoint"), "S3 endpoint URL (for non-AWS S3 services)")
	region := pflag.StringP("region", "r", viper.GetString("s3.region"), "AWS region")
	accessKey := pflag.StringP("access-key", "a", viper.GetString("s3.access_key"), "AWS access key ID")
	secretKey := pflag.StringP("secret-key", "s", viper.GetString("s3.secret_key"), "AWS secret access key")

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

	if *bucket == "" {
		fmt.Println("Error: S3 bucket name is required")
		fmt.Println("Set it using --bucket flag, config.yaml file, or IMGOOD_S3_BUCKET environment variable")
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

	// Compress the image if requested
	imageData := buffer
	if *compress {
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

		// Use the compressed image for upload
		imageData = newImage
		fmt.Printf("Compressed image: %d bytes (%.2f%% of original)\n", len(newImage), float64(len(newImage))/float64(len(buffer))*100)
	}

	// Set up S3 client
	cfg, err := configureAWS(*region, *accessKey, *secretKey, *endpoint)
	if err != nil {
		fmt.Printf("Error configuring AWS: %s\n", err)
		os.Exit(1)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Set default key if not provided
	if *key == "" {
		*key = filepath.Base(*inputPath)
	}

	// Upload to S3
	ctx := context.Background()
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(*bucket),
		Key:    aws.String(*key),
		Body:   strings.NewReader(string(imageData)),
	})

	if err != nil {
		fmt.Printf("Error uploading to S3: %s\n", err)
		os.Exit(1)
	}

	// Construct the S3 URL
	var s3URL string
	if *endpoint != "" {
		// For custom S3 endpoints
		s3URL = fmt.Sprintf("%s/%s/%s", *endpoint, *bucket, *key)
	} else {
		// For AWS S3
		s3URL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", *bucket, *region, *key)
	}

	fmt.Printf("Successfully uploaded to S3: %s\n", s3URL)
}

// configureAWS sets up the AWS configuration with the provided credentials and region
func configureAWS(region, accessKey, secretKey, endpoint string) (aws.Config, error) {
	configOptions := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	// Add credentials if provided
	if accessKey != "" && secretKey != "" {
		configOptions = append(configOptions, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		))
	}

	// Add custom endpoint if provided
	if endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					URL:               endpoint,
					HostnameImmutable: true,
					SigningRegion:     region,
				}, nil
			}
			// Fallback to default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})
		configOptions = append(configOptions, config.WithEndpointResolverWithOptions(customResolver))
	}

	return config.LoadDefaultConfig(context.Background(), configOptions...)
}
