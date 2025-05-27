package s3

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/mingeme/imgood/internal/config"
)

// Client represents an S3 client
type Client struct {
	s3Client *s3.Client
	config   config.S3Config
}

// NewClient creates a new S3 client with the provided configuration
func NewClient(cfg config.S3Config) (*Client, error) {
	// Validate required configuration
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("S3 bucket name is required")
	}

	// Configure AWS
	awsCfg, err := configureAWS(cfg.Region, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error configuring AWS: %w", err)
	}

	// Create S3 client with options
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	return &Client{
		s3Client: s3Client,
		config:   cfg,
	}, nil
}

// UploadFile uploads a file to S3
func (c *Client) UploadFile(key string, data []byte) error {
	ctx := context.Background()
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return fmt.Errorf("error uploading to S3: %w", err)
	}

	return nil
}

// GetObject downloads an object from S3
func (c *Client) GetObject(key string) ([]byte, error) {
	ctx := context.Background()
	result, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("error getting object from S3: %w", err)
	}

	defer result.Body.Close()

	// Read the object body
	data := bytes.Buffer{}
	_, err = data.ReadFrom(result.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading object body: %w", err)
	}

	return data.Bytes(), nil
}

// ObjectExists checks if an object exists in S3
func (c *Client) ObjectExists(key string) (bool, error) {
	ctx := context.Background()
	_, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		// Check if the error is because the object doesn't exist
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "NoSuchKey") {
			return false, nil
		}
		return false, fmt.Errorf("error checking if object exists: %w", err)
	}

	return true, nil
}

// GetFileURL returns the URL for an uploaded file
func (c *Client) GetFileURL(key string) string {
	if c.config.Endpoint != "" {
		// For custom S3 endpoints
		// Format: https://imgood.s3.example.com/path/to/file.jpeg
		return fmt.Sprintf("https://%s.%s/%s", c.config.Bucket, strings.TrimPrefix(c.config.Endpoint, "https://"), key)
	}
	// For AWS S3
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.config.Bucket, c.config.Region, key)
}

// S3Object represents an object in S3
type S3Object struct {
	Key          string
	Size         int64
	LastModified time.Time
	URL          string
}

// ListObjects lists objects in the S3 bucket with an optional prefix
func (c *Client) ListObjects(prefix string, maxKeys int32) ([]S3Object, error) {
	ctx := context.Background()

	// Create the input for listing objects
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.config.Bucket),
	}

	// Add prefix if provided
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	// Set max keys if provided
	if maxKeys > 0 {
		input.MaxKeys = aws.Int32(maxKeys)
	}

	// List objects
	result, err := c.s3Client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error listing objects in S3: %w", err)
	}

	// Convert to S3Object slice
	objects := make([]S3Object, 0, len(result.Contents))
	for _, item := range result.Contents {
		objects = append(objects, S3Object{
			Key:          *item.Key,
			Size:         *item.Size,
			LastModified: *item.LastModified,
			URL:          c.GetFileURL(*item.Key),
		})
	}

	return objects, nil
}

// configureAWS sets up the AWS configuration with the provided credentials and region
func configureAWS(region, accessKey, secretKey string) (aws.Config, error) {
	configOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(region),
	}

	// Add credentials if provided
	if accessKey != "" && secretKey != "" {
		configOptions = append(configOptions, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		))
	}

	// Load the AWS configuration
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(), configOptions...)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
